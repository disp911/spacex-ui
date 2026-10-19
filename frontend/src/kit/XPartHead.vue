<script setup lang="ts">
/**
 * The start of a part of a long window, under a line: its title and, at the
 * right, the control that picks what the part holds, e.g. the protocol's own
 * settings or the transport with its select. wide makes room for three
 * choices side by side, e.g. the security's none, Reality and TLS.
 */
defineProps<{ title: string; mono?: boolean; wide?: boolean }>()
</script>

<template>
  <div class="x-part-head">
    <h3 :class="['x-part-head__title', { 'x-part-head__title--mono': mono }]">{{ title }}</h3>
    <div v-if="$slots.default" :class="['x-part-head__control', { 'x-part-head__control--wide': wide }]"><slot /></div>
  </div>
</template>

<style scoped>
.x-part-head {
  display: flex;
  align-items: center;
  gap: 10px;
  min-height: calc(var(--ctl-h) + 16px);
  margin-top: 4px;
  padding-top: 16px;
  border-top: 1px solid var(--div);
}

.x-part-head__title {
  flex: 1;
  min-width: 0;
  margin: 0;
  color: var(--fg);
  font: 600 var(--fs-14) / var(--lh-tight) var(--font-sans);
  letter-spacing: -.15px;
  white-space: nowrap;
}

.x-part-head__title--mono {
  font: 600 var(--fs-13) / var(--lh-tight) var(--font-mono);
}

.x-part-head__control {
  width: 230px;
  min-width: 0;
}

.x-part-head__control--wide {
  width: 250px;
}

@media (max-width: 767px) {
  .x-part-head__title {
    font-size: var(--fs-16);
  }

  .x-part-head__title--mono {
    font-size: var(--fs-14);
  }

  .x-part-head__control {
    width: 180px;
  }

  .x-part-head__control--wide {
    width: 210px;
  }
}
</style>
