import type { App } from 'vue'

import pinia from '@/stores'
import i18n from './locale'

const pluginList = [pinia, router, i18n]
const directiveList = [{ name: 'ripple', directive: ripple }]

const plugins = {
  install(app: App<Element>) {
    // 批量注册插件  改用自动引入
    Object.entries(pluginList).forEach(([, v]) => {
      app.use(v)
    })
    // 批量注册指令
    Object.entries(directiveList).forEach(([, v]) => {
      app.directive(v.name, v.directive)
    })

    // 执行初始化后的函数
    setupDiscreteApi()
    setupRouterGuards(router)
  }
}

export default plugins
