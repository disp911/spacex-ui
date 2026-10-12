<script setup lang="ts">
/** A labelled figure; short is the label phones show. */
defineProps<{ label: string; short?: string; value: string }>()
</script>

<template>
  <div class="stat">
    <div class="stat__label">
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
  margin-bottom: 5px;
  overflow: hidden;
  color: var(--dim);
  font: 500 var(--fs-10) / normal var(--font-mono);
  letter-spacing: .07em;
  text-transform: uppercase;
  white-space: nowrap;
  text-overflow: ellipsis;
}

.stat__short {
  display: none;
}

.stat__value {
  overflow: hidden;
  color: var(--fg);
  font: 500 var(--fs-16) / normal var(--font-mono);
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
