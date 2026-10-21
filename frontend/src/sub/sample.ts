import type { SubData } from './data'

// The subscriptions mockup frames 1a, 1b and 1c draw, and a delayed start, for the Vite dev server.

const GB = 1024 ** 3
const at = (s: string) => new Date(s).getTime()
const base = 'https://de.example.net:2096'
const link = (host: string, name: string) =>
  `vless://5f8e2a41-9c3d-4b7e-a1f2-6d0c8b3e7a95@${host}:443?type=tcp&security=reality&pbk=Xk3vQ8mN2pL5rT9wY1zA4bC7dE0fG6hJ8kM2nP5qR3s&fp=chrome&sni=www.microsoft.com&sid=6ba85179e30d4fc2&spx=%2F&flow=xtls-rprx-vision#${name}`

export function sample(state: string | null): SubData {
  const common = {
    subId: 'alice-7f3k2',
    enabled: true,
    subUrl: `${base}/sub/alice-7f3k2`,
    subJsonUrl: `${base}/json/alice-7f3k2`,
    subClashUrl: `${base}/clash/alice-7f3k2`,
    links: [link('de.example.net', 'Germany-VLESS-alice'), link('de-cdn.example.net', 'Germany-VLESS-alice-CDN')],
  }
  if (state === 'unlimited') {
    return { ...common, up: 18.6 * GB, down: 412.8 * GB, total: 0, expiry: 0, lastOnline: at('2026-10-10T02:41:37'), subJsonUrl: '', subClashUrl: '' }
  }
  if (state === 'delayed') {
    return { ...common, up: 0, down: 0, total: 100 * GB, expiry: -30 * 86_400_000, lastOnline: 0 }
  }
  if (state === 'inactive') {
    return { ...common, up: 1.1 * GB, down: 48.9 * GB, total: 50 * GB, expiry: at('2026-09-30T00:00:00'), lastOnline: 0 }
  }
  return { ...common, up: 1.2 * GB, down: 30.5 * GB, total: 100 * GB, expiry: at('2026-12-31T23:59:59'), lastOnline: at('2026-10-09T13:58:12') }
}
