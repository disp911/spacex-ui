<script setup lang="ts">
withDefaults(defineProps<{ size?: number }>(), { size: 14 })
</script>

<template>
  <svg class="x-spinner" :width="size" :height="size" viewBox="0 0 14 14" fill="none" aria-hidden="true">
    <circle cx="7" cy="7" r="5.6" stroke="currentColor" stroke-width="1.8" opacity=".3" />
    <path d="M7 1.4A5.6 5.6 0 0 1 12.6 7" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" />
  </svg>
</template>

<style scoped>
.x-spinner {
  display: block;
  flex: none;
  animation: x-spin .7s linear infinite;
}
</style>
