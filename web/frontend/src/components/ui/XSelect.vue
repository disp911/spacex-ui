<script setup lang="ts" generic="V extends string | number">
import { computed, ref } from 'vue'
import XIcon from './XIcon.vue'
import XSheet from './XSheet.vue'
import { useEscape, useIsMobile } from '@/composables/useMedia'
import type { SelectOption } from './types'

// Compact dropdown. Desktop: popover under the trigger. Mobile: bottom sheet.
const model = defineModel<V>({ required: true })
const props = withDefaults(
  defineProps<{
    options: SelectOption<V>[]
    /** Caps prefix inside the trigger, e.g. "ДАТА". */
    label?: string
    /** Title of the mobile sheet (defaults to label). */
    sheetTitle?: string
    size?: 'xs' | 'sm' | 'md'
    block?: boolean
    mono?: boolean
    menuWidth?: number
    align?: 'left' | 'right'
    /** Always use the popover, even on mobile. */
    inline?: boolean
  }>(),
  { size: 'sm', align: 'left' },
)

const open = ref(false)
const isMobile = useIsMobile()
const asSheet = computed(() => isMobile.value && !props.inline)
const current = computed(() => props.options.find((o) => o.value === model.value))

function pick(v: V) {
  model.value = v
  open.value = false
}
useEscape(() => (open.value = false))
</script>

<template>
  <div class="x-select" :class="{ 'x-select--block': block }">
    <button
      type="button"
      class="x-select__trigger"
      :class="[`x-select__trigger--${size}`, { 'is-open': open }]"
      :aria-expanded="open"
      aria-haspopup="listbox"
      @click="open = !open"
    >
      <span v-if="label && !(block && isMobile)" class="caps x-select__label">{{ label }}</span>
      <span v-if="current?.dot" class="x-select__dot" :style="{ background: current.dot }" />
      <span class="x-select__value" :class="{ mono }"><slot name="value" :option="current">{{ current?.label ?? '—' }}</slot></span>
      <XIcon name="chevron-down" :size="10" :stroke-width="1.6" class="x-select__caret" />
    </button>

    <template v-if="!asSheet">
      <div v-if="open" class="x-select__scrim" @click="open = false" />
      <Transition name="x-pop">
        <ul
          v-if="open"
          class="x-select__menu"
          :class="`is-${align}`"
          role="listbox"
          :style="{ width: menuWidth ? menuWidth + 'px' : undefined }"
        >
          <li
            v-for="o in options"
            :key="String(o.value)"
            role="option"
            :aria-selected="o.value === model"
            class="x-select__opt"
            :class="{ 'is-on': o.value === model }"
            @click="pick(o.value)"
          >
            <span v-if="o.dot" class="x-select__dot" :style="{ background: o.dot }" />
            <span class="x-select__opt-label" :class="{ mono }">{{ o.label }}</span>
            <span v-if="o.meta" class="x-select__opt-meta">{{ o.meta }}</span>
            <XIcon name="check" :size="12" :stroke-width="1.9" class="x-select__check" />
          </li>
        </ul>
      </Transition>
    </template>

    <XSheet v-else v-model:open="open" :title="sheetTitle ?? label">
      <div
        v-for="o in options"
        :key="String(o.value)"
        class="x-select__sheet-opt"
        :class="{ 'is-on': o.value === model }"
        role="option"
        :aria-selected="o.value === model"
        @click="pick(o.value)"
      >
        <span v-if="o.dot" class="x-select__dot" :style="{ background: o.dot }" />
        <span class="x-select__opt-label" :class="{ mono }">{{ o.label }}</span>
        <span v-if="o.meta" class="x-select__opt-meta">{{ o.meta }}</span>
        <XIcon name="check" :size="14" :stroke-width="1.9" class="x-select__check" />
      </div>
    </XSheet>
  </div>
</template>

<style scoped>
.x-select {
  position: relative;
  flex: none;
}
.x-select--block {
  flex: 1;
  min-width: 0;
}
.x-select__trigger {
  position: relative;
  z-index: 6;
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
  height: var(--h-sm);
  padding: 0 9px 0 11px;
  border-radius: var(--r-sm);
  border: 1px solid var(--border-control);
  background: var(--control);
  color: var(--text-2);
  transition: border-color var(--dur-base) ease, background var(--dur-base) ease, box-shadow var(--dur-base) ease;
}
.x-select__trigger--xs {
  height: var(--h-xs);
  padding: 0 7px 0 9px;
  gap: 6px;
}
.x-select__trigger--md {
  height: var(--h-md);
  border-radius: var(--r-md);
}
.x-select__trigger:hover,
.x-select__trigger.is-open {
  border-color: var(--accent-focus);
  background: var(--accent-soft-hover);
  box-shadow: 0 0 0 3px var(--accent-ring);
  color: var(--text);
}
.x-select__label {
  font-size: var(--fs-2xs);
}
.x-select__value {
  font: var(--fw-medium) var(--fs-sm) var(--font-sans);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  min-width: 0;
}
.x-select--block .x-select__value {
  flex: 1;
  text-align: left;
}
.x-select__trigger--xs .x-select__value {
  font-size: var(--fs-xs);
}
.mono {
  font-family: var(--font-mono);
}
.x-select__caret {
  color: var(--text-4);
  transition: transform var(--dur-fast) ease;
}
.is-open .x-select__caret {
  transform: rotate(180deg);
}
.x-select__dot {
  width: 6px;
  height: 6px;
  flex: none;
  border-radius: 50%;
}

.x-select__scrim {
  position: fixed;
  inset: 0;
  z-index: 4;
  background: var(--scrim-soft);
  backdrop-filter: blur(3px) saturate(0.92);
  animation: spx-fade-in 0.16s ease-out;
}
.x-select__menu {
  position: absolute;
  top: calc(100% + 5px);
  z-index: 7;
  min-width: 100%;
  max-height: 280px;
  overflow-y: auto;
  margin: 0;
  padding: 5px;
  list-style: none;
  background: var(--surface-raised);
  border: 1px solid var(--border);
  border-radius: var(--r-md);
  box-shadow: var(--shadow-pop);
  backdrop-filter: blur(10px);
}
.x-select__menu.is-left { left: 0; }
.x-select__menu.is-right { right: 0; }
.x-select__opt,
.x-select__sheet-opt {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 8px;
  border-radius: var(--r-sm);
  color: var(--text-2);
  cursor: pointer;
  white-space: nowrap;
}
.x-select__opt:hover {
  background: var(--control-hover);
}
.x-select__opt.is-on,
.x-select__sheet-opt.is-on {
  background: var(--accent-soft);
  color: var(--accent);
}
.x-select__opt-label {
  flex: 1;
  min-width: 0;
  font: var(--fw-regular) var(--fs-sm) var(--font-sans);
  overflow: hidden;
  text-overflow: ellipsis;
}
.x-select__opt-label.mono {
  font-family: var(--font-mono);
  font-weight: var(--fw-medium);
}
.x-select__opt-meta {
  font: var(--fw-regular) var(--fs-xs) var(--font-mono);
  color: var(--text-4);
}
.x-select__check {
  color: transparent;
}
.is-on .x-select__check {
  color: var(--accent);
}
.x-select__sheet-opt {
  gap: 10px;
  padding: 12px 11px;
  border-radius: var(--r-md);
}
.x-select__sheet-opt .x-select__opt-label {
  font-size: var(--fs-lg);
}

.x-pop-enter-active {
  transition: opacity 0.14s ease-out, transform 0.14s ease-out;
}
.x-pop-enter-from {
  opacity: 0;
  transform: translateY(6px);
}
</style>
