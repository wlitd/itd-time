import path from 'node:path'
import VueI18nPlugin from '@intlify/unplugin-vue-i18n/vite'
import vue from '@vitejs/plugin-vue'
import VueJSx from '@vitejs/plugin-vue-jsx'
import UnoCSS from 'unocss/vite'
import AutoImport from 'unplugin-auto-import/vite'
import { NaiveUiResolver } from 'unplugin-vue-components/resolvers'
import Components from 'unplugin-vue-components/vite'
import { defineConfig } from 'vite'

const COMPONENT_INCLUDE_PATTERN = [/\.vue$/, /\.vue\?vue/, /\.tsx$/]

// https://vitejs.dev/config/
export default defineConfig({
  resolve: {
    alias: {
      '@': path.resolve(import.meta.dirname, 'src')
    }
  },

  plugins: [
    vue(),
    VueJSx(),
    UnoCSS(),
    VueI18nPlugin({
      include: []
    }),
    AutoImport({
      imports: [
        'vue',
        '@vueuse/core',
        'pinia',
        'vue-router',
        {
          'naive-ui': ['useDialog', 'useMessage', 'useNotification', 'useLoadingBar']
        }
      ],
      dirs: [
        './src/composables/**',
        './src/plugins/**',
        './src/router/**',
        './src/stores/**',
        './src/utils/**'
      ],
      vueTemplate: true,
      dts: 'src/auto-imports.d.ts'
    }),
    Components({
      resolvers: [NaiveUiResolver()],
      extensions: ['vue'],
      include: COMPONENT_INCLUDE_PATTERN,
      dts: 'src/components.d.ts'
    })
  ]
})
