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

/** Per-client traffic counters stored next to an inbound. */
export interface ClientStat {
  id: number
  inboundId: number
  enable: boolean
  email: string
  up: number
  down: number
  allTime: number
  expiryTime: number
  total: number
  lastOnline: number
}

/** Row of /panel/api/inbounds/list (database/model.Inbound). */
export interface DbInbound {
  id: number
  up: number
  down: number
  total: number
  allTime: number
  remark: string
  enable: boolean
  expiryTime: number
  trafficReset: string
  clientStats: ClientStat[] | null
  listen: string
  port: number
  protocol: string
  settings: string
  streamSettings: string
  tag: string
  sniffing: string
}

/** Panel-wide defaults used by the inbounds page (subset of defaultSettings). */
export interface PanelDefaults {
  expireDiff: number
  trafficDiff: number
  subEnable: boolean
  subURI: string
  subJsonEnable: boolean
  subJsonURI: string
  remarkModel: string
  ipLimitEnable: boolean
}

const I = 'panel/api/inbounds/'
const enc = encodeURIComponent

export const inbounds = {
  list: () => unwrap(http.get<DbInbound[] | null>(I + 'list')),
  /** Emails of clients currently online. */
  onlines: () => unwrap(http.post<string[] | null>(I + 'onlines')),
  setEnable: (id: number, enable: boolean) => http.post<null>(I + 'setEnable/' + id, { enable }),
  remove: (id: number) => http.post<null>(I + 'del/' + id),
  /** `client` is the raw client object from the inbound settings JSON. */
  updateClient: (inboundId: number, clientId: string, client: Record<string, unknown>) =>
    http.post<null>(I + 'updateClient/' + enc(clientId), { id: inboundId, settings: JSON.stringify({ clients: [client] }) }),
  removeClient: (inboundId: number, email: string) => http.post<null>(I + inboundId + '/delClientByEmail/' + enc(email)),
  resetClientTraffic: (inboundId: number, email: string) => http.post<null>(I + inboundId + '/resetClientTraffic/' + enc(email)),
  resetClientTraffics: (inboundId: number) => http.post<null>(I + 'resetAllClientTraffics/' + inboundId),
  resetAllTraffics: () => http.post<null>(I + 'resetAllTraffics'),
  /** -1 cleans every inbound. */
  delDepletedClients: (inboundId: number) => http.post<null>(I + 'delDepletedClients/' + inboundId),
}

export const settings = {
  defaults: () => unwrap(http.post<PanelDefaults>('panel/setting/defaultSettings')),
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
