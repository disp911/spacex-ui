<script setup lang="ts">
import { ref } from 'vue'
import XModal from '@/components/ui/XModal.vue'
import XAlert from '@/components/ui/XAlert.vue'
import XButton from '@/components/ui/XButton.vue'
import type { PanelUpdateInfo } from '@/api/panel'
import { server } from '@/api/panel'
import { toast, toastError, toastMsg } from '@/composables/useToast'
import { withV } from '@/utils/format'
import { t } from '@/i18n'

const open = defineModel<boolean>('open', { default: false })
defineProps<{ info: PanelUpdateInfo | null }>()

const busy = ref(false)

async function run() {
  busy.value = true
  try {
    const m = await server.updatePanel()
    if (!m.success) {
      toastMsg(m)
      busy.value = false
      return
    }
    toast(t('panelUpdate.reloading'), 'info', 15000)
    setTimeout(() => location.reload(), 15000)
  } catch (e) {
    toastError(e)
    busy.value = false
  }
}
</script>

<template>
  <XModal v-model:open="open" :title="t('panelUpdate.title')" :width="436" :closable="!busy">
    <XAlert v-if="info?.updateAvailable" tone="warning">{{ t('panelUpdate.warn') }}</XAlert>
    <div class="versions">
      <div class="row">
        <span class="row__label">{{ t('panelUpdate.current') }}</span>
        <span class="row__value mono">{{ withV(info?.currentVersion ?? '') || '—' }}</span>
      </div>
      <div class="row" :class="{ 'is-new': info?.updateAvailable }">
        <span class="row__label">{{ t('panelUpdate.latest') }}</span>
        <span v-if="info?.updateAvailable" class="row__pulse" />
        <span class="row__value mono">{{ withV(info?.latestVersion ?? '') || '—' }}</span>
      </div>
    </div>
    <template #footer>
      <XButton v-if="info?.updateAvailable" variant="primary" icon="upload" :loading="busy" @click="run">
        {{ busy ? t('common.updating') : t('panelUpdate.cta') }}
      </XButton>
      <span v-else class="uptodate">{{ t('panelUpdate.upToDate') }}</span>
    </template>
  </XModal>
</template>

<style scoped>
.versions {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.row {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 11px 13px;
  background: var(--field);
  border: 1px solid var(--border-field);
  border-radius: var(--r-md);
}
.row__label {
  flex: 1;
  min-width: 0;
  font: var(--fw-regular) var(--fs-sm) var(--font-sans);
  color: var(--text-3);
}
.row__value {
  font-size: var(--fs-md);
  font-weight: var(--fw-medium);
  color: var(--text-2);
}
.row.is-new {
  background: var(--accent-soft);
  border-color: var(--accent-border);
}
.row.is-new .row__value {
  color: var(--accent);
}
.row__pulse {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--accent);
  animation: spx-pulse 2.4s infinite;
}
.uptodate {
  margin-right: auto;
  font: var(--fw-regular) var(--fs-sm) var(--font-sans);
  color: var(--text-3);
}
</style>
