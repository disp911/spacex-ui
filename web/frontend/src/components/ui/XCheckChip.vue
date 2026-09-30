<script setup lang="ts">
import XIcon from './XIcon.vue'

// Toggle chip with a checkbox square, tinted when on (tone decides the tint).
const on = defineModel<boolean>({ default: false })
withDefaults(defineProps<{ label: string; tone?: 'accent' | 'danger' | 'info' | 'neutral'; size?: 'sm' | 'lg' }>(), {
  tone: 'accent',
  size: 'sm',
})
</script>

<template>
  <button
    type="button"
    role="checkbox"
    :aria-checked="on"
    class="x-chip"
    :class="[`x-chip--${tone}`, `x-chip--${size}`, { 'is-on': on }]"
    @click="on = !on"
  >
    <span class="x-chip__box"><XIcon name="check" :size="10" :stroke-width="2.2" /></span>
    <span class="x-chip__label">{{ label }}</span>
  </button>
</template>

<style scoped>
.x-chip {
  --tone: var(--accent);
  --tone-bg: var(--accent-soft);
  --tone-border: var(--accent-border);
  flex: none;
  display: inline-flex;
  align-items: center;
  gap: 8px;
  height: var(--h-sm);
  padding: 0 11px 0 9px;
  border-radius: var(--r-sm);
  border: 1px solid var(--border-control);
  background: var(--control);
  white-space: nowrap;
  transition: border-color var(--dur-fast) ease, background var(--dur-fast) ease;
}
.x-chip--lg {
  height: var(--h-lg);
  width: 100%;
  padding: 0 11px;
  gap: 11px;
  border-radius: var(--r-md);
}
.x-chip--danger { --tone: var(--danger); --tone-bg: var(--danger-soft); --tone-border: var(--danger-border); }
.x-chip--info { --tone: var(--info); --tone-bg: var(--info-soft); --tone-border: var(--info-border); }
.x-chip--neutral { --tone: var(--text-2); --tone-bg: var(--control); --tone-border: var(--border-control); }
.x-chip:hover {
  border-color: var(--accent-focus);
}
.x-chip.is-on {
  background: var(--tone-bg);
  border-color: var(--tone-border);
}
.x-chip__box {
  width: 15px;
  height: 15px;
  flex: none;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--r-xs);
  border: 1px solid var(--scrollbar);
  color: transparent;
}
.x-chip.is-on .x-chip__box {
  background: var(--tone);
  border-color: var(--tone);
  color: var(--bg);
}
.x-chip__label {
  font: var(--fw-medium) var(--fs-sm) var(--font-sans);
  color: var(--text-3);
}
.x-chip--lg .x-chip__label {
  font-size: var(--fs-md);
  flex: 1;
  text-align: left;
}
.x-chip.is-on .x-chip__label {
  color: var(--tone);
}
</style>
