<script setup lang="ts">
/** Dashboard card: a small caps title, optional header content and the body. tool: the head holds a button. */
defineProps<{ title: string; tool?: boolean }>()
</script>

<template>
  <section class="dash-card">
    <header :class="['dash-card__head', { 'dash-card__head--tool': tool }]">
      <h2 class="x-caps x-caps--lg dash-card__title">{{ title }}</h2>
      <div v-if="$slots.head" class="dash-card__extra"><slot name="head" /></div>
    </header>
    <slot />
  </section>
</template>

<style scoped>
.dash-card {
  display: flex;
  flex-direction: column;
  min-width: 0;
  padding: 14px 20px;
  border: 1px solid var(--line);
  border-radius: var(--r-16);
  background: var(--card);
}

.dash-card__head {
  display: flex;
  align-items: center;
  gap: 10px;
  /* One 14px line; badges and buttons on the right overflow it evenly. */
  box-sizing: content-box;
  height: 14px;
  margin: 0 -20px 12px;
  padding: 0 20px 14px;
  border-bottom: 1px solid var(--div);
}

.dash-card__title {
  margin: 0;
  white-space: nowrap;
}

.dash-card__extra {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-left: auto;
}

@media (max-width: 767px) {
  .dash-card {
    padding: 12px 16px;
  }

  .dash-card__head {
    margin: 0 -16px 12px;
    padding: 0 16px 12px;
  }

  /* A 44px button does not fit the card padding around a 14px line: the line grows. */
  .dash-card__head--tool {
    height: auto;
  }
}
</style>
