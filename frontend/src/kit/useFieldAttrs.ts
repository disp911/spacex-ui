import { computed, useAttrs } from 'vue'

/**
 * Field components put class and style on their wrapper and every other
 * attribute (id, name, autocomplete, aria-*) on the native control inside.
 * Components using it must set defineOptions({ inheritAttrs: false }).
 */
export function useFieldAttrs() {
  const attrs = useAttrs()
  const wrapper = computed(() => ({ class: attrs.class, style: attrs.style }))
  const control = computed(() => {
    const { class: _c, style: _s, ...rest } = attrs
    return rest
  })
  return { wrapper, control }
}
