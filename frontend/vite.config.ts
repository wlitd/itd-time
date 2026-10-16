import process from 'node:process'
import { fileURLToPath, URL } from 'node:url'
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import wails from '@wailsio/runtime/plugins/vite'
import UnoCSS from 'unocss/vite'
import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite'
import VueI18nPlugin from '@intlify/unplugin-vue-i18n/vite'

// https://vitejs.dev/config/
export default defineConfig({
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
  server: {
    host: '127.0.0.1',
    port: Number(process.env.WAILS_VITE_PORT) || 9245,
    strictPort: true,
  },
  plugins: [
    UnoCSS(),
    AutoImport({
      // 自动导入常用组合式 API，减少样板 import。
      imports: ['vue', 'pinia', '@vueuse/core'],
      dts: 'src/auto-imports.d.ts',
      dirs: ['src/shared/composables', 'src/core/stores'],
      vueTemplate: true,
    }),
    Components({
      // 自动注册自研组件（src 下的 *.vue）；PrimeVue 组件按需显式导入。
      dirs: ['src'],
      extensions: ['vue'],
      deep: true,
      dts: 'src/components.d.ts',
    }),
    // 编译各组件的 <i18n> 自定义块；仅保留应用支持的语言，裁剪多余 locale。
    VueI18nPlugin({ onlyLocales: ['zh-CN', 'en-US'] }),
    vue(),
    wails('./bindings'),
  ],
})
