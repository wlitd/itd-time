// 中英双语文案。此处仅保留「跨组件共享」的键（被 ≥2 个组件引用），
// 单组件独占的文案已下沉至各组件的 <i18n> 自定义块自行维护。
// 组件内 useI18n() 为 local 作用域，本地短键优先、缺失时回落到此处全局键。
export const messages = {
  'zh-CN': {
    todos: {
      titlePlaceholder: '待办事项',
      deadlinePlaceholder: '截止时间',
    },
    remind: {
      offworkTitle: '我下班了，蟹老板！',
    },
    update: {
      cancel: '取消',
    },
  },
  'en-US': {
    todos: {
      titlePlaceholder: 'Todo title',
      deadlinePlaceholder: 'Deadline',
    },
    remind: {
      offworkTitle: "I'm off work, Mr. Krabs!",
    },
    update: {
      cancel: 'Cancel',
    },
  },
} as const

export type AppLocale = keyof typeof messages
