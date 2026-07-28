<script setup lang="ts">
import type { Component } from 'vue'
import type { RemindMode, RemindStyle } from '@/stores/offWork'
import SquidwardRemind from '@/components/reminds/SquidwardRemind.vue'
import { CloseRemind, GetRemindMode, GetRemindStyle } from '../../wailsjs/go/main/App'
import { WindowShow } from '../../wailsjs/runtime/runtime'

/** 提醒弹窗自动关闭时间（毫秒） */
const AUTO_CLOSE_MS = 30 * 1000

/** 提醒样式与组件的映射，后续扩充新样式时在此注册 */
const remindComponents: Record<RemindStyle, Component> = {
  squidward: SquidwardRemind
}

const route = useRoute()

const style = ref<RemindStyle>('squidward')
const mode = ref<RemindMode>('offwork')
const ready = ref<boolean>(false)

/** 关闭提醒弹窗（退出提醒进程） */
function close(): void {
  if (isWails())
    CloseRemind()
}

onMounted(async () => {
  if (isWails()) {
    style.value = (await GetRemindStyle() || 'squidward') as RemindStyle
    mode.value = (await GetRemindMode() || 'offwork') as RemindMode
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
  <component :is="remindComponents[style]" v-if="ready" :mode="mode" @close="close" />
</template>
