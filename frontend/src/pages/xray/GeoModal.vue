<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { api } from '../../api'
import { formatTimestamp } from '../../i18n/format'
import {
  toast, XActionList, XButton, XChip, XConfirm, XEmpty, XGridTable, XIcon, XInlineButton, XModal, XNotice, XSection, XSpinner, type GridColumn,
} from '../../kit'
import { cause } from '../settings/form'
import GeoSourceModal from './GeoSourceModal.vue'
import type { GeoSource } from './model'

/**
 * Geo files: the standard ones Xray ships with, updated one by one or all at
 * once, and the user's own sources. Everything here is done on the server at
 * once, apart from the template and its Save.
 */
const open = defineModel<boolean>('open', { default: false })
const { t } = useI18n()

const standard = ['geosite.dat', 'geoip.dat', 'geosite_IR.dat', 'geoip_IR.dat', 'geosite_RU.dat', 'geoip_RU.dat']
const openStandard = ref(true)
const openCustom = ref(true)

// --- Standard files -------------------------------------------------------------

/** The file waiting for a yes, '' for all of them. */
const asking = ref<string | null>(null)
/** The file being updated, '' for all of them. */
const updating = ref<string | null>(null)

async function update() {
  const file = asking.value
  asking.value = null
  if (file === null) return
  updating.value = file
  try {
    await api.post(file ? `panel/api/server/updateGeofile/${file}` : 'panel/api/server/updateGeofile')
    toast('ok', t('xray.geoUpdated'))
  } catch (e) {
    toast('err', cause(e))
  } finally {
    updating.value = null
  }
}

// --- The user's sources ---------------------------------------------------------

const sources = ref<GeoSource[]>([])
const loading = ref(false)

async function load() {
  loading.value = true
  try {
    sources.value = (await api.get<GeoSource[] | null>('panel/api/custom-geo/list')) ?? []
  } catch (e) {
    toast('err', cause(e))
  } finally {
    loading.value = false
  }
}
watch(open, (on) => on && load())

const columns = computed<GridColumn[]>(() => [
  { key: 'ext', label: t('xray.geoExt'), size: '2fr' },
  { key: 'updated', label: t('xray.geoLastUpdated') },
  { key: 'actions', label: t('xray.geoActions'), size: '96px' },
])
const ext = (s: GeoSource) => `ext:${s.type}_${s.alias}.dat:tag`

const editing = ref<GeoSource | null>(null)
const formOpen = ref(false)
function edit(s: GeoSource | null) {
  editing.value = s
  formOpen.value = true
}

/** The source being downloaded again, or all of them. */
const pulling = ref<number | 'all' | null>(null)

async function pull(id: number | 'all') {
  pulling.value = id
  try {
    await api.post(id === 'all' ? 'panel/api/custom-geo/update-all' : `panel/api/custom-geo/download/${id}`)
    if (id === 'all') toast('ok', t('xray.geoUpdatedAll'))
  } catch (e) {
    toast('err', cause(e))
  } finally {
    pulling.value = null
    load()
  }
}

const deleting = ref<GeoSource | null>(null)
const deleteOpen = computed({ get: () => !!deleting.value, set: (on) => !on && (deleting.value = null) })

async function remove() {
  const s = deleting.value
  deleting.value = null
  if (!s) return
  try {
    await api.post(`panel/api/custom-geo/delete/${s.id}`)
  } catch (e) {
    toast('err', cause(e))
  }
  load()
}
</script>

<template>
  <XModal v-model:open="open" :title="t('xray.geoFiles')" size="xl" :cancel-text="t('xray.close')">
    <div class="geo">
      <XNotice tone="ok">{{ t('xray.geoAppliedAtOnce') }}</XNotice>

      <XSection v-model:open="openStandard" :title="t('xray.geoStandard')">
        <div class="geo__block">
          <XActionList>
            <div v-for="(file, i) in standard" :key="file" class="geo__file">
              <XChip :tone="i % 2 ? 'ok' : 'unl'" mono>{{ file }}</XChip>
              <XSpinner v-if="updating === file" :size="14" class="geo__busy" />
              <XInlineButton v-else icon="refresh" :label="t('xray.geoDownload')" :disabled="updating !== null" @click="asking = file" />
            </div>
          </XActionList>
          <div class="geo__end">
            <XButton size="sm" icon="refresh" :loading="updating === ''" :disabled="updating !== null && updating !== ''" @click="asking = ''">
              {{ t('xray.geoUpdateAll') }}
            </XButton>
          </div>
        </div>
      </XSection>

      <XSection v-model:open="openCustom" :title="t('xray.geoCustom')">
        <div class="geo__block">
          <p class="geo__hint"><XIcon name="info" :size="13" :stroke="1.4" />{{ t('xray.geoHint') }}</p>
          <div class="geo__tools">
            <XButton variant="primary" size="sm" icon="plus" @click="edit(null)">{{ t('xray.geoAdd') }}</XButton>
            <XButton size="sm" icon="refresh" :loading="pulling === 'all'" :disabled="!sources.length || pulling !== null" @click="pull('all')">
              {{ t('xray.geoUpdateAll') }}
            </XButton>
          </div>
          <XGridTable :columns="columns" :rows="sources" :row-key="(s) => s.id">
            <template #cell-ext="{ row }"><code class="geo__ext">{{ ext(row) }}</code></template>
            <template #cell-updated="{ row }">
              <span class="geo__time">{{ row.lastUpdatedAt ? formatTimestamp(row.lastUpdatedAt) : '—' }}</span>
            </template>
            <template #cell-actions="{ row }">
              <div class="geo__actions">
                <XInlineButton icon="edit" :label="t('xray.geoEdit')" @click="edit(row)" />
                <XSpinner v-if="pulling === row.id || pulling === 'all'" :size="14" class="geo__busy" />
                <XInlineButton v-else icon="refresh" :label="t('xray.geoDownload')" :disabled="pulling !== null" @click="pull(row.id)" />
                <XInlineButton icon="delete" :label="t('xray.delete')" danger @click="deleting = row" />
              </div>
            </template>
            <template #empty>
              <XSpinner v-if="loading" :size="16" class="geo__wait" />
              <XEmpty v-else icon="globe" compact />
            </template>
          </XGridTable>
        </div>
      </XSection>
    </div>
  </XModal>

  <GeoSourceModal v-model:open="formOpen" :source="editing" @saved="load" />

  <XConfirm :open="asking !== null" icon="refresh" :title="t('xray.geoUpdateTitle')" :cta="t('xray.confirm')" @update:open="!$event && (asking = null)" @confirm="update">
    {{ asking ? t('xray.geoUpdateFile', { file: asking }) : t('xray.geoUpdateAllText') }}
  </XConfirm>
  <XConfirm v-model:open="deleteOpen" icon="delete" :title="t('xray.delete')" :cta="t('xray.confirm')" @confirm="remove">
    {{ t('xray.geoDeleteText') }}
  </XConfirm>
</template>

<style scoped>
.geo {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.geo__block {
  display: flex;
  flex-direction: column;
  gap: 10px;
  min-width: 0;
}

.geo__file {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  height: 40px;
  padding: 0 6px 0 12px;
}

.geo__busy {
  margin: 0 7px;
  color: var(--acc);
}

.geo__end {
  display: flex;
  justify-content: flex-end;
}

.geo__hint {
  display: flex;
  align-items: flex-start;
  gap: 7px;
  margin: 0;
  color: var(--mut);
  font: 400 var(--fs-12) / var(--lh-text) var(--font-sans);
}

.geo__hint > :first-child {
  flex: none;
  margin-top: 2px;
  color: var(--dim);
}

.geo__tools {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.geo__ext {
  padding: 2px 6px;
  border-radius: 5px;
  background: var(--hov);
  color: var(--fg2);
  font: 400 var(--fs-11) / var(--lh-tight) var(--font-mono);
  white-space: nowrap;
}

.geo__time {
  color: var(--mut);
  font: 400 var(--fs-11) / 22px var(--font-mono);
  white-space: nowrap;
}

.geo__actions {
  display: flex;
  align-items: center;
  gap: 2px;
}

.geo__wait {
  display: block;
  margin: 40px auto;
  color: var(--dim);
}
</style>
