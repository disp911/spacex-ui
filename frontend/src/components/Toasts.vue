<script setup lang="ts">
import { toasts } from '../lib/toast'
import Icon from './Icon.vue'
</script>

<template>
  <div class="toasts" aria-live="polite">
    <TransitionGroup name="toast">
      <div v-for="m in toasts" :key="m.id" class="toast" :class="m.kind">
        <Icon :name="m.kind === 'success' ? 'check' : m.kind === 'error' ? 'alert' : 'info'" :size="15" />
        <span>{{ m.text }}</span>
      </div>
    </TransitionGroup>
  </div>
</template>

<style scoped>
.toasts {
  position: fixed;
  top: 16px;
  left: 50%;
  transform: translateX(-50%);
  z-index: 1000;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  pointer-events: none;
  width: max-content;
  max-width: calc(100vw - 32px);
}
.toast {
  display: flex;
  align-items: center;
  gap: 9px;
  padding: 10px 14px;
  border-radius: 12px;
  background: var(--popover);
  border: 1px solid var(--line-card);
  box-shadow: var(--shadow-pop);
  backdrop-filter: blur(10px);
  font: 500 12.5px var(--font-sans);
  color: var(--text);
  pointer-events: auto;
}
.toast.success svg {
  color: var(--accent);
}
.toast.error {
  border-color: rgba(var(--danger-rgb), 0.32);
}
.toast.error svg {
  color: var(--danger-text);
}
.toast.info svg {
  color: var(--info);
}
.toast-enter-active {
  animation: sheetUp 0.16s ease-out;
}
.toast-leave-active {
  transition: opacity 0.2s ease;
}
.toast-leave-to {
  opacity: 0;
}
</style>
