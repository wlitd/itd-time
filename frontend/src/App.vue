<script setup lang="ts">
import { useTodoRemind } from '@/composables/useTodoRemind'

const { nTheme, nThemeOverride } = storeToRefs(useThemeStore())
const { locale } = storeToRefs(useLocaleStore())

// 启动 TODO 提醒轮询（主窗口全局生效）
useTodoRemind()
</script>

<template>
  <NConfigProvider
    class="h-full" :theme="nTheme" :theme-overrides="nThemeOverride" :locale="locale.nLocale"
    :date-locale="locale.nDateLocale"
  >
    <RouterView />
  </NConfigProvider>
</template>
