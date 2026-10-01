<script setup lang="ts">
import { computed, ref } from 'vue'

// Six-cell one-time code input backed by a single hidden-caret <input>,
// so paste, autofill (autocomplete=one-time-code) and backspace just work.
const model = defineModel<string>({ default: '' })
const props = withDefaults(defineProps<{ length?: number; id?: string }>(), { length: 6 })
defineEmits<{ complete: [string] }>()

const focused = ref(false)
const cells = computed(() => Array.from({ length: props.length }, (_, i) => model.value[i] ?? ''))
const active = computed(() => Math.min(model.value.length, props.length - 1))

function onInput(e: Event) {
  const el = e.target as HTMLInputElement
  model.value = el.value.replace(/\D/g, '').slice(0, props.length)
  el.value = model.value
}
</script>

<template>
  <div class="x-otp">
    <input
      :id="id"
      class="x-otp__native"
      :value="model"
      inputmode="numeric"
      autocomplete="one-time-code"
      :maxlength="length"
      @input="onInput"
      @focus="focused = true"
      @blur="focused = false"
      @keyup.enter="model.length === length && $emit('complete', model)"
    />
    <div
      v-for="(c, i) in cells"
      :key="i"
      class="x-otp__cell"
      :class="{ 'is-filled': c, 'is-active': focused && i === active }"
      aria-hidden="true"
    >
      {{ c }}
    </div>
  </div>
</template>

<style scoped>
.x-otp {
  position: relative;
  display: flex;
  gap: 7px;
}
.x-otp__native {
  position: absolute;
  inset: 0;
  width: 100%;
  opacity: 0;
  border: 0;
  cursor: text;
  z-index: 1;
  font-size: 16px;
}
.x-otp__cell {
  flex: 1;
  height: var(--h-lg);
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--r-md);
  background: var(--field);
  border: 1px solid var(--border-subtle);
  color: var(--text-5);
  font: var(--fw-medium) 15px var(--font-mono);
  transition: border-color var(--dur-fast) ease, box-shadow var(--dur-fast) ease;
}
.x-otp__cell.is-filled {
  border-color: var(--border-field);
  color: var(--text);
}
.x-otp__cell.is-active {
  border-color: var(--accent);
  box-shadow: 0 0 0 3px var(--accent-ring);
}
</style>
