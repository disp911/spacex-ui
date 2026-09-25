<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { api } from '../lib/api'
import { confirm } from '../lib/confirm'
import { dateTime, version as fmtVersion } from '../lib/format'
import { t, tn } from '../lib/i18n'
import { isMobile } from '../lib/media'
import { toast } from '../lib/toast'
import Icon from '../components/Icon.vue'
import Modal from '../components/Modal.vue'
import Spinner from '../components/Spinner.vue'
import CustomGeoModal, { type CustomGeo } from './CustomGeoModal.vue'

// Xray core versions, the bundled geofiles and user geo sources, as three
// accordion sections with one open at a time.
const props = defineProps<{ current: string }>()
const emit = defineEmits<{ close: [] }>()

type Section = 'xray' | 'geo' | 'custom' | ''
const section = ref<Section>('xray')
const toggle = (s: Section) => (section.value = section.value === s ? '' : s)

const installed = computed(() => fmtVersion(props.current))

// ───── Xray ─────
const releases = ref<{ version: string; publishedAt: number }[] | null>(null)
async function loadReleases() {
  const msg = await api.get<{ version: string; publishedAt: number }[]>('panel/api/server/getXrayReleases')
  releases.value = msg.success && Array.isArray(msg.obj) ? msg.obj : []
}
const day = (ts: number) => dateTime(ts).slice(0, 10)

function pickVersion(v: string) {
  if (v === installed.value) return
  confirm({
    title: t('spx.xraySwitchTitle'),
    text: t('spx.xraySwitchText', { version: v, current: installed.value || '—' }),
    cta: t('spx.updateXray'),
    busy: t('spx.installing'),
    action: async () => {
      const msg = await api.post(`panel/api/server/installXray/${encodeURIComponent(v)}`)
      if (msg.success) emit('close')
    },
  })
}

// ───── Geofiles ─────
const geofiles = ref<{ name: string; updatedAt: number }[]>([])
async function loadGeofiles() {
  const msg = await api.get<{ name: string; updatedAt: number }[]>('panel/api/server/geofiles', { notify: false })
  if (msg.success && Array.isArray(msg.obj)) geofiles.value = msg.obj
}

function updateGeofile(name?: string) {
  confirm({
    title: name ? t('spx.geofileUpdateTitle', { name }) : t('spx.geofilesUpdateTitle'),
    text: name ? t('spx.geofileUpdateText') : t('spx.geofilesUpdateText', { n: geofiles.value.length }),
    cta: name ? t('update') : t('pages.index.geofilesUpdateAll'),
    busy: t('spx.updating'),
    action: async () => {
      await api.post(name ? `panel/api/server/updateGeofile/${encodeURIComponent(name)}` : 'panel/api/server/updateGeofile')
      await loadGeofiles()
    },
  })
}

// ───── Custom sources ─────
const sources = ref<CustomGeo[] | null>(null)
const editing = ref<CustomGeo | 'new' | null>(null)
const rowBusy = ref<number | null>(null)

async function loadSources() {
  const msg = await api.get<CustomGeo[]>('panel/api/custom-geo/list', { notify: false })
  sources.value = msg.success && Array.isArray(msg.obj) ? msg.obj : []
}
const fileOf = (s: CustomGeo) => `${s.type}_${s.alias}.dat`
const route = (s: CustomGeo) => `ext:${fileOf(s)}:tag`

function refreshSource(s: CustomGeo) {
  confirm({
    title: t('spx.sourceUpdateTitle', { name: s.alias }),
    text: t('spx.sourceUpdateText'),
    cta: t('update'),
    busy: t('spx.updating'),
    action: async () => {
      rowBusy.value = s.id
      try {
        await api.post(`panel/api/custom-geo/download/${s.id}`)
        await loadSources()
      } finally {
        rowBusy.value = null
      }
    },
  })
}

function deleteSource(s: CustomGeo) {
  confirm({
    title: t('spx.sourceDeleteTitle', { name: s.alias }),
    text: t('spx.sourceDeleteText', { file: fileOf(s) }),
    cta: t('delete'),
    busy: t('spx.deleting'),
    tone: 'danger',
    action: async () => {
      await api.post(`panel/api/custom-geo/delete/${s.id}`)
      await loadSources()
    },
  })
}

function updateAllSources() {
  confirm({
    title: t('spx.sourcesUpdateTitle'),
    text: t('spx.sourcesUpdateText'),
    cta: t('pages.index.geofilesUpdateAll'),
    busy: t('spx.updating'),
    action: async () => {
      const msg = await api.post('panel/api/custom-geo/update-all', undefined, { notify: false })
      toast(msg.msg || t('pages.index.customGeoToastUpdateAll'), msg.success ? 'success' : 'error')
      await loadSources()
    },
  })
}

function onSaved() {
  editing.value = null
  loadSources()
}

const sourcesCount = computed(() =>
  sources.value === null ? '' : sources.value.length ? tn('spx.sourcesCount', sources.value.length) : t('spx.none'),
)

onMounted(() => {
  loadReleases()
  loadGeofiles()
  loadSources()
})
</script>

<template>
  <Modal :title="t('pages.index.xrayUpdates')" :subtitle="t('spx.updatesSubtitle')" :width="640" :closable="true" @close="emit('close')">
    <div class="sections">
      <!-- Xray core -->
      <div>
        <button class="sec-head" :class="{ open: section === 'xray' }" @click="toggle('xray')">
          <span class="sec-caret">▸</span>
          <span class="sec-name">Xray</span>
          <span class="sec-meta">{{ installed }}</span>
        </button>
        <div v-if="section === 'xray'" class="sec-body">
          <div class="notice warn">
            <Icon name="alert" />
            <div>
              <div class="notice-title">{{ t('spx.important') }}</div>
              <div class="notice-text">{{ t('spx.xraySwitchWarning') }}</div>
            </div>
          </div>
          <div class="list-box">
            <div class="list">
              <div v-if="releases === null" class="empty-inline"><Spinner /></div>
              <div v-else-if="!releases.length" class="empty-inline">{{ t('noData') }}</div>
              <template v-for="(r, i) in releases ?? []" :key="r.version">
                <div v-if="i" class="sep" />
                <button class="ver" :class="{ cur: r.version === installed }" @click="pickVersion(r.version)">
                  <span class="ver-name">{{ r.version }}</span>
                  <span v-if="r.version === installed" class="ver-chip">{{ t('spx.installed') }}</span>
                  <span class="ver-date">{{ day(r.publishedAt) }}</span>
                </button>
              </template>
            </div>
          </div>
        </div>
      </div>

      <!-- Geofiles -->
      <div>
        <button class="sec-head" :class="{ open: section === 'geo' }" @click="toggle('geo')">
          <span class="sec-caret">▸</span>
          <span class="sec-name">Geofiles</span>
          <span class="sec-meta">{{ tn('spx.filesCount', geofiles.length) }}</span>
        </button>
        <div v-if="section === 'geo'" class="sec-body">
          <div class="bar">
            <div class="notice info">
              <Icon name="info" :size="15" />
              <div class="notice-text muted">{{ t('spx.geofilesHint') }} <span class="mono hl">bin/</span></div>
            </div>
            <button class="ctl small" @click="updateGeofile()">{{ t('pages.index.geofilesUpdateAll') }}</button>
          </div>
          <div class="table">
            <div class="tr th">
              <span class="grow">{{ t('spx.file') }}</span>
              <span class="w-date">{{ t('pages.index.customGeoLastUpdated') }}</span>
              <span class="w-act right">{{ t('update') }}</span>
            </div>
            <template v-for="(g, i) in geofiles" :key="g.name">
              <div v-if="i" class="sep wide" />
              <div class="tr">
                <span class="grow mono strong ellipsis">{{ g.name }}</span>
                <span class="w-date mono dim">{{ g.updatedAt ? dateTime(g.updatedAt) : '—' }}</span>
                <span class="w-act right">
                  <button class="mini" :title="t('pages.index.customGeoDownload')" @click="updateGeofile(g.name)">
                    <Icon name="restart" :size="13" />
                  </button>
                </span>
              </div>
            </template>
          </div>
        </div>
      </div>

      <!-- Custom sources -->
      <div>
        <button class="sec-head" :class="{ open: section === 'custom' }" @click="toggle('custom')">
          <span class="sec-caret">▸</span>
          <span class="sec-name">{{ isMobile ? t('spx.customGeoShort') : t('pages.index.customGeoTitle') }}</span>
          <span class="sec-meta">{{ sourcesCount }}</span>
        </button>
        <div v-if="section === 'custom'" class="sec-body">
          <div class="bar">
            <div class="notice info">
              <Icon name="info" :size="15" />
              <div class="notice-text muted">{{ t('spx.customGeoHint') }} <span class="mono accent">ext:{{ t('spx.fileTag') }}</span></div>
            </div>
            <div class="bar-actions">
              <button v-if="sources?.length" class="ctl small" @click="updateAllSources">{{ t('pages.index.geofilesUpdateAll') }}</button>
              <button class="ctl small primary" @click="editing = 'new'"><Icon name="plus" :size="13" />{{ t('pages.index.customGeoAdd') }}</button>
            </div>
          </div>
          <div class="table">
            <template v-if="!isMobile">
              <div class="tr th center">
                <span class="w-type">{{ t('pages.index.customGeoType') }}</span>
                <span class="w-alias">{{ t('pages.index.customGeoAlias') }}</span>
                <span class="grow">{{ t('spx.routing') }}</span>
                <span class="w-upd">{{ t('pages.index.customGeoLastUpdated') }}</span>
                <span class="w-acts">{{ t('pages.index.customGeoActions') }}</span>
              </div>
            </template>
            <div v-if="sources === null" class="empty-inline"><Spinner /></div>
            <template v-for="(s, i) in sources ?? []" :key="s.id">
              <div v-if="i || !isMobile" class="sep wide" :class="{ hidden: !i }" />
              <div class="tr center" :class="{ card: isMobile }">
                <span class="w-type"><span class="type" :class="s.type">{{ s.type }}</span></span>
                <span class="w-alias mono strong ellipsis">{{ s.alias }}</span>
                <span class="grow mono accent ellipsis route">{{ route(s) }}</span>
                <span class="w-upd mono dim">{{ s.lastUpdatedAt ? dateTime(s.lastUpdatedAt) : '—' }}</span>
                <span class="w-acts acts">
                  <button class="mini" :title="t('pages.index.customGeoEdit')" @click="editing = s"><Icon name="edit" :size="13" /></button>
                  <button class="mini" :title="t('pages.index.customGeoDownload')" @click="refreshSource(s)">
                    <Spinner v-if="rowBusy === s.id" :size="12" /><Icon v-else name="restart" :size="13" />
                  </button>
                  <button class="mini danger" :title="t('pages.index.customGeoDelete')" @click="deleteSource(s)"><Icon name="trash" :size="13" /></button>
                </span>
              </div>
            </template>
            <div v-if="sources && !sources.length" class="empty">
              <div class="empty-icon"><Icon name="file" /></div>
              <div class="empty-title">{{ t('spx.noSources') }}</div>
              <div class="empty-text">{{ t('spx.noSourcesHint') }}</div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <template #footer>
      <button class="ctl cancel" @click="emit('close')">{{ t('close') }}</button>
    </template>
  </Modal>

  <CustomGeoModal v-if="editing" :source="editing === 'new' ? null : editing" @close="editing = null" @saved="onSaved" />
</template>

<style scoped>
.sections {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.sec-head {
  width: 100%;
  display: flex;
  align-items: center;
  gap: 11px;
  height: 44px;
  padding: 0 13px;
  border-radius: 12px;
  background: var(--fill);
  border: 1px solid var(--border);
  color: var(--text-2);
  text-align: left;
  transition:
    background 0.18s ease,
    border-color 0.18s ease;
}
.sec-head:hover {
  background: var(--row-hover);
}
.sec-head.open {
  background: var(--chip);
  border-color: rgba(var(--accent-rgb), 0.28);
  color: var(--text);
}
.sec-caret {
  font: 400 15px/1 var(--font-mono);
  width: 14px;
  display: flex;
  justify-content: center;
  color: var(--muted);
  transition: transform 0.18s ease;
}
.open .sec-caret {
  color: var(--accent);
  transform: rotate(90deg);
}
.sec-name {
  flex: 1;
  min-width: 0;
  font: 500 13.5px var(--font-sans);
}
.sec-meta {
  flex: none;
  white-space: nowrap;
  font: 500 11.5px var(--font-mono);
  color: var(--text-3b);
}
.sec-body {
  margin-top: 8px;
  display: flex;
  flex-direction: column;
  gap: 10px;
  animation: fadeIn 0.18s ease-out;
}

.notice {
  display: flex;
  gap: 10px;
  align-items: flex-start;
  border-radius: 11px;
  padding: 11px 12px;
}
.notice.warn {
  background: rgba(var(--warn-rgb), 0.1);
  border: 1px solid rgba(var(--warn-rgb), 0.3);
  color: var(--warn);
}
:root[data-theme='light'] .notice.warn {
  color: var(--warn-title);
}
.notice.info {
  flex: 1;
  min-width: 0;
  align-items: center;
  gap: 9px;
  background: var(--fill);
  border: 1px solid var(--border);
  padding: 9px 12px;
  color: var(--info);
}
.notice-title {
  font: 600 12px var(--font-sans);
  color: var(--warn-title);
  margin-bottom: 2px;
}
.notice-text {
  font: 400 11.5px/1.5 var(--font-sans);
  color: var(--warn-body);
}
.notice-text.muted {
  line-height: 1.4;
  color: var(--text-3);
}
.hl {
  color: var(--text-2);
}
.accent {
  color: var(--accent);
}

.list-box {
  background: var(--fill);
  border: 1px solid var(--border);
  border-radius: 12px;
  padding: 5px 3px 5px 0;
  clip-path: inset(0 round 12px);
}
.list {
  max-height: 258px;
  overflow-y: auto;
  padding: 0 2px 0 5px;
}
.sep {
  height: 1px;
  margin: 0 10px;
  background: var(--line-soft);
}
.sep.wide {
  margin: 0 13px;
}
.sep.hidden {
  display: none;
}
.ver {
  width: 100%;
  display: flex;
  align-items: center;
  gap: 10px;
  height: 40px;
  padding: 0 10px;
  border-radius: 9px;
  text-align: left;
  transition: background 0.16s ease;
}
.ver:hover {
  background: var(--row-hover-soft);
}
.ver.cur {
  background: var(--accent-soft);
  cursor: default;
}
.ver-name {
  flex: 1;
  min-width: 0;
  font: 500 12.5px var(--font-mono);
  color: var(--text-2);
}
.cur .ver-name {
  color: var(--accent);
}
.ver-chip {
  flex: none;
  font: 500 9.5px var(--font-mono);
  letter-spacing: 0.06em;
  color: var(--accent);
  border: 1px solid rgba(var(--accent-rgb), 0.34);
  border-radius: 6px;
  padding: 2px 6px;
}
.ver-date {
  flex: none;
  font: 400 11px var(--font-mono);
  color: var(--muted);
}
.empty-inline {
  display: flex;
  justify-content: center;
  padding: 18px;
  color: var(--muted);
  font: 400 12px var(--font-sans);
}

.bar {
  display: flex;
  align-items: center;
  gap: 12px;
}
.bar-actions {
  flex: none;
  display: flex;
  gap: 8px;
}
.ctl.small {
  height: 30px;
  padding: 0 12px;
  border-radius: 9px;
  background: var(--chip);
  font: 500 12px var(--font-sans);
  gap: 6px;
}
.ctl.small.primary {
  background: var(--accent);
  font-weight: 600;
}

.table {
  background: var(--fill);
  border: 1px solid var(--border);
  border-radius: 12px;
  overflow: hidden;
}
.tr {
  display: flex;
  align-items: center;
  gap: 10px;
  height: 44px;
  padding: 0 13px;
}
.tr.th {
  height: auto;
  padding: 9px 13px;
  border-bottom: 1px solid var(--border);
  font: 500 9.5px var(--font-mono);
  letter-spacing: 0.07em;
  color: var(--muted);
  white-space: nowrap;
  text-transform: uppercase;
}
.tr.center > span {
  text-align: center;
}
.grow {
  flex: 1;
  min-width: 0;
}
.w-date {
  width: 150px;
  flex: none;
}
.w-act {
  width: 72px;
  flex: none;
  display: flex;
  justify-content: flex-end;
}
.th .w-act {
  display: block;
}
.right {
  text-align: right;
}
.w-type {
  width: 62px;
  flex: none;
}
.w-alias {
  width: 96px;
  flex: none;
}
.w-upd {
  width: 104px;
  flex: none;
}
.w-acts {
  width: 88px;
  flex: none;
}
.acts {
  display: flex;
  gap: 6px;
  justify-content: center;
}
.mono {
  font: 500 12px var(--font-mono);
}
.strong {
  color: var(--text);
}
.dim {
  font: 400 11px var(--font-mono);
  color: var(--text-3b);
  white-space: nowrap;
}
.route {
  font-size: 11.5px;
  font-weight: 400;
}
.ellipsis {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.type {
  display: block;
  font: 500 10px var(--font-mono);
  letter-spacing: 0.04em;
  border-radius: 6px;
  padding: 3px 0;
  text-align: center;
  color: var(--accent);
  border: 1px solid rgba(var(--accent-rgb), 0.3);
}
.type.geoip {
  color: var(--info);
  border-color: rgba(var(--info-rgb), 0.3);
}
.mini {
  height: 26px;
  width: 26px;
  border-radius: 8px;
  border: 1px solid var(--line-control);
  background: var(--chip);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-3b);
  transition:
    box-shadow 0.22s ease,
    border-color 0.22s ease,
    background 0.22s ease,
    color 0.22s ease;
}
.mini:hover {
  border-color: var(--glow-border);
  background: var(--accent-fill);
  color: var(--text);
}
.mini.danger {
  border-color: rgba(var(--danger-rgb), 0.28);
  background: rgba(var(--danger-rgb), 0.08);
  color: var(--danger-text);
}
.mini.danger:hover {
  border-color: rgba(var(--danger-rgb), 0.5);
  background: rgba(var(--danger-rgb), 0.16);
  box-shadow: 0 0 0 3px rgba(var(--danger-rgb), 0.12);
}
.empty {
  padding: 30px 20px 32px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  text-align: center;
}
.empty-icon {
  width: 38px;
  height: 38px;
  border-radius: 12px;
  border: 1px dashed var(--scroll-thumb);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--faint);
  margin-bottom: 4px;
}
.empty-title {
  font: 500 12.5px var(--font-sans);
  color: var(--text-3);
}
.empty-text {
  font: 400 11.5px/1.5 var(--font-sans);
  color: var(--muted);
  max-width: 320px;
}
.cancel {
  background: var(--chip);
  padding: 0 15px;
}

/* Phones: rows become stacked cards. */
@media (max-width: 760px) {
  .bar {
    flex-direction: column;
    align-items: stretch;
  }
  .bar-actions .ctl {
    flex: 1;
    height: 38px;
  }
  .bar > .ctl {
    height: 38px;
  }
  .w-date {
    width: auto;
  }
  .tr.card {
    height: auto;
    flex-wrap: wrap;
    padding: 12px 13px;
    gap: 6px 10px;
  }
  .tr.card > span {
    text-align: left;
  }
  .tr.card .w-alias {
    width: auto;
    flex: 1;
  }
  .tr.card .route {
    order: 5;
    flex-basis: 100%;
  }
  .tr.card .w-upd {
    order: 6;
    width: auto;
    flex: 1;
  }
  .tr.card .acts {
    order: 7;
    width: auto;
  }
}
</style>
