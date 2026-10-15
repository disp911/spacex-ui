<script setup lang="ts">
/** A 2px line with a stripe running along it: work of unknown length, e.g. a restart. */
defineProps<{ label?: string }>()
</script>

<template>
  <div class="x-progress-line" role="progressbar" :aria-label="label"><span class="x-progress-line__run" /></div>
</template>

<style scoped>
.x-progress-line {
  position: relative;
  height: 2px;
  overflow: hidden;
  background: var(--acc-soft);
}

.x-progress-line__run {
  position: absolute;
  top: 0;
  bottom: 0;
  width: 35%;
  background: var(--acc);
  animation: x-progress-line 1.1s ease-in-out infinite;
}

@keyframes x-progress-line {
  from {
    left: -35%;
  }

  to {
    left: 100%;
  }
}

@media (prefers-reduced-motion: reduce) {
  .x-progress-line__run {
    left: 0;
    width: 100%;
    animation: none;
    opacity: .6;
  }
}
</style>
