<script setup lang="ts">
import type { ButtonSize, FormValidationStatus, SelectInst } from 'naive-ui'
import type { FormItemInjection } from 'naive-ui/es/_mixins/use-form-item'
import type { Value } from 'naive-ui/es/select/src/interface'
import type { FloatSelectProps } from './types'
import { useThemeVars } from 'naive-ui'
import { formItemInjectionKey } from 'naive-ui/es/_mixins/use-form-item'

defineOptions({
  inheritAttrs: false
})

const {
  placeholder,
  size = 'medium',
  defaultValue,
  status,
  disabled = false
} = defineProps<FloatSelectProps>()

const emit = defineEmits<{
  focus: []
  focusInput: []
  blur: []
  blurInput: []
}>()

defineSlots<{
  header?: () => any
  action?: () => any
  empty?: () => any
  arrow?: () => any
}>()

const value = defineModel<Value | null>('value')

// 在 NFormItem 中使用时可以注入表单项上下文
const formItemContext = inject<FormItemInjection | null>(formItemInjectionKey, null)
const _size = computed<ButtonSize>(() => formItemContext?.mergedSize.value || size)
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
  '--float-select-label-font-size': fontSize.value,
  '--float-select-label-float-font-size': smallerFontSize.value,
  '--float-select-label-color': statusTextColor.value || themeVars.value.placeholderColor,
  '--float-select-label-float-color': statusTextColor.value || themeVars.value.textColor1,
  '--float-select-label-bg': themeVars.value.tagColor,
  '--float-select-label-padding-x': '2px',
  '--float-select-label-float-top': '-7px',
  '--float-select-label-left': `${paddingX.value}px`,
  '--float-select-label-initial-left': `${paddingX.value + 4}px`
}))

// 标记是否已经测量完成，用于控制 label 的过渡效果
const isMeasured = ref<boolean>(false)
onMounted(() => {
  isMeasured.value = true
})

const isFocused = ref<boolean>(false)
const shouldFloat = computed<boolean>(() => {
  if (isFocused.value)
    return true

  if (Array.isArray(defaultValue)) {
    return defaultValue.length > 0
  } else if (Array.isArray(value.value)) {
    return value.value.length > 0
  } else {
    return !!defaultValue || !!value.value
  }
})

function handleFocus(): void {
  isFocused.value = true
  emit('focus')
}

function handleFocusInput(): void {
  isFocused.value = true
  emit('focusInput')
}

function handleBlur(): void {
  isFocused.value = false
  emit('blur')
}

function handleBlurInput(): void {
  isFocused.value = false
  emit('blurInput')
}

const nSelectRef = ref<SelectInst | null>(null)

defineExpose({
  focus: () => nSelectRef.value?.focus?.(),
  focusInput: () => nSelectRef.value?.focusInput?.(),
  blur: () => nSelectRef.value?.blur?.(),
  blurInput: () => nSelectRef.value?.blurInput?.()
})
</script>

<template>
  <div class="relative inline-flex flex-1" :style="cssVars">
    <NSelect
      ref="nSelectRef" v-model:value="value" placeholder=""
      v-bind="{ defaultValue, size: _size, disabled: _disabled, status: _status, ...$attrs }" @focus="handleFocus"
      @blur="handleBlur" @focus-input="handleFocusInput" @blur-input="handleBlurInput"
    >
      <template #header>
        <slot name="header" />
      </template>
      <template #action>
        <slot name="action" />
      </template>
      <template #empty>
        <slot name="empty" />
      </template>
      <template #arrow>
        <slot name="arrow" />
      </template>
    </NSelect>
    <label
      :style="commonFont" class="float-select__label absolute pointer-events-none z-2 float-select__label--input"
      :class="[{ 'float-select__label--float': shouldFloat, 'float-select__label--ready': isMeasured }]"
    >
      {{ placeholder }}
    </label>
  </div>
</template>

<style lang="css" scoped>
.float-select__label {
  font-size: var(--float-select-label-font-size);
  color: var(--float-select-label-color);
  left: var(--float-select-label-initial-left);
  /* 防止 label 过长被 suffix 和按钮遮挡 */
  max-width: calc(100% - var(--float-select-label-left) - 44px);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.float-select__label--input {
  top: 50%;
  transform: translateY(-50%);
}

.float-select__label--float {
  font-size: var(--float-select-label-float-font-size);
  color: var(--float-select-label-float-color);
  top: var(--float-select-label-float-top);
  padding: 0 var(--float-select-label-padding-x);
  line-height: 1 !important;
  background-color: var(--float-select-label-bg);
  transform: translateY(0);
}

/* 测量完成后才为 label 添加过渡效果，避免初始位置闪烁 */
.float-select__label--ready {
  transition: all 0.2s ease;
}
</style>
