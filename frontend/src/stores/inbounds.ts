import { defineStore } from 'pinia'
import { computed, ref, shallowRef, triggerRef } from 'vue'
import { api } from '../api'
import { subscribe } from '../live'
import { countClients, type ClientCounts, type ClientStat, type DbInbound, type InboundDefaults } from '../pages/inbounds/model'

/** Delta the server pushes as client_stats: absolute counters of the clients active this cycle and inbound totals. */
interface StatsDelta {
  clients?: (Partial<ClientStat> & { inboundId: number; email: string })[]
  inbounds?: { id: number; up: number; down: number; total?: number; allTime?: number; enable?: boolean }[]
}

/** What the server pushes as traffic; only who is online and when each client was last seen matter here. */
interface TrafficPush {
  onlineClients?: string[] | null
  lastOnlineMap?: Record<string, number> | null
}

const refreshKey = 'inbounds.refresh'
const intervalKey = 'inbounds.refreshInterval'

function read(key: string): string | null {
  try {
    return localStorage.getItem(key)
  } catch {
    return null
  }
}

function write(key: string, value: string) {
  try {
    localStorage.setItem(key, value)
  } catch {
    // the choice lasts for this visit only
  }
}

/**
 * The inbounds with their clients' traffic, who is online and the panel's
 * defaults, while a page watches them. The server pushes changes over the
 * WebSocket; with auto refresh on, the list is also fetched again every
 * interval, as the old panel did.
 */
export const useInbounds = defineStore('inbounds', () => {
  const list = shallowRef<DbInbound[]>([])
  const online = shallowRef(new Set<string>())
  const lastOnline = shallowRef<Record<string, number>>({})
  const defaults = ref<InboundDefaults | null>(null)
  const loaded = ref(false)
  const failed = ref(false)
  const refreshing = ref(false)

  const autoRefresh = ref(read(refreshKey) === '1')
  const interval = ref(Number(read(intervalKey)) || 5000)

  const counts = computed(() => {
    const d = defaults.value ?? { expireDiff: 0, trafficDiff: 0 }
    const now = Date.now()
    return new Map<number, ClientCounts>(list.value.map((ib) => [ib.id, countClients(ib, online.value, d, now)]))
  })

  async function loadDefaults() {
    defaults.value = await api.post<InboundDefaults>('panel/setting/defaultSettings')
  }

  /** The list, who is online and when each was last seen, fetched again. */
  async function load() {
    refreshing.value = true
    try {
      const [inbounds, onlines, last] = await Promise.all([
        api.get<DbInbound[]>('panel/api/inbounds/list'),
        api.post<string[] | null>('panel/api/inbounds/onlines'),
        api.post<Record<string, number> | null>('panel/api/inbounds/lastOnline'),
      ])
      list.value = inbounds ?? []
      online.value = new Set(onlines ?? [])
      lastOnline.value = last ?? {}
      loaded.value = true
      failed.value = false
    } catch (e) {
      if (!loaded.value) failed.value = true
      throw e
    } finally {
      refreshing.value = false
    }
  }

  function applyDelta(delta: StatsDelta) {
    const byId = new Map(list.value.map((ib) => [ib.id, ib]))
    for (const c of delta.clients ?? []) {
      const s = byId.get(c.inboundId)?.clientStats?.find((x) => x.email === c.email)
      if (s) Object.assign(s, Object.fromEntries(Object.entries(c).filter(([, v]) => v !== undefined)))
    }
    for (const t of delta.inbounds ?? []) {
      const ib = byId.get(t.id)
      if (!ib) continue
      ib.up = t.up
      ib.down = t.down
      if (t.total !== undefined) ib.total = t.total
      if (t.allTime !== undefined) ib.allTime = t.allTime
      if (t.enable !== undefined) ib.enable = t.enable
    }
    triggerRef(list)
  }

  /** Changes an inbound in place, e.g. its switch applied before the server answers. */
  function patch(id: number, fields: Partial<DbInbound>) {
    const ib = list.value.find((x) => x.id === id)
    if (!ib) return
    Object.assign(ib, fields)
    triggerRef(list)
  }

  function applyTraffic(push: TrafficPush) {
    // Null is everyone gone offline, not "no news".
    if ('onlineClients' in push) online.value = new Set(push.onlineClients ?? [])
    if (push.lastOnlineMap && typeof push.lastOnlineMap === 'object') lastOnline.value = push.lastOnlineMap
  }

  // --- Watching -------------------------------------------------------------------

  let watchers = 0
  let stop: (() => void)[] = []
  let loop = 0
  let invalidated = 0

  function schedule() {
    clearTimeout(loop)
    if (!autoRefresh.value || !watchers) return
    loop = window.setTimeout(async () => {
      await load().catch(() => {})
      schedule()
    }, interval.value)
  }

  /** Starts live updates for the calling page's lifetime; returns the function that ends them. */
  function watch() {
    if (watchers++ === 0) {
      loadDefaults().catch(() => {})
      load()
        .catch(() => {})
        .finally(() => {
          stop = [
            subscribe('inbounds', (p) => Array.isArray(p) && (list.value = p as DbInbound[])),
            subscribe('client_stats', (p) => p && typeof p === 'object' && applyDelta(p as StatsDelta)),
            subscribe('traffic', (p) => p && typeof p === 'object' && applyTraffic(p as TrafficPush)),
            // Too large to push: fetched instead, once the burst of changes is over.
            subscribe('invalidate', (p) => {
              const type = (p as { type?: string } | null)?.type
              if (type !== 'inbounds' && type !== 'traffic') return
              clearTimeout(invalidated)
              invalidated = window.setTimeout(() => load().catch(() => {}), 1000)
            }),
            // Pushes sent while the connection was down are lost.
            subscribe('reconnected', () => load().catch(() => {})),
          ]
          schedule()
        })
    }
    return () => {
      if (--watchers > 0) return
      stop.forEach((off) => off())
      stop = []
      clearTimeout(loop)
      clearTimeout(invalidated)
    }
  }

  function setAutoRefresh(on: boolean) {
    autoRefresh.value = on
    write(refreshKey, on ? '1' : '0')
    schedule()
  }

  function setRefreshInterval(ms: number) {
    interval.value = ms
    write(intervalKey, String(ms))
    schedule()
  }

  return { list, online, lastOnline, defaults, counts, loaded, failed, refreshing, autoRefresh, interval, load, patch, watch, setAutoRefresh, setRefreshInterval }
})
