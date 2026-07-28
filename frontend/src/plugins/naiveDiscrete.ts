import type { DiscreteApi, GlobalThemeOverrides } from 'naive-ui'
import type { DialogApiInjection } from 'naive-ui/es/dialog/src/DialogProvider'
import type { LoadingBarInst } from 'naive-ui/es/loading-bar/src/LoadingBarProvider'
import type { MessageApiInjection } from 'naive-ui/es/message/src/MessageProvider'
import type { NotificationApiInjection } from 'naive-ui/es/notification/src/NotificationProvider'
import type { BuiltInGlobalTheme } from 'naive-ui/es/themes/interface'
import { createDiscreteApi } from 'naive-ui'

let discrete: DiscreteApi<'message' | 'notification' | 'dialog' | 'loadingBar'> | null = null

/**
 * 初始化 Naive UI 离散 API
 * ⚠️ 必须在 Pinia 安装后调用
 */
export function setupDiscreteApi(): void {
  if (discrete) {
    console.warn('Discrete API already initialized!')
    return
  }

  const themeStore = useThemeStore()

  discrete = createDiscreteApi(
    ['message', 'notification', 'dialog', 'loadingBar'],
    {
      configProviderProps: computed<{ theme: BuiltInGlobalTheme, themeOverrides: GlobalThemeOverrides }>(() => ({
        theme: themeStore.nTheme,
        themeOverrides: themeStore.nThemeOverride
      }))
    }
  )
}

function createLazyDiscreteProxy<T>(key: keyof DiscreteApi): T {
  return new Proxy(
    {},
    {
      get(_, prop) {
        if (!discrete) {
          console.warn(`Discrete API not initialized. Call setupDiscreteApi() first!`)
          return () => { }
        }
        return discrete![key]![prop as keyof typeof discrete[typeof key]]
      }
    }
  ) as T
}

export const message = createLazyDiscreteProxy<MessageApiInjection>('message')
export const notification = createLazyDiscreteProxy<NotificationApiInjection>('notification')
export const dialog = createLazyDiscreteProxy<DialogApiInjection>('dialog')
export const loadingBar = createLazyDiscreteProxy<LoadingBarInst>('loadingBar')
