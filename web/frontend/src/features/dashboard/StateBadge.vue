<script setup lang="ts">
import { computed } from 'vue'
import { t } from '@/i18n'

// Pulsing dot + large state label ("Запущен").
const props = defineProps<{ state?: string }>()

const tone = computed(() => {
  switch (props.state) {
    case 'running':
      return 'ok'
    case 'stop':
      return 'warn'
    case 'error':
      return 'err'
    default:
      return 'unknown'
  }
})
const label = computed(() => t(`dash.state.${props.state === 'running' || props.state === 'stop' || props.state === 'error' ? props.state : 'unknown'}`))
</script>

<template>
  <div class="state" :class="`is-${tone}`">
    <span class="state__dot" />
    <span class="state__label">{{ label }}</span>
  </div>
</template>

<style scoped>
.state {
  --c: var(--text-4);
  display: flex;
  align-items: center;
  gap: 10px;
  min-width: 0;
  color: var(--c);
}
.state.is-ok { --c: var(--accent); }
.state.is-warn { --c: var(--warning); }
.state.is-err { --c: var(--danger); }
.state__dot {
  width: 9px;
  height: 9px;
  flex: none;
  border-radius: 50%;
  background: var(--c);
  animation: spx-pulse 2.4s infinite;
}
.state.is-unknown .state__dot {
  animation: none;
}
.state__label {
  font: var(--fw-semibold) 18px / 1.2 var(--font-sans);
  letter-spacing: -0.4px;
  white-space: nowrap;
}
</style>
