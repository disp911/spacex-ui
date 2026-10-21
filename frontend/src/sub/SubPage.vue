<script setup lang="ts">
import { watchEffect } from 'vue'
import { useI18n } from 'vue-i18n'
import { XToaster } from '../kit'
import type { SubData } from './data'

/** The subscription page of mockup «xpanel Subscription». */
defineProps<{ data: SubData | null }>()
const { t } = useI18n()

watchEffect(() => (document.title = t('sub.title')))
</script>

<template>
  <main class="sub">
    <h1 class="sub__title">{{ t('sub.title') }}</h1>
    <p v-if="data" class="mono">{{ data.subId }}</p>
  </main>
  <XToaster />
</template>

<style scoped>
.sub {
  padding: var(--sp-24) var(--sp-16);
}

.sub__title {
  margin: 0;
  font: 600 var(--fs-16) / var(--lh-tight) var(--font-sans);
}
</style>
