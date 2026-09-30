<script setup lang="ts">
import XIcon from '@/components/ui/XIcon.vue'
import { useTheme } from '@/composables/useTheme'
import { t } from '@/i18n'

const { theme, toggle } = useTheme()
</script>

<template>
  <button type="button" class="theme-toggle" :aria-label="t('common.toggleTheme')" :title="t('common.toggleTheme')" @click="toggle">
    <span class="theme-toggle__icons">
      <XIcon name="moon" :class="{ 'is-shown': theme === 'dark' }" />
      <XIcon name="sun" :class="{ 'is-shown': theme === 'light' }" />
    </span>
  </button>
</template>

<style scoped>
.theme-toggle {
  flex: none;
  width: var(--h-md);
  height: var(--h-md);
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--r-md);
  border: 1px solid var(--border-control);
  background: var(--field);
  color: var(--text-2);
  transition: border-color var(--dur-base) ease, background var(--dur-base) ease, box-shadow var(--dur-base) ease, color var(--dur-fast) ease;
}
.theme-toggle:hover {
  border-color: var(--accent-focus);
  background: var(--accent-soft-hover);
  box-shadow: 0 0 0 3px var(--accent-ring), 0 0 14px var(--accent-glow);
  color: var(--text);
}
.theme-toggle__icons {
  position: relative;
  width: 16px;
  height: 16px;
}
.theme-toggle__icons :deep(svg) {
  position: absolute;
  inset: 0;
  opacity: 0;
  transform: rotate(-90deg) scale(0.55);
  transition: transform 0.42s cubic-bezier(0.34, 1.32, 0.5, 1), opacity 0.26s ease;
}
.theme-toggle__icons :deep(svg.is-shown) {
  opacity: 1;
  transform: none;
}
</style>
