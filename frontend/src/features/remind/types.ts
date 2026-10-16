import type SquidwardRemind from './SquidwardRemind.vue'

/** 提醒卡片通用入参：本地化后的标题与描述。 */
export interface RemindCardProps {
  title: string
  description: string
}

export interface RemindCardEmits {
  confirm: []
  cancel: []
}

export type SquidwardRemindInst = InstanceType<typeof SquidwardRemind>
