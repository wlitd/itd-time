import dayjs from 'dayjs'

/**
 * TODO 项类型
 */
export interface TodoItem {
  id: string
  /** TODO 标题 */
  title: string
  /** TODO 截至时间 */
  deadline: number
  /** TODO 是否完成 */
  completed: boolean
}

/**
 * 显示的 TODO 项类型
 */
export interface DisplayTodo extends TodoItem {
  /** TODO 是否过期 */
  isOverdue: boolean
}

/**
 * TODO 状态管理 Store
 * 用于管理应用的 TODO 列表设置
 */
export const useTodoStore = defineStore('todo', () => {
  /** TODO 列表 */
  const todos = ref<TodoItem[]>([])

  /**
   * 添加 TODO
   * @param todo TODO 项
   */
  function addTodo(todo: Pick<TodoItem, 'title' | 'deadline'>): void {
    todos.value.push({ id: Date.now().toString(), ...todo, completed: false })
  }

  /**
   * 更新 TODO
   * @param todo TODO 项
   */
  function updateTodo(todo: TodoItem): void {
    const index = todos.value.findIndex(t => t.id === todo.id)
    if (index !== -1) {
      todos.value[index] = todo
    }
  }

  /**
   * 切换 TODO 完成状态
   * @param id TODO ID
   */
  function toggleTodo(id: string): void {
    const todo = todos.value.find(t => t.id === id)
    if (todo)
      todo.completed = !todo.completed
  }

  /**
   * 删除 TODO
   * @param id TODO ID
   */
  function removeTodo(id: string): void {
    todos.value = todos.value.filter(t => t.id !== id)
  }

  /**
   * 获取排序后的 TODO 列表
   */
  const now = useNow({ interval: 1000 })
  const sortedTodos = computed<DisplayTodo[]>(() => {
    return [...todos.value]
      .map(todo => ({
        ...todo,
        isOverdue: !todo.completed && dayjs(todo.deadline).isBefore(dayjs(now.value))
      }))
      .sort((a, b) => {
        if (a.completed !== b.completed) {
          return a.completed ? 1 : -1
        }
        return a.deadline - b.deadline
      })
  })

  return { todos, addTodo, updateTodo, toggleTodo, removeTodo, sortedTodos }
}, {
  /** 持久化配置 */
  persist: {
    /** 持久化存储的键名 */
    key: 'todo',
    /** 使用 localStorage 进行持久化存储 */
    storage: localStorage,
    /** 需要持久化的状态字段 */
    pick: ['todos']
  }
})
