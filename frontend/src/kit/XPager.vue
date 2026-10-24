<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import XIcon from './XIcon.vue'
import XMenu, { type MenuItem } from './XMenu.vue'
import { useSlidingThumb } from './slidingThumb'

/**
 * Pages of a list. With sizes, a choice of rows per page at the end; with
 * summary, phones show «Page 2 of 5» at the start in place of the compact count;
 * strip makes it one line of every page number with the current one lit, the
 * whole width on phones: a tap picks a page, a finger or the mouse drawn along
 * the line picks the one it is let go on, and a line too long for its place
 * scrolls sideways with the current page in its middle.
 */
const props = defineProps<{ pages: number; sizes?: number[]; summary?: boolean; strip?: boolean }>()
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

// The current page's accent travels between the numbers like liquid.
const nav = ref<HTMLElement | null>(null)
const thumb = useSlidingThumb(nav, '.x-pager__number.is-current', [page, items])

function go(to: number) {
  page.value = Math.min(Math.max(to, 1), props.pages)
}

// --- Strip -------------------------------------------------------------------

const line = ref<HTMLElement | null>(null)
/** The page lit on the strip: the current one, or the one under a moving finger. */
const lit = ref(page.value)
watch(page, (p) => (lit.value = p))
const stripThumb = useSlidingThumb(line, '.x-pager__cell.is-current', [lit, () => props.pages])

// A line that fits is drawn along to pick; one too long scrolls instead.
const fits = ref(true)
let sizer: ResizeObserver | null = null
const measure = () => (fits.value = !line.value || line.value.scrollWidth <= line.value.clientWidth + 1)
onMounted(() => {
  if (!line.value) return
  sizer = new ResizeObserver(measure)
  sizer.observe(line.value)
  center(false)
})
onBeforeUnmount(() => sizer?.disconnect())

function center(smooth = true) {
  const box = line.value
  const el = box?.querySelector<HTMLElement>('.x-pager__cell.is-current')
  if (!box || !el || fits.value) return
  box.scrollTo({ left: el.offsetLeft - (box.clientWidth - el.offsetWidth) / 2, behavior: smooth ? 'smooth' : 'auto' })
}
watch([page, () => props.pages], () => nextTick(() => (measure(), center())))

let drawing = false
let moved = false
let startX = 0
function drawStart(e: PointerEvent) {
  drawing = true
  moved = false
  startX = e.clientX
}
function drawMove(e: PointerEvent) {
  if (!drawing || !fits.value) return
  if (Math.abs(e.clientX - startX) > 6) moved = true
  if (!moved) return
  const cell = document.elementFromPoint(e.clientX, e.clientY)?.closest<HTMLElement>('.x-pager__cell')
  if (cell?.dataset.page) lit.value = Number(cell.dataset.page)
}
function drawEnd() {
  if (!drawing) return
  drawing = false
  if (moved) go(lit.value)
}
</script>

<template>
  <nav ref="nav" :class="['x-pager', { 'x-pager--summary': summary, 'x-pager--strip': strip }]" :aria-label="t('kit.pages')">
    <span
      v-if="thumb"
      :class="['x-pager__thumb', thumb.ahead && `x-slide--${thumb.ahead}`]"
      :style="{ left: `${thumb.left}px`, right: `${thumb.right}px`, top: `${thumb.top}px`, height: `${thumb.height}px` }"
      aria-hidden="true"
    />
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
        <span class="x-pager__n">{{ item }}</span>
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
    <div
      v-if="strip"
      ref="line"
      :class="['x-pager__line', { 'is-fit': fits }]"
      @pointerdown="drawStart"
      @pointermove="drawMove"
      @pointerup="drawEnd"
      @pointercancel="drawing = false"
    >
      <span
        v-if="stripThumb"
        :class="['x-pager__line-thumb', stripThumb.ahead && `x-slide--${stripThumb.ahead}`]"
        :style="{ left: `${stripThumb.left}px`, right: `${stripThumb.right}px` }"
        aria-hidden="true"
      />
      <button
        v-for="n in pages"
        :key="n"
        type="button"
        :data-page="n"
        :class="['x-pager__cell', { 'is-current': n === lit }]"
        :aria-current="n === page ? 'page' : undefined"
        @click="go(n)"
      >
        {{ n }}
      </button>
    </div>
  </nav>
</template>

<style scoped>
.x-pager {
  position: relative;
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

/* Over the chips' backgrounds and under their numbers, which are lifted above it. */
.x-pager__thumb {
  position: absolute;
  z-index: 1;
  border-radius: var(--ctl-sm-r);
  background: var(--acc);
  pointer-events: none;
}

.x-pager__n {
  position: relative;
  z-index: 2;
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

  /* On phones the strip takes the whole line, its numbers sharing it. */
  .x-pager--strip {
    width: 100%;
  }

  .x-pager--strip .x-pager__line {
    flex: 1;
  }
}

/* The strip: the line of numbers stands in for the steps, the chips and the count. */
.x-pager--strip {
  min-width: 0;
}

.x-pager--strip > :not(.x-pager__line) {
  display: none;
}

.x-pager__line {
  position: relative;
  display: flex;
  flex: 0 1 auto;
  gap: 2px;
  min-width: 0;
  height: var(--ctl-h);
  padding: 2px;
  overflow-x: auto;
  border: 1px solid var(--field-line);
  border-radius: var(--ctl-r);
  background: var(--field);
  scrollbar-width: none;
  touch-action: pan-x;
  user-select: none;
}

.x-pager__line::-webkit-scrollbar {
  display: none;
}

.x-pager__line.is-fit {
  touch-action: none;
}

/* Under the numbers, which are lifted above it, as in the chips. */
.x-pager__line-thumb {
  position: absolute;
  top: 2px;
  bottom: 2px;
  z-index: 1;
  border-radius: var(--ctl-sm-r);
  background: var(--acc);
  pointer-events: none;
}

/* A number is at least a square of the line's inner height. */
.x-pager__cell {
  position: relative;
  z-index: 2;
  flex: 1 0 auto;
  min-width: calc(var(--ctl-h) - 6px);
  padding: 0 8px;
  border: 0;
  border-radius: var(--ctl-sm-r);
  background: transparent;
  color: var(--fg2);
  font: 500 var(--fs-13) / 1 var(--font-mono);
  outline: none;
  transition: color var(--t-med);
}

.x-pager__cell:hover:not(.is-current) {
  color: var(--fg);
}

.x-pager__cell:focus-visible {
  box-shadow: inset 0 0 0 2px var(--acc-line);
}

.x-pager__cell.is-current {
  color: var(--acc-ink);
  font-weight: 600;
}
</style>
