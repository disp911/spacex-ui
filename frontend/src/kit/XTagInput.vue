<script setup lang="ts">
import { computed, ref, useId, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import XIcon from './XIcon.vue'
import { useFloating } from './floating'
import { useFieldAttrs } from './useFieldAttrs'

defineOptions({ inheritAttrs: false })

const props = defineProps<{
  /** Suggested values offered in a dropdown while typing. */
  presets?: string[]
  placeholder?: string
  mono?: boolean
  disabled?: boolean
  invalid?: boolean
  modified?: boolean
}>()

const model = defineModel<string[]>({ default: () => [] })
const { t } = useI18n()
const { wrapper, control } = useFieldAttrs()

const id = useId()
const box = ref<HTMLElement | null>(null)
const input = ref<HTMLInputElement | null>(null)
const list = ref<HTMLElement | null>(null)
const text = ref('')
const focused = ref(false)
const active = ref(0)

const suggestions = computed(() => {
  const q = text.value.trim().toLowerCase()
  return (props.presets ?? []).filter((p) => !model.value.includes(p) && p.toLowerCase().includes(q))
})
const open = computed(() => focused.value && suggestions.value.length > 0)
const { style } = useFloating(box, list, open, { matchWidth: true, onDismiss: () => (focused.value = false) })
watch(text, () => (active.value = 0))

function add(value: string) {
  const v = value.trim()
  if (v && !model.value.includes(v)) model.value = [...model.value, v]
  text.value = ''
  active.value = 0
}

function remove(i: number) {
  model.value = model.value.filter((_, j) => j !== i)
  input.value?.focus()
}

function onKeydown(e: KeyboardEvent) {
  const n = suggestions.value.length
  // Enter takes the highlighted suggestion, a comma always takes the typed text.
  if (e.key === 'Enter' && open.value) {
    e.preventDefault()
    add(suggestions.value[active.value] ?? text.value)
  } else if ((e.key === 'Enter' || e.key === ',') && text.value.trim()) {
    e.preventDefault()
    add(text.value)
  } else if (e.key === 'Backspace' && !text.value && model.value.length) {
    remove(model.value.length - 1)
  } else if ((e.key === 'ArrowDown' || e.key === 'ArrowUp') && n) {
    e.preventDefault()
    active.value = (active.value + (e.key === 'ArrowDown' ? 1 : -1) + n) % n
  }
}

function onBlur() {
  if (text.value.trim()) add(text.value)
}
</script>

<template>
  <div
    ref="box"
    v-bind="wrapper"
    :class="['x-field-box', 'x-tags', { 'x-field-box--mono': mono, 'is-disabled': disabled, 'is-invalid': invalid, 'is-modified': modified }]"
    @click="input?.focus()"
  >
    <span v-for="(tag, i) in model" :key="tag" class="x-tags__chip">
      <span class="x-tags__text">{{ tag }}</span>
      <button v-if="!disabled" type="button" class="x-tags__remove" :aria-label="t('kit.remove', { value: tag })" @click.stop="remove(i)">
        <XIcon name="close" :size="10" :stroke="1.8" />
      </button>
    </span>
    <input
      ref="input"
      v-model="text"
      v-bind="control"
      class="x-tags__control"
      role="combobox"
      aria-autocomplete="list"
      :aria-expanded="open"
      :aria-controls="`${id}-list`"
      :aria-activedescendant="open ? `${id}-${active}` : undefined"
      :placeholder="model.length ? undefined : placeholder"
      :disabled="disabled"
      autocomplete="off"
      @focus="focused = true"
      @blur="onBlur"
      @keydown="onKeydown"
    />
  </div>
  <Teleport to="body">
    <Transition name="x-popup">
      <ul v-if="open" :id="`${id}-list`" ref="list" role="listbox" class="x-popup x-tags__list" :style="style">
        <li
          v-for="(s, i) in suggestions"
          :id="`${id}-${i}`"
          :key="s"
          role="option"
          :aria-selected="i === active"
          :class="['x-popup-item', { 'is-active': i === active }]"
          @pointermove="active = i"
          @pointerdown.prevent
          @click="add(s)"
        >
          {{ s }}
        </li>
      </ul>
    </Transition>
  </Teleport>
</template>

<style scoped>
.x-tags {
  flex-wrap: wrap;
  gap: 5px;
  height: auto;
  min-height: var(--ctl-h);
  padding: 4px 5px;
  cursor: text;
}

.x-tags__chip {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  max-width: 100%;
  height: var(--chip-h);
  padding: 0 4px 0 8px;
  border: 1px solid var(--chip-line);
  border-radius: var(--r-6);
  background: var(--chip);
  color: var(--fg2);
  font: 500 var(--fs-11) / 1 var(--font-sans);
  white-space: nowrap;
}

.x-field-box--mono .x-tags__chip {
  font-family: var(--font-mono);
}

.x-tags__text {
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
}

.x-tags__remove {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 16px;
  height: 16px;
  padding: 0;
  border: 0;
  border-radius: var(--r-4);
  background: transparent;
  color: var(--dim);
  outline: none;
}

.x-tags__remove:hover,
.x-tags__remove:focus-visible {
  background: var(--press);
  color: var(--fg);
}

.x-tags__control {
  flex: 1;
  min-width: 80px;
  height: calc(var(--ctl-h) - 10px);
  padding: 0 6px;
  border: 0;
  outline: none;
  background: transparent;
  font: inherit;
  color: inherit;
}

.x-tags__control::placeholder {
  color: var(--faint);
}

.x-tags__list {
  margin: 0;
  list-style: none;
}
</style>
