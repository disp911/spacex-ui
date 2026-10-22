<script setup lang="ts">
/**
 * A grey bar in place of a figure or a line of text that is still loading, as
 * tall as the text line it stands in, so nothing moves when the value comes.
 */
withDefaults(defineProps<{ width?: number | string }>(), { width: 64 })
</script>

<template>
  <span class="x-skeleton-text x-skeleton" aria-hidden="true">
    <span class="x-skeleton__bar" :style="{ width: typeof width === 'number' ? `${width}px` : width }" />
  </span>
</template>

<style scoped>
.x-skeleton-text {
  display: flex;
  align-items: center;
  max-width: 100%;
  height: 1lh;
}
</style>
