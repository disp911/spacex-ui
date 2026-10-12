<script setup lang="ts">
import { nextTick, ref, useId } from 'vue'
import XIcon from './XIcon.vue'
import { useFloating, type Placement } from './floating'
import type { IconName } from './icons'

export interface MenuItem {
  key: string
  label: string
  icon?: IconName
  danger?: boolean
  disabled?: boolean
  /** Thin line above the item. */
  divided?: boolean
  /** Short text at the right edge, e.g. a language code. */
  hint?: string
  /** Makes the item a choice: a check marks the chosen one. */
  checked?: boolean
  /** Colour of a small dot before the label, e.g. a log level. */
  dot?: string
}

/**
 * Action menu. The trigger is any button that binds the slot props:
 *   <XMenu :items="items" @select="run" v-slot="{ props }">
 *     <XButton v-bind="props" variant="ghost" icon="dots" label="More" />
 *   </XMenu>
 */
const props = withDefaults(
  defineProps<{
    items: MenuItem[]
    placement?: Placement
    title?: string
    /** On phones opens as a bottom sheet over a scrim: filters in windows. */
    sheet?: boolean
  }>(),
  { placement: 'bottom-end' },
)
const emit = defineEmits<{ select: [key: string] }>()

const id = useId()
const anchor = ref<HTMLElement | null>(null)
const menu = ref<HTMLElement | null>(null)
const open = ref(false)
const asSheet = ref(false)
const { style } = useFloating(anchor, menu, open, { placement: props.placement, onDismiss: () => close(!!menu.value?.contains(document.activeElement)) })

function enabledItems() {
  return [...(menu.value?.querySelectorAll<HTMLElement>('[role^="menuitem"]:not([aria-disabled="true"])') ?? [])]
}

function focusItem(index: number) {
  const list = enabledItems()
  list[(index + list.length) % list.length]?.focus()
}

function show(first: 'first' | 'last' | 'none') {
  asSheet.value = props.sheet && matchMedia('(max-width: 767px)').matches
  open.value = true
  if (first !== 'none') nextTick(() => focusItem(first === 'first' ? 0 : -1))
}

function close(refocus: boolean) {
  open.value = false
  if (refocus) anchor.value?.focus()
}

function onTriggerClick(e: MouseEvent) {
  anchor.value = e.currentTarget as HTMLElement
  if (open.value) close(false)
  else show(e.detail === 0 ? 'first' : 'none')
}

function onTriggerKeydown(e: KeyboardEvent) {
  anchor.value = e.currentTarget as HTMLElement
  if (e.key === 'ArrowDown' || e.key === 'ArrowUp') {
    e.preventDefault()
    show(e.key === 'ArrowDown' ? 'first' : 'last')
  }
}

function onMenuKeydown(e: KeyboardEvent) {
  const list = enabledItems()
  const i = list.indexOf(document.activeElement as HTMLElement)
  if (e.key === 'ArrowDown' || e.key === 'ArrowUp') {
    e.preventDefault()
    focusItem(i + (e.key === 'ArrowDown' ? 1 : -1))
  } else if (e.key === 'Home' || e.key === 'End') {
    e.preventDefault()
    focusItem(e.key === 'Home' ? 0 : -1)
  } else if (e.key === 'Tab') {
    close(false)
  }
}

function choose(item: MenuItem) {
  if (item.disabled) return
  close(true)
  emit('select', item.key)
}
</script>

<template>
  <slot
    :open="open"
    :props="{
      'aria-haspopup': 'menu',
      'aria-expanded': open,
      'aria-controls': `${id}-menu`,
      onClick: onTriggerClick,
      onKeydown: onTriggerKeydown,
    }"
  />
  <Teleport to="body">
    <Transition name="x-menu-scrim">
      <div v-if="open && asSheet" class="x-menu-scrim" aria-hidden="true" />
    </Transition>
    <Transition :name="asSheet ? 'x-menu-sheet' : 'x-popup'">
      <div
        v-if="open"
        :id="`${id}-menu`"
        ref="menu"
        role="menu"
        :aria-label="title"
        :class="['x-popup', 'x-menu', { 'x-menu--sheet': asSheet }]"
        :style="asSheet ? undefined : style"
        @keydown="onMenuKeydown"
      >
        <div v-if="title" class="x-caps x-menu__title" aria-hidden="true">{{ title }}</div>
        <template v-for="item in items" :key="item.key">
          <div v-if="item.divided" class="x-menu__divider" role="separator" />
          <button
            type="button"
            :role="item.checked === undefined ? 'menuitem' : 'menuitemradio'"
            tabindex="-1"
            :aria-checked="item.checked"
            :aria-disabled="item.disabled || undefined"
            :class="['x-popup-item', { 'is-danger': item.danger, 'is-selected': item.checked }]"
            @click="choose(item)"
          >
            <XIcon v-if="item.icon" :name="item.icon" :size="14" :stroke="1.4" />
            <span v-if="item.dot" class="x-menu__dot" :style="{ background: item.dot }" />
            <span class="x-menu__label">{{ item.label }}</span>
            <span v-if="item.hint" class="x-menu__hint">{{ item.hint }}</span>
            <XIcon v-if="item.checked !== undefined" name="check" :size="13" :stroke="1.8" :class="['x-menu__check', { 'is-on': item.checked }]" />
          </button>
        </template>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.x-menu--sheet {
  position: fixed;
  right: 0;
  bottom: 0;
  left: 0;
  gap: 2px;
  max-height: calc(100% - 48px);
  padding: 10px 12px calc(16px + env(safe-area-inset-bottom));
  border-width: 1px 0 0;
  border-radius: var(--r-22) var(--r-22) 0 0;
  box-shadow: var(--shadow-sheet);
}

.x-menu--sheet::before {
  content: '';
  flex: none;
  align-self: center;
  width: 38px;
  height: 4px;
  margin-bottom: 6px;
  border-radius: var(--r-pill);
  background: var(--track);
}

.x-menu--sheet .x-menu__title {
  margin-bottom: 6px;
  padding: 0 6px;
}

.x-menu--sheet .x-popup-item {
  padding: 0 12px;
  border-radius: var(--r-12);
}

.x-menu--sheet .x-menu__check {
  width: 14px;
  height: 14px;
}

.x-menu-scrim {
  position: fixed;
  inset: 0;
  z-index: 999;
  background: var(--scrim);
  backdrop-filter: blur(3px) saturate(.92);
}

.x-menu-scrim-enter-active,
.x-menu-scrim-leave-active {
  transition: opacity var(--t-med);
}

.x-menu-scrim-enter-from,
.x-menu-scrim-leave-to {
  opacity: 0;
}

.x-menu-sheet-enter-active,
.x-menu-sheet-leave-active {
  transition: transform var(--t-med) var(--ease);
}

.x-menu-sheet-enter-from,
.x-menu-sheet-leave-to {
  transform: translateY(100%);
}

.x-menu__title {
  padding: 8px 10px 6px;
}

.x-menu__dot {
  flex: none;
  width: 6px;
  height: 6px;
  border-radius: 50%;
}

.x-menu__label {
  flex: 1;
  min-width: 0;
}

.x-menu__hint {
  color: var(--dim);
  font: 500 var(--fs-10) / 1 var(--font-mono);
}

.x-menu__check {
  visibility: hidden;
}

.x-menu__check.is-on {
  visibility: visible;
}

.x-menu__divider {
  height: 1px;
  margin: 4px 6px;
  background: var(--div);
}
</style>
