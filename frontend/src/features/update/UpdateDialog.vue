<script setup lang="ts">
import { computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import Dialog from 'primevue/dialog'
import Button from 'primevue/button'
import ProgressBar from 'primevue/progressbar'
import { useToast } from 'primevue/usetoast'
import { useUpdate } from './useUpdate'

const { t } = useI18n()
const { state, isBusy, check, download, install, closeDialog, consumeUpToDate } = useUpdate()
const toast = useToast()

// 无更新时不弹窗：监听阶段，转成轻量的 Toast 提示。
watch(() => state.phase, (phase) => {
  if (phase !== 'upToDate')
    return
  const version = state.info?.currentVersion
  toast.add({
    severity: 'success',
    summary: t('upToDate'),
    detail: version ? t('upToDateDetail', { version }) : t('upToDateDetailFallback'),
    life: 3000,
  })
  consumeUpToDate()
})

const visible = computed<boolean>({
  get: () => state.dialogVisible,
  set: (v: boolean): void => {
    if (!v)
      closeDialog()
  },
})

const phaseText = computed<string>(() => {
  switch (state.phase) {
    case 'checking':
      return t('checking')
    case 'downloading':
      return t('downloadingDesc', { percent: state.percent })
    case 'ready':
      return t('installed')
    case 'error':
      return state.error === 'web' ? t('isWeb') : t('checkFailed')
    default:
      return state.info ? t('newVersionFound', { version: state.info.version }) : ''
  }
})
</script>

<template>
  <Dialog
    v-model:visible="visible"
    modal
    :style="{ width: 'min(24rem, 92vw)' }"
    :closable="state.phase !== 'downloading'"
    :header="t('checkUpdate')"
  >
    <div class="flex flex-col gap-4 py-2">
      <!-- 检查中 -->
      <div v-if="state.phase === 'checking'" class="flex-center py-6">
        <i class="i-lucide-loader-circle animate-spin text-2xl text-primary" />
        <span class="ml-3 text-sm text-[var(--app-muted)]">{{ t('checking') }}</span>
      </div>

      <!-- 出错 -->
      <div v-else-if="state.phase === 'error'" class="flex items-start gap-3">
        <i class="i-lucide-triangle-alert text-xl text-[var(--p-red-500)]" />
        <span class="text-sm">{{ phaseText }}</span>
      </div>

      <!-- 下载 / 就绪 -->
      <div v-else-if="state.phase === 'downloading' || state.phase === 'ready'" class="flex flex-col gap-3">
        <!-- 确定态进度条：隐藏条上的百分比（showValue=false），下方文字仍展示进度 -->
        <ProgressBar :value="state.percent" :show-value="false" />
        <p class="text-sm text-[var(--app-muted)]">{{ phaseText }}</p>
        <p v-if="state.phase === 'downloading'" class="text-xs text-[var(--app-muted)]">
          {{ t('downloadTip') }}
        </p>
      </div>

      <!-- 发现新版本 -->
      <div v-else class="flex flex-col gap-3">
        <div class="flex items-center gap-3">
          <i class="i-lucide-sparkles text-xl text-primary" />
          <span class="text-base font-semibold">{{ phaseText }}</span>
        </div>
        <div v-if="state.info?.notes" class="flex flex-col gap-1">
          <span class="text-xs uppercase tracking-wide text-[var(--app-muted)]">{{ t('releaseNotes') }}</span>
          <p class="whitespace-pre-wrap text-sm text-[var(--app-text)]">{{ state.info.notes }}</p>
        </div>
      </div>
    </div>

    <template #footer>
      <div class="flex w-full items-center justify-end gap-2">
        <Button
          v-if="state.phase === 'available'"
          :label="t('download')"
          :loading="isBusy"
          @click="download"
        >
          <template #icon>
            <i class="i-lucide-download" />
          </template>
        </Button>
        <Button
          v-if="state.phase === 'ready'"
          :label="t('installRestart')"
          severity="success"
          @click="install"
        >
          <template #icon>
            <i class="i-lucide-rotate-cw" />
          </template>
        </Button>
        <Button
          v-if="state.phase === 'error'"
          :label="t('checkUpdate')"
          @click="check()"
        >
          <template #icon>
            <i class="i-lucide-refresh-cw" />
          </template>
        </Button>
      </div>
    </template>
  </Dialog>
</template>

<i18n lang="yaml">
zh-CN:
  checking: 检查中…
  downloadingDesc: "正在下载更新 {percent}%…"
  installed: 已就绪，等待安装
  isWeb: 浏览器环境不支持检查更新！
  checkFailed: 检查更新失败！
  upToDate: 已是最新版本
  upToDateDetail: "当前版本 v{version}，无需更新"
  upToDateDetailFallback: 当前已是最新版本
  newVersionFound: "发现新版本 {version}"
  checkUpdate: 检查更新
  downloadTip: 下载完成后点击「安装并重启」完成更新
  releaseNotes: 更新说明
  download: 下载
  installRestart: 安装并重启
en-US:
  checking: Checking…
  downloadingDesc: "Downloading update {percent}%…"
  installed: Ready to install
  isWeb: Browser environment does not support checking for updates!
  checkFailed: Update check failed!
  upToDate: You're up to date
  upToDateDetail: "Version v{version} is the latest"
  upToDateDetailFallback: You're running the latest version
  newVersionFound: "New version {version} found"
  checkUpdate: Check for Updates
  downloadTip: 'Click "Install & Restart" once download completes'
  releaseNotes: Release Notes
  download: Download
  installRestart: Install & Restart
</i18n>
