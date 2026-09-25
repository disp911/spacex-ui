<script setup lang="ts">
defineProps<{ label: string; value?: string | number; color?: string; size?: 'md' | 'sm' }>()
</script>

<template>
  <div class="metric">
    <div class="caps-sm">{{ label }}</div>
    <div class="val" :class="size ?? 'md'" :style="color ? { color } : undefined"><slot>{{ value }}</slot></div>
  </div>
</template>

<style scoped>
.metric {
  min-width: 0;
}
.caps-sm {
  margin-bottom: 5px;
}
.val {
  font: 500 16px var(--font-mono);
  color: var(--text);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.val.sm {
  font-size: 14px;
}
</style>
