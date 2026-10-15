<script setup lang="ts">
/**
 * One item of a list inside an XNest, e.g. a noise: a title with a blue dot
 * when changed, actions at the right (a remove button) and its fields; with
 * columns, the fields fill a grid of at least that many px each.
 */
defineProps<{ title: string; modified?: boolean; columns?: number }>()
</script>

<template>
  <div class="x-nest-card">
    <div class="x-nest-card__head">
      <span class="x-nest-card__title">{{ title }}</span>
      <span v-if="modified" class="x-nest-card__dot" aria-hidden="true" />
      <div v-if="$slots.actions" class="x-nest-card__actions"><slot name="actions" /></div>
    </div>
    <div :class="{ 'x-nest-card__grid': columns }" :style="columns ? { '--nest-col': `${columns}px` } : undefined"><slot /></div>
  </div>
</template>

<style scoped>
.x-nest-card {
  min-width: 0;
  padding: 10px 12px 12px;
  border: 1px solid var(--div);
  border-radius: var(--r-12);
  background: var(--card);
}

.x-nest-card__head {
  display: flex;
  align-items: center;
  gap: 8px;
  min-height: 26px;
  margin-bottom: 8px;
}

.x-nest-card__title {
  color: var(--fg);
  font: 600 var(--fs-13) / var(--lh-tight) var(--font-sans);
}

.x-nest-card__dot {
  flex: none;
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--mod);
}

.x-nest-card__grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(min(var(--nest-col), 100%), 1fr));
  gap: 10px 12px;
}

.x-nest-card__actions {
  display: flex;
  gap: 6px;
  margin-left: auto;
}
</style>
