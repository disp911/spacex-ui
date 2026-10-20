<script setup lang="ts">
/**
 * One fact of a details window: its name on the left and its values on the
 * right (the slot), chips that wrap, e.g. a client's ID with a copy chip
 * beside it. The name sits level with the first line of chips.
 */
defineProps<{ label: string }>()
</script>

<template>
  <div class="x-fact-row">
    <span class="x-fact-row__label">{{ label }}</span>
    <div class="x-fact-row__value"><slot /></div>
  </div>
</template>

<style scoped>
.x-fact-row {
  display: flex;
  align-items: flex-start;
  gap: var(--sp-12);
  min-width: 0;
  min-height: 24px;
}

.x-fact-row__label {
  flex: none;
  width: 150px;
  padding-top: calc((var(--chip-h) - var(--lh-text) * 1em) / 2);
  color: var(--mut);
  font: 400 var(--fs-12) / var(--lh-text) var(--font-sans);
}

.x-fact-row__value {
  display: flex;
  flex: 1;
  flex-wrap: wrap;
  align-items: center;
  gap: var(--sp-6);
  min-width: 0;
}

@media (max-width: 767px) {
  .x-fact-row {
    min-height: 28px;
  }

  .x-fact-row__label {
    width: 112px;
  }
}
</style>
