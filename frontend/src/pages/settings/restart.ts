// After a restart the panel may answer at another address: the saved port,
// base path and certificate take effect. The address is worked out as the old
// settings page did it.

import { config } from '../../config'
import type { AllSetting } from './model'

function isIp(host: string) {
  const h = host.replace(/^\[|\]$/g, '')
  return /^(\d{1,3}\.){3}\d{1,3}$/.test(h) || (h.includes(':') && /^[\da-f:.]+$/i.test(h))
}

/**
 * Where the panel will be: its root and this page there. Through a domain only
 * the scheme and the path change (a proxy may stand in front); through an IP
 * also the port, and the host when the panel's domain is an IP. The dev server
 * keeps its own address.
 */
export function addressAfterRestart(s: AllSetting) {
  const url = new URL(location.href)
  let base = s.webBasePath.replace(/^\//, '')
  if (base && !base.endsWith('/')) base += '/'
  url.pathname = `/${base}`
  url.search = ''
  url.hash = ''
  if (!import.meta.env.DEV) {
    url.protocol = s.webCertFile || s.webKeyFile ? 'https:' : 'http:'
    if (isIp(location.hostname)) {
      if (s.webDomain && isIp(s.webDomain)) url.hostname = s.webDomain
      if (s.webPort && String(s.webPort) !== location.port) url.port = String(s.webPort)
    }
  }
  const panel = url.toString()
  return { panel, page: new URL(`${config.mount.slice(config.base.length)}settings`, panel).toString() }
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
