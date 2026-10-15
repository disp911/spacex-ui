<script setup lang="ts">
defineProps<{
  title?: string
  /** A line under the title. */
  description?: string
  /** The heading runs into the body with no line between, the actions at its bottom right. */
  seamless?: boolean
  /** Card inside a card: smaller radius, darker surface. */
  nested?: boolean
  /** Body without padding, for tables and lists that reach the edges. */
  flush?: boolean
}>()
</script>

<template>
  <section :class="['x-card', { 'x-card--nested': nested, 'x-card--seamless': seamless }]">
    <header v-if="title || $slots.title || $slots.actions" class="x-card__head">
      <div v-if="description" class="x-card__heading">
        <h2 class="x-card__title"><slot name="title">{{ title }}</slot></h2>
        <p class="x-card__description">{{ description }}</p>
      </div>
      <h2 v-else class="x-card__title"><slot name="title">{{ title }}</slot></h2>
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

.x-card__heading {
  flex: 1 1 0;
  min-width: 220px;
}

.x-card__description {
  margin: 3px 0 0;
  color: var(--mut);
  font: 400 var(--fs-12) / var(--lh-text) var(--font-sans);
}

.x-card--seamless .x-card__head {
  flex-wrap: wrap;
  align-items: flex-end;
  gap: 12px 20px;
  min-height: 0;
  padding: 16px 20px 0;
  border-bottom: 0;
}

.x-card--seamless .x-card__actions {
  max-width: 100%;
}

.x-card--seamless .x-card__body {
  padding: 14px 20px 20px;
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

@media (max-width: 767px) {
  .x-card--seamless .x-card__head {
    padding: 16px 16px 0;
  }

  .x-card--seamless .x-card__body {
    padding: 14px 16px 16px;
  }
}
</style>
