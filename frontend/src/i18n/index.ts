import { createI18n } from 'vue-i18n'
import type { MessageSchema } from './schema.gen'

export type Locale = 'en-US' | 'ru-RU'

export const locales: { value: Locale; label: string; flag: string }[] = [
  { value: 'ru-RU', label: 'Русский', flag: '🇷🇺' },
  { value: 'en-US', label: 'English', flag: '🇬🇧' },
]

declare module 'vue-i18n' {
  // eslint-disable-next-line @typescript-eslint/no-empty-object-type
  export interface DefineLocaleMessage extends MessageSchema {}
}

/** Russian has three plural forms: 1 файл, 2 файла, 5 файлов. */
function russianPlural(choice: number, choicesLength: number) {
  if (choicesLength < 3) return choice === 1 ? 0 : 1
  const n = Math.abs(choice) % 100
  const last = n % 10
  if (n > 10 && n < 20) return 2
  if (last === 1) return 0
  if (last >= 2 && last <= 4) return 1
  return 2
}

export const i18n = createI18n<[MessageSchema], Locale, false>({
  legacy: false,
  locale: 'en-US',
  fallbackLocale: 'en-US',
  messages: {} as Record<Locale, MessageSchema>,
  pluralRules: { 'ru-RU': russianPlural },
  missingWarn: import.meta.env.DEV,
  fallbackWarn: false,
})

// Only the active language is downloaded.
const loaders: Record<Locale, () => Promise<{ default: MessageSchema }>> = {
  'en-US': () => import('virtual:locale/en-US'),
  'ru-RU': () => import('virtual:locale/ru-RU'),
}

function readCookie(name: string) {
  return document.cookie.split('; ').find((c) => c.startsWith(`${name}=`))?.slice(name.length + 1)
}

/** The panel keeps the language in the lang cookie, which the server reads too. */
export function initialLocale(): Locale {
  const saved = readCookie('lang')
  if (saved === 'en-US' || saved === 'ru-RU') return saved
  return navigator.language.toLowerCase().startsWith('ru') ? 'ru-RU' : 'en-US'
}

export async function setLocale(locale: Locale) {
  if (!i18n.global.availableLocales.includes(locale)) {
    i18n.global.setLocaleMessage(locale, (await loaders[locale]()).default)
  }
  i18n.global.locale.value = locale
  document.documentElement.lang = locale.slice(0, 2)
  document.cookie = `lang=${locale}; path=/; max-age=31536000; samesite=lax`
}
