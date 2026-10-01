const UNITS = ['B', 'KB', 'MB', 'GB', 'TB', 'PB']

/** Split a byte count into a display value and unit, e.g. [10.36, 'GB']. */
export function sizeParts(bytes: number, digits = 2): [string, string] {
  if (!Number.isFinite(bytes) || bytes <= 0) return ['0', 'B']
  let i = 0
  let v = bytes
  while (v >= 1024 && i < UNITS.length - 1) {
    v /= 1024
    i++
  }
  return [v.toFixed(i === 0 ? 0 : digits), UNITS[i]]
}

export function formatSize(bytes: number, digits = 2): string {
  const [v, u] = sizeParts(bytes, digits)
  return `${v} ${u}`
}

export function formatSpeed(bytesPerSec: number): string {
  return formatSize(bytesPerSec) + '/s'
}

/** Compact uptime such as "3d 4h", "5h 12m" or "42s". */
export function formatUptime(seconds: number): string {
  if (!Number.isFinite(seconds) || seconds < 0) return '—'
  const d = Math.floor(seconds / 86400)
  const h = Math.floor((seconds % 86400) / 3600)
  const m = Math.floor((seconds % 3600) / 60)
  if (d > 0) return `${d}d ${h}h`
  if (h > 0) return `${h}h ${m}m`
  if (m > 0) return `${m}m`
  return `${Math.floor(seconds)}s`
}

export function percent(current: number, total: number): number {
  if (!total) return 0
  return Math.min(100, Math.max(0, (current / total) * 100))
}

/** Thousands separated with a thin space, e.g. 6 200. */
export function formatCount(n: number): string {
  return String(Math.round(n)).replace(/\B(?=(\d{3})+(?!\d))/g, ' ')
}

const pad = (n: number) => String(n).padStart(2, '0')

/** dd.MM.yyyy HH:mm */
export function formatDateTime(d: Date): string {
  return `${pad(d.getDate())}.${pad(d.getMonth() + 1)}.${d.getFullYear()} ${pad(d.getHours())}:${pad(d.getMinutes())}`
}

/** dd.MM HH:mm:ss — used in dense log tables. */
export function formatLogTime(iso: string): string {
  const d = new Date(iso)
  if (Number.isNaN(d.getTime())) return iso
  return `${pad(d.getDate())}.${pad(d.getMonth() + 1)} ${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())}`
}

export function formatClock(iso: string): string {
  const d = new Date(iso)
  if (Number.isNaN(d.getTime())) return iso
  return `${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())}`
}

/** "2026-09-20" → "20.09.2026" (the backend keys log days in ISO form). */
export function formatDay(day: string): string {
  const m = /^(\d{4})-(\d{2})-(\d{2})$/.exec(day)
  return m ? `${m[3]}.${m[2]}.${m[1]}` : day
}

/** "26.4.25" → "v26.4.25"; non-numeric values (e.g. "Unknown") give "". */
export function withV(version: string): string {
  if (!version || !/^v?\d/i.test(version)) return ''
  return /^v/i.test(version) ? version : 'v' + version
}
