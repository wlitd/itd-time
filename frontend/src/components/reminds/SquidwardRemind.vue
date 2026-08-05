<script setup lang="ts">
import { useThemeVars } from 'naive-ui'
import squidwardPng from '@/assets/images/squidward.png'

/** 提醒弹窗通用展示组件，所有文案通过 props 传入，不内置业务逻辑 */
defineProps<{
  /** 提醒标题 */
  title: string
  /** 提醒描述 */
  description: string
  /** 确认按钮文案 */
  confirmText: string
  /** 取消按钮文案 */
  cancelText: string
}>()

const emit = defineEmits<{
  /** 点击确认按钮 */
  confirm: []
  /** 点击取消按钮 */
  cancel: []
}>()

const themeVars = useThemeVars()

/** 半透明卡片背景色（窗口本身全透明） */
const cardBgColor = computed<string>(() => `color-mix(in srgb, ${themeVars.value.cardColor} 88%, transparent)`)
</script>

<template>
  <div
    class="h-screen w-screen box-border flex items-center gap-2 px-2 py-3 rounded-2xl overflow-hidden select-none"
    :style="{ 'backgroundColor': cardBgColor, 'color': themeVars.textColorBase, 'border': `1px solid ${themeVars.dividerColor}`, '--wails-draggable': 'drag' }"
  >
    <!-- 章鱼哥头像 -->
    <NAvatar :src="squidwardPng" round :size="56" class="shrink-0 squidward-shake" />

    <!-- 提醒文案 -->
    <div class="flex flex-col gap-3 flex-1 min-w-0">
      <div class="flex items-center gap-2 min-w-0">
        <!-- 纯 CSS 截断兜底，不使用 NEllipsis，避免 tooltip 在小窗内弹出打乱布局 -->
        <div class="text-sm font-semibold truncate min-w-0">{{ title }}</div>
      </div>
      <NText depth="3" class="text-xs truncate block min-w-0">{{ description }}</NText>
    </div>

    <!-- 操作按钮 -->
    <div class="flex items-center gap-2 shrink-0" style="--wails-draggable: no-drag">
      <div
        class="size-8 rounded-full flex-center cursor-pointer bg-#fa5151 hover:bg-#e64340 transition-colors"
        :title="cancelText" @click="emit('cancel')"
      >
        <div class="i-lucide:x text-white text-base" />
      </div>
      <div
        class="size-8 rounded-full flex-center cursor-pointer bg-#07c160 hover:bg-#06ad56 transition-colors"
        :title="confirmText" @click="emit('confirm')"
      >
        <div class="i-lucide:check text-white text-base" />
      </div>
    </div>
  </div>
</template>

<style scoped>
/* 章鱼哥左右摇摆动画，模拟来电抖动效果 */
@keyframes squidward-shake {
  0%, 100% { transform: rotate(0deg); }
  20% { transform: rotate(-8deg); }
  40% { transform: rotate(8deg); }
  60% { transform: rotate(-5deg); }
  80% { transform: rotate(5deg); }
}

.squidward-shake {
  animation: squidward-shake 1.2s ease-in-out infinite;
}
</style>
