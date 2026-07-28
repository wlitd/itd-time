export type RemindStyle = 'squidward'

/** 提醒模式：advance=提前提醒 offwork=下班提醒 preview=预览 */
export type RemindMode = 'advance' | 'offwork' | 'preview'

/**
 * 下班倒计时状态管理 Store
 * 用于管理应用的下班倒计时设置
 */
export const useOffWorkStore = defineStore('offWork', () => {
  /** 下班时间 */
  const offTime = ref<string>('18:00')
  /** 提前分钟数 */
  const advanceMinutes = ref<number>(10)
  /** 重复的天数 */
  const repeatDays = ref<number[]>([1, 2, 3, 4, 5])
  /** 提醒样式 */
  const remindStyle = ref<RemindStyle>('squidward')

  /**
   * 设置下班时间
   * @param _offTime 下班时间
   */
  function setOffWorkTime(_offTime: string): void {
    offTime.value = _offTime
  }

  /**
   * 设置提前分钟数
   * @param _advanceMinutes 提前分钟数
   */
  function setAdvanceMinutes(_advanceMinutes: number): void {
    advanceMinutes.value = _advanceMinutes
  }

  /**
   * 切换重复的天数
   * @param day 天数
   */
  function toggleRepeatDay(day: number): void {
    const index = repeatDays.value.indexOf(day)
    if (index > -1) {
      repeatDays.value.splice(index, 1)
    } else {
      repeatDays.value.push(day)
    }
  }

  /**
   * 设置提醒样式
   * @param _remindStyle 提醒样式
   */
  function setRemindStyle(_remindStyle: RemindStyle): void {
    remindStyle.value = _remindStyle
  }

  return { offTime, setOffWorkTime, advanceMinutes, setAdvanceMinutes, repeatDays, toggleRepeatDay, remindStyle, setRemindStyle }
}, {
  /** 持久化配置 */
  persist: {
    /** 持久化存储的键名 */
    key: 'offWork',
    /** 使用 localStorage 进行持久化存储 */
    storage: localStorage,
    /** 需要持久化的状态字段 */
    pick: ['offTime', 'advanceMinutes', 'repeatDays', 'remindStyle']
  }
})
