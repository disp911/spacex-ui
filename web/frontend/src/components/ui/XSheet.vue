<script setup lang="ts">
import { toRef } from 'vue'
import { useEscape } from '@/composables/useMedia'
import { useScrollLock } from '@/composables/useScrollLock'

// Mobile bottom sheet with a grab handle and optional caps title.
const open = defineModel<boolean>('open', { default: false })
defineProps<{ title?: string }>()

useScrollLock(toRef(open))
useEscape(() => (open.value = false))
</script>

<template>
  <Teleport to="body">
    <Transition name="x-sheet">
      <div v-if="open" class="x-sheet" @mousedown.self="open = false">
        <div class="x-sheet__panel" role="dialog" aria-modal="true" :aria-label="title">
          <div class="x-sheet__handle" />
          <div v-if="title" class="caps x-sheet__title">{{ title }}</div>
          <div class="x-sheet__body no-scrollbar"><slot /></div>
          <div v-if="$slots.footer" class="x-sheet__foot"><slot name="footer" /></div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.x-sheet {
  position: fixed;
  inset: 0;
  z-index: 60;
  display: flex;
  align-items: flex-end;
  background: var(--scrim);
  backdrop-filter: blur(3px) saturate(0.92);
}
.x-sheet__panel {
  width: 100%;
  max-height: 85vh;
  display: flex;
  flex-direction: column;
  gap: 9px;
  padding: 10px 12px calc(18px + env(safe-area-inset-bottom));
  background: var(--surface-raised);
  border-top: 1px solid var(--border);
  border-radius: var(--r-sheet) var(--r-sheet) 0 0;
  box-shadow: var(--shadow-sheet);
  color: var(--text);
}
.x-sheet__handle {
  align-self: center;
  width: 38px;
  height: 4px;
  border-radius: 2px;
  background: var(--scrollbar);
}
.x-sheet__title {
  padding: 0 5px;
}
.x-sheet__body {
  min-height: 0;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 6px;
}
.x-sheet__foot {
  display: flex;
  gap: 9px;
  margin-top: 3px;
}
.x-sheet__foot :deep(.x-btn) {
  flex: 1;
}

.x-sheet-enter-active,
.x-sheet-leave-active {
  transition: opacity 0.22s ease;
}
.x-sheet-enter-active .x-sheet__panel {
  transition: transform 0.34s var(--ease-out);
}
.x-sheet-leave-active .x-sheet__panel {
  transition: transform 0.22s cubic-bezier(0.4, 0, 1, 1);
}
.x-sheet-enter-from,
.x-sheet-leave-to {
  opacity: 0;
}
.x-sheet-enter-from .x-sheet__panel,
.x-sheet-leave-to .x-sheet__panel {
  transform: translateY(100%);
}
</style>
