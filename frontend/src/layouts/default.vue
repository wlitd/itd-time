<script setup lang="ts">
import { useThemeVars } from 'naive-ui'

const themeVars = useThemeVars()
</script>

<template>
  <div :style="{ backgroundColor: themeVars.bodyColor }" class="h-full flex flex-col gap-3 px-6 pb-3 box-border">
    <DefaultHeader />
    <RouterView class="flex-1" />
  </div>
</template>
