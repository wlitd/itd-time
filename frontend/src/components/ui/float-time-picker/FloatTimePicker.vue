<script setup lang="ts">
import type { FormValidationStatus, TimePickerInst, TimePickerSize } from 'naive-ui'
import type { FormItemInjection } from 'naive-ui/es/_mixins/use-form-item'
import type { FloatTimePickerProps } from './types'
import { useThemeVars } from 'naive-ui'
import { formItemInjectionKey } from 'naive-ui/es/_mixins/use-form-item'

defineOptions({
  inheritAttrs: false
})

const { placeholder, size = 'medium', disabled = false, status, defaultValue } = defineProps<FloatTimePickerProps>()

const emit = defineEmits<{
  focus: []
  blur: []
}>()

defineSlots<{
  icon?: () => any
}>()

const slots = useSlots()

const value = defineModel<number | null>('value')

// 在 NFormItem 中使用时可以注入表单项上下文
const formItemContext = inject<FormItemInjection | null>(formItemInjectionKey, null)
const _size = computed<TimePickerSize>(() => formItemContext?.mergedSize.value || size)
const _disabled = computed<boolean>(() => formItemContext?.disabled.value || disabled)
const _status = computed<FormValidationStatus | undefined>(() => formItemContext?.mergedValidationStatus.value || status)

// 使用适配 Naive UI 的主题变量
const themeVars = useThemeVars()
const { fontSizeMap, fontSize, paddingX, commonFont } = useNaiveSize(_size.value)

// label 浮动后减小字体
const smallerFontSize = computed<string>(() => _size.value === 'large' ? fontSizeMap.value.medium : fontSizeMap.value.tiny)

const statusTextColor = computed<string>(() => {
  if (_status.value === 'warning') {
    return themeVars.value.warningColor
  } else if (_status.value === 'error') {
    return themeVars.value.errorColor
  } else {
    return ''
  }
})

// CSS 变量
const cssVars = computed<Record<string, string>>(() => ({
  '--float-time-picker-label-font-size': fontSize.value,
  '--float-time-picker-label-float-font-size': smallerFontSize.value,
  '--float-time-picker-label-color': statusTextColor.value || themeVars.value.placeholderColor,
  '--float-time-picker-label-float-color': statusTextColor.value || themeVars.value.textColor1,
  '--float-time-picker-label-bg': themeVars.value.tagColor,
  '--float-time-picker-label-padding-x': '2px',
  '--float-time-picker-label-float-top': '-7px',
  '--float-time-picker-label-left': `${paddingX.value}px`,
  '--float-time-picker-label-initial-left': `${paddingX.value + 4}px`
}))

// 标记是否已经测量完成，用于控制 label 的过渡效果
const isMeasured = ref<boolean>(false)
// TimePicker 没有 prefix slot，直接添加过渡效果
onMounted(() => {
  isMeasured.value = true
})

const isFocused = ref<boolean>(false)
const shouldFloat = computed<boolean>(() => isFocused.value || !!defaultValue || !!value.value)

function handleFocus(): void {
  isFocused.value = true
  emit('focus')
}

function handleBlur(): void {
  isFocused.value = false
  emit('blur')
}

const nDatePickerRef = ref<TimePickerInst | null>(null)

defineExpose({
  focus: () => nDatePickerRef.value?.focus?.(),
  blur: () => nDatePickerRef.value?.blur?.()
})
</script>

<template>
  <div class="relative inline-flex flex-1" :style="cssVars">
    <NTimePicker
      ref="nDatePickerRef" v-model:value="value" placeholder=""
      v-bind="{ defaultValue, size: _size, disabled: _disabled, status: _status, ...$attrs }" @focus="handleFocus"
      @blur="handleBlur"
    >
      <!-- NTimePicker 判断显示 icon 的逻辑和 DatePicker 不一致，这里需要自行判断是否传递 icon 插槽 -->
      <template v-if="slots.icon" #icon>
        <slot name="icon" />
      </template>
    </NTimePicker>
    <label
      :style="commonFont"
      class="float-time-picker__label absolute pointer-events-none z-2 float-time-picker__label--input"
      :class="[{ 'float-time-picker__label--float': shouldFloat, 'float-time-picker__label--ready': isMeasured }]"
    >
      {{ placeholder }}
    </label>
  </div>
</template>

<style lang="css" scoped>
.float-time-picker__label {
  font-size: var(--float-time-picker-label-font-size);
  color: var(--float-time-picker-label-color);
  left: var(--float-time-picker-label-initial-left);
  /* 防止 label 过长被 suffix 和按钮遮挡 */
  max-width: calc(100% - var(--float-time-picker-label-left) - 44px);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.float-time-picker__label--input {
  top: 50%;
  transform: translateY(-50%);
}

.float-time-picker__label--float {
  font-size: var(--float-time-picker-label-float-font-size);
  color: var(--float-time-picker-label-float-color);
  top: var(--float-time-picker-label-float-top);
  padding: 0 var(--float-time-picker-label-padding-x);
  line-height: 1 !important;
  background-color: var(--float-time-picker-label-bg);
  transform: translateY(0);
}

/* 测量完成后才为 label 添加过渡效果，避免初始位置闪烁 */
.float-time-picker__label--ready {
  transition: all 0.2s ease;
}
</style>
