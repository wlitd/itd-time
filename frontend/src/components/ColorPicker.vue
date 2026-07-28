<script setup lang="ts">
import { builtInThemeColor } from '@/stores/theme'

const { primaryColor } = storeToRefs(useThemeStore())
</script>

<template>
  <div class="size-7">
    <NColorPicker v-model:value="primaryColor" :render-label="() => ''" :swatches="builtInThemeColor" />
  </div>
</template>

<style lang="css" scoped>
:deep(.n-color-picker) {
  height: 28px !important;
}
</style>
