<script setup lang="ts">
import { computed, onBeforeUnmount, reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { api } from '../api'
import {
  toast, XButton, XCard, XChip, XEmpty, XExpandTable, XInlineButton, XKeyValue, XMenu, XNotice, XPager, XPopover, XSpinner, XSummaryStat, XSwitch, XTooltip,
  type ExpandColumn, type MenuItem,
} from '../kit'
import { formatBytes, formatDate, formatRelative } from '../i18n/format'
import AppPage from '../layout/AppPage.vue'
import { useInbounds } from '../stores/inbounds'
import { insecure } from './xray/model'
import ClientCounts from './inbounds/ClientCounts.vue'
import InboundsTools from './inbounds/InboundsTools.vue'
import {
  deepSearch, expiryState, filterEmails, hasClients, pageSizes, protocolChips, trafficState, type ClientCounts as Counts, type DbInbound, type Filter, type State,
} from './inbounds/model'

const { t } = useI18n()
const store = useInbounds()
const stop = store.watch()
onBeforeUnmount(stop)

const thresholds = computed(() => store.defaults ?? { expireDiff: 0, trafficDiff: 0 })

// --- Summary ---------------------------------------------------------------------

const totals = computed(() => {
  let up = 0
  let down = 0
  let allTime = 0
  for (const ib of store.list) {
    up += ib.up
    down += ib.down
    allTime += ib.allTime || ib.up + ib.down
  }
  return { up, down, allTime }
})

/** The clients of every inbound together, as the summary counts them. */
const allCounts = computed<Counts>(() => {
  const sum: Counts = { clients: 0, active: [], deactive: [], depleted: [], expiring: [], online: [], comments: new Map() }
  for (const c of store.counts.values()) {
    sum.clients += c.clients
    sum.active.push(...c.active)
    sum.deactive.push(...c.deactive)
    sum.depleted.push(...c.depleted)
    sum.expiring.push(...c.expiring)
    sum.online.push(...c.online)
    c.comments.forEach((v, k) => sum.comments.set(k, v))
  }
  return sum
})

// --- Search and filter ---------------------------------------------------------------

const mode = ref<'search' | 'filter'>('search')
const search = ref('')
const filter = ref<Filter>('none')

// Typed text settles before the list is searched, as on the old panel.
const query = ref('')
let typing = 0
watch(search, (q) => {
  clearTimeout(typing)
  typing = window.setTimeout(() => (query.value = q.trim()), 500)
})
onBeforeUnmount(() => clearTimeout(typing))

const shown = computed<DbInbound[]>(() => {
  const q = query.value
  if (mode.value === 'search' && q) return store.list.filter((ib) => deepSearch(ib, q))
  if (mode.value === 'filter' && filter.value !== 'none') {
    const f = filter.value
    return store.list.filter((ib) => {
      const c = store.counts.get(ib.id)
      return !!c && filterEmails(c, f).length > 0
    })
  }
  return store.list
})

// --- Pages -------------------------------------------------------------------------

const page = ref(1)
const size = ref(25)
watch(
  () => store.defaults?.pageSize,
  (s) => s && (size.value = s),
  { immediate: true },
)
watch([query, filter, mode], () => (page.value = 1))

const pages = computed(() => Math.max(1, Math.ceil(shown.value.length / size.value)))
watch(pages, (n) => page.value > n && (page.value = n))
const rows = computed(() =>
  shown.value.length > size.value ? shown.value.slice((page.value - 1) * size.value, page.value * size.value) : shown.value,
)
const sizes = computed(() => pageSizes(store.defaults?.pageSize || 25, shown.value.length))

// --- Rows --------------------------------------------------------------------------

const columns = computed<ExpandColumn[]>(() => [
  { key: 'menu', label: t('inbounds.colMenu'), size: '28px', muted: 'none' },
  { key: 'enable', label: t('inbounds.colEnable'), size: '52px', muted: 'none' },
  { key: 'remark', label: t('inbounds.colRemark'), size: 'minmax(0, 1.6fr)', muted: 'soft' },
  { key: 'port', label: t('inbounds.colPort'), size: '56px' },
  { key: 'protocol', label: t('inbounds.colProtocol'), size: 'minmax(0, 1.15fr)' },
  { key: 'clients', label: t('inbounds.clients'), size: 'minmax(0, 1.3fr)' },
  { key: 'traffic', label: t('inbounds.colTraffic'), size: '150px' },
  { key: 'allTime', label: t('inbounds.colAllTime'), size: '86px' },
  { key: 'expiry', label: t('inbounds.colExpiry'), size: '128px' },
])
const phoneColumns = computed<ExpandColumn[]>(() => [
  { key: 'menu', label: t('inbounds.colMenu'), size: '44px', muted: 'none' },
  { key: 'remark', label: t('inbounds.colRemark'), muted: 'none' },
  { key: 'info', label: t('inbounds.colInfo'), size: '76px', align: 'end', muted: 'none' },
])

const toneOf = (s: State) => (s === 'unl' ? 'neutral' : s)
const used = (ib: DbInbound) => ib.up + ib.down
const traffic = (ib: DbInbound) => trafficState(used(ib), ib.total, thresholds.value)
const expiry = (ib: DbInbound) => expiryState(ib.expiryTime, thresholds.value)
const security = (s: string) => (s === 'tls' ? 'TLS' : s === 'reality' ? 'Reality' : s)
const countsOf = (ib: DbInbound) => store.counts.get(ib.id)

/** The menu of an inbound, as the mockup orders it; on phones it ends with a switch. */
function menu(ib: DbInbound, phone = false): MenuItem[] {
  const items: MenuItem[] = [{ key: 'edit', label: t('inbounds.edit'), icon: 'edit' }]
  const settings = JSON.parse(ib.settings || '{}')
  if ((ib.protocol === 'shadowsocks' && !hasClients(ib, settings)) || ib.protocol === 'wireguard') items.push({ key: 'qr', label: t('inbounds.qr'), icon: 'qr' })
  if (hasClients(ib, settings)) {
    items.push(
      { key: 'addClient', label: t('inbounds.addClient'), icon: 'userPlus', divided: true },
      { key: 'addBulk', label: t('inbounds.addBulk'), icon: 'users' },
      { key: 'copyClients', label: t('inbounds.copyClients'), icon: 'copy' },
      { key: 'resetClients', label: t('inbounds.resetClients'), icon: 'reset' },
      { key: 'exportLinks', label: t('inbounds.exportLinks'), icon: 'link' },
      ...(store.defaults?.subEnable ? [{ key: 'exportSubs', label: t('inbounds.exportSubs'), icon: 'link' as const }] : []),
      { key: 'delDepleted', label: t('inbounds.delDepleted'), icon: 'userX', danger: true },
    )
  } else items.push({ key: 'info', label: t('inbounds.info'), icon: 'info' })
  items.push(
    { key: 'exportInbound', label: t('inbounds.exportInbound'), icon: 'upload', divided: true },
    { key: 'resetTraffic', label: t('inbounds.resetTraffic'), icon: 'reset' },
    ...(ib.protocol !== 'mtproto' ? [{ key: 'clone', label: t('inbounds.clone'), icon: 'clone' as const }] : []),
    { key: 'delete', label: t('inbounds.delete'), icon: 'delete', danger: true },
  )
  if (phone) items.push({ key: 'enable', label: t('inbounds.enable'), switch: true, checked: ib.enable, busy: toggling.has(ib.id), divided: true })
  return items
}

function act(key: string, ib: DbInbound) {
  if (key === 'enable') setEnable(ib, !ib.enable)
}

// --- Enabling ----------------------------------------------------------------------

/** Inbounds whose switch is being applied, and those whose change failed a moment ago. */
const toggling = reactive(new Set<number>())
const refused = reactive(new Set<number>())

/** Applied at once and rolled back if the server refuses, as on the old panel. */
async function setEnable(ib: DbInbound, on: boolean) {
  if (toggling.has(ib.id)) return
  const before = ib.enable
  store.patch(ib.id, { enable: on })
  toggling.add(ib.id)
  const form = new FormData()
  form.set('enable', String(on))
  try {
    await api.post(`panel/api/inbounds/setEnable/${ib.id}`, form)
  } catch {
    store.patch(ib.id, { enable: before })
    refused.add(ib.id)
    window.setTimeout(() => refused.delete(ib.id), 3400)
    toast('err', t('inbounds.toggleFailed'))
  } finally {
    toggling.delete(ib.id)
  }
}

// --- Page ----------------------------------------------------------------------------

const bannerClosed = ref(false)
const showBanner = computed(() => store.loaded && !bannerClosed.value && insecure())
const subscription = computed(() => !!store.defaults?.subEnable)
const open = ref<(string | number)[]>([])
</script>

<template>
  <AppPage :title="t('inbounds.title')">
    <template v-if="showBanner" #banner>
      <XNotice tone="err" :title="t('settings.securityTitle')" closable @close="bannerClosed = true">{{ t('xray.securityText') }}</XNotice>
    </template>

    <div v-if="!store.loaded && !store.failed" class="inbounds__loading" role="status"><XSpinner :size="16" /></div>
    <XCard v-else-if="store.failed">
      <XEmpty tone="err" icon="alert" :title="t('settings.loadFailed')" />
    </XCard>

    <template v-else>
      <XCard class="inbounds__summary">
        <div class="inbounds__stats">
          <XSummaryStat icon="updown" :label="t('inbounds.statSentReceived')" :value="`${formatBytes(totals.up)} / ${formatBytes(totals.down)}`" />
          <XSummaryStat icon="database" :label="t('inbounds.statTraffic')" :value="formatBytes(totals.up + totals.down)" />
          <XSummaryStat icon="clock" :label="t('inbounds.statAllTime')" :value="formatBytes(totals.allTime)" />
          <XSummaryStat icon="swap" :label="t('inbounds.statInbounds')" :value="String(store.list.length)" />
          <XSummaryStat icon="user" :label="t('inbounds.clients')" class="inbounds__clients">
            <ClientCounts :counts="allCounts" labelled />
          </XSummaryStat>
        </div>
      </XCard>

      <XExpandTable
        v-model:open="open"
        :columns="columns"
        :phone-columns="phoneColumns"
        :rows="rows"
        :row-key="(ib) => ib.id"
        :muted="(ib) => !ib.enable"
        :toggle-label="t('inbounds.clients')"
      >
        <template #tools>
          <InboundsTools v-model:mode="mode" v-model:search="search" v-model:filter="filter" :subscription="subscription" />
        </template>

        <template #cell-menu="{ row }">
          <XMenu v-slot="{ props: trigger }" :items="menu(row)" placement="bottom-start" @select="act($event, row)">
            <XInlineButton v-bind="trigger" icon="dots" :label="t('inbounds.colMenu')" />
          </XMenu>
        </template>
        <template #cell-enable="{ row }">
          <XSwitch
            :model-value="row.enable"
            :loading="toggling.has(row.id)"
            :invalid="refused.has(row.id)"
            :label="t('inbounds.enable')"
            @update:model-value="setEnable(row, $event)"
          />
        </template>
        <template #cell-remark="{ row }">
          <div class="inbounds__remark" :title="row.remark">{{ row.remark }}</div>
        </template>
        <template #cell-port="{ row }">
          <span class="inbounds__port">{{ row.port }}</span>
        </template>
        <template #cell-protocol="{ row }">
          <div class="inbounds__chips">
            <XChip mono>{{ protocolChips(row).protocol }}</XChip>
            <XChip v-if="protocolChips(row).network" mono outline>{{ protocolChips(row).network }}</XChip>
            <XChip v-if="protocolChips(row).security" tone="ok" mono outline>{{ security(protocolChips(row).security!) }}</XChip>
          </div>
        </template>
        <template #cell-clients="{ row }">
          <ClientCounts v-if="hasClients(row) && countsOf(row)" :counts="countsOf(row)!" />
        </template>
        <template #cell-traffic="{ row }">
          <XTooltip :delay="0">
            <XChip :tone="toneOf(traffic(row))" mono interactive>{{ formatBytes(used(row)) }} / {{ row.total ? formatBytes(row.total) : '∞' }}</XChip>
            <template #content>
              <div class="inbounds__updown">↑{{ formatBytes(row.up) }} ↓{{ formatBytes(row.down) }}</div>
              <div v-if="row.total && used(row) < row.total" class="inbounds__remain">
                <span>{{ t('inbounds.remain') }}</span><span>{{ formatBytes(row.total - used(row)) }}</span>
              </div>
            </template>
          </XTooltip>
        </template>
        <template #cell-allTime="{ row }">
          <XChip mono>{{ formatBytes(row.allTime || used(row)) }}</XChip>
        </template>
        <template #cell-expiry="{ row }">
          <XChip v-if="!row.expiryTime" icon="infinity" :aria-label="t('inbounds.noExpiry')" />
          <XTooltip v-else :text="formatDate(row.expiryTime)" :delay="0">
            <XChip :tone="toneOf(expiry(row))" interactive>{{ formatRelative(row.expiryTime) }}</XChip>
          </XTooltip>
        </template>

        <template #phone-menu="{ row }">
          <XMenu v-slot="{ props: trigger }" :items="menu(row, true)" placement="bottom-start" @select="act($event, row)">
            <XButton v-bind="trigger" variant="ghost" icon="dots" :label="t('inbounds.colMenu')" />
          </XMenu>
        </template>
        <template #phone-remark="{ row }">
          <div :class="['inbounds__remark', 'inbounds__remark--phone', { 'is-off': !row.enable }]">{{ row.remark }}</div>
        </template>
        <template #phone-info="{ row }">
          <XPopover :label="t('inbounds.colInfo')">
            <template #default="{ props: trigger }">
              <XButton v-bind="trigger" icon="info" round :mark="row.enable ? undefined : 'pause'" :label="t('inbounds.colInfo')" />
            </template>
            <template #content>
              <XKeyValue
                plain
                class="inbounds__facts"
                :rows="[
                  { key: 'protocol', label: t('inbounds.colProtocol') },
                  { label: t('inbounds.colPort'), value: String(row.port) },
                  ...(hasClients(row) ? [{ key: 'clients', label: t('inbounds.clients') }] : []),
                  { key: 'traffic', label: t('inbounds.colTraffic') },
                  { key: 'expiry', label: t('inbounds.colExpiry') },
                  { label: t('inbounds.trafficReset'), value: row.trafficReset },
                ]"
              >
                <template #value-protocol>
                  <XChip mono>{{ protocolChips(row).protocol }}</XChip>
                  <XChip v-if="protocolChips(row).network" mono outline>{{ protocolChips(row).network }}</XChip>
                  <XChip v-if="protocolChips(row).security" tone="ok" mono outline>{{ security(protocolChips(row).security!) }}</XChip>
                </template>
                <template #value-clients>
                  <ClientCounts v-if="countsOf(row)" :counts="countsOf(row)!" />
                </template>
                <template #value-traffic>
                  <XTooltip :delay="0">
                    <XChip :tone="toneOf(traffic(row))" mono interactive>{{ formatBytes(used(row)) }} / {{ row.total ? formatBytes(row.total) : '∞' }}</XChip>
                    <template #content>
                      <div class="inbounds__updown">↑{{ formatBytes(row.up) }} ↓{{ formatBytes(row.down) }}</div>
                      <div v-if="row.total && used(row) < row.total" class="inbounds__remain">
                        <span>{{ t('inbounds.remain') }}</span><span>{{ formatBytes(row.total - used(row)) }}</span>
                      </div>
                    </template>
                  </XTooltip>
                </template>
                <template #value-expiry>
                  <XChip v-if="!row.expiryTime" icon="infinity" :aria-label="t('inbounds.noExpiry')" />
                  <XChip v-else :tone="expiry(row) === 'err' ? 'err' : 'on'" mono>{{ formatDate(row.expiryTime) }}</XChip>
                </template>
              </XKeyValue>
            </template>
          </XPopover>
        </template>

        <template #empty>
          <XEmpty icon="inbox" :title="t('inbounds.noData')" tab />
        </template>
        <template v-if="shown.length > 25" #foot>
          <XPager v-model="page" v-model:size="size" :pages="pages" :sizes="sizes" summary class="inbounds__pager" />
        </template>
      </XExpandTable>
    </template>
  </AppPage>
</template>

<style scoped>
.inbounds__loading {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 640px;
  color: var(--acc);
}

.inbounds__stats {
  display: grid;
  grid-template-columns: minmax(0, 1.35fr) minmax(0, .95fr) minmax(0, 1.15fr) minmax(0, 1fr) minmax(0, 1.6fr);
  gap: 18px;
}

.inbounds__remark {
  display: -webkit-box;
  overflow: hidden;
  color: var(--fg);
  font: 500 var(--fs-13) / 1.4 var(--font-sans);
  word-break: break-word;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
}

.inbounds__remark--phone {
  padding: 8px 8px 8px 28px;
  font-size: var(--fs-14);
}

.inbounds__remark--phone.is-off {
  color: var(--mut);
}

.inbounds__port {
  color: var(--fg2);
  font: 400 var(--fs-13) / var(--lh-tight) var(--font-mono);
}

.inbounds__chips {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}

.inbounds__updown {
  color: var(--fg);
  font: 500 var(--fs-12) / var(--lh-tight) var(--font-mono);
  white-space: nowrap;
}

.inbounds__remain {
  display: flex;
  justify-content: space-between;
  gap: 16px;
  margin-top: 6px;
  padding-top: 6px;
  border-top: 1px solid var(--div);
  color: var(--dim);
  font: 400 var(--fs-12) / var(--lh-tight) var(--font-sans);
}

.inbounds__remain span:last-child {
  color: var(--fg);
  font: 500 var(--fs-12) / var(--lh-tight) var(--font-mono);
}

.inbounds__facts {
  width: 276px;
}

.inbounds__pager {
  justify-content: flex-end;
}

@media (max-width: 767px) {
  .inbounds__stats {
    grid-template-columns: repeat(2, minmax(0, 1fr));
    gap: 14px 12px;
  }

  .inbounds__clients {
    grid-column: 1 / -1;
  }
}
</style>
