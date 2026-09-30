import { onBeforeUnmount, ref } from 'vue'

/** Mobile layout breakpoint — keep in sync with --bp-mobile in tokens.css. */
export const MOBILE_QUERY = '(max-width: 767px)'

const mq = window.matchMedia(MOBILE_QUERY)
const isMobile = ref(mq.matches)
mq.addEventListener('change', (e) => (isMobile.value = e.matches))

export function useIsMobile() {
  return isMobile
}

/** Close something on Escape while mounted. */
export function useEscape(fn: () => void) {
  const h = (e: KeyboardEvent) => {
    if (e.key === 'Escape') fn()
  }
  window.addEventListener('keydown', h)
  onBeforeUnmount(() => window.removeEventListener('keydown', h))
}
