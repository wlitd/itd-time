<script setup lang="ts">
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const remaindListRef = ref()
</script>

<template>
  <div class="flex flex-col gap-3 overflow-hidden">
    <div class="flex justify-between items-center gap-2">
      <NH3>
        <NText type="info" strong>{{ t('moreReminders') }}</NText>
      </NH3>
      <NButton v-ripple circle @click="remaindListRef.scrollAndAdd()">
        <template #icon>
          <div class="i-lucide-plus" />
        </template>
      </NButton>
    </div>
    <RemindList ref="remaindListRef" />
  </div>
</template>

<i18n lang="json">
{
  "zh": {
    "moreReminders": "更多提醒"
  },
  "en": {
    "moreReminders": "More Reminders"
  }
}
</i18n>
