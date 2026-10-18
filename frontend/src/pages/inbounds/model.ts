// The panel's inbounds as GET panel/api/inbounds/list answers them, read the
// way the old panel read them (web/assets/js/model/dbinbound.js, inbound.js and
// the inbounds page): clients, their traffic and their state.

type Json = Record<string, any>

export const protocols = ['vmess', 'vless', 'trojan', 'shadowsocks', 'wireguard', 'hysteria', 'mixed', 'http', 'tunnel', 'tun', 'mtproto'] as const
export type Protocol = (typeof protocols)[number]

/** Traffic of one client as the server counts it (xray/client_traffic.go). */
export interface ClientStat {
  id: number
  inboundId: number
  enable: boolean
  email: string
  up: number
  down: number
  allTime: number
  /** 0 never; below 0 a delayed start: that many milliseconds from the first use. */
  expiryTime: number
  /** Traffic limit in bytes, 0 without one. */
  total: number
  /** Days of the period that renews itself, 0 when it does not. */
  reset: number
  lastOnline: number
}

/** An inbound as stored (database/model/model.go Inbound); settings, streamSettings and sniffing are JSON text. */
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
  lastTrafficResetTime: number
  clientStats: ClientStat[] | null
  listen: string
  port: number
  protocol: Protocol
  settings: string
  streamSettings: string
  tag: string
  sniffing: string
}

/** A client as the inbound's settings hold it; the key that tells it apart depends on the protocol. */
export interface Client extends Json {
  email: string
  enable: boolean
  limitIp?: number
  totalGB?: number
  expiryTime?: number
  subId?: string
  comment?: string
  reset?: number
}

/** What POST panel/setting/defaultSettings tells the inbounds page. */
export interface InboundDefaults {
  /** Days and gigabytes before the end at which a client counts as nearly depleted. */
  expireDiff: number
  trafficDiff: number
  pageSize: number
  defaultCert: string
  defaultKey: string
  subEnable: boolean
  subJsonEnable: boolean
  subClashEnable: boolean
  subTitle: string
  subURI: string
  subJsonURI: string
  subClashURI: string
  remarkModel: string
  ipLimitEnable: boolean
}

const GB = 1024 ** 3
const DAY = 86_400_000

const parse = (text: string): Json => {
  if (!text) return {}
  try {
    return JSON.parse(text) as Json
  } catch {
    return {}
  }
}

export const settingsOf = (ib: DbInbound) => parse(ib.settings)
export const streamOf = (ib: DbInbound) => parse(ib.streamSettings)
export const sniffingOf = (ib: DbInbound) => parse(ib.sniffing)

/** Shadowsocks with this method has one key for everyone, so no clients. */
const singleUserMethod = '2022-blake3-chacha20-poly1305'

/** Whether the inbound has clients of its own. */
export function hasClients(ib: DbInbound, settings = settingsOf(ib)): boolean {
  switch (ib.protocol) {
    case 'vmess':
    case 'vless':
    case 'trojan':
    case 'hysteria':
    case 'mtproto':
      return true
    case 'shadowsocks':
      return settings.method !== singleUserMethod
    default:
      return false
  }
}

export function clientsOf(ib: DbInbound, settings = settingsOf(ib)): Client[] {
  return hasClients(ib, settings) && Array.isArray(settings.clients) ? (settings.clients as Client[]) : []
}

/** The key a client is changed and deleted by on the server. */
export function clientKey(protocol: Protocol, c: Client): string {
  switch (protocol) {
    case 'trojan':
      return String(c.password ?? '')
    case 'shadowsocks':
      return c.email
    case 'hysteria':
      return String(c.auth ?? '')
    default:
      return String(c.id ?? '')
  }
}

/** Whether the inbound's clients get share links. */
export const hasLink = (p: Protocol) => ['vmess', 'vless', 'trojan', 'shadowsocks', 'hysteria', 'mtproto'].includes(p)

/** The address its links point at: the listen address when it is a specific one, the panel's host otherwise. */
export const addressOf = (ib: DbInbound) => (ib.listen && ib.listen !== '0.0.0.0' ? ib.listen : location.hostname)

export function statsMap(ib: DbInbound): Map<string, ClientStat> {
  return new Map((ib.clientStats ?? []).map((s) => [s.email, s]))
}

// --- Client state ---------------------------------------------------------------

export const usedOf = (s: ClientStat | undefined) => (s ? s.up + s.down : 0)

/** All-time traffic, never below what the live counters show. */
export const allTimeOf = (s: ClientStat | undefined) => Math.max(s?.allTime ?? 0, usedOf(s))

export const remainingOf = (s: ClientStat | undefined) => (s && s.total > 0 ? Math.max(0, s.total - usedOf(s)) : 0)

/** Out of traffic or past its date. */
export function depleted(s: ClientStat | undefined, now = Date.now()): boolean {
  if (!s) return false
  const exhausted = s.total > 0 && usedOf(s) >= s.total
  const expired = s.expiryTime > 0 && s.expiryTime <= now
  return exhausted || expired
}

/** Close to its traffic limit or its date, by the panel's settings; not yet depleted. */
export function expiring(s: ClientStat | undefined, d: Pick<InboundDefaults, 'expireDiff' | 'trafficDiff'>, now = Date.now()): boolean {
  if (!s || depleted(s, now)) return false
  return (s.expiryTime > 0 && s.expiryTime - now < d.expireDiff * DAY) || (s.total > 0 && s.total - usedOf(s) < d.trafficDiff * GB)
}

/** Days of a delayed start, or null when the date is fixed or absent. */
export const delayedDays = (expiryTime: number) => (expiryTime < 0 ? expiryTime / -DAY : null)

/** The clients of an inbound by state, as the old panel counted them; a disabled inbound counts every client as disabled. */
export interface ClientCounts {
  clients: number
  active: string[]
  deactive: string[]
  depleted: string[]
  expiring: string[]
  online: string[]
  comments: Map<string, string>
}

export function countClients(ib: DbInbound, online: Set<string>, d: Pick<InboundDefaults, 'expireDiff' | 'trafficDiff'>, now = Date.now()): ClientCounts {
  const clients = clientsOf(ib)
  const counts: ClientCounts = { clients: clients.length, active: [], deactive: [], depleted: [], expiring: [], online: [], comments: new Map() }
  if (!ib.enable) {
    counts.deactive = clients.map((c) => c.email)
    return counts
  }
  for (const c of clients) {
    if (c.comment) counts.comments.set(c.email, c.comment)
    if (c.enable) {
      counts.active.push(c.email)
      if (online.has(c.email)) counts.online.push(c.email)
    } else counts.deactive.push(c.email)
  }
  for (const s of ib.clientStats ?? []) {
    if (depleted(s, now)) counts.depleted.push(s.email)
    else if (expiring(s, d, now)) counts.expiring.push(s.email)
  }
  return counts
}

// --- Search -----------------------------------------------------------------------

/** Whether text occurs, case aside, in any value of a JSON tree, as the old panel searched. */
export function deepSearch(value: unknown, text: string): boolean {
  if (value == null || value === '') return false
  if (Array.isArray(value)) return value.some((v) => deepSearch(v, text))
  if (typeof value === 'object') return Object.values(value).some((v) => deepSearch(v, text))
  return String(value).toLowerCase().includes(text.toLowerCase())
}

// --- How the list shows them -----------------------------------------------------

export type State = 'ok' | 'warn' | 'err' | 'unl'

/** Traffic against its limit: none, plenty, close by the panel's threshold, or used up. */
export function trafficState(used: number, total: number, d: Pick<InboundDefaults, 'trafficDiff'>): State {
  if (!total) return 'unl'
  if (used >= total) return 'err'
  return total - used <= d.trafficDiff * GB ? 'warn' : 'ok'
}

/** A date against now: none, far, close by the panel's threshold, or past. */
export function expiryState(at: number, d: Pick<InboundDefaults, 'expireDiff'>, now = Date.now()): State {
  if (!at) return 'unl'
  if (at <= now) return 'err'
  return at - now <= d.expireDiff * DAY ? 'warn' : 'ok'
}

/** Network and security beside the protocol, for the protocols that have them. */
export function protocolChips(ib: DbInbound): { protocol: string; network?: string; security?: string } {
  if (!['vmess', 'vless', 'trojan', 'shadowsocks'].includes(ib.protocol)) return { protocol: ib.protocol }
  const st = streamOf(ib)
  const security = st.security && st.security !== 'none' ? String(st.security) : undefined
  return { protocol: ib.protocol, network: st.network ? String(st.network) : undefined, security }
}

/** Rows per page to choose from, as the old panel offered them: the panel's size and its multiples up to 1000, and all. */
export function pageSizes(size: number, total: number): number[] {
  const out = [size, ...[2, 5, 10, 20].map((m) => size * m).filter((v) => v < total && v <= 1000)]
  return out.includes(total) ? out : [...out, total]
}

/** The client states the list can be narrowed to; none shows everything. */
export type Filter = 'none' | 'active' | 'disabled' | 'depleted' | 'expiring' | 'online'

/** The emails of an inbound's clients in a state, from its counts. */
export function filterEmails(c: ClientCounts, f: Exclude<Filter, 'none'>): string[] {
  switch (f) {
    case 'active':
      return c.active.filter((e) => !c.depleted.includes(e))
    case 'disabled':
      return c.deactive
    case 'depleted':
      return c.depleted
    case 'expiring':
      return c.expiring
    case 'online':
      return c.online
  }
}
