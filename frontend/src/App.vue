<script setup lang="ts">
import { defineAsyncComponent } from 'vue'
import { config } from './lib/config'
import Toasts from './components/Toasts.vue'

// One bundle entry for both page shells the server renders; each page is
// its own chunk so the login screen doesn't load the dashboard.
const Page = config.page === 'dashboard'
  ? defineAsyncComponent(() => import('./pages/Dashboard.vue'))
  : defineAsyncComponent(() => import('./pages/Login.vue'))
</script>

<template>
  <Page />
  <Toasts />
</template>
