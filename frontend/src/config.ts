// Settings the panel puts into index.html (see web/spa). The Vite dev server
// serves index.html as is, so development falls back to the defaults.

export interface AppConfig {
  /** Panel base path, e.g. "/" or "/secret/". */
  base: string
  /** Where this frontend is served, e.g. "/secret/v3/". */
  mount: string
  loggedIn: boolean
  version: string
}

function read(): AppConfig {
  const fallback: AppConfig = { base: '/', mount: '/', loggedIn: true, version: 'dev' }
  const raw = document.querySelector<HTMLMetaElement>('meta[name="app-config"]')?.content
  if (!raw) return fallback
  try {
    return { ...fallback, ...(JSON.parse(raw) as Partial<AppConfig>) }
  } catch {
    return fallback
  }
}

export const config = read()

/** URL of a panel path outside this frontend: the API, old pages, logout. */
export function panelUrl(path: string) {
  return config.base + path.replace(/^\//, '')
}
