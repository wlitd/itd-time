import Aura from '@primevue/themes/aura'
import { definePreset } from '@primevue/themes'

// 品牌主色「玫红」rose 阶（锚在品红侧 ~330°，贴合背景粉带并与 danger 纯红拉开距离），
// 与 DESIGN.md / uno theme 保持一致。
const rose = {
  50: '#fdf2f8',
  100: '#fce7f3',
  200: '#fbcfe8',
  300: '#f9a8d4',
  400: '#f472b6',
  500: '#ec4899',
  600: '#db2777',
  700: '#be185d',
  800: '#9d174d',
  900: '#831843',
  950: '#500724',
}

// 基于 Aura 定制：仅替换 primary 语义色阶，其余（surface/圆角/动效）沿用 Aura 默认。
// Aura 会依据 primary 色阶自动推导 color/hoverColor/activeColor 等 token。
export const appPreset = definePreset(Aura, {
  semantic: {
    primary: {
      50: rose[50],
      100: rose[100],
      200: rose[200],
      300: rose[300],
      400: rose[400],
      500: rose[500],
      600: rose[600],
      700: rose[700],
      800: rose[800],
      900: rose[900],
      950: rose[950],
    },
    colorScheme: {
      light: {
        primary: {
          color: '{primary.600}',
          contrastColor: '#ffffff',
          hoverColor: '{primary.700}',
          activeColor: '{primary.800}',
        },
        focusRing: '{primary.color}',
      },
      dark: {
        primary: {
          color: '{primary.400}',
          contrastColor: '{surface.900}',
          hoverColor: '{primary.300}',
          activeColor: '{primary.200}',
        },
        focusRing: '{primary.color}',
      },
    },
  },
})

// 深色模式选择器：由 store 在 <html> 上切换 .dark 类驱动。
export const DARK_MODE_SELECTOR = '.dark'
