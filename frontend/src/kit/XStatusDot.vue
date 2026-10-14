<script setup lang="ts">
import type { Tone } from './types'

/**
 * Status dot. pulse — a ring for live things (clients online);
 * blink — the slow fade of a service state or an offer. sm is the 6px
 * dot beside short text; color overrides the tone, e.g. a log level.
 */
withDefaults(
  defineProps<{ tone?: Tone; color?: string; size?: 'sm' | 'md'; pulse?: boolean; blink?: boolean; label?: string }>(),
  { tone: 'ok', size: 'md' },
)
</script>

<template>
  <span
    :class="['x-dot', `x-dot--${size}`, `x-dot--${tone}`, { 'is-pulse': pulse, 'is-blink': blink }]"
    :style="color ? { '--dot': color } : undefined"
    :role="label ? 'img' : undefined"
    :aria-label="label"
    :title="label"
  />
</template>

<style scoped>
.x-dot {
  display: inline-block;
  flex: none;
  border-radius: 50%;
  background: var(--dot);
}

.x-dot--md {
  width: 8px;
  height: 8px;
}

.x-dot--sm {
  width: 6px;
  height: 6px;
}

.x-dot.is-pulse {
  --pulse-color: var(--dot);
  animation: x-pulse 1.6s ease-out infinite;
}

.x-dot.is-blink {
  animation: x-blink 2.4s infinite;
}

.x-dot--ok { --dot: var(--ok); }
.x-dot--warn { --dot: var(--warn); }
.x-dot--err { --dot: var(--err); }
.x-dot--unl { --dot: var(--unl); }
.x-dot--on { --dot: var(--on); }
.x-dot--off { --dot: var(--off); }
.x-dot--neutral { --dot: var(--bar-unl); }
</style>
