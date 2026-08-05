import type { RemindMode } from '@/stores/offWork'
import { ShowRemind, SyncRemindConfig } from '../../wailsjs/go/main/App'

/**
 * 下班提醒 composable
 * 负责将提醒配置同步到 Go 端定时器，并提供手动弹出提醒的能力
 * ⚠️ 仅在主窗口中使用，提醒时机由 Go 端定时器判断（不受 WebView 后台节流影响）
 */
export function useOffWorkRemind() {
  const { offTime, advanceMinutes, repeatDays, remindStyle } = storeToRefs(useOffWorkStore())

  // 配置变更时同步到 Go 端
  watch(
    [offTime, advanceMinutes, repeatDays, remindStyle],
    () => {
      if (!isWails())
        return

      SyncRemindConfig(offTime.value, advanceMinutes.value || 0, [...repeatDays.value], remindStyle.value)
    },
    { immediate: true, deep: true }
  )

  /**
   * 弹出提醒窗口（独立进程，显示在屏幕右下角）
   * @param mode 提醒模式
   */
  function showRemind(mode: RemindMode): void {
    if (!isWails())
      return

    ShowRemind(remindStyle.value, mode, '')
  }

  return { showRemind }
}
