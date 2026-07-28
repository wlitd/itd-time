<script setup lang="ts">
import type { RemindMode } from '@/stores/offWork'
import { useThemeVars } from 'naive-ui'
import { useI18n } from 'vue-i18n'
import squidwardPng from '@/assets/images/squidward.png'

const { mode } = defineProps<{
  /** 提醒模式 */
  mode: RemindMode
}>()

const emit = defineEmits<{
  /** 关闭提醒弹窗 */
  close: []
}>()

const { t } = useI18n()
const themeVars = useThemeVars()

const { offTime, advanceMinutes } = storeToRefs(useOffWorkStore())

/** 是否为提前提醒（预览按下班提醒展示） */
const isAdvance = computed<boolean>(() => mode === 'advance')

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
        <div class="text-sm font-semibold truncate min-w-0">
          {{ isAdvance ? t('advanceTitle', { minutes: advanceMinutes }) : t('offworkTitle') }}
        </div>
      </div>
      <NText depth="3" class="text-xs truncate block min-w-0">
        {{ isAdvance ? t('advanceDesc', { offTime }) : t('offworkDesc') }}
      </NText>
    </div>

    <!-- 操作按钮 -->
    <div class="flex items-center gap-2 shrink-0" style="--wails-draggable: no-drag">
      <div
        class="size-8 rounded-full flex-center cursor-pointer bg-#fa5151 hover:bg-#e64340 transition-colors"
        :title="t('later')" @click="emit('close')"
      >
        <div class="i-lucide:x text-white text-base" />
      </div>
      <div
        class="size-8 rounded-full flex-center cursor-pointer bg-#07c160 hover:bg-#06ad56 transition-colors"
        :title="t('gotIt')" @click="emit('close')"
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

<i18n lang="json">
{
  "zh": {
    "offworkTitle": "我下班了，蟹老板！",
    "offworkDesc": "时间到，章鱼哥已准时挂钟离岗 🐙",
    "advanceTitle": "还有 {minutes} 分钟下班",
    "advanceDesc": "{offTime} 准时跑路，先收拾收拾吧",
    "later": "再忍忍",
    "gotIt": "下班！"
  },
  "en": {
    "offworkTitle": "I'm off work, Mr. Krabs!",
    "offworkDesc": "Time's up, Squidward has clocked out 🐙",
    "advanceTitle": "{minutes} minutes to go",
    "advanceDesc": "Leaving at {offTime}, time to pack up",
    "later": "Hold on",
    "gotIt": "Off work!"
  }
}
</i18n>
