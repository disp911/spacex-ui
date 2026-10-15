<script setup lang="ts">
/** A caps label with a line, splitting the rows of a settings section into groups. */
defineProps<{ label: string }>()
</script>

<template>
  <div class="x-setting-divider" role="separator" :aria-label="label">
    <span class="x-caps x-caps--lg">{{ label }}</span>
    <span class="x-setting-divider__line" />
  </div>
</template>

<style scoped>
.x-setting-divider {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 20px 20px 4px;
  border-top: 1px solid var(--div);
}

.x-setting-divider .x-caps {
  white-space: nowrap;
}

.x-setting-divider__line {
  flex: 1;
  height: 1px;
  background: var(--div);
}

@media (max-width: 767px) {
  .x-setting-divider {
    padding: 18px 16px 2px;
  }
}
</style>
