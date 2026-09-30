<script setup lang="ts">
import XIcon, { type IconName } from './XIcon.vue'

// Row action inside a card: icon · label · chevron.
defineProps<{ icon: IconName; label: string }>()
defineEmits<{ click: [] }>()
</script>

<template>
  <button type="button" class="x-list-link" @click="$emit('click')">
    <XIcon :name="icon" class="x-list-link__icon" />
    <span class="x-list-link__label">{{ label }}</span>
    <XIcon name="chevron-right" :size="12" class="x-list-link__chev" />
  </button>
</template>

<style scoped>
.x-list-link {
  display: flex;
  align-items: center;
  gap: 11px;
  width: calc(100% + 22px);
  height: var(--h-md);
  margin: 0 -11px;
  padding: 0 11px;
  border: 0;
  border-radius: var(--r-md);
  background: transparent;
  color: var(--text-2);
  font: var(--fw-regular) var(--fs-md) var(--font-sans);
  text-align: left;
  transition: background var(--dur-fast) ease, color var(--dur-fast) ease;
}
.x-list-link:hover {
  background: var(--hover);
  color: var(--text);
}
.x-list-link__icon {
  color: var(--text-3);
}
.x-list-link__label {
  flex: 1;
  min-width: 0;
}
.x-list-link__chev {
  color: var(--text-5);
}
@media (max-width: 767px) {
  .x-list-link {
    height: 40px;
    font-size: var(--fs-lg);
  }
}
</style>
