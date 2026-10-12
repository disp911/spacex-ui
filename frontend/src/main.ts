import { createPinia } from 'pinia'
import { createApp } from 'vue'
import '@fontsource/ibm-plex-sans/400.css'
import '@fontsource/ibm-plex-sans/500.css'
import '@fontsource/ibm-plex-sans/600.css'
import '@fontsource/ibm-plex-mono/400.css'
import '@fontsource/ibm-plex-mono/500.css'
import '@fontsource/ibm-plex-mono/600.css'
import './styles/tokens.css'
import './styles/base.css'
import './kit/kit.css'
import { toLogin } from './api'
import App from './App.vue'
import { config } from './config'
import { i18n, initialLocale, setLocale } from './i18n'
import { router } from './router'
import { useTheme } from './stores/theme'

if (!config.loggedIn) toLogin()

const app = createApp(App).use(createPinia()).use(i18n).use(router)
useTheme()
await setLocale(initialLocale())
app.mount('#app')
