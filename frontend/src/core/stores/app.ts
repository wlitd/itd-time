import { defineStore } from 'pinia'
import { SettingsService, SystemService, TodoService } from '../../core/bindings'
import type {
  OffWorkConfig,
  RemindSettings,
  ThemePrefs,
  Todo,
  UpdatePrefs,
} from '../../core/bindings'
import { RemindStyle, Weekday } from '../../core/bindings'
import { setAppLocale } from '../../core/i18n'
import type { AppLocale } from '../../core/i18n/messages'
import { isWails } from '../../shared/env'
import { floorToMinute } from '../../shared/datetime'

/** 展示用代办：在领域模型上补派生的逾期标记。 */
export interface DisplayTodo extends Todo {
  isOverdue: boolean
}

/** 兜底默认值：仅用于纯浏览器预览（无 Go 后端时）。 */
function browserDefaults(): {
  offWork: OffWorkConfig
  remind: RemindSettings
  theme: ThemePrefs
  update: UpdatePrefs
} {
  return {
    offWork: {
      enabled: true,
      offTime: '17:30',
      repeatDays: [Weekday.Monday, Weekday.Tuesday, Weekday.Wednesday, Weekday.Thursday, Weekday.Friday],
    } as OffWorkConfig,
    remind: { advanceMinutes: 10, style: RemindStyle.RemindStyleSquidward } as RemindSettings,
    theme: { dark: false, locale: 'zh-CN' } as ThemePrefs,
    update: { checkOnLaunch: true } as UpdatePrefs,
  }
}

export const useAppStore = defineStore('app', () => {
  const ready = ref<boolean>(false)
  const offWork = ref<OffWorkConfig>(browserDefaults().offWork)
  const remind = ref<RemindSettings>(browserDefaults().remind)
  const theme = ref<ThemePrefs>(browserDefaults().theme)
  const update = ref<UpdatePrefs>(browserDefaults().update)
  const todos = ref<Todo[]>([])

  // ==================== 派生 ====================

  /** 排序 + 逾期标记：未完成在前，按截止时间升序。 */
  const sortedTodos = computed<DisplayTodo[]>(() => {
    const nowMs = Date.now()
    return todos.value
      .map(t => ({ ...t, isOverdue: !t.completed && t.deadline < nowMs }))
      .sort((a, b) => {
        if (a.completed !== b.completed)
          return a.completed ? 1 : -1
        return a.deadline - b.deadline
      })
  })

  // ==================== 主题 ====================

  function applyTheme(): void {
    document.documentElement.classList.toggle('dark', theme.value.dark)
    setAppLocale((theme.value.locale || 'zh-CN') as AppLocale)
  }

  // ==================== 水合 ====================

  async function hydrate(): Promise<void> {
    if (!isWails()) {
      applyTheme()
      ready.value = true
      return
    }
    const data = await SettingsService.GetAppData()
    offWork.value = data.offWork
    remind.value = data.remind
    theme.value = data.theme
    // 旧版持久化 JSON 无 update 字段时回退默认值（绑定层为 undefined）。
    update.value = data.update ?? browserDefaults().update
    todos.value = data.todos ?? []
    applyTheme()
    ready.value = true
  }

  // ==================== 下班配置 ====================

  async function saveOffWork(next: OffWorkConfig): Promise<void> {
    offWork.value = next
    if (isWails())
      await SettingsService.SaveOffWorkConfig(next)
  }

  async function setOffTime(offTime: string): Promise<void> {
    await saveOffWork({ ...offWork.value, offTime })
  }

  async function toggleRepeatDay(day: number): Promise<void> {
    const set = new Set(offWork.value.repeatDays ?? [])
    if (set.has(day))
      set.delete(day)
    else
      set.add(day)
    const repeatDays = [...set].sort((a, b) => a - b)
    await saveOffWork({ ...offWork.value, repeatDays })
  }

  // ==================== 提醒偏好 ====================

  async function saveRemind(next: RemindSettings): Promise<void> {
    remind.value = next
    if (isWails())
      await SettingsService.SaveRemindSettings(next)
  }

  async function setAdvanceMinutes(minutes: number): Promise<void> {
    await saveRemind({ ...remind.value, advanceMinutes: minutes })
  }

  async function setRemindStyle(style: RemindStyle): Promise<void> {
    await saveRemind({ ...remind.value, style })
  }

  // ==================== 主题 / 语言持久化 ====================

  async function toggleDark(): Promise<void> {
    const next: ThemePrefs = { ...theme.value, dark: !theme.value.dark }
    theme.value = next
    applyTheme()
    if (isWails())
      await SettingsService.SaveTheme(next)
  }

  async function setLocale(locale: AppLocale): Promise<void> {
    const next: ThemePrefs = { ...theme.value, locale }
    theme.value = next
    applyTheme()
    if (isWails())
      await SettingsService.SaveTheme(next)
  }

  // ==================== 更新偏好 ====================

  async function saveUpdatePrefs(next: UpdatePrefs): Promise<void> {
    update.value = next
    if (isWails())
      await SettingsService.SaveUpdatePrefs(next)
  }

  async function setCheckOnLaunch(checkOnLaunch: boolean): Promise<void> {
    await saveUpdatePrefs({ ...update.value, checkOnLaunch })
  }

  // ==================== 代办 CRUD ====================

  async function refreshTodos(): Promise<void> {
    if (isWails())
      todos.value = (await TodoService.List()) ?? []
  }

  async function addTodo(title: string, deadline: number): Promise<void> {
    if (!isWails())
      return
    await TodoService.Add({ title, deadline: floorToMinute(deadline) })
    await refreshTodos()
  }

  async function updateTodo(id: string, title: string, deadline: number): Promise<void> {
    if (!isWails())
      return
    await TodoService.Update(id, { title, deadline: floorToMinute(deadline) })
    await refreshTodos()
  }

  async function setCompleted(id: string, completed: boolean): Promise<void> {
    if (!isWails())
      return
    await TodoService.SetCompleted(id, completed)
    await refreshTodos()
  }

  async function removeTodo(id: string): Promise<void> {
    if (!isWails())
      return
    await TodoService.Delete(id)
    await refreshTodos()
  }

  // ==================== 系统集成 ====================

  const autostart = ref<boolean>(false)

  async function loadAutostart(): Promise<void> {
    if (isWails())
      autostart.value = await SystemService.IsAutostartEnabled()
  }

  async function setAutostart(enabled: boolean): Promise<void> {
    if (!isWails()) {
      autostart.value = enabled
      return
    }
    await SystemService.SetAutostart(enabled)
    autostart.value = enabled
  }

  /** 关窗即隐藏：隐藏主窗口到托盘。 */
  async function hideToTray(): Promise<void> {
    if (isWails())
      await SystemService.HideMainWindow()
  }

  /** 退出应用。 */
  async function quit(): Promise<void> {
    if (isWails())
      await SystemService.Quit()
  }

  return {
    ready,
    offWork,
    remind,
    theme,
    update,
    todos,
    autostart,
    sortedTodos,
    hydrate,
    saveOffWork,
    setOffTime,
    toggleRepeatDay,
    saveRemind,
    setAdvanceMinutes,
    setRemindStyle,
    toggleDark,
    setLocale,
    saveUpdatePrefs,
    setCheckOnLaunch,
    refreshTodos,
    addTodo,
    updateTodo,
    setCompleted,
    removeTodo,
    loadAutostart,
    setAutostart,
    hideToTray,
    quit,
  }
})
