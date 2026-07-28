# FloatSelect

A floating label select component that wraps the Naive UI Select with Material Design-style floating labels.

## Features

- Inherits all properties from Naive UI's Select component
- Floating label effect that moves up when focused or has content
- Automatic adaptation to different sizes (tiny, small, medium, large)
- Form integration (inherits size, disabled, and validation status from NFormItem)
- Theme-aware styling using Naive UI theme variables
- CSS variable based customization
- Long label truncation with ellipsis
- Custom slots support (header, action, empty, arrow)
- Multiple selection support
- Filterable option support

## Props

| Name | Type | Default | Description |
|------|------|---------|-------------|
| placeholder | `string` | - | The placeholder text that floats above the input when focused or has content |
| size | `'tiny' \| 'small' \| 'medium' \| 'large'` | `'medium'` | Select size |
| defaultValue | `Value \| null` | - | Default value |
| status | `'warning' \| 'error'` | - | Validation status |
| disabled | `boolean` | `false` | Whether the select is disabled |

> **Note:** All other props are inherited from Naive UI's [SelectProps](https://www.naiveui.com/en-US/os-theme/components/select#Select-Props).

## Events

| Name | Parameters | Description |
|------|------------|-------------|
| focus | `()` | Emitted when the select gains focus |
| blur | `()` | Emitted when the select loses focus |
| focusInput | `()` | Emitted when the input inside the select gains focus |
| blurInput | `()` | Emitted when the input inside the select loses focus |

> **Note:** All other events are inherited from Naive UI's Select component.

## Slots

| Name | Props | Description |
|------|-------|-------------|
| header | - | Custom header content |
| action | - | Custom action content |
| empty | - | Custom empty content |
| arrow | - | Custom arrow icon |

## Methods

| Name | Parameters | Description |
|------|------------|-------------|
| focus | `()` | Focus the select |
| blur | `()` | Blur the select |
| focusInput | `()` | Focus the input inside the select |
| blurInput | `()` | Blur the input inside the select |

## Usage

```vue
<script setup lang="ts">
import type { FormInst, FormRules, FormSize, SelectOption, SelectRenderLabel, SelectRenderTag } from 'naive-ui'
import type { Value } from 'naive-ui/es/select/src/interface'
import type { FloatSelectInst } from '@/components/ui/float-select/types'
import { NAvatar, NTag, NText } from 'naive-ui'
import { reactive, ref } from 'vue'

// Basic
const opts: SelectOption[] = [
  {
    label: 'Everybody\'s Got Something to Hide Except Me and My Monkey',
    value: 'song0',
    disabled: true
  },
  {
    label: 'Drive My Car',
    value: 'song1'
  },
  {
    label: 'Norwegian Wood',
    value: 'song2'
  }
]
const basicVal = ref<string>('')

// Events
const eventVal = ref(null)

// Methods
const selectRef = ref<FloatSelectInst>()

// Form
const size = ref<FormSize>('medium')
const disabled = ref<boolean>(false)
const formRef = ref<FormInst>()
const formData: { val: Value | null } = reactive({
  val: null
})
const formRules: FormRules = {
  val: [
    { required: true, message: 'Please select val', trigger: 'blur' }
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

// Other
const show = ref<boolean>(false)
const emptyOpts: SelectOption[] = []

const renderLabel: SelectRenderLabel = (option) => {
  return h(
    'div',
    {
      style: {
        display: 'flex',
        alignItems: 'center'
      }
    },
    [
      h(NAvatar, {
        round: true,
        size: 'small'
      }, { default: () => 'Tag' }),
      h(
        'div',
        {
          style: {
            marginLeft: '12px',
            padding: '4px 0'
          }
        },
        [
          h('div', null, [option.label as string]),
          h(
            NText,
            { depth: 3, tag: 'div' },
            {
              default: () => 'description'
            }
          )
        ]
      )
    ]
  )
}

const renderMultipleSelectTag: SelectRenderTag = ({ option, handleClose }) => {
  return h(
    NTag,
    {
      style: {
        padding: '0 6px 0 4px'
      },
      round: true,
      closable: true,
      onClose: (e) => {
        e.stopPropagation()
        handleClose()
      }
    },
    {
      default: () =>
        h(
          'div',
          {
            style: {
              display: 'flex',
              alignItems: 'center'
            }
          },
          [
            h(NAvatar, {
              round: true,
              style: {
                marginRight: '4px'
              }
            }, { default: () => 'Tag' }),
            option.label as string
          ]
        )
    }
  )
}

const options = [
  {
    label: 'Tag1',
    value: 'tag1'
  },
  {
    label: 'Tag2',
    value: 'tag2'
  },
  {
    label: 'Tag3',
    value: 'tag3'
  }
]
</script>

<template>
  <!-- Basic -->
  <FloatSelect v-model:value="basicVal" :options="opts" placeholder="Please Select" />

  <!-- Sizes -->
  <NFlex vertical>
    <FloatSelect :options="opts" size="tiny" placeholder="Tiny Select" />
    <FloatSelect :options="opts" size="small" placeholder="Small Select" />
    <FloatSelect :options="opts" placeholder="Medium Select" />
    <FloatSelect :options="opts" size="large" placeholder="Large Select" />
  </NFlex>

  <!-- Status -->
  <NFlex vertical>
    <FloatSelect :options="opts" loading placeholder="Loading" />
    <FloatSelect :options="opts" disabled placeholder="Disabled" />
    <FloatSelect :options="opts" status="warning" placeholder="Warning" />
    <FloatSelect :options="opts" status="error" placeholder="Error" />
  </NFlex>

  <!-- Slots -->
  <FloatSelect v-model:show="show" :options="emptyOpts" placeholder="Please Select">
    <template #header>
      <NText strong>Header</NText>
    </template>
    <template #action>
      <div class="flex justify-end">
        <NButton>Action</NButton>
      </div>
    </template>
    <template #empty>
      <NText>Empty</NText>
    </template>
    <template v-if="show" #arrow>
      <div class="i-lucide:lasso-select" />
    </template>
  </FloatSelect>

  <!-- Events -->
  <FloatSelect
    v-model:value="eventVal" placeholder="Please Select" @focus="() => message.info('Focus')"
    @blur="() => message.info('Blur')"
  />

  <!-- Methods -->
  <NFlex>
    <NButton @click="selectRef?.focus()">Focus</NButton>
    <NButton @click="selectRef?.focusInput()">Focus Input</NButton>
    <NButton @click="selectRef?.blur()">Blur</NButton>
    <NButton @click="selectRef?.blurInput()">Blur Input</NButton>
  </NFlex>
  <FloatSelect ref="selectRef" tag multiple :show="false" filterable placeholder="Please Select" />

  <!-- Other - Multiple with Custom Rendering -->
  <FloatSelect
    multiple :options="options" :render-label="renderLabel"
    :render-tag="renderMultipleSelectTag" filterable placeholder="Please Select"
  />

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
    <NFormItem path="val">
      <FloatSelect
        v-model:value="formData.val" clearable class="w-60" :options="opts"
        placeholder="Please Select"
      />
    </NFormItem>
    <NFormItem>
      <RippleButton @click="handleSubmit">
        Submit
      </RippleButton>
    </NFormItem>
  </NForm>
</template>
```