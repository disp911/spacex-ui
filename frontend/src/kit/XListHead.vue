<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import XAddButton from './XAddButton.vue'

/** The title over a group of fields or a list, e.g. Final Rules, with a dashed + at its right that adds an item. */
const props = defineProps<{ title: string; add?: boolean }>()
const emit = defineEmits<{ add: [] }>()
const { t } = useI18n()
</script>

<template>
  <div class="x-list-head">
    <span class="x-list-head__title">{{ title }}</span>
    <XAddButton v-if="add" :label="t('kit.addRow', { list: props.title })" @click="emit('add')" />
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
</style>
