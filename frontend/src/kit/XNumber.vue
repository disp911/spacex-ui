<script setup lang="ts">
import { ref, watch } from 'vue'
import { useFieldAttrs } from './useFieldAttrs'

defineOptions({ inheritAttrs: false })

const props = defineProps<{
  min?: number
  max?: number
  /** Allow fractions; whole numbers otherwise. */
  decimal?: boolean
  unit?: string
  /** What the arrow keys add or take away; Shift makes it ten times more. 1 by default. */
  step?: number
  placeholder?: string
  disabled?: boolean
  readonly?: boolean
  invalid?: boolean
  modified?: boolean
}>()

const model = defineModel<number | null>({ default: null })
const { wrapper, control } = useFieldAttrs()
const text = ref(format(model.value))

watch(model, (v) => {
  if (parse(text.value) !== v) text.value = format(v)
})

function format(v: number | null) {
  return v == null ? '' : String(v)
}

function parse(s: string): number | null {
  const t = s.trim().replace(',', '.')
  if (t === '' || t === '-') return null
  const n = props.decimal ? Number(t) : Number.parseInt(t, 10)
  return Number.isFinite(n) ? n : null
}

function onInput(e: Event) {
  const raw = (e.target as HTMLInputElement).value
  const allowed = props.decimal ? /[^\d.,-]/g : /[^\d-]/g
  text.value = raw.replace(allowed, '')
  ;(e.target as HTMLInputElement).value = text.value
  model.value = parse(text.value)
}

function clamp() {
  let v = parse(text.value)
  if (v != null) {
    if (props.min != null && v < props.min) v = props.min
    if (props.max != null && v > props.max) v = props.max
  }
  model.value = v
  text.value = format(v)
}

function nudge(e: KeyboardEvent, by: number) {
  if (props.readonly) return
  e.preventDefault()
  const size = (props.step ?? 1) * (e.shiftKey ? 10 : 1)
  text.value = format((parse(text.value) ?? props.min ?? 0) + by * size)
  clamp()
}
</script>

<template>
  <div
    v-bind="wrapper"
    :class="['x-field-box', 'x-field-box--mono', 'x-number', { 'is-disabled': disabled, 'is-invalid': invalid, 'is-modified': modified }]"
  >
    <input
      v-bind="control"
      class="x-number__control"
      type="text"
      :inputmode="decimal ? 'decimal' : 'numeric'"
      :value="text"
      :placeholder="placeholder"
      :disabled="disabled"
      :readonly="readonly"
      :aria-invalid="invalid || undefined"
      :aria-valuemin="min"
      :aria-valuemax="max"
      autocomplete="off"
      @input="onInput"
      @blur="clamp"
      @keydown.up="nudge($event, 1)"
      @keydown.down="nudge($event, -1)"
    />
    <span v-if="unit" class="x-number__unit">{{ unit }}</span>
  </div>
</template>

<style scoped>
.x-number {
  gap: 8px;
  width: 140px;
}

.x-number__control {
  flex: 1;
  min-width: 0;
  height: 100%;
  padding: 0;
  border: 0;
  outline: none;
  background: transparent;
  font: inherit;
  color: inherit;
}

.x-number__control::placeholder {
  color: var(--faint);
}

.x-number__unit {
  flex: none;
  color: var(--dim);
  font: 400 var(--fs-12) var(--font-mono);
}
</style>
