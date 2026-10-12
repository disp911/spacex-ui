<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { api, ApiError } from '../api'
import { config } from '../config'
import { bytesParts, formatDuration, formatNumber, formatSpeed } from '../i18n/format'
import { toast, XButton, XIcon, XNotice } from '../kit'
import AppPage from '../layout/AppPage.vue'
import { useServerStatus } from '../stores/server'
import DashCard from './dashboard/DashCard.vue'
import DashLink from './dashboard/DashLink.vue'
import DashStat from './dashboard/DashStat.vue'
import PanelLogModal from './dashboard/PanelLogModal.vue'
import ResourcesCard from './dashboard/ResourcesCard.vue'
import ServiceCard from './dashboard/ServiceCard.vue'

const { t } = useI18n()
const server = useServerStatus()
server.watch()
const status = computed(() => server.status)
const host = location.hostname

// --- Services -------------------------------------------------------------

type Service = 'xray' | 'telegram'
const busy = ref<Record<Service, 'restart' | 'stop' | null>>({ xray: null, telegram: null })

const endpoints = {
  xray: { restart: 'restartXrayService', stop: 'stopXrayService' },
  telegram: { restart: 'restartTelemtService', stop: 'stopTelemtService' },
} as const
const done = {
  xray: { restart: 'dashboard.xrayRestarted', stop: 'dashboard.xrayStopped' },
  telegram: { restart: 'dashboard.telegramRestarted', stop: 'dashboard.telegramStopped' },
} as const

async function act(service: Service, action: 'restart' | 'stop') {
  busy.value[service] = action
  try {
    await api.post(`panel/api/server/${endpoints[service][action]}`)
    toast('ok', t(done[service][action]))
  } catch (e) {
    toast('err', t('dashboard.actionFailed', { error: e instanceof ApiError ? e.message : String(e) }))
  } finally {
    await server.refresh().catch(() => {})
    busy.value[service] = null
  }
}

// --- Notices and small state ----------------------------------------------

const localHosts = ['localhost', '127.0.0.1', '[::1]']
const insecure = location.protocol === 'http:' && !localHosts.includes(location.hostname)
const httpDismissed = ref(read('dash.httpDismissed', sessionStorage) === '1')
function dismissHttp() {
  httpDismissed.value = true
  write('dash.httpDismissed', '1', sessionStorage)
}

const ipHidden = ref(read('dash.ipShown', localStorage) !== '1')
function toggleIp() {
  ipHidden.value = !ipHidden.value
  write('dash.ipShown', ipHidden.value ? '0' : '1', localStorage)
}

function read(key: string, storage: Storage) {
  try {
    return storage.getItem(key)
  } catch {
    return null
  }
}
function write(key: string, value: string, storage: Storage) {
  try {
    storage.setItem(key, value)
  } catch {
    // only this visit remembers it
  }
}

const ip = (v?: string) => (!v || v === 'N/A' ? '—' : v)

// --- Windows ----------------------------------------------------------------

const panelLogOpen = ref(false)

// --- Panel version ----------------------------------------------------------

const update = ref<{ latestVersion: string; updateAvailable: boolean } | null>(null)
onMounted(async () => {
  try {
    update.value = await api.get('panel/api/server/getPanelUpdateInfo')
  } catch {
    // no answer from GitHub: nothing to offer
  }
})

const sent = computed(() => bytesParts(status.value?.netTraffic.sent ?? 0))
const received = computed(() => bytesParts(status.value?.netTraffic.recv ?? 0))
const dash = '—'
</script>

<template>
  <AppPage :title="t('nav.dashboard')">
    <template #subtitle>
      <i18n-t keypath="common.server" tag="span">
        <template #host><span class="mono">{{ host }}</span></template>
      </i18n-t>
    </template>
    <template v-if="insecure && !httpDismissed" #aside>
      <XNotice tone="warn" :title="t('dashboard.httpTitle')" class="http">
        {{ t('dashboard.httpText') }}
        <template #actions>
          <XButton variant="ghost" size="sm" icon="close" :label="t('kit.close')" class="http__close" @click="dismissHttp" />
        </template>
      </XNotice>
    </template>

    <div class="dash">
      <div class="dash__col">
        <ServiceCard
          title="Xray"
          version-prefix="v"
          class="area-xray"
          :service="status?.xray"
          :error-title="t('dashboard.xrayError')"
          :busy="busy.xray"
          @restart="act('xray', 'restart')"
          @stop="act('xray', 'stop')"
          @log="panelLogOpen = true"
        >
          <DashLink icon="doc" :label="t('dashboard.xrayLog')" />
          <DashLink icon="download" :label="t('dashboard.xrayVersion')" />
        </ServiceCard>

        <ServiceCard
          v-if="status?.telemt"
          :title="t('dashboard.telegram')"
          class="area-telegram"
          :service="status.telemt"
          :error-title="t('dashboard.telegramError')"
          :busy="busy.telegram"
          @restart="act('telegram', 'restart')"
          @stop="act('telegram', 'stop')"
          @log="panelLogOpen = true"
        />

        <DashCard :title="t('dashboard.manage')" class="grow area-manage">
          <div class="links">
            <DashLink icon="doc" :label="t('dashboard.panelLog')" @click="panelLogOpen = true" />
            <DashLink icon="code" :label="t('dashboard.config')" />
            <DashLink icon="database" :label="t('dashboard.backup')" />
          </div>
          <div class="version">
            <span class="version__label">{{ t('dashboard.panelVersion') }}</span>
            <span class="mono version__value">v{{ config.version }}</span>
            <button v-if="update?.updateAvailable" type="button" class="version__update">
              <span class="version__dot" />
              {{ t('dashboard.update', { version: update.latestVersion }) }}
            </button>
          </div>
        </DashCard>

        <DashCard :title="t('dashboard.ip')" class="area-ip">
          <template #head>
            <XButton
              variant="ghost"
              size="sm"
              :icon="ipHidden ? 'eye' : 'eyeOff'"
              :label="ipHidden ? t('dashboard.showIp') : t('dashboard.hideIp')"
              :aria-pressed="!ipHidden"
              class="ip-toggle"
              @click="toggleIp"
            />
          </template>
          <div class="tiles">
            <div class="tile">
              <div class="tile__label">IPv4</div>
              <div :class="['tile__ip', { 'is-hidden': ipHidden }]">{{ ip(status?.publicIP.ipv4) }}</div>
            </div>
            <div class="tile">
              <div class="tile__label">IPv6</div>
              <div :class="['tile__ip', { 'is-hidden': ipHidden }]">{{ ip(status?.publicIP.ipv6) }}</div>
            </div>
          </div>
        </DashCard>
      </div>

      <div class="dash__col">
        <ResourcesCard :status="status" class="grow area-resources" />

        <DashCard :title="t('dashboard.network')" class="area-network">
          <div class="tiles">
            <div class="tile tile--total">
              <div class="tile__label">
                <XIcon name="arrowUp" :size="13" class="up" />
                <span class="long">{{ t('dashboard.sent') }}</span><span class="short">{{ t('dashboard.sentShort') }}</span>
              </div>
              <div class="tile__total up">{{ status ? sent.n : dash }} <small v-if="status">{{ sent.unit }}</small></div>
            </div>
            <div class="tile tile--total">
              <div class="tile__label">
                <XIcon name="arrowDown" :size="13" class="down" />
                <span class="long">{{ t('dashboard.received') }}</span><span class="short">{{ t('dashboard.receivedShort') }}</span>
              </div>
              <div class="tile__total down">{{ status ? received.n : dash }} <small v-if="status">{{ received.unit }}</small></div>
            </div>
          </div>
          <div class="net-figures">
            <div class="net-figures__pair">
              <DashStat :label="t('dashboard.upload')" :value="status ? formatSpeed(status.netIO.up) : dash" class="small up" />
              <DashStat :label="t('dashboard.download')" :value="status ? formatSpeed(status.netIO.down) : dash" class="small down" />
            </div>
            <div class="net-figures__pair">
              <DashStat label="TCP" :value="status ? formatNumber(status.tcpCount) : dash" class="small" />
              <DashStat label="UDP" :value="status ? formatNumber(status.udpCount) : dash" class="small" />
            </div>
          </div>
        </DashCard>

        <DashCard :title="t('dashboard.system')" class="area-system">
          <div class="system">
            <DashStat :label="t('dashboard.osUptime')" :value="status ? formatDuration(status.uptime) : dash" class="system__uptime" />
            <div class="system__load">
              <DashStat :label="t('dashboard.load')" :short="t('dashboard.loadShort')" value="">
                <template v-if="status?.loads?.length">
                  <template v-for="(l, i) in status.loads" :key="i">
                    <span v-if="i" class="system__sep"> | </span>{{ formatNumber(l, 2) }}
                  </template>
                </template>
                <template v-else>{{ dash }}</template>
              </DashStat>
              <p class="system__hint">{{ t('dashboard.loadHint') }}</p>
            </div>
          </div>
        </DashCard>
      </div>
    </div>

    <PanelLogModal v-model:open="panelLogOpen" />
  </AppPage>
</template>

<style scoped>
.mono {
  color: var(--fg2);
  font-family: var(--font-mono);
}

.http {
  align-items: center;
  padding: 6px 6px 6px 12px;
  border-radius: var(--r-12);
  font-size: var(--fs-12);
  line-height: var(--lh-tight);
}

.http :deep(.x-notice__icon) {
  margin-top: 0;
}

.http :deep(.x-notice__title) {
  margin-bottom: 1px;
  font-size: var(--fs-12);
}

.dash {
  display: grid;
  flex: 1 0 auto;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
  align-items: stretch;
}

.dash__col {
  display: flex;
  flex-direction: column;
  gap: 12px;
  min-width: 0;
}

.grow {
  flex: 1 0 auto;
}

/* Management */
.links {
  display: flex;
  flex-direction: column;
  gap: 2px;
  margin-bottom: 12px;
}

.version {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 8px;
  min-height: 28px;
  margin-top: auto;
  padding-top: 12px;
  border-top: 1px solid var(--div);
}

.version__label {
  color: var(--dim);
  font-size: var(--fs-12);
}

.version__value {
  font-size: var(--fs-12);
}

.version__update {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  height: var(--ctl-sm-h);
  margin-left: auto;
  padding: 0 10px;
  border: 1px solid var(--acc-line);
  border-radius: var(--r-8);
  background: var(--acc-soft);
  color: var(--acc);
  font: 500 var(--fs-12) / 1 var(--font-sans);
  white-space: nowrap;
  outline: none;
  transition: box-shadow var(--t-fast), border-color var(--t-fast);
}

.version__update:hover,
.version__update:focus-visible {
  border-color: var(--acc);
  box-shadow: 0 0 0 3px var(--acc-ring);
}

.version__dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--acc);
  animation: dash-pulse 2.4s infinite;
}

@keyframes dash-pulse {
  50% {
    opacity: .35;
  }
}

/* Tiles: IP addresses and traffic totals */
.ip-toggle {
  margin: -4px 0;
}

.tiles {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

.tile {
  min-width: 0;
  padding: 10px 14px;
  border: 1px solid var(--field-line);
  border-radius: var(--r-12);
  background: var(--field);
}

.tile--total {
  padding: 12px 14px;
}

.tile__label {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 8px;
  color: var(--dim);
  font: 500 var(--fs-10) / var(--lh-tight) var(--font-mono);
  letter-spacing: .07em;
  text-transform: uppercase;
  white-space: nowrap;
}

.tile__ip {
  overflow: hidden;
  color: var(--fg);
  font: 500 var(--fs-14) / var(--lh-tight) var(--font-mono);
  white-space: nowrap;
  text-overflow: ellipsis;
  transition: filter var(--t-med);
}

.tile__ip.is-hidden {
  filter: blur(5px);
  user-select: none;
}

.tile__total {
  font: 600 var(--fs-20) / var(--lh-tight) var(--font-mono);
  letter-spacing: -.5px;
  white-space: nowrap;
}

.tile__total small {
  color: var(--mut);
  font-size: var(--fs-13);
  font-weight: 500;
  letter-spacing: 0;
}

.up {
  color: var(--upload);
}

.down {
  color: var(--acc);
}

.short {
  display: none;
}

/* Network figures: two pairs under the two totals. */
.net-figures,
.net-figures__pair {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

.net-figures {
  margin-top: 16px;
}

.net-figures__pair {
  padding: 0 14px;
}

.small :deep(.stat__value) {
  font-size: var(--fs-14);
}

.small.up :deep(.stat__value) {
  color: var(--upload);
}

.small.down :deep(.stat__value) {
  color: var(--acc);
}

.small:not(.up, .down) :deep(.stat__value) {
  color: var(--fg2);
}

/* System */
.system {
  display: flex;
  align-items: flex-start;
  gap: 24px;
}

.system__uptime {
  flex: none;
}

.system__load {
  flex: 1;
  min-width: 0;
}

.system__sep {
  color: var(--faint);
}

.system__hint {
  margin: 4px 0 0;
  color: var(--dim);
  font-size: var(--fs-11);
}

@media (max-width: 1100px) {
  .dash {
    grid-template-columns: minmax(0, 1fr);
  }

  .grow {
    flex: none;
  }
}

/* Phones: one column in the order of the mockup. */
@media (max-width: 767px) {
  .dash {
    display: flex;
    flex-direction: column;
  }

  .dash__col {
    display: contents;
  }

  .area-xray { order: 1; }
  .area-telegram { order: 2; }
  .area-resources { order: 3; }
  .area-network { order: 4; }
  .area-ip { order: 5; }
  .area-manage { order: 6; }
  .area-system { order: 7; }

  .net-figures {
    grid-template-columns: repeat(4, minmax(0, 1fr));
  }

  .net-figures__pair {
    display: contents;
  }

  .tile__total {
    font-size: var(--fs-18);
  }

  .long {
    display: none;
  }

  .short {
    display: inline;
  }

  .system {
    gap: 20px;
  }

  .system__hint {
    display: none;
  }

  /* On phones the update offer is a full-width button under the version. */
  .version {
    margin-top: 0;
    padding-top: 12px;
  }

  .version__update {
    flex-basis: 100%;
    justify-content: center;
    height: var(--ctl-h);
    margin: 2px 0 0;
    border-radius: var(--ctl-r);
    font-size: var(--fs-13);
  }
}
</style>
