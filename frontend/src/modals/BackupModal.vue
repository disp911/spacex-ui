<script setup lang="ts">
import { ref } from 'vue'
import { api, sleep } from '../lib/api'
import { panelUrl } from '../lib/config'
import { t } from '../lib/i18n'
import Icon from '../components/Icon.vue'
import Modal from '../components/Modal.vue'
import Spinner from '../components/Spinner.vue'

const emit = defineEmits<{ close: [] }>()
const importing = ref(false)
const fileInput = ref<HTMLInputElement>()

function exportDb() {
  location.assign(panelUrl('panel/api/server/getDb'))
}

async function onFile(e: Event) {
  const file = (e.target as HTMLInputElement).files?.[0]
  if (!file) return
  importing.value = true
  const form = new FormData()
  form.append('db', file)
  const msg = await api.post('panel/api/server/importDB', form)
  if (!msg.success) {
    importing.value = false
    return
  }
  // The imported database only takes effect after a panel restart.
  const restart = await api.post('panel/setting/restartPanel')
  if (restart.success) {
    await sleep(5000)
    location.reload()
  } else {
    importing.value = false
  }
}
</script>

<template>
  <Modal :title="t('pages.index.backup')" :width="436" :closable="!importing" @close="emit('close')">
    <div class="box">
      <div class="action">
        <div class="texts">
          <div class="name">{{ t('pages.index.exportDatabase') }}</div>
          <div class="desc">{{ t('spx.exportDesc') }}</div>
        </div>
        <button class="go" :title="t('pages.index.exportDatabase')" @click="exportDb"><Icon name="download" /></button>
      </div>
      <div class="action">
        <div class="texts">
          <div class="name">{{ t('pages.index.importDatabase') }}</div>
          <div class="desc">{{ importing ? t('pages.index.dontRefresh') : t('spx.importDesc') }}</div>
        </div>
        <button class="go" :title="t('pages.index.importDatabase')" :disabled="importing" @click="fileInput?.click()">
          <Spinner v-if="importing" :size="15" on-accent /><Icon v-else name="upload" />
        </button>
        <input ref="fileInput" type="file" accept=".db" hidden @change="onFile" />
      </div>
    </div>
  </Modal>
</template>

<style scoped>
.box {
  background: var(--fill);
  border: 1px solid var(--border);
  border-radius: 12px;
  overflow: hidden;
}
.action {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 14px 15px;
}
.action + .action {
  border-top: 1px solid var(--border);
}
.texts {
  flex: 1;
  min-width: 0;
}
.name {
  font: 600 13px var(--font-sans);
  margin-bottom: 3px;
}
.desc {
  font: 400 11.5px/1.5 var(--font-sans);
  color: var(--text-3);
}
.go {
  flex: none;
  width: 36px;
  height: 36px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--accent);
  color: var(--on-accent);
  transition:
    background 0.2s ease,
    box-shadow 0.22s ease;
}
.go:hover {
  background: var(--accent-hover);
  box-shadow: var(--btn-glow);
}
.go:disabled {
  cursor: progress;
}
</style>
