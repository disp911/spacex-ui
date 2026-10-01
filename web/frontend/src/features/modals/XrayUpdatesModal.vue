<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue'
import XModal from '@/components/ui/XModal.vue'
import XAccordion from '@/components/ui/XAccordion.vue'
import XAlert from '@/components/ui/XAlert.vue'
import XButton from '@/components/ui/XButton.vue'
import XTag from '@/components/ui/XTag.vue'
import XEmpty from '@/components/ui/XEmpty.vue'
import XSpinner from '@/components/ui/XSpinner.vue'
import XConfirm from '@/components/ui/XConfirm.vue'
import SourceFormModal from './SourceFormModal.vue'
import { GEOFILES, customGeo, server, type CustomGeo } from '@/api/panel'
import type { Msg } from '@/api/http'
import { toastError, toastMsg } from '@/composables/useToast'
import { formatDateTime, withV } from '@/utils/format'
import { t } from '@/i18n'

const open = defineModel<boolean>('open', { default: false })
const props = defineProps<{ installed?: string; section?: 'xray' | 'geofiles' | 'custom' }>()
const emit = defineEmits<{ changed: [] }>()

const section = ref<'xray' | 'geofiles' | 'custom' | ''>('xray')
const versions = ref<string[]>([])
const versionsLoading = ref(false)
const sources = ref<CustomGeo[]>([])

async function loadVersions() {
  versionsLoading.value = true
  try {
    versions.value = (await server.xrayVersions()) ?? []
  } catch (e) {
    toastError(e)
  } finally {
    versionsLoading.value = false
  }
}
async function loadSources() {
  try {
    sources.value = (await customGeo.list()) ?? []
  } catch (e) {
    toastError(e)
  }
}

watch(open, (v) => {
  if (!v) return
  section.value = props.section ?? 'xray'
  loadVersions()
  loadSources()
})

const toggle = (s: 'xray' | 'geofiles' | 'custom') => (section.value = section.value === s ? '' : s)
const isInstalled = (v: string) => withV(v) === withV(props.installed ?? '')

// ---- Confirm dialog ---------------------------------------------------
interface Pending {
  title: string
  text: string
  cta: string
  tone?: 'warning' | 'danger'
  busyLabel?: string
  run: () => Promise<Msg<unknown>>
  after?: () => void
}
const confirm = reactive({ open: false, busy: false, pending: null as Pending | null })

function ask(p: Pending) {
  confirm.pending = p
  confirm.busy = false
  confirm.open = true
}
async function runConfirm() {
  const p = confirm.pending
  if (!p || confirm.busy) return
  confirm.busy = true
  try {
    const m = await p.run()
    toastMsg(m)
    if (m.success) p.after?.()
  } catch (e) {
    toastError(e)
  } finally {
    confirm.busy = false
    confirm.open = false
  }
}

function pickVersion(v: string) {
  if (isInstalled(v)) return
  ask({
    title: t('updates.confirmXrayTitle'),
    text: t('updates.confirmXrayText', { v: withV(v), cur: withV(props.installed ?? '') || '—' }),
    cta: t('updates.confirmXrayCta'),
    run: () => server.installXray(v),
    after: () => emit('changed'),
  })
}
const askGeo = (name?: string) =>
  ask(
    name
      ? { title: t('updates.confirmGeoTitle', { name }), text: t('updates.confirmGeoText'), cta: t('common.refresh'), run: () => server.updateGeofile(name) }
      : { title: t('updates.confirmAllGeoTitle'), text: t('updates.confirmAllGeoText'), cta: t('common.updateAll'), run: () => server.updateGeofile() },
  )
const askAllCustom = () =>
  ask({
    title: t('updates.confirmAllCustomTitle'),
    text: t('updates.confirmAllCustomText'),
    cta: t('common.updateAll'),
    run: () => customGeo.updateAll() as Promise<Msg<unknown>>,
    after: loadSources,
  })
const askSource = (s: CustomGeo) =>
  ask({ title: t('updates.confirmSrcTitle', { alias: s.alias }), text: t('updates.confirmSrcText'), cta: t('common.refresh'), run: () => customGeo.download(s.id), after: loadSources })
const askDelete = (s: CustomGeo) =>
  ask({
    title: t('updates.confirmDelTitle', { alias: s.alias }),
    text: t('updates.confirmDelText', { file: `${s.type}_${s.alias}.dat` }),
    cta: t('common.delete'),
    tone: 'danger',
    busyLabel: t('common.deleting'),
    run: () => customGeo.remove(s.id),
    after: loadSources,
  })

// ---- Source form -------------------------------------------------------
const form = reactive({ open: false, source: null as CustomGeo | null })
const openForm = (s: CustomGeo | null) => {
  form.source = s
  form.open = true
}

const sourcesMeta = computed(() => (sources.value.length ? t('updates.sources', { n: sources.value.length }) : t('updates.sourcesNone')))
const fmtTs = (ts: number) => (ts ? formatDateTime(new Date(ts * 1000)) : t('common.never'))
</script>

<template>
  <XModal v-model:open="open" :title="t('updates.title')" :subtitle="t('updates.subtitle')" :width="640">
    <div class="sections">
      <!-- Xray core -->
      <XAccordion :title="'Xray'" :meta="withV(installed ?? '')" :open="section === 'xray'" @toggle="toggle('xray')">
        <XAlert tone="warning" :title="t('updates.important')">{{ t('updates.versionWarn') }}</XAlert>
        <div class="list">
          <div class="list__scroll">
            <div v-if="versionsLoading" class="list__loading"><XSpinner /></div>
            <div v-else-if="!versions.length" class="list__loading">{{ t('common.nothingFound') }}</div>
            <button
              v-for="v in versions"
              v-else
              :key="v"
              type="button"
              class="ver"
              :class="{ 'is-current': isInstalled(v) }"
              :disabled="isInstalled(v)"
              @click="pickVersion(v)"
            >
              <span class="ver__name">{{ withV(v) }}</span>
              <XTag v-if="isInstalled(v)" tone="accent">{{ t('updates.installed') }}</XTag>
            </button>
          </div>
        </div>
      </XAccordion>

      <!-- Geo files -->
      <XAccordion :title="t('updates.geofiles')" :meta="t('updates.files', { n: GEOFILES.length })" :open="section === 'geofiles'" @toggle="toggle('geofiles')">
        <div class="bar">
          <XAlert tone="info" compact class="bar__info">
            {{ t('updates.geoInfo') }} <span class="mono em">bin/</span>
          </XAlert>
          <XButton size="sm" @click="askGeo()">{{ t('common.updateAll') }}</XButton>
        </div>
        <div class="table">
          <div class="table__head caps">
            <span class="c-grow">{{ t('updates.colFile') }}</span>
            <span class="c-act">{{ t('updates.colUpdate') }}</span>
          </div>
          <div v-for="g in GEOFILES" :key="g" class="table__row">
            <span class="c-grow mono strong ellipsis">{{ g }}</span>
            <span class="c-act">
              <XButton icon="restart" icon-only size="xs" :label="t('common.updateNow')" @click="askGeo(g)" />
            </span>
          </div>
        </div>
      </XAccordion>

      <!-- Custom sources -->
      <XAccordion :title="t('updates.custom')" :meta="sourcesMeta" :open="section === 'custom'" @toggle="toggle('custom')">
        <div class="bar">
          <XAlert tone="info" compact class="bar__info">
            {{ t('updates.customInfo') }} <span class="mono accent">ext:file.dat:tag</span>
          </XAlert>
          <div class="bar__actions">
            <XButton size="sm" :disabled="!sources.length" @click="askAllCustom">{{ t('common.updateAll') }}</XButton>
            <XButton size="sm" variant="primary" icon="plus" @click="openForm(null)">{{ t('common.add') }}</XButton>
          </div>
        </div>
        <div class="table">
          <template v-if="sources.length">
            <div class="table__head caps">
              <span class="c-type">{{ t('updates.colType') }}</span>
              <span class="c-alias">{{ t('updates.colAlias') }}</span>
              <span class="c-grow c-route">{{ t('updates.colRoute') }}</span>
              <span class="c-date">{{ t('updates.colUpdated') }}</span>
              <span class="c-acts">{{ t('updates.colActions') }}</span>
            </div>
            <div v-for="s in sources" :key="s.id" class="table__row">
              <span class="c-type"><XTag :tone="s.type === 'geosite' ? 'accent' : 'info'">{{ s.type }}</XTag></span>
              <span class="c-alias mono strong ellipsis">{{ s.alias }}</span>
              <span class="c-grow c-route mono accent ellipsis">ext:{{ s.type }}_{{ s.alias }}.dat:tag</span>
              <span class="c-date mono muted">{{ fmtTs(s.lastUpdatedAt) }}</span>
              <span class="c-acts">
                <XButton icon="edit" icon-only size="xs" :label="t('common.edit')" @click="openForm(s)" />
                <XButton icon="restart" icon-only size="xs" :label="t('common.updateNow')" @click="askSource(s)" />
                <XButton icon="trash" icon-only size="xs" variant="danger" :label="t('common.delete')" @click="askDelete(s)" />
              </span>
            </div>
          </template>
          <XEmpty v-else icon="doc" :title="t('updates.emptyTitle')" :text="t('updates.emptyText')" />
        </div>
      </XAccordion>
    </div>

    <template #footer>
      <XButton @click="open = false">{{ t('common.close') }}</XButton>
    </template>
  </XModal>

  <XConfirm
    v-model:open="confirm.open"
    :title="confirm.pending?.title ?? ''"
    :text="confirm.pending?.text"
    :cta="confirm.pending?.cta ?? ''"
    :tone="confirm.pending?.tone"
    :busy="confirm.busy"
    :busy-label="confirm.pending?.busyLabel"
    @confirm="runConfirm"
  />
  <SourceFormModal v-model:open="form.open" :source="form.source" @saved="loadSources" />
</template>

<style scoped>
.sections {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.list {
  padding: 5px 3px 5px 0;
  background: var(--field);
  border: 1px solid var(--border-field);
  border-radius: var(--r-md);
  overflow: hidden;
}
.list__scroll {
  max-height: 258px;
  overflow-y: auto;
  padding: 0 2px 0 5px;
}
.list__loading {
  display: flex;
  justify-content: center;
  padding: 18px;
  color: var(--text-4);
  font: var(--fw-regular) var(--fs-sm) var(--font-sans);
}
.ver {
  width: 100%;
  display: flex;
  align-items: center;
  gap: 10px;
  height: 40px;
  padding: 0 10px;
  border: 0;
  border-radius: var(--r-sm);
  background: transparent;
  color: var(--text-2);
  text-align: left;
  transition: background var(--dur-fast) ease;
}
.ver + .ver {
  box-shadow: inset 0 1px 0 var(--divider);
}
.ver:not(:disabled):hover {
  background: var(--hover);
}
.ver.is-current {
  background: var(--accent-soft);
  color: var(--accent);
  cursor: default;
  box-shadow: none;
}
.ver.is-current + .ver {
  box-shadow: none;
}
.ver__name {
  flex: 1;
  font: var(--fw-medium) var(--fs-md) var(--font-mono);
}

.bar {
  display: flex;
  align-items: center;
  gap: 12px;
}
.bar__info {
  flex: 1;
  min-width: 0;
}
.bar__actions {
  flex: none;
  display: flex;
  gap: 8px;
}

.table {
  background: var(--field);
  border: 1px solid var(--border-field);
  border-radius: var(--r-md);
  overflow: hidden;
}
.table__head,
.table__row {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 0 13px;
}
.table__head {
  height: 34px;
  border-bottom: 1px solid var(--border-field);
  font-size: var(--fs-2xs);
}
.table__row {
  height: 44px;
  font-size: var(--fs-sm);
}
.table__row + .table__row {
  border-top: 1px solid var(--divider);
}
.c-grow { flex: 1; min-width: 0; }
.c-act { width: 72px; flex: none; display: flex; justify-content: flex-end; }
.c-type { width: 72px; flex: none; }
.c-type :deep(.x-tag) { width: 100%; }
.c-alias { width: 96px; flex: none; }
.c-date { width: 118px; flex: none; font-size: var(--fs-xs); }
.c-acts { width: 92px; flex: none; display: flex; gap: 6px; justify-content: flex-end; }
.mono { font-family: var(--font-mono); }
.strong { font-weight: var(--fw-medium); color: var(--text); }
.accent { color: var(--accent); }
.muted { color: var(--text-3); }
.em { color: var(--text-2); }

@media (max-width: 767px) {
  .bar {
    flex-direction: column;
    align-items: stretch;
  }
  .bar__actions :deep(.x-btn) {
    flex: 1;
    height: var(--h-md);
  }
  .c-date,
  .c-route {
    display: none;
  }
  .c-alias {
    flex: 1;
    min-width: 0;
  }
  .list__scroll {
    max-height: 290px;
  }
}
</style>
