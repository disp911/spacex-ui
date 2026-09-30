<script setup lang="ts">
import { computed, ref } from 'vue'
import XIcon from './XIcon.vue'
import XSheet from './XSheet.vue'
import { useEscape, useIsMobile } from '@/composables/useMedia'
import type { MenuItem } from './types'

// "⋯" action menu. Desktop: popover. Mobile: bottom sheet.
const props = withDefaults(
  defineProps<{ items: MenuItem[]; label: string; title?: string; align?: 'left' | 'right'; size?: 'xs' | 'sm' | 'md' }>(),
  { align: 'right', size: 'sm' },
)
const emit = defineEmits<{ select: [string] }>()

const open = ref(false)
const isMobile = useIsMobile()
useEscape(() => (open.value = false))

const btnSize = computed(() => ({ xs: 'var(--h-xs)', sm: 'var(--h-sm)', md: 'var(--h-md)' })[props.size])

function pick(item: MenuItem) {
  if (item.disabled) return
  open.value = false
  if (!item.href) emit('select', item.key)
}
</script>

<template>
  <div class="x-menu" @click.stop>
    <button
      type="button"
      class="x-menu__trigger"
      :class="{ 'is-open': open }"
      :style="{ width: btnSize, height: btnSize }"
      :aria-label="label"
      :title="label"
      aria-haspopup="menu"
      :aria-expanded="open"
      @click="open = !open"
    >
      <svg width="14" height="14" viewBox="0 0 16 16" aria-hidden="true">
        <circle cx="3.5" cy="8" r="1.3" fill="currentColor" />
        <circle cx="8" cy="8" r="1.3" fill="currentColor" />
        <circle cx="12.5" cy="8" r="1.3" fill="currentColor" />
      </svg>
    </button>

    <template v-if="!isMobile">
      <div v-if="open" class="x-menu__scrim" @click="open = false" />
      <Transition name="x-pop">
        <div v-if="open" class="x-menu__list" :class="`is-${align}`" role="menu">
          <template v-for="it in items" :key="it.key">
            <div v-if="it.divided" class="x-menu__sep" />
            <component
              :is="it.href ? 'a' : 'button'"
              :href="it.href"
              :type="it.href ? undefined : 'button'"
              role="menuitem"
              class="x-menu__item"
              :class="{ 'is-danger': it.danger }"
              :disabled="it.disabled"
              @click="pick(it)"
            >
              <XIcon v-if="it.icon" :name="it.icon" :size="14" />
              <span>{{ it.label }}</span>
            </component>
          </template>
        </div>
      </Transition>
    </template>

    <XSheet v-else v-model:open="open" :title="title">
      <template v-for="it in items" :key="it.key">
        <div v-if="it.divided" class="x-menu__sep" />
        <component
          :is="it.href ? 'a' : 'button'"
          :href="it.href"
          :type="it.href ? undefined : 'button'"
          class="x-menu__item x-menu__item--lg"
          :class="{ 'is-danger': it.danger }"
          :disabled="it.disabled"
          @click="pick(it)"
        >
          <XIcon v-if="it.icon" :name="it.icon" :size="16" />
          <span>{{ it.label }}</span>
        </component>
      </template>
    </XSheet>
  </div>
</template>

<style scoped>
.x-menu {
  position: relative;
  flex: none;
}
.x-menu__trigger {
  position: relative;
  z-index: 6;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--r-sm);
  border: 1px solid var(--border-control);
  background: var(--control);
  color: var(--text-3);
  transition: border-color var(--dur-base) ease, background var(--dur-base) ease, color var(--dur-fast) ease, box-shadow var(--dur-base) ease;
}
.x-menu__trigger:hover,
.x-menu__trigger.is-open {
  border-color: var(--accent-focus);
  background: var(--accent-soft-hover);
  box-shadow: 0 0 0 3px var(--accent-ring);
  color: var(--text);
}
.x-menu__scrim {
  position: fixed;
  inset: 0;
  z-index: 4;
  background: var(--scrim-soft);
  backdrop-filter: blur(3px) saturate(0.92);
  animation: spx-fade-in 0.16s ease-out;
}
.x-menu__list {
  position: absolute;
  top: calc(100% + 5px);
  z-index: 7;
  min-width: 220px;
  padding: 5px;
  display: flex;
  flex-direction: column;
  gap: 1px;
  background: var(--surface-raised);
  border: 1px solid var(--border);
  border-radius: var(--r-md);
  box-shadow: var(--shadow-pop);
  backdrop-filter: blur(10px);
}
.x-menu__list.is-right { right: 0; }
.x-menu__list.is-left { left: 0; }
.x-menu__item {
  display: flex;
  align-items: center;
  gap: 9px;
  width: 100%;
  padding: 7px 9px;
  border: 0;
  border-radius: var(--r-sm);
  background: transparent;
  color: var(--text-2);
  font: var(--fw-regular) var(--fs-sm) var(--font-sans);
  text-align: left;
  text-decoration: none;
  white-space: nowrap;
}
.x-menu__item :deep(svg) {
  color: var(--text-3);
}
.x-menu__item:hover {
  background: var(--control-hover);
  color: var(--text);
}
.x-menu__item.is-danger,
.x-menu__item.is-danger :deep(svg) {
  color: var(--danger-text);
}
.x-menu__item.is-danger:hover {
  background: var(--danger-soft);
}
.x-menu__item:disabled {
  opacity: 0.4;
  cursor: not-allowed;
  background: transparent;
}
.x-menu__item--lg {
  gap: 11px;
  padding: 12px 11px;
  border-radius: var(--r-md);
  font-size: var(--fs-lg);
}
.x-menu__sep {
  height: 1px;
  margin: 4px 8px;
  background: var(--divider);
}
.x-pop-enter-active {
  transition: opacity 0.14s ease-out, transform 0.14s ease-out;
}
.x-pop-enter-from {
  opacity: 0;
  transform: translateY(6px);
}
</style>
