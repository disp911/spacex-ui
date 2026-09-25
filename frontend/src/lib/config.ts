// Boot config the Go server writes into index.html (web/controller/spa.go).

export interface BootConfig {
  page: 'login' | 'dashboard'
  basePath: string
  version: string
  host: string
  lang: string
  languages: { code: string; name: string }[]
  twoFactor: boolean
  messages: Record<string, string>
}

const fallback: BootConfig = {
  page: 'login',
  basePath: '/',
  version: '',
  host: location.hostname,
  lang: 'en-US',
  languages: [],
  twoFactor: false,
  messages: {},
}

export const config: BootConfig = { ...fallback, ...(window.__SPX__ ?? {}) }

/** Absolute URL of a panel path such as "panel/inbounds". */
export function panelUrl(path = ''): string {
  return config.basePath + path.replace(/^\//, '')
}
