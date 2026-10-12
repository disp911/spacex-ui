// Numbers, sizes, speeds and durations in the interface language:
// «48,2 МБ», «1,24 МБ/с», «3 д 4 ч» in Russian, "48.2 MB", "3d 4h" in English.
// Intl formats the numbers; unit labels come from the translations because
// Intl's Russian labels differ from ours («кБ», «3 д. 4 ч»).

import { i18n } from '.'

const t = (key: string, params: Record<string, string>) => i18n.global.t(key, params)
const locale = () => i18n.global.locale.value

const byteUnits = ['units.byte', 'units.kilobyte', 'units.megabyte', 'units.gigabyte', 'units.terabyte', 'units.petabyte'] as const

export function formatNumber(value: number, fractionDigits = 0) {
  return new Intl.NumberFormat(locale(), { maximumFractionDigits: fractionDigits }).format(value)
}

/** Number and unit apart, for layouts that style them differently: { n: '10,36', unit: 'ГБ' }. */
export function bytesParts(bytes: number) {
  let value = Math.max(0, bytes)
  let unit = 0
  while (value >= 1024 && unit < byteUnits.length - 1) {
    value /= 1024
    unit++
  }
  // Binary steps; two decimals below 10, one above.
  const digits = unit === 0 ? 0 : value < 10 ? 2 : 1
  return { n: formatNumber(value, digits), unit: t(byteUnits[unit]!, {}) }
}

/** 1536 → «1,5 КБ». */
export function formatBytes(bytes: number) {
  return t('units.value', bytesParts(bytes))
}

export function formatSpeed(bytesPerSecond: number) {
  return t('units.perSecond', { value: formatBytes(bytesPerSecond) })
}

export function formatGigahertz(mhz: number) {
  return t('units.value', { n: formatNumber(mhz / 1000, 2), unit: t('units.gigahertz', {}) })
}

export function formatMilliseconds(ms: number) {
  return t('units.milliseconds', { n: formatNumber(Math.round(ms)) })
}

const durationUnits = [
  ['units.days', 86400],
  ['units.hours', 3600],
  ['units.minutes', 60],
  ['units.seconds', 1],
] as const

/**
 * A duration in its largest unit and the next ones, zeros left out:
 * 273600 s → «3 д 4 ч», 7200 s → «2 ч», 86460 s → «1 д».
 */
export function formatDuration(seconds: number, parts = 2) {
  let rest = Math.max(0, Math.floor(seconds))
  const out: string[] = []
  let taken = 0
  for (const [key, size] of durationUnits) {
    const n = Math.floor(rest / size)
    rest -= n * size
    if (n === 0 && taken === 0) continue
    if (n > 0) out.push(t(key, { n: formatNumber(n) }))
    if (++taken === parts) break
  }
  return out.length ? out.join(' ') : t('units.seconds', { n: formatNumber(0) })
}
