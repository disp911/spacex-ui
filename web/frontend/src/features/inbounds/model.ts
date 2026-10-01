import type { ClientStat, DbInbound, PanelDefaults } from '@/api/panel'

export const PROTOCOL_LABELS: Record<string, string> = {
  vless: 'VLESS',
  vmess: 'VMess',
  trojan: 'Trojan',
  shadowsocks: 'Shadowsocks',
  hysteria: 'Hysteria',
  wireguard: 'WireGuard',
  http: 'HTTP',
  mixed: 'Mixed',
  mtproto: 'MTProto',
  tunnel: 'Tunnel',
  tun: 'TUN',
}

/** Protocols whose settings carry a client list. */
const MULTI_USER = new Set(['vless', 'vmess', 'trojan', 'shadowsocks', 'hysteria', 'mtproto'])

/** Usage health: ok, close to a limit, or over it / expired / disabled. */
export type Health = 'ok' | 'warn' | 'depleted' | 'off'

export interface ClientRow {
  email: string
  raw: Record<string, unknown>
  /** Id the updateClient endpoint expects for this protocol. */
  clientId: string
  enable: boolean
  subId: string
  limitIp: number
  comment: string
  up: number
  down: number
  /** Traffic limit in bytes, 0 = unlimited. */
  total: number
  /** Unix ms; 0 = never, negative = starts on first use (days). */
  expiryTime: number
  online: boolean
  lastOnline: number
  health: Health
}

export interface InboundRow {
  db: DbInbound
  id: number
  remark: string
  protocol: string
  protocolLabel: string
  network: string
  security: string
  port: number
  enable: boolean
  up: number
  down: number
  total: number
  allTime: number
  expiryTime: number
  multiUser: boolean
  clients: ClientRow[]
  online: number
  health: Health
  /** Clients that are over quota or expired. */
  depleted: number
  /** Clients that are close to their quota or expiry. */
  expiring: number
}

function parse<T>(s: string, fallback: T): T {
  try {
    return (s ? JSON.parse(s) : fallback) ?? fallback
  } catch {
    return fallback
  }
}

function clientIdFor(protocol: string, c: Record<string, any>): string {
  switch (protocol) {
    case 'trojan':
      return c.password ?? ''
    case 'shadowsocks':
      return c.email ?? ''
    case 'hysteria':
      return c.auth ?? ''
    default:
      return c.id ?? ''
  }
}

export function healthOf(
  x: { enable: boolean; up: number; down: number; total: number; expiryTime: number },
  d: Pick<PanelDefaults, 'expireDiff' | 'trafficDiff'> | null,
  now = Date.now(),
): Health {
  const used = x.up + x.down
  // Checked before `enable`: the panel disables depleted clients itself,
  // and those should still read as depleted, not as switched off.
  if (x.total > 0 && used >= x.total) return 'depleted'
  if (x.expiryTime > 0 && x.expiryTime <= now) return 'depleted'
  if (!x.enable) return 'off'
  // Panel settings define "soon"; fall back to 3 days / 90 % when unset.
  const expireDiff = (d?.expireDiff || 3) * 86_400_000
  const trafficLeft = d?.trafficDiff ? d.trafficDiff * 1_073_741_824 : x.total * 0.1
  if (x.total > 0 && x.total - used <= trafficLeft) return 'warn'
  if (x.expiryTime > 0 && x.expiryTime - now <= expireDiff) return 'warn'
  return 'ok'
}

export function toRow(db: DbInbound, online: Set<string>, d: PanelDefaults | null): InboundRow {
  const settings = parse<Record<string, any>>(db.settings, {})
  const stream = parse<Record<string, any>>(db.streamSettings, {})
  const stats = new Map<string, ClientStat>((db.clientStats ?? []).map((s) => [s.email, s]))
  const multiUser = MULTI_USER.has(db.protocol) && Array.isArray(settings.clients)

  const clients: ClientRow[] = multiUser
    ? (settings.clients as Record<string, any>[]).map((c) => {
        const st = stats.get(c.email)
        const row = {
          email: c.email ?? '',
          raw: c,
          clientId: clientIdFor(db.protocol, c),
          enable: c.enable !== false && (st?.enable ?? true),
          subId: c.subId ?? '',
          limitIp: c.limitIp ?? 0,
          comment: c.comment ?? '',
          up: st?.up ?? 0,
          down: st?.down ?? 0,
          total: st?.total ?? (c.totalGB ?? 0),
          expiryTime: st?.expiryTime ?? (c.expiryTime ?? 0),
          online: online.has(c.email),
          lastOnline: st?.lastOnline ?? 0,
          health: 'ok' as Health,
        }
        row.health = healthOf(row, d)
        return row
      })
    : []

  const row: InboundRow = {
    db,
    id: db.id,
    remark: db.remark,
    protocol: db.protocol,
    protocolLabel: PROTOCOL_LABELS[db.protocol] ?? db.protocol,
    network: stream.network ?? '',
    security: stream.security && stream.security !== 'none' ? stream.security : '',
    port: db.port,
    enable: db.enable,
    up: db.up,
    down: db.down,
    total: db.total,
    allTime: db.allTime,
    expiryTime: db.expiryTime,
    multiUser,
    clients,
    online: clients.filter((c) => c.online).length,
    health: 'ok',
    depleted: clients.filter((c) => c.health === 'depleted').length,
    expiring: clients.filter((c) => c.health === 'warn').length,
  }
  const own = healthOf(row, d)
  row.health = own !== 'ok' ? own : row.depleted ? 'warn' : row.expiring ? 'warn' : 'ok'
  return row
}

/** "до 12.10.2026", "через 5 д.", "∞" — expiry label for a timestamp. */
export function expiryParts(ms: number, now = Date.now()): { date: string; days: number | null } {
  if (!ms) return { date: '', days: null }
  if (ms < 0) return { date: '', days: Math.round(-ms / 86_400_000) } // starts on first use
  const d = new Date(ms)
  const pad = (n: number) => String(n).padStart(2, '0')
  return { date: `${pad(d.getDate())}.${pad(d.getMonth() + 1)}.${d.getFullYear()}`, days: Math.ceil((ms - now) / 86_400_000) }
}
