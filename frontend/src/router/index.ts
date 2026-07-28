import { createRouter, createWebHistory } from 'vue-router'
import { routes } from './routes'

declare module 'vue-router' {
  interface RouteMeta {
    titleKey?: string // 用于 i18n 翻译键
    i18n?: Record<string, string> // i18n 翻译数据
  }
}

const router = createRouter({
  history: createWebHistory(getEnv('VITE_PUBLIC_PATH')),
  routes
})

export default router
