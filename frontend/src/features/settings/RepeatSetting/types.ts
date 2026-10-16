import type RepeatSetting from './RepeatSetting.vue'

/** 重复日选择：modelValue 为星期序号数组（0=周日 .. 6=周六，与 Go time.Weekday 对齐）。 */
export interface RepeatSettingProps {
  modelValue: number[]
}

export interface RepeatSettingEmits {
  'update:modelValue': [days: number[]]
}

export type RepeatSettingInst = InstanceType<typeof RepeatSetting>
