<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { WEEKDAY_KEYS } from '../../../shared/datetime'
import type { RepeatSettingEmits, RepeatSettingProps } from './types'

const { modelValue } = defineProps<RepeatSettingProps>()
const emit = defineEmits<RepeatSettingEmits>()
const { t } = useI18n()

// 展示顺序：周一优先，周日殿后。
const orderedDays: number[] = [1, 2, 3, 4, 5, 6, 0]

function isSelected(day: number): boolean {
  return modelValue.includes(day)
}

function toggle(day: number): void {
  const next = isSelected(day)
    ? modelValue.filter(d => d !== day)
    : [...modelValue, day].sort((a, b) => a - b)
  emit('update:modelValue', next)
}
</script>

<template>
  <div class="flex flex-wrap gap-2">
    <button
      v-for="day in orderedDays"
      :key="day"
      type="button"
      class="h-9 w-9 rounded-full border text-sm font-medium transition-colors duration-200 cursor-pointer"
      :class="isSelected(day)
        ? 'border-transparent bg-primary text-white'
        : 'border-[var(--app-border)] text-[var(--app-muted)] hover:border-primary hover:text-primary'"
      @click="toggle(day)"
    >
      {{ t(WEEKDAY_KEYS[day]) }}
    </button>
  </div>
</template>

<i18n lang="yaml">
zh-CN:
  mon: 一
  tue: 二
  wed: 三
  thu: 四
  fri: 五
  sat: 六
  sun: 日
en-US:
  mon: Mon
  tue: Tue
  wed: Wed
  thu: Thu
  fri: Fri
  sat: Sat
  sun: Sun
</i18n>
