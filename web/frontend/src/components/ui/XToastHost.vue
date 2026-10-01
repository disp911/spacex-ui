<script setup lang="ts">
import XIcon from './XIcon.vue'
import { dismissToast, toasts } from '@/composables/useToast'
</script>

<template>
  <Teleport to="body">
    <TransitionGroup tag="div" name="x-toast" class="x-toasts" aria-live="polite">
      <div v-for="m in toasts" :key="m.id" class="x-toast" :class="`is-${m.tone}`" @click="dismissToast(m.id)">
        <XIcon :name="m.tone === 'success' ? 'check' : m.tone === 'info' ? 'info' : 'alert'" :size="15" />
        <span>{{ m.text }}</span>
      </div>
    </TransitionGroup>
  </Teleport>
</template>

<style scoped>
.x-toasts {
  position: fixed;
  top: 16px;
  left: 50%;
  transform: translateX(-50%);
  z-index: 100;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  width: min(420px, calc(100vw - 32px));
  pointer-events: none;
}
.x-toast {
  pointer-events: auto;
  display: flex;
  align-items: flex-start;
  gap: 10px;
  max-width: 100%;
  padding: 10px 14px;
  border-radius: var(--r-md);
  border: 1px solid var(--border);
  background: var(--surface-raised);
  box-shadow: var(--shadow-pop);
  backdrop-filter: blur(10px);
  font: var(--fw-medium) var(--fs-md) / 1.4 var(--font-sans);
  color: var(--text);
  cursor: pointer;
}
.x-toast :deep(svg) {
  margin-top: 2px;
}
.x-toast.is-success :deep(svg) { color: var(--accent); }
.x-toast.is-danger { border-color: var(--danger-border); }
.x-toast.is-danger :deep(svg) { color: var(--danger); }
.x-toast.is-warning :deep(svg) { color: var(--warning); }
.x-toast.is-info :deep(svg) { color: var(--info); }

.x-toast-enter-active,
.x-toast-leave-active {
  transition: opacity 0.2s ease, transform 0.2s var(--ease-out);
}
.x-toast-enter-from,
.x-toast-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}
</style>
