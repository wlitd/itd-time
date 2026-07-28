# FloatInput

A floating label input component that wraps the Naive UI input with Material Design-style floating labels.

## Features

- Inherits all properties from Naive UI's input component
- Floating label effect that moves up when focused or has content
- Supports text, password, and textarea types
- Automatic adaptation to different input sizes
- Built-in prefix slot support
- Theme-aware styling using Naive UI theme variables
- CSS variable based customization

## Props

| Name | Type | Default | Description |
|------|------|---------|-------------|
| placeholder | `string` | - | The placeholder text that floats above the input when focused or has content |
| type | `'text' \| 'password' \| 'textarea'` | `'text'` | Input type |
| size | `'tiny' \| 'small' \| 'medium' \| 'large'` | `'medium'` | Input size |
| pair | `boolean` | `false` | Whether to use pair input (not supported yet) |
| defaultValue | `string` | - | Default value of the input |
| status | `'success' \| 'warning' \| 'error'` | - | Validation status |
| disabled | `boolean` | `false` | Whether the input is disabled |

> **Note:** All other props are inherited from Naive UI's [InputProps](https://www.naiveui.com/en-US/os-theme/components/input#Input-Props).

## Events

All events are inherited from Naive UI's input component.

## Slots

| Name | Description |
|------|-------------|
| prefix | Content placed before the input field |
| suffix | Content placed after the input field |
| clear-icon | Custom clear icon |
| count | Custom count display (receives `{ value: string }` slot props) |
| password-visible-icon | Custom password visible icon |
| password-invisible-icon | Custom password invisible icon |
| default | The input content (rarely used) |

> **Note:** The examples use UnoCSS icons (e.g., `i-mdi:lock`) which are automatically transformed into icon elements. You can use any icon component or SVG in the prefix slot.

## Methods

| Name | Parameters | Description |
|------|------------|-------------|
| clear | `()` | Clear the input value |
| focus | `()` | Focus the input |
| blur | `()` | Blur the input |
| select | `()` | Select all text in the input |
| scrollToEnd | `(options: { left?: number, top?: number, behavior: 'auto' \| 'smooth' })` | Scroll to the end of the input |

## Usage

```vue
<script setup lang="ts">
import type { FormInst, FormItemRule, FormRules } from 'naive-ui'
import type { Size } from 'naive-ui/es/form/src/interface'
import type { FloatInputInst } from '@/components/ui/float-input/types'
import { reactive, ref } from 'vue'

const value = ref<string>('')
const text = ref<string>('')

const methodValue = ref<string>('')
const floatInputRef = ref<FloatInputInst>()

const size = ref<Size>('medium')
const disabled = ref<boolean>(false)
const formRef = ref<FormInst>()
const formData: { email: string, password: string } = reactive({
  email: '',
  password: ''
})
const formRules: FormRules = {
  email: [
    { required: true, message: 'Please enter email', trigger: 'blur' }
  ],
  password: [
    {
      required: true,
      validator: (__: FormItemRule, value: string) => {
        if (!value) {
          return new Error('Please enter password')
        } else if (value.length < 6) {
          return new Error('Password must be at least 6 characters')
        }
        return true
      },
      trigger: 'blur'
    }
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
  <FloatInput v-model:value="value" placeholder="Please Input" />
  <FloatInput v-model:value="text" type="textarea" placeholder="Please Input" />

  <NFlex vertical>
    <FloatInput size="tiny" placeholder="Tiny Input" />
    <FloatInput size="small" placeholder="Small Input" />
    <FloatInput placeholder="Medium Input" />
    <FloatInput size="large" placeholder="Large Input" />
  </NFlex>

  <NFlex vertical>
    <FloatInput round size="small" placeholder="Small Input" />
    <FloatInput round placeholder="Medium Input" />
    <FloatInput round size="large" placeholder="Large Input" />
  </NFlex>

  <NFlex vertical>
    <FloatInput placeholder="Prefix">
      <template #prefix>
        <div class="i-mdi:email" />
      </template>
    </FloatInput>
    <FloatInput placeholder="Suffix">
      <template #suffix>
        <div class="i-mdi:magnify" />
      </template>
    </FloatInput>
    <FloatInput placeholder="ClearIcon" clearable>
      <template #clear-icon>
        <div class="i-mdi:delete" />
      </template>
    </FloatInput>
    <FloatInput placeholder="Password" type="password" show-password-on="click">
      <template #password-visible-icon>
        <div class="i-mdi:eye" />
      </template>
      <template #password-invisible-icon>
        <div class="i-mdi:eye-off" />
      </template>
    </FloatInput>
    <FloatInput placeholder="Count" show-count>
      <template #count="{ value: inputStr }">
        <NBadge :value="inputStr.length" />
      </template>
    </FloatInput>
  </NFlex>

  <FloatInput loading placeholder="Loading" />
  <FloatInput disabled placeholder="Disabled" />
  <FloatInput placeholder="Warning" status="warning" />
  <FloatInput placeholder="Error" status="error" />

  <NFlex>
    <RippleButton @click="() => { floatInputRef?.clear() }">Clear</RippleButton>
    <RippleButton @click="() => { floatInputRef?.focus() }">Focus</RippleButton>
    <RippleButton @click="() => { floatInputRef?.blur() }">Blur</RippleButton>
    <RippleButton @click="() => { floatInputRef?.select() }">Select</RippleButton>
    <RippleButton @click="() => { floatInputRef?.scrollToEnd({ behavior: 'smooth', left: 1000 }) }">ScrollToEnd</RippleButton>
  </NFlex>
  <FloatInput ref="floatInputRef" v-model:value="methodValue" placeholder="Method" />

  <FloatInput placeholder="Event" @focus="() => { message.info('Focus') }" @blur="() => { message.info('Blur') }" @input="() => { message.info('Input') }" @change="() => { message.info('Change') }" @keyup="() => { message.info('Keyup') }" />

  <FloatInput autosize style="min-width: 50%" placeholder="Autosize" />
  <NInputGroup>
    <NInputGroupLabel>https://www.</NInputGroupLabel>
    <FloatInput placeholder="Group(It's so ugly)" />
    <NInputGroupLabel>.com</NInputGroupLabel>
  </NInputGroup>
  <FloatInput maxlength="30" show-count placeholder="Count" />
  <FloatInput :input-props="{ type: 'url' }" placeholder="Native attribute" />
  <FloatInput :allow-input="(value) => !value || /^\d+$/.test(value)" placeholder="Only Allow input digit" />

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
    <NFormItem path="email">
      <FloatInput v-model:value="formData.email" placeholder="Email(Required)">
        <template #prefix>
          <div class="i-mdi:email" />
        </template>
      </FloatInput>
    </NFormItem>
    <NFormItem path="password">
      <FloatInput v-model:value="formData.password" type="password" show-password-on="click" placeholder="Password(Required)">
        <template #prefix>
          <div class="i-mdi:lock" />
        </template>
      </FloatInput>
    </NFormItem>
    <NFormItem>
      <RippleButton @click="handleSubmit">Submit</RippleButton>
    </NFormItem>
  </NForm>
</template>
```
