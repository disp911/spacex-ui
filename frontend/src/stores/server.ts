import { defineStore } from 'pinia'
import { onScopeDispose, ref } from 'vue'
import { api } from '../api'
import { subscribe } from '../live'

export type ProcessState = 'running' | 'stop' | 'error'

export interface ServiceStatus {
  state: ProcessState
  errorMsg: string
  version: string
  mem: number
  uptime: number
  online: number
}

/** Answer of panel/api/server/status (web/service/server.go Status). */
export interface ServerStatus {
  cpu: number
  cpuCores: number
  logicalPro: number
  cpuSpeedMhz: number
  mem: { current: number; total: number }
  swap: { current: number; total: number }
  disk: { current: number; total: number }
  xray: ServiceStatus
  /** Absent when there is no mtproto inbound. */
  telemt?: ServiceStatus
  uptime: number
  loads: number[]
  tcpCount: number
  udpCount: number
  netIO: { up: number; down: number }
  netTraffic: { sent: number; recv: number }
  publicIP: { ipv4: string; ipv6: string }
  /** The panel's own process: uptime in seconds, memory, goroutines. */
  appStats: { uptime: number; mem: number; threads: number }
}

/** Live server status while someone watches it. */
export const useServerStatus = defineStore('serverStatus', () => {
  const status = ref<ServerStatus | null>(null)
  /** The first status could not be fetched; screens stop showing it as loading. */
  const failed = ref(false)
  let watchers = 0
  let unsubscribe: (() => void) | null = null
  let updatedAt = 0

  function set(next: ServerStatus) {
    status.value = next
    failed.value = false
    updatedAt = Date.now()
  }

  async function refresh() {
    try {
      set(await api.get<ServerStatus>('panel/api/server/status'))
    } catch (e) {
      if (!status.value) failed.value = true
      throw e
    }
  }

  /** Starts live updates for the calling component's lifetime. */
  function watch() {
    if (watchers++ === 0) {
      refresh().catch(() => {})
      const off1 = subscribe('status', (payload) => set(payload as ServerStatus))
      // After a reconnect the last pushed status may be stale.
      const off2 = subscribe('reconnected', () => refresh().catch(() => {}))
      // Without the socket (blocked by a proxy) the status is polled.
      const poll = window.setInterval(() => {
        if (Date.now() - updatedAt > 6000) refresh().catch(() => {})
      }, 3000)
      unsubscribe = () => (off1(), off2(), clearInterval(poll))
    }
    onScopeDispose(() => {
      if (--watchers === 0) {
        unsubscribe?.()
        unsubscribe = null
      }
    })
  }

  return { status, failed, refresh, watch }
})
