<script setup lang="ts">
import { watchEffect } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'
import { XToaster } from './kit'

const { t } = useI18n()
const route = useRoute()

watchEffect(() => {
  const title = route.meta.titleKey ? t(route.meta.titleKey) : ''
  document.title = [title, location.hostname].filter(Boolean).join(' · ')
})
</script>

<template>
  <RouterView />
  <XToaster />
</template>
