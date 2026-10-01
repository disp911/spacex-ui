<script setup lang="ts">
import XIcon, { type IconName } from './XIcon.vue'

withDefaults(defineProps<{ icon?: IconName; title: string; text?: string }>(), { icon: 'search' })
</script>

<template>
  <div class="x-empty">
    <div class="x-empty__icon"><XIcon :name="icon" :size="17" :stroke-width="1.3" /></div>
    <div class="x-empty__title">{{ title }}</div>
    <div v-if="text || $slots.default" class="x-empty__text"><slot>{{ text }}</slot></div>
  </div>
</template>

<style scoped>
.x-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 30px 20px 32px;
  text-align: center;
  height: 100%;
}
.x-empty__icon {
  width: 38px;
  height: 38px;
  margin-bottom: 4px;
  border-radius: var(--r-md);
  border: 1px dashed var(--scrollbar);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-5);
}
.x-empty__title {
  font: var(--fw-medium) var(--fs-md) var(--font-sans);
  color: var(--text-3);
}
.x-empty__text {
  font: var(--fw-regular) var(--fs-sm) / 1.5 var(--font-sans);
  color: var(--text-4);
  max-width: 320px;
  text-wrap: pretty;
}
</style>
