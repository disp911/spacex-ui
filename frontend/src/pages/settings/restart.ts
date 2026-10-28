// After a restart the panel may answer at another address: the saved port and
// base path take effect, or the website's mode and domain.

import { config } from '../../config'
import type { AllSetting } from './model'
import type { SiteSettings } from './site'

function isIp(host: string) {
  const h = host.replace(/^\[|\]$/g, '')
  return /^(\d{1,3}\.){3}\d{1,3}$/.test(h) || (h.includes(':') && /^[\da-f:.]+$/i.test(h))
}

/** The panel's root at url with the base path, and this page there. */
function placesAt(url: URL, basePath: string) {
  let base = basePath.replace(/^\//, '')
  if (base && !base.endsWith('/')) base += '/'
  url.pathname = `/${base}`
  url.search = ''
  url.hash = ''
  const panel = url.toString()
  return { panel, page: new URL(`${config.mount.slice(config.base.length)}settings`, panel).toString() }
}

/**
 * Where the panel will be: its root and this page there. The scheme stays:
 * the certificate is set on the website tab, which moves the panel itself.
 * Through a domain only the path changes (a proxy may stand in front);
 * through an IP also the port, and the host when the panel's domain is an
 * IP. The dev server keeps its own address.
 */
export function addressAfterRestart(s: AllSetting) {
  const url = new URL(location.href)
  if (!import.meta.env.DEV && isIp(location.hostname)) {
    if (s.webDomain && isIp(s.webDomain)) url.hostname = s.webDomain
    if (s.webPort && String(s.webPort) !== location.port) url.port = String(s.webPort)
  }
  return placesAt(url, s.webBasePath)
}

/**
 * Where the panel will be with the website's settings applied. Behind REALITY
 * it answers at the domain on 443, and on its own port only when that stays
 * open and this page came through it. On its own port it answers at the
 * domain or IP set, over https when it has a certificate.
 */
export function addressAfterSite(s: SiteSettings, port: number, basePath: string) {
  const url = new URL(location.href)
  if (!import.meta.env.DEV) {
    if (s.mode === 'reality') {
      if (!(s.panelPort && location.port === String(port))) {
        url.protocol = 'https:'
        url.hostname = s.domain
        url.port = ''
      }
    } else {
      url.protocol = s.domain || s.certFile ? 'https:' : 'http:'
      if (s.domain) url.hostname = s.domain.includes(':') ? `[${s.domain}]` : s.domain
      url.port = String(port)
    }
  }
  return placesAt(url, basePath)
}

const sleep = (ms: number) => new Promise((r) => setTimeout(r, ms))

/** Whether the panel answers: any reply to its login endpoint but a proxy's 5xx counts. */
async function answers(panel: string) {
  const url = new URL('login', panel)
  const sameOrigin = url.origin === location.origin
  try {
    const res = await fetch(url, { mode: sameOrigin ? 'same-origin' : 'no-cors', cache: 'no-store', credentials: 'omit' })
    // Another origin hides the reply; that it came is enough.
    return !sameOrigin || res.status < 500
  } catch {
    return false
  }
}

/**
 * Resolves when the panel answers again. It restarts three seconds after the
 * request, so asking starts after four; then once a second, for as long as it
 * takes. A browser on https cannot ask an http address, so such a move only
 * waits the restart out.
 */
export async function waitForPanel(panel: string) {
  await sleep(4000)
  if (location.protocol === 'https:' && panel.startsWith('http:')) return
  while (!(await answers(panel))) await sleep(1000)
}
