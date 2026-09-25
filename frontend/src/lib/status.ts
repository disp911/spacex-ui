import { t } from './i18n'

export type ServiceState = 'running' | 'stop' | 'error' | 'unknown'

export function normState(s: string | undefined): ServiceState {
  return s === 'running' || s === 'stop' || s === 'error' ? s : 'unknown'
}

export function stateColor(s: ServiceState): string {
  return {
    running: 'var(--accent)',
    stop: 'var(--warn-state)',
    error: 'var(--danger)',
    unknown: 'var(--unknown)',
  }[s]
}

export function stateLabel(s: ServiceState): string {
  return {
    running: t('pages.index.xrayStatusRunning'),
    stop: t('pages.index.xrayStatusStop'),
    error: t('pages.index.xrayStatusError'),
    unknown: t('pages.index.xrayStatusUnknown'),
  }[s]
}

/** Server status as returned by panel/api/server/status (web/service/server.go). */
export interface ServerStatus {
  cpu: number
  cpuCores: number
  logicalPro: number
  cpuSpeedMhz: number
  mem: { current: number; total: number }
  swap: { current: number; total: number }
  disk: { current: number; total: number }
  xray: { state: string; errorMsg: string; version: string }
  uptime: number
  loads: number[]
  tcpCount: number
  udpCount: number
  netIO: { up: number; down: number }
  netTraffic: { sent: number; recv: number }
  publicIP: { ipv4: string; ipv6: string }
  appStats: { threads: number; mem: number; uptime: number }
}
