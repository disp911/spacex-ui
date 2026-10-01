<script setup lang="ts">
withDefaults(defineProps<{ size?: number; wordmark?: boolean }>(), { size: 28, wordmark: true })
</script>

<template>
  <div class="x-logo">
    <svg :width="size" :height="size" viewBox="0 0 32 32" aria-hidden="true">
      <rect width="32" height="32" rx="9" fill="var(--logo-bg)" />
      <ellipse cx="16" cy="16" rx="11" ry="4.6" transform="rotate(-32 16 16)" fill="none" stroke="var(--logo-fg)" stroke-width="1.7" opacity=".9" />
      <circle cx="16" cy="16" r="3.2" fill="var(--logo-fg)" />
      <circle cx="25.2" cy="10.6" r="1.9" fill="var(--logo-fg)" />
    </svg>
    <span v-if="wordmark" class="x-logo__word">SpaceX <span>panel</span></span>
  </div>
</template>

<style scoped>
.x-logo {
  display: flex;
  align-items: center;
  gap: 10px;
  min-width: 0;
}
.x-logo svg {
  display: block;
  flex: none;
}
.x-logo__word {
  font: var(--fw-semibold) var(--fs-xl) var(--font-sans);
  letter-spacing: var(--tracking-tight);
  white-space: nowrap;
  color: var(--text);
}
.x-logo__word span {
  font-weight: var(--fw-regular);
  color: var(--text-3);
}
</style>
