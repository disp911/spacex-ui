<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import XIcon from './XIcon.vue'

/** The title over a group of fields or a list, e.g. Final Rules, with a dashed + at its right that adds an item. */
const props = defineProps<{ title: string; add?: boolean }>()
const emit = defineEmits<{ add: [] }>()
const { t } = useI18n()
</script>

<template>
  <div class="x-list-head">
    <span class="x-list-head__title">{{ title }}</span>
    <button v-if="add" type="button" class="x-list-head__add" :aria-label="t('kit.addRow', { list: props.title })" @click="emit('add')">
      <XIcon name="plus" :size="12" :stroke="1.8" />
    </button>
  </div>
</template>

<style scoped>
.x-list-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  min-height: 24px;
}

.x-list-head__title {
  color: var(--fg2);
  font: 600 var(--fs-13) / var(--lh-tight) var(--font-sans);
}

.x-list-head__add {
  display: flex;
  flex: none;
  align-items: center;
  justify-content: center;
  width: 22px;
  height: 22px;
  padding: 0;
  border: 1px dashed var(--acc);
  border-radius: var(--r-6);
  background: transparent;
  color: var(--acc);
  outline: none;
  transition: background-color var(--t-fast);
}

.x-list-head__add:hover {
  background: var(--acc-soft);
}

.x-list-head__add:focus-visible {
  box-shadow: 0 0 0 3px var(--acc-ring);
}
</style>
