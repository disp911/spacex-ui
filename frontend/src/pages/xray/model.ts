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
