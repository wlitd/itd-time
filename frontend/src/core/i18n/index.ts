import { createI18n } from 'vue-i18n'
import { messages } from './messages'
import type { AppLocale } from './messages'

// 默认语言；实际以 Go 端持久化的 ThemePrefs.locale 为准（store 水合后覆盖）。
export const DEFAULT_LOCALE: AppLocale = 'zh-CN'

export const i18n = createI18n({
  legacy: false,
  globalInjection: true,
  locale: DEFAULT_LOCALE,
  fallbackLocale: 'zh-CN',
  messages,
})

// 供 store 在非组件上下文中切换语言。
export function setAppLocale(locale: AppLocale): void {
  i18n.global.locale.value = locale
}
