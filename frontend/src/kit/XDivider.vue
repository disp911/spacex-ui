<script setup lang="ts">
/**
 * A line across with a word in its middle, splitting a window into parts, e.g.
 * Settings and Account info; without a word, just the line. spaced keeps 16px
 * above and below it in a stack of fields 12px apart, as between the parts of
 * the inbound window.
 */
defineProps<{ label?: string; spaced?: boolean }>()
</script>

<template>
  <div :class="['x-divider', { 'x-divider--spaced': spaced }]" role="separator" :aria-label="label">
    <span v-if="label" class="x-divider__label">{{ label }}</span>
  </div>
</template>

<style scoped>
.x-divider {
  display: flex;
  align-items: center;
  gap: var(--sp-12);
  color: var(--mut);
  font: 600 var(--fs-12) / var(--lh-tight) var(--font-sans);
}

.x-divider--spaced {
  margin-block: 4px;
}

.x-divider:not(:has(.x-divider__label))::after {
  display: none;
}

.x-divider::before,
.x-divider::after {
  content: '';
  flex: 1;
  height: 1px;
  background: var(--div);
}

.x-divider__label {
  flex: none;
  max-width: 70%;
  text-align: center;
}
</style>
