import { computed, ref, watch } from 'vue'
import ru from './ru'
import en from './en'

export type Locale = 'ru' | 'en'

export const LOCALES: { code: Locale; label: string; cookie: string }[] = [
  { code: 'ru', label: 'Русский', cookie: 'ru-RU' },
  { code: 'en', label: 'English', cookie: 'en-US' },
]

const dicts: Record<Locale, unknown> = { ru, en }

// The Go backend localizes its own messages from the `lang` cookie
// (web/locale); keep the two in sync.
function readCookieLocale(): Locale | null {
  const m = /(?:^|;\s*)lang=([^;]+)/.exec(document.cookie)
  if (!m) return null
  const v = decodeURIComponent(m[1]).toLowerCase()
  return v.startsWith('ru') ? 'ru' : v.startsWith('en') ? 'en' : null
}

const initial: Locale = readCookieLocale() ?? (navigator.language?.toLowerCase().startsWith('ru') ? 'ru' : 'en')
export const locale = ref<Locale>(initial)

function applyLocale(l: Locale) {
  document.documentElement.lang = l
  const cookie = LOCALES.find((x) => x.code === l)!.cookie
  document.cookie = `lang=${encodeURIComponent(cookie)};path=/;max-age=31536000`
}
applyLocale(initial)
watch(locale, applyLocale)

function lookup(path: string, l: Locale): string | undefined {
  let cur: any = dicts[l]
  for (const part of path.split('.')) {
    if (cur == null) return undefined
    cur = cur[part]
  }
  return typeof cur === 'string' ? cur : undefined
}

export function t(key: string, params?: Record<string, string | number>): string {
  let s = lookup(key, locale.value) ?? lookup(key, 'ru') ?? key
  if (params) {
    for (const [k, v] of Object.entries(params)) s = s.split(`{${k}}`).join(String(v))
  }
  return s
}

export function useI18n() {
  return { t, locale, localeCode: computed(() => locale.value.toUpperCase()) }
}
