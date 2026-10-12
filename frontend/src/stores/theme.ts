import { defineStore } from 'pinia'
import { ref, watchEffect } from 'vue'

export type Theme = 'dark' | 'light'

const KEY = 'theme'

function load(): Theme {
  try {
    const saved = localStorage.getItem(KEY)
    if (saved === 'dark' || saved === 'light') return saved
  } catch {
    // storage blocked: follow the system
  }
  return matchMedia('(prefers-color-scheme: light)').matches ? 'light' : 'dark'
}

/** Dark or light; the choice is kept in this browser. */
export const useTheme = defineStore('theme', () => {
  const theme = ref<Theme>(load())

  watchEffect(() => {
    document.documentElement.dataset.theme = theme.value
  })

  function set(next: Theme) {
    theme.value = next
    try {
      localStorage.setItem(KEY, next)
    } catch {
      // the theme still applies for this visit
    }
  }

  return { theme, set }
})
