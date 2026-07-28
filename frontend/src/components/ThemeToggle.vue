<script setup lang="ts">
import { useThemeVars } from 'naive-ui'
import { useThemeStore } from '@/stores/theme'

const themeVars = useThemeVars()
const { isDark } = storeToRefs(useThemeStore())

function onToggleTheme(e: MouseEvent) {
  if (!document.startViewTransition) {
    useToggle(isDark)()
    return
  }

  const transition = document.startViewTransition(() => {
    useToggle(isDark)()
  })

  const x = e.clientX
  const y = e.clientY
  const endRadius = Math.hypot(
    Math.max(x, window.innerWidth - x),
    Math.max(y, window.innerHeight - y)
  )

  transition.ready.then(() => {
    const clipPath = [
      `circle(0px at ${x}px ${y}px)`,
      `circle(${endRadius}px at ${x}px ${y}px)`
    ]

    document.documentElement.animate(
      {
        clipPath: isDark.value ? (clipPath as any).toReversed() : clipPath
      },
      {
        duration: 400,
        easing: 'ease-in',
        pseudoElement: `::view-transition-${isDark.value ? 'old' : 'new'}(root)`,
        fill: 'both'
      }
    )
  })
}
</script>

<template>
  <div
    :style="{ color: themeVars.textColorBase }" class="size-5 cursor-pointer"
    :class="isDark ? 'i-lucide:moon' : 'i-lucide:sun'"
    @click="onToggleTheme"
  />
</template>
