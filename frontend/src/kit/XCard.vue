<script setup lang="ts">
defineProps<{
  title?: string
  /** Card inside a card: smaller radius, darker surface. */
  nested?: boolean
  /** Body without padding, for tables and lists that reach the edges. */
  flush?: boolean
}>()
</script>

<template>
  <section :class="['x-card', { 'x-card--nested': nested }]">
    <header v-if="title || $slots.title || $slots.actions" class="x-card__head">
      <h2 class="x-card__title"><slot name="title">{{ title }}</slot></h2>
      <div v-if="$slots.actions" class="x-card__actions"><slot name="actions" /></div>
    </header>
    <div :class="['x-card__body', { 'x-card__body--flush': flush }]"><slot /></div>
  </section>
</template>

<style scoped>
.x-card {
  min-width: 0;
  border: 1px solid var(--line);
  border-radius: var(--r-16);
  background: var(--card);
}

.x-card--nested {
  border-radius: var(--r-12);
  background: var(--nest);
}

.x-card__head {
  display: flex;
  align-items: center;
  gap: 12px;
  min-height: 52px;
  padding: 10px 16px;
  border-bottom: 1px solid var(--div);
}

.x-card__title {
  flex: 1;
  min-width: 0;
  margin: 0;
  color: var(--fg);
  font: 600 var(--fs-14) / var(--lh-tight) var(--font-sans);
}

.x-card__actions {
  display: flex;
  flex: none;
  align-items: center;
  gap: 8px;
}

.x-card__body {
  padding: 16px;
}

.x-card__body--flush {
  padding: 0;
}
</style>
