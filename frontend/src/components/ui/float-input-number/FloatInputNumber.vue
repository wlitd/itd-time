<script setup lang="ts">
import type { ButtonSize, FormValidationStatus, InputNumberInst } from 'naive-ui'
import type { FormItemInjection } from 'naive-ui/es/_mixins/use-form-item'
import type { FloatInputNumberProps } from './types'
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
  buttonPlacement
} = defineProps<FloatInputNumberProps>()

const emit = defineEmits<{
  focus: []
  blur: []
}>()

const slots = defineSlots<{
  prefix?: () => any
  suffix?: () => any
  'minus-icon'?: () => any
  'add-icon'?: () => any
}>()

const value = defineModel<number | null>('value')

// 在 NFormItem 中使用时可以注入表单项上下文
const formItemContext = inject<FormItemInjection | null>(formItemInjectionKey, null)
const _size = computed<ButtonSize>(() => formItemContext?.mergedSize.value || size)
const _disabled = computed<boolean>(() => formItemContext?.disabled.value || disabled)
const _status = computed<FormValidationStatus | undefined>(() => formItemContext?.mergedValidationStatus.value || status)

// 获取前缀宽度
const prefixEl = useTemplateRef('prefixRef')
const { width: prefixWidth } = useElementSize(prefixEl)

// 获取自定义减号按钮宽度
const minusButtonEl = useTemplateRef('minusButtonRef')
const { width: customMinusButtonWidth } = useElementSize(minusButtonEl)

// 左侧按钮宽度计算逻辑：
// 1. button-placement !== 'both': 只考虑 prefix 宽度，左侧按钮宽度为 0
// 2. button-placement === 'both' 且无自定义图标: 左侧按钮宽度为 18px（默认值）
// 3. button-placement === 'both' 且有自定义图标: 测量自定义图标宽度
const minusButtonWidth = computed<number>(() => {
  if (buttonPlacement !== 'both') {
    return 0
  }
  // buttonPlacement === 'both'
  if (slots['minus-icon']) {
    // 有自定义图标，测量宽度
    return customMinusButtonWidth.value > 0 ? customMinusButtonWidth.value : 18
  }
  // 无自定义图标，使用默认值
  return 18
})

// Prefix 和减号按钮之间的间距（当两者都存在时）
const prefixMarginLeft = computed<number>(() => {
  // 只有当同时存在减号按钮和 prefix 时，Naive UI 才会添加 10px 的 margin-left
  if (buttonPlacement === 'both' && slots.prefix) {
    return 10
  }
  return 0
})

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
  '--float-input-number-label-font-size': fontSize.value,
  '--float-input-number-label-float-font-size': smallerFontSize.value,
  '--float-input-number-label-color': statusTextColor.value || themeVars.value.placeholderColor,
  '--float-input-number-label-float-color': statusTextColor.value || themeVars.value.textColor1,
  '--float-input-number-label-bg': themeVars.value.tagColor,
  '--float-input-number-label-padding-x': '2px',
  '--float-input-number-label-float-top': '-7px',
  // label 位置 = prefix宽度 + 减号按钮宽度 + prefix与按钮之间的间距 + 左侧内边距
  '--float-input-number-label-left': `${prefixWidth.value + minusButtonWidth.value + prefixMarginLeft.value + paddingX.value}px`,
  '--float-input-number-label-initial-left': `${prefixWidth.value + minusButtonWidth.value + prefixMarginLeft.value + paddingX.value + 4}px`
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

const nInputNumberRef = ref<InputNumberInst | null>(null)

defineExpose({
  focus: () => nInputNumberRef.value?.focus(),
  blur: () => nInputNumberRef.value?.blur()
})
</script>

<template>
  <div class="relative inline-flex flex-1" :style="cssVars">
    <NInputNumber
      ref="nInputNumberRef" v-model:value="value" placeholder=""
      v-bind="{ defaultValue, size: _size, disabled: _disabled, status: _status, buttonPlacement, ...$attrs }"
      @focus="handleFocus" @blur="handleBlur"
    >
      <template #prefix>
        <div v-if="$slots.prefix" ref="prefixRef">
          <slot name="prefix" />
        </div>
      </template>
      <template #minus-icon>
        <div v-if="$slots['minus-icon']" ref="minusButtonRef" class="inline-flex items-center justify-center">
          <slot name="minus-icon" />
        </div>
      </template>
      <template #add-icon>
        <slot name="add-icon" />
      </template>
      <template #suffix>
        <slot name="suffix" />
      </template>
    </NInputNumber>
    <label
      :style="commonFont"
      class="float-input-number__label absolute pointer-events-none z-2 float-input-number__label--input"
      :class="[{ 'float-input-number__label--float': shouldFloat, 'float-input-number__label--ready': isMeasured }]"
    >
      {{ placeholder }}
    </label>
  </div>
</template>

<style lang="css" scoped>
.float-input-number__label {
  font-size: var(--float-input-number-label-font-size);
  color: var(--float-input-number-label-color);
  left: var(--float-input-number-label-initial-left);
  /* 防止 label 过长被 suffix 和按钮遮挡 */
  max-width: calc(100% - var(--float-input-number-label-left) - 44px);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.float-input-number__label--input {
  top: 50%;
  transform: translateY(-50%);
}

.float-input-number__label--float {
  font-size: var(--float-input-number-label-float-font-size);
  color: var(--float-input-number-label-float-color);
  top: var(--float-input-number-label-float-top);
  padding: 0 var(--float-input-number-label-padding-x);
  line-height: 1 !important;
  background-color: var(--float-input-number-label-bg);
  transform: translateY(0);
}

/* 测量完成后才为 label 添加过渡效果，避免初始位置闪烁 */
.float-input-number__label--ready {
  transition: all 0.2s ease;
}
</style>
