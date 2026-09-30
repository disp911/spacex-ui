<script setup lang="ts">
import { computed } from 'vue'
import XIcon, { type IconName } from './XIcon.vue'
import XSpinner from './XSpinner.vue'

const props = withDefaults(
  defineProps<{
    variant?: 'primary' | 'secondary' | 'soft' | 'danger' | 'danger-solid' | 'ghost'
    size?: 'xs' | 'sm' | 'md' | 'lg'
    icon?: IconName
    /** Square icon-only button; pass a label for accessibility. */
    iconOnly?: boolean
    label?: string
    loading?: boolean
    disabled?: boolean
    block?: boolean
    href?: string
    type?: 'button' | 'submit'
  }>(),
  { variant: 'secondary', size: 'md', type: 'button' },
)

defineEmits<{ click: [MouseEvent] }>()

const iconSize = computed(() => ({ xs: 13, sm: 14, md: 16, lg: 16 })[props.size])
</script>

<template>
  <component
    :is="href ? 'a' : 'button'"
    :href="href"
    :type="href ? undefined : type"
    class="x-btn"
    :class="[`x-btn--${variant}`, `x-btn--${size}`, { 'x-btn--icon': iconOnly, 'x-btn--block': block, 'is-loading': loading }]"
    :disabled="href ? undefined : disabled || loading"
    :aria-label="iconOnly ? label : undefined"
    :title="iconOnly ? label : undefined"
    @click="$emit('click', $event)"
  >
    <XSpinner v-if="loading" :size="iconSize" />
    <XIcon v-else-if="icon" :name="icon" :size="iconSize" />
    <span v-if="!iconOnly && ($slots.default || label)" class="x-btn__label"><slot>{{ label }}</slot></span>
  </component>
</template>

<style scoped>
.x-btn {
  --btn-h: var(--h-md);
  --btn-r: var(--r-md);
  --btn-px: 14px;
  --btn-fs: var(--fs-md);

  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: var(--sp-2);
  height: var(--btn-h);
  padding: 0 var(--btn-px);
  border-radius: var(--btn-r);
  border: 1px solid transparent;
  font: var(--fw-medium) var(--btn-fs) / 1 var(--font-sans);
  white-space: nowrap;
  text-decoration: none;
  user-select: none;
  flex: none;
  transition:
    background var(--dur-base) ease,
    border-color var(--dur-base) ease,
    box-shadow var(--dur-base) ease,
    color var(--dur-fast) ease;
}
.x-btn:disabled {
  cursor: not-allowed;
  opacity: 0.5;
}
.x-btn.is-loading {
  cursor: progress;
  opacity: 1;
}
.x-btn__label {
  overflow: hidden;
  text-overflow: ellipsis;
  line-height: 1.3;
}

/* Sizes */
.x-btn--xs { --btn-h: var(--h-xs); --btn-r: var(--r-sm); --btn-px: 9px; --btn-fs: var(--fs-xs); }
.x-btn--sm { --btn-h: var(--h-sm); --btn-r: var(--r-sm); --btn-px: 12px; --btn-fs: var(--fs-sm); }
.x-btn--md { --btn-h: var(--h-md); --btn-r: var(--r-md); --btn-px: 14px; --btn-fs: var(--fs-md); }
.x-btn--lg { --btn-h: var(--h-lg); --btn-r: var(--r-md); --btn-px: 16px; --btn-fs: var(--fs-lg); }
.x-btn--icon { width: var(--btn-h); padding: 0; }
.x-btn--block { display: flex; width: 100%; }

/* Primary — filled accent */
.x-btn--primary {
  background: var(--accent);
  color: var(--on-accent);
  font-weight: var(--fw-semibold);
}
.x-btn--primary:not(:disabled):hover {
  background: var(--accent-hover);
  box-shadow: 0 0 0 4px var(--accent-ring), 0 0 22px var(--accent-glow);
}
.x-btn--primary.is-loading {
  background: var(--accent-press);
  box-shadow: none;
}

/* Secondary — neutral outlined control */
.x-btn--secondary {
  background: var(--control);
  border-color: var(--border-control);
  color: var(--text-2);
}
.x-btn--secondary:not(:disabled):hover {
  background: var(--accent-soft-hover);
  border-color: var(--accent-focus);
  box-shadow: 0 0 0 3px var(--accent-ring), 0 0 14px var(--accent-glow);
  color: var(--text);
}

/* Soft — tinted accent, used for call-outs like "Update panel" */
.x-btn--soft {
  background: var(--accent-soft);
  border-color: var(--accent-border);
  color: var(--accent);
}
.x-btn--soft:not(:disabled):hover {
  border-color: var(--accent-focus);
  box-shadow: 0 0 0 3px var(--accent-ring), 0 0 14px var(--accent-glow);
}

/* Danger — tinted, for destructive secondary actions */
.x-btn--danger {
  background: var(--danger-soft);
  border-color: var(--danger-border);
  color: var(--danger-text);
}
.x-btn--danger:not(:disabled):hover {
  border-color: var(--danger);
  box-shadow: 0 0 0 3px var(--danger-ring);
}

/* Danger solid — confirm destructive action */
.x-btn--danger-solid {
  background: var(--danger);
  color: var(--on-danger);
  font-weight: var(--fw-semibold);
}
.x-btn--danger-solid:not(:disabled):hover {
  box-shadow: 0 0 0 4px var(--danger-ring);
}

/* Ghost — no chrome until hovered */
.x-btn--ghost {
  background: transparent;
  color: var(--text-3);
}
.x-btn--ghost:not(:disabled):hover {
  background: var(--hover);
  color: var(--text);
}
</style>
