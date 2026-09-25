<script setup lang="ts">
withDefaults(defineProps<{ size?: number; onAccent?: boolean }>(), { size: 16 })
</script>

<template>
  <svg class="spinner" :width="size" :height="size" viewBox="0 0 18 18" fill="none" aria-hidden="true">
    <circle cx="9" cy="9" r="7.5" :stroke="onAccent ? 'var(--on-accent)' : 'currentColor'" stroke-width="2" opacity=".28" />
    <path d="M9 1.5a7.5 7.5 0 0 1 7.5 7.5" :stroke="onAccent ? 'var(--on-accent)' : 'currentColor'" stroke-width="2" stroke-linecap="round" />
  </svg>
</template>

<style scoped>
.spinner {
  transform-origin: 50% 50%;
  animation: spin 0.7s linear infinite;
}
</style>
