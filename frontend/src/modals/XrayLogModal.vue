<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { api, download } from '../lib/api'
import { t, tn } from '../lib/i18n'
import { isMobile } from '../lib/media'
import { toast } from '../lib/toast'
import Dropdown from '../components/Dropdown.vue'
import Icon from '../components/Icon.vue'
import Modal from '../components/Modal.vue'
import Spinner from '../components/Spinner.vue'

// Xray's access log as stored by the panel: one page per request, filtered
// server-side by day, client, text and connection kind.
const emit = defineEmits<{ close: [] }>()

interface Entry {
  DateTime: string
  FromAddress: string
  ToAddress: string
  Inbound: string
  Outbound: string
  Email: string
  Event: number // 0 direct, 1 blocked, 2 proxy
}
interface Page {
  entries: Entry[] | null
  total: number
  page: number
  pageSize: number
  date: string
  dates: string[] | null
  email: string
  clients: string[] | null
}

const data = ref<Page | null>(null)
const loading = ref(false)
const downloading = ref(false)
const filter = reactive({ date: '', email: '', query: '', page: 1 })
const kinds = reactive({ direct: true, blocked: true, proxy: true })
const kindsSheet = ref(false)

function body() {
  return {
    date: filter.date,
    email: filter.email,
    filter: filter.query.trim(),
    showDirect: kinds.direct,
    showBlocked: kinds.blocked,
    showProxy: kinds.proxy,
  }
}

let seq = 0
async function load() {
  const my = ++seq
  loading.value = true
  const msg = await api.post<Page>('panel/api/server/xraylogs', { ...body(), page: filter.page }, { notify: false })
  if (my !== seq) return
  loading.value = false
  if (msg.success && msg.obj) {
    data.value = msg.obj
    filter.date = msg.obj.date
    filter.page = msg.obj.page || 1
  } else if (!data.value) {
    data.value = { entries: [], total: 0, page: 1, pageSize: 200, date: '', dates: [], email: '', clients: [] }
  }
}

let debounce: number | undefined
watch(
  () => filter.query,
  () => {
    clearTimeout(debounce)
    debounce = window.setTimeout(() => {
      filter.page = 1
      load()
    }, 300)
  },
)
watch([() => filter.date, () => filter.email, () => ({ ...kinds })], (_n, o) => {
  if (!o) return
  filter.page = 1
  load()
}, { deep: true })

// dates arrive as YYYY-MM-DD; show them as DD.MM.YYYY
const dayLabel = (d: string) => (/^\d{4}-\d{2}-\d{2}$/.test(d) ? d.split('-').reverse().join('.') : d)
const dateOptions = computed(() => (data.value?.dates ?? []).map((d) => ({ value: d, label: dayLabel(d) })))
const clientOptions = computed(() => [
  { value: '', label: t('spx.allClients') },
  ...(data.value?.clients ?? []).map((c) => ({ value: c, label: c })),
])

const entries = computed(() => data.value?.entries ?? [])
const pageSize = computed(() => data.value?.pageSize || 200)
const pages = computed(() => Math.max(1, Math.ceil((data.value?.total ?? 0) / pageSize.value)))
const range = computed(() => {
  const total = data.value?.total ?? 0
  if (!total) return ''
  const from = (filter.page - 1) * pageSize.value + 1
  return t('spx.shownRange', { from, to: Math.min(total, from + entries.value.length - 1), total, page: filter.page, pages: pages.value })
})
const pageButtons = computed(() => {
  const n = pages.value
  const p = filter.page
  if (n <= 7) return Array.from({ length: n }, (_, i) => i + 1)
  const set = new Set([1, n, p - 1, p, p + 1].filter((x) => x >= 1 && x <= n))
  const out: (number | '…')[] = []
  ;[...set].sort((a, b) => a - b).forEach((x, i, arr) => {
    if (i && x - arr[i - 1] > 1) out.push('…')
    out.push(x)
  })
  return out
})
function go(p: number) {
  if (p < 1 || p > pages.value || p === filter.page) return
  filter.page = p
  load()
}

const pad = (n: number) => String(n).padStart(2, '0')
function when(iso: string) {
  const d = new Date(iso)
  if (isNaN(d.getTime())) return { day: '', time: iso }
  return { day: `${pad(d.getDate())}.${pad(d.getMonth() + 1)}`, time: `${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())}` }
}
const kindOf = (e: Entry) => (e.Event === 1 ? 'blocked' : e.Event === 2 ? 'proxy' : 'direct')
// Xray tags its own traffic xray.system.<uuid>; the uuid only widens the column.
const inboundOf = (e: Entry) => (e.Inbound.startsWith('xray.system.') ? 'xray.system' : e.Inbound)

const kindsLabel = computed(() => {
  const on = (['direct', 'blocked', 'proxy'] as const).filter((k) => kinds[k])
  return on.length === 3 ? t('spx.allTypes') : on.length ? on.map((k) => k[0].toUpperCase() + k.slice(1)).join(', ') : t('spx.none')
})

async function save() {
  downloading.value = true
  try {
    await download('panel/api/server/xraylogs/download', body(), 'xray.log')
  } catch (e) {
    toast(e instanceof Error ? e.message : String(e), 'error')
  } finally {
    downloading.value = false
  }
}

onMounted(load)
</script>

<template>
  <Modal :title="isMobile ? t('spx.connectionsLogShort') : t('spx.connectionsLog')" :width="1320" :height="820" body-class="xl-body" @close="emit('close')">
    <template #subtitle>
      <template v-if="data">{{ tn('spx.recordsFor', data.total, { date: dayLabel(data.date) }) }}</template>
    </template>

    <template #toolbar>
      <template v-if="!isMobile">
        <button class="ctl icon" :title="t('spx.refresh')" @click="load">
          <Icon name="restart" :size="15" :class="{ spinning: loading }" />
        </button>
        <Dropdown v-model="filter.date" :options="dateOptions" :label="t('spx.date')" :title="t('spx.date')" />
        <Dropdown v-model="filter.email" :options="clientOptions" :label="t('spx.client')" :title="t('spx.client')" :menu-width="240" />
        <label class="search">
          <Icon name="search" :size="13" />
          <input v-model="filter.query" :placeholder="t('spx.searchPlaceholder')" spellcheck="false" />
          <button v-if="filter.query" class="clear" :title="t('spx.clear')" @click="filter.query = ''"><Icon name="close" :size="12" /></button>
        </label>
        <button v-for="k in (['direct', 'blocked', 'proxy'] as const)" :key="k" class="kind" :class="[k, { on: kinds[k] }]" @click="kinds[k] = !kinds[k]">
          <span class="kbox"><Icon v-if="kinds[k]" name="check" :size="11" /></span>{{ k[0].toUpperCase() + k.slice(1) }}
        </button>
        <button class="ctl dl" :disabled="downloading" @click="save">
          <Spinner v-if="downloading" :size="14" /><Icon v-else name="download" :size="14" />{{ t('download') }}
        </button>
      </template>
      <div v-else class="mfilters">
        <label class="search">
          <Icon name="search" :size="13" />
          <input v-model="filter.query" :placeholder="t('spx.searchPlaceholderShort')" spellcheck="false" />
          <button v-if="filter.query" class="clear" @click="filter.query = ''"><Icon name="close" :size="12" /></button>
        </label>
        <div class="mrow">
          <Dropdown v-model="filter.date" :options="dateOptions" :title="t('spx.date')" block />
          <Dropdown v-model="filter.email" :options="clientOptions" :title="t('spx.client')" block />
          <button class="ctl mkinds" @click="kindsSheet = true"><span class="value">{{ kindsLabel }}</span><span class="caret">▼</span></button>
        </div>
      </div>
    </template>

    <div v-if="!data" class="state"><Spinner :size="20" /></div>
    <div v-else-if="!entries.length" class="state empty">
      <div class="empty-title">{{ t('spx.nothingFound') }}</div>
      <div class="empty-text">{{ t('spx.nothingFoundHint') }}</div>
    </div>

    <div v-else-if="!isMobile" class="table-wrap" :class="{ dim: loading }">
      <table>
        <thead>
          <tr>
            <th>{{ t('spx.date') }}</th>
            <th>{{ t('spx.from') }}</th>
            <th class="grow">{{ t('spx.to') }}</th>
            <th>{{ t('spx.inbound') }}</th>
            <th>{{ t('spx.outbound') }}</th>
            <th>Email</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(e, i) in entries" :key="i" :class="kindOf(e)">
            <td>{{ when(e.DateTime).day }} {{ when(e.DateTime).time }}</td>
            <td>{{ e.FromAddress }}</td>
            <td class="grow">{{ e.ToAddress }}</td>
            <td :title="e.Inbound">{{ inboundOf(e) }}</td>
            <td><span class="odot" />{{ e.Outbound }}</td>
            <td>{{ e.Email }}</td>
          </tr>
        </tbody>
      </table>
    </div>

    <div v-else class="cards" :class="{ dim: loading }">
      <div v-for="(e, i) in entries" :key="i" class="entry" :class="kindOf(e)">
        <div class="e-top">
          <span class="odot" />
          <span class="e-to">{{ e.ToAddress }}</span>
          <span class="e-time">{{ when(e.DateTime).time }}</span>
        </div>
        <div class="e-bottom">
          <span class="e-email">{{ e.Email || e.FromAddress }}</span>
          <span class="e-in">{{ inboundOf(e) }}</span>
          <span class="e-out">{{ e.Outbound }}</span>
        </div>
      </div>
    </div>

    <template v-if="data && pages > 0 && entries.length" #footer>
      <span class="range">{{ isMobile ? t('spx.pageOf', { page: filter.page, pages }) : range }}</span>
      <div class="pager">
        <button class="pg" :disabled="filter.page <= 1" @click="go(filter.page - 1)"><Icon name="chevronLeft" :size="13" /></button>
        <template v-if="!isMobile">
          <template v-for="(p, i) in pageButtons" :key="i">
            <span v-if="p === '…'" class="gap">…</span>
            <button v-else class="pg" :class="{ on: p === filter.page }" @click="go(p)">{{ p }}</button>
          </template>
        </template>
        <button class="pg" :disabled="filter.page >= pages" @click="go(filter.page + 1)"><Icon name="chevronRight" :size="13" /></button>
      </div>
    </template>
  </Modal>

  <Teleport v-if="kindsSheet" to="body">
    <div class="sheet-scrim" @click="kindsSheet = false" />
    <div class="sheet">
      <div class="grab" />
      <div class="sheet-title">{{ t('spx.connectionTypes') }}</div>
      <button v-for="k in (['direct', 'blocked', 'proxy'] as const)" :key="k" class="sheet-row" :class="k" @click="kinds[k] = !kinds[k]">
        <span class="odot" />
        <span class="sl">{{ k[0].toUpperCase() + k.slice(1) }}</span>
        <span class="kbox" :class="{ on: kinds[k] }"><Icon v-if="kinds[k]" name="check" :size="12" /></span>
      </button>
    </div>
  </Teleport>
</template>

<style scoped>
:deep(.xl-body) {
  display: flex;
  flex-direction: column;
  padding-top: 12px;
}
.spinning {
  animation: spin 0.7s linear infinite;
}
.search {
  flex: 1;
  min-width: 140px;
  height: 34px;
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 0 10px 0 12px;
  border-radius: 10px;
  border: 1px solid var(--line-control);
  background: var(--fill);
  color: var(--muted);
  transition:
    border-color 0.22s ease,
    box-shadow 0.22s ease;
}
.search:focus-within,
.search:hover {
  border-color: var(--glow-border);
  box-shadow: var(--glow);
}
.search input {
  flex: 1;
  min-width: 0;
  border: 0;
  outline: 0;
  background: transparent;
  font: 400 12px var(--font-mono);
  color: var(--text);
}
.search input::placeholder {
  color: var(--faint);
}
.clear {
  display: flex;
  color: var(--muted);
}
.kind {
  --k: var(--text-3);
  --k-rgb: var(--accent-rgb);
  height: 34px;
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 0 11px;
  border-radius: 10px;
  border: 1px solid var(--line-control);
  background: var(--fill);
  font: 500 12px var(--font-sans);
  color: var(--text-3);
  transition:
    border-color 0.2s ease,
    background 0.2s ease,
    color 0.2s ease;
}
.kind.direct {
  --k: var(--accent);
  --k-rgb: var(--accent-rgb);
}
.kind.blocked {
  --k: var(--danger-text);
  --k-rgb: var(--danger-rgb);
}
.kind.proxy {
  --k: var(--info);
  --k-rgb: var(--info-rgb);
}
.kind.on {
  color: var(--k);
  border-color: rgba(var(--k-rgb), 0.4);
  background: rgba(var(--k-rgb), 0.1);
}
.kbox {
  width: 15px;
  height: 15px;
  border-radius: 5px;
  border: 1px solid var(--scroll-thumb);
  display: flex;
  align-items: center;
  justify-content: center;
}
.kind.on .kbox {
  border-color: transparent;
  background: var(--k);
  color: var(--bg);
}
.dl {
  background: var(--chip);
}

.state {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 6px;
  color: var(--accent);
}
.empty-title {
  font: 500 13px var(--font-sans);
  color: var(--text-3);
}
.empty-text {
  font: 400 12px var(--font-sans);
  color: var(--muted);
}
.dim {
  opacity: 0.55;
  transition: opacity 0.2s ease;
}

.table-wrap {
  flex: 1;
  min-height: 0;
  overflow: auto;
  background: var(--fill);
  border: 1px solid var(--border);
  border-radius: 12px;
}
table {
  width: 100%;
  border-collapse: collapse;
  font: 400 12px var(--font-mono);
}
th {
  position: sticky;
  top: 0;
  z-index: 1;
  background: var(--chip);
  text-align: left;
  padding: 9px 14px;
  font: 500 9.5px var(--font-mono);
  letter-spacing: 0.07em;
  color: var(--muted);
  text-transform: uppercase;
  white-space: nowrap;
  border-bottom: 1px solid var(--border);
}
td {
  padding: 0 14px;
  height: 34px;
  white-space: nowrap;
  color: var(--text-2);
  border-bottom: 1px solid var(--line-soft);
}
td:first-child {
  color: var(--text-3b);
}
.grow {
  width: 100%;
}
.odot {
  display: inline-block;
  width: 6px;
  height: 6px;
  border-radius: 3px;
  margin-right: 8px;
  vertical-align: 1px;
  background: var(--faint);
  flex: none;
}
tr.proxy td {
  color: var(--info);
  background: rgba(var(--info-rgb), 0.04);
}
tr.blocked td {
  color: var(--danger-text);
  background: rgba(var(--danger-rgb), 0.05);
}
tr.proxy .odot,
.entry.proxy .odot,
.sheet-row.proxy .odot {
  background: var(--info);
}
tr.blocked .odot,
.entry.blocked .odot,
.sheet-row.blocked .odot {
  background: var(--danger);
}
.sheet-row.direct .odot {
  background: var(--accent);
}
tbody tr:hover td {
  background: var(--row-hover);
}

.range {
  margin-right: auto;
  font: 400 11.5px var(--font-sans);
  color: var(--text-3b);
}
.pager {
  display: flex;
  align-items: center;
  gap: 6px;
}
.pg {
  min-width: 26px;
  height: 26px;
  padding: 0 7px;
  border-radius: 8px;
  border: 1px solid var(--line-control);
  background: var(--chip);
  color: var(--text-2);
  font: 500 11.5px var(--font-mono);
  display: flex;
  align-items: center;
  justify-content: center;
  transition:
    border-color 0.18s ease,
    background 0.18s ease;
}
.pg:hover:not(:disabled) {
  border-color: var(--glow-border);
  background: var(--accent-fill);
}
.pg.on {
  background: var(--accent);
  border-color: transparent;
  color: var(--on-accent);
}
.pg:disabled {
  opacity: 0.4;
  cursor: default;
}
.gap {
  color: var(--muted);
  font: 400 11px var(--font-mono);
}

/* Phones */
.mfilters {
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.mfilters .search {
  flex: none;
  height: 36px;
}
.mrow {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 8px;
}
.mrow :deep(.trigger),
.mkinds {
  height: 34px;
}
.mkinds {
  justify-content: flex-start;
  background: var(--chip);
  padding: 0 9px 0 11px;
  min-width: 0;
}
.mkinds .value {
  font: 500 12px var(--font-sans);
  color: var(--text);
  overflow: hidden;
  text-overflow: ellipsis;
}
.mkinds .caret {
  margin-left: auto;
}
.cards {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.entry {
  background: var(--fill);
  border: 1px solid var(--border);
  border-radius: 12px;
  padding: 11px 13px;
  font-family: var(--font-mono);
}
.entry.blocked {
  background: rgba(var(--danger-rgb), 0.06);
  border-color: rgba(var(--danger-rgb), 0.25);
}
.e-top {
  display: flex;
  align-items: center;
  gap: 2px;
  margin-bottom: 7px;
}
.e-to {
  flex: 1;
  min-width: 0;
  font: 500 13.5px var(--font-mono);
  color: var(--text);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.entry.proxy .e-to {
  color: var(--info);
}
.entry.blocked .e-to {
  color: var(--danger-text);
}
.e-time {
  font: 400 11.5px var(--font-mono);
  color: var(--text-3b);
}
.e-bottom {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
}
.e-email {
  font: 450 12px var(--font-sans);
  color: var(--text-2);
}
.e-in {
  flex: 1;
  min-width: 0;
  font: 400 11px var(--font-mono);
  color: var(--muted);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.e-out {
  flex: none;
  font: 500 10.5px var(--font-mono);
  color: var(--text-2);
  border: 1px solid var(--border);
  border-radius: 6px;
  padding: 2px 7px;
}
.entry.proxy .e-out {
  color: var(--info);
  border-color: rgba(var(--info-rgb), 0.4);
}
.entry.blocked .e-out {
  color: var(--danger-text);
  border-color: rgba(var(--danger-rgb), 0.4);
}
.sheet-scrim {
  position: fixed;
  inset: 0;
  z-index: 300;
  background: var(--scrim-strong);
  backdrop-filter: blur(2px);
  animation: scrimIn 0.18s ease-out;
}
.sheet {
  position: fixed;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 301;
  background: var(--sheet);
  border-top: 1px solid var(--line-card);
  border-radius: 22px 22px 0 0;
  box-shadow: var(--shadow-sheet);
  padding: 8px 12px max(14px, env(safe-area-inset-bottom));
  animation: sheetUp 0.16s ease-out;
}
.grab {
  width: 38px;
  height: 4px;
  border-radius: 2px;
  background: var(--scroll-thumb);
  margin: 2px auto 10px;
}
.sheet-title {
  font: 600 15px var(--font-sans);
  padding: 2px 11px 10px;
}
.sheet-row {
  width: 100%;
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 11px;
  border-radius: 11px;
  color: var(--text-2);
}
.sheet-row .sl {
  flex: 1;
  text-align: left;
  font: 450 14px var(--font-sans);
}
.sheet-row .kbox.on {
  border-color: transparent;
  background: var(--accent);
  color: var(--on-accent);
}
</style>
