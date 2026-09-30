import { ref, watch } from 'vue'

export type Theme = 'dark' | 'light'

const KEY = 'spx-theme'
const initial = (document.documentElement.dataset.theme as Theme) || 'dark'
const theme = ref<Theme>(initial === 'light' ? 'light' : 'dark')

watch(theme, (t) => {
  document.documentElement.dataset.theme = t
  try {
    localStorage.setItem(KEY, t)
  } catch {
    /* storage unavailable */
  }
})

export function useTheme() {
  return {
    theme,
    toggle: () => (theme.value = theme.value === 'dark' ? 'light' : 'dark'),
  }
}
