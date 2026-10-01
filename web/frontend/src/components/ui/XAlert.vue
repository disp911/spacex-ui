<script setup lang="ts">
import XIcon from './XIcon.vue'

withDefaults(defineProps<{ tone?: 'warning' | 'danger' | 'info'; title?: string; dismissible?: boolean; compact?: boolean }>(), {
  tone: 'warning',
})
defineEmits<{ dismiss: [] }>()
</script>

<template>
  <div class="x-alert" :class="[`x-alert--${tone}`, { 'x-alert--compact': compact }]" role="status">
    <XIcon :name="tone === 'info' ? 'info' : 'alert'" class="x-alert__icon" />
    <div class="x-alert__body">
      <div v-if="title" class="x-alert__title">{{ title }}</div>
      <div class="x-alert__text"><slot /></div>
    </div>
    <button v-if="dismissible" class="x-alert__close" type="button" aria-label="Close" @click="$emit('dismiss')">
      <XIcon name="close" :size="12" :stroke-width="1.6" />
    </button>
  </div>
</template>

<style scoped>
.x-alert {
  display: flex;
  gap: 10px;
  align-items: flex-start;
  padding: 11px 12px;
  border-radius: var(--r-md);
  border: 1px solid;
}
.x-alert__icon {
  margin-top: 1px;
}
.x-alert__body {
  flex: 1;
  min-width: 0;
}
.x-alert__title {
  font: var(--fw-semibold) var(--fs-sm) / 1.3 var(--font-sans);
  margin-bottom: 2px;
}
.x-alert__text {
  font: var(--fw-regular) var(--fs-sm) / 1.5 var(--font-sans);
  text-wrap: pretty;
}
.x-alert__close {
  flex: none;
  display: flex;
  padding: 2px;
  margin: -1px -2px 0 0;
  border: 0;
  background: none;
  color: inherit;
  opacity: 0.8;
}
.x-alert__close:hover {
  opacity: 1;
}

.x-alert--warning {
  background: var(--warning-soft);
  border-color: var(--warning-border);
  color: var(--warning-text);
}
.x-alert--warning .x-alert__icon { color: var(--warning); }
.x-alert--warning .x-alert__title { color: var(--warning-title); }

.x-alert--danger {
  background: var(--danger-soft);
  border-color: var(--danger-border);
  color: var(--danger-text);
}
.x-alert--danger .x-alert__icon { color: var(--danger-text); }
.x-alert--danger .x-alert__text { font-weight: var(--fw-medium); }

.x-alert--info {
  background: var(--field);
  border-color: var(--border-field);
  color: var(--text-3);
}
.x-alert--info .x-alert__icon { color: var(--info); }

.x-alert--compact {
  padding: 9px 12px;
  align-items: center;
}
.x-alert--compact .x-alert__icon { margin-top: 0; }
</style>
