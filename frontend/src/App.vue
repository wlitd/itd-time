<script setup lang="ts">
const { nTheme, nThemeOverride } = storeToRefs(useThemeStore())
const { locale } = storeToRefs(useLocaleStore())
</script>

<template>
  <NConfigProvider
    class="h-full" :theme="nTheme" :theme-overrides="nThemeOverride" :locale="locale.nLocale"
    :date-locale="locale.nDateLocale"
  >
    <RouterView />
  </NConfigProvider>
</template>
