<script setup lang="ts">
import { XIcon, type IconName } from '../../kit'

/** A row that opens a window: icon, label and a chevron. */
defineProps<{ icon: IconName; label: string }>()
</script>

<template>
  <button type="button" class="dash-link">
    <XIcon :name="icon" :size="16" :stroke="1.4" />
    <span class="dash-link__label">{{ label }}</span>
    <XIcon name="chevron" :size="12" :stroke="1.6" class="dash-link__chevron" />
  </button>
</template>

<style scoped>
.dash-link {
  display: flex;
  align-items: center;
  gap: 12px;
  width: calc(100% + 24px);
  height: var(--ctl-h);
  margin: 0 -12px;
  padding: 0 12px;
  border: 0;
  border-radius: var(--r-10);
  background: transparent;
  color: var(--fg2);
  font: 400 var(--ctl-fs) / 1 var(--font-sans);
  text-align: left;
  outline: none;
  transition: background-color var(--t-fast), color var(--t-fast);
}

.dash-link:hover {
  background: var(--hov);
  color: var(--fg);
}

.dash-link:focus-visible {
  box-shadow: inset 0 0 0 2px var(--acc-line);
}

.dash-link__label {
  flex: 1;
  min-width: 0;
}

.dash-link__chevron {
  color: var(--faint);
}
</style>
