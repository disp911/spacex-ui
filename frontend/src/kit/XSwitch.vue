<script setup lang="ts">
import { computed, useAttrs } from 'vue'
import XSpinner from './XSpinner.vue'

defineOptions({ inheritAttrs: false })

const props = defineProps<{
  disabled?: boolean
  /** The change is being applied; the knob shows a spinner. */
  loading?: boolean
  /** Applying failed and the switch was rolled back. */
  invalid?: boolean
  modified?: boolean
  label?: string
  /**
   * A word of state to the left, e.g. «On» for a setting that applies at once:
   * accent when on, dim when off, amber while the change is applied.
   */
  text?: string
}>()

const model = defineModel<boolean>({ default: false })
const attrs = useAttrs()

function toggle() {
  model.value = !model.value
}

const button = computed(() => ({
  ...attrs,
  type: 'button' as const,
  role: 'switch',
  'aria-checked': model.value,
  'aria-label': props.label,
  'aria-busy': props.loading || undefined,
  'aria-invalid': props.invalid || undefined,
  disabled: props.disabled,
  class: ['x-switch', attrs.class, { 'is-on': model.value, 'is-loading': props.loading, 'is-invalid': props.invalid, 'is-modified': props.modified }],
}))
</script>

<template>
  <span v-if="text" class="x-switch-line">
    <span :class="['x-switch-line__text', { 'is-on': model, 'is-loading': loading }]" aria-hidden="true">{{ text }}</span>
    <button v-bind="button" @click="loading || toggle()">
      <span class="x-switch__knob"><XSpinner v-if="loading" :size="10" /></span>
    </button>
  </span>
  <button v-else v-bind="button" @click="loading || toggle()">
    <span class="x-switch__knob"><XSpinner v-if="loading" :size="10" /></span>
  </button>
</template>

<style scoped>
.x-switch-line {
  display: inline-flex;
  align-items: center;
  gap: 10px;
}

.x-switch-line__text {
  color: var(--dim);
  font: 500 var(--fs-12) / var(--lh-tight) var(--font-sans);
  white-space: nowrap;
  transition: color var(--t-fast);
}

.x-switch-line__text.is-on {
  color: var(--acc);
}

.x-switch-line__text.is-loading {
  color: var(--warn-t);
}

.x-switch {
  --knob-size: calc(var(--switch-h) - 6px);
  /* The knob sits evenly inside the track. */
  --inset: calc((var(--switch-h) - var(--knob-size)) / 2);
  position: relative;
  flex: none;
  width: var(--switch-w);
  height: var(--switch-h);
  padding: 0;
  border: 0;
  border-radius: var(--r-pill);
  background: var(--track);
  outline: none;
  transition: background-color var(--t-med), box-shadow var(--t-fast), filter var(--t-fast);
}

.x-switch__knob {
  position: absolute;
  top: var(--inset);
  left: var(--inset);
  display: flex;
  align-items: center;
  justify-content: center;
  width: var(--knob-size);
  height: var(--knob-size);
  border-radius: var(--r-pill);
  background: var(--knob);
  color: var(--dim);
  box-shadow: var(--shadow-knob);
  transition: left var(--t-med) var(--ease), width var(--t-fast) var(--ease), background-color var(--t-med);
}

.x-switch.is-on {
  background: var(--acc);
}

.x-switch.is-on .x-switch__knob {
  left: calc(100% - var(--knob-size) - var(--inset));
  background: var(--knob-on);
  color: var(--acc);
}

.x-switch:hover:not(:disabled),
.x-switch.is-hover {
  filter: brightness(1.12);
}

.x-switch:active:not(:disabled) .x-switch__knob,
.x-switch.is-active .x-switch__knob {
  width: calc(var(--knob-size) + var(--inset));
}

.x-switch.is-on:active:not(:disabled) .x-switch__knob,
.x-switch.is-on.is-active .x-switch__knob {
  left: calc(100% - var(--knob-size) - 2 * var(--inset));
}

.x-switch:focus-visible,
.x-switch.is-focus {
  box-shadow: 0 0 0 2px var(--bg), 0 0 0 4px var(--acc-line);
}

.x-switch.is-modified {
  box-shadow: 0 0 0 2px var(--bg), 0 0 0 4px var(--mod-line);
}

.x-switch.is-invalid {
  box-shadow: 0 0 0 1px var(--err), 0 0 0 4px var(--err-ring);
}

.x-switch.is-loading {
  cursor: progress;
  opacity: .8;
}

.x-switch:disabled {
  cursor: not-allowed;
  opacity: .4;
  filter: grayscale(.6);
}
</style>
