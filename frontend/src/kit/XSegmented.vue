<script setup lang="ts" generic="T extends OptionValue">
import { ref } from 'vue'
import XIcon from './XIcon.vue'
import type { IconName } from './icons'
import type { OptionValue } from './types'

const props = defineProps<{
  options: { value: T; label: string; icon?: IconName; disabled?: boolean }[]
  /** Show only icons; labels become tooltips and accessible names. */
  iconOnly?: boolean
  /** Takes the width it is given; the options share it equally. */
  block?: boolean
  disabled?: boolean
  label?: string
}>()

const model = defineModel<T>({ required: true })
const root = ref<HTMLElement | null>(null)

function onKeydown(e: KeyboardEvent) {
  const dir = e.key === 'ArrowRight' || e.key === 'ArrowDown' ? 1 : e.key === 'ArrowLeft' || e.key === 'ArrowUp' ? -1 : 0
  if (!dir) return
  e.preventDefault()
  const n = props.options.length
  let i = props.options.findIndex((o) => o.value === model.value)
  for (let k = 0; k < n; k++) {
    i = (i + dir + n) % n
    if (!props.options[i]?.disabled) break
  }
  const next = props.options[i]
  if (!next) return
  model.value = next.value
  root.value?.querySelectorAll<HTMLElement>('[role="radio"]')[i]?.focus()
}
</script>

<template>
  <div ref="root" role="radiogroup" :aria-label="label" :class="['x-seg', { 'x-seg--block': block, 'is-disabled': disabled }]" @keydown="onKeydown">
    <button
      v-for="o in options"
      :key="String(o.value)"
      type="button"
      role="radio"
      :aria-checked="o.value === model"
      :aria-label="iconOnly ? o.label : undefined"
      :title="iconOnly ? o.label : undefined"
      :tabindex="o.value === model ? 0 : -1"
      :disabled="disabled || o.disabled"
      :class="['x-seg__item', { 'is-selected': o.value === model, 'x-seg__item--icon': iconOnly }]"
      @click="model = o.value"
    >
      <XIcon v-if="o.icon" :name="o.icon" :size="15" />
      <span v-if="!iconOnly">{{ o.label }}</span>
    </button>
  </div>
</template>

<style scoped>
.x-seg {
  display: inline-flex;
  align-items: stretch;
  gap: 2px;
  height: var(--ctl-h);
  padding: 2px;
  border: 1px solid var(--field-line);
  border-radius: var(--ctl-r);
  background: var(--field);
}

.x-seg--block {
  display: flex;
  width: 100%;
}

.x-seg--block .x-seg__item {
  flex: 1;
  min-width: 0;
}

.x-seg__item {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 0 12px;
  border: 0;
  border-radius: var(--r-8);
  background: transparent;
  color: var(--mut);
  font: 500 var(--fs-12) / 1 var(--font-sans);
  white-space: nowrap;
  outline: none;
  transition: background-color var(--t-fast), color var(--t-fast), box-shadow var(--t-fast);
}

.x-seg__item--icon {
  width: calc(var(--ctl-h) - 4px);
  padding: 0;
  color: var(--dim);
}

.x-seg__item:hover:not(:disabled),
.x-seg__item.is-hover {
  color: var(--fg);
}

.x-seg__item:focus-visible,
.x-seg__item.is-focus {
  box-shadow: inset 0 0 0 2px var(--acc-line);
  color: var(--fg);
}

.x-seg__item:active:not(:disabled),
.x-seg__item.is-active {
  opacity: .75;
}

.x-seg__item.is-selected {
  background: var(--seg);
  color: var(--fg);
  font-weight: 600;
  box-shadow: var(--shadow-seg);
}

.x-seg__item--icon.is-selected {
  color: var(--acc);
}

.x-seg__item:disabled {
  color: var(--faint);
  cursor: not-allowed;
}

.x-seg.is-disabled {
  opacity: .45;
}

@media (max-width: 767px) {
  .x-seg__item {
    font-size: var(--fs-13);
  }
}
</style>
