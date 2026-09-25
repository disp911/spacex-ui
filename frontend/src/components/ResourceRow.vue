<script setup lang="ts">
import { computed } from 'vue'
import { pct } from '../lib/format'

const props = defineProps<{ label: string; detail: string; percent: number }>()

// Green below 80%, amber to 90%, red above.
const color = computed(() =>
  props.percent >= 90 ? 'var(--danger)' : props.percent >= 80 ? 'var(--warn-state)' : 'var(--accent)',
)
</script>

<template>
  <div class="res">
    <div class="line">
      <span class="name">{{ label }}</span>
      <span class="detail">{{ detail }}</span>
      <span class="value" :style="{ color }">{{ pct(percent) }}<small>%</small></span>
    </div>
    <div class="track"><div class="bar" :style="{ width: `${percent}%`, background: color }" /></div>
  </div>
</template>

<style scoped>
.res {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.line {
  display: flex;
  align-items: baseline;
  gap: 9px;
  min-width: 0;
}
.name {
  font: 500 11.5px var(--font-mono);
  letter-spacing: 0.06em;
  color: var(--text-2);
  white-space: nowrap;
  text-transform: uppercase;
}
.detail {
  font: 400 10.5px var(--font-mono);
  color: var(--muted);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  min-width: 0;
}
.value {
  margin-left: auto;
  font: 600 16px var(--font-mono);
  white-space: nowrap;
}
.value small {
  font-size: 10.5px;
  font-weight: 500;
  color: var(--muted);
}
.track {
  height: 6px;
  border-radius: 3px;
  background: var(--track);
  overflow: hidden;
}
.bar {
  height: 100%;
  border-radius: 3px;
  transition: width 0.6s ease;
}
</style>
