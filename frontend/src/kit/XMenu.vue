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
}

/**
 * Action menu. The trigger is any button that binds the slot props:
 *   <XMenu :items="items" @select="run" v-slot="{ props }">
 *     <XButton v-bind="props" variant="ghost" icon="dots" label="More" />
 *   </XMenu>
 */
const props = withDefaults(defineProps<{ items: MenuItem[]; placement?: Placement }>(), { placement: 'bottom-end' })
const emit = defineEmits<{ select: [key: string] }>()

const id = useId()
const anchor = ref<HTMLElement | null>(null)
const menu = ref<HTMLElement | null>(null)
const open = ref(false)
const { style } = useFloating(anchor, menu, open, { placement: props.placement, onDismiss: () => close(!!menu.value?.contains(document.activeElement)) })

function enabledItems() {
  return [...(menu.value?.querySelectorAll<HTMLElement>('[role="menuitem"]:not([aria-disabled="true"])') ?? [])]
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
      <div v-if="open" :id="`${id}-menu`" ref="menu" role="menu" class="x-popup x-menu" :style="style" @keydown="onMenuKeydown">
        <template v-for="item in items" :key="item.key">
          <div v-if="item.divided" class="x-menu__divider" role="separator" />
          <button
            type="button"
            role="menuitem"
            tabindex="-1"
            :aria-disabled="item.disabled || undefined"
            :class="['x-popup-item', { 'is-danger': item.danger }]"
            @click="choose(item)"
          >
            <XIcon v-if="item.icon" :name="item.icon" :size="14" :stroke="1.4" />
            <span>{{ item.label }}</span>
          </button>
        </template>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.x-menu__divider {
  height: 1px;
  margin: 4px 6px;
  background: var(--div);
}
</style>
