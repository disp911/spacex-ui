<script setup lang="ts" generic="T extends string">
// Equal-width option buttons, one selected (e.g. geosite | geoip).
const model = defineModel<T>({ required: true })
defineProps<{ options: readonly T[]; mono?: boolean }>()
</script>

<template>
  <div class="x-seg" role="radiogroup">
    <button
      v-for="o in options"
      :key="o"
      type="button"
      role="radio"
      :aria-checked="model === o"
      class="x-seg__opt"
      :class="{ 'is-on': model === o, mono }"
      @click="model = o"
    >
      {{ o }}
    </button>
  </div>
</template>

<style scoped>
.x-seg {
  display: flex;
  gap: 7px;
}
.x-seg__opt {
  flex: 1;
  height: var(--h-md);
  border-radius: var(--r-md);
  border: 1px solid var(--border-field);
  background: var(--field);
  color: var(--text-3);
  font: var(--fw-medium) var(--fs-md) var(--font-sans);
  transition: border-color var(--dur-base) ease, color var(--dur-base) ease, box-shadow var(--dur-base) ease;
}
.x-seg__opt.mono {
  font-family: var(--font-mono);
}
.x-seg__opt:hover {
  border-color: var(--accent-focus);
  color: var(--text-2);
}
.x-seg__opt.is-on {
  border-color: var(--accent);
  background: var(--accent-soft);
  color: var(--accent);
  box-shadow: 0 0 0 3px var(--accent-ring);
}
</style>
