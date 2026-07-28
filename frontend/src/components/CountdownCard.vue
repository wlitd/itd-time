<script setup lang="ts">
import dayjs from 'dayjs'
import duration from 'dayjs/plugin/duration'
import { useI18n } from 'vue-i18n'

dayjs.extend(duration)

const { t } = useI18n()

const { offTime, advanceMinutes } = storeToRefs(useOffWorkStore())

const now = useNow({ interval: 1000 })

const countdown = computed<string>(() => {
  if (!offTime.value)
    return '00:00:00'

  const target = dayjs(`${dayjs().format('YYYY-MM-DD')}${offTime.value}`)

  const diff = target.diff(now.value)
  if (diff <= 0)
    return '00:00:00'

  return dayjs.duration(diff).format('HH:mm:ss')
})

const notifyTime = computed<string>(() => {
  if (!offTime.value)
    return '--:--'

  return dayjs(`2000-01-01 ${offTime.value}`)
    .subtract(advanceMinutes.value, 'minute')
    .format('HH:mm')
})

const showModal = ref<boolean>(false)

// 同步提醒配置到 Go 端定时器，到点后自动弹出提醒窗口
const { showRemind } = useOffWorkRemind()
</script>

<template>
  <NCard
    embedded hoverable class="cursor-pointer rounded-xl" :content-style="{
      padding: '12px 8px',
    }" @click="showModal = true"
  >
    <div class="flex flex-col gap-1">
      <NText depth="3" class="text-xs pl-2">{{ t('freedomCountdown') }}</NText>
      <div class="flex justify-between items-center gap-2">
        <div class="flex-1 led text-5xl tracking-wider text-center font-medium">{{ countdown }}</div>
        <div class="flex flex-col gap-1 items-end">
          <NText depth="3" class="text-xs">{{ t('leave', { offTime }) }}</NText>
          <NText depth="3" class="text-xs">{{ t('remind', { notifyTime }) }}</NText>
        </div>
      </div>
    </div>
  </NCard>

  <NModal
    v-model:show="showModal" :auto-focus="false" transform-origin="center" class="w-92%" :style="{
      position: 'absolute',
      top: '12px',
      left: '50%',
      transform: 'translateX(-50%)',
      borderRadius: '12px',
    }"
  >
    <NCard :content-style="{ padding: '12px 8px' }">
      <NH2 prefix="bar">
        <NText type="primary">{{ t('setting') }}</NText>
      </NH2>
      <NDivider />
      <NGrid x-gap="8" y-gap="12" :cols="24">
        <NGi :span="12">
          <OffTimeSetting />
        </NGi>
        <NGi :span="12">
          <FloatInputNumber v-model:value="advanceMinutes" :min="0" :placeholder="t('advanceMinutes')" clearable />
        </NGi>
        <NGi :span="24">
          <RepeatSetting />
        </NGi>

        <NGi :span="24">
          <RemindStyleSelect />
        </NGi>
      </NGrid>

      <NDivider />

      <div class="flex justify-end">
        <NButton @click="showRemind('preview')">
          <template #icon>
            <div class="i-lucide:eye" />
          </template>
          {{ t('previewRemind') }}
        </NButton>
      </div>
    </NCard>
  </NModal>
</template>

<i18n lang="json">
{
  "zh": {
    "freedomCountdown": "自由倒计时",
    "leave": "今天 {offTime} 下班",
    "remind": "将在 {notifyTime} 提醒您",
    "setting": "设置",
    "advanceMinutes": "提前提醒",
    "repeat": "重复",
    "previewRemind": "预览提醒"
  },
  "en": {
    "freedomCountdown": "Freedom Countdown",
    "leave": "Leave at {offTime}",
    "remind": "Remind at {notifyTime}",
    "setting": "Setting",
    "advanceMinutes": "Advance Minutes",
    "repeat": "Repeat",
    "previewRemind": "Preview Remind"
  }
}
</i18n>
