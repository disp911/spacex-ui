<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { XButton, XCard } from '../kit'
import AppPage from '../layout/AppPage.vue'

const { t } = useI18n()
</script>

<template>
  <AppPage :title="t('common.notFound')">
    <XCard>
      <p class="hint">{{ t('common.notFoundHint') }}</p>
      <XButton :to="{ name: 'dashboard' }" variant="primary" icon="dashboard">{{ t('common.toDashboard') }}</XButton>
    </XCard>
  </AppPage>
</template>

<style scoped>
.hint {
  margin: 0 0 14px;
  color: var(--mut);
}
</style>
