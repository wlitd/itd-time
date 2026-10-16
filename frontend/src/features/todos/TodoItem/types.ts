import type TodoItem from './TodoItem.vue'
import type { DisplayTodo } from '../../../core/stores/app'

/** 单条代办展示项。editing 为真时进入行内编辑态。 */
export interface TodoItemProps {
  todo: DisplayTodo
  editing?: boolean
}

/** 行内编辑保存的载荷。 */
export interface TodoSavePayload {
  id: string
  title: string
  deadline: number
}

export interface TodoItemEmits {
  toggle: [id: string, completed: boolean]
  edit: [todo: DisplayTodo]
  remove: [id: string]
  save: [payload: TodoSavePayload]
  cancel: []
}

export type TodoItemInst = InstanceType<typeof TodoItem>
