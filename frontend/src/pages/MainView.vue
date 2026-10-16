<script setup lang="ts">
import AppHeader from '../features/shell/AppHeader.vue'
import CountdownCard from '../features/countdown/CountdownCard.vue'
import TodoList from '../features/todos/TodoList.vue'
import UpdateDialog from '../features/update/UpdateDialog.vue'
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
