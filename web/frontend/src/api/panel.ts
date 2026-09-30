import { http, unwrap } from './http'

export type XrayState = 'running' | 'stop' | 'error'

export interface ServerStatus {
  cpu: number
  cpuCores: number
  logicalPro: number
  cpuSpeedMhz: number
  mem: { current: number; total: number }
  swap: { current: number; total: number }
  disk: { current: number; total: number }
  xray: { state: XrayState; errorMsg: string; version: string }
  uptime: number
  loads: number[]
  tcpCount: number
  udpCount: number
  netIO: { up: number; down: number }
  netTraffic: { sent: number; recv: number }
  publicIP: { ipv4: string; ipv6: string }
  appStats: { threads: number; mem: number; uptime: number }
  /**
   * Optional Telegram (telemt) proxy block. Not reported by the backend yet;
   * the dashboard card renders only when present.
   */
  telemt?: { state: XrayState; errorMsg?: string; version: string; online: number; uptime: number; mem: number }
}

export interface CpuPoint {
  t: number
  cpu: number
}

export interface PanelUpdateInfo {
  currentVersion: string
  latestVersion: string
  updateAvailable: boolean
}

export interface XrayLogEntry {
  DateTime: string
  FromAddress: string
  ToAddress: string
  Inbound: string
  Outbound: string
  Email: string
  /** 0 direct, 1 blocked, 2 proxy */
  Event: number
}

export interface XrayLogPage {
  entries: XrayLogEntry[]
  total: number
  page: number
  pageSize: number
  date: string
  dates: string[]
  email: string
  clients: string[]
}

export interface XrayLogQuery {
  date?: string
  page?: number
  email?: string
  filter?: string
  showDirect: boolean
  showBlocked: boolean
  showProxy: boolean
}

export interface CustomGeo {
  id: number
  type: 'geosite' | 'geoip'
  alias: string
  url: string
  lastUpdatedAt: number
}

export const auth = {
  login(username: string, password: string, twoFactorCode?: string) {
    return http.post<null>('login', { username, password, twoFactorCode })
  },
  twoFactorEnabled() {
    return unwrap(http.post<boolean>('getTwoFactorEnable'))
  },
}

const S = 'panel/api/server/'

export const server = {
  status: () => unwrap(http.get<ServerStatus>(S + 'status')),
  cpuHistory: (bucketSeconds: number) => unwrap(http.get<CpuPoint[]>(S + 'cpuHistory/' + bucketSeconds)),
  xrayVersions: () => unwrap(http.get<string[]>(S + 'getXrayVersion')),
  panelUpdateInfo: () => unwrap(http.get<PanelUpdateInfo>(S + 'getPanelUpdateInfo')),
  configJson: () => unwrap(http.get<unknown>(S + 'getConfigJson')),
  stopXray: () => http.post<null>(S + 'stopXrayService'),
  restartXray: () => http.post<null>(S + 'restartXrayService'),
  installXray: (version: string) => http.post<null>(S + 'installXray/' + encodeURIComponent(version)),
  updatePanel: () => http.post<null>(S + 'updatePanel'),
  updateGeofile: (fileName?: string) =>
    http.post<null>(S + 'updateGeofile' + (fileName ? '/' + encodeURIComponent(fileName) : '')),
  logs: (count: number, level: string, syslog: boolean) =>
    unwrap(http.post<string[]>(S + 'logs/' + count, { level, syslog })),
  xrayLogs: (q: XrayLogQuery) => unwrap(http.post<XrayLogPage>(S + 'xraylogs', { ...q })),
  downloadXrayLogs: (q: XrayLogQuery) => http.download(S + 'xraylogs/download', { ...q }),
  importDb: (file: File) => {
    const fd = new FormData()
    fd.append('db', file)
    return http.upload<null>(S + 'importDB', fd)
  },
  restartPanel: () => http.post<null>('panel/setting/restartPanel'),
}

export const inbounds = {
  /** Emails of clients currently online. */
  onlines: () => unwrap(http.post<string[] | null>('panel/api/inbounds/onlines')),
}

const G = 'panel/api/custom-geo/'

export const customGeo = {
  list: () => unwrap(http.get<CustomGeo[]>(G + 'list')),
  add: (f: Pick<CustomGeo, 'type' | 'alias' | 'url'>) => http.post<null>(G + 'add', f),
  update: (id: number, f: Pick<CustomGeo, 'type' | 'alias' | 'url'>) => http.post<null>(G + 'update/' + id, f),
  remove: (id: number) => http.post<null>(G + 'delete/' + id),
  download: (id: number) => http.post<null>(G + 'download/' + id),
  updateAll: () => http.post<unknown>(G + 'update-all'),
}

/** Files the panel knows how to refresh (mirrors the legacy UI list). */
export const GEOFILES = ['geosite.dat', 'geoip.dat', 'geosite_IR.dat', 'geoip_IR.dat', 'geosite_RU.dat', 'geoip_RU.dat']
