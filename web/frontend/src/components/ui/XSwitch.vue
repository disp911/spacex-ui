<script setup lang="ts">
// On/off toggle. `loading` keeps the knob where it is and dims the track.
const on = defineModel<boolean>({ default: false })
defineProps<{ label?: string; loading?: boolean; disabled?: boolean }>()
</script>

<template>
  <button
    type="button"
    role="switch"
    class="x-switch"
    :class="{ 'is-on': on, 'is-loading': loading }"
    :aria-checked="on"
    :aria-label="label"
    :title="label"
    :disabled="disabled || loading"
    @click.stop="on = !on"
  >
    <span class="x-switch__knob" />
  </button>
</template>

<style scoped>
.x-switch {
  position: relative;
  flex: none;
  width: 34px;
  height: 20px;
  padding: 0;
  border-radius: var(--r-pill);
  border: 1px solid var(--border-control);
  background: var(--control);
  transition: background var(--dur-base) ease, border-color var(--dur-base) ease, box-shadow var(--dur-base) ease;
}
.x-switch:not(:disabled):hover {
  border-color: var(--accent-focus);
  box-shadow: 0 0 0 3px var(--accent-ring);
}
.x-switch__knob {
  position: absolute;
  top: 2px;
  left: 2px;
  width: 14px;
  height: 14px;
  border-radius: 50%;
  background: var(--text-4);
  transition: transform var(--dur-base) var(--ease-out), background var(--dur-base) ease;
}
.x-switch.is-on {
  background: var(--accent);
  border-color: var(--accent);
}
.x-switch.is-on .x-switch__knob {
  transform: translateX(14px);
  background: var(--on-accent);
}
.x-switch.is-loading {
  opacity: 0.6;
  cursor: progress;
}
.x-switch:disabled:not(.is-loading) {
  opacity: 0.45;
  cursor: not-allowed;
}
</style>
