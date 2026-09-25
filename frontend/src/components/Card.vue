<script setup lang="ts">
defineProps<{ label: string }>()
</script>

<template>
  <section class="card">
    <header class="card-head">
      <div class="caps">{{ label }}</div>
      <div v-if="$slots.extra" class="extra"><slot name="extra" /></div>
    </header>
    <slot />
  </section>
</template>

<style scoped>
.card {
  background: var(--card);
  border: 1px solid var(--line-card);
  border-radius: 16px;
  padding: 14px 20px;
  display: flex;
  flex-direction: column;
  min-width: 0;
}
.card-head {
  display: flex;
  align-items: center;
  gap: 10px;
  min-height: 14px;
  padding: 0 20px 14px;
  margin: 0 -20px 13px;
  border-bottom: 1px solid var(--line);
}
.extra {
  margin-left: auto;
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: -5px;
  margin-bottom: -5px;
}
@media (max-width: 760px) {
  .card {
    padding: 13px 15px;
  }
  .card-head {
    padding: 0 15px 13px;
    margin: 0 -15px 13px;
  }
}
</style>
