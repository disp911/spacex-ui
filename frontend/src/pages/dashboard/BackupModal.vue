<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { api, ApiError, saveFile } from '../../api'
import { toast, XActionList, XActionRow, XModal } from '../../kit'

/** Export and import of the panel database. */
const open = defineModel<boolean>('open', { default: false })
const { t } = useI18n()

const exporting = ref(false)

async function exportDb() {
  exporting.value = true
  try {
    saveFile(await api.file('panel/api/server/getDb'), 'x-ui.db')
  } catch (e) {
    toast('err', t('dashboard.actionFailed', { error: e instanceof ApiError ? e.message : String(e) }))
  } finally {
    exporting.value = false
  }
}
</script>

<template>
  <XModal v-model:open="open" :title="t('dashboard.backup')" size="sm">
    <XActionList>
      <XActionRow :title="t('dashboard.backupExport')" icon="download" :loading="exporting" @click="exportDb">
        <i18n-t keypath="dashboard.backupExportText" tag="span"><template #ext><code>.db</code></template></i18n-t>
      </XActionRow>
      <XActionRow :title="t('dashboard.backupImport')" icon="upload">
        <i18n-t keypath="dashboard.backupImportText" tag="span"><template #ext><code>.db</code></template></i18n-t>
      </XActionRow>
    </XActionList>
  </XModal>
</template>
