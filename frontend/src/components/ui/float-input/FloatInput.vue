<script setup lang="ts">
import type { ButtonSize, FormValidationStatus, InputInst } from 'naive-ui'
import type { FormItemInjection } from 'naive-ui/es/_mixins/use-form-item'
import type { FloatInputProps } from './types'
import { useThemeVars } from 'naive-ui'
import { formItemInjectionKey } from 'naive-ui/es/_mixins/use-form-item'

defineOptions({
  inheritAttrs: false
})

const {
  placeholder,
  pair = false,
  type = 'text',
  size = 'medium',
  defaultValue,
  status,
  disabled = false
} = defineProps<FloatInputProps>()

const emit = defineEmits<{
  focus: []
  blur: []
}>()

const slots = defineSlots<{
  prefix?: () => any
  'clear-icon'?: () => any
  count?: (props: { value: string }) => any
  'password-invisible-icon'?: () => any
  'password-visible-icon'?: () => any
  suffix?: () => any
}>()

if (pair) {
  console.error('FloatInput does not support pair for now, please use the native NInput!')
}

const value = defineModel<string | null>('value')

// 在 NFormItem 中使用时可以注入表单项上下文
const formItemContext = inject<FormItemInjection | null>(formItemInjectionKey, null)
const _size = computed<ButtonSize>(() => formItemContext?.mergedSize.value || size)
const _disabled = computed<boolean>(() => formItemContext?.disabled.value || disabled)
const _status = computed<FormValidationStatus | undefined>(() => formItemContext?.mergedValidationStatus.value || status)

// 获取前缀宽度
const prefixEl = useTemplateRef('prefixRef')
const { width: prefixWidth } = useElementSize(prefixEl)

// 使用适配 Naive UI 的主题变量
const themeVars = useThemeVars()
const { fontSizeMap, fontSize, paddingX, commonFont } = useNaiveSize(_size.value)

const textareaPaddingYMap: Record<ButtonSize, number> = { tiny: 2, small: 3.5, medium: 6.5, large: 8.75 }

const textareaPaddingY = computed<number>(() => textareaPaddingYMap[_size.value])

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
  '--float-input-label-font-size': fontSize.value,
  '--float-input-label-float-font-size': smallerFontSize.value,
  '--float-input-label-color': statusTextColor.value || themeVars.value.placeholderColor,
  '--float-input-label-float-color': statusTextColor.value || themeVars.value.textColor1,
  '--float-input-label-bg': themeVars.value.tagColor,
  '--float-input-label-padding-x': '2px',
  '--float-input-label-float-top': type === 'textarea' ? `${-2 - textareaPaddingY.value}px` : '-2px',
  '--float-input-label-left': `${prefixWidth.value + paddingX.value}px`,
  '--float-input-label-initial-left': `${prefixWidth.value + paddingX.value + 4}px`
}))

// 标记 prefix slot 是否已经测量完成，用于控制 label 的过渡效果
const isMeasured = ref<boolean>(false)
// 当没有 prefix slot 时，直接添加过渡效果
onMounted(() => {
  if (!slots.prefix) {
    isMeasured.value = true
  }
})

// 如果有 prefix slot，则等待其测量完成
watch(() => prefixWidth.value, () => {
  nextTick(() => {
    requestAnimationFrame(() => {
      isMeasured.value = true
    })
  })
}, { once: true })

// 非浮动状态样式（通过 class 切换，减少内联样式计算）
const initialLabelClass = computed<Record<string, boolean>>(() => ({
  'float-input__label--textarea': type === 'textarea',
  'float-input__label--input': type !== 'textarea'
}))

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

const nInputRef = ref<InputInst | null>(null)

defineExpose({
  clear: () => nInputRef.value?.clear(),
  focus: () => nInputRef.value?.focus(),
  blur: () => nInputRef.value?.blur(),
  select: () => nInputRef.value?.select(),
  scrollToEnd: ({ left, top, behavior }: { left?: number, top?: number, behavior: 'auto' | 'smooth' }) => nInputRef.value?.scrollTo({ left, top, behavior })
})
</script>

<template>
  <div v-if="!pair" class="relative inline-flex flex-1" :style="cssVars">
    <NInput
      ref="nInputRef" v-model:value="value" placeholder=""
      v-bind="{ defaultValue, type, size: _size, disabled: _disabled, status: _status, ...$attrs }" @focus="handleFocus"
      @blur="handleBlur"
    >
      <template #prefix>
        <div v-if="$slots.prefix" ref="prefixRef">
          <slot name="prefix" />
        </div>
      </template>
      <template #clear-icon>
        <slot name="clear-icon" />
      </template>
      <template #count="{ value: inputStr }">
        <slot name="count" :value="inputStr" />
      </template>
      <template #password-invisible-icon>
        <slot name="password-invisible-icon" />
      </template>
      <template #password-visible-icon>
        <slot name="password-visible-icon" />
      </template>
      <template #suffix>
        <slot name="suffix" />
      </template>
    </NInput>
    <label
      :style="commonFont" class="float-input__label absolute pointer-events-none z-2"
      :class="[initialLabelClass, { 'float-input__label--float': shouldFloat, 'float-input__label--ready': isMeasured }]"
    >
      {{ placeholder }}
    </label>
  </div>
</template>

<style lang="css" scoped>
.float-input__label {
  font-size: var(--float-input-label-font-size);
  color: var(--float-input-label-color);
  left: var(--float-input-label-initial-left);
  /* 防止 label 过长被 suffix 和按钮遮挡 */
  max-width: calc(100% - var(--float-input-label-left) - 44px);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.float-input__label--input {
  top: 50%;
  transform: translateY(-50%);
}

.float-input__label--textarea {
  top: v-bind('`${textareaPaddingY}px`');
}

.float-input__label--float {
  font-size: var(--float-input-label-float-font-size);
  color: var(--float-input-label-float-color);
  top: var(--float-input-label-float-top);
  padding: 0 var(--float-input-label-padding-x);
  line-height: 1 !important;
  background-color: var(--float-input-label-bg);
}

/* 测量完成后才为 label 添加过渡效果，避免初始位置闪烁 */
.float-input__label--ready {
  transition: all 0.2s ease;
}
</style>
