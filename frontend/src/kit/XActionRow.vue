<script setup lang="ts">
import XIcon from './XIcon.vue'
import XSpinner from './XSpinner.vue'
import type { IconName } from './icons'

/**
 * A row of XActionList: a title, a description (the slot) and a square
 * accent button. The whole row is the button.
 */
defineProps<{ title: string; icon: IconName; loading?: boolean; disabled?: boolean }>()
</script>

<template>
  <button
    type="button"
    class="x-action-row"
    :disabled="disabled"
    :aria-disabled="loading || undefined"
    :aria-busy="loading || undefined"
  >
    <span class="x-action-row__body">
      <span class="x-action-row__title">{{ title }}</span>
      <span v-if="$slots.default" class="x-action-row__text"><slot /></span>
    </span>
    <span class="x-action-row__icon" aria-hidden="true">
      <XSpinner v-if="loading" :size="15" />
      <XIcon v-else :name="icon" :size="16" :stroke="1.6" class="x-action-row__glyph" />
    </span>
  </button>
</template>

<style scoped>
.x-action-row {
  display: flex;
  align-items: center;
  gap: 14px;
  width: 100%;
  padding: 14px 16px;
  border: 0;
  background: transparent;
  color: var(--fg);
  text-align: left;
  cursor: pointer;
  outline: none;
  transition: background-color var(--t-fast);
}

.x-action-row:hover:not(:disabled),
.x-action-row.is-hover {
  background: var(--ctrl);
}

.x-action-row:focus-visible,
.x-action-row.is-focus {
  box-shadow: inset 0 0 0 2px var(--acc-line);
}

.x-action-row[aria-disabled='true'] {
  pointer-events: none;
}

.x-action-row:disabled {
  cursor: not-allowed;
  opacity: .5;
}

.x-action-row__body {
  display: flex;
  flex: 1;
  flex-direction: column;
  gap: 4px;
  min-width: 0;
}

.x-action-row__title {
  font: 600 var(--fs-13) / var(--lh-tight) var(--font-sans);
  letter-spacing: -.15px;
}

.x-action-row__text {
  color: var(--mut);
  font: 400 var(--fs-12) / var(--lh-text) var(--font-sans);
  text-wrap: pretty;
}

/* Inline code in the description, e.g. the .db extension. */
.x-action-row__text :deep(code) {
  color: var(--fg2);
  font-family: var(--font-mono);
}

.x-action-row__glyph {
  width: var(--btn-icon);
  height: var(--btn-icon);
}

.x-action-row__icon {
  display: flex;
  flex: none;
  align-items: center;
  justify-content: center;
  width: var(--ctl-h);
  height: var(--ctl-h);
  border-radius: var(--ctl-r);
  background: var(--acc);
  color: var(--acc-ink);
  transition: background-color var(--t-fast), box-shadow var(--t-fast);
}

.x-action-row:hover:not(:disabled) .x-action-row__icon,
.x-action-row.is-hover .x-action-row__icon {
  background: var(--acc-hover);
  box-shadow: 0 0 0 4px var(--acc-ring), 0 0 18px var(--acc-glow);
}

@media (max-width: 767px) {
  .x-action-row {
    padding: 16px 14px;
  }

  .x-action-row__title {
    font-size: var(--fs-14);
  }
}
</style>
