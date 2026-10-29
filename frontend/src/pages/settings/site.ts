// The website the panel serves to visitors without its secret path, and its
// certificate, as POST panel/setting/site returns them.

import { api } from '../../api'

/** reality: Xray holds port 443 and passes visitors without a key to the site; '' : the site is on the panel's own port. */
export type SiteMode = '' | 'reality'

/** What the website tab sets; POST panel/setting/siteApply takes it. The mode is set on the server only. */
export interface SiteSettings {
  domain: string
  /** On the panel's own port: the admin's certificate; empty, the panel obtains one for the domain. */
  certFile: string
  keyFile: string
}

export interface SiteCert {
  /** The panel obtains and renews it itself; otherwise it is read from the admin's files. */
  auto: boolean
  domain: string
  /** Who issued the certificate served; empty while there is none and a self-signed one stands in. */
  issuer: string
  notAfter: string
  /** The last failure to obtain or renew it, cleared by success. */
  error: string
  errorAt: string
}

export interface SiteCheck {
  ok: boolean
  at: string
  error: string
}

export interface SiteInfo extends SiteSettings {
  mode: SiteMode
  /** The panel's own port. */
  port: number
  /** Where the site, the panel and the subscriptions are found, e.g. https://example.com; empty while there is no site. */
  url: string
  /** Behind REALITY: the socket the site listens on, a REALITY target. */
  socket: string
  /** Null while the site has no certificate. */
  cert: SiteCert | null
  /** Null before the first look at the site from outside. */
  check: SiteCheck | null
}

export const loadSite = () => api.post<SiteInfo>('panel/setting/site')

export interface SiteTrouble {
  /** cert: no valid certificate; check: not reachable from outside; renew: renewal failed while the certificate still holds. */
  kind: 'cert' | 'check' | 'renew'
  tone: 'err' | 'warn'
  text: string
}

/** What on the website needs the admin, the worst first; null when nothing. */
export function siteTrouble(s: SiteInfo | null): SiteTrouble | null {
  const c = s?.cert
  if (!c?.auto) return null
  const valid = !!c.issuer && Date.parse(c.notAfter) > Date.now()
  // While the first certificate is on its way there is no failure yet.
  if (!valid && (c.error || c.issuer)) return { kind: 'cert', tone: 'err', text: c.error }
  if (s?.check && !s.check.ok) return { kind: 'check', tone: 'err', text: s.check.error }
  if (c.error) return { kind: 'renew', tone: 'warn', text: c.error }
  return null
}
