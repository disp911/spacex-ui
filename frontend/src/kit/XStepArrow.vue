<script setup lang="ts">
import XIcon from './XIcon.vue'

/** A chevron between two buttons done in turn, e.g. Save, then Restart; accent once the first is done. */
defineProps<{ done?: boolean }>()
</script>

<template>
  <XIcon name="chevron" :size="14" :stroke="1.7" :class="['x-step-arrow', { 'is-done': done }]" aria-hidden="true" />
</template>

<style scoped>
.x-step-arrow {
  flex: none;
  color: var(--faint);
  transition: color var(--t-med);
}

.x-step-arrow.is-done {
  color: var(--acc);
}
</style>
