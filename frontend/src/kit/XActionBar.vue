<script setup lang="ts">
/**
 * Phones: a bar pinned to the bottom of a page with a line of state above the
 * page's main buttons (the slot), e.g. Save and Restart; they share the row.
 */
</script>

<template>
  <div class="x-action-bar">
    <div class="x-action-bar__state"><slot name="state" /></div>
    <div class="x-action-bar__buttons"><slot /></div>
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
  padding: 10px 14px calc(14px + env(safe-area-inset-bottom));
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
  gap: 8px;
}

.x-action-bar__buttons > :deep(*) {
  flex: 1;
  min-width: 0;
}
</style>
