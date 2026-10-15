<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref } from 'vue'
import XIcon from './XIcon.vue'

/**
 * A round-cornered button at the bottom right that scrolls a long page back to
 * its top; it shows once the page is scrolled past 200px. On phones it rides
 * above the bar of buttons, as toasts do.
 */
defineProps<{ label: string }>()

const shown = ref(false)
const check = () => (shown.value = window.scrollY > 200)
onMounted(() => {
  check()
  window.addEventListener('scroll', check, { passive: true })
})
onBeforeUnmount(() => window.removeEventListener('scroll', check))

function toTop() {
  const smooth = !matchMedia('(prefers-reduced-motion: reduce)').matches
  window.scrollTo({ top: 0, behavior: smooth ? 'smooth' : 'auto' })
}
</script>

<template>
  <Transition name="x-back-to-top">
    <button v-if="shown" type="button" class="x-back-to-top" :aria-label="label" :title="label" @click="toTop">
      <XIcon name="arrowUp" :size="16" :stroke="1.6" />
    </button>
  </Transition>
</template>

<style scoped>
.x-back-to-top {
  position: fixed;
  right: 24px;
  bottom: 24px;
  z-index: 30;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  border: 1px solid var(--line);
  border-radius: var(--r-12);
  background: var(--dlg);
  color: var(--fg2);
  box-shadow: var(--shadow-pop);
  cursor: pointer;
  outline: none;
  transition: color var(--t-fast), border-color var(--t-fast);
}

.x-back-to-top:hover {
  border-color: var(--field-hover);
  color: var(--fg);
}

.x-back-to-top:focus-visible {
  box-shadow: 0 0 0 2px var(--bg), 0 0 0 4px var(--acc);
}

.x-back-to-top-enter-active,
.x-back-to-top-leave-active {
  transition: opacity var(--t-med), transform var(--t-med) var(--ease);
}

.x-back-to-top-enter-from,
.x-back-to-top-leave-to {
  opacity: 0;
  transform: translateY(8px);
}

@media (max-width: 767px) {
  /* 12px above the bar of buttons when a page has one, where toasts stand 20px above it. */
  .x-back-to-top {
    right: 14px;
    bottom: calc(var(--toast-bottom, calc(12px + env(safe-area-inset-bottom))) + 12px);
  }
}
</style>
