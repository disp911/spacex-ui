import { onBeforeUnmount, onMounted, ref, shallowRef, watch } from 'vue'
import { inbounds, server, type CpuPoint, type ServerStatus, type XrayState } from '@/api/panel'
import { PanelSocket } from '@/api/ws'

export type CpuRange = '2m' | '30m' | '1h' | '2h' | '3h' | '5h'

/** Range → aggregation bucket in seconds (the backend returns 60 buckets). */
export const CPU_RANGES: Record<CpuRange, number> = {
  '2m': 2,
  '30m': 30,
  '1h': 60,
  '2h': 120,
  '3h': 180,
  '5h': 300,
}

/**
 * Live server status for the dashboard: WebSocket push with a polling
 * fallback, the online-clients count and the CPU history for the chart.
 */
export function useServerStatus() {
  const status = shallowRef<ServerStatus | null>(null)
  const online = ref<number | null>(null)
  const cpuRange = ref<CpuRange>('2m')
  const cpuSeries = ref<CpuPoint[]>([])

  let socket: PanelSocket | null = null
  let pollTimer: number | undefined
  let onlineTimer: number | undefined
  let historyTimer: number | undefined
  let fallbackTimer: number | undefined
  let disposed = false

  async function refresh() {
    try {
      status.value = await server.status()
    } catch {
      /* keep the last known status */
    }
  }

  async function refreshOnline() {
    try {
      online.value = (await inbounds.onlines())?.length ?? 0
    } catch {
      /* ignore */
    }
  }

  async function refreshHistory() {
    try {
      cpuSeries.value = (await server.cpuHistory(CPU_RANGES[cpuRange.value])) ?? []
    } catch {
      cpuSeries.value = []
    }
  }

  function startPolling() {
    if (pollTimer || disposed) return
    pollTimer = window.setInterval(refresh, 2000)
  }

  function setXrayState(state: XrayState, errorMsg = '') {
    if (!status.value) return
    status.value = { ...status.value, xray: { ...status.value.xray, state, errorMsg } }
  }

  function scheduleHistory() {
    clearInterval(historyTimer)
    // Short ranges move quickly; long ones only gain a point every few minutes.
    const every = cpuRange.value === '2m' ? 4000 : 30000
    historyTimer = window.setInterval(refreshHistory, every)
  }

  watch(cpuRange, () => {
    refreshHistory()
    scheduleHistory()
  })

  onMounted(async () => {
    await refresh()
    if (disposed) return
    refreshOnline()
    refreshHistory()
    scheduleHistory()
    onlineTimer = window.setInterval(refreshOnline, 10000)

    socket = new PanelSocket()
    socket.on('status', (p: ServerStatus) => (status.value = p))
    socket.on('xray_state', (p: { state: XrayState; errorMsg?: string }) => setXrayState(p.state, p.errorMsg))
    socket.on('down', startPolling)
    socket.connect()
    // If the socket never delivers, poll anyway.
    fallbackTimer = window.setTimeout(() => {
      if (!socket?.connected) startPolling()
    }, 5000)
  })

  onBeforeUnmount(() => {
    disposed = true
    clearTimeout(fallbackTimer)
    socket?.close()
    clearInterval(pollTimer)
    clearInterval(onlineTimer)
    clearInterval(historyTimer)
  })

  return { status, online, cpuRange, cpuSeries, refresh, setXrayState }
}
