<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'
import { XNotice } from '../kit'
import AppPage from '../layout/AppPage.vue'

/** Stands in for a screen that is not built yet. */
const { t } = useI18n()
const route = useRoute()
</script>

<template>
  <AppPage :title="route.meta.titleKey ? t(route.meta.titleKey) : ''">
    <XNotice tone="info">{{ t('common.screenDraft') }}</XNotice>
  </AppPage>
</template>
