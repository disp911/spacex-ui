<script setup lang="ts">
/**
 * The start of a part of a long window, under a line: its title. What the part
 * holds is chosen in its first row, e.g. the transport's select or the
 * security's none, Reality and TLS.
 */
defineProps<{ title: string; mono?: boolean }>()
</script>

<template>
  <div class="x-part-head">
    <h3 :class="['x-part-head__title', { 'x-part-head__title--mono': mono }]">{{ title }}</h3>
  </div>
</template>

<style scoped>
.x-part-head {
  margin-top: 4px;
  padding-top: 16px;
  border-top: 1px solid var(--div);
}

.x-part-head__title {
  margin: 0;
  color: var(--fg);
  font: 600 var(--fs-14) / var(--lh-tight) var(--font-sans);
  letter-spacing: -.15px;
}

.x-part-head__title--mono {
  font: 600 var(--fs-13) / var(--lh-tight) var(--font-mono);
}

@media (max-width: 767px) {
  .x-part-head__title {
    font-size: var(--fs-16);
  }

  .x-part-head__title--mono {
    font-size: var(--fs-14);
  }
}
</style>
