import { createApp } from 'vue'
import './styles/tokens.css'
import './styles/base.css'
import App from './App.vue'
import { config } from './lib/config'

document.documentElement.lang = config.lang.slice(0, 2)
createApp(App).mount('#app')
