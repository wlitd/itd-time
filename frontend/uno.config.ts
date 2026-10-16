import { defineConfig, presetUno, presetIcons, presetTypography } from 'unocss'

// UnoCSS 配置：工具类 + Lucide/Flag 图标 + 基础排版。
// 图标集合由 @iconify-json/lucide、@iconify-json/flag 自动提供，直接写 i-lucide:xxx / i-flag:xxx。
export default defineConfig({
  presets: [
    presetUno(),
    presetIcons({
      scale: 1.2,
      extraProperties: {
        display: 'inline-block',
        'vertical-align': 'middle',
      },
    }),
    presetTypography(),
  ],
  theme: {
    fontFamily: {
      // Plus Jakarta Sans 优先（如后续本地内嵌），回退系统西文，再回退中文系统栈。
      sans: ['"Plus Jakarta Sans"', 'system-ui', '-apple-system', 'Segoe UI', 'PingFang SC', 'Microsoft YaHei', 'sans-serif'],
      led: ['"LED"', 'ui-monospace', 'monospace'],
    },
    colors: {
      primary: {
        DEFAULT: '#DB2777',
        50: '#FDF2F8',
        100: '#FCE7F3',
        200: '#FBCFE8',
        300: '#F9A8D4',
        400: '#F472B6',
        500: '#EC4899',
        600: '#DB2777',
        700: '#BE185D',
        800: '#9D174D',
        900: '#831843',
      },
      accent: '#F97316',
      // 语义危险色（红）。作为真实色值注册，才能使用 bg-danger/10 这类透明度修饰；
      // 文字色仍走 CSS 变量 var(--app-danger) 以适配 light/dark。
      danger: '#DC2626',
    },
  },
  shortcuts: {
    'flex-center': 'flex items-center justify-center',
    'flex-between': 'flex items-center justify-between',
    // LED 数字微光特效：主色发光，深色下更明显。
    'led-glow': 'text-primary-400 [text-shadow:0_0_6px_rgba(244,114,182,0.55),0_0_16px_rgba(236,72,153,0.35)]',
  },
})
