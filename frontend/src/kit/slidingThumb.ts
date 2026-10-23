import { nextTick, onBeforeUnmount, onMounted, ref, watch, type Ref, type WatchSource } from 'vue'

/** Where a sliding thumb stands, and which edge leads when it travels. */
export interface Thumb {
  left: number
  right: number
  top: number
  bottom: number
  height: number
  ahead: 'left' | 'right' | 'up' | 'down' | null
}

/**
 * One highlight that travels to the chosen item of a group (segments, tabs,
 * pages, a side menu) instead of jumping. The thumb is placed in the box's own
 * coordinates: the box is its positioned parent. On a change the leading edge
 * gets x-slide--left/right, or x-slide--up/down for a column (kit.css), so the
 * thumb stretches on the way and settles with a little spring; a resize moves
 * it without animation.
 */
export function useSlidingThumb(box: Ref<HTMLElement | null>, selector: string, changes: WatchSource[], axis: 'x' | 'y' = 'x') {
  const thumb = ref<Thumb | null>(null)

  function place(animate: boolean) {
    const root = box.value
    const el = root?.querySelector<HTMLElement>(selector)
    if (!root || !el) {
      thumb.value = null
      return
    }
    const left = el.offsetLeft
    const top = el.offsetTop
    // From the box's inner edges, as an absolute child is placed; in a scrolling strip it may lie past them.
    const right = root.clientWidth - left - el.offsetWidth
    const bottom = root.clientHeight - top - el.offsetHeight
    const was = thumb.value
    let ahead: Thumb['ahead'] = null
    if (animate && was) {
      if (axis === 'x') ahead = left > was.left ? 'right' : left < was.left ? 'left' : null
      else ahead = top > was.top ? 'down' : top < was.top ? 'up' : null
    }
    thumb.value = { left, right, top, bottom, height: el.offsetHeight, ahead }
  }

  let observer: ResizeObserver | null = null
  onMounted(() => {
    place(false)
    observer = new ResizeObserver(() => place(false))
    if (box.value) observer.observe(box.value)
  })
  onBeforeUnmount(() => observer?.disconnect())
  watch(changes, () => nextTick(() => place(true)), { deep: true })

  return thumb
}
