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
}

/**
 * Action menu. The trigger is any button that binds the slot props:
 *   <XMenu :items="items" @select="run" v-slot="{ props }">
 *     <XButton v-bind="props" variant="ghost" icon="dots" label="More" />
 *   </XMenu>
 */
const props = withDefaults(defineProps<{ items: MenuItem[]; placement?: Placement; title?: string }>(), { placement: 'bottom-end' })
const emit = defineEmits<{ select: [key: string] }>()

const id = useId()
const anchor = ref<HTMLElement | null>(null)
const menu = ref<HTMLElement | null>(null)
const open = ref(false)
const { style } = useFloating(anchor, menu, open, { placement: props.placement, onDismiss: () => close(!!menu.value?.contains(document.activeElement)) })

function enabledItems() {
  return [...(menu.value?.querySelectorAll<HTMLElement>('[role^="menuitem"]:not([aria-disabled="true"])') ?? [])]
}

function focusItem(index: number) {
  const list = enabledItems()
  list[(index + list.length) % list.length]?.focus()
}

function show(first: 'first' | 'last' | 'none') {
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
    <Transition name="x-popup">
      <div
        v-if="open"
        :id="`${id}-menu`"
        ref="menu"
        role="menu"
        :aria-label="title"
        class="x-popup x-menu"
        :style="style"
        @keydown="onMenuKeydown"
      >
        <div v-if="title" class="x-menu__title" aria-hidden="true">{{ title }}</div>
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
.x-menu__title {
  padding: 7px 10px 6px;
  color: var(--dim);
  font: 500 var(--fs-10) / 1 var(--font-mono);
  letter-spacing: .09em;
  text-transform: uppercase;
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
