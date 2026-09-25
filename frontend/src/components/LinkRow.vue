<script setup lang="ts">
import Icon, { type IconName } from './Icon.vue'

defineProps<{ icon: IconName; label: string }>()
defineEmits<{ click: [] }>()
</script>

<template>
  <button type="button" class="link-row" @click="$emit('click')">
    <Icon :name="icon" class="ico" />
    <span class="lbl">{{ label }}</span>
    <slot />
    <span class="chev">›</span>
  </button>
</template>

<style scoped>
.link-row {
  display: flex;
  align-items: center;
  gap: 11px;
  height: 34px;
  padding: 0 11px;
  margin: 0 -11px;
  border-radius: 10px;
  color: var(--text-2);
  font: 450 13px var(--font-sans);
  text-align: left;
  transition:
    background 0.18s ease,
    color 0.18s ease;
}
.link-row:hover {
  background: var(--row-hover-soft);
  color: var(--text);
}
.ico {
  color: var(--text-3b);
}
.lbl {
  flex: 1;
  min-width: 0;
}
.chev {
  color: var(--faint);
  font: 400 12px var(--font-mono);
}
@media (max-width: 760px) {
  .link-row {
    height: 44px;
    font-size: 14px;
  }
}
</style>
