<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, reactive, ref, shallowRef, watch } from 'vue'
import AppShell from '@/components/shell/AppShell.vue'
import XButton from '@/components/ui/XButton.vue'
import XInput from '@/components/ui/XInput.vue'
import XSelect from '@/components/ui/XSelect.vue'
import XMenu from '@/components/ui/XMenu.vue'
import XCheckChip from '@/components/ui/XCheckChip.vue'
import XEmpty from '@/components/ui/XEmpty.vue'
import XSpinner from '@/components/ui/XSpinner.vue'
import XConfirm from '@/components/ui/XConfirm.vue'
import XIcon from '@/components/ui/XIcon.vue'
import type { MenuItem } from '@/components/ui/types'
import InboundList from '@/features/inbounds/InboundList.vue'
import QrModal from '@/features/inbounds/QrModal.vue'
import { PROTOCOL_LABELS, toRow, type ClientRow, type InboundRow } from '@/features/inbounds/model'
import { inbounds, settings, type DbInbound, type PanelDefaults } from '@/api/panel'
import type { Msg } from '@/api/http'
import { toastError, toastMsg } from '@/composables/useToast'
import { env, panelUrl } from '@/env'
import { sizeParts } from '@/utils/format'
import { t } from '@/i18n'

// ---- Data ------------------------------------------------------------
const raw = shallowRef<DbInbound[]>([])
const online = shallowRef(new Set<string>())
const defaults = ref<PanelDefaults | null>(null)
const loaded = ref(false)
const loading = ref(false)

async function load() {
  loading.value = true
  try {
    const [list, on] = await Promise.all([inbounds.list(), inbounds.onlines().catch(() => [])])
    raw.value = list ?? []
    online.value = new Set(on ?? [])
  } catch (e) {
    toastError(e)
  } finally {
    loading.value = false
    loaded.value = true
  }
}

onMounted(async () => {
  settings.defaults().then((d) => (defaults.value = d)).catch(() => {})
  await load()
})

const rows = computed(() => raw.value.map((db) => toRow(db, online.value, defaults.value)))

// ---- Auto refresh ----------------------------------------------------
const AUTO_KEY = 'spx-inbounds-auto'
const auto = ref(localStorage.getItem(AUTO_KEY) === '1')
let timer: number | undefined
watch(
  auto,
  (v) => {
    localStorage.setItem(AUTO_KEY, v ? '1' : '0')
    clearInterval(timer)
    if (v) timer = window.setInterval(() => !loading.value && load(), 10_000)
  },
  { immediate: true },
)
onBeforeUnmount(() => clearInterval(timer))

// ---- Filters ---------------------------------------------------------
type StateFilter = 'all' | 'on' | 'off' | 'attention'
const q = ref('')
const protocol = ref('')
const state = ref<StateFilter>('all')

const protocolOptions = computed(() => [
  { value: '', label: t('ib.allProtocols') },
  ...[...new Set(raw.value.map((r) => r.protocol))].sort().map((p) => ({ value: p, label: PROTOCOL_LABELS[p] ?? p })),
])
const stateOptions = computed(() => [
  { value: 'all' as StateFilter, label: t('ib.stateAll') },
  { value: 'on' as StateFilter, label: t('ib.stateOn'), dot: 'var(--accent)' },
  { value: 'off' as StateFilter, label: t('ib.stateOff'), dot: 'var(--muted-dot)' },
  { value: 'attention' as StateFilter, label: t('ib.stateAttention'), dot: 'var(--warning)' },
])

const expanded = reactive(new Set<number>())
const clientMatches = (r: InboundRow, s: string) => r.clients.some((c) => c.email.toLowerCase().includes(s))

const visible = computed(() => {
  const s = q.value.trim().toLowerCase()
  return rows.value.filter((r) => {
    if (protocol.value && r.protocol !== protocol.value) return false
    if (state.value === 'on' && !r.enable) return false
    if (state.value === 'off' && r.enable) return false
    if (state.value === 'attention' && r.health !== 'warn' && r.health !== 'depleted') return false
    if (!s) return true
    return r.remark.toLowerCase().includes(s) || String(r.port).includes(s) || r.protocol.includes(s) || clientMatches(r, s)
  })
})

// Searching for a client opens the inbounds that contain it.
watch(q, (v) => {
  const s = v.trim().toLowerCase()
  if (s.length < 2) return
  for (const r of rows.value) if (clientMatches(r, s)) expanded.add(r.id)
})
const toggleExpand = (id: number) => (expanded.has(id) ? expanded.delete(id) : expanded.add(id))

// ---- Summary ---------------------------------------------------------
const summary = computed(() => {
  const all = rows.value
  const clients = all.reduce((n, r) => n + r.clients.length, 0)
  const on = all.reduce((n, r) => n + r.online, 0)
  const up = all.reduce((n, r) => n + r.up, 0)
  const down = all.reduce((n, r) => n + r.down, 0)
  const total = all.reduce((n, r) => n + (r.allTime || r.up + r.down), 0)
  return {
    inbounds: all.length,
    enabled: all.filter((r) => r.enable).length,
    clients,
    online: on,
    total: sizeParts(total),
    up: sizeParts(up),
    down: sizeParts(down),
    expiring: all.reduce((n, r) => n + r.expiring, 0),
    depleted: all.reduce((n, r) => n + r.depleted, 0),
  }
})

// ---- Actions ---------------------------------------------------------
const busy = reactive(new Set<string>())

async function setInboundEnable(r: InboundRow, v: boolean) {
  const key = 'i' + r.id
  busy.add(key)
  try {
    const m = await inbounds.setEnable(r.id, v)
    if (m.success) raw.value = raw.value.map((x) => (x.id === r.id ? { ...x, enable: v } : x))
    else toastMsg(m)
  } catch (e) {
    toastError(e)
  } finally {
    busy.delete(key)
  }
}

async function setClientEnable(r: InboundRow, c: ClientRow, v: boolean) {
  busy.add(c.email)
  try {
    const m = await inbounds.updateClient(r.id, c.clientId, { ...c.raw, enable: v })
    toastMsg(m)
    if (m.success) await load()
  } catch (e) {
    toastError(e)
  } finally {
    busy.delete(c.email)
  }
}

interface Pending {
  title: string
  text: string
  cta: string
  tone?: 'warning' | 'danger'
  run: () => Promise<Msg<unknown>>
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
    if (m.success) await load()
  } catch (e) {
    toastError(e)
  } finally {
    confirm.busy = false
    confirm.open = false
  }
}

function inboundAction(key: string, r: InboundRow) {
  const name = r.remark || '#' + r.id
  if (key === 'delete')
    ask({ title: t('ib.confirmDelTitle', { name }), text: t('ib.confirmDelText'), cta: t('common.delete'), tone: 'danger', run: () => inbounds.remove(r.id) })
  else if (key === 'resetClients')
    ask({ title: t('ib.confirmResetClientsTitle', { name }), text: t('ib.confirmResetClientsText'), cta: t('ib.reset'), run: () => inbounds.resetClientTraffics(r.id) })
  else if (key === 'delDepleted')
    ask({ title: t('ib.confirmDepletedTitle'), text: t('ib.confirmDepletedText'), cta: t('common.delete'), tone: 'danger', run: () => inbounds.delDepletedClients(r.id) })
}

function clientAction(key: string, r: InboundRow, c: ClientRow) {
  if (key === 'reset')
    ask({ title: t('ib.confirmResetClientTitle', { email: c.email }), text: t('ib.confirmResetClientText'), cta: t('ib.reset'), run: () => inbounds.resetClientTraffic(r.id, c.email) })
  else if (key === 'delete')
    ask({ title: t('ib.confirmDelClientTitle', { email: c.email }), text: t('ib.confirmDelClientText'), cta: t('common.delete'), tone: 'danger', run: () => inbounds.removeClient(r.id, c.email) })
}

const legacy = panelUrl('panel/inbounds')
const pageMenu = computed<MenuItem[]>(() => [
  { key: 'resetAll', label: t('ib.resetAll'), icon: 'restart' },
  { key: 'delDepletedAll', label: t('ib.delDepletedAll'), icon: 'trash', disabled: !summary.value.depleted },
  { key: 'legacy', label: t('ib.legacyNote'), icon: 'external', href: legacy, divided: true },
])
function pageAction(key: string) {
  if (key === 'resetAll')
    ask({ title: t('ib.confirmResetAllTitle'), text: t('ib.confirmResetAllText'), cta: t('ib.reset'), tone: 'danger', run: () => inbounds.resetAllTraffics() })
  else if (key === 'delDepletedAll')
    ask({ title: t('ib.confirmDepletedTitle'), text: t('ib.confirmDepletedText'), cta: t('common.delete'), tone: 'danger', run: () => inbounds.delDepletedClients(-1) })
}

// ---- QR --------------------------------------------------------------
const qr = reactive({ open: false, inbound: null as InboundRow | null, client: null as ClientRow | null })
function showQr(r: InboundRow, c: ClientRow | null) {
  qr.inbound = r
  qr.client = c
  qr.open = true
}
</script>

<template>
  <AppShell :title="t('ib.title')" :subtitle="env.host">
    <!-- Summary -->
    <div class="sum">
      <div class="tile">
        <div class="caps">{{ t('ib.inbounds') }}</div>
        <div class="tile__value">{{ summary.inbounds }}<small> / {{ summary.enabled }} {{ t('ib.stateOn').toLowerCase() }}</small></div>
      </div>
      <div class="tile">
        <div class="caps">{{ t('ib.clients') }}</div>
        <div class="tile__value"><span class="accent">{{ summary.online }}</span><small> / {{ summary.clients }}</small></div>
      </div>
      <div class="tile">
        <div class="caps">{{ t('ib.traffic') }}</div>
        <div class="tile__value">{{ summary.total[0] }}<small> {{ summary.total[1] }}</small></div>
        <div class="tile__hint mono">
          <XIcon name="arrow-up" :size="11" :stroke-width="1.6" class="info" />{{ summary.up.join(' ') }}
          <XIcon name="arrow-down" :size="11" :stroke-width="1.6" class="accent" />{{ summary.down.join(' ') }}
        </div>
      </div>
      <button type="button" class="tile tile--btn" :class="{ 'is-warn': summary.expiring + summary.depleted > 0 }" @click="state = state === 'attention' ? 'all' : 'attention'">
        <div class="caps">{{ t('ib.attention') }}</div>
        <div class="tile__value">{{ summary.expiring + summary.depleted }}</div>
        <div class="tile__hint">{{ t('ib.attentionHint', { e: summary.expiring, d: summary.depleted }) }}</div>
      </button>
    </div>

    <!-- Toolbar -->
    <div class="tb">
      <XInput v-model="q" size="sm" icon="search" clearable :placeholder="t('ib.search')" class="tb__search" />
      <div class="tb__filters no-scrollbar">
        <XSelect v-model="protocol" :options="protocolOptions" :label="t('ib.protocol')" :sheet-title="t('ib.protocol')" :menu-width="180" />
        <XSelect v-model="state" :options="stateOptions" :label="t('ib.state')" :sheet-title="t('ib.state')" :menu-width="190" />
        <XCheckChip v-model="auto" :label="t('ib.autoRefresh')" />
        <XButton icon="restart" icon-only size="sm" :label="t('common.refresh')" :loading="loading" @click="load" />
      </div>
      <div class="tb__end">
        <XMenu :items="pageMenu" :label="t('ib.more')" :title="t('ib.title')" @select="pageAction" />
        <XButton variant="primary" size="sm" icon="plus" :href="legacy">{{ t('ib.add') }}</XButton>
      </div>
    </div>

    <!-- List -->
    <div v-if="!loaded" class="state"><XSpinner :size="22" /></div>
    <div v-else-if="!rows.length" class="state state--card">
      <XEmpty icon="inbounds" :title="t('ib.emptyTitle')" :text="t('ib.emptyText')" />
      <XButton variant="primary" icon="plus" :href="legacy">{{ t('ib.add') }}</XButton>
    </div>
    <div v-else-if="!visible.length" class="state state--card">
      <XEmpty :title="t('ib.notFound')" :text="t('ib.notFoundText')" />
    </div>
    <InboundList
      v-else
      :rows="visible"
      :expanded="expanded"
      :busy="busy"
      :show-ip="defaults?.ipLimitEnable"
      :filter="q"
      @expand="toggleExpand"
      @toggle="setInboundEnable"
      @qr="showQr"
      @action="inboundAction"
      @client-toggle="setClientEnable"
      @client-action="clientAction"
    />

    <XConfirm
      v-model:open="confirm.open"
      :title="confirm.pending?.title ?? ''"
      :text="confirm.pending?.text"
      :cta="confirm.pending?.cta ?? ''"
      :tone="confirm.pending?.tone"
      :busy="confirm.busy"
      :busy-label="confirm.pending?.tone === 'danger' ? t('common.deleting') : t('common.updating')"
      @confirm="runConfirm"
    />
    <QrModal v-model:open="qr.open" :inbound="qr.inbound" :client="qr.client" :defaults="defaults" />
  </AppShell>
</template>

<style scoped>
.sum {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 12px;
  margin-bottom: 14px;
}
.tile {
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 6px;
  padding: 13px 16px 14px;
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--r-lg);
  text-align: left;
  color: var(--text);
}
.tile--btn {
  cursor: pointer;
  transition: border-color var(--dur-base) ease, box-shadow var(--dur-base) ease;
}
.tile--btn:hover {
  border-color: var(--accent-focus);
  box-shadow: 0 0 0 3px var(--accent-ring);
}
.tile--btn.is-warn .tile__value {
  color: var(--warning);
}
.tile__value {
  font: var(--fw-semibold) var(--fs-2xl) var(--font-mono);
  letter-spacing: -0.5px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.tile__value small {
  margin-left: 3px;
  font-size: var(--fs-md);
  font-weight: var(--fw-medium);
  color: var(--text-4);
  letter-spacing: 0;
}
.tile__hint {
  display: flex;
  align-items: center;
  gap: 4px;
  font: var(--fw-regular) var(--fs-xs) var(--font-sans);
  color: var(--text-4);
  white-space: nowrap;
}
.tile__hint.mono {
  font-family: var(--font-mono);
}
.tile__hint :deep(svg + *),
.tile__hint :deep(svg) {
  flex: none;
}
.tile__hint .accent,
.tile__hint .info {
  margin-left: 4px;
}
.accent { color: var(--accent); }
.info { color: var(--info); }

.tb {
  display: flex;
  align-items: center;
  gap: 9px;
  margin-bottom: 12px;
}
.tb__search {
  flex: 1;
  min-width: 200px;
  max-width: 420px;
}
.tb__filters {
  display: flex;
  align-items: center;
  gap: 9px;
}
.tb__end {
  margin-left: auto;
  display: flex;
  align-items: center;
  gap: 9px;
}
.state {
  display: flex;
  justify-content: center;
  padding: 60px 0;
  color: var(--text-4);
}
.state--card {
  flex-direction: column;
  align-items: center;
  gap: 4px;
  padding: 20px 0 36px;
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--r-lg);
}

@media (max-width: 1180px) {
  .sum {
    grid-template-columns: repeat(2, 1fr);
  }
  .tb {
    flex-wrap: wrap;
  }
  .tb__search {
    max-width: none;
    flex-basis: 100%;
  }
}
@media (max-width: 767px) {
  .sum {
    gap: 10px;
    margin-bottom: 12px;
  }
  .tile {
    padding: 11px 13px 12px;
  }
  .tile__value {
    font-size: var(--fs-xl);
  }
  .tile__hint {
    display: none;
  }
  .tb {
    gap: 8px;
  }
  .tb__search :deep(input) {
    font-size: 16px;
  }
  .tb__filters {
    flex: 1;
    min-width: 0;
    overflow-x: auto;
    margin: 0 -14px;
    padding: 0 14px;
    order: 2;
  }
  .tb__end {
    order: 1;
    margin-left: 0;
    flex-basis: 100%;
  }
  .tb__end :deep(.x-btn--primary) {
    flex: 1;
    height: var(--h-md);
  }
  .tb__filters :deep(.x-select__trigger) {
    height: var(--h-md);
  }
}
</style>
