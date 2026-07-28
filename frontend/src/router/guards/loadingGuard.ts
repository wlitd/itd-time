import type { Router } from 'vue-router'

export function loadingGuard(router: Router): void {
  router.beforeEach(() => {
    loadingBar.start()
  })

  router.afterEach(() => {
    loadingBar.finish()
  })

  router.onError(() => {
    loadingBar.error()
  })
}
