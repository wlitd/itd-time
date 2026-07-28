# FloatInputNumber

A floating label input number component that wraps the Naive UI input number with Material Design-style floating labels.

## Features

- Inherits all properties from Naive UI's input number component
- Floating label effect that moves up when focused or has content
- Automatic adaptation to different input sizes
- Built-in prefix, suffix, and custom button slot support
- Theme-aware styling using Naive UI theme variables
- CSS variable based customization
- Custom parse and format function support (e.g., thousand separators)

## Props

| Name | Type | Default | Description |
|------|------|---------|-------------|
| placeholder | `string` | - | The placeholder text that floats above the input when focused or has content |
| size | `'tiny' \| 'small' \| 'medium' \| 'large'` | `'medium'` | Input size |
| defaultValue | `number \| null` | - | Default value of the input |
| status | `'success' \| 'warning' \| 'error'` | - | Validation status |
| disabled | `boolean` | `false` | Whether the input is disabled |
| buttonPlacement | `'right' \| 'both'` | `'right'` | Placement of increment/decrement buttons |

> **Note:** All other props are inherited from Naive UI's [InputNumberProps](https://www.naiveui.com/en-US/os-theme/components/input-number#InputNumber-Props).

## Events

| Name | Parameters | Description |
|------|------------|-------------|
| focus | `()` | Emitted when the input receives focus |
| blur | `()` | Emitted when the input loses focus |

> **Note:** All other events are inherited from Naive UI's input number component.

## Slots

| Name | Description |
|------|-------------|
| prefix | Content placed before the input field |
| suffix | Content placed after the input field |
| minus-icon | Custom decrement button icon (requires `button-placement="both"`) |
| add-icon | Custom increment button icon |

> **Note:** The examples use UnoCSS icons (e.g., `i-mdi:arrow-up-bold`) which are automatically transformed into icon elements. You can use any icon component or SVG in the slots.

## Methods

| Name | Parameters | Description |
|------|------------|-------------|
| focus | `()` | Focus the input |
| blur | `()` | Blur the input |

## Usage

```vue
<script setup lang="ts">
import type { FormInst, FormRules, FormSize } from 'naive-ui'
import type { FloatInputNumberInst } from '@/components/ui/float-input-number/types'
import { reactive, ref } from 'vue'

const value = ref<number>(0)
const methodValue = ref<number | null>(null)
const floatInputNumberRef = ref<FloatInputNumberInst>()

// Custom parse function (for thousand separator demo)
function parseThousand(input: string): number | null {
  const nums = input.replace(/,/g, '').trim()
  if (/^\d+(?:\.\d+)?$/.test(nums)) {
    return Number(nums)
  }
  return nums === '' ? null : Number.NaN
}

function formatThousand(value: number | null): string {
  if (value === null) {
    return ''
  }
  return value.toLocaleString('en-US')
}

const size = ref<FormSize>('medium')
const disabled = ref<boolean>(false)
const formRef = ref<FormInst>()
const formData: { amount: number | null } = reactive({
  amount: null
})
const formRules: FormRules = {
  amount: [
    { required: true, type: 'number', message: 'Please enter amount', trigger: 'blur' }
  ]
}
function handleSubmit(): void {
  formRef.value?.validate((errors) => {
    if (!errors) {
      message.success('Submit successfully')
    } else {
      message.error('Submit failed')
    }
  })
}
</script>

<template>
  <!-- Basic -->
  <FloatInputNumber v-model:value="value" placeholder="Please Input Number" />

  <!-- Size -->
  <NFlex vertical>
    <FloatInputNumber size="tiny" placeholder="Tiny Input" />
    <FloatInputNumber size="small" placeholder="Small Input" />
    <FloatInputNumber placeholder="Medium Input" />
    <FloatInputNumber size="large" placeholder="Large Input" />
  </NFlex>

  <!-- Slots -->
  <NFlex vertical>
    <FloatInputNumber placeholder="Prefix">
      <template #prefix>
        <span>￥</span>
      </template>
    </FloatInputNumber>
    <FloatInputNumber placeholder="Suffix">
      <template #suffix>
        <span>%</span>
      </template>
    </FloatInputNumber>
    <FloatInputNumber placeholder="Custom Buttons with Prefix" button-placement="both">
      <template #prefix>
        <span>￥</span>
      </template>
      <template #minus-icon>
        <div class="i-mdi:arrow-down-bold" />
      </template>
      <template #add-icon>
        <div class="i-mdi:arrow-up-bold" />
      </template>
    </FloatInputNumber>
  </NFlex>

  <!-- Status -->
  <FloatInputNumber loading placeholder="Loading" />
  <FloatInputNumber disabled placeholder="Disabled" />
  <FloatInputNumber placeholder="Warning" status="warning" />
  <FloatInputNumber placeholder="Error" status="error" />

  <!-- Method -->
  <NFlex>
    <RippleButton @click="() => { floatInputNumberRef?.focus() }">Focus</RippleButton>
    <RippleButton @click="() => { floatInputNumberRef?.blur() }">Blur</RippleButton>
  </NFlex>
  <FloatInputNumber ref="floatInputNumberRef" v-model:value="methodValue" placeholder="Method" />

  <!-- Event -->
  <FloatInputNumber
    placeholder="Event"
    @focus="() => { message.info('Focus') }"
    @blur="() => { message.info('Blur') }"
  />

  <!-- Other -->
  <FloatInputNumber :show-button="false" placeholder="Hide Buttons" />
  <FloatInputNumber :min="0" :max="100" :step="5" placeholder="With Range (0-100, step 5)" />
  <FloatInputNumber :precision="2" placeholder="2 Decimal Places" />
  <FloatInputNumber
    placeholder="Thousand Separator"
    :default-value="1075"
    :parse="parseThousand"
    :format="formatThousand"
  />

  <!-- Form -->
  <NFlex align="center">
    <span>Size:</span>
    <NRadioGroup v-model:value="size">
      <NRadioButton value="small">Small</NRadioButton>
      <NRadioButton value="medium">Medium</NRadioButton>
      <NRadioButton value="large">Large</NRadioButton>
    </NRadioGroup>
  </NFlex>
  <NFlex align="center">
    <span>Disabled:</span>
    <NSwitch v-model:value="disabled" />
  </NFlex>
  <NForm ref="formRef" inline :size="size" :disabled="disabled" :model="formData" :rules="formRules">
    <NFormItem path="amount">
      <FloatInputNumber v-model:value="formData.amount" clearable placeholder="Amount(Required)">
        <template #prefix>
          <span>￥</span>
        </template>
      </FloatInputNumber>
    </NFormItem>
    <NFormItem>
      <RippleButton @click="handleSubmit">Submit</RippleButton>
    </NFormItem>
  </NForm>
</template>
```
