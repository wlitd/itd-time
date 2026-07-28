<script setup lang="ts">
import type { SelectOption, SelectRenderLabel, SelectRenderTag } from 'naive-ui'
import { NAvatar, NText } from 'naive-ui'
import { useI18n } from 'vue-i18n'

import squidwardPng from '@/assets/images/squidward.png'

const { t } = useI18n()

const { remindStyle } = storeToRefs(useOffWorkStore())

const remindStyleOpts: SelectOption[] = reactive([
  { label: t('squidward'), value: 'squidward', description: t('squidwardDescription'), avatar: squidwardPng }
])

const renderSingleSelectTag: SelectRenderTag = ({ option }) => {
  return h(
    'div',
    {
      style: {
        display: 'flex',
        alignItems: 'center',
        gap: '8px'
      }
    },
    [
      h(NAvatar, {
        src: option.avatar as string,
        round: true,
        size: 24
      }),
      option.label as string
    ]
  )
}

const renderLabel: SelectRenderLabel = (option) => {
  return h(
    'div',
    {
      style: {
        display: 'flex',
        alignItems: 'center',
        gap: '8px'
      }
    },
    [
      h(NAvatar, {
        src: option.avatar as string,
        round: true,
        size: 'small'
      }),
      h(
        'div',
        {
          style: {
            display: 'flex',
            flexDirection: 'column',
            gap: '4px',
            padding: '4px 0'
          }
        },
        [
          h('div', null, [option.label as string]),
          h(
            NText,
            { depth: 3, tag: 'div' },
            {
              default: () => option.description as string
            }
          )
        ]
      )
    ]
  )
}
</script>

<template>
  <div class="flex w-full">
    <FloatSelect
      v-model:value="remindStyle" :placeholder="t('remindStyle')" :options="remindStyleOpts"
      :render-tag="renderSingleSelectTag" :render-label="renderLabel"
    />
  </div>
</template>

<i18n lang="json">
{
    "zh": {
        "remindStyle": "提醒样式",
        "squidward": "我下班了，蟹老板！",
        "squidwardDescription": "你我皆是章鱼哥"
    },
    "en": {
        "remindStyle": "Remind Style",
        "squidward": "I'm off work, Mr. Krabs!",
        "squidwardDescription": "You and I are both Squidward"
    }
}
</i18n>
