// The Xray config template as POST panel/xray/ returns it and panel/xray/update
// takes it back, with the URL that outbound tests open.

import type { MessageSchema } from '../../i18n/schema.gen'

/** A message of the Xray settings screen. */
export type Text = `xray.${keyof MessageSchema['xray']}`

// The template is Xray's own JSON: the screen reads and writes the parts it
// knows and keeps everything else as it is.
type Json = Record<string, any>
export type XrayConfig = Json

export interface XraySetting {
  config: XrayConfig
  testUrl: string
}

/** What POST panel/xray/ answers, as a JSON string. */
export interface XrayReply {
  xraySetting: XrayConfig
  inboundTags: string[]
  clientReverseTags: string[]
  outboundTestUrl: string
}

export const defaultTestUrl = 'https://www.google.com/generate_204'

export type TabKey = 'general' | 'routing' | 'outbounds' | 'balancer' | 'dns' | 'advanced'

export const tabs: { key: TabKey; label: Text }[] = [
  { key: 'general', label: 'xray.tabGeneral' },
  { key: 'routing', label: 'xray.tabRouting' },
  { key: 'outbounds', label: 'xray.tabOutbounds' },
  { key: 'balancer', label: 'xray.tabBalancer' },
  { key: 'dns', label: 'xray.tabDns' },
  { key: 'advanced', label: 'xray.tabAdvanced' },
]

export const options = {
  freedom: ['AsIs', 'UseIP', 'UseIPv4', 'UseIPv6', 'UseIPv6v4', 'UseIPv4v6', 'ForceIP', 'ForceIPv6v4', 'ForceIPv6', 'ForceIPv4v6', 'ForceIPv4'],
  routing: ['AsIs', 'IPIfNonMatch', 'IPOnDemand'],
  logLevel: ['none', 'debug', 'info', 'warning', 'error'],
  access: ['', 'none', './access.log'],
  error: ['', 'none', './error.log'],
  mask: ['', 'quarter', 'half', 'full'],
}

// --- Fields of the General tab ------------------------------------------------

/** A value inside the template with the default it has when the template leaves it out. */
interface Field<T> {
  get(c: XrayConfig): T
  set(c: XrayConfig, v: T): void
}

/** The object at key, made when it is missing. */
const part = (o: Json, key: string): Json => (o[key] ??= {})

/** The freedom outbound tagged direct, which carries the freedom strategy. */
const direct = (c: XrayConfig): Json | undefined =>
  (c.outbounds as Json[] | undefined)?.find((o) => o.protocol === 'freedom' && o.tag === 'direct')

const logField = <T>(key: string, fallback: T): Field<T> => ({
  get: (c) => c.log?.[key] || fallback,
  set: (c, v) => (part(c, 'log')[key] = v),
})

const statsField = (key: string): Field<boolean> => ({
  get: (c) => !!c.policy?.system?.[key],
  set: (c, v) => (part(part(c, 'policy'), 'system')[key] = v),
})

export const fields = {
  freedom: {
    get: (c) => direct(c)?.settings?.domainStrategy || 'AsIs',
    set: (c, v) => {
      const o = direct(c)
      if (o) part(o, 'settings').domainStrategy = v
      else ((c.outbounds ??= []) as Json[]).push({ protocol: 'freedom', tag: 'direct', settings: { domainStrategy: v } })
    },
  } as Field<string>,
  routing: {
    get: (c) => c.routing?.domainStrategy || 'AsIs',
    set: (c, v) => (part(c, 'routing').domainStrategy = v),
  } as Field<string>,
  statsInboundUplink: statsField('statsInboundUplink'),
  statsInboundDownlink: statsField('statsInboundDownlink'),
  statsOutboundUplink: statsField('statsOutboundUplink'),
  statsOutboundDownlink: statsField('statsOutboundDownlink'),
  logLevel: logField('loglevel', 'warning'),
  access: logField('access', ''),
  error: logField('error', ''),
  mask: logField('maskAddress', ''),
  dnsLog: logField('dnsLog', false),
}

export type FieldKey = keyof typeof fields | 'testUrl'

export function fieldValue(s: XraySetting, key: FieldKey): unknown {
  return key === 'testUrl' ? s.testUrl : fields[key].get(s.config)
}

// --- Changes ----------------------------------------------------------------

const same = (a: unknown, b: unknown) => JSON.stringify(a ?? null) === JSON.stringify(b ?? null)

/** Entries of two lists that differ, position by position. */
function listChanges(a: unknown, b: unknown) {
  const x = Array.isArray(a) ? a : []
  const y = Array.isArray(b) ? b : []
  let n = 0
  for (let i = 0; i < Math.max(x.length, y.length); i++) if (!same(x[i], y[i])) n++
  return n
}

/** The outbounds without the freedom strategy, which the General tab counts. */
function outbounds(c: XrayConfig) {
  return ((c.outbounds ?? []) as Json[]).map((o) =>
    o.protocol === 'freedom' && o.tag === 'direct' && o.settings ? { ...o, settings: { ...o.settings, domainStrategy: undefined } } : o,
  )
}

/** Without empty objects, which the fields of the General tab may leave behind. */
function prune(v: unknown): unknown {
  if (!v || typeof v !== 'object' || Array.isArray(v)) return v
  const out = Object.fromEntries(Object.entries(v).map(([k, x]) => [k, prune(x)]).filter(([, x]) => x !== undefined))
  return Object.keys(out).length ? out : undefined
}

/** The parts of the template no other tab shows. */
function rest(c: XrayConfig) {
  const { outbounds: _o, routing, dns: _d, fakedns: _f, observatory: _ob, burstObservatory: _b, log, policy, ...other } = c
  const { domainStrategy: _s, rules: _r, balancers: _bl, ...routingRest } = routing ?? {}
  const { loglevel: _l, access: _a, error: _e, maskAddress: _m, dnsLog: _dl, ...logRest } = log ?? {}
  const { statsInboundUplink: _1, statsInboundDownlink: _2, statsOutboundUplink: _3, statsOutboundDownlink: _4, ...systemRest } = policy?.system ?? {}
  return prune({ other, routingRest, logRest, policy: { ...policy, system: systemRest } })
}

/** How many things a tab has changed. */
export function changes(tab: TabKey, draft: XraySetting, saved: XraySetting): number {
  const a = draft.config
  const b = saved.config
  switch (tab) {
    case 'general':
      return (['testUrl', ...Object.keys(fields)] as FieldKey[]).filter((k) => !same(fieldValue(draft, k), fieldValue(saved, k))).length
    case 'routing':
      return listChanges(a.routing?.rules, b.routing?.rules)
    case 'outbounds':
      return listChanges(outbounds(a), outbounds(b))
    case 'balancer':
      return listChanges(a.routing?.balancers, b.routing?.balancers) + +!same(a.observatory, b.observatory) + +!same(a.burstObservatory, b.burstObservatory)
    case 'dns':
      return +!same(a.dns, b.dns) + +!same(a.fakedns, b.fakedns)
    case 'advanced':
      return +!same(rest(a), rest(b))
  }
}

/** Whether the panel is opened over plain HTTP. */
export const insecure = () => location.protocol !== 'https:'

// --- Routing rules ----------------------------------------------------------------

export type Rule = Json

/** A rule's field as a list: Xray takes some as arrays and some as comma-separated text. */
export const listOf = (v: unknown): string[] =>
  Array.isArray(v) ? v.map(String) : typeof v === 'string' ? v.split(',').map((s) => s.trim()).filter(Boolean) : []

/** What a rule matches, grouped as the table shows it: a label and the values of each part. */
export const ruleGroups: { key: string; label: Text; parts: [string, (r: Rule) => string[]][] }[] = [
  {
    key: 'source',
    label: 'xray.ruleSource',
    parts: [
      ['IP', (r) => listOf(r.sourceIP ?? r.source)],
      ['Port', (r) => listOf(r.sourcePort)],
      ['VLESS', (r) => listOf(r.vlessRoute)],
    ],
  },
  {
    key: 'network',
    label: 'xray.ruleNetwork',
    parts: [
      ['L4', (r) => (r.network ? [String(r.network)] : [])],
      ['Protocol', (r) => listOf(r.protocol)],
      ['Attrs', (r) => (r.attrs && typeof r.attrs === 'object' ? Object.entries(r.attrs).map(([k, v]) => `${k}=${v}`) : listOf(r.attrs))],
    ],
  },
  {
    key: 'destination',
    label: 'xray.ruleDestination',
    parts: [
      ['IP', (r) => listOf(r.ip)],
      ['Domain', (r) => listOf(r.domain)],
      ['Port', (r) => listOf(r.port)],
    ],
  },
  {
    key: 'inbound',
    label: 'xray.ruleInbound',
    parts: [
      ['Tag', (r) => listOf(r.inboundTag)],
      ['User', (r) => listOf(r.user)],
    ],
  },
]

/** Where a rule sends what it matches: an outbound, or a balancer that picks one. */
export function ruleTarget(r: Rule): { kind: 'outbound' | 'balancer'; tag: string } | null {
  if (r.outboundTag) return { kind: 'outbound', tag: String(r.outboundTag) }
  if (r.balancerTag) return { kind: 'balancer', tag: String(r.balancerTag) }
  return null
}

// --- Geo files ------------------------------------------------------------------

/** A geo file of the user's own: where it comes from and when it last came. */
export interface GeoSource {
  id: number
  type: 'geosite' | 'geoip'
  alias: string
  url: string
  lastUpdatedAt?: number
}

// --- Outbounds ------------------------------------------------------------------

export type Outbound = Json

/** Where an outbound connects: host:port for each server it knows, none for freedom and blackhole. */
export function outboundAddresses(o: Outbound): string[] {
  const s = o.settings ?? {}
  const hostPort = (x: Json) => (x?.address ? `${x.address}${x.port != null ? `:${x.port}` : ''}` : '')
  let list: string[] = []
  if (Array.isArray(s.vnext)) list = s.vnext.map(hostPort)
  else if (Array.isArray(s.servers)) list = s.servers.map(hostPort)
  else if (Array.isArray(s.peers)) list = s.peers.map((p: Json) => String(p?.endpoint ?? ''))
  else if (o.protocol !== 'freedom' && o.protocol !== 'blackhole') list = [hostPort(s)]
  return list.filter(Boolean)
}

/** Protocols whose transport and security the table shows beside the protocol. */
const streamed = ['vmess', 'vless', 'trojan', 'shadowsocks']

/** The protocol, and for proxies the transport and, when there is one, tls or reality. */
export function outboundChips(o: Outbound): { protocol: string; network?: string; security?: string } {
  const st = o.streamSettings ?? {}
  if (!streamed.includes(o.protocol)) return { protocol: String(o.protocol ?? '') }
  const security = st.security === 'tls' || st.security === 'reality' ? st.security : undefined
  return { protocol: o.protocol, network: st.network ?? 'tcp', security }
}

/** A blackhole drops everything: nothing to test. */
export const untestable = (o: Outbound) => o.protocol === 'blackhole' || o.tag === 'blocked'

/** What GET panel/xray/getOutboundsTraffic lists. */
export interface OutboundTraffic {
  tag: string
  up: number
  down: number
}

/** What POST panel/xray/testOutbound answers. */
export interface OutboundTest {
  success: boolean
  delay: number
  error?: string
  statusCode?: number
}

// --- Balancers ------------------------------------------------------------------

export type Balancer = Json

/** The strategies a balancer picks an outbound by; random when the template names none. */
export const strategies = [
  { value: 'random', label: 'Random' },
  { value: 'roundRobin', label: 'Round Robin' },
  { value: 'leastLoad', label: 'Least Load' },
  { value: 'leastPing', label: 'Least Ping' },
]

export const strategyOf = (b: Balancer): string => b.strategy?.type ?? 'random'

const defaultObservatory = (): Json => ({
  subjectSelector: [],
  probeURL: 'https://www.google.com/generate_204',
  probeInterval: '1m',
  enableConcurrency: true,
})

const defaultBurstObservatory = (): Json => ({
  subjectSelector: [],
  pingConfig: {
    destination: 'https://www.google.com/generate_204',
    interval: '1m',
    connectivity: 'http://connectivitycheck.platform.hicloud.com/generate_204',
    timeout: '5s',
    sampling: 2,
  },
})

/**
 * Keeps the observers in step with the balancers, as the old panel did: Least
 * Ping balancers have the Observatory watch their selectors, Least Load, Round
 * Robin and Random ones the Burst Observatory; an observer no balancer needs goes away. Settings made
 * by hand stay, the selectors are counted anew.
 */
export function syncObservatories(c: XrayConfig) {
  const balancers = (c.routing?.balancers ?? []) as Balancer[]
  const watch = (key: 'observatory' | 'burstObservatory', list: Balancer[], made: () => Json) => {
    if (!list.length) {
      delete c[key]
      return
    }
    c[key] ??= made()
    c[key].subjectSelector = [...new Set(list.flatMap((b) => listOf(b.selector)))]
  }
  watch('observatory', balancers.filter((b) => strategyOf(b) === 'leastPing'), defaultObservatory)
  watch('burstObservatory', balancers.filter((b) => ['leastLoad', 'roundRobin', 'random'].includes(strategyOf(b))), defaultBurstObservatory)
}

// --- DNS ------------------------------------------------------------------------

/** A DNS server: an address alone, or an object with its own settings. */
export type DnsServer = string | Json

export const dnsStrategies = ['UseSystem', 'UseIP', 'UseIPv4', 'UseIPv6']

/** The DNS block the old panel writes when DNS is switched on. */
export const newDns = (): Json => ({ servers: [], queryStrategy: 'UseIP', tag: 'dns_inbound', enableParallelQuery: false })

export const defaultDnsServer = (): Json => ({
  address: 'localhost',
  port: 53,
  domains: [],
  expectIPs: [],
  unexpectedIPs: [],
  queryStrategy: 'UseIP',
  skipFallback: true,
  disableCache: false,
  finalQuery: false,
})

export const defaultFakeDns = (): Json => ({ ipPool: '198.18.0.0/16', poolSize: 65535 })

/** The DNS server lists of the old panel's templates; installing one replaces the servers. */
export const dnsPresets: { name: string; family: boolean; servers: string[] }[] = [
  { name: 'Google DNS', family: false, servers: ['8.8.8.8', '8.8.4.4', '2001:4860:4860::8888', '2001:4860:4860::8844'] },
  { name: 'Cloudflare DNS', family: false, servers: ['1.1.1.1', '1.0.0.1', '2606:4700:4700::1111', '2606:4700:4700::1001'] },
  { name: 'Adguard DNS', family: false, servers: ['94.140.14.14', '94.140.15.15', '2a10:50c0::ad1:ff', '2a10:50c0::ad2:ff'] },
  { name: 'Adguard Family DNS', family: true, servers: ['94.140.14.14', '94.140.15.15', '2a10:50c0::ad1:ff', '2a10:50c0::ad2:ff'] },
  { name: 'Cloudflare Family DNS', family: true, servers: ['1.1.1.3', '1.0.0.3', '2606:4700:4700::1113', '2606:4700:4700::1003'] },
]

export const dnsAddress = (s: DnsServer) => (typeof s === 'string' ? s : String(s?.address ?? ''))
export const dnsList = (s: DnsServer, key: 'domains' | 'expectIPs') => (typeof s === 'string' ? [] : listOf(s?.[key]))
