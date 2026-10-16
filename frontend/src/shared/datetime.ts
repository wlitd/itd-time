import dayjs from 'dayjs'

/** 解析 "HH:mm" 为当日的时间点；非法则回退 17:30。 */
export function parseOffTime(offTime: string, base: dayjs.Dayjs): dayjs.Dayjs {
  const m = /^(\d{1,2}):(\d{2})$/.exec(offTime.trim())
  const hour = m ? Number(m[1]) : 17
  const minute = m ? Number(m[2]) : 30
  return base.hour(hour).minute(minute).second(0).millisecond(0)
}

/**
 * 计算下一次「下班」发生的日期时间：
 * - repeatDays 为空视为每天；
 * - 今天的下班时间若尚未过期且命中重复日，则取今天；否则向后逐日找首个命中日。
 * 返回 null 表示下班提醒未启用。
 */
export function nextOffWork(
  offTime: string,
  repeatDays: number[],
  enabled: boolean,
  now: dayjs.Dayjs = dayjs(),
): dayjs.Dayjs | null {
  if (!enabled)
    return null

  const everyDay = repeatDays.length === 0
  for (let offset = 0; offset <= 7; offset++) {
    const day = now.add(offset, 'day')
    // dayjs.day(): 0=周日 .. 6=周六，与 Go time.Weekday 对齐。
    const hit = everyDay || repeatDays.includes(day.day())
    if (!hit)
      continue
    const target = parseOffTime(offTime, day)
    if (offset > 0 || target.isAfter(now))
      return target
  }
  return null
}

/** 距目标还剩的整秒数（不小于 0）。 */
export function secondsUntil(target: dayjs.Dayjs, now: dayjs.Dayjs = dayjs()): number {
  const diff = target.diff(now, 'second')
  return diff > 0 ? diff : 0
}

/** 秒数格式化为 HH:MM:SS。 */
export function formatHMS(totalSeconds: number): string {
  const h = Math.floor(totalSeconds / 3600)
  const m = Math.floor((totalSeconds % 3600) / 60)
  const s = totalSeconds % 60
  return [h, m, s].map(n => String(n).padStart(2, '0')).join(':')
}

/** 星期序号（0=周日）到 RepeatSetting <i18n> 块内星期短键的映射。 */
export const WEEKDAY_KEYS: string[] = ['sun', 'mon', 'tue', 'wed', 'thu', 'fri', 'sat']

/** 毫秒时间戳按当前语言格式化（用于代办截止时间展示）。 */
export function formatDeadline(ms: number, pattern: string): string {
  return dayjs(ms).format(pattern)
}

/**
 * 将毫秒时间戳归整到「整分」（秒/毫秒清零）。
 * DatePicker 选时/分时保留了初始 Date 的秒与毫秒，会让提醒时刻偏离用户预期的整分；
 * 存入前归零，保证「提前提醒」按用户选定的分钟准时触发。
 */
export function floorToMinute(ms: number): number {
  return dayjs(ms).second(0).millisecond(0).valueOf()
}

/** "HH:mm" -> 今天的 Date（供 PrimeVue DatePicker timeOnly 绑定）。非法回退 17:30。 */
export function hhmmToDate(hhmm: string): Date {
  const m = /^(\d{1,2}):(\d{2})$/.exec(hhmm.trim())
  const d = new Date()
  d.setHours(m ? Number(m[1]) : 17, m ? Number(m[2]) : 30, 0, 0)
  return d
}

/** Date -> "HH:mm"（24 小时制）。 */
export function dateToHHMM(date: Date): string {
  return `${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}`
}
