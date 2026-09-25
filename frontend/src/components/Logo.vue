<script setup lang="ts">
withDefaults(defineProps<{ size?: number; label?: boolean }>(), { size: 32, label: true })
</script>

<template>
  <div class="logo">
    <svg :width="size" :height="size" viewBox="0 0 32 32" aria-hidden="true">
      <rect width="32" height="32" rx="9" fill="var(--accent)" />
      <ellipse cx="16" cy="16" rx="11" ry="4.6" transform="rotate(-32 16 16)" fill="none" stroke="var(--on-accent)" stroke-width="1.7" opacity=".9" />
      <circle cx="16" cy="16" r="3.2" fill="var(--on-accent)" />
      <circle cx="25.2" cy="10.6" r="1.9" fill="var(--on-accent)" />
    </svg>
    <div v-if="label" class="name">SpaceX <span>panel</span></div>
  </div>
</template>

<style scoped>
.logo {
  display: flex;
  align-items: center;
  gap: 11px;
}
.name {
  font: 600 15.5px var(--font-sans);
  letter-spacing: -0.3px;
  white-space: nowrap;
}
.name span {
  font-weight: 450;
  color: var(--text-3b);
}
</style>
