<script setup lang="ts">
import { computed, defineAsyncComponent, onBeforeUnmount, onMounted, reactive, ref, watch } from 'vue'
import { api } from '../lib/api'
import { config } from '../lib/config'
import { confirm } from '../lib/confirm'
import { cpuSpeed, duration, percent, size, sizeParts, version } from '../lib/format'
import { t, tn } from '../lib/i18n'
import { isMobile } from '../lib/media'
import { normState, stateColor, stateLabel, type ServerStatus } from '../lib/status'
import { ws } from '../lib/ws'
import AppShell from '../components/AppShell.vue'
import Card from '../components/Card.vue'
import Chip from '../components/Chip.vue'
import ConfirmDialog from '../components/ConfirmDialog.vue'
import CpuChart from '../components/CpuChart.vue'
import Dropdown from '../components/Dropdown.vue'
import Icon from '../components/Icon.vue'
import LinkRow from '../components/LinkRow.vue'
import Metric from '../components/Metric.vue'
import ResourceRow from '../components/ResourceRow.vue'
import Spinner from '../components/Spinner.vue'

const XrayUpdatesModal = defineAsyncComponent(() => import('../modals/XrayUpdatesModal.vue'))
const PanelUpdateModal = defineAsyncComponent(() => import('../modals/PanelUpdateModal.vue'))
const PanelLogModal = defineAsyncComponent(() => import('../modals/PanelLogModal.vue'))
const BackupModal = defineAsyncComponent(() => import('../modals/BackupModal.vue'))
const ConfigModal = defineAsyncComponent(() => import('../modals/ConfigModal.vue'))
const XrayLogModal = defineAsyncComponent(() => import('../modals/XrayLogModal.vue'))

type ModalName = 'updates' | 'panelUpdate' | 'panelLog' | 'backup' | 'config' | 'xrayLog'
const modal = ref<ModalName | null>(null)

// ───── server status ─────
const status = ref<ServerStatus | null>(null)
const onlineClients = ref(0)
const xrayState = computed(() => normState(status.value?.xray.state))
const xrayError = computed(() => (xrayState.value === 'error' ? status.value?.xray.errorMsg || '' : ''))

async function fetchStatus() {
  const msg = await api.get<ServerStatus>('panel/api/server/status', { notify: false })
  if (msg.success && msg.obj) applyStatus(msg.obj)
}

function applyStatus(s: ServerStatus) {
  status.value = s
  if (range.value === 2) {
    cpuSeries.value = [...cpuSeries.value, Math.max(0, Math.min(100, s.cpu))].slice(-60)
  }
}

async function fetchOnlines() {
  const msg = await api.post<string[] | null>('panel/api/inbounds/onlines', undefined, { notify: false })
  if (msg.success) onlineClients.value = msg.obj?.length ?? 0
}

const res = computed(() => {
  const s = status.value
  if (!s) return null
  const pair = (c: number, tot: number) => `${size(c)} / ${size(tot)}`
  const cores = tn('spx.cores', s.cpuCores)
  return {
    cpu: { detail: `${cores} · ${cpuSpeed(s.cpuSpeedMhz)}`, pct: Math.min(100, Math.max(0, s.cpu)) },
    mem: { detail: pair(s.mem.current, s.mem.total), pct: percent(s.mem.current, s.mem.total) },
    swap: { detail: pair(s.swap.current, s.swap.total), pct: percent(s.swap.current, s.swap.total) },
    disk: { detail: pair(s.disk.current, s.disk.total), pct: percent(s.disk.current, s.disk.total) },
  }
})

// ───── CPU history ─────
// The server keeps 60 points per bucket size: 2 s buckets cover 2 minutes,
// 300 s buckets cover 5 hours.
const rangeOptions = [
  { value: 2, label: t('spx.range2m') },
  { value: 30, label: t('spx.range30m') },
  { value: 60, label: t('spx.range1h') },
  { value: 120, label: t('spx.range2h') },
  { value: 180, label: t('spx.range3h') },
  { value: 300, label: t('spx.range5h') },
]
const range = ref(2)
const cpuSeries = ref<number[]>([])
const cpuPeak = computed(() => (cpuSeries.value.length ? Math.max(...cpuSeries.value).toFixed(1) : '0.0'))
const rangeDisplay = computed(() =>
  t('spx.rangeFor', { range: rangeOptions.find((o) => o.value === range.value)?.label ?? '' }),
)

async function fetchCpuHistory() {
  const bucket = range.value
  const msg = await api.get<{ t: number; cpu: number }[]>(`panel/api/server/cpuHistory/${bucket}`, { notify: false })
  if (msg.success && Array.isArray(msg.obj) && bucket === range.value) {
    cpuSeries.value = msg.obj.map((p) => Math.max(0, Math.min(100, p.cpu)))
  }
}
watch(range, fetchCpuHistory)

// ───── Xray actions ─────
const busy = reactive({ restart: false, stop: false })

async function restartXray() {
  if (busy.restart) return
  busy.restart = true
  await api.post('panel/api/server/restartXrayService')
  busy.restart = false
  fetchStatus()
}

async function stopXray() {
  await confirm({
    title: t('spx.stopXrayTitle'),
    text: t('spx.stopXrayText'),
    cta: t('pages.index.stopXray'),
    busy: t('spx.stopping'),
    tone: 'danger',
    action: async () => {
      await api.post('panel/api/server/stopXrayService')
      await fetchStatus()
    },
  })
}

// ───── panel version ─────
const updateInfo = ref({ currentVersion: config.version, latestVersion: '', updateAvailable: false })
async function fetchUpdateInfo() {
  const msg = await api.get<typeof updateInfo.value>('panel/api/server/getPanelUpdateInfo', { notify: false })
  if (msg.success && msg.obj) updateInfo.value = msg.obj
}

// ───── misc ─────
const ipHidden = ref(true)
const xrayLogEnabled = ref(false)
const httpWarning = ref(location.protocol !== 'https:' && !sessionStorageGet('spx-http-warning-dismissed'))
function dismissHttpWarning() {
  httpWarning.value = false
  try {
    sessionStorage.setItem('spx-http-warning-dismissed', '1')
  } catch {
    /* ignore */
  }
}
function sessionStorageGet(key: string) {
  try {
    return sessionStorage.getItem(key)
  } catch {
    return null
  }
}

const net = computed(() => {
  const s = status.value
  return {
    sent: sizeParts(s?.netTraffic.sent ?? 0),
    recv: sizeParts(s?.netTraffic.recv ?? 0),
    up: `${size(s?.netIO.up ?? 0)}/s`,
    down: `${size(s?.netIO.down ?? 0)}/s`,
  }
})
const loads = computed(() => (status.value?.loads ?? [0, 0, 0]).map((l) => l.toFixed(2)))

// ───── lifecycle ─────
let poll: number | undefined
let slow: number | undefined
const offs: (() => void)[] = []

onMounted(async () => {
  fetchStatus()
  fetchCpuHistory()
  fetchOnlines()
  fetchUpdateInfo()
  api.post<{ ipLimitEnable: boolean }>('panel/setting/defaultSettings', undefined, { notify: false }).then((m) => {
    if (m.success) xrayLogEnabled.value = !!m.obj?.ipLimitEnable
  })

  offs.push(
    ws.on('status', (p: ServerStatus) => applyStatus(p)),
    ws.on('traffic', (p: { onlineClients?: string[] }) => {
      if (Array.isArray(p?.onlineClients)) onlineClients.value = p.onlineClients.length
    }),
    ws.on('xray_state', (p: { state: string; errorMsg?: string }) => {
      if (status.value) status.value.xray = { ...status.value.xray, state: p.state, errorMsg: p.errorMsg ?? '' }
    }),
  )
  ws.connect()

  // Poll while the socket is down; refresh the slow-moving bits regardless.
  poll = window.setInterval(() => {
    if (!ws.connected) fetchStatus()
  }, 2000)
  slow = window.setInterval(() => {
    if (range.value !== 2) fetchCpuHistory()
    if (!ws.connected) fetchOnlines()
  }, 30_000)
})

onBeforeUnmount(() => {
  clearInterval(poll)
  clearInterval(slow)
  offs.forEach((off) => off())
})
</script>

<template>
  <AppShell active="dashboard" :title="t('pages.index.title')" :xray-state="xrayState">
    <template #head>
      <div v-if="httpWarning" class="http-warn">
        <Icon name="alert" />
        <div class="hw-text">
          <div class="hw-title">{{ t('spx.httpWarningTitle') }}</div>
          <div class="hw-body">{{ t('spx.httpWarningText') }}</div>
        </div>
        <button class="hw-close" :title="t('close')" @click="dismissHttpWarning">✕</button>
      </div>
    </template>

    <div v-if="isMobile && httpWarning" class="http-warn mobile">
      <Icon name="alert" />
      <div class="hw-text">
        <div class="hw-title">{{ t('spx.httpWarningTitle') }}</div>
        <div class="hw-body">{{ t('spx.httpWarningText') }}</div>
      </div>
      <button class="hw-close" :title="t('close')" @click="dismissHttpWarning">✕</button>
    </div>

    <div v-if="!status" class="loading"><Spinner :size="22" /></div>

    <div v-else class="grid" :class="{ mobile: isMobile }">
      <div class="col">
      <!-- Xray -->
      <Card label="Xray" class="c-xray">
        <template #extra>
          <Chip v-if="status.xray.version && status.xray.version !== 'Unknown'">{{ version(status.xray.version) }}</Chip>
        </template>
        <div class="state-row">
          <span class="dot state-dot" :style="{ color: stateColor(xrayState) }" />
          <span class="state" :style="{ color: stateColor(xrayState) }">{{ stateLabel(xrayState) }}</span>
          <div v-if="isMobile" class="actions">
            <button class="ctl icon" :title="t('pages.index.restartXray')" @click="restartXray">
              <Spinner v-if="busy.restart" :size="15" /><Icon v-else name="restart" />
            </button>
            <button class="ctl icon danger" :title="t('pages.index.stopXray')" @click="stopXray"><Icon name="stop" /></button>
          </div>
        </div>
        <div v-if="xrayError" class="err-box mono">{{ xrayError }}</div>
        <div class="metrics">
          <Metric :label="isMobile ? t('spx.online') : t('spx.clientsOnline')" :value="onlineClients" />
          <Metric :label="isMobile ? t('spx.runs') : t('spx.uptime')" :value="duration(status.appStats.uptime)" />
          <Metric :label="t('spx.memory')" :value="size(status.appStats.mem)" />
          <div v-if="!isMobile" class="actions">
            <button class="ctl icon" :title="t('pages.index.restartXray')" @click="restartXray">
              <Spinner v-if="busy.restart" :size="15" /><Icon v-else name="restart" />
            </button>
            <button class="ctl icon danger" :title="t('pages.index.stopXray')" @click="stopXray"><Icon name="stop" /></button>
          </div>
        </div>
        <div class="links">
          <LinkRow v-if="xrayLogEnabled" icon="log" :label="t('spx.xrayLog')" @click="modal = 'xrayLog'" />
          <LinkRow icon="version" :label="t('pages.index.xraySwitch')" @click="modal = 'updates'" />
        </div>
      </Card>

      <!-- Management -->
      <Card :label="t('menu.link')" class="c-mgmt">
        <div class="mgmt-links">
          <LinkRow icon="log" :label="t('spx.panelLog')" @click="modal = 'panelLog'" />
          <LinkRow icon="code" :label="t('pages.index.config')" @click="modal = 'config'" />
          <LinkRow icon="database" :label="t('pages.index.backup')" @click="modal = 'backup'" />
        </div>
        <div class="version-row">
          <span class="v-label">{{ t('spx.panelVersion') }}</span>
          <button class="v-value" @click="modal = 'panelUpdate'">{{ version(updateInfo.currentVersion || config.version) }}</button>
          <button v-if="updateInfo.updateAvailable" class="update-cta" @click="modal = 'panelUpdate'">
            <span class="dot" />{{ t('spx.updatePanelTo') }} <span class="mono">{{ version(updateInfo.latestVersion) }}</span>
          </button>
        </div>
      </Card>

      <!-- Server IPs -->
      <Card :label="t('pages.index.ipAddresses')" class="c-ip">
        <template #extra>
          <button class="ctl sm icon" :title="t('pages.index.toggleIpVisibility')" @click="ipHidden = !ipHidden">
            <Icon :name="ipHidden ? 'eye' : 'eyeOff'" :size="15" />
          </button>
        </template>
        <div class="ip-grid">
          <div class="well">
            <div class="caps-sm keep-case">IPv4</div>
            <div class="ip" :class="{ blur: ipHidden }">{{ status.publicIP.ipv4 || 'N/A' }}</div>
          </div>
          <div class="well">
            <div class="caps-sm keep-case">IPv6</div>
            <div class="ip" :class="{ blur: ipHidden }">{{ status.publicIP.ipv6 || 'N/A' }}</div>
          </div>
        </div>
      </Card>

      </div>
      <div class="col">
      <!-- Resources + CPU history -->
      <Card :label="t('spx.serverResources')" class="c-res">
        <div class="res-list">
          <ResourceRow :label="t('pages.index.cpu')" :detail="res!.cpu.detail" :percent="res!.cpu.pct" />
          <ResourceRow :label="t('pages.index.memory')" :detail="res!.mem.detail" :percent="res!.mem.pct" />
          <ResourceRow :label="t('spx.swap')" :detail="res!.swap.detail" :percent="res!.swap.pct" />
          <ResourceRow :label="t('pages.index.storage')" :detail="res!.disk.detail" :percent="res!.disk.pct" />
        </div>
        <div class="cpu">
          <div class="cpu-head">
            <span class="caps-sm">{{ t('spx.cpuLoad') }}</span>
            <Dropdown v-model="range" :options="rangeOptions" size="sm" :display="rangeDisplay" :menu-width="124" :title="t('spx.cpuLoad')" />
            <span class="peak">{{ t('spx.peak', { value: cpuPeak }) }}</span>
          </div>
          <CpuChart :values="cpuSeries" />
        </div>
      </Card>

      <!-- Network -->
      <Card :label="t('spx.network')" class="c-net">
        <div class="net-totals">
          <div class="well">
            <div class="net-cap"><Icon name="up" :size="13" class="up" /><span class="caps-sm">{{ t('spx.totalSent') }}</span></div>
            <div class="net-big up">{{ net.sent[0] }} <small>{{ net.sent[1] }}</small></div>
          </div>
          <div class="well">
            <div class="net-cap"><Icon name="down" :size="13" class="down" /><span class="caps-sm">{{ t('spx.totalReceived') }}</span></div>
            <div class="net-big down">{{ net.recv[0] }} <small>{{ net.recv[1] }}</small></div>
          </div>
        </div>
        <div class="net-rates">
          <div class="pair">
            <Metric :label="t('pages.index.upload')" :value="net.up" size="sm" color="var(--info)" />
            <Metric :label="t('pages.index.download')" :value="net.down" size="sm" color="var(--accent)" />
          </div>
          <div class="pair">
            <Metric label="TCP" :value="status.tcpCount" size="sm" color="var(--text-2)" />
            <Metric label="UDP" :value="status.udpCount" size="sm" color="var(--text-2)" />
          </div>
        </div>
      </Card>

      <!-- System -->
      <Card :label="t('spx.system')" class="c-sys">
        <div class="sys">
          <Metric :label="t('spx.osUptime')" :value="duration(status.uptime)" class="sys-uptime" />
          <div class="sys-load">
            <div class="caps-sm">{{ t('spx.loadAverage') }}</div>
            <div class="load mono">
              {{ loads[0] }} <span class="sep">|</span> {{ loads[1] }} <span class="sep">|</span> {{ loads[2] }}
            </div>
            <div class="load-hint">{{ t('pages.index.systemLoadDesc') }}</div>
          </div>
        </div>
      </Card>
      </div>
    </div>

    <XrayUpdatesModal v-if="modal === 'updates'" :current="status?.xray.version ?? ''" @close="modal = null" />
    <PanelUpdateModal v-if="modal === 'panelUpdate'" :info="updateInfo" @close="modal = null" />
    <PanelLogModal v-if="modal === 'panelLog'" @close="modal = null" />
    <BackupModal v-if="modal === 'backup'" @close="modal = null" />
    <ConfigModal v-if="modal === 'config'" @close="modal = null" />
    <XrayLogModal v-if="modal === 'xrayLog'" @close="modal = null" />
    <ConfirmDialog />
  </AppShell>
</template>

<style scoped>
.loading {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--accent);
}

/* HTTP warning: in the header on desktop, aligned with the right column. */
.http-warn {
  margin-left: auto;
  width: calc(50% - 7px);
  min-height: 45px;
  flex: none;
  display: flex;
  gap: 10px;
  align-items: center;
  background: rgba(var(--warn-rgb), 0.1);
  border: 1px solid rgba(var(--warn-rgb), 0.3);
  border-radius: 11px;
  padding: 6px 11px 6px 12px;
  color: var(--warn);
}
:root[data-theme='light'] .http-warn {
  color: var(--warn-title);
}
.hw-text {
  flex: 1;
  min-width: 0;
}
.hw-title {
  font: 600 12px var(--font-sans);
  color: var(--warn-title);
  margin-bottom: 1px;
}
.hw-body {
  font: 400 11.5px var(--font-sans);
  color: var(--warn-body);
}
.hw-close {
  flex: none;
  align-self: flex-start;
  margin-top: 6px;
  color: var(--warn-body);
  font: 400 12px var(--font-mono);
  padding: 0 2px;
}
.http-warn.mobile {
  width: auto;
  margin: 0 0 12px;
  align-items: flex-start;
  padding: 11px 12px;
}
.http-warn.mobile svg {
  margin-top: 1px;
}
.http-warn.mobile .hw-title {
  font-size: 13px;
  margin-bottom: 3px;
}
.http-warn.mobile .hw-body {
  font-size: 12px;
  line-height: 1.5;
}
.http-warn.mobile .hw-close {
  margin-top: 0;
}

/* Two equal columns on desktop; on phones one column in a fixed order. */
.grid {
  flex: 1 0 auto;
  display: grid;
  grid-template-columns: minmax(0, 1fr) minmax(0, 1fr);
  gap: 12px;
}
.col {
  display: flex;
  flex-direction: column;
  gap: 12px;
  min-width: 0;
}
.c-mgmt,
.c-res {
  flex: 1;
}
.grid.mobile {
  display: flex;
  flex-direction: column;
}
.grid.mobile .col {
  display: contents;
}
.grid.mobile .c-xray {
  order: 1;
}
.grid.mobile .c-res {
  order: 2;
}
.grid.mobile .c-net {
  order: 3;
}
.grid.mobile .c-ip {
  order: 4;
}
.grid.mobile .c-mgmt {
  order: 5;
}
.grid.mobile .c-sys {
  order: 6;
}

.state-row {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 12px;
}
.state-dot {
  width: 9px;
  height: 9px;
}
.state {
  font: 600 18px var(--font-sans);
  letter-spacing: -0.4px;
}
.err-box {
  background: rgba(var(--danger-rgb), 0.1);
  border: 1px solid rgba(var(--danger-rgb), 0.3);
  border-radius: 10px;
  padding: 10px 12px;
  margin-bottom: 16px;
  font: 400 11.5px/1.5 var(--font-mono);
  color: var(--danger-text);
  max-height: 76px;
  overflow: auto;
  word-break: break-word;
}
.metrics {
  display: flex;
  align-items: flex-end;
  gap: 12px;
}
.metrics > .metric {
  flex: 1;
}
.actions {
  flex: none;
  display: flex;
  gap: 8px;
  margin-left: auto;
}
.links {
  display: flex;
  flex-direction: column;
  gap: 2px;
  margin-top: 12px;
  padding-top: 11px;
  border-top: 1px solid var(--line);
}
.mobile .links {
  margin-top: 14px;
  padding-top: 4px;
  gap: 0;
}
.mobile .actions .ctl {
  width: 36px;
  height: 36px;
}

.mgmt-links {
  display: flex;
  flex-direction: column;
  gap: 2px;
  margin-bottom: 13px;
}
.version-row {
  margin-top: auto;
  padding-top: 13px;
  border-top: 1px solid var(--line);
  display: flex;
  align-items: center;
  gap: 8px;
  min-height: 28px;
  flex-wrap: wrap;
}
.v-label {
  font: 400 11.5px var(--font-sans);
  color: var(--muted);
  white-space: nowrap;
}
.v-value {
  font: 500 12px var(--font-mono);
  color: var(--text-2);
}
.v-value:hover {
  color: var(--accent);
}
.update-cta {
  margin-left: auto;
  display: flex;
  align-items: center;
  gap: 7px;
  height: 28px;
  padding: 0 10px;
  border-radius: 9px;
  border: 1px solid rgba(var(--accent-rgb), 0.34);
  background: var(--accent-soft);
  color: var(--accent);
  font: 500 11.5px var(--font-sans);
  white-space: nowrap;
  transition:
    box-shadow 0.22s ease,
    border-color 0.22s ease;
}
.update-cta:hover {
  border-color: rgba(var(--accent-rgb), 0.6);
  box-shadow: var(--glow);
}
.update-cta .dot {
  width: 6px;
  height: 6px;
}

.ip-grid,
.net-totals {
  display: grid;
  grid-template-columns: minmax(0, 1fr) minmax(0, 1fr);
  gap: 12px;
}
.well {
  background: var(--fill);
  border: 1px solid var(--border);
  border-radius: 12px;
  padding: 10px 14px;
  min-width: 0;
}
.well .caps-sm {
  display: block;
  margin-bottom: 6px;
}
.keep-case {
  text-transform: none;
}
.ip {
  font: 500 14px var(--font-mono);
  color: var(--text);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  transition: filter 0.2s ease;
}
.ip.blur {
  filter: blur(5px);
  user-select: none;
}

.res-list {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 16px;
}
.cpu {
  flex: none;
  margin-top: 14px;
  padding-top: 13px;
  border-top: 1px solid var(--line);
}
.cpu-head {
  display: flex;
  align-items: center;
  gap: 9px;
  margin-bottom: 8px;
}
.peak {
  margin-left: auto;
  font: 500 10px var(--font-mono);
  color: var(--text-3b);
  white-space: nowrap;
}

.net-totals {
  margin-bottom: 18px;
}
.net-totals .well {
  padding: 12px 14px;
}
.net-cap {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 7px;
}
.net-cap .caps-sm {
  margin: 0;
  overflow: hidden;
  text-overflow: ellipsis;
}
.up {
  color: var(--info);
}
.down {
  color: var(--accent);
}
.net-big {
  font: 600 22px var(--font-mono);
  letter-spacing: -0.5px;
  white-space: nowrap;
}
.net-big small {
  font-size: 13px;
  font-weight: 500;
  color: var(--text-3b);
  letter-spacing: 0;
}
.net-rates {
  display: grid;
  grid-template-columns: minmax(0, 1fr) minmax(0, 1fr);
  gap: 12px;
}
.pair {
  display: grid;
  grid-template-columns: minmax(0, 1fr) minmax(0, 1fr);
  gap: 12px;
  padding: 0 14px;
}
.mobile .net-rates {
  grid-template-columns: 1fr;
  gap: 14px;
}
.mobile .net-big {
  font-size: 19px;
}

.sys {
  display: flex;
  align-items: flex-start;
  gap: 24px;
}
.sys-uptime {
  flex: none;
}
.sys-load {
  flex: 1;
  min-width: 0;
}
.sys-load .caps-sm {
  display: block;
  margin-bottom: 5px;
}
.load {
  font: 500 16px var(--font-mono);
  color: var(--text);
}
.sep {
  color: var(--separator);
}
.load-hint {
  font: 400 11px var(--font-sans);
  color: var(--muted);
  margin-top: 5px;
}
</style>
