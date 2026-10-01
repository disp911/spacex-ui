<script setup lang="ts">
import XIcon from './XIcon.vue'

// Label + control + hint/error. The control goes in the default slot.
defineProps<{ label: string; hint?: string; error?: string; locked?: boolean; for?: string }>()
</script>

<template>
  <div class="x-field">
    <label class="caps x-field__label" :for="$props.for">
      {{ label }}
      <XIcon v-if="locked" name="lock" :size="11" class="x-field__lock" />
    </label>
    <slot />
    <div v-if="error" class="x-field__error" role="alert">
      <XIcon name="alert" :size="13" />
      <span>{{ error }}</span>
    </div>
    <div v-else-if="hint || $slots.hint" class="x-field__hint"><slot name="hint">{{ hint }}</slot></div>
  </div>
</template>

<style scoped>
.x-field {
  display: flex;
  flex-direction: column;
  min-width: 0;
}
.x-field__label {
  display: flex;
  align-items: center;
  gap: 7px;
  margin-bottom: 7px;
}
.x-field__lock {
  color: var(--text-5);
}
.x-field__hint {
  margin-top: 6px;
  font: var(--fw-regular) var(--fs-xs) / 1.4 var(--font-sans);
  color: var(--text-4);
}
.x-field__error {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-top: 7px;
  font: var(--fw-medium) var(--fs-sm) var(--font-sans);
  color: var(--danger-text);
}
</style>
