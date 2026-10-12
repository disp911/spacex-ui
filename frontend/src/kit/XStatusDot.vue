<script setup lang="ts">
import type { Tone } from './types'

withDefaults(defineProps<{ tone?: Tone; pulse?: boolean; label?: string }>(), { tone: 'ok' })
</script>

<template>
  <span class="x-dot" :role="label ? 'img' : undefined" :aria-label="label" :title="label">
    <span :class="['x-dot__core', `x-dot__core--${tone}`, { 'is-pulse': pulse }]" />
  </span>
</template>

<style scoped>
.x-dot {
  display: inline-flex;
  flex: none;
  align-items: center;
  justify-content: center;
  width: 16px;
  height: 16px;
}

.x-dot__core {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--dot);
}

.x-dot__core.is-pulse {
  --pulse-color: var(--dot);
  animation: x-pulse 1.6s ease-out infinite;
}

.x-dot__core--ok { --dot: var(--ok); }
.x-dot__core--warn { --dot: var(--warn); }
.x-dot__core--err { --dot: var(--err); }
.x-dot__core--unl { --dot: var(--unl); }
.x-dot__core--on { --dot: var(--on); }
.x-dot__core--off { --dot: var(--off); }
.x-dot__core--neutral { --dot: var(--bar-unl); }
</style>
