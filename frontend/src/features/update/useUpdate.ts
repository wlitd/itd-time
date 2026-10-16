import { computed, reactive, readonly } from 'vue'
import { Events, Updater } from '@wailsio/runtime'
import { UpdateService } from '../../core/bindings'
import type { UpdateInfo } from '../../core/bindings'
import { isWails } from '../../shared/env'

/** 更新流程阶段。 */
export type UpdatePhase = 'idle' | 'checking' | 'available' | 'upToDate' | 'downloading' | 'ready' | 'error'

/** 下载进度事件负载（对齐 Go updater 的 Progress）。 */
interface ProgressPayload {
  written: number
  total: number
}

interface UpdateState {
  phase: UpdatePhase
  info: UpdateInfo | null
  percent: number
  error: string
  dialogVisible: boolean
}

// 模块级单例：头部按钮与更新弹窗共享同一份状态与订阅。
const state = reactive<UpdateState>({
  phase: 'idle',
  info: null,
  percent: 0,
  error: '',
  dialogVisible: false,
})

let bound = false

/** 订阅 updater 事件（幂等，仅绑定一次）。 */
function bindEvents(): void {
  if (bound || !isWails())
    return
  bound = true

  Events.On(Updater.Events.DownloadProgress, (e): void => {
    const p = e.data as unknown as ProgressPayload
    state.percent = p.total > 0 ? Math.round((p.written / p.total) * 100) : 0
    if (state.phase !== 'ready')
      state.phase = 'downloading'
  })

  Events.On(Updater.Events.UpdateReady, (): void => {
    state.percent = 100
    state.phase = 'ready'
  })

  Events.On(Updater.Events.Error, (e): void => {
    const info = e.data as unknown as { message?: string }
    state.error = info?.message ?? '更新失败'
    state.phase = 'error'
  })
}

/** 检查更新选项。silent：启动自动检查模式，仅发现新版本时弹窗，无更新/失败均静默不打扰。 */
export interface CheckOptions {
  silent?: boolean
}

/** 主动检查更新。 */
async function check(options: CheckOptions = {}): Promise<void> {
  const silent = options.silent === true
  if (!isWails()) {
    if (silent)
      return
    state.error = 'web'
    state.phase = 'error'
    state.dialogVisible = true
    return
  }
  bindEvents()
  state.phase = 'checking'
  state.error = ''
  state.percent = 0
  state.info = null
  try {
    const info = await UpdateService.CheckForUpdate()
    state.info = info
    if (info.hasUpdate) {
      // 发现新版本：无论是否静默都走弹窗（下载 / 安装 / 跳过）
      state.phase = 'available'
      state.dialogVisible = true
    }
    else if (silent) {
      // 静默模式无更新：直接回空闲，不触发 Toast。
      state.phase = 'idle'
    }
    else {
      // 无更新：不弹窗，交由 UpdateDialog 监听转 Toast 提示
      state.phase = 'upToDate'
    }
  }
  catch (err) {
    if (silent) {
      // 静默模式失败：不打扰用户，回空闲。
      state.phase = 'idle'
      return
    }
    state.error = String(err)
    state.phase = 'error'
    state.dialogVisible = true
  }
}

/** 下载并暂存更新。 */
async function download(): Promise<void> {
  state.phase = 'downloading'
  state.percent = 0
  await UpdateService.DownloadAndInstall()
}

/** 用已暂存版本替换二进制并重启。 */
async function install(): Promise<void> {
  await UpdateService.ApplyUpdate()
}

/** 跳过当前版本，持久化并不再提醒。 */
async function skip(): Promise<void> {
  const version = state.info?.version
  if (version)
    await UpdateService.SkipVersion(version)
  closeDialog()
}

function closeDialog(): void {
  state.dialogVisible = false
  if (state.phase === 'error' || state.phase === 'available' || state.phase === 'ready' || state.phase === 'upToDate')
    state.phase = 'idle'
}

/** 消费“已是最新”态：Toast 提示后回到空闲，供监听方调用。 */
function consumeUpToDate(): void {
  if (state.phase === 'upToDate')
    state.phase = 'idle'
}

function openDialog(): void {
  state.dialogVisible = true
}

export function useUpdate() {
  return {
    state: readonly(state),
    isBusy: computed<boolean>(() => state.phase === 'checking' || state.phase === 'downloading'),
    hasResult: computed<boolean>(() => state.phase !== 'idle' && state.phase !== 'checking'),
    check,
    download,
    install,
    skip,
    openDialog,
    closeDialog,
    consumeUpToDate,
  }
}
