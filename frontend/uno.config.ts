import { defineConfig, presetIcons, presetWind3 } from 'unocss'

export default defineConfig({
  presets: [
    presetWind3({}),
    presetIcons({
      collections: {
        flag: () => import('@iconify-json/flag/icons.json').then(i => i.default),
        lucide: () => import('@iconify-json/lucide/icons.json').then(i => i.default)
      }
    })
  ],

  rules: [
    [/^size-(\d+)$/, ([, d]) => ({ height: `${Number(d) / 4}rem`, width: `${Number(d) / 4}rem` })],
    ['led', { 'font-family': 'LED' }]
  ],

  shortcuts: [
    ['flex-center', 'flex justify-center items-center'],
    ['absolute-y-center', 'absolute top-1/2 transform -translate-y-1/2'],
    ['absolute-x-center', 'absolute left-1/2 transform -translate-x-1/2'],
    ['absolute-center', 'absolute top-1/2 left-1/2 transform -translate-y-1/2 -translate-x-1/2']
  ]
})
