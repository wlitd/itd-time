import { IsRemindWindow, ShowRemind } from '../../wailsjs/go/main/App'

/**
 * TODO 待办提醒 composable
 * 每秒轮询所有未完成 TODO 的截止时间，在「截止时间 - advanceMinutes」时刻触发提醒弹窗
 * 提醒标记持久化在各 TODO 项的 reminded 字段中，重启应用不会重复弹窗
 * ⚠️ 仅在主窗口中调用，提醒弹窗子进程不参与轮询（避免无限套娃）
 */
export function useTodoRemind() {
  const todoStore = useTodoStore()
  const { advanceMinutes, remindStyle } = storeToRefs(useOffWorkStore())

  const now = useNow({ interval: 1000 })

  /** 是否为提醒弹窗子进程（子进程不触发新提醒） */
  const isRemindProcess = ref<boolean>(true)

  onMounted(async () => {
    if (isWails())
      isRemindProcess.value = await IsRemindWindow()
    else
      isRemindProcess.value = false
  })

  // 每秒检查是否有 TODO 到达提醒时间
  watch(now, () => {
    if (!isWails() || isRemindProcess.value)
      return

    const advanceMs = advanceMinutes.value * 60 * 1000
    const style = remindStyle.value
    const nowMs = now.value.getTime()

    for (const todo of todoStore.todos) {
      if (todo.completed || todo.reminded)
        continue

      const remindAt = todo.deadline - advanceMs
      if (nowMs < remindAt)
        continue

      todoStore.markReminded(todo.id)
      // 每 tick 仅触发一个提醒弹窗，多个待办会依次在下个 tick 弹出（间隔 ≥1s）
      ShowRemind(style, 'todoAdvance', todo.title)
      return
    }
  })
}
