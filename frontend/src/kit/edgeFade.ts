import { computed, nextTick, onBeforeUnmount, onMounted, ref, type Ref } from 'vue'

/**
 * A box that scrolls sideways fades its content out at an edge with more
 * beyond it, and only there: at the start only the end fades, in the middle
 * both, at the end only the start; a box that fits does not fade. Bind
 * classes to the scrolling box (kit.css draws the fade with a mask, so put it
 * on what scrolls, not on a frame around it) and call update when what it
 * holds changes; scrolling and resizing update it on their own.
 */
export function useEdgeFade(box: Ref<HTMLElement | null>) {
  const more = ref({ left: false, right: false })

  function update() {
    const el = box.value
    if (!el) return
    more.value = { left: el.scrollLeft > 1, right: el.scrollLeft + el.clientWidth < el.scrollWidth - 1 }
  }

  let sizer: ResizeObserver | null = null
  onMounted(() => {
    const el = box.value
    if (!el) return
    el.addEventListener('scroll', update, { passive: true })
    sizer = new ResizeObserver(update)
    sizer.observe(el)
    nextTick(update)
    // Labels widen once the web font is in.
    document.fonts?.ready.then(update)
  })
  onBeforeUnmount(() => {
    box.value?.removeEventListener('scroll', update)
    sizer?.disconnect()
  })

  const classes = computed(() => ({ 'x-fade-start': more.value.left, 'x-fade-end': more.value.right }))
  return { more, classes, update }
}
