const UNITS = ['B', 'KB', 'MB', 'GB', 'TB', 'PB']

/** Splits a byte count into a value and a unit: 1536 → ["1.50", "KB"]. */
export function sizeParts(bytes: number): [string, string] {
  if (!bytes || bytes <= 0) return ['0', 'B']
  let i = 0
  let v = bytes
  while (v >= 1024 && i < UNITS.length - 1) {
    v /= 1024
    i++
  }
  return [i === 0 ? v.toFixed(0) : v.toFixed(2), UNITS[i]]
}

export function size(bytes: number): string {
  return sizeParts(bytes).join(' ')
}

/** Uptime in seconds as "3d 4h", "5h", "12m" or "40s". */
export function duration(sec: number): string {
  if (!sec || sec < 0) return '0s'
  if (sec < 60) return `${Math.floor(sec)}s`
  if (sec < 3600) return `${Math.floor(sec / 60)}m`
  if (sec < 86400) return `${Math.floor(sec / 3600)}h`
  const d = Math.floor(sec / 86400)
  const h = Math.floor((sec % 86400) / 3600)
  return h ? `${d}d ${h}h` : `${d}d`
}

export function percent(current: number, total: number): number {
  if (!total) return 0
  return Math.min(100, Math.max(0, (current / total) * 100))
}

/** 14.2 → "14.2", 1.3 → "1.30": three significant figures like the design. */
export function pct(v: number): string {
  return v >= 10 ? v.toFixed(1) : v.toFixed(2)
}

export function cpuSpeed(mhz: number): string {
  return mhz > 1000 ? `${(mhz / 1000).toFixed(2)} GHz` : `${mhz.toFixed(0)} MHz`
}

const pad = (n: number) => String(n).padStart(2, '0')

/** Unix seconds → "18.09.2026 04:12". */
export function dateTime(ts: number): string {
  if (!ts) return ''
  const d = new Date(ts * 1000)
  return `${pad(d.getDate())}.${pad(d.getMonth() + 1)}.${d.getFullYear()} ${pad(d.getHours())}:${pad(d.getMinutes())}`
}

export function version(v: string): string {
  if (!v) return ''
  return /^\d/.test(v) ? `v${v}` : v
}
