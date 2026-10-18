<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import XIcon from './XIcon.vue'
import XMenu, { type MenuItem } from './XMenu.vue'

/**
 * Pages of a list. With sizes, a choice of rows per page at the end; with
 * summary, phones show «Page 2 of 5» at the start in place of the compact count.
 */
const props = defineProps<{ pages: number; sizes?: number[]; summary?: boolean }>()
const page = defineModel<number>({ default: 1 })
const size = defineModel<number>('size', { default: 0 })
const { t } = useI18n()

const sizeItems = computed<MenuItem[]>(() => (props.sizes ?? []).map((n) => ({ key: String(n), label: t('kit.perPage', { n }), checked: n === size.value })))

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
  <nav :class="['x-pager', { 'x-pager--summary': summary }]" :aria-label="t('kit.pages')">
    <span v-if="summary" class="x-pager__summary">{{ t('kit.pageOf', { page, pages }) }}</span>
    <button type="button" class="x-pager__chip x-pager__step" :aria-label="t('kit.previousPage')" :disabled="page <= 1" @click="go(page - 1)">
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
    <button type="button" class="x-pager__chip x-pager__step" :aria-label="t('kit.nextPage')" :disabled="page >= pages" @click="go(page + 1)">
      <XIcon name="chevron" :size="12" :stroke="1.8" />
    </button>
    <XMenu v-if="sizes?.length" v-slot="{ props: trigger }" :items="sizeItems" placement="bottom-end" @select="(size = Number($event)), (page = 1)">
      <button type="button" class="x-pager__chip x-pager__size" v-bind="trigger">
        {{ t('kit.perPage', { n: size }) }}
        <XIcon name="caret" :size="12" :stroke="1.7" />
      </button>
    </XMenu>
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

.x-pager__size {
  gap: 6px;
  margin-left: 4px;
  padding: 0 9px 0 10px;
  white-space: nowrap;
}

.x-pager__size :deep(.x-icon) {
  color: var(--dim);
}

.x-pager__summary {
  display: none;
  flex: 1;
  min-width: 0;
  color: var(--mut);
  font: 400 var(--fs-12) / var(--lh-tight) var(--font-sans);
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

  .x-pager__summary {
    display: block;
  }

  .x-pager--summary .x-pager__compact {
    display: none;
  }

  .x-pager__size {
    height: var(--ctl-h);
    border-radius: var(--ctl-r);
  }

  .x-pager__step {
    min-width: var(--ctl-h);
    height: var(--ctl-h);
    border-radius: var(--ctl-r);
  }
}
</style>
