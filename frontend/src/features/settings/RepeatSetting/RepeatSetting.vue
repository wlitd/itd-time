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

/** 星期短名：共享词表在 core/i18n 全局 messages，local 作用域缺失时自动回落。 */
function dayLabel(day: number): string {
  return t(`weekdays.${WEEKDAY_KEYS[day]}`)
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
      {{ dayLabel(day) }}
    </button>
  </div>
</template>
