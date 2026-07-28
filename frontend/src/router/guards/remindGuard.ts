import type { Router } from 'vue-router'
import { IsRemindWindow } from '../../../wailsjs/go/main/App'

/**
 * 提醒弹窗守卫：
 * 提醒进程（--remind 启动）加载的仍是首页路由，
 * 此守卫负责将其重定向到提醒弹窗页面
 */
export function remindGuard(router: Router): void {
  router.beforeEach(async (to) => {
    if (!isWails() || to.name === 'Remind')
      return true

    if (await IsRemindWindow())
      return { name: 'Remind', replace: true }

    return true
  })
}
