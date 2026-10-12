<script setup lang="ts">
import XIcon from './XIcon.vue'
import type { IconName } from './icons'

/** Nothing to show yet: an icon in a dashed frame, a title and a hint. Fills and centres in its box. */
withDefaults(defineProps<{ title: string; hint?: string; icon?: IconName }>(), { icon: 'doc' })
</script>

<template>
  <div class="x-empty">
    <span class="x-empty__icon"><XIcon :name="icon" :size="16" :stroke="1.3" /></span>
    <span class="x-empty__title">{{ title }}</span>
    <span v-if="hint" class="x-empty__hint">{{ hint }}</span>
    <slot />
  </div>
</template>

<style scoped>
.x-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  height: 100%;
  min-height: 200px;
  padding: 0 16px;
  text-align: center;
}

.x-empty__icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  border: 1px dashed var(--track);
  border-radius: var(--r-12);
  color: var(--faint);
}

.x-empty__title {
  color: var(--mut);
  font: 500 var(--fs-13) / var(--lh-tight) var(--font-sans);
}

.x-empty__hint {
  color: var(--dim);
  font: 400 var(--fs-12) / var(--lh-text) var(--font-sans);
  text-wrap: pretty;
}
</style>
