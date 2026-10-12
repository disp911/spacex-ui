<script setup lang="ts">
import { ref, useId } from 'vue'
import { useFloating } from './floating'

/** Short hint shown on hover and keyboard focus. */
const props = withDefaults(defineProps<{ text?: string; delay?: number }>(), { delay: 300 })

const id = useId()
const anchor = ref<HTMLElement | null>(null)
const tip = ref<HTMLElement | null>(null)
const open = ref(false)
const { style } = useFloating(anchor, tip, open, { placement: 'top', onDismiss: hide })
let timer = 0

function show() {
  clearTimeout(timer)
  timer = window.setTimeout(() => (open.value = true), props.delay)
}

function hide() {
  clearTimeout(timer)
  open.value = false
}
</script>

<template>
  <span
    ref="anchor"
    class="x-tooltip-anchor"
    :aria-describedby="open ? id : undefined"
    @mouseenter="show"
    @mouseleave="hide"
    @focusin="show"
    @focusout="hide"
  >
    <slot />
  </span>
  <Teleport to="body">
    <Transition name="x-tooltip">
      <div v-if="open" :id="id" ref="tip" role="tooltip" class="x-tooltip" :style="style">
        <slot name="content">{{ text }}</slot>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.x-tooltip-anchor {
  display: inline-flex;
  min-width: 0;
}

.x-tooltip {
  z-index: 1100;
  max-width: 280px;
  padding: 6px 10px;
  border: 1px solid var(--line);
  border-radius: var(--r-8);
  background: var(--dlg);
  color: var(--fg2);
  box-shadow: var(--shadow-pop);
  font: 400 var(--fs-12) / var(--lh-text) var(--font-sans);
  pointer-events: none;
}

.x-tooltip-enter-active,
.x-tooltip-leave-active {
  transition: opacity var(--t-fast);
}

.x-tooltip-enter-from,
.x-tooltip-leave-to {
  opacity: 0;
}
</style>
