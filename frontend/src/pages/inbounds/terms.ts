import { computed, ref, watch } from 'vue'

type Json = Record<string, any>

const GB = 1024 ** 3
const DAY = 86_400_000

/**
 * A client's traffic limit and end as its forms edit them, the old panel's
 * way: the limit in gigabytes; a delayed start keeps the length in days as a
 * negative end, and turning it either way starts over without an end.
 */
export function useClientTerms(client: () => Json) {
  const totalGb = computed({
    get: () => Math.round(((client().totalGB ?? 0) / GB) * 100) / 100,
    set: (gb: number | null) => (client().totalGB = Math.round((gb ?? 0) * GB)),
  })

  const delayed = ref(false)
  watch(client, (c) => (delayed.value = (c.expiryTime ?? 0) < 0), { immediate: true })
  function setDelayed(on: boolean) {
    delayed.value = on
    client().expiryTime = 0
  }

  const days = computed({
    get: () => Math.round(-(client().expiryTime ?? 0) / DAY),
    set: (n: number | null) => (client().expiryTime = -(n ?? 0) * DAY),
  })
  const expiry = computed({
    get: () => Math.max(0, client().expiryTime ?? 0),
    set: (ms: number) => (client().expiryTime = ms),
  })

  return { totalGb, delayed, setDelayed, days, expiry }
}
