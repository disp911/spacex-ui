<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue'
import XModal from '@/components/ui/XModal.vue'
import XButton from '@/components/ui/XButton.vue'
import XSelect from '@/components/ui/XSelect.vue'
import XInput from '@/components/ui/XInput.vue'
import XCheckChip from '@/components/ui/XCheckChip.vue'
import XEmpty from '@/components/ui/XEmpty.vue'
import XSheet from '@/components/ui/XSheet.vue'
import XPagination from '@/components/ui/XPagination.vue'
import XIcon from '@/components/ui/XIcon.vue'
import { server, type XrayLogEntry, type XrayLogPage, type XrayLogQuery } from '@/api/panel'
import { saveBlob } from '@/api/http'
import { toastError } from '@/composables/useToast'
import { useIsMobile } from '@/composables/useMedia'
import { formatClock, formatCount, formatDay, formatLogTime } from '@/utils/format'
import { t } from '@/i18n'

const open = defineModel<boolean>('open', { default: false })
const isMobile = useIsMobile()

type Kind = 'direct' | 'blocked' | 'proxy'
const KIND_OF_EVENT: Kind[] = ['direct', 'blocked', 'proxy']
const KINDS: { id: Kind; label: string; tone: 'accent' | 'danger' | 'info' }[] = [
  { id: 'direct', label: 'Direct', tone: 'accent' },
  { id: 'blocked', label: 'Blocked', tone: 'danger' },
  { id: 'proxy', label: 'Proxy', tone: 'info' },
]

const f = reactive({ date: '', email: '', filter: '', page: 1, kinds: { direct: true, blocked: true, proxy: true } as Record<Kind, boolean> })
const data = ref<XrayLogPage | null>(null)
const loading = ref(false)
const downloading = ref(false)
const kindsSheet = ref(false)

const query = (): XrayLogQuery => ({
  date: f.date,
  page: f.page,
  email: f.email,
  filter: f.filter.trim(),
  showDirect: f.kinds.direct,
  showBlocked: f.kinds.blocked,
  showProxy: f.kinds.proxy,
})

let seq = 0
// Set while adopting the server's default day, so the filter watcher
// doesn't fire a second, identical request.
let adopting = false
async function load() {
  const my = ++seq
  loading.value = true
  try {
    const res = await server.xrayLogs(query())
    if (my !== seq) return
    data.value = res
    if (!f.date && res.date) {
      adopting = true
      f.date = res.date
    }
  } catch (e) {
    toastError(e)
  } finally {
    if (my === seq) loading.value = false
  }
}

watch(open, (v) => v && load())
// Any filter change goes back to page 1; the search box is debounced.
let debounce: number | undefined
watch(
  () => f.filter,
  () => {
    clearTimeout(debounce)
    debounce = window.setTimeout(() => (f.page === 1 ? load() : (f.page = 1)), 300)
  },
)
watch(
  () => [f.date, f.email, f.kinds.direct, f.kinds.blocked, f.kinds.proxy],
  () => {
    if (adopting) {
      adopting = false
      return
    }
    if (f.page === 1) load()
    else f.page = 1
  },
)
watch(() => f.page, load)

const total = computed(() => data.value?.total ?? 0)
const pageSize = computed(() => data.value?.pageSize || 200)
const pages = computed(() => Math.max(1, Math.ceil(total.value / pageSize.value)))
const rows = computed(() =>
  (data.value?.entries ?? []).map((e: XrayLogEntry, i) => ({ ...e, key: i, kind: KIND_OF_EVENT[e.Event] ?? 'direct' })),
)

const dateOptions = computed(() => (data.value?.dates ?? []).map((d) => ({ value: d, label: formatDay(d) })))
const clientOptions = computed(() => [
  { value: '', label: t('xlog.allClients') },
  ...(data.value?.clients ?? []).map((c) => ({ value: c, label: c })),
])
const kindsOn = computed(() => KINDS.filter((k) => f.kinds[k.id]).length)
const kindsLabel = computed(() =>
  kindsOn.value === 3 ? t('xlog.allKinds') : kindsOn.value === 0 ? t('xlog.noKinds') : kindsOn.value === 1 ? KINDS.find((k) => f.kinds[k.id])!.label : t('xlog.kindsOf', { n: kindsOn.value }),
)

const subtitle = computed(() => {
  if (!data.value?.date) return ''
  return isMobile.value
    ? t('xlog.subtitleShort', { n: formatCount(total.value), day: formatDay(data.value.date) })
    : t('xlog.subtitle', { n: formatCount(total.value), day: formatDay(data.value.date) })
})
const rangeLabel = computed(() => {
  if (!total.value) return t('xlog.noRows')
  const from = (f.page - 1) * pageSize.value + 1
  return t('xlog.range', {
    from: formatCount(from),
    to: formatCount(from + rows.value.length - 1),
    total: formatCount(total.value),
    page: f.page,
    pages: pages.value,
  })
})

async function download() {
  downloading.value = true
  try {
    const { blob, filename } = await server.downloadXrayLogs(query())
    saveBlob(blob, filename)
  } catch (e) {
    toastError(e)
  } finally {
    downloading.value = false
  }
}
</script>

<template>
  <XModal v-model:open="open" :title="isMobile ? t('xlog.titleShort') : t('xlog.title')" :subtitle="subtitle" :width="1368" flush tall>
    <template #toolbar>
      <!-- Desktop toolbar -->
      <div v-if="!isMobile" class="tb">
        <XButton icon="restart" icon-only size="sm" :label="t('common.refresh')" :loading="loading" @click="load" />
        <XSelect v-model="f.date" :options="dateOptions" :label="t('xlog.date')" mono :menu-width="176" />
        <XSelect v-model="f.email" :options="clientOptions" :label="t('xlog.client')" mono :menu-width="212" />
        <XInput v-model="f.filter" size="sm" icon="search" mono clearable :placeholder="t('xlog.search')" class="tb__search" />
        <XCheckChip v-for="k in KINDS" :key="k.id" v-model="f.kinds[k.id]" :label="k.label" :tone="k.tone" />
        <XButton icon="download" size="sm" :loading="downloading" @click="download">{{ t('common.download') }}</XButton>
      </div>
      <!-- Mobile toolbar -->
      <template v-else>
        <XInput v-model="f.filter" size="md" icon="search" mono clearable :placeholder="t('xlog.searchShort')" />
        <div class="tb-m">
          <XSelect v-model="f.date" :options="dateOptions" :sheet-title="t('xlog.date')" block mono size="md" />
          <XSelect v-model="f.email" :options="clientOptions" :sheet-title="t('xlog.client')" block size="md" />
          <button type="button" class="tb-m__kinds" @click="kindsSheet = true">
            <span class="ellipsis">{{ kindsLabel }}</span>
            <XIcon name="chevron-down" :size="10" :stroke-width="1.6" />
          </button>
        </div>
      </template>
    </template>

    <!-- Desktop table -->
    <div v-if="!isMobile" class="table">
      <div class="table__head caps">
        <span class="c-date">{{ t('xlog.colDate') }}</span>
        <span class="c-from">{{ t('xlog.colFrom') }}</span>
        <span class="c-to">{{ t('xlog.colTo') }}</span>
        <span class="c-in">{{ t('xlog.colInbound') }}</span>
        <span class="c-out">{{ t('xlog.colOutbound') }}</span>
        <span class="c-email">{{ t('xlog.colEmail') }}</span>
      </div>
      <div class="table__body" :class="{ 'is-loading': loading }">
        <div v-for="r in rows" :key="r.key" class="row" :class="`is-${r.kind}`">
          <span class="c-date dim">{{ formatLogTime(r.DateTime) }}</span>
          <span class="c-from ellipsis">{{ r.FromAddress }}</span>
          <span class="c-to ellipsis">{{ r.ToAddress }}</span>
          <span class="c-in ellipsis dim2">{{ r.Inbound }}</span>
          <span class="c-out"><span class="dot" /><span class="ellipsis">{{ r.Outbound }}</span></span>
          <span class="c-email ellipsis strong">{{ r.Email }}</span>
        </div>
        <XEmpty v-if="!rows.length && !loading" :title="t('xlog.emptyTitle')" :text="data?.dates?.length ? t('xlog.emptyText') : t('xlog.noLog')" />
      </div>
    </div>

    <!-- Mobile cards -->
    <div v-else class="cards" :class="{ 'is-loading': loading }">
      <div v-for="r in rows" :key="r.key" class="card" :class="`is-${r.kind}`">
        <div class="card__top">
          <span class="dot" />
          <span class="card__dst ellipsis">{{ r.ToAddress }}</span>
          <span class="card__time">{{ formatClock(r.DateTime) }}</span>
        </div>
        <div class="card__bottom">
          <span class="card__email ellipsis">{{ r.Email }}</span>
          <span class="card__in ellipsis">{{ r.Inbound }}</span>
          <span class="card__out">{{ r.Outbound }}</span>
        </div>
      </div>
      <XEmpty v-if="!rows.length && !loading" :title="t('xlog.emptyTitle')" :text="data?.dates?.length ? t('xlog.emptyText') : t('xlog.noLog')" />
    </div>

    <template #footer>
      <span class="range">{{ isMobile ? (total ? t('xlog.pageShort', { page: f.page, pages }) : t('xlog.noRows')) : rangeLabel }}</span>
      <XPagination v-model="f.page" :pages="pages" :compact="isMobile" />
    </template>
  </XModal>

  <XSheet v-model:open="kindsSheet" :title="t('xlog.kinds')">
    <XCheckChip v-for="k in KINDS" :key="k.id" v-model="f.kinds[k.id]" :label="k.label" :tone="k.tone" size="lg" />
    <template #footer>
      <XButton variant="primary" size="lg" block @click="kindsSheet = false">
        {{ total ? t('xlog.show', { n: formatCount(total) }) : t('common.nothingFound') }}
      </XButton>
    </template>
  </XSheet>
</template>

<style scoped>
.tb {
  width: 100%;
  display: flex;
  align-items: center;
  gap: 9px;
}
.tb__search {
  flex: 1;
  min-width: 160px;
}
.tb-m {
  display: flex;
  gap: 6px;
}
.tb-m__kinds {
  flex: 1;
  min-width: 0;
  height: var(--h-md);
  padding: 0 8px 0 10px;
  display: flex;
  align-items: center;
  gap: 6px;
  border-radius: var(--r-md);
  border: 1px solid var(--border-control);
  background: var(--control);
  color: var(--text-2);
  font: var(--fw-medium) var(--fs-sm) var(--font-sans);
}
.tb-m__kinds span {
  flex: 1;
  text-align: left;
}
.tb-m__kinds :deep(svg) {
  color: var(--text-4);
}

/* Kind colors shared by table rows and cards */
.is-direct { --ink: var(--text-2); --dot: var(--muted-dot); --tint: transparent; --out-border: var(--border-control); }
.is-blocked { --ink: var(--danger-text); --dot: var(--danger); --tint: var(--danger-soft); --out-border: var(--danger-border); }
.is-proxy { --ink: var(--info-text); --dot: var(--info); --tint: var(--info-soft); --out-border: var(--info-border); }
.dot {
  width: 6px;
  height: 6px;
  flex: none;
  border-radius: 50%;
  background: var(--dot);
}

.table {
  flex: 1;
  min-height: 0;
  display: flex;
  flex-direction: column;
  background: var(--bg);
  border: 1px solid var(--border-subtle);
  border-radius: var(--r-md);
  overflow: hidden;
}
.table__head,
.row {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 0 14px;
  white-space: nowrap;
}
.table__head {
  flex: none;
  height: 34px;
  background: var(--field);
  border-bottom: 1px solid var(--border-subtle);
  font-size: var(--fs-2xs);
}
.table__body {
  flex: 1;
  min-height: 0;
  overflow: auto;
  transition: opacity var(--dur-fast) ease;
}
.is-loading {
  opacity: 0.55;
}
.row {
  height: 34px;
  border-bottom: 1px solid var(--divider);
  background: color-mix(in srgb, var(--tint) 60%, transparent);
  color: var(--ink);
  font: var(--fw-regular) var(--fs-sm) var(--font-mono);
  transition: background var(--dur-fast) ease;
}
.row:hover {
  background: var(--hover);
}
.c-date { width: 124px; flex: none; }
.c-from { width: 150px; flex: none; }
.c-to { flex: 1; min-width: 0; }
.c-in { width: 150px; flex: none; }
.c-out { width: 132px; flex: none; display: flex; align-items: center; gap: 6px; overflow: hidden; }
.c-email { width: 150px; flex: none; }
.dim { opacity: 0.7; }
.dim2 { opacity: 0.85; }
.strong { font-weight: var(--fw-medium); }

.range {
  flex: 1;
  min-width: 0;
  font: var(--fw-regular) var(--fs-sm) var(--font-sans);
  color: var(--text-3);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.cards {
  display: flex;
  flex-direction: column;
  gap: 8px;
  transition: opacity var(--dur-fast) ease;
}
.card {
  padding: 11px 12px;
  border-radius: var(--r-md);
  border: 1px solid var(--border-field);
  background: var(--field);
  color: var(--ink);
}
.card.is-blocked,
.card.is-proxy {
  background: var(--tint);
}
.card__top {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 7px;
}
.card__top .dot {
  width: 7px;
  height: 7px;
}
.card__dst {
  flex: 1;
  min-width: 0;
  font: var(--fw-medium) var(--fs-md) var(--font-mono);
}
.card__time {
  flex: none;
  font: var(--fw-regular) var(--fs-xs) var(--font-mono);
  color: var(--text-4);
}
.card__bottom {
  display: flex;
  align-items: center;
  gap: 8px;
}
.card__email {
  flex: none;
  max-width: 120px;
  font: var(--fw-medium) var(--fs-xs) var(--font-sans);
  color: var(--text-2);
}
.card__in {
  flex: 1;
  min-width: 0;
  font: var(--fw-regular) var(--fs-xs) var(--font-mono);
  color: var(--text-4);
}
.card__out {
  flex: none;
  padding: 3px 7px;
  border-radius: var(--r-xs);
  border: 1px solid var(--out-border);
  font: var(--fw-medium) var(--fs-2xs) var(--font-mono);
  letter-spacing: 0.04em;
  white-space: nowrap;
}
</style>
