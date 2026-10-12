<script setup lang="ts">
/**
 * Page frame: the title with an optional subtitle on desktop (the top bar
 * shows the title on phones), actions on the right and an aside — a notice
 * that sits next to the title on desktop and above the content on phones.
 */
defineProps<{ title: string }>()
</script>

<template>
  <div class="page">
    <header class="page__head">
      <div class="page__heading">
        <h1 class="page__title">{{ title }}</h1>
        <p v-if="$slots.subtitle" class="page__subtitle"><slot name="subtitle" /></p>
      </div>
      <div v-if="$slots.actions" class="page__actions"><slot name="actions" /></div>
      <div v-if="$slots.aside" class="page__aside"><slot name="aside" /></div>
    </header>
    <slot />
  </div>
</template>

<style scoped>
.page {
  display: flex;
  flex: 1 0 auto;
  flex-direction: column;
  gap: 16px;
}

.page__head {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 12px 20px;
  min-height: var(--ctl-h);
}

.page__heading {
  flex: 1;
  min-width: 0;
}

.page__title {
  margin: 0;
  font: 600 var(--fs-20) / var(--lh-tight) var(--font-sans);
  letter-spacing: -.5px;
}

.page__subtitle {
  margin: 4px 0 0;
  color: var(--dim);
  font-size: var(--fs-12);
  white-space: nowrap;
}

.page__actions {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.page__aside {
  flex: 0 1 calc(50% - 6px);
  min-width: 0;
}

@media (max-width: 767px) {
  .page__head {
    min-height: 0;
  }

  /* The top bar already shows the title and the server. */
  .page__heading {
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

  .page__aside {
    flex-basis: 100%;
  }
}
</style>
