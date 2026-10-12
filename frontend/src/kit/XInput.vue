<script setup lang="ts">
import { computed, ref } from 'vue'
import XIcon from './XIcon.vue'
import type { IconName } from './icons'
import { useFieldAttrs } from './useFieldAttrs'

defineOptions({ inheritAttrs: false })

const props = withDefaults(
  defineProps<{
    type?: 'text' | 'password' | 'email' | 'url' | 'search'
    placeholder?: string
    /** Monospace text for technical values: paths, keys, addresses. */
    mono?: boolean
    disabled?: boolean
    readonly?: boolean
    invalid?: boolean
    /** Unsaved change. */
    modified?: boolean
    icon?: IconName
    clearable?: boolean
    /** Text after the value, e.g. a unit. */
    suffix?: string
  }>(),
  { type: 'text' },
)

const model = defineModel<string>({ default: '' })
const { wrapper, control } = useFieldAttrs()
const input = ref<HTMLInputElement | null>(null)
const revealed = ref(false)
const nativeType = computed(() => (props.type === 'password' && revealed.value ? 'text' : props.type))

function clear() {
  model.value = ''
  input.value?.focus()
}

defineExpose({ focus: () => input.value?.focus() })
</script>

<template>
  <div
    v-bind="wrapper"
    :class="['x-field-box', 'x-input', { 'x-field-box--mono': mono, 'is-disabled': disabled, 'is-invalid': invalid, 'is-modified': modified }]"
  >
    <XIcon v-if="icon" :name="icon" :size="15" class="x-input__icon" />
    <input
      ref="input"
      v-model="model"
      v-bind="control"
      class="x-input__control"
      :type="nativeType"
      :placeholder="placeholder"
      :disabled="disabled"
      :readonly="readonly"
      :aria-invalid="invalid || undefined"
      :spellcheck="mono ? false : undefined"
    />
    <span v-if="suffix" class="x-input__suffix">{{ suffix }}</span>
    <button
      v-if="clearable && model && !disabled && !readonly"
      type="button"
      class="x-input__action"
      aria-label="Clear"
      @click="clear"
    >
      <XIcon name="close" :size="13" :stroke="1.6" />
    </button>
    <button
      v-if="type === 'password'"
      type="button"
      class="x-input__action"
      :aria-label="revealed ? 'Hide' : 'Show'"
      :aria-pressed="revealed"
      :disabled="disabled"
      @click="revealed = !revealed"
    >
      <XIcon :name="revealed ? 'eyeOff' : 'eye'" :size="15" :stroke="1.35" />
    </button>
  </div>
</template>

<style scoped>
.x-input {
  gap: 8px;
}

.x-input__control {
  flex: 1;
  min-width: 0;
  height: 100%;
  padding: 0;
  border: 0;
  outline: none;
  background: transparent;
  font: inherit;
  color: inherit;
  text-overflow: ellipsis;
}

.x-input__control::placeholder {
  color: var(--faint);
}

.x-input__icon,
.x-input__suffix {
  flex: none;
  color: var(--dim);
}

.x-input__suffix {
  font: 400 var(--fs-12) var(--font-mono);
}

.x-input__action {
  flex: none;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  margin-right: -6px;
  padding: 0;
  border: 0;
  border-radius: var(--r-8);
  background: transparent;
  color: var(--dim);
  outline: none;
}

.x-input__action:hover {
  background: var(--hov);
  color: var(--fg);
}

.x-input__action:focus-visible {
  box-shadow: 0 0 0 2px var(--acc-line);
  color: var(--fg);
}
</style>
