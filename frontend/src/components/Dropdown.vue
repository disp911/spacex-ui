<script setup lang="ts" generic="T extends string | number">
import { computed, onBeforeUnmount, ref, watch } from 'vue'
import { isMobile } from '../lib/media'
import Icon from './Icon.vue'

// A select: outlined trigger, glass popover with a blurring scrim behind it.
// On phones it opens as a bottom sheet instead.
export interface Option<V> {
  value: V
  label: string
  hint?: string
  dot?: string
}

const props = withDefaults(
  defineProps<{
    options: Option<T>[]
    /** Uppercase mono prefix inside the trigger, e.g. "СТРОК". */
    label?: string
    /** Heading of the bottom sheet on phones. */
    title?: string
    size?: 'sm' | 'md'
    align?: 'left' | 'right'
    menuWidth?: number
    /** Text for the trigger instead of the selected option's label. */
    display?: string
    /** Stretch the trigger to its container (phone filter rows). */
    block?: boolean
  }>(),
  { size: 'md', align: 'left' },
)
const model = defineModel<T>()
const open = ref(false)

const current = computed(() => props.options.find((o) => o.value === model.value))
const sheet = computed(() => isMobile.value)

function pick(v: T) {
  model.value = v
  open.value = false
}

function onKey(e: KeyboardEvent) {
  if (e.key === 'Escape') open.value = false
}
watch(open, (v) => (v ? document.addEventListener('keydown', onKey) : document.removeEventListener('keydown', onKey)))
onBeforeUnmount(() => document.removeEventListener('keydown', onKey))
</script>

<template>
  <div class="dd" :class="{ open: open && !sheet, block }">
    <button type="button" class="ctl trigger" :class="[size, { 'is-open': open }]" @click="open = !open">
      <span v-if="label" class="prefix">{{ label }}</span>
      <span v-if="current?.dot" class="odot" :style="{ background: current.dot }" />
      <span class="value">{{ display ?? current?.label }}</span>
      <span class="caret">{{ open ? '▲' : '▼' }}</span>
    </button>

    <template v-if="open && !sheet">
      <div class="scrim" @click="open = false" />
      <div class="menu" :class="align" :style="menuWidth ? { width: `${menuWidth}px` } : undefined" role="listbox">
        <button
          v-for="o in options"
          :key="o.value"
          type="button"
          class="row"
          :class="{ on: o.value === model }"
          role="option"
          :aria-selected="o.value === model"
          @click="pick(o.value)"
        >
          <span v-if="o.dot" class="odot" :style="{ background: o.dot }" />
          <span class="rl">{{ o.label }}</span>
          <span v-if="o.hint" class="rh">{{ o.hint }}</span>
          <Icon name="check" :size="12" class="rc" />
        </button>
      </div>
    </template>

    <Teleport v-if="open && sheet" to="body">
      <div class="sheet-scrim" @click="open = false" />
      <div class="sheet" role="listbox">
        <div class="grab" />
        <div v-if="title" class="sheet-title">{{ title }}</div>
        <div class="sheet-rows">
          <button
            v-for="o in options"
            :key="o.value"
            type="button"
            class="row"
            :class="{ on: o.value === model }"
            @click="pick(o.value)"
          >
            <span v-if="o.dot" class="odot" :style="{ background: o.dot }" />
            <span class="rl">{{ o.label }}</span>
            <span v-if="o.hint" class="rh">{{ o.hint }}</span>
            <Icon name="check" :size="15" class="rc" />
          </button>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<style scoped>
.dd {
  position: relative;
  flex: none;
}
.dd.open {
  z-index: 50;
}
.dd.block {
  flex: 1;
  min-width: 0;
}
.trigger {
  position: relative;
  z-index: 3;
  gap: 7px;
  padding: 0 9px 0 11px;
  background: var(--chip);
}
.block .trigger {
  width: 100%;
  justify-content: flex-start;
}
.trigger.md {
  height: 34px;
  font: 500 12.5px var(--font-sans);
}
.prefix {
  font: 500 9.5px var(--font-mono);
  letter-spacing: 0.07em;
  color: var(--muted);
  text-transform: uppercase;
}
.value {
  font: 500 12px var(--font-mono);
  color: var(--text);
  overflow: hidden;
  text-overflow: ellipsis;
}
.trigger.sm .value {
  font-size: 10.5px;
  color: inherit;
}
.block .caret {
  margin-left: auto;
}
.odot {
  width: 6px;
  height: 6px;
  border-radius: 3px;
  flex: none;
}
.scrim {
  position: fixed;
  inset: 0;
  z-index: 1;
  background: var(--scrim);
  backdrop-filter: blur(3px) saturate(0.92);
  animation: scrimIn 0.16s ease-out;
}
.menu {
  position: absolute;
  top: calc(100% + 5px);
  z-index: 2;
  min-width: 100%;
  max-height: 320px;
  overflow-y: auto;
  background: var(--popover);
  border: 1px solid var(--line-card);
  border-radius: 10px;
  backdrop-filter: blur(10px);
  box-shadow: var(--shadow-pop);
  padding: 5px;
  display: flex;
  flex-direction: column;
  gap: 1px;
  animation: sheetUp 0.14s ease-out;
}
.menu.left {
  left: 0;
}
.menu.right {
  right: 0;
}
.row {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 8px;
  border-radius: 7px;
  color: var(--text-2);
  text-align: left;
  white-space: nowrap;
  transition: background 0.16s ease;
}
.row:hover {
  background: var(--row-hover);
}
.row.on {
  background: var(--accent-soft);
  color: var(--accent);
}
.rl {
  flex: 1;
  min-width: 0;
  font: 450 11.5px var(--font-sans);
  overflow: hidden;
  text-overflow: ellipsis;
}
.rh {
  font: 400 10.5px var(--font-mono);
  color: var(--muted);
}
.row.on .rh {
  color: inherit;
  opacity: 0.75;
}
.rc {
  color: transparent;
}
.row.on .rc {
  color: var(--accent);
}

/* Bottom sheet (phones) */
.sheet-scrim {
  position: fixed;
  inset: 0;
  z-index: 300;
  background: var(--scrim-strong);
  backdrop-filter: blur(2px);
  animation: scrimIn 0.18s ease-out;
}
.sheet {
  position: fixed;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 301;
  max-height: 75vh;
  display: flex;
  flex-direction: column;
  background: var(--sheet);
  border-top: 1px solid var(--line-card);
  border-radius: 22px 22px 0 0;
  box-shadow: var(--shadow-sheet);
  padding: 8px 12px max(14px, env(safe-area-inset-bottom));
  animation: sheetUp 0.16s ease-out;
}
.grab {
  width: 38px;
  height: 4px;
  border-radius: 2px;
  background: var(--scroll-thumb);
  margin: 2px auto 10px;
  flex: none;
}
.sheet-title {
  font: 600 15px var(--font-sans);
  padding: 2px 11px 10px;
}
.sheet-rows {
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 2px;
}
.sheet .row {
  padding: 12px 11px;
  border-radius: 11px;
}
.sheet .rl {
  font-size: 14px;
}
.sheet .rh {
  font-size: 12px;
}
</style>
