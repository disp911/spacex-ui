<script setup lang="ts">
import XTooltip from './XTooltip.vue'
import type { Tone } from './types'

/**
 * Status dot. pulse — a ring for live things (clients online);
 * blink — the slow fade of a service state or an offer. sm is the 6px
 * dot beside short text; color overrides the tone, e.g. a log level. With
 * tip the dot sits in a 16px spot and the label, if any, shows in a tooltip
 * the keyboard reaches too, e.g. a client's state.
 */
withDefaults(
  defineProps<{ tone?: Tone; color?: string; size?: 'sm' | 'md'; pulse?: boolean; blink?: boolean; label?: string; tip?: boolean }>(),
  { tone: 'ok', size: 'md' },
)
</script>

<template>
  <XTooltip v-if="tip && label" :text="label" :delay="0">
    <span class="x-dot-spot" tabindex="0" role="img" :aria-label="label">
      <span
        :class="['x-dot', `x-dot--${size}`, `x-dot--${tone}`, { 'is-pulse': pulse, 'is-blink': blink }]"
        :style="color ? { '--dot': color } : undefined"
      />
    </span>
  </XTooltip>
  <span v-else-if="tip" class="x-dot-spot">
    <span
      :class="['x-dot', `x-dot--${size}`, `x-dot--${tone}`, { 'is-pulse': pulse, 'is-blink': blink }]"
      :style="color ? { '--dot': color } : undefined"
    />
  </span>
  <span
    v-else
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

.x-dot-spot {
  display: flex;
  flex: none;
  align-items: center;
  justify-content: center;
  width: 16px;
  height: 16px;
  border-radius: var(--r-8);
  outline: none;
}

.x-dot-spot:focus-visible {
  box-shadow: 0 0 0 2px var(--acc-line);
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
