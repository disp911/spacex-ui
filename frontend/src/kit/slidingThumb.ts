import { nextTick, onBeforeUnmount, onMounted, ref, watch, type Ref, type WatchSource } from 'vue'

/** Where a sliding thumb stands, and which edge leads when it travels. */
export interface Thumb {
  left: number
  right: number
  top: number
  height: number
  ahead: 'left' | 'right' | null
}

/**
 * One highlight that travels to the chosen item of a group (segments, tabs,
 * pages) instead of jumping. The thumb is placed in the box's own coordinates:
 * the box is its positioned parent. On a change the leading edge gets
 * x-slide--left or x-slide--right (kit.css), so the thumb stretches on the way
 * and settles with a little spring; a resize moves it without animation.
 */
export function useSlidingThumb(box: Ref<HTMLElement | null>, selector: string, changes: WatchSource[]) {
  const thumb = ref<Thumb | null>(null)

  function place(animate: boolean) {
    const root = box.value
    const el = root?.querySelector<HTMLElement>(selector)
    if (!root || !el) {
      thumb.value = null
      return
    }
    const left = el.offsetLeft
    // From the box's inner edge, as an absolute child is placed; in a scrolling strip it may lie past it.
    const right = root.clientWidth - el.offsetLeft - el.offsetWidth
    const was = thumb.value
    const ahead = animate && was ? (left > was.left ? 'right' : left < was.left ? 'left' : null) : null
    thumb.value = { left, right, top: el.offsetTop, height: el.offsetHeight, ahead }
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
