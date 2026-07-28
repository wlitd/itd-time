import { createApp } from 'vue'
import App from './App.vue'

import '@/styles/index.css'
import '@/styles/variables.css'

import 'virtual:uno.css'

createApp(App).use(plugins).mount('#app')
