import type { GlobalThemeOverrides } from 'naive-ui'
import type { BuiltInGlobalTheme } from 'naive-ui/es/themes/interface'
import { TinyColor } from '@ctrl/tinycolor'
import { darkTheme, lightTheme } from 'naive-ui'

/**
 * 预设主题颜色数组
 */
const PRESET_COLORS = [
  '#FF6699', // 粉色
  '#fa8c16', // 橙色
  '#9FDB1DFF', // 绿色
  '#66FFCC', // 青色
  '#14C9C9FF', // 蓝色
  '#722ed1' // 紫色
] as const

/** 内置主题颜色 */
export const builtInThemeColor = [...PRESET_COLORS]

/**
 * 基础主题覆盖样式
 */
const BASE_THEME_OVERRIDES: GlobalThemeOverrides = {
  common: {
    borderRadius: '6px'
  }
}

/**
 * 根据主色调生成动态主题覆盖样式
 * @param primaryColor 主题颜色
 * @returns 主题覆盖配置对象
 */
const generateDynamicThemeOverride = (primaryColor: string): GlobalThemeOverrides => {
  const tc = new TinyColor(primaryColor)
  if (!tc.isValid) {
    console.warn(`Invalid color: ${primaryColor}, fallback to default`)
    return BASE_THEME_OVERRIDES
  }

  return {
    ...BASE_THEME_OVERRIDES,
    common: {
      ...BASE_THEME_OVERRIDES.common,
      primaryColor: tc.toHexString(),
      primaryColorHover: tc.clone().lighten(10).toHexString(),
      primaryColorPressed: tc.clone().darken(15).toHexString(),
      primaryColorSuppl: tc.clone().setAlpha(1).toHex8String()
    }
  }
}

/**
 * 主题状态管理 Store
 * 用于管理应用的主题设置，包括暗黑模式和主色调
 */
export const useThemeStore = defineStore('theme', () => {
  /** 是否为暗黑模式 */
  const isDark = useDark()
  /** Naive UI 主题对象 */
  const nTheme = computed<BuiltInGlobalTheme>(() => isDark.value ? darkTheme : lightTheme)

  /** 主题主色调 */
  const primaryColor = ref<string>(PRESET_COLORS[0])
  /** 主题覆盖配置 */
  const nThemeOverride = computed<GlobalThemeOverrides>(() => generateDynamicThemeOverride(primaryColor.value))

  /** 监听主色调变化并更新CSS变量 */
  watch(primaryColor, (newColor) => {
    document.documentElement.style.setProperty('--primary-color', newColor)
  }, { immediate: true })

  return { isDark, nTheme, primaryColor, nThemeOverride }
}, {
  /** 持久化配置 */
  persist: {
    /** 持久化存储的键名 */
    key: 'theme',
    /** 使用 localStorage 进行持久化存储 */
    storage: localStorage,
    /** 需要持久化的状态字段 */
    pick: ['primaryColor']
  }
})
