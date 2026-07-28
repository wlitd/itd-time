import type { NDateLocale, NLocale } from 'naive-ui'
import { dateEnUS, dateZhCN, enUS, zhCN } from 'naive-ui'
import i18n from '@/plugins/locale'

/**
 * 语言环境键名类型
 */
export type LocaleKey = 'zh' | 'en'

/**
 * 语言环境信息接口
 */
interface LocaleInfo {
  /** 语言环境键名 */
  key: LocaleKey
  /** Naive UI 语言环境 */
  nLocale: NLocale
  /** Naive UI 日期语言环境 */
  nDateLocale: NDateLocale
}

/**
 * 语言环境映射表
 */
export const localeMap: LocaleInfo[] = [
  {
    key: 'zh',
    nLocale: zhCN,
    nDateLocale: dateZhCN
  },
  {
    key: 'en',
    nLocale: enUS,
    nDateLocale: dateEnUS
  }
]

/**
 * 语言环境状态管理 Store
 * 用于管理应用的国际化设置
 */
export const useLocaleStore = defineStore('locale', () => {
  /** 当前语言环境键名 */
  const currentLocaleKey = ref<LocaleKey>(DEFAULT_LOCALE_KEY)

  /** 当前语言环境信息 */
  const locale = computed<LocaleInfo>(() => {
    return localeMap.find(l => l.key === currentLocaleKey.value) as LocaleInfo
  })

  /** 监听语言环境变化并更新全局国际化配置 */
  watch(currentLocaleKey, (newLocaleKey) => {
    i18n.global.locale.value = newLocaleKey
  }, { immediate: true })

  /**
   * 设置语言环境键名
   * @param localeKey 语言环境键名
   */
  function setLocaleKey(localeKey: LocaleKey): void {
    currentLocaleKey.value = localeKey
  }

  return { locale, setLocaleKey, currentLocaleKey }
}, {
  /** 持久化配置 */
  persist: {
    /** 持久化存储的键名 */
    key: 'locale',
    /** 使用 localStorage 进行持久化存储 */
    storage: localStorage,
    /** 需要持久化的状态字段 */
    pick: ['currentLocaleKey']
  }
})
