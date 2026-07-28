import type { RouteRecordRaw } from 'vue-router'

export const routes: RouteRecordRaw[] = [
  {
    path: '/',
    component: () => import('@/layouts/default.vue'),
    children: [
      {
        path: '',
        name: 'Home',
        component: () => import('@/pages/index.vue'),
        meta: {
          titleKey: 'home',
          i18n: {
            zh: '首页',
            en: 'Home'
          }
        }
      }
    ]
  },
  {
    // 提醒弹窗页面（独立提醒进程使用，不套默认布局）
    path: '/remind',
    name: 'Remind',
    component: () => import('@/pages/remind.vue'),
    meta: {
      titleKey: 'remind',
      i18n: {
        zh: '提醒',
        en: 'Remind'
      }
    }
  }
]
