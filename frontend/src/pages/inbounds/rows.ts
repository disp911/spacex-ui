import { ref, watch, type Ref } from 'vue'

type Json = Record<string, any>
export type Pair = [string, string]

/**
 * Headers edited as rows, as the old panel did, so empty or repeated names
 * stay rows of their own while typed. Xray's object is written from the rows
 * that have both a name and a value: a list of values per name (TCP's
 * request and response), or one value with the later row winning. The rows
 * are read again whenever owner() gives another object, e.g. another
 * transport or inbound.
 */
export function headerRows(owner: () => Json | undefined, key: string, asLists = false): Ref<Pair[]> {
  const rows = ref<Pair[]>([])
  watch(
    owner,
    (o) => {
      const h = (o?.[key] ?? {}) as Record<string, string | string[]>
      rows.value = Object.entries(h).flatMap(([name, v]) => (Array.isArray(v) ? v : [v]).map((value) => [name, String(value)] as Pair))
    },
    { immediate: true },
  )
  watch(
    rows,
    (pairs) => {
      const o = owner()
      if (!o) return
      const out: Record<string, string | string[]> = {}
      for (const [name, value] of pairs) {
        if (!name.trim() || !value.trim()) continue
        if (asLists) ((out[name] ??= []) as string[]).push(value)
        else out[name] = value
      }
      o[key] = out
    },
    { deep: true },
  )
  return rows
}
