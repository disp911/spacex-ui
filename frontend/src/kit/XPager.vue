<script setup lang="ts">
import { computed } from 'vue'
import XIcon from './XIcon.vue'

const props = defineProps<{ pages: number }>()
const page = defineModel<number>({ default: 1 })

/** Page numbers with null for a gap: 1 … 4 5 6 … 12. */
const items = computed(() => {
  const n = props.pages
  const p = page.value
  if (n <= 7) return Array.from({ length: n }, (_, i) => i + 1)
  const out: (number | null)[] = [1]
  const from = Math.max(2, Math.min(p - 1, n - 4))
  const to = Math.min(n - 1, Math.max(p + 1, 5))
  if (from > 2) out.push(null)
  for (let i = from; i <= to; i++) out.push(i)
  if (to < n - 1) out.push(null)
  out.push(n)
  return out
})

function go(to: number) {
  page.value = Math.min(Math.max(to, 1), props.pages)
}
</script>

<template>
  <nav class="x-pager" aria-label="Pages">
    <button type="button" class="x-pager__chip x-pager__step" aria-label="Previous page" :disabled="page <= 1" @click="go(page - 1)">
      <XIcon name="chevronLeft" :size="12" :stroke="1.8" />
    </button>
    <template v-for="(item, i) in items" :key="i">
      <span v-if="item === null" class="x-pager__gap">…</span>
      <button
        v-else
        type="button"
        :class="['x-pager__chip', 'x-pager__number', { 'is-current': item === page }]"
        :aria-current="item === page ? 'page' : undefined"
        @click="go(item)"
      >
        {{ item }}
      </button>
    </template>
    <span class="x-pager__compact">{{ page }} / {{ pages }}</span>
    <button type="button" class="x-pager__chip x-pager__step" aria-label="Next page" :disabled="page >= pages" @click="go(page + 1)">
      <XIcon name="chevron" :size="12" :stroke="1.8" />
    </button>
  </nav>
</template>

<style scoped>
.x-pager {
  display: flex;
  align-items: center;
  gap: 4px;
}

.x-pager__chip {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: var(--ctl-sm-h);
  height: var(--ctl-sm-h);
  padding: 0 6px;
  border: 1px solid var(--field-line);
  border-radius: var(--ctl-sm-r);
  background: var(--ctrl);
  color: var(--fg2);
  font: 500 var(--fs-12) / 1 var(--font-mono);
  outline: none;
  transition: border-color var(--t-fast), background-color var(--t-fast), box-shadow var(--t-fast);
}

.x-pager__chip:hover:not(:disabled),
.x-pager__chip.is-hover {
  border-color: var(--field-hover);
}

.x-pager__chip:focus-visible,
.x-pager__chip.is-focus {
  border-color: var(--acc-line);
  box-shadow: 0 0 0 3px var(--acc-ring);
}

.x-pager__chip:active:not(:disabled),
.x-pager__chip.is-active {
  background: var(--press);
}

.x-pager__chip.is-current {
  border-color: var(--acc);
  background: var(--acc);
  color: var(--acc-ink);
}

.x-pager__chip:disabled {
  color: var(--faint);
  opacity: .55;
  cursor: not-allowed;
}

.x-pager__gap {
  min-width: 20px;
  color: var(--faint);
  font: 500 var(--fs-12) / 1 var(--font-mono);
  text-align: center;
}

.x-pager__compact {
  display: none;
  padding: 0 8px;
  color: var(--fg2);
  font: 500 var(--fs-13) / 1 var(--font-mono);
}

@media (max-width: 767px) {
  .x-pager__number,
  .x-pager__gap {
    display: none;
  }

  .x-pager__compact {
    display: inline;
  }

  .x-pager__step {
    min-width: var(--ctl-h);
    height: var(--ctl-h);
    border-radius: var(--ctl-r);
  }
}
</style>
