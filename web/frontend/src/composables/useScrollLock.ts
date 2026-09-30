import { onBeforeUnmount, watch, type Ref } from 'vue'

let locks = 0

/** Lock page scrolling while `active` is true (reference counted). */
export function useScrollLock(active: Ref<boolean>) {
  let held = false
  const set = (on: boolean) => {
    if (on === held) return
    held = on
    locks += on ? 1 : -1
    document.documentElement.style.overflow = locks > 0 ? 'hidden' : ''
  }
  watch(active, set, { immediate: true })
  onBeforeUnmount(() => set(false))
}
