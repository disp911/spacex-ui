import { onBeforeUnmount, ref } from 'vue'

/** True while the screen is a phone's: the width where the kit switches to its phone layout. */
export function usePhone() {
  const query = matchMedia('(max-width: 767px)')
  const phone = ref(query.matches)
  const update = () => (phone.value = query.matches)
  query.addEventListener('change', update)
  onBeforeUnmount(() => query.removeEventListener('change', update))
  return phone
}
