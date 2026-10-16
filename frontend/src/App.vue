<script setup lang="ts">
import Toast from 'primevue/toast'
import ConfirmPopup from 'primevue/confirmpopup'
import { isRemindWindow } from './shared/env'
import MainView from './pages/MainView.vue'
import RemindWindow from './features/remind/RemindWindow.vue'

// 主窗与提醒窗共用同一入口，由 URL query（?remind=1）区分窗型。
const isRemind = isRemindWindow()
</script>

<template>
  <Toast position="top-right" />
  <ConfirmPopup />
  <RemindWindow v-if="isRemind" />
  <MainView v-else />
</template>
