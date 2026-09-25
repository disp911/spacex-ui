<script setup lang="ts">
import { ref } from 'vue'
import { api, sleep } from '../lib/api'
import { version } from '../lib/format'
import { t } from '../lib/i18n'
import Icon from '../components/Icon.vue'
import Modal from '../components/Modal.vue'
import Spinner from '../components/Spinner.vue'

const props = defineProps<{ info: { currentVersion: string; latestVersion: string; updateAvailable: boolean } }>()
const emit = defineEmits<{ close: [] }>()
const updating = ref(false)

async function update() {
  if (updating.value || !props.info.updateAvailable) return
  updating.value = true
  const msg = await api.post('panel/api/server/updatePanel')
  if (!msg.success) {
    updating.value = false
    return
  }
  // The panel restarts itself after the update; give it time, then reload.
  await sleep(15000)
  location.reload()
}
</script>

<template>
  <Modal :title="t('pages.index.updatePanel')" :width="436" :closable="!updating" @close="emit('close')">
    <div class="stack">
      <div v-if="info.updateAvailable" class="notice">
        <Icon name="alert" />
        <span>{{ updating ? t('pages.index.dontRefresh') : t('spx.panelUpdateNotice') }}</span>
      </div>
      <div v-else class="notice ok">
        <Icon name="check" />
        <span>{{ t('pages.index.panelUpToDate') }}</span>
      </div>
      <div class="row">
        <span>{{ t('spx.currentVersion') }}</span>
        <span class="v">{{ version(info.currentVersion) }}</span>
      </div>
      <div class="row" :class="{ fresh: info.updateAvailable }">
        <span>{{ t('spx.latestVersion') }}</span>
        <span class="v"><span v-if="info.updateAvailable" class="dot" />{{ version(info.latestVersion || info.currentVersion) }}</span>
      </div>
    </div>
    <template #footer>
      <button v-if="info.updateAvailable" class="ctl primary cta" @click="update">
        <Spinner v-if="updating" :size="15" on-accent /><Icon v-else name="arrowUp" :size="15" />
        {{ updating ? t('spx.updating') : t('pages.index.updatePanel') }}
      </button>
      <button v-else class="ctl cancel" @click="emit('close')">{{ t('close') }}</button>
    </template>
  </Modal>
</template>

<style scoped>
.stack {
  display: flex;
  flex-direction: column;
  gap: 10px;
}
.notice {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 11px 13px;
  border-radius: 11px;
  background: rgba(var(--warn-rgb), 0.1);
  border: 1px solid rgba(var(--warn-rgb), 0.3);
  color: var(--warn);
  font: 400 12px/1.45 var(--font-sans);
}
.notice span {
  color: var(--warn-body);
}
.notice.ok {
  background: var(--accent-soft);
  border-color: rgba(var(--accent-rgb), 0.3);
  color: var(--accent);
}
.notice.ok span {
  color: var(--accent);
}
.row {
  display: flex;
  align-items: center;
  height: 40px;
  padding: 0 13px;
  border-radius: 11px;
  background: var(--fill);
  border: 1px solid var(--border);
  font: 400 12.5px var(--font-sans);
  color: var(--text-3);
}
.row.fresh {
  background: var(--accent-soft);
  border-color: rgba(var(--accent-rgb), 0.3);
}
.v {
  margin-left: auto;
  display: flex;
  align-items: center;
  gap: 10px;
  font: 500 13px var(--font-mono);
  color: var(--text);
}
.fresh .v {
  color: var(--accent);
}
.v .dot {
  width: 6px;
  height: 6px;
}
.cta {
  height: 36px;
  padding: 0 16px;
  font-size: 13px;
}
.cancel {
  background: var(--chip);
  padding: 0 15px;
}
</style>
