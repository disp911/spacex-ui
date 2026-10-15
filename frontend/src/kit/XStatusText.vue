<script setup lang="ts">
/** A dot and a short line of state, e.g. unsaved changes (mod) or errors (err) on a page with a Save button. */
defineProps<{ tone: 'mod' | 'err' | 'ok' }>()
</script>

<template>
  <span :class="['x-status-text', `x-status-text--${tone}`]" role="status">
    <span class="x-status-text__dot" aria-hidden="true" />
    <span class="x-status-text__label"><slot /></span>
  </span>
</template>

<style scoped>
.x-status-text {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
  font: 500 var(--fs-12) / var(--lh-tight) var(--font-sans);
}

.x-status-text__dot {
  flex: none;
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--tone);
}

.x-status-text__label {
  min-width: 0;
  overflow: hidden;
  color: var(--tone-t);
  white-space: nowrap;
  text-overflow: ellipsis;
}

.x-status-text--mod { --tone: var(--mod); --tone-t: var(--mod); }
.x-status-text--err { --tone: var(--err); --tone-t: var(--err-t); }
.x-status-text--ok { --tone: var(--ok); --tone-t: var(--ok-t); }
</style>
