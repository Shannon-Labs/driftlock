import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { createHead } from '@vueuse/head'
import App from './App.vue'
import router from './router'
import './style.css'

// Note: Cloudflare Pages DevTools overlay can cause 'mce-autosize-textarea already defined' errors
// This is handled automatically by Cloudflare in production deployments
// The overlay is only used for preview builds and development

const app = createApp(App)
const pinia = createPinia()
const head = createHead()

app.use(pinia)
app.use(head)
app.use(router)
app.mount('#app')