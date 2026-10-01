<script setup lang="ts">
import { computed } from 'vue'
import XIcon from './XIcon.vue'

// Numbered pagination with ellipses; `compact` shows only prev/next.
const page = defineModel<number>({ required: true })
const props = defineProps<{ pages: number; compact?: boolean }>()

const items = computed<(number | '…')[]>(() => {
  const n = props.pages
  const p = page.value
  if (n <= 7) return Array.from({ length: n }, (_, i) => i + 1)
  const out: (number | '…')[] = [1]
  if (p > 3) out.push('…')
  for (let i = Math.max(2, p - 1); i <= Math.min(n - 1, p + 1); i++) out.push(i)
  if (p < n - 2) out.push('…')
  out.push(n)
  return out
})
</script>

<template>
  <nav class="x-pag" :class="{ 'is-compact': compact }" aria-label="Pagination">
    <button type="button" class="x-pag__btn" :disabled="page <= 1" aria-label="Previous" @click="page--">
      <XIcon name="chevron-left" :size="13" :stroke-width="1.6" />
    </button>
    <template v-if="!compact">
      <template v-for="(it, i) in items" :key="i">
        <span v-if="it === '…'" class="x-pag__gap">…</span>
        <button v-else type="button" class="x-pag__btn" :class="{ 'is-on': it === page }" :aria-current="it === page ? 'page' : undefined" @click="page = it">
          {{ it }}
        </button>
      </template>
    </template>
    <button type="button" class="x-pag__btn" :disabled="page >= pages" aria-label="Next" @click="page++">
      <XIcon name="chevron-right" :size="13" :stroke-width="1.6" />
    </button>
  </nav>
</template>

<style scoped>
.x-pag {
  display: flex;
  align-items: center;
  gap: 6px;
}
.x-pag__btn {
  flex: none;
  min-width: 28px;
  height: 28px;
  padding: 0 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--r-sm);
  border: 1px solid var(--border-control);
  background: var(--control);
  color: var(--text-2);
  font: var(--fw-medium) var(--fs-sm) var(--font-mono);
  transition: border-color var(--dur-fast) ease, background var(--dur-fast) ease, color var(--dur-fast) ease;
}
.x-pag__btn:not(:disabled):not(.is-on):hover {
  border-color: var(--accent-focus);
  color: var(--text);
}
.x-pag__btn:disabled {
  color: var(--text-disabled);
  cursor: not-allowed;
}
.x-pag__btn.is-on {
  background: var(--accent);
  border-color: var(--accent);
  color: var(--on-accent);
}
.x-pag__gap {
  min-width: 20px;
  text-align: center;
  color: var(--text-5);
  font-family: var(--font-mono);
}
.x-pag.is-compact .x-pag__btn {
  width: var(--h-md);
  height: var(--h-md);
  border-radius: var(--r-md);
}
</style>
