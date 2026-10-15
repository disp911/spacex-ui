<script setup lang="ts">
import { computed, onBeforeUnmount, ref } from 'vue'
import { encode } from 'uqr'
import XIcon from './XIcon.vue'
import { copyText } from './clipboard'

/**
 * A QR code on a white plate, black on white in both themes so cameras read
 * it. With copy, a click puts that text on the clipboard: the corner badge
 * turns from a copy sign into a check for a moment.
 */
const props = withDefaults(defineProps<{ value: string; size?: number; copy?: string; copyLabel?: string }>(), { size: 180 })
const emit = defineEmits<{ copied: [ok: boolean] }>()

const qr = computed(() => encode(props.value, { ecc: 'L', border: 2 }))
const path = computed(() => {
  let d = ''
  qr.value.data.forEach((row, y) => row.forEach((on, x) => on && (d += `M${x} ${y}h1v1h-1z`)))
  return d
})

const done = ref(false)
let timer = 0
async function onCopy() {
  if (props.copy == null) return
  const ok = await copyText(props.copy)
  emit('copied', ok)
  if (!ok) return
  done.value = false
  requestAnimationFrame(() => (done.value = true))
  clearTimeout(timer)
  timer = window.setTimeout(() => (done.value = false), 1800)
}
onBeforeUnmount(() => clearTimeout(timer))
</script>

<template>
  <component
    :is="copy != null ? 'button' : 'div'"
    :type="copy != null ? 'button' : undefined"
    :class="['x-qr', { 'x-qr--copy': copy != null }]"
    :style="{ width: `${size}px`, height: `${size}px` }"
    :title="copyLabel"
    :aria-label="copyLabel"
    @click="onCopy"
  >
    <svg :viewBox="`0 0 ${qr.size} ${qr.size}`" shape-rendering="crispEdges" aria-hidden="true">
      <path :d="path" fill="#000" />
    </svg>
    <span v-if="copy != null" :class="['x-qr__badge', { 'is-done': done }]" aria-hidden="true">
      <XIcon :name="done ? 'check' : 'copy'" :size="14" :stroke="done ? 1.9 : 1.5" />
    </span>
  </component>
</template>

<style scoped>
.x-qr {
  position: relative;
  display: block;
  flex: none;
  padding: 0;
  border: 0;
  border-radius: var(--r-10);
  background: #fff;
  box-shadow: 0 0 0 1px var(--line);
}

.x-qr svg {
  display: block;
  width: 100%;
  height: 100%;
}

.x-qr--copy {
  cursor: pointer;
  outline: none;
  transition: box-shadow var(--t-med);
}

.x-qr--copy:hover,
.x-qr--copy:focus-visible {
  box-shadow: 0 0 0 2px var(--acc), 0 0 0 6px var(--acc-ring);
}

.x-qr__badge {
  position: absolute;
  right: -9px;
  bottom: -9px;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border-radius: var(--r-8);
  background: var(--acc);
  color: var(--acc-ink);
  box-shadow: 0 0 0 3px var(--dlg);
}

.x-qr__badge.is-done {
  animation: x-qr-pop .42s cubic-bezier(.22, 1, .36, 1);
}

@keyframes x-qr-pop {
  40% {
    transform: scale(1.12);
  }
}
</style>
