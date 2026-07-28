<script setup lang="ts">
import dayjs from 'dayjs'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const { offTime } = storeToRefs(useOffWorkStore())
const { setOffWorkTime } = useOffWorkStore()

const offTimeValue = computed<number>({
  get: () => {
    const [hours, minutes] = offTime.value.split(':').map(Number)
    return dayjs().hour(hours).minute(minutes).second(0).valueOf()
  },
  set: (val: number | null) => {
    if (val) {
      setOffWorkTime(dayjs(val).format('HH:mm'))
    }
  }
})
</script>

<template>
  <FloatTimePicker v-model:value="offTimeValue" format="HH:mm" :placeholder="t('offTime')" clearable>
    <template #icon>
      <div class="i-lucide:alarm-clock" />
    </template>
  </FloatTimePicker>
</template>

<i18n lang="json">
{
  "zh": {
    "offTime": "下班时间"
  },
  "en": {
    "offTime": "Off Time"
  }
}
</i18n>
