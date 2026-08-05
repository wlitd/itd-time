import { ShowRemind } from '../../wailsjs/go/main/App'

/**
 * TODO 待办提醒 composable
 * 每秒轮询所有未完成 TODO 的截止时间，在「截止时间 - advanceMinutes」时刻触发提醒弹窗
 * ⚠️ 仅在主窗口中调用，每个 TODO 每会话仅提醒一次（remindedIds 去重）
 */
export function useTodoRemind() {
  const todoStore = useTodoStore()
  const { advanceMinutes, remindStyle } = storeToRefs(useOffWorkStore())

  const now = useNow({ interval: 1000 })

  /** 本会话内已提醒过的 TODO ID 集合，避免重复弹窗 */
  const remindedIds = reactive<Set<string>>(new Set())

  // 清理已删除/已完成的 TODO 对应的提醒记录
  watch(
    () => todoStore.todos.map(t => `${t.id}:${t.completed}`),
    (entries) => {
      const completedOrMissing = new Set(
        entries
          .filter(e => e.endsWith(':true') || !todoStore.todos.some(t => `${t.id}:${t.completed}` === e))
          .map(e => e.split(':')[0])
      )
      for (const id of completedOrMissing) {
        remindedIds.delete(id)
      }
    }
  )

  // 每秒检查是否有 TODO 到达提醒时间
  watch(now, () => {
    if (!isWails())
      return

    const advanceMs = advanceMinutes.value * 60 * 1000
    const style = remindStyle.value
    const nowMs = now.value.getTime()

    for (const todo of todoStore.todos) {
      if (todo.completed || remindedIds.has(todo.id))
        continue

      const remindAt = todo.deadline - advanceMs
      if (nowMs < remindAt)
        continue

      remindedIds.add(todo.id)
      // 每 tick 仅触发一个提醒弹窗，多个待办会依次在下个 tick 弹出（间隔 ≥1s）
      ShowRemind(style, 'advance')
      return
    }
  })
}
