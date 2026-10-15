<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref } from 'vue'

/** fit: buttons share the row by the length of their labels, e.g. a short Save beside Restart Xray. */
defineProps<{ fit?: boolean }>()

/**
 * Phones: a bar pinned to the bottom of a page with a line of state above the
 * page's main buttons (the slot), e.g. Save and Restart; buttons with text
 * share the row, square ones and anything else keep their size.
 * Without anything in the state slot the line takes no room.
 * While it is shown, toasts rise 20px above it (--toast-bottom, read by XToaster).
 */
const bar = ref<HTMLElement | null>(null)
const root = document.documentElement
let observer: ResizeObserver | undefined

function lift() {
  const h = bar.value?.offsetHeight ?? 0
  if (h) root.style.setProperty('--toast-bottom', `${h + 20}px`)
  else root.style.removeProperty('--toast-bottom')
}

onMounted(() => {
  observer = new ResizeObserver(lift)
  if (bar.value) observer.observe(bar.value)
})
onBeforeUnmount(() => {
  observer?.disconnect()
  root.style.removeProperty('--toast-bottom')
})
</script>

<template>
  <div ref="bar" class="x-action-bar">
    <div v-if="$slots.state" class="x-action-bar__state"><slot name="state" /></div>
    <div :class="['x-action-bar__buttons', { 'x-action-bar__buttons--fit': fit }]"><slot /></div>
  </div>
</template>

<style scoped>
.x-action-bar {
  position: sticky;
  bottom: 0;
  z-index: 20;
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 12px 14px calc(12px + env(safe-area-inset-bottom));
  border-top: 1px solid var(--div);
  background: var(--side);
}

.x-action-bar__state {
  display: flex;
  align-items: center;
  min-height: 16px;
  color: var(--dim);
  font: 400 var(--fs-12) / var(--lh-tight) var(--font-sans);
}

.x-action-bar__buttons {
  display: flex;
  align-items: center;
  gap: 8px;
}

.x-action-bar__buttons > :deep(.x-btn:not(.x-btn--icon)) {
  flex: 1;
  min-width: 0;
}

.x-action-bar__buttons--fit > :deep(.x-btn:not(.x-btn--icon)) {
  flex: 1 1 auto;
}
</style>
