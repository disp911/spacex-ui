<script setup lang="ts">
import XIcon from './XIcon.vue'

defineProps<{
  disabled?: boolean
  indeterminate?: boolean
  invalid?: boolean
}>()

const model = defineModel<boolean>({ default: false })
</script>

<template>
  <label :class="['x-check', { 'is-disabled': disabled }]">
    <input
      v-model="model"
      class="x-check__input sr-only"
      type="checkbox"
      :disabled="disabled"
      :aria-invalid="invalid || undefined"
      :indeterminate="indeterminate"
    />
    <span :class="['x-check__box', { 'is-on': model || indeterminate, 'is-invalid': invalid }]" aria-hidden="true">
      <span v-if="indeterminate" class="x-check__dash" />
      <XIcon v-else-if="model" name="check" :size="12" :stroke="2.2" />
    </span>
    <span v-if="$slots.default" class="x-check__label"><slot /></span>
  </label>
</template>

<style scoped>
.x-check {
  position: relative;
  display: inline-flex;
  align-items: center;
  gap: 8px;
  min-height: var(--check);
  color: var(--fg2);
  font-size: var(--ctl-fs);
  cursor: pointer;
  user-select: none;
}

.x-check__box {
  display: flex;
  flex: none;
  align-items: center;
  justify-content: center;
  width: var(--check);
  height: var(--check);
  border: 1px solid var(--field-hover);
  border-radius: var(--check-r);
  background: var(--field);
  color: var(--acc-ink);
  transition: background-color var(--t-fast), border-color var(--t-fast), box-shadow var(--t-fast);
}

.x-check:hover .x-check__box,
.x-check.is-hover .x-check__box {
  border-color: var(--acc-line);
}

.x-check__box.is-on {
  border-color: var(--acc);
  background: var(--acc);
}

.x-check__dash {
  width: 8px;
  height: 2px;
  border-radius: 1px;
  background: currentColor;
}

.x-check__input:focus-visible + .x-check__box,
.x-check.is-focus .x-check__box {
  box-shadow: 0 0 0 2px var(--bg), 0 0 0 4px var(--acc-line);
}

.x-check__box.is-invalid {
  border-color: var(--err);
  box-shadow: 0 0 0 3px var(--err-ring);
}

.x-check.is-disabled {
  cursor: not-allowed;
  opacity: .45;
}
</style>
