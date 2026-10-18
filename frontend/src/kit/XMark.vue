<script setup lang="ts">
import { computed } from 'vue'

/**
 * Text with what a search found in it marked, case aside, e.g. a remark or a
 * client's email while the list is searched. Without a query it is the text.
 */
const props = defineProps<{ text: string | number; query?: string }>()

const parts = computed(() => {
  const text = String(props.text)
  const q = props.query?.toLowerCase() ?? ''
  if (!q) return [{ text, hit: false }]
  const lower = text.toLowerCase()
  const out: { text: string; hit: boolean }[] = []
  let at = 0
  for (let k = lower.indexOf(q); k >= 0; k = lower.indexOf(q, at)) {
    if (k > at) out.push({ text: text.slice(at, k), hit: false })
    out.push({ text: text.slice(k, k + q.length), hit: true })
    at = k + q.length
  }
  if (at < text.length) out.push({ text: text.slice(at), hit: false })
  return out
})
</script>

<template>
  <template v-for="(p, i) in parts" :key="i"><mark v-if="p.hit" class="x-mark">{{ p.text }}</mark><template v-else>{{ p.text }}</template></template>
</template>

<style scoped>
.x-mark {
  border-radius: var(--r-4);
  background: var(--hl);
  box-shadow: 0 0 0 1px var(--hl);
  color: var(--fg);
}
</style>
