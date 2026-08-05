<script setup lang="ts">
import type { Component } from 'vue'
import type { RemindMode, RemindStyle } from '@/stores/offWork'
import { useI18n } from 'vue-i18n'
import SquidwardRemind from '@/components/reminds/SquidwardRemind.vue'
import { CloseRemind, GetRemindAdvanceMinutes, GetRemindMode, GetRemindOffTime, GetRemindStyle, GetTodoTitle } from '../../wailsjs/go/main/App'
import { WindowShow } from '../../wailsjs/runtime/runtime'

/** 提醒弹窗自动关闭时间（毫秒） */
const AUTO_CLOSE_MS = 30 * 1000

/** 提醒样式与组件的映射，后续扩充新样式时在此注册 */
const remindComponents: Record<RemindStyle, Component> = {
  squidward: SquidwardRemind
}

const { t } = useI18n()
const route = useRoute()

const style = ref<RemindStyle>('squidward')
const mode = ref<RemindMode>('offwork')
const todoTitle = ref<string>('')
const ready = ref<boolean>(false)

// 提醒弹窗是独立进程，不能依赖主窗口的 store（localStorage 隔离）
// offTime / advanceMinutes 由 Go 端 ShowRemind 通过 CLI 参数传入
const offTime = ref<string>('18:00')
const advanceMinutes = ref<number>(10)

/** 提醒标题（按模式计算） */
const remindTitle = computed<string>(() => {
  switch (mode.value) {
    case 'todoAdvance':
      return t('todoTitle', { title: todoTitle.value || t('untitled') })
    case 'advance':
      return t('advanceTitle', { minutes: advanceMinutes.value })
    default:
      return t('offworkTitle')
  }
})

/** 提醒描述（按模式计算） */
const remindDesc = computed<string>(() => {
  switch (mode.value) {
    case 'todoAdvance':
      return t('todoDesc')
    case 'advance':
      return t('advanceDesc', { offTime: offTime.value })
    default:
      return t('offworkDesc')
  }
})

/** 关闭提醒弹窗（退出提醒进程） */
function close(): void {
  if (isWails())
    CloseRemind()
}

onMounted(async () => {
  if (isWails()) {
    style.value = (await GetRemindStyle() || 'squidward') as RemindStyle
    mode.value = (await GetRemindMode() || 'offwork') as RemindMode
    todoTitle.value = (await GetTodoTitle()) || ''
    offTime.value = (await GetRemindOffTime()) || '18:00'
    advanceMinutes.value = await GetRemindAdvanceMinutes()
  } else {
    // 浏览器调试：通过 query 参数指定样式与模式
    style.value = (route.query.style as RemindStyle) || 'squidward'
    mode.value = (route.query.mode as RemindMode) || 'offwork'
  }

  ready.value = true

  // 等待提醒内容渲染完成后再显示窗口，避免白屏闪烁
  await nextTick()
  if (isWails())
    WindowShow()

  setTimeout(close, AUTO_CLOSE_MS)
})
</script>

<template>
  <component
    :is="remindComponents[style]"
    v-if="ready"
    :title="remindTitle"
    :description="remindDesc"
    :confirm-text="t('confirmText')"
    :cancel-text="t('cancelText')"
    @confirm="close"
    @cancel="close"
  />
</template>

<i18n lang="json">
{
  "zh": {
    "offworkTitle": "我下班了，蟹老板！",
    "offworkDesc": "时间到，章鱼哥已准时挂钟离岗 🐙",
    "advanceTitle": "还有 {minutes} 分钟下班",
    "advanceDesc": "{offTime} 准时跑路，先收拾收拾吧",
    "todoTitle": "待办: {title}",
    "todoDesc": "该处理了，别摸鱼了 🐟",
    "confirmText": "朕知道了",
    "cancelText": "我就硬拖",
    "untitled": "未命名"
  },
  "en": {
    "offworkTitle": "I'm off work, Mr. Krabs!",
    "offworkDesc": "Time's up, Squidward has clocked out 🐙",
    "advanceTitle": "{minutes} minutes to go",
    "advanceDesc": "Leaving at {offTime}, time to pack up",
    "todoTitle": "Todo: {title}",
    "todoDesc": "Time to get it done 🐟",
    "confirmText": "Yeah, I know",
    "cancelText": "I'll wing it",
    "untitled": "Untitled"
  }
}
</i18n>
