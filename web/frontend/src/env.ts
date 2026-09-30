// Runtime configuration injected by the Go panel into index.html
// (see web/spa.go). Falls back to sane defaults for `npm run dev`.
interface Boot {
  basePath: string
  host: string
  version: string
}

declare global {
  interface Window {
    __SPX__?: Partial<Boot>
  }
}

const boot = window.__SPX__ ?? {}

function normalize(path: string | undefined): string {
  if (!path) return '/'
  let p = path.startsWith('/') ? path : '/' + path
  if (!p.endsWith('/')) p += '/'
  return p
}

export const env: Boot = {
  basePath: normalize(boot.basePath),
  host: boot.host || window.location.hostname,
  version: boot.version || '',
}

/** URL of the SPA root, e.g. "/secret/ui/". */
export const appBase = env.basePath + 'ui/'

/** Build an absolute URL under the panel base path. */
export function panelUrl(path: string): string {
  return env.basePath + path.replace(/^\//, '')
}
