import { config } from './config'

type Params = Record<string, string | number>

/**
 * Translates a key from the panel's TOML translations. Placeholders may be
 * written either as {{ .name }} (go-i18n style) or #name# (legacy 3x-ui).
 */
export function t(key: string, params?: Params): string {
  let s = config.messages[key] ?? key
  if (params) {
    for (const [k, v] of Object.entries(params)) {
      s = s.replace(new RegExp(`\\{\\{\\s*\\.${k}\\s*\\}\\}|#${k}#`, 'g'), String(v))
    }
  }
  return s
}

export function setLanguage(code: string) {
  document.cookie = `lang=${encodeURIComponent(code)}; path=/; max-age=${60 * 60 * 24 * 365}; SameSite=Lax`
  location.reload()
}

const rules = new Intl.PluralRules(config.lang)
const ORDER: Record<string, string[]> = {
  ru: ['one', 'few', 'many'],
  en: ['one', 'other'],
}

/**
 * Plural-aware t(): the message holds its forms separated by "|" in the
 * language's CLDR order (ru: one|few|many, en: one|other); {{ .n }} is the count.
 */
export function tn(key: string, n: number, params?: Params): string {
  const forms = t(key, { n, ...params }).split('|')
  const order = ORDER[config.lang.slice(0, 2)] ?? ['one', 'other']
  const i = order.indexOf(rules.select(n))
  return (forms[i >= 0 ? i : forms.length - 1] ?? forms[forms.length - 1]).trim()
}
