<script setup lang="ts">
/**
 * Facts in a frame, one per line with every other line tinted: a name and a
 * value in mono that breaks anywhere, e.g. the keys of an account. On phones
 * the name goes over its value.
 */
defineProps<{ rows: { label: string; value: string }[] }>()
</script>

<template>
  <dl class="x-kv">
    <div v-for="r in rows" :key="r.label" class="x-kv__row">
      <dt class="x-kv__label">{{ r.label }}</dt>
      <dd class="x-kv__value">{{ r.value }}</dd>
    </div>
  </dl>
</template>

<style scoped>
.x-kv {
  margin: 0;
  overflow: clip;
  border: 1px solid var(--div);
  border-radius: var(--r-10);
}

.x-kv__row {
  display: grid;
  grid-template-columns: 130px minmax(0, 1fr);
  align-items: baseline;
  gap: var(--sp-12);
  padding: var(--sp-8) var(--sp-12);
}

.x-kv__row + .x-kv__row {
  border-top: 1px solid var(--div);
}

.x-kv__row:nth-child(even) {
  background: var(--tint);
}

.x-kv__label {
  color: var(--mut);
  font: 400 var(--fs-12) / var(--lh-text) var(--font-sans);
}

.x-kv__value {
  margin: 0;
  color: var(--fg2);
  font: 400 var(--fs-12) / var(--lh-text) var(--font-mono);
  word-break: break-all;
}

@media (max-width: 767px) {
  .x-kv {
    border-radius: var(--r-12);
  }

  .x-kv__row {
    grid-template-columns: minmax(0, 1fr);
    gap: var(--sp-4);
  }
}
</style>
