<script setup lang="ts">
import XIcon from './XIcon.vue'

/**
 * A short warning under a field, e.g. the address the panel will have after a
 * restart. Values in the text go in <code>: mono, brighter, broken anywhere.
 */
</script>

<template>
  <p class="x-field-note">
    <XIcon name="warning" :size="13" :stroke="1.4" class="x-field-note__icon" />
    <span class="x-field-note__text"><slot /></span>
  </p>
</template>

<style scoped>
.x-field-note {
  display: flex;
  align-items: flex-start;
  gap: 6px;
  max-width: 100%;
  margin: 0;
  color: var(--warn-x);
  font: 400 var(--fs-12) / var(--lh-text) var(--font-sans);
}

.x-field-note__icon {
  flex: none;
  margin-top: 2px;
  color: var(--warn);
}

.x-field-note__text {
  min-width: 0;
}

.x-field-note__text :slotted(code) {
  color: var(--warn-t);
  font: inherit;
  font-family: var(--font-mono);
  word-break: break-all;
}
</style>
