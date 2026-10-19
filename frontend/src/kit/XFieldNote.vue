<script setup lang="ts">
import XIcon from './XIcon.vue'

/**
 * A short line under a field or a form: a warning, e.g. the address the panel
 * will have after a restart, or an error the server gave; info is a quiet hint,
 * e.g. how to use what is listed under it. Values in the text go in <code>:
 * mono, brighter, broken anywhere. attached keeps it 6px under what it is
 * about in a stack of fields 12px apart, e.g. a switch in the inbound window.
 */
withDefaults(defineProps<{ tone?: 'warn' | 'err' | 'info'; attached?: boolean }>(), { tone: 'warn' })
</script>

<template>
  <p :class="['x-field-note', `x-field-note--${tone}`, { 'x-field-note--attached': attached }]" :role="tone === 'err' ? 'alert' : undefined">
    <XIcon :name="tone === 'warn' ? 'warning' : 'alert'" :size="12" :stroke="1.4" class="x-field-note__icon" />
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
  font: 400 var(--fs-12) / var(--lh-text) var(--font-sans);
}

.x-field-note--attached {
  margin-top: -6px;
}

.x-field-note--warn {
  color: var(--warn-x);
}

.x-field-note--err {
  color: var(--err-t);
  font-weight: 500;
}

.x-field-note--info {
  color: var(--mut);
}

.x-field-note--info .x-field-note__icon {
  color: var(--dim);
}

.x-field-note__icon {
  flex: none;
  margin-top: 2px;
}

.x-field-note--warn .x-field-note__icon {
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
