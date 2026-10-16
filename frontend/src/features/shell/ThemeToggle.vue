<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import Button from 'primevue/button'
import { useAppStore } from '../../core/stores/app'

const { t } = useI18n()
const store = useAppStore()

const dark = computed<boolean>(() => store.theme.dark)

function toggle(): void {
  void store.toggleDark()
}
</script>

<template>
  <Button
    severity="secondary"
    text
    rounded
    v-tooltip.bottom="t('theme')"
    :aria-label="t('theme')"
    @click="toggle"
  >
    <template #icon>
      <i :class="dark ? 'i-lucide-sun' : 'i-lucide-moon'" class="text-lg" />
    </template>
  </Button>
</template>

<i18n lang="yaml">
zh-CN:
  theme: 应用主题
en-US:
  theme: Theme
</i18n>
