<script setup lang="ts">
import { ref } from 'vue'
import XModal from '@/components/ui/XModal.vue'
import XIcon from '@/components/ui/XIcon.vue'
import XSpinner from '@/components/ui/XSpinner.vue'
import { server } from '@/api/panel'
import { panelUrl } from '@/env'
import { toast, toastError, toastMsg } from '@/composables/useToast'
import { t } from '@/i18n'

const open = defineModel<boolean>('open', { default: false })
const importing = ref(false)
const fileInput = ref<HTMLInputElement>()

function exportDb() {
  window.location.href = panelUrl('panel/api/server/getDb')
}

async function onFile(e: Event) {
  const file = (e.target as HTMLInputElement).files?.[0]
  ;(e.target as HTMLInputElement).value = ''
  if (!file) return
  importing.value = true
  try {
    const up = await server.importDb(file)
    if (!up.success) {
      toastMsg(up)
      importing.value = false
      return
    }
    toast(t('backup.importing'), 'info', 6000)
    const r = await server.restartPanel()
    if (r.success) setTimeout(() => location.reload(), 5000)
    else {
      toastMsg(r)
      importing.value = false
    }
  } catch (err) {
    toastError(err)
    importing.value = false
  }
}
</script>

<template>
  <XModal v-model:open="open" :title="t('backup.title')" :width="436" :closable="!importing">
    <div class="backup">
      <button type="button" class="backup__row" @click="exportDb">
        <div class="backup__texts">
          <div class="backup__title">{{ t('backup.exportTitle') }}</div>
          <div class="backup__text">{{ t('backup.exportText') }}</div>
        </div>
        <span class="backup__btn"><XIcon name="download" :stroke-width="1.6" /></span>
      </button>
      <button type="button" class="backup__row" :disabled="importing" @click="fileInput?.click()">
        <div class="backup__texts">
          <div class="backup__title">{{ t('backup.importTitle') }}</div>
          <div class="backup__text">{{ t('backup.importText') }}</div>
        </div>
        <span class="backup__btn">
          <XSpinner v-if="importing" />
          <XIcon v-else name="upload" :stroke-width="1.6" />
        </span>
      </button>
      <input ref="fileInput" type="file" accept=".db" hidden @change="onFile" />
    </div>
  </XModal>
</template>

<style scoped>
.backup {
  background: var(--field);
  border: 1px solid var(--border-field);
  border-radius: var(--r-md);
  overflow: hidden;
}
.backup__row {
  width: 100%;
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 14px 15px;
  border: 0;
  background: transparent;
  text-align: left;
  transition: background var(--dur-fast) ease;
}
.backup__row + .backup__row {
  border-top: 1px solid var(--divider);
}
.backup__row:hover {
  background: var(--control);
}
.backup__row:disabled {
  cursor: progress;
}
.backup__texts {
  flex: 1;
  min-width: 0;
}
.backup__title {
  margin-bottom: 4px;
  font: var(--fw-semibold) var(--fs-md) var(--font-sans);
  letter-spacing: -0.15px;
  color: var(--text);
}
.backup__text {
  font: var(--fw-regular) var(--fs-sm) / 1.5 var(--font-sans);
  color: var(--text-3);
  text-wrap: pretty;
}
.backup__btn {
  flex: none;
  width: var(--h-md);
  height: var(--h-md);
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--r-md);
  background: var(--accent);
  color: var(--on-accent);
  transition: background 0.2s ease, box-shadow var(--dur-base) ease;
}
.backup__row:hover .backup__btn {
  background: var(--accent-hover);
  box-shadow: 0 0 0 4px var(--accent-ring), 0 0 18px var(--accent-glow);
}
</style>
