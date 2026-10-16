<script setup lang="ts">
import { computed, ref } from 'vue'
import dayjs from 'dayjs'
import { useIntervalFn } from '@vueuse/core'
import { useI18n } from 'vue-i18n'
import Button from 'primevue/button'
import { useAppStore } from '../../core/stores/app'
import { formatHMS, nextOffWork, parseOffTime, secondsUntil } from '../../shared/datetime'
import SettingsDialog from '../settings/SettingsDialog.vue'

const { t } = useI18n()
const store = useAppStore()

const settingsVisible = ref<boolean>(false)

const now = ref<dayjs.Dayjs>(dayjs())
useIntervalFn(() => {
  now.value = dayjs()
}, 1000)

const todayOff = computed<dayjs.Dayjs>(() => parseOffTime(store.offWork.offTime, now.value))

const todayMatches = computed<boolean>(() => {
  const days = store.offWork.repeatDays ?? []
  return days.length === 0 || days.includes(now.value.day())
})

/** 今日已到点且命中重复日：进入「下班」状态。 */
const isOffNow = computed<boolean>(
  () => todayMatches.value && !now.value.isBefore(todayOff.value),
)

const target = computed<dayjs.Dayjs | null>(() =>
  nextOffWork(store.offWork.offTime, [...(store.offWork.repeatDays ?? [])], true, now.value),
)

const countdown = computed<string>(() => (target.value ? formatHMS(secondsUntil(target.value, now.value)) : '--:--:--'))

/** 剩余不足一小时：进入冲刺态，用暖色强调。 */
const isUrgent = computed<boolean>(() => {
  if (!target.value)
    return false
  return secondsUntil(target.value, now.value) <= 3600
})

/** 提醒时刻：下班时间减去提前分钟数。 */
const notifyTime = computed<string>(() => {
  const advance = store.remind.advanceMinutes ?? 0
  return todayOff.value.subtract(advance, 'minute').format('HH:mm')
})
</script>

<template>
  <section class="surface-card relative flex flex-col overflow-hidden px-4 py-2">
    <!-- 顶部光晕，呼应倒计时状态色 -->
    <div
      class="pointer-events-none absolute inset-x-0 -top-16 h-24 blur-2xl transition-colors duration-500"
      :class="isOffNow ? 'bg-primary/25' : isUrgent ? 'bg-accent/20' : 'bg-primary/10'"
    />

    <!-- 顶行：标题左、齿轮右上 -->
    <div class="relative flex items-center justify-between">
      <p class="my-0 text-[11px] font-medium uppercase tracking-[0.18em] text-[var(--app-muted)]">
        {{ t('freedomCountdown') }}
      </p>
      <Button
        text
        rounded
        severity="secondary"
        class="-mr-2 -mt-1 !py-0 wails-no-drag"
        :aria-label="t('setting')"
        @click="settingsVisible = true"
      >
        <template #icon>
          <i class="i-lucide-settings text-xl" />
        </template>
      </Button>
    </div>

    <!-- 主体行：大号 LED 左、两段描述右 -->
    <div class="relative flex items-center justify-between gap-3">
      <p
        v-if="isOffNow"
        class="my-0 text-[2.1rem] font-extrabold leading-none tracking-tight text-primary"
      >
        {{ t('freed') }}
      </p>
      <p
        v-else
        class="led my-0 whitespace-nowrap text-[2.9rem] font-bold leading-none tracking-tighter"
        :class="isUrgent
          ? 'text-[var(--app-accent)] [text-shadow:0_0_12px_rgba(249,115,22,0.5),0_0_28px_rgba(249,115,22,0.28)]'
          : 'led-glow'"
      >
        {{ countdown }}
      </p>

      <div class="flex shrink-0 flex-col items-end gap-1 text-right text-[12px] text-[var(--app-muted)]">
        <span>{{ t('leave', { offTime: store.offWork.offTime }) }}</span>
        <span v-if="!isOffNow && (store.remind.advanceMinutes ?? 0) > 0">{{ t('remind', { notifyTime }) }}</span>
      </div>
    </div>

    <SettingsDialog v-model:visible="settingsVisible" />
  </section>
</template>

<i18n lang="yaml">
zh-CN:
  freedomCountdown: 自由倒计时
  setting: 设置
  freed: 下班！
  leave: "今天 {offTime} 下班"
  remind: "将在 {notifyTime} 提醒您"
en-US:
  freedomCountdown: Freedom Countdown
  setting: Setting
  freed: Off Duty!
  leave: "Leave at {offTime}"
  remind: "Remind at {notifyTime}"
</i18n>
