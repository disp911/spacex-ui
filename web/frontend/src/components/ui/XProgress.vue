<script setup lang="ts">
import { computed } from 'vue'

// Labeled usage bar: "ЦП  14 cores · 3.60 GHz ........ 14.2%".
const props = withDefaults(defineProps<{ label: string; meta?: string; value: number; dangerAt?: number }>(), { dangerAt: 90 })
const danger = computed(() => props.value >= props.dangerAt)
const shown = computed(() => (props.value >= 10 || props.value === 0 ? props.value.toFixed(1) : props.value.toFixed(2)))
</script>

<template>
  <div class="x-progress" :class="{ 'is-danger': danger }">
    <div class="x-progress__row">
      <span class="x-progress__label">{{ label }}</span>
      <span v-if="meta" class="x-progress__meta">{{ meta }}</span>
      <span class="x-progress__value">{{ shown }}<small>%</small></span>
    </div>
    <div class="x-progress__track" role="progressbar" :aria-valuenow="value" aria-valuemin="0" aria-valuemax="100" :aria-label="label">
      <div class="x-progress__bar" :style="{ width: Math.min(100, Math.max(0, value)) + '%' }" />
    </div>
  </div>
</template>

<style scoped>
.x-progress {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.x-progress__row {
  display: flex;
  align-items: baseline;
  gap: 9px;
  min-width: 0;
}
.x-progress__label {
  font: var(--fw-medium) var(--fs-sm) var(--font-mono);
  letter-spacing: 0.06em;
  text-transform: uppercase;
  color: var(--text-2);
  white-space: nowrap;
}
.x-progress__meta {
  font: var(--fw-regular) var(--fs-xs) var(--font-mono);
  color: var(--text-4);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  min-width: 0;
}
.x-progress__value {
  margin-left: auto;
  font: var(--fw-semibold) var(--fs-xl) var(--font-mono);
  color: var(--accent);
  white-space: nowrap;
}
.x-progress__value small {
  font-size: var(--fs-xs);
  font-weight: var(--fw-medium);
  color: var(--text-4);
}
.x-progress__track {
  height: 6px;
  border-radius: 3px;
  background: var(--track);
  overflow: hidden;
}
.x-progress__bar {
  height: 100%;
  border-radius: 3px;
  background: var(--accent);
  transition: width 0.6s var(--ease-out);
}
.is-danger .x-progress__value { color: var(--danger); }
.is-danger .x-progress__bar { background: var(--danger); }
</style>
