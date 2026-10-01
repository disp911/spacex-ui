<script setup lang="ts">
defineProps<{ size?: number }>()
</script>

<template>
  <svg class="x-spinner" :width="size ?? 16" :height="size ?? 16" viewBox="0 0 18 18" fill="none" aria-hidden="true">
    <circle cx="9" cy="9" r="7.5" stroke="currentColor" stroke-width="2" opacity=".3" />
    <path d="M9 1.5a7.5 7.5 0 0 1 7.5 7.5" stroke="currentColor" stroke-width="2" stroke-linecap="round" />
  </svg>
</template>

<style scoped>
.x-spinner {
  display: block;
  flex: none;
  transform-origin: 50% 50%;
  animation: spx-spin 0.7s linear infinite;
}
</style>
