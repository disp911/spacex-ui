import { panelUrl } from '@/env'
import type { DbInbound } from './panel'

// Share links (vless://, vmess://, trojan://, WireGuard configs …) are built
// by the legacy client model in web/assets/js/model. It is large and covers
// every protocol/transport combination, so instead of a second
// implementation we load those scripts on demand and call them. The result
// is byte-for-byte what the legacy panel produces.
const SCRIPTS = [
  'assets/moment/moment.min.js',
  'assets/js/util/index.js',
  'assets/qrcode/qrious2.min.js',
  'assets/uri/URI.min.js',
  'assets/js/model/reality_targets.js',
  'assets/js/model/inbound.js',
  'assets/js/model/dbinbound.js',
]

let loading: Promise<void> | null = null

function loadScript(src: string): Promise<void> {
  return new Promise((resolve, reject) => {
    const s = document.createElement('script')
    s.src = panelUrl(src)
    s.async = false
    s.onload = () => resolve()
    s.onerror = () => reject(new Error('Failed to load ' + src))
    document.head.appendChild(s)
  })
}

function ensureLoaded(): Promise<void> {
  if (!loading) {
    loading = (async () => {
      for (const src of SCRIPTS) await loadScript(src)
    })().catch((e) => {
      loading = null
      throw e
    })
  }
  return loading
}

// The scripts declare top-level classes: global lexical bindings, not
// window properties, so they are reached through a global-scope Function.
function global<T = any>(name: string): T {
  return new Function(`return typeof ${name} === 'undefined' ? undefined : ${name}`)() as T
}

export interface ShareLink {
  remark: string
  link: string
}

/**
 * Links for one client of an inbound (or all peers of a WireGuard inbound).
 * `client` is the raw client object from the inbound settings JSON.
 */
export async function shareLinks(db: DbInbound, client: Record<string, unknown> | null, remarkModel: string): Promise<ShareLink[]> {
  await ensureLoaded()
  const DBInbound = global('DBInbound')
  if (!DBInbound) throw new Error('Link model is unavailable')
  const dbInbound = new DBInbound(db)
  const inbound = dbInbound.toInbound()

  if (db.protocol === 'wireguard') {
    return String(inbound.genInboundLinks(db.remark) || '')
      .split('\r\n')
      .filter(Boolean)
      .map((link, i) => ({ remark: `Peer ${i + 1}`, link }))
  }
  // genAllLinks matches clients by object; hand it the parsed instance.
  const match = client ? (inbound.clients ?? []).find((c: any) => c.email === client.email) ?? client : undefined
  return (inbound.genAllLinks(db.remark, remarkModel || '-ieo', match) as ShareLink[]).filter((l) => l.link)
}

/** Render `text` as a QR code into a canvas (QRious from the panel assets). */
export async function drawQr(canvas: HTMLCanvasElement, text: string, size: number, dark: boolean) {
  await ensureLoaded()
  const QRious = global('QRious')
  if (!QRious) return
  new QRious({
    element: canvas,
    value: text,
    size,
    level: 'L',
    background: dark ? '#111520' : '#ffffff',
    foreground: dark ? '#e6e9ef' : '#13171e',
    padding: 10,
  })
}
