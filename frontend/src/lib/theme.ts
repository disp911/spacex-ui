import { ref, watch } from 'vue'

// Shares the "dark-mode" key with the old pages, so the theme carries over
// while the panel is half old and half new.
const KEY = 'dark-mode'

function initial(): boolean {
  try {
    const v = localStorage.getItem(KEY)
    return v === null ? true : v === 'true'
  } catch {
    return true
  }
}

export const isDark = ref(initial())

function apply() {
  document.documentElement.dataset.theme = isDark.value ? 'dark' : 'light'
}
apply()

watch(isDark, (v) => {
  apply()
  try {
    localStorage.setItem(KEY, String(v))
  } catch {
    /* private mode */
  }
})

export function toggleTheme() {
  isDark.value = !isDark.value
}
