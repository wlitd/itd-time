<script setup lang="ts">
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const { repeatDays } = storeToRefs(useOffWorkStore())
const { toggleRepeatDay } = useOffWorkStore()

// 周数据定义
const weekDays = [
  { key: 1, label: 'week.mon' },
  { key: 2, label: 'week.tue' },
  { key: 3, label: 'week.wed' },
  { key: 4, label: 'week.thu' },
  { key: 5, label: 'week.fri' },
  { key: 6, label: 'week.sat' },
  { key: 7, label: 'week.sun' }
]
</script>

<template>
  <div class="flex flex-col gap-3">
    <NText depth="3" class="text-xs">{{ t('repeat') }}</NText>

    <!-- 自适应宽度的周日历选择器 -->
    <div class="flex flex-wrap justify-between gap-2">
      <div
        v-for="day in weekDays" :key="day.key"
        class="cursor-pointer flex-center rounded-xl border-solid border-1 p-1 text-xs transition-all duration-200 flex-1 size-7 shrink-0"
        :class="repeatDays.includes(day.key)
          ? 'bg-[var(--primary-color)] border-transparent'
          : 'hover:border-[var(--primary-color)] hover:text-[var(--primary-color)]'
        " @click="toggleRepeatDay(day.key)"
      >
        {{ t(day.label) }}
      </div>
    </div>
  </div>
</template>

<i18n lang="json">
{
  "zh": {
    "repeat": "重复：",
    "week": {
      "mon": "一",
      "tue": "二",
      "wed": "三",
      "thu": "四",
      "fri": "五",
      "sat": "六",
      "sun": "日"
    }
  },
  "en": {
    "repeat": "Repeat:",
    "week": {
      "mon": "Mon",
      "tue": "Tue",
      "wed": "Wed",
      "thu": "Thu",
      "fri": "Fri",
      "sat": "Sat",
      "sun": "Sun"
    }
  }
}
</i18n>
