<script setup lang="ts">
// Dashboard card: caps-label header with an optional right slot, then body.
defineProps<{ title?: string; fill?: boolean }>()
</script>

<template>
  <section class="x-card" :class="{ 'x-card--fill': fill }">
    <header v-if="title || $slots.aside" class="x-card__head">
      <h2 class="caps x-card__title">{{ title }}</h2>
      <div v-if="$slots.aside" class="x-card__aside"><slot name="aside" /></div>
    </header>
    <div class="x-card__body"><slot /></div>
    <footer v-if="$slots.footer" class="x-card__foot"><slot name="footer" /></footer>
  </section>
</template>

<style scoped>
.x-card {
  --card-px: var(--sp-5);
  display: flex;
  flex-direction: column;
  min-width: 0;
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--r-lg);
  box-shadow: var(--shadow-card);
  padding: 14px var(--card-px) 16px;
}
.x-card--fill {
  flex: 1;
  min-height: 0;
}
.x-card__head {
  display: flex;
  align-items: center;
  gap: 10px;
  min-height: var(--h-xs);
  margin: -4px calc(-1 * var(--card-px)) 13px;
  padding: 0 var(--card-px) 10px;
  border-bottom: 1px solid var(--divider);
}
.x-card__title {
  margin: 0;
  font-size: var(--fs-2xs);
}
.x-card__aside {
  margin-left: auto;
  display: flex;
  align-items: center;
  gap: var(--sp-2);
}
.x-card__body {
  display: flex;
  flex-direction: column;
  flex: 1;
  min-height: 0;
}
.x-card__foot {
  margin-top: auto;
  padding-top: 13px;
  border-top: 1px solid var(--divider);
}
.x-card__body + .x-card__foot {
  margin-top: 12px;
}

@media (max-width: 767px) {
  .x-card {
    --card-px: 15px;
    padding: 13px var(--card-px) 14px;
  }
}
</style>
