import type { Router } from 'vue-router'

/**
 * 注册所有路由守卫
 * ⚠️ 必须在 Pinia 安装后调用
 * @param router Vue Router 实例
 */
export function setupRouterGuards(router: Router): void {
  remindGuard(router)
  loadingGuard(router)
}
