// What the subscription server puts into sub.html (sub.PageData in Go):
// amounts in bytes, times in milliseconds, 0 when there is none.

export interface SubData {
  subId: string
  enabled: boolean
  up: number
  down: number
  /** The traffic limit; 0 is none. */
  total: number
  expiry: number
  lastOnline: number
  subUrl: string
  /** Empty while the JSON or Clash subscription is off. */
  subJsonUrl: string
  subClashUrl: string
  links: string[]
}

/** The page's data; the Vite dev server has none, so development shows a sample (?state=unlimited|inactive|delayed). */
export async function loadData(): Promise<SubData | null> {
  const raw = document.querySelector<HTMLMetaElement>('meta[name="app-config"]')?.content
  if (raw) {
    try {
      const data = JSON.parse(raw) as SubData
      return { ...data, links: data.links ?? [] }
    } catch {
      return null
    }
  }
  if (import.meta.env.DEV) return (await import('./sample')).sample(new URLSearchParams(location.search).get('state'))
  return null
}
