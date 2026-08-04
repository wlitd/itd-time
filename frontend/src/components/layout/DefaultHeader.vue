<script setup lang="tsx">
import type { DropdownOption } from 'naive-ui'
import { NEllipsis, NSwitch, useThemeVars } from 'naive-ui'
import { useI18n } from 'vue-i18n'
import {
  CheckForUpdate,
  DisableAutoStart,
  EnableAutoStart,
  GetVersion,
  IsAutoStartEnabled
} from '@/../wailsjs/go/main/App'
import { Quit } from '@/../wailsjs/runtime/runtime'

const { t } = useI18n()
const themeVars = useThemeVars()

const { currentLocaleKey } = storeToRefs(useLocaleStore())
const { setLocaleKey } = useLocaleStore()

// ==================== 状态 ====================

/** 开机自启状态 */
const selfStartEnabled = ref<boolean>(false)
/** 启动时自动检查更新 */
const autoCheckUpdate = useStorage<boolean>('autoCheckUpdate', false)
/** 当前版本号 */
const currentVersion = ref<string>('')
/** 是否正在检查更新 */
const checking = ref<boolean>(false)

// ==================== 初始化 ====================

async function init(): Promise<void> {
  if (isWails()) {
    try {
      selfStartEnabled.value = await IsAutoStartEnabled()
    } catch { /* 忽略 */ }

    try {
      currentVersion.value = await GetVersion()
    } catch { /* 忽略 */ }

    // 如果开启了自动检查更新，启动时检查一次
    if (autoCheckUpdate.value) {
      doCheckUpdate(true)
    }
  }
}

// ==================== 自启动 ====================

async function toggleSelfStart(val: boolean): Promise<void> {
  if (!isWails()) {
    selfStartEnabled.value = val
    return
  }
  try {
    if (val) {
      await EnableAutoStart()
    } else {
      await DisableAutoStart()
    }
    selfStartEnabled.value = val
    message.success(val ? t('selfStartEnabled') : t('selfStartDisabled'))
  } catch {
    message.error(t('selfStartFailed'))
  }
}

// ==================== 更新检查 ====================

interface UpdateResult {
  info: { version: string, downloadUrl: string, releaseNote: string } | null
  hasUpdate: boolean
}

/** 执行更新检查 */
async function doCheckUpdate(silent = false): Promise<void> {
  if (!isWails()) {
    if (!silent)
      message.info(t('isWeb'))
    return
  }

  if (checking.value)
    return
  checking.value = true

  try {
    const result = (await CheckForUpdate()) as unknown as UpdateResult

    if (result && result.hasUpdate && result.info) {
      const info = result.info
      dialog.warning({
        title: t('updateAvailable'),
        content: () => (
          <div class="flex flex-col gap-2">
            <div>{t('newVersionFound', { version: info.version })}</div>
            {info.releaseNote ? <div class="text-xs opacity-70">{info.releaseNote}</div> : null}
          </div>
        ),
        positiveText: t('download'),
        negativeText: t('cancel'),
        onPositiveClick: async () => {
          try {
            const { DownloadAndInstall } = await import('@/../wailsjs/go/main/App')
            await DownloadAndInstall(info.downloadUrl)
          } catch {
            message.error(t('updateFailed'))
          }
        }
      })
    } else if (!silent) {
      message.success(t('upToDate'))
    }
  } catch {
    if (!silent)
      message.error(t('updateFailed'))
  } finally {
    checking.value = false
  }
}

// ==================== 菜单选项 ====================

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
    key: 'version',
    type: 'render',
    render: () => (
      <div class="px-1">
        <div class="rounded h-34px flex items-center cursor-pointer px-2 justify-between gap-2 hover:bg-[var(--itd-hover)]">
          <div>{t('version')}</div>
          <NTag bordered={false} size="small" round type="info">
            { `V${currentVersion.value}` || 'dev' }
          </NTag>
        </div>
      </div>
    )
  },
  {
    key: 'selfStart',
    type: 'render',
    render: () => (
      <div class="px-1">
        <div class="rounded h-34px flex items-center cursor-pointer px-2 justify-between gap-2 hover:bg-[var(--itd-hover)]">
          <div>{t('selfStart')}</div>
          <NSwitch
            round={false}
            value={selfStartEnabled.value}
            onUpdate:value={toggleSelfStart}
          />
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
          <NSwitch
            round={false}
            value={autoCheckUpdate.value}
            onUpdate:value={(v: boolean) => {
              autoCheckUpdate.value = v
            }}
          />
        </div>
      </div>
    )
  },
  {
    key: 'checkUpdate',
    type: 'render',
    render: () => (
      <div class="px-1">
        <div
          class="rounded h-34px flex items-center cursor-pointer px-2 justify-between gap-2 hover:bg-[var(--itd-hover)]"
          onClick={() => doCheckUpdate(false)}
        >
          <div>{t('checkUpdate')}</div>
          {checking.value
            ? <div class="i-lucide:loader-circle animate-spin text-sm" />
            : null}
        </div>
      </div>
    )
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

onMounted(() => {
  init()
})
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
          trigger="click"
          placement="bottom-end"
          :style="{
            'width': '180px',
            '--itd-hover': themeVars.hoverColor,
            'color': themeVars.textColor2,
            'border-radius': '12px',
          }"
          :options="opts"
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
    "version": "当前版本",
    "close": "关闭",
    "isWeb": "浏览器环境不支持检查更新！",
    "selfStartEnabled": "已开启开机自启！",
    "selfStartDisabled": "已关闭开机自启！",
    "selfStartFailed": "自启动设置失败！",
    "updateAvailable": "发现新版本！",
    "upToDate": "已是最新版本！",
    "updateFailed": "检查更新失败！",
    "newVersionFound": "发现新版本 {version}，是否下载？",
    "download": "下载",
    "cancel": "取消"
  },
  "en": {
    "off": "I'm free!",
    "languageSetting": "Language Setting",
    "selfStart": "Auto Start",
    "autoCheckUpdate": "Check Updates at Startup",
    "checkUpdate": "Check for Updates",
    "theme": "Theme",
    "themeColor": "Theme Color",
    "version": "Version",
    "close": "Close",
    "isWeb": "Browser environment does not support checking for updates!",
    "selfStartEnabled": "Auto start enabled!",
    "selfStartDisabled": "Auto start disabled!",
    "selfStartFailed": "Auto start setting failed!",
    "updateAvailable": "Update Available!",
    "upToDate": "Up to date!",
    "updateFailed": "Update check failed!",
    "newVersionFound": "New version {version} found. Download?",
    "download": "Download",
    "cancel": "Cancel"
  }
}
</i18n>
