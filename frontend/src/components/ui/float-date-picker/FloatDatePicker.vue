<script setup lang="ts">
import type { DatePickerInst, DatePickerProps, DatePickerSize, FormValidationStatus } from 'naive-ui'
import type { FormItemInjection } from 'naive-ui/es/_mixins/use-form-item'
import type { FloatDatePickerProps } from './types'
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
  disabled = false,
  type = 'date',
  panel
} = defineProps<FloatDatePickerProps>()

const emit = defineEmits<{
  focus: []
  blur: []
}>()

defineSlots<{
  'date-icon'?: () => any
  footer?: () => any
  'next-month'?: () => any
  'next-year'?: () => any
  'prev-month'?: () => any
  'prev-year'?: () => any
  // Date, Year, QuarterRange, Week 专用插槽
  clear?: (props: { onClear: () => void, text: string }) => any
  now?: (props: { onNow: () => void, text: string }) => any
  // DateTime, Month, Quarter 专用插槽
  confirm?: (props: { onConfirm: () => void, disabled: boolean, text: string }) => any
}>()

// 检查不支持的类型
const unsupportedTypes: DatePickerProps['type'][] = ['daterange', 'datetimerange', 'monthrange', 'yearrange', 'quarterrange']
if (unsupportedTypes.includes(type)) {
  console.error(`FloatDatePicker does not support range types: ${type}. Please use non-range types only.`)
}

// 检查不支持的 panel 模式
if (panel) {
  console.error('FloatDatePicker does not support panel mode. Please use it as a regular input.')
}

const value = defineModel<number | null>('value')

// 在 NFormItem 中使用时可以注入表单项上下文
const formItemContext = inject<FormItemInjection | null>(formItemInjectionKey, null)
const _size = computed<DatePickerSize>(() => formItemContext?.mergedSize.value || size)
const _disabled = computed<boolean>(() => formItemContext?.disabled.value || disabled)
const _status = computed<FormValidationStatus | undefined>(() => formItemContext?.mergedValidationStatus.value || status)

// DatePicker 没有 prefix slot，不需要测量前缀宽度

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
  '--float-date-picker-label-font-size': fontSize.value,
  '--float-date-picker-label-float-font-size': smallerFontSize.value,
  '--float-date-picker-label-color': statusTextColor.value || themeVars.value.placeholderColor,
  '--float-date-picker-label-float-color': statusTextColor.value || themeVars.value.textColor1,
  '--float-date-picker-label-bg': themeVars.value.tagColor,
  '--float-date-picker-label-padding-x': '2px',
  '--float-date-picker-label-float-top': '-7px',
  '--float-date-picker-label-left': `${paddingX.value}px`,
  '--float-date-picker-label-initial-left': `${paddingX.value + 4}px`
}))

// 标记是否已经测量完成，用于控制 label 的过渡效果
const isMeasured = ref<boolean>(false)
// DatePicker 没有 prefix slot，直接添加过渡效果
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

const nDatePickerRef = ref<DatePickerInst | null>(null)

defineExpose({
  focus: () => nDatePickerRef.value?.focus?.(),
  blur: () => nDatePickerRef.value?.blur?.()
})
</script>

<template>
  <div class="relative inline-flex flex-1" :style="cssVars">
    <NDatePicker
      ref="nDatePickerRef" v-model:value="value" placeholder=""
      v-bind="{ defaultValue, type, size: _size, disabled: _disabled, status: _status, panel: false, ...$attrs }"
      @focus="handleFocus" @blur="handleBlur"
    >
      <template #date-icon>
        <slot name="date-icon" />
      </template>
      <template v-if="$slots.footer" #footer>
        <slot name="footer" />
      </template>
      <template #next-month>
        <slot name="next-month" />
      </template>
      <template #next-year>
        <slot name="next-year" />
      </template>
      <template #prev-month>
        <slot name="prev-month" />
      </template>
      <template #prev-year>
        <slot name="prev-year" />
      </template>
      <template #clear="props">
        <slot name="clear" v-bind="props" />
      </template>
      <template #now="props">
        <slot name="now" v-bind="props" />
      </template>
      <template #confirm="props">
        <slot name="confirm" v-bind="props" />
      </template>
    </NDatePicker>
    <label
      :style="commonFont"
      class="float-date-picker__label absolute pointer-events-none z-2 float-date-picker__label--input"
      :class="[{ 'float-date-picker__label--float': shouldFloat, 'float-date-picker__label--ready': isMeasured }]"
    >
      {{ placeholder }}
    </label>
  </div>
</template>

<style lang="css" scoped>
.float-date-picker__label {
  font-size: var(--float-date-picker-label-font-size);
  color: var(--float-date-picker-label-color);
  left: var(--float-date-picker-label-initial-left);
  /* 防止 label 过长被 suffix 和按钮遮挡 */
  max-width: calc(100% - var(--float-date-picker-label-left) - 44px);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.float-date-picker__label--input {
  top: 50%;
  transform: translateY(-50%);
}

.float-date-picker__label--float {
  font-size: var(--float-date-picker-label-float-font-size);
  color: var(--float-date-picker-label-float-color);
  top: var(--float-date-picker-label-float-top);
  padding: 0 var(--float-date-picker-label-padding-x);
  line-height: 1 !important;
  background-color: var(--float-date-picker-label-bg);
  transform: translateY(0);
}

/* 测量完成后才为 label 添加过渡效果，避免初始位置闪烁 */
.float-date-picker__label--ready {
  transition: all 0.2s ease;
}
</style>
