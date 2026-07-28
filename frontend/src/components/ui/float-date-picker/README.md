# FloatDatePicker

A floating label date picker component that wraps the Naive UI DatePicker with Material Design-style floating labels.

## Features

- Inherits all properties from Naive UI's DatePicker component
- Floating label effect that moves up when focused or has content
- Supports multiple types: date, datetime, month, year, quarter, week
- Custom date format configuration
- Automatic adaptation to different sizes (small, medium, large)
- Form integration (inherits size, disabled, and validation status from NFormItem)
- Theme-aware styling using Naive UI theme variables
- CSS variable based customization
- Long label truncation with ellipsis
- Scoped slots for clear, now, and confirm buttons

## Props

| Name | Type | Default | Description |
|------|------|---------|-------------|
| placeholder | `string` | - | The placeholder text that floats above the input when focused or has content |
| type | `'date' \| 'datetime' \| 'month' \| 'year' \| 'quarter' \| 'week'` | `'date'` | Picker type (range types are not supported) |
| size | `'small' \| 'medium' \| 'large'` | `'medium'` | Picker size |
| defaultValue | `number \| null` | - | Default value (timestamp in milliseconds) |
| status | `'warning' \| 'error'` | - | Validation status |
| disabled | `boolean` | `false` | Whether the picker is disabled |
| panel | `boolean` | - | Panel mode (not supported, will log error if set to true) |

> **Note:** All other props are inherited from Naive UI's [DatePickerProps](https://www.naiveui.com/en-US/os-theme/components/date-picker#DatePicker-Props). Range types (`daterange`, `datetimerange`, `monthrange`, `yearrange`, `quarterrange`) are **not supported**.

## Events

| Name | Parameters | Description |
|------|------------|-------------|
| focus | `()` | Emitted when the picker gains focus |
| blur | `()` | Emitted when the picker loses focus |
| next-month | `()` | Emitted when clicking the next month button |
| prev-month | `()` | Emitted when clicking the previous month button |
| next-year | `()` | Emitted when clicking the next year button |
| prev-year | `()` | Emitted when clicking the previous year button |

> **Note:** All other events are inherited from Naive UI's DatePicker component.

## Slots

| Name | Props | Description |
|------|-------|-------------|
| date-icon | - | Custom date icon |
| footer | - | Custom footer content |
| next-month | - | Custom next month button |
| next-year | - | Custom next year button |
| prev-month | - | Custom previous month button |
| prev-year | - | Custom previous year button |
| clear | `{ onClear: () => void, text: string }` | Custom clear button (Date, Year, QuarterRange, Week types) |
| now | `{ onNow: () => void, text: string }` | Custom "Now" button (Date, Year, QuarterRange, Week types) |
| confirm | `{ onConfirm: () => void, disabled: boolean, text: string }` | Custom confirm button (DateTime, Month, Quarter types) |

> **Note:** The examples use UnoCSS icons (e.g., `i-mdi-calendar-clock`) which are automatically transformed into icon elements. Use `<div>` tags instead of `<i>` tags for icon containers.

## Methods

| Name | Parameters | Description |
|------|------------|-------------|
| focus | `()` | Focus the picker |
| blur | `()` | Blur the picker |

## Usage

```vue
<script setup lang="ts">
import type { FormInst, FormRules } from 'naive-ui'
import type { FloatDatePickerInst } from '@/components/ui/float-date-picker/types'
import { reactive, ref } from 'vue'

// Basic
const dateValue = ref<number | null>(1780329600000)

// Types with format
const yearValue = ref<number | null>(1780329600000)
const monthValue = ref<number | null>(1780329600000)
const datetimeValue = ref<number | null>(1780329600000)

// Events
const eventValue = ref<number | null>(1780329600000)

// Methods
const datePickerRef = ref<FloatDatePickerInst>()

// Form
const size = ref('medium')
const disabled = ref<boolean>(false)
const formRef = ref<FormInst>()
const formData: { startDate: number | null } = reactive({
  startDate: null
})
const formRules: FormRules = {
  startDate: [
    { required: true, message: 'Please select start date', trigger: 'blur' }
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
  <FloatDatePicker v-model:value="dateValue" placeholder="Please Select Date" />

  <!-- Sizes -->
  <NFlex vertical>
    <FloatDatePicker size="small" placeholder="Small Date Picker" />
    <FloatDatePicker placeholder="Medium Date Picker" />
    <FloatDatePicker size="large" placeholder="Large Date Picker" />
  </NFlex>

  <!-- Status -->
  <NFlex vertical>
    <FloatDatePicker disabled placeholder="Disabled" />
    <FloatDatePicker status="warning" placeholder="Warning" />
    <FloatDatePicker status="error" placeholder="Error" />
  </NFlex>

  <!-- Custom Slots -->
  <FloatDatePicker placeholder="Custom Date Icon & Navigation">
    <template #date-icon>
      <div class="i-mdi-calendar-clock" />
    </template>
    <template #prev-month>
      <div class="i-mdi-chevron-left" />
    </template>
    <template #next-month>
      <div class="i-mdi-chevron-right" />
    </template>
    <template #prev-year>
      <div class="i-mdi-chevron-double-left" />
    </template>
    <template #next-year>
      <div class="i-mdi-chevron-double-right" />
    </template>
  </FloatDatePicker>

  <!-- Scoped Slots -->
  <FloatDatePicker type="month" placeholder="Custom Now Button">
    <template #now="{ onNow, text }">
      <NButton size="tiny" @click="onNow">{{ text }} 📅</NButton>
    </template>
  </FloatDatePicker>

  <FloatDatePicker type="datetime" placeholder="Custom Confirm Button">
    <template #confirm="{ onConfirm, disabled: _disabled, text }">
      <NButton type="primary" size="tiny" :disabled="_disabled" @click="onConfirm">
        {{ text }} ✓
      </NButton>
    </template>
  </FloatDatePicker>

  <FloatDatePicker placeholder="Custom Footer">
    <template #footer>
      <NText italic>Custom footer content</NText>
    </template>
  </FloatDatePicker>

  <!-- Events -->
  <FloatDatePicker
    v-model:value="eventValue"
    placeholder="Event Demo"
    @focus="() => { message.info('Focus') }"
    @blur="() => { message.info('Blur') }"
    @next-month="() => { message.info('Next Month') }"
    @prev-month="() => { message.info('Prev Month') }"
    @next-year="() => { message.info('Next Year') }"
    @prev-year="() => { message.info('Prev Year') }"
  />

  <!-- Methods -->
  <NFlex>
    <NButton @click="datePickerRef?.focus()">Focus</NButton>
    <NButton @click="datePickerRef?.blur()">Blur</NButton>
  </NFlex>
  <FloatDatePicker ref="datePickerRef" placeholder="Method Demo" />

  <!-- Types with Format -->
  <FloatDatePicker v-model:value="yearValue" type="year" placeholder="Year Picker" />
  <FloatDatePicker v-model:value="yearValue" type="year" format="y年" year-format="y年" placeholder="Year with Format" />
  <FloatDatePicker v-model:value="monthValue" type="month" format="y年 M月" year-format="y年" month-format="M月" placeholder="Month with Format" />
  <FloatDatePicker v-model:value="dateValue" type="date" format="y年 M月 d日" placeholder="Date with Format" />
  <FloatDatePicker v-model:value="datetimeValue" type="datetime" format="y年 M月 d日 HH:mm:ss" placeholder="DateTime with Format" />

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
    <NFormItem path="startDate">
      <FloatDatePicker v-model:value="formData.startDate" placeholder="Start Date(Required)">
        <template #date-icon>
          <div class="i-mdi-calendar-clock" />
        </template>
      </FloatDatePicker>
    </NFormItem>
    <NFormItem>
      <RippleButton @click="handleSubmit">Submit</RippleButton>
    </NFormItem>
  </NForm>
</template>
```
