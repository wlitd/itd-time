<script setup lang="ts">
import { onMounted } from 'vue'
import { useAppStore } from '../core/stores/app'
import { useUpdate } from '../features/update/useUpdate'
import AppHeader from '../features/shell/AppHeader.vue'
import CountdownCard from '../features/countdown/CountdownCard.vue'
import TodoList from '../features/todos/TodoList.vue'
import UpdateDialog from '../features/update/UpdateDialog.vue'

const store = useAppStore()
const { check } = useUpdate()

// 启动静默检查更新（偏好开关控制）：仅发现新版本时弹窗，其余一律不打扰。
// 水合已在 main.ts 挂载前 await 完成，此处读偏好无竞态。
onMounted((): void => {
  if (store.update.checkOnLaunch)
    void check({ silent: true })
})
</script>

<template>
  <div class="flex h-full flex-col text-[var(--app-text)]">
    <AppHeader />

    <!-- 整页不滚动：倒计时卡片固定，代办面板吃掉剩余高度并只在内部滚动 -->
    <main class="flex min-h-0 flex-1 flex-col gap-3 px-3 pb-3 pt-1">
      <CountdownCard />
      <TodoList class="min-h-0 flex-1" />
    </main>

    <UpdateDialog />
  </div>
</template>
