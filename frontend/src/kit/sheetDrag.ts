/**
 * Pulling a phone sheet down by its handle or title: the sheet follows the
 * finger; let go past a quarter of its height, or with a flick, and it slides
 * away and closes, otherwise it springs back. The zone the drag starts on
 * needs touch-action: none, or the browser takes the gesture for a scroll.
 */
export function startSheetDrag(e: PointerEvent, sheet: HTMLElement, close: () => void) {
  if (e.button !== 0) return
  const zone = e.currentTarget as HTMLElement
  const startY = e.clientY
  let dy = 0
  let lastY = startY
  let lastT = performance.now()
  let speed = 0

  // Keeps the moves coming when the finger leaves the zone; a pointer that is already gone cannot be captured.
  try {
    zone.setPointerCapture(e.pointerId)
  } catch {
    // the drag still follows moves over the zone
  }
  sheet.style.transition = 'none'

  function move(ev: PointerEvent) {
    const now = performance.now()
    speed = (ev.clientY - lastY) / Math.max(1, now - lastT)
    lastY = ev.clientY
    lastT = now
    dy = Math.max(0, ev.clientY - startY)
    sheet.style.transform = dy ? `translateY(${dy}px)` : ''
  }

  function end() {
    zone.removeEventListener('pointermove', move)
    zone.removeEventListener('pointerup', end)
    zone.removeEventListener('pointercancel', end)
    if (dy > sheet.offsetHeight / 4 || (dy > 24 && speed > 0.6)) {
      // Slides on from where the finger left it; the closing transition fades the scrim meanwhile.
      sheet.style.transition = 'transform var(--t-enter) ease-in'
      sheet.style.transform = 'translateY(100%)'
      close()
      return
    }
    sheet.style.transition = 'transform var(--t-enter) var(--ease-out)'
    sheet.style.transform = ''
    sheet.addEventListener('transitionend', () => (sheet.style.transition = ''), { once: true })
  }

  zone.addEventListener('pointermove', move)
  zone.addEventListener('pointerup', end)
  zone.addEventListener('pointercancel', end)
}
