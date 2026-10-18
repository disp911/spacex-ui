<script setup lang="ts">
import { computed } from 'vue'
import type { Tone } from './types'

/** Used amount, a bar and the limit: client traffic or subscription period. */
const props = withDefaults(
  defineProps<{
    used: string
    limit: string
    /** Share used, 0..1; null when there is no limit. */
    ratio: number | null
    tone?: Tone
    /** Limit reached: the whole pill turns red. */
    exhausted?: boolean
    disabled?: boolean
    /** Used value in sans instead of mono, e.g. «через 9 дней». */
    textual?: boolean
    /** As wide as its place, e.g. a table cell. */
    block?: boolean
    /** Taller with a thicker bar, e.g. alone in a popover. */
    large?: boolean
    /** Renders a button: opens a tooltip. */
    interactive?: boolean
  }>(),
  { tone: 'ok' },
)

const width = computed(() => `${Math.round(Math.min(Math.max(props.ratio ?? 1, 0), 1) * 100)}%`)
const fill = computed(() => (props.disabled ? 'off' : props.ratio == null ? 'neutral' : props.tone))
</script>

<template>
  <component
    :is="interactive ? 'button' : 'span'"
    :type="interactive ? 'button' : undefined"
    :class="['x-usage', { 'is-exhausted': exhausted, 'is-disabled': disabled, 'x-usage--textual': textual, 'x-usage--block': block, 'x-usage--large': large, 'x-usage--interactive': interactive }]"
  >
    <span :class="['x-usage__used', `x-usage__used--${tone}`]">{{ used }}</span>
    <span class="x-usage__track" role="presentation">
      <span :class="['x-usage__bar', `x-usage__bar--${fill}`]" :style="{ width }" />
    </span>
    <span class="x-usage__limit">{{ limit }}</span>
  </component>
</template>

<style scoped>
.x-usage {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  width: 200px;
  height: calc(var(--chip-h) + 2px);
  padding: 0 10px;
  border: 1px solid var(--field-line);
  border-radius: var(--r-pill);
  background: var(--field);
  font: 500 var(--fs-11) / 1 var(--font-mono);
  white-space: nowrap;
  outline: none;
}

.x-usage--block {
  width: 100%;
}

.x-usage--large {
  height: 32px;
  padding: 0 12px;
}

.x-usage--interactive {
  cursor: default;
  transition: box-shadow var(--t-fast);
}

.x-usage--interactive:focus-visible,
.x-usage--interactive.is-focus {
  box-shadow: 0 0 0 2px var(--nest), 0 0 0 4px var(--acc-line);
}

.x-usage__used {
  flex: none;
  width: 60px;
  color: var(--fg2);
  text-align: right;
}

.x-usage--textual .x-usage__used {
  width: auto;
  font-family: var(--font-sans);
  text-align: left;
}

.x-usage--textual .x-usage__used--ok { color: var(--ok-t); }
.x-usage--textual .x-usage__used--warn { color: var(--warn-t); }
.x-usage--textual .x-usage__used--err { color: var(--err-t); }

.x-usage__track {
  display: flex;
  flex: 1;
  min-width: 16px;
  height: 4px;
  overflow: hidden;
  border-radius: var(--r-pill);
  background: var(--bar-track);
}

.x-usage--large .x-usage__track {
  height: 5px;
}

.x-usage__bar {
  height: 100%;
  border-radius: var(--r-pill);
  transition: width var(--t-med) var(--ease);
}

.x-usage__bar--ok { background: var(--ok); }
.x-usage__bar--warn { background: var(--warn); }
.x-usage__bar--err { background: var(--err); }
.x-usage__bar--unl { background: var(--unl); }
.x-usage__bar--on { background: var(--on); }
.x-usage__bar--off { background: var(--off); }
.x-usage__bar--neutral { background: var(--bar-unl); }

.x-usage__limit {
  flex: none;
  width: 42px;
  color: var(--mut);
}

.x-usage--textual .x-usage__limit {
  width: auto;
}

.x-usage.is-exhausted {
  border-color: var(--err-line);
  background: var(--err-bg);
}

.x-usage.is-exhausted .x-usage__used {
  color: var(--err-t);
}

.x-usage.is-disabled .x-usage__used {
  color: var(--faint);
}

@media (max-width: 767px) {
  .x-usage__used {
    width: 66px;
  }

  .x-usage__limit {
    width: 46px;
  }
}
</style>
