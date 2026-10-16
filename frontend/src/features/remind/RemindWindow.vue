<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '../../core/stores/app'
import { Coordinator, RemindKind } from '../../core/bindings'
import { readRemindParams } from '../../shared/env'
import { resolveRemind } from './registry'

interface RemindContent {
  title: string
  description: string
}

const { t } = useI18n()
// 下班标题与设置弹窗预览共享，留在全局，用 global composer 直读。
const { t: gt } = useI18n({ useScope: 'global' })
const store = useAppStore()
const params = readRemindParams()

const card = computed(() => resolveRemind(store.remind.style))

// Coordinator 仅回传 kind + 兜底中文标题；此处按 kind 用当前语言与配置重建文案。
const content = computed<RemindContent>(() => {
  const minutes = store.remind.advanceMinutes ?? 0
  const offTime = store.offWork.offTime
  switch (params.kind) {
    case RemindKind.RemindKindOffWorkAdvance:
      return { title: t('advanceTitle', { minutes }), description: t('advanceDesc', { offTime }) }
    case RemindKind.RemindKindTodo:
      return {
        title: t('todoTitle', { title: params.title || t('untitled') }),
        description: t('todoDesc'),
      }
    case RemindKind.RemindKindOffWork:
      return { title: gt('remind.offworkTitle'), description: t('offworkDesc') }
    default:
      return { title: params.title || gt('remind.offworkTitle'), description: '' }
  }
})

function dismiss(): void {
  void Coordinator.Dismiss(params.id)
}
</script>

<template>
  <div class="h-screen w-screen">
    <component
      :is="card"
      :title="content.title"
      :description="content.description"
      @confirm="dismiss"
      @cancel="dismiss"
    />
  </div>
</template>

<i18n lang="yaml">
zh-CN:
  advanceTitle: "还有 {minutes} 分钟下班"
  advanceDesc: "{offTime} 准时跑路，先收拾收拾吧"
  todoTitle: "待办：{title}"
  todoDesc: 该处理了，别摸鱼了
  offworkDesc: 时间到，章鱼哥已准时挂钟离岗
  untitled: 未命名
en-US:
  advanceTitle: "{minutes} minutes to go"
  advanceDesc: "Leaving at {offTime}, time to pack up"
  todoTitle: "Todo: {title}"
  todoDesc: Time to get it done
  offworkDesc: "Time's up, Squidward has clocked out"
  untitled: Untitled
</i18n>
