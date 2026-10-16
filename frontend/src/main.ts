import { createApp } from 'vue'
import { createPinia } from 'pinia'
import PrimeVue from 'primevue/config'
import ToastService from 'primevue/toastservice'
import ConfirmationService from 'primevue/confirmationservice'
import Tooltip from 'primevue/tooltip'
import { i18n } from './core/i18n'
import { appPreset, DARK_MODE_SELECTOR } from './core/theme/preset'
import { useAppStore } from './core/stores/app'
import { isRemindWindow } from './shared/env'
import App from './App.vue'

import 'virtual:uno.css'
import './styles/main.css'

// 提醒窗需要整窗透明：给根节点打标，交由 CSS 将 body 背景置为透明，
// 避免 CSS 背景覆盖原生 BackgroundTypeTransparent 设置。
if (isRemindWindow())
  document.documentElement.classList.add('remind-window')

async function bootstrap(): Promise<void> {
  const app = createApp(App)

  app.use(createPinia())
  app.use(i18n)
  app.use(PrimeVue, {
    theme: {
      preset: appPreset,
      options: { darkModeSelector: DARK_MODE_SELECTOR, cssLayer: false },
    },
  })
  app.use(ToastService)
  app.use(ConfirmationService)
  app.directive('tooltip', Tooltip)

  // 先水合（从 Go 拉取全量快照并应用主题/语言），就绪后再挂载，避免首屏闪烁。
  await useAppStore().hydrate()

  app.mount('#app')
}

void bootstrap()
