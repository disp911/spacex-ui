<script setup lang="ts">
import { computed } from 'vue'
import XIcon from './XIcon.vue'
import type { IconName } from './icons'
import { usePhone } from './media'

/** A figure of a summary card: an icon in a small frame, the label (two lines at most) and the value, or the slot in its place, e.g. counts. Compact as on a phone, which it is there unless told otherwise. */
const props = withDefaults(defineProps<{ icon: IconName; label: string; value?: string; compact?: boolean }>(), { value: undefined, compact: undefined })
const phone = usePhone()
const compact = computed(() => props.compact ?? phone.value)
</script>

<template>
  <div :class="['x-sum', { 'x-sum--compact': compact }]">
    <span class="x-sum__icon"><XIcon :name="icon" :size="16" /></span>
    <div class="x-sum__body">
      <div class="x-sum__label">{{ label }}</div>
      <slot><div class="x-sum__value" :title="value">{{ value }}</div></slot>
    </div>
  </div>
</template>

<style scoped>
.x-sum {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  min-width: 0;
}

.x-sum__icon {
  display: flex;
  flex: none;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border: 1px solid var(--div);
  border-radius: var(--r-8);
  background: var(--nest);
  color: var(--mut);
}

.x-sum__body {
  display: flex;
  flex: 1;
  flex-direction: column;
  gap: 4px;
  min-width: 0;
}

.x-sum__label {
  display: -webkit-box;
  min-height: 32px;
  overflow: hidden;
  color: var(--dim);
  font: 400 var(--fs-12) / 1.35 var(--font-sans);
  text-wrap: pretty;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
}

.x-sum__value {
  overflow: hidden;
  color: var(--fg);
  font: 500 var(--fs-14) / 20px var(--font-mono);
  letter-spacing: -.2px;
  white-space: nowrap;
  text-overflow: ellipsis;
}

/* Compact ones stand one a row, so a label needs no room for a second line; the square grows with its icon. */
.x-sum--compact .x-sum__icon {
  width: 40px;
  height: 40px;
  border-radius: var(--r-10);
}

.x-sum--compact .x-sum__label {
  min-height: 0;
}
</style>
