import { inject, onBeforeUnmount, ref, type InjectionKey, type Ref } from 'vue'

/** True while the screen is a phone's: the width where the kit switches to its phone layout. */
export function usePhone() {
  const query = matchMedia('(max-width: 767px)')
  const phone = ref(query.matches)
  const update = () => (phone.value = query.matches)
  query.addEventListener('change', update)
  onBeforeUnmount(() => query.removeEventListener('change', update))
  return phone
}

/** Set by a list that has gone to its phone columns for want of room; what it holds follows it. */
export const compactKey: InjectionKey<Ref<boolean>> = Symbol('compact')

/** True on a phone, or inside a list that shows its phone columns because it is narrow. */
export function useCompact() {
  const phone = usePhone()
  return inject(compactKey, phone)
}
