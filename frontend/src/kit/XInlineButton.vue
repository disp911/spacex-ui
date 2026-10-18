<script setup lang="ts">
import XIcon from './XIcon.vue'
import XSpinner from './XSpinner.vue'
import type { IconName } from './icons'

/**
 * A 28px icon button that lives inside another element, such as a field or a
 * tile: show a password, clear a value, reveal an address. Square buttons on
 * their own are XButton.
 */
withDefaults(
  defineProps<{
    icon: IconName
    label: string
    iconSize?: number
    stroke?: number
    /** Removes something: red. */
    danger?: boolean
    /** At work: an accent spinner in place of the icon, and no clicks. */
    loading?: boolean
  }>(),
  { iconSize: 14, stroke: 1.35 },
)
</script>

<template>
  <button
    type="button"
    :class="['x-inline-btn', { 'x-inline-btn--danger': danger, 'is-loading': loading }]"
    :aria-label="label"
    :aria-busy="loading || undefined"
    :title="label"
  >
    <XSpinner v-if="loading" :size="12" />
    <XIcon v-else :name="icon" :size="iconSize" :stroke="stroke" />
  </button>
</template>

<style scoped>
.x-inline-btn {
  display: flex;
  flex: none;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  padding: 0;
  border: 0;
  border-radius: var(--r-8);
  background: transparent;
  color: var(--dim);
  outline: none;
  transition: background-color var(--t-fast), color var(--t-fast);
}

.x-inline-btn:hover:not(:disabled),
.x-inline-btn.is-hover {
  background: var(--hov);
  color: var(--fg);
}

.x-inline-btn--danger {
  color: var(--err-t);
}

.x-inline-btn--danger:hover:not(:disabled) {
  background: var(--err-bg);
  color: var(--err-t);
}

.x-inline-btn:focus-visible,
.x-inline-btn.is-focus {
  box-shadow: 0 0 0 2px var(--acc-line);
  color: var(--fg);
}

.x-inline-btn:disabled {
  cursor: not-allowed;
  opacity: .5;
}

/* At work it stays bright, even while its neighbours are held back. */
.x-inline-btn.is-loading {
  color: var(--acc);
  opacity: 1;
  pointer-events: none;
}
</style>
