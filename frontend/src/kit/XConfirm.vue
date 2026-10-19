<script setup lang="ts">
import XButton from './XButton.vue'
import XIcon from './XIcon.vue'
import XModal from './XModal.vue'
import type { IconName } from './icons'

/**
 * Asks before an action that changes a lot: a warning icon, a title, the
 * text (the slot), Cancel and the action. While busy the action shows a
 * spinner and its busy label, and the window cannot be closed. icon names the
 * action in place of the warning sign, e.g. refresh for a restart. danger is
 * for what cannot be undone, e.g. deleting or resetting traffic: the sign and
 * the action turn red.
 */
const props = withDefaults(defineProps<{ title: string; cta: string; busy?: boolean; busyLabel?: string; icon?: IconName; danger?: boolean }>(), {
  icon: 'alert',
})
const emit = defineEmits<{ confirm: [] }>()
const open = defineModel<boolean>('open', { default: false })

function run() {
  if (!props.busy) emit('confirm')
}
</script>

<template>
  <XModal v-model:open="open" :title="title" size="xs" headless cancel :dismissible="!busy">
    <div class="x-confirm">
      <span :class="['x-confirm__icon', { 'x-confirm__icon--danger': danger }]"><XIcon :name="icon" :size="16" :stroke="1.4" /></span>
      <div class="x-confirm__body">
        <h2 class="x-confirm__title">{{ title }}</h2>
        <div class="x-confirm__text"><slot /></div>
      </div>
    </div>
    <template #footer>
      <XButton :variant="danger ? 'danger' : 'primary'" :loading="busy" @click="run">{{ busy && busyLabel ? busyLabel : cta }}</XButton>
    </template>
  </XModal>
</template>

<style scoped>
.x-confirm {
  display: flex;
  gap: 12px;
}

.x-confirm__icon {
  display: flex;
  flex: none;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border: 1px solid var(--warn-line);
  border-radius: var(--r-10);
  background: var(--warn-bg);
  color: var(--warn);
}

.x-confirm__icon--danger {
  border-color: var(--err-line);
  background: var(--err-bg);
  color: var(--err-t);
}

.x-confirm__body {
  flex: 1;
  min-width: 0;
}

.x-confirm__title {
  margin: 0 0 4px;
  color: var(--fg);
  font: 600 var(--fs-14) / var(--lh-tight) var(--font-sans);
  letter-spacing: -.2px;
  text-wrap: pretty;
}

.x-confirm__text {
  color: var(--mut);
  font: 400 var(--fs-12) / var(--lh-text) var(--font-sans);
  text-wrap: pretty;
}

@media (max-width: 767px) {
  .x-confirm__title {
    margin-bottom: 6px;
    font-size: var(--fs-16);
  }

  .x-confirm__text {
    font-size: var(--fs-13);
  }
}
</style>
