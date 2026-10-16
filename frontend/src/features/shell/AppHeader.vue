<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import Button from 'primevue/button'
import Popover from 'primevue/popover'
import Divider from 'primevue/divider'
import ToggleSwitch from 'primevue/toggleswitch'
import Tag from 'primevue/tag'
import { useAppStore } from '../../core/stores/app'
import { UpdateService } from '../../core/bindings'
import { isWails } from '../../shared/env'
import type { AppLocale } from '../../core/i18n/messages'
import { useUpdate } from '../update/useUpdate'
import ThemeToggle from './ThemeToggle.vue'

const store = useAppStore()
const { locale, t } = useI18n()
const { check } = useUpdate()

const langPopover = ref<InstanceType<typeof Popover>>()
const morePopover = ref<InstanceType<typeof Popover>>()

const version = ref<string>('')

const autostart = computed<boolean>({
  get: () => store.autostart,
  set: (v: boolean): Promise<void> => store.setAutostart(v),
})

interface LanguageMeta {
  code: AppLocale
  label: string
  flag: string
}

const languages: LanguageMeta[] = [
  { code: 'zh-CN', label: '简体中文', flag: 'i-flag:cn-4x3' },
  { code: 'en-US', label: 'English', flag: 'i-flag:us-4x3' },
]

// 菜单项统一样式：透明底、主色轻 hover、图标左对齐
const itemCls = 'flex w-full cursor-pointer items-center gap-2.5 rounded-xl border-0 bg-transparent px-2.5 py-2 text-left text-[14px] text-[var(--app-text)] transition-colors duration-150 hover:bg-primary/10'

onMounted(async (): Promise<void> => {
  if (!isWails())
    return
  await store.loadAutostart()
  version.value = await UpdateService.CurrentVersion()
})

function toggleLang(event: Event): void {
  langPopover.value?.toggle(event)
}

function selectLang(code: AppLocale): void {
  void store.setLocale(code)
  langPopover.value?.hide()
}

function toggleMore(event: Event): void {
  morePopover.value?.toggle(event)
}

async function onCheckUpdate(): Promise<void> {
  morePopover.value?.hide()
  await check()
}

async function onHide(): Promise<void> {
  morePopover.value?.hide()
  await store.hideToTray()
}

async function onQuit(): Promise<void> {
  morePopover.value?.hide()
  await store.quit()
}
</script>

<template>
  <header class="wails-drag flex select-none items-center justify-between px-3 pb-2 pt-3">
    <div class="flex items-center gap-2">
      <span class="flex-center h-8 w-8 rounded-lg bg-primary/15 text-primary">
        <i class="i-lucide-fish text-lg" />
      </span>
      <span class="text-[15px] font-semibold tracking-tight">{{ t('title') }}</span>
    </div>

    <div class="wails-no-drag flex items-center gap-1 rounded-full bg-[var(--app-surface)] p-1 shadow-[var(--app-shadow-sm)] ring-1 ring-[var(--app-border)]">
      <!-- 语言 -->
      <Button text rounded severity="secondary" @click="toggleLang">
        <template #icon>
          <i class="i-lucide-languages text-lg" />
        </template>
      </Button>

      <ThemeToggle />

      <!-- 更多 -->
      <Button text rounded severity="secondary" @click="toggleMore">
        <template #icon>
          <i class="i-lucide-ellipsis text-lg" />
        </template>
      </Button>
    </div>

    <!-- 语言面板：与「更多」面板同款样式 -->
    <Popover ref="langPopover">
      <div class="flex w-44 flex-col gap-0.5 p-1.5">
        <button
          v-for="lang in languages"
          :key="lang.code"
          :class="[itemCls, locale === lang.code ? '!text-primary' : '']"
          @click="selectLang(lang.code)"
        >
          <span class="h-[14px] w-[19px] shrink-0 overflow-hidden rounded-[3px] ring-1 ring-black/10">
            <i :class="[lang.flag, 'block h-full w-full']" />
          </span>
          <span class="flex-1">{{ lang.label }}</span>
          <i v-if="locale === lang.code" class="i-lucide-check shrink-0 text-base text-primary" />
        </button>
      </div>
    </Popover>

    <!-- 更多面板 -->
    <Popover ref="morePopover">
      <div class="flex w-60 flex-col gap-0.5 p-1.5">
        <!-- 版本信息 -->
        <div class="flex items-center justify-between px-2.5 py-2">
          <span class="flex items-center gap-2.5 text-[14px] text-[var(--app-text)]">
            <i class="i-lucide-info text-base text-[var(--app-muted)]" />
            {{ t('version') }}
          </span>
          <Tag v-if="version" :value="`v${version}`" severity="info" class="!text-[11px]" />
        </div>

        <Divider class="!my-1" />

        <!-- 开机自启 -->
        <div class="flex items-center justify-between px-2.5 py-2">
          <span class="flex items-center gap-2.5 text-[14px] text-[var(--app-text)]">
            <i class="i-lucide-rocket text-base text-[var(--app-muted)]" />
            {{ t('selfStart') }}
          </span>
          <ToggleSwitch v-model="autostart" />
        </div>

        <!-- 检查更新 -->
        <button :class="itemCls" @click="onCheckUpdate">
          <i class="i-lucide-cloud-download text-base text-[var(--app-muted)]" />
          {{ t('checkUpdate') }}
        </button>

        <Divider class="!my-1" />

        <!-- 隐藏到托盘 -->
        <button :class="itemCls" @click="onHide">
          <i class="i-lucide-minimize text-base text-[var(--app-muted)]" />
          {{ t('hideToTray') }}
        </button>

        <!-- 退出 -->
        <button :class="[itemCls, '!text-[var(--app-danger)] hover:!bg-danger/10']" @click="onQuit">
          <i class="i-lucide-power text-base" />
          {{ t('quit') }}
        </button>
      </div>
    </Popover>
  </header>
</template>

<i18n lang="yaml">
zh-CN:
  title: 我免费啦！
  version: 版本信息
  selfStart: 开机自启
  checkUpdate: 检查更新
  hideToTray: 隐藏到托盘
  quit: 退出
en-US:
  title: I'm free!
  version: Version Info
  selfStart: Auto Start
  checkUpdate: Check for Updates
  hideToTray: Hide to Tray
  quit: Quit
</i18n>
