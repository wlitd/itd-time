# FloatTimePicker

A floating label time picker component that wraps the Naive UI TimePicker with Material Design-style floating labels.

## Features

- Inherits all properties from Naive UI's TimePicker component
- Floating label effect that moves up when focused or has content
- Automatic adaptation to different sizes (small, medium, large)
- Form integration (inherits size, disabled, and validation status from NFormItem)
- Theme-aware styling using Naive UI theme variables
- CSS variable based customization
- Long label truncation with ellipsis
- Custom icon slot support

## Props

| Name | Type | Default | Description |
|------|------|---------|-------------|
| placeholder | `string` | - | The placeholder text that floats above the input when focused or has content |
| size | `'small' \| 'medium' \| 'large'` | `'medium'` | Time picker size |
| defaultValue | `number \| null` | - | Default value (timestamp in milliseconds) |
| status | `'warning' \| 'error'` | - | Validation status |
| disabled | `boolean` | `false` | Whether the picker is disabled |

> **Note:** All other props are inherited from Naive UI's [TimePickerProps](https://www.naiveui.com/en-US/os-theme/components/time-picker#TimePicker-Props).

## Events

| Name | Parameters | Description |
|------|------------|-------------|
| focus | `()` | Emitted when the picker gains focus |
| blur | `()` | Emitted when the picker loses focus |

> **Note:** All other events are inherited from Naive UI's TimePicker component.

## Slots

| Name | Props | Description |
|------|-------|-------------|
| icon | - | Custom icon |

> **Note:** The examples use UnoCSS icons (e.g., `i-mdi-clock`) which are automatically transformed into icon elements. Use `<div>` tags instead of `<i>` tags for icon containers.

## Methods

| Name | Parameters | Description |
|------|------------|-------------|
| focus | `()` | Focus the picker |
| blur | `()` | Blur the picker |

## Usage

```vue
<script setup lang="ts">
import type { FormInst, FormRules, FormSize } from 'naive-ui'
import type { FloatTimePickerInst } from '@/components/ui/float-time-picker/types'
import { reactive, ref } from 'vue'

// Basic
const dateValue = ref<number | null>(1183135260000)

// Events
const eventValue = ref<number | null>(861333934000)

// Methods
const timePickerRef = ref<FloatTimePickerInst>()

// Form
const size = ref<FormSize>('medium')
const disabled = ref<boolean>(false)
const formRef = ref<FormInst>()
const formData: { time: number | null } = reactive({
  time: null
})
const formRules: FormRules = {
  time: [
    { required: true, message: 'Please select time', trigger: 'blur' }
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
  <FloatTimePicker v-model:value="dateValue" placeholder="Please Select Time" />

  <!-- Sizes -->
  <NFlex vertical>
    <FloatTimePicker size="small" placeholder="Small Time Picker" />
    <FloatTimePicker placeholder="Medium Time Picker" />
    <FloatTimePicker size="large" placeholder="Large Time Picker" />
  </NFlex>

  <!-- Status -->
  <NFlex vertical>
    <FloatTimePicker disabled placeholder="Disabled" />
    <FloatTimePicker status="warning" placeholder="Warning" />
    <FloatTimePicker status="error" placeholder="Error" />
  </NFlex>

  <!-- Custom Slot -->
  <FloatTimePicker placeholder="Custom Icon">
    <template #icon>
      <div class="i-mdi-clock" />
    </template>
  </FloatTimePicker>

  <!-- Events -->
  <FloatTimePicker
    v-model:value="eventValue"
    placeholder="Event Demo"
    clearable
    @focus="() => { message.info('Focus') }"
    @blur="() => { message.info('Blur') }"
    @clear="() => { message.info('Clear') }"
    @confirm="(val, formatVal) => { message.info(`val: ${val}, formatVal: ${formatVal}`) }"
  />

  <!-- Methods -->
  <NFlex>
    <NButton @click="timePickerRef?.focus()">Focus</NButton>
    <NButton @click="timePickerRef?.blur()">Blur</NButton>
  </NFlex>
  <FloatTimePicker ref="timePickerRef" placeholder="Method Demo" />

  <!-- Form Integration -->
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
    <NFormItem path="time">
      <FloatTimePicker v-model:value="formData.time" placeholder="Time(Required)">
        <template #icon>
          <div class="i-mdi-clock" />
        </template>
      </FloatTimePicker>
    </NFormItem>
    <NFormItem>
      <RippleButton @click="handleSubmit">Submit</RippleButton>
    </NFormItem>
  </NForm>
</template>
```
