import type { RouteRecordRaw } from 'vue-router'
import { createI18n } from 'vue-i18n'

export const DEFAULT_LOCALE_KEY = 'zh'

const i18n = createI18n({
  legacy: false,
  locale: DEFAULT_LOCALE_KEY,
  fallbackLocale: 'en',
  messages: {}
})

// 将路由的 i18n 数据添加到 i18n 实例中
;(() => {
  function processRouteForI18n(route: RouteRecordRaw) {
    if (route.meta && route.meta.i18n) {
      const i18nData = route.meta.i18n
      Object.keys(i18nData).forEach((locale) => {
        // 使用 titleKey 作为键来存储翻译文本
        if (route.meta?.titleKey) {
          i18n.global.mergeLocaleMessage(locale, { [route.meta.titleKey]: i18nData[locale] })
        } else {
          console.warn(`[i18n] titleKey is missing for ${route.path}`)
        }
      })
    }

    if (route.children && Array.isArray(route.children)) {
      route.children.forEach(processRouteForI18n)
    }
  }

  routes.forEach(processRouteForI18n)
})()

export default i18n
