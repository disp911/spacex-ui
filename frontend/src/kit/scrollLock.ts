// Windows can open over windows: the page scrolls again only when the last one closes.
let locks = 0

export function lockScroll(lock: boolean) {
  locks = Math.max(0, locks + (lock ? 1 : -1))
  document.documentElement.style.overflow = locks ? 'hidden' : ''
}
