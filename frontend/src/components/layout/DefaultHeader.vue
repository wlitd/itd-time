<script setup lang="tsx">
import type { DropdownOption } from 'naive-ui'
import { NEllipsis, NSwitch, useThemeVars } from 'naive-ui'
import { useI18n } from 'vue-i18n'
import { Quit } from '@/../wailsjs/runtime/runtime'

const { t } = useI18n()

const themeVars = useThemeVars()

const { currentLocaleKey } = storeToRefs(useLocaleStore())
const { setLocaleKey } = useLocaleStore()

const opts = computed<DropdownOption[]>(() => ([
  {
    key: 'locale',
    label: t('languageSetting'),
    children: [
      {
        key: 'zh',
        type: 'render',
        render: () => (
          <div class="px-1">
            <div class={`rounded w-24 h-34px flex items-center cursor-pointer px-2 justify-between gap-2 ${currentLocaleKey.value === 'zh' ? 'bg-[var(--primary-color)]' : 'hover:bg-[var(--itd-hover)]'}`}>
              <div class="i-flag:cn-4x3" />
              <div>简体中文</div>
            </div>
          </div>
        ),
        props: {
          onClick: () => {
            if (currentLocaleKey.value !== 'zh')
              setLocaleKey('zh')
          }
        }
      },
      {
        key: 'en',
        type: 'render',
        render: () => (
          <div class="px-1">
            <div class={`rounded w-24 h-34px flex items-center cursor-pointer px-2 justify-between gap-2 ${currentLocaleKey.value === 'en' ? 'bg-[var(--primary-color)]' : 'hover:bg-[var(--itd-hover)]'}`}>
              <div class="i-flag:us-4x3" />
              <div>English</div>
            </div>
          </div>
        ),
        props: {
          onClick: () => {
            if (currentLocaleKey.value !== 'en')
              setLocaleKey('en')
          }
        }
      }
    ]
  },
  {
    key: 'divider-1',
    type: 'divider'
  },
  {
    key: 'selfStart',
    type: 'render',
    render: () => (
      <div class="px-1">
        <div class="rounded h-34px flex items-center cursor-pointer px-2 justify-between gap-2 hover:bg-[var(--itd-hover)]">
          <div>{t('selfStart')}</div>
          <NSwitch round={false} />
        </div>
      </div>
    )
  },
  {
    key: 'autoCheckUpdate',
    type: 'render',
    render: () => (
      <div class="px-1">
        <div class="rounded h-34px flex items-center cursor-pointer px-2 justify-between gap-2 hover:bg-[var(--itd-hover)]">
          <NEllipsis class="flex-1">{t('autoCheckUpdate')}</NEllipsis>
          <NSwitch round={false} />
        </div>
      </div>
    )
  },
  {
    key: 'checkUpdate',
    label: t('checkUpdate')
  },
  {
    key: 'divider-2',
    type: 'divider'
  },
  {
    key: 'theme',
    type: 'render',
    render: () => (
      <div class="px-1">
        <div class="rounded h-34px flex items-center cursor-pointer px-2 justify-between gap-2 hover:bg-[var(--itd-hover)]">
          <div>{t('theme')}</div>
          <ThemeToggle />
        </div>
      </div>
    )
  },
  {
    key: 'themeColor',
    type: 'render',
    render: () => (
      <div class="px-1">
        <div class="rounded h-34px flex items-center cursor-pointer px-2 justify-between gap-2 hover:bg-[var(--itd-hover)]">
          <div>{t('themeColor')}</div>
          <ColorPicker />
        </div>
      </div>
    )
  },
  {
    key: 'divider-3',
    type: 'divider'
  },
  {
    key: 'close',
    label: t('close'),
    props: {
      onClick: () => {
        Quit()
      }
    }
  }
]))
</script>

<template>
  <div class="flex flex-col">
    <!-- 拖拽区域 -->
    <div class="h-6" style="--wails-draggable: drag;" />

    <div class="flex flex-col gap-3">
      <div class="flex items-center gap-2 justify-between">
        <NH2 prefix="bar">
          <NText strong italic>{{ t('off') }}</NText>
        </NH2>

        <NDropdown
          trigger="click" placement="bottom-end" :style="{
            'width': '180px',
            '--itd-hover': themeVars.hoverColor,
            'color': themeVars.textColor2,
            'border-radius': '12px',
          }" :options="opts"
        >
          <NButton v-ripple circle type="primary">
            <template #icon>
              <div class="i-lucide:settings" />
            </template>
          </NButton>
        </NDropdown>
      </div>

      <CountdownCard />
    </div>
  </div>
</template>

<i18n lang="json">
{
  "zh": {
    "off": "我免费啦！",
    "languageSetting": "语言设置",
    "selfStart": "开机自启",
    "autoCheckUpdate": "启动时检查更新",
    "checkUpdate": "检查更新",
    "theme": "应用主题",
    "themeColor": "主题色",
    "close": "关闭"
  },
  "en": {
    "off": "I’m free!",
    "languageSetting": "Language Setting",
    "selfStart": "Self Start",
    "autoCheckUpdate": " Check For Updates At Startup",
    "checkUpdate": "Check For Updates",
    "theme": "Theme",
    "themeColor": "Theme Color",
    "close": "Close"
  }
}
</i18n>
