<script setup lang="ts">
import XIcon from './XIcon.vue'

// One collapsible section: tile-like header with caret, title and meta.
defineProps<{ title: string; meta?: string; open: boolean }>()
defineEmits<{ toggle: [] }>()
</script>

<template>
  <div class="x-acc" :class="{ 'is-open': open }">
    <button type="button" class="x-acc__head" :aria-expanded="open" @click="$emit('toggle')">
      <XIcon name="chevron-right" :size="13" :stroke-width="1.6" class="x-acc__caret" />
      <span class="x-acc__title">{{ title }}</span>
      <span v-if="meta" class="x-acc__meta">{{ meta }}</span>
    </button>
    <div v-if="open" class="x-acc__body"><slot /></div>
  </div>
</template>

<style scoped>
.x-acc__head {
  width: 100%;
  display: flex;
  align-items: center;
  gap: 11px;
  height: var(--h-lg);
  padding: 0 13px;
  border-radius: var(--r-md);
  border: 1px solid var(--border-field);
  background: var(--field);
  color: var(--text-2);
  text-align: left;
  transition: background var(--dur-fast) ease, border-color var(--dur-fast) ease;
}
.x-acc__head:hover {
  background: var(--control-hover);
}
.is-open .x-acc__head {
  background: var(--control);
  border-color: var(--accent-border);
  color: var(--text);
}
.x-acc__caret {
  color: var(--text-4);
  transition: transform var(--dur-fast) ease;
}
.is-open .x-acc__caret {
  color: var(--accent);
  transform: rotate(90deg);
}
.x-acc__title {
  flex: 1;
  min-width: 0;
  font: var(--fw-medium) var(--fs-md) var(--font-sans);
}
.x-acc__meta {
  flex: none;
  font: var(--fw-medium) var(--fs-sm) var(--font-mono);
  color: var(--text-3);
  white-space: nowrap;
}
.x-acc__body {
  margin-top: 8px;
  display: flex;
  flex-direction: column;
  gap: 10px;
  animation: spx-pop-in 0.18s var(--ease-out);
}
</style>
