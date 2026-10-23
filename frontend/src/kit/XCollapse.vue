<script setup lang="ts">
/**
 * Content that opens and closes in place, e.g. an inbound's clients or a
 * section's body: it grows from nothing to its height and fades in, and
 * shrinks back to close. lazy builds the content only while it is open;
 * otherwise it stays in the page, hidden, and keeps what was typed in it.
 */
defineProps<{ open: boolean; lazy?: boolean }>()

const ENTER = 'height var(--t-enter) var(--ease-out), opacity var(--t-enter) var(--ease-out)'
const LEAVE = 'height var(--t-leave) ease-in, opacity var(--t-leave) ease-in'

// Ends on the height's transition, or after a while if none runs (reduced motion, nothing to grow).
function run(el: HTMLElement, done: () => void) {
  const finish = () => {
    clearTimeout(timer)
    el.removeEventListener('transitionend', onEnd)
    done()
  }
  const onEnd = (e: TransitionEvent) => {
    if (e.target === el && e.propertyName === 'height') finish()
  }
  const timer = setTimeout(finish, 400)
  el.addEventListener('transitionend', onEnd)
}

function enter(node: Element, done: () => void) {
  const el = node as HTMLElement
  el.style.overflow = 'clip'
  // Its full height with every margin inside: scrollHeight would leave out the last one.
  const height = el.offsetHeight
  el.style.height = '0'
  el.style.opacity = '0'
  void el.offsetHeight
  el.style.transition = ENTER
  el.style.height = `${height}px`
  el.style.opacity = '1'
  run(el, done)
}

function leave(node: Element, done: () => void) {
  const el = node as HTMLElement
  el.style.overflow = 'clip'
  el.style.height = `${el.offsetHeight}px`
  void el.offsetHeight
  el.style.transition = LEAVE
  el.style.height = '0'
  el.style.opacity = '0'
  run(el, done)
}

function settle(node: Element) {
  const el = node as HTMLElement
  el.style.height = ''
  el.style.opacity = ''
  el.style.overflow = ''
  el.style.transition = ''
}
</script>

<template>
  <Transition :css="false" @enter="enter" @after-enter="settle" @enter-cancelled="settle" @leave="leave" @after-leave="settle" @leave-cancelled="settle">
    <!-- One directive each: on the same element both would start the same entrance twice. -->
    <div v-if="lazy && open" class="x-collapse"><slot /></div>
    <div v-else-if="!lazy" v-show="open" class="x-collapse"><slot /></div>
  </Transition>
</template>

<style scoped>
/* Holds its content's margins inside, so the height it grows to is the whole of it; it clips only while it moves. */
.x-collapse {
  display: flow-root;
}
</style>
