<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'
import { panelUrl } from '../config'
import { XButton, XNotice } from '../kit'
import AppPage from '../layout/AppPage.vue'

/** Stands in for a screen that has not moved to the new interface yet. */
const { t } = useI18n()
const route = useRoute()
</script>

<template>
  <AppPage :title="route.meta.titleKey ? t(route.meta.titleKey) : ''">
    <XNotice tone="info">
      {{ t('common.screenMoving') }}
      <template v-if="route.meta.oldPath" #actions>
        <XButton :href="panelUrl(route.meta.oldPath)" size="sm" icon="arrowRight">{{ t('common.openOldScreen') }}</XButton>
      </template>
    </XNotice>
  </AppPage>
</template>
