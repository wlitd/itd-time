import type { CSSProperties } from 'vue'
import { useThemeVars } from 'naive-ui'

type CommonSize = 'tiny' | 'small' | 'medium' | 'large'

export function useNaiveSize(size: CommonSize = 'medium') {
  const themeVars = useThemeVars()

  const fontSizeMap = computed<Record<CommonSize, string>>(() => ({
    tiny: themeVars.value.fontSizeTiny,
    small: themeVars.value.fontSizeSmall,
    medium: themeVars.value.fontSizeMedium,
    large: themeVars.value.fontSizeLarge
  }))

  const heightMap = computed<Record<CommonSize, string>>(() => ({
    tiny: themeVars.value.heightTiny,
    small: themeVars.value.heightSmall,
    medium: themeVars.value.heightMedium,
    large: themeVars.value.heightLarge
  }))

  const paddingXMap: Record<CommonSize, number> = {
    tiny: 8,
    small: 10,
    medium: 12,
    large: 14
  }

  const fontSize = computed<string>(() => fontSizeMap.value[size])
  const height = computed<string>(() => heightMap.value[size])
  const paddingX = computed<number>(() => paddingXMap[size])

  function createFontStyle(isStrong: boolean): CSSProperties {
    return {
      fontWeight: isStrong ? themeVars.value.fontWeightStrong : themeVars.value.fontWeight,
      lineHeight: themeVars.value.lineHeight
    }
  }

  const commonFont = computed<CSSProperties>(() => createFontStyle(false))
  const strongFont = computed<CSSProperties>(() => createFontStyle(true))

  return { fontSizeMap, fontSize, paddingX, height, commonFont, strongFont }
}
