import { nextTick, onBeforeUnmount, ref, watch, type Ref } from 'vue'

export type Placement = 'bottom-start' | 'bottom-end' | 'top'

const GAP = 6
const EDGE = 8
const MAX_HEIGHT = 360

/**
 * Keeps a teleported popup (dropdown, menu, tooltip) next to its anchor while
 * it is open: below the anchor, or above it when there is no room below.
 * Clicks outside both elements and Escape call onDismiss.
 */
export function useFloating(
  anchor: Ref<HTMLElement | null>,
  popup: Ref<HTMLElement | null>,
  open: Ref<boolean>,
  options: { placement?: Placement; matchWidth?: boolean; onDismiss?: () => void } = {},
) {
  // Fixed at the corner until measured, so the popup takes its natural size.
  const initial = { position: 'fixed', top: '0px', left: '0px' }
  const style = ref<Record<string, string>>(initial)
  let frame = 0

  function update() {
    const a = anchor.value?.getBoundingClientRect()
    const p = popup.value
    if (!a || !p) return
    const width = options.matchWidth ? a.width : p.offsetWidth
    const height = p.offsetHeight
    const vw = window.innerWidth
    const vh = window.innerHeight
    const placement = options.placement ?? 'bottom-start'

    const below = vh - a.bottom - GAP - EDGE
    const above = a.top - GAP - EDGE
    const up = placement === 'top' ? height <= above || above > below : height > below && above > below
    const room = Math.min(MAX_HEIGHT, up ? above : below)
    const top = up ? a.top - GAP - Math.min(height, room) : a.bottom + GAP

    let left = placement === 'bottom-end' ? a.right - width : placement === 'top' ? a.left + a.width / 2 - width / 2 : a.left
    left = Math.min(Math.max(left, EDGE), vw - width - EDGE)

    style.value = {
      position: 'fixed',
      top: `${Math.round(top)}px`,
      left: `${Math.round(left)}px`,
      ...(options.matchWidth ? { minWidth: `${Math.round(a.width)}px` } : {}),
      maxHeight: `${Math.max(0, room)}px`,
    }
  }

  function schedule() {
    cancelAnimationFrame(frame)
    frame = requestAnimationFrame(update)
  }

  function onPointerDown(e: PointerEvent) {
    const t = e.target as Node
    if (anchor.value?.contains(t) || popup.value?.contains(t)) return
    options.onDismiss?.()
  }

  function onKeyDown(e: KeyboardEvent) {
    if (e.key === 'Escape') {
      e.stopPropagation()
      options.onDismiss?.()
    }
  }

  function bind(on: boolean) {
    if (on) {
      window.addEventListener('scroll', schedule, true)
      window.addEventListener('resize', schedule)
      document.addEventListener('pointerdown', onPointerDown, true)
      document.addEventListener('keydown', onKeyDown, true)
    } else {
      window.removeEventListener('scroll', schedule, true)
      window.removeEventListener('resize', schedule)
      document.removeEventListener('pointerdown', onPointerDown, true)
      document.removeEventListener('keydown', onKeyDown, true)
    }
  }

  watch(open, async (on) => {
    bind(on)
    if (!on) return
    style.value = initial
    await nextTick()
    update()
  }, { flush: 'post' })

  onBeforeUnmount(() => {
    cancelAnimationFrame(frame)
    bind(false)
  })

  return { style, update: schedule }
}
