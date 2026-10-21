import type { SubData } from './data'

export type Status = 'active' | 'unlimited' | 'inactive'

/**
 * As the old page judged it: off is inactive; on without a traffic limit is
 * unlimited; otherwise active while time and traffic last. A delayed start
 * (a negative expiry, the days after the first connection) has not run out.
 */
export function statusOf(d: SubData, now = Date.now()): Status {
  if (!d.enabled) return 'inactive'
  if (!d.total) return 'unlimited'
  const timeLeft = d.expiry <= 0 || d.expiry >= now
  const trafficLeft = d.up + d.down <= d.total
  return timeLeft && trafficLeft ? 'active' : 'inactive'
}

/** Days after the first connection a delayed start runs for, null for a date or none. */
export const delayedDays = (expiry: number) => (expiry < 0 ? Math.round(expiry / -86_400_000) : null)

/** A link's name as the old page found it: the vmess ps, the #remark, a remark or email parameter, the user. */
export function linkName(link: string): string | null {
  try {
    if (link.startsWith('vmess://')) {
      const json = JSON.parse(atob(link.slice('vmess://'.length))) as { ps?: string; add?: string; id?: string }
      return json.ps || (json.add && json.id ? json.add : null)
    }
    const hash = link.indexOf('#')
    if (hash !== -1) return decodeURIComponent(link.slice(hash + 1))
    if (link.startsWith('vless://') || link.startsWith('trojan://')) {
      const query = link.indexOf('?')
      if (query !== -1) {
        const params = new URLSearchParams(link.slice(query + 1))
        const named = params.get('remark') || params.get('email')
        if (named) return named
      }
      const at = link.indexOf('@')
      const sep = link.indexOf('://')
      if (at !== -1 && sep !== -1) return link.slice(sep + 3, at)
    }
  } catch {
    // a broken link keeps the numbered name
  }
  return null
}

export interface App {
  name: string
  /** Opens the app with the subscription; without it the link is copied for pasting. */
  open?: string
}

/** The apps the old page offered and how each takes a subscription. */
export function apps(d: SubData): { android: App[]; ios: App[] } {
  const url = d.subUrl
  const sub = encodeURIComponent(url)
  const name = encodeURIComponent(d.subId)
  const v2box = { name: 'V2Box', open: `v2box://install-sub?url=${sub}&name=${name}` }
  const happ = { name: 'Happ', open: `happ://add/${url}` }
  return {
    android: [
      v2box,
      { name: 'V2RayNG', open: `v2rayng://install-config?url=${sub}` },
      { name: 'Sing-box' },
      { name: 'V2RayTun' },
      { name: 'NPV Tunnel' },
      happ,
    ],
    ios: [
      { name: 'Shadowrocket', open: `shadowrocket://add/sub/${btoa(`${url}?flag=shadowrocket`)}?remark=${encodeURIComponent(d.subId || 'Subscription')}` },
      v2box,
      { name: 'Streisand', open: `streisand://import/${sub}` },
      { name: 'V2RayTun' },
      { name: 'NPV Tunnel' },
      happ,
    ],
  }
}
