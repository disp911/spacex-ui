<script setup lang="ts">
import { computed, nextTick, ref, useId, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import XIcon from './XIcon.vue'
import { useFloating } from './floating'
import { useFieldAttrs } from './useFieldAttrs'

defineOptions({ inheritAttrs: false })

type Preset = string | { value: string; label: string }

/**
 * Free values as chips, typed in and taken with Enter or a comma. Presets open
 * under the field on focus, filtered by the typed text: choosing one toggles
 * it, a chosen one is ticked, and a chip shows a preset's label. Typed text
 * that is no preset is offered at the bottom as «Add …».
 */
const props = defineProps<{
  presets?: Preset[]
  /** A caps title over the presets. */
  presetsTitle?: string
  placeholder?: string
  /** The placeholder once something is chosen. */
  morePlaceholder?: string
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

const all = computed(() => (props.presets ?? []).map((p) => (typeof p === 'string' ? { value: p, label: p } : p)))
const label = (v: string) => all.value.find((p) => p.value === v)?.label ?? v

const query = computed(() => text.value.trim())
const shown = computed(() => {
  const q = query.value.toLowerCase()
  return all.value.filter((p) => !q || p.label.toLowerCase().includes(q) || p.value.toLowerCase().includes(q))
})
/** Typed text that is neither a preset nor chosen already. */
const custom = computed(() => {
  const q = query.value
  return q && !model.value.includes(q) && !all.value.some((p) => p.value === q || p.label === q) ? q : ''
})
const rows = computed(() => shown.value.length + (custom.value ? 1 : 0))
const open = computed(() => focused.value && !props.disabled && rows.value > 0)
const { style } = useFloating(box, list, open, { matchWidth: true, onDismiss: () => (focused.value = false) })
watch(query, () => (active.value = 0))

function add(value: string) {
  const v = value.trim()
  if (v && !model.value.includes(v)) model.value = [...model.value, v]
  text.value = ''
}

function toggle(value: string) {
  model.value = model.value.includes(value) ? model.value.filter((x) => x !== value) : [...model.value, value]
  text.value = ''
}

function remove(i: number) {
  model.value = model.value.filter((_, j) => j !== i)
  input.value?.focus()
}

function pick(i: number) {
  const p = shown.value[i]
  if (p) toggle(p.value)
  else if (custom.value) add(custom.value)
}

function onKeydown(e: KeyboardEvent) {
  if (e.key === 'Enter' && open.value) {
    e.preventDefault()
    // A typed value goes in as it is unless a preset row is under the cursor.
    if (query.value && !shown.value[active.value]) add(query.value)
    else pick(active.value)
  } else if ((e.key === 'Enter' || e.key === ',') && query.value) {
    e.preventDefault()
    add(query.value)
  } else if (e.key === 'Backspace' && !text.value && model.value.length) {
    remove(model.value.length - 1)
  } else if ((e.key === 'ArrowDown' || e.key === 'ArrowUp') && rows.value) {
    e.preventDefault()
    if (!focused.value) focused.value = true
    active.value = (active.value + (e.key === 'ArrowDown' ? 1 : -1) + rows.value) % rows.value
    nextTick(() => list.value?.querySelector('.is-active')?.scrollIntoView({ block: 'nearest' }))
  } else if (e.key === 'Escape') {
    focused.value = false
  } else if (e.key === 'Tab') {
    focused.value = false
  }
}

function onBlur() {
  if (query.value) add(query.value)
}
</script>

<template>
  <div
    ref="box"
    v-bind="wrapper"
    :class="['x-field-box', 'x-tags', { 'x-field-box--mono': mono, 'is-open': open, 'is-disabled': disabled, 'is-invalid': invalid, 'is-modified': modified }]"
    @click="input?.focus()"
  >
    <span v-for="(tag, i) in model" :key="tag" class="x-tag-chip">
      <span class="x-tag-chip__text">{{ label(tag) }}</span>
      <button v-if="!disabled" type="button" class="x-tag-chip__remove" :aria-label="t('kit.remove', { value: label(tag) })" @click.stop="remove(i)">
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
      :placeholder="model.length ? morePlaceholder : placeholder"
      :disabled="disabled"
      autocomplete="off"
      spellcheck="false"
      @focus="focused = true"
      @blur="onBlur"
      @keydown="onKeydown"
    />
  </div>
  <Teleport to="body">
    <Transition name="x-popup">
      <ul v-if="open" :id="`${id}-list`" ref="list" role="listbox" aria-multiselectable="true" class="x-popup x-tags__list" :style="style">
        <li v-if="presetsTitle && shown.length" class="x-caps x-tags__title" role="presentation">{{ presetsTitle }}</li>
        <li
          v-for="(p, i) in shown"
          :id="`${id}-${i}`"
          :key="p.value"
          role="option"
          :aria-selected="model.includes(p.value)"
          :class="['x-popup-item', { 'is-active': i === active, 'is-selected': model.includes(p.value), 'x-tags__option--mono': mono && p.label === p.value }]"
          @pointermove="active = i"
          @pointerdown.prevent
          @click="toggle(p.value)"
        >
          <span class="x-tags__label">{{ p.label }}</span>
          <XIcon v-if="model.includes(p.value)" name="check" :size="14" :stroke="1.8" />
        </li>
        <li
          v-if="custom"
          :id="`${id}-${shown.length}`"
          role="option"
          :aria-selected="false"
          :class="['x-popup-item', 'x-tags__add', { 'is-active': active === shown.length }]"
          @pointermove="active = shown.length"
          @pointerdown.prevent
          @click="add(custom)"
        >
          <XIcon name="plus" :size="13" :stroke="1.7" />
          <span class="x-tags__label">{{ t('kit.addValue', { value: custom }) }}</span>
          <span class="x-tags__key">Enter</span>
        </li>
      </ul>
    </Transition>
  </Teleport>
</template>

<style scoped>
.x-tags {
  flex-wrap: wrap;
  gap: 4px;
  height: auto;
  min-height: var(--ctl-h);
  padding: 4px;
  cursor: text;
}

.x-tags__control {
  flex: 1;
  min-width: 110px;
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
  max-height: 264px;
  margin: 0;
  list-style: none;
}

.x-tags__title {
  padding: 6px 10px 5px;
}

.x-tags__label {
  flex: 1;
  min-width: 0;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
}

.x-tags__option--mono {
  font-family: var(--font-mono);
}

.x-tags__add {
  color: var(--acc);
}

.x-tags__key {
  flex: none;
  color: var(--dim);
  font: 400 var(--fs-10) / 1 var(--font-mono);
}
</style>
