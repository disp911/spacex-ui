<script setup lang="ts" generic="T extends OptionValue">
import XCheckbox from './XCheckbox.vue'
import type { OptionValue } from './types'

/** Checkboxes in a row for a list of choices, e.g. Sniffing's http, tls, quic and fakedns; wraps when narrow. */
const props = defineProps<{
  options: { value: T; label: string }[]
  /** Values and codes: monospace labels. */
  mono?: boolean
  disabled?: boolean
  label?: string
}>()

const model = defineModel<T[]>({ default: () => [] })

const checked = (v: T) => model.value.includes(v)
function toggle(v: T, on: boolean) {
  // Kept in the order of the options, as the old panel's checkbox group did.
  const next = new Set(model.value)
  if (on) next.add(v)
  else next.delete(v)
  model.value = props.options.map((o) => o.value).filter((x) => next.has(x))
}
</script>

<template>
  <div class="x-check-group" role="group" :aria-label="label">
    <XCheckbox v-for="o in options" :key="String(o.value)" :model-value="checked(o.value)" :disabled="disabled" @update:model-value="toggle(o.value, $event)">
      <span :class="{ 'x-check-group__mono': mono }">{{ o.label }}</span>
    </XCheckbox>
  </div>
</template>

<style scoped>
.x-check-group {
  display: flex;
  flex-wrap: wrap;
  gap: 8px 16px;
}

.x-check-group__mono {
  font-family: var(--font-mono);
  font-size: var(--fs-12);
}
</style>
