<script setup lang="ts">
/** A labelled figure; short is the label phones show. Small figures are secondary, coloured by direction. */
defineProps<{ label: string; short?: string; value: string; small?: boolean; tone?: 'up' | 'down' }>()
</script>

<template>
  <div :class="['stat', { 'stat--small': small }, tone && `stat--${tone}`]">
    <div class="x-caps stat__label">
      <span :class="{ 'stat__long': short }">{{ label }}</span>
      <span v-if="short" class="stat__short">{{ short }}</span>
    </div>
    <div class="stat__value"><slot>{{ value }}</slot></div>
  </div>
</template>

<style scoped>
.stat {
  min-width: 0;
}

.stat__label {
  margin-bottom: 4px;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
}

.stat--small .stat__value {
  color: var(--fg2);
  font-size: var(--fs-14);
}

.stat--up .stat__value { color: var(--upload); }
.stat--down .stat__value { color: var(--acc); }

.stat__short {
  display: none;
}

.stat__value {
  overflow: hidden;
  color: var(--fg);
  font: 500 var(--fs-16) / var(--lh-tight) var(--font-mono);
  white-space: nowrap;
  text-overflow: ellipsis;
}

@media (max-width: 767px) {
  .stat__long {
    display: none;
  }

  .stat__short {
    display: inline;
  }
}
</style>
