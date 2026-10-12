<script setup lang="ts">
/** Page frame: the title on desktop (the top bar shows it on phones) and actions. */
defineProps<{ title: string }>()
</script>

<template>
  <div class="page">
    <header class="page__head">
      <h1 class="page__title">{{ title }}</h1>
      <div v-if="$slots.actions" class="page__actions"><slot name="actions" /></div>
    </header>
    <slot />
  </div>
</template>

<style scoped>
.page {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.page__head {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 12px;
  min-height: var(--ctl-h);
}

.page__title {
  flex: 1;
  min-width: 0;
  margin: 0;
  font: 600 var(--fs-20) / var(--lh-tight) var(--font-sans);
  letter-spacing: -.5px;
}

.page__actions {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

@media (max-width: 767px) {
  .page__head {
    min-height: 0;
  }

  .page__title {
    position: absolute;
    width: 1px;
    height: 1px;
    overflow: hidden;
    clip: rect(0 0 0 0);
    white-space: nowrap;
  }

  .page__head:not(:has(.page__actions)) {
    display: contents;
  }
}
</style>
