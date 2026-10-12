<script setup lang="ts">
import { useFieldAttrs } from './useFieldAttrs'

defineOptions({ inheritAttrs: false })

withDefaults(
  defineProps<{
    placeholder?: string
    rows?: number
    mono?: boolean
    disabled?: boolean
    readonly?: boolean
    invalid?: boolean
    modified?: boolean
  }>(),
  { rows: 3 },
)

const model = defineModel<string>({ default: '' })
const { wrapper, control } = useFieldAttrs()
</script>

<template>
  <div
    v-bind="wrapper"
    :class="['x-field-box', 'x-textarea', { 'x-field-box--mono': mono, 'is-disabled': disabled, 'is-invalid': invalid, 'is-modified': modified }]"
  >
    <textarea
      v-model="model"
      v-bind="control"
      class="x-textarea__control"
      :rows="rows"
      :placeholder="placeholder"
      :disabled="disabled"
      :readonly="readonly"
      :aria-invalid="invalid || undefined"
      :spellcheck="mono ? false : undefined"
    />
  </div>
</template>

<style scoped>
.x-textarea {
  height: auto;
  min-height: 76px;
  padding: 8px 10px 8px var(--input-px);
  align-items: stretch;
}

.x-textarea__control {
  flex: 1;
  min-width: 0;
  min-height: 58px;
  padding: 0;
  border: 0;
  outline: none;
  background: transparent;
  font: inherit;
  line-height: var(--lh-text);
  color: inherit;
  resize: vertical;
}

.x-field-box--mono .x-textarea__control {
  line-height: var(--lh-code);
}

.x-textarea__control::placeholder {
  color: var(--faint);
}
</style>
