<script setup lang="ts">
import { computed, useSlots } from 'vue'
import XIcon from './XIcon.vue'
import XSpinner from './XSpinner.vue'
import type { IconName } from './icons'

const props = withDefaults(
  defineProps<{
    variant?: 'primary' | 'secondary' | 'soft' | 'danger' | 'ghost' | 'ghost-danger'
    size?: 'md' | 'sm'
    icon?: IconName
    loading?: boolean
    disabled?: boolean
    block?: boolean
    type?: 'button' | 'submit' | 'reset'
    /** Accessible name for an icon-only button. */
    label?: string
  }>(),
  { variant: 'secondary', size: 'md', type: 'button' },
)

const slots = useSlots()
const iconOnly = computed(() => !slots.default && !!props.icon)
const iconSize = computed(() => (props.size === 'sm' ? 14 : 15))
</script>

<template>
  <button
    :type="type"
    :class="['x-btn', `x-btn--${variant}`, `x-btn--${size}`, { 'x-btn--icon': iconOnly, 'x-btn--block': block, 'is-loading': loading }]"
    :disabled="disabled"
    :aria-disabled="loading || undefined"
    :aria-busy="loading || undefined"
    :aria-label="label"
    :title="iconOnly ? label : undefined"
  >
    <XSpinner v-if="loading" :size="iconSize - 1" />
    <XIcon v-else-if="icon" :name="icon" :size="iconSize" />
    <span v-if="$slots.default" class="x-btn__text"><slot /></span>
  </button>
</template>

<style scoped>
.x-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 7px;
  height: var(--ctl-h);
  padding: 0 var(--btn-px);
  border: 1px solid transparent;
  border-radius: var(--ctl-r);
  font: 600 var(--ctl-fs) / 1 var(--font-sans);
  white-space: nowrap;
  user-select: none;
  outline: none;
  transition: background-color var(--t-fast), border-color var(--t-fast), color var(--t-fast), box-shadow var(--t-fast), transform var(--t-fast);
}

.x-btn--sm {
  height: var(--ctl-sm-h);
  padding: 0 10px;
  gap: 6px;
  border-radius: var(--ctl-sm-r);
  font-size: var(--ctl-sm-fs);
}

.x-btn--icon {
  width: var(--ctl-h);
  padding: 0;
}

.x-btn--sm.x-btn--icon {
  width: var(--ctl-sm-h);
}

.x-btn--block {
  display: flex;
  width: 100%;
}

.x-btn__text {
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
}

.x-btn.is-loading {
  pointer-events: none;
  opacity: .8;
}

.x-btn:focus-visible,
.x-btn.is-focus {
  box-shadow: 0 0 0 2px var(--bg), 0 0 0 4px var(--acc);
}

/* primary */
.x-btn--primary {
  background: var(--acc);
  color: var(--acc-ink);
}
.x-btn--primary:hover:not(:disabled),
.x-btn--primary.is-hover {
  background: var(--acc-hover);
  box-shadow: 0 0 0 3px var(--acc-ring), 0 0 14px var(--acc-glow);
}
.x-btn--primary:active:not(:disabled),
.x-btn--primary.is-active {
  background: var(--acc-press);
  box-shadow: none;
}

/* secondary */
.x-btn--secondary {
  background: var(--ctrl);
  border-color: var(--field-line);
  color: var(--fg2);
  font-weight: 500;
}
.x-btn--secondary:hover:not(:disabled),
.x-btn--secondary.is-hover {
  border-color: var(--field-hover);
  color: var(--fg);
}
.x-btn--secondary:focus-visible,
.x-btn--secondary.is-focus {
  border-color: var(--acc-line);
  box-shadow: 0 0 0 3px var(--acc-ring);
}
.x-btn--secondary:active:not(:disabled),
.x-btn--secondary.is-active {
  background: var(--press);
  border-color: var(--field-hover);
  color: var(--fg);
}

/* soft */
.x-btn--soft {
  background: var(--acc-soft);
  border-color: var(--acc-line);
  color: var(--acc);
}
.x-btn--soft:hover:not(:disabled),
.x-btn--soft.is-hover {
  border-color: var(--acc);
}
.x-btn--soft:active:not(:disabled),
.x-btn--soft.is-active {
  transform: scale(.98);
}

/* danger */
.x-btn--danger {
  background: var(--err-bg);
  border-color: var(--err-line);
  color: var(--err-t);
}
.x-btn--danger:hover:not(:disabled),
.x-btn--danger.is-hover {
  box-shadow: 0 0 0 3px var(--err-ring);
}
.x-btn--danger:focus-visible,
.x-btn--danger.is-focus {
  box-shadow: 0 0 0 2px var(--bg), 0 0 0 4px var(--err);
}
.x-btn--danger:active:not(:disabled),
.x-btn--danger.is-active {
  background: var(--err-ring);
  border-color: var(--err);
}

/* ghost: row actions, "more" buttons */
.x-btn--ghost {
  background: transparent;
  color: var(--dim);
  font-weight: 500;
}
.x-btn--ghost:hover:not(:disabled),
.x-btn--ghost.is-hover {
  background: var(--hov);
  color: var(--fg);
}
.x-btn--ghost:focus-visible,
.x-btn--ghost.is-focus {
  box-shadow: 0 0 0 2px var(--acc-line);
  color: var(--fg);
}
.x-btn--ghost:active:not(:disabled),
.x-btn--ghost.is-active {
  background: var(--press);
  color: var(--fg);
}

/* ghost-danger: delete action in a row */
.x-btn--ghost-danger {
  background: transparent;
  color: var(--err-t);
  font-weight: 500;
}
.x-btn--ghost-danger:hover:not(:disabled),
.x-btn--ghost-danger.is-hover {
  background: var(--err-bg);
}
.x-btn--ghost-danger:focus-visible,
.x-btn--ghost-danger.is-focus {
  box-shadow: 0 0 0 2px var(--err-line);
}
.x-btn--ghost-danger:active:not(:disabled),
.x-btn--ghost-danger.is-active {
  background: var(--err-ring);
}

.x-btn:disabled {
  cursor: not-allowed;
  background: var(--ctrl);
  border-color: var(--line);
  color: var(--faint);
  box-shadow: none;
}
.x-btn--ghost:disabled,
.x-btn--ghost-danger:disabled {
  background: transparent;
  border-color: transparent;
}
</style>
