<script setup lang="ts">
import { ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { api, ApiError } from '../../api'
import { toast, XButton, XChip, XInfoRow, XModal, XNotice, XStatusDot } from '../../kit'

export interface PanelUpdateInfo {
  currentVersion: string
  latestVersion: string
  updateAvailable: boolean
}

/** The current and the latest panel version; updating runs the updater and reloads once the panel is back. */
const open = defineModel<boolean>('open', { default: false })
const current = defineModel<PanelUpdateInfo>('info', { required: true })
const { t } = useI18n()

// Asked again on opening, so a version installed meanwhile shows as up to date, here and on the dashboard.
watch(open, async (on) => {
  if (!on) return
  try {
    current.value = await api.get<PanelUpdateInfo>('panel/api/server/getPanelUpdateInfo')
  } catch {
    // keep what the dashboard already knows
  }
})

const updating = ref(false)

async function update() {
  updating.value = true
  try {
    await api.post('panel/api/server/updatePanel')
    // The updater restarts the panel; come back once it answers again.
    setTimeout(() => location.reload(), 15000)
  } catch (e) {
    toast('err', t('dashboard.actionFailed', { error: e instanceof ApiError ? e.message : String(e) }))
    updating.value = false
  }
}
</script>

<template>
  <XModal v-model:open="open" :title="t('dashboard.panelUpdate')" size="sm" :dismissible="!updating">
    <XNotice v-if="current.updateAvailable" tone="warn">{{ t('dashboard.panelUpdateText') }}</XNotice>

    <div class="panel-update__rows">
      <XInfoRow :label="t('dashboard.panelCurrent')">{{ current.currentVersion }}</XInfoRow>
      <XInfoRow v-if="current.updateAvailable" :label="t('dashboard.panelLatest')" accent>
        <XStatusDot size="sm" blink />{{ current.latestVersion }}
      </XInfoRow>
      <XInfoRow v-else :label="t('dashboard.panelUpToDate')" accent>
        <XChip tone="ok">{{ t('dashboard.panelUpdated') }}</XChip>
      </XInfoRow>
    </div>

    <template #footer>
      <XButton variant="primary" icon="arrowUp" :loading="updating" :disabled="!current.updateAvailable" @click="update">
        {{ updating ? t('dashboard.panelUpdating') : t('dashboard.panelUpdateRun') }}
      </XButton>
    </template>
  </XModal>
</template>

<style scoped>
.panel-update__rows {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
</style>
