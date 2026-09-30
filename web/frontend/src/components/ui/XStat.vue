<script setup lang="ts">
// Caps label over a monospace value.
withDefaults(defineProps<{ label: string; value?: string | number; tone?: 'default' | 'secondary' | 'accent' | 'info'; size?: 'md' | 'lg' }>(), {
  tone: 'default',
  size: 'md',
})
</script>

<template>
  <div class="x-stat">
    <div class="caps x-stat__label"><slot name="label">{{ label }}</slot></div>
    <div class="x-stat__value" :class="[`x-stat__value--${tone}`, `x-stat__value--${size}`]">
      <slot>{{ value }}</slot>
    </div>
  </div>
</template>

<style scoped>
.x-stat {
  min-width: 0;
}
.x-stat__label {
  margin-bottom: 5px;
  display: flex;
  align-items: center;
  gap: 6px;
}
.x-stat__value {
  font: var(--fw-medium) var(--fs-xl) / 1.25 var(--font-mono);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.x-stat__value--lg {
  font-size: var(--fs-2xl);
  font-weight: var(--fw-semibold);
  letter-spacing: -0.5px;
}
.x-stat__value--default { color: var(--text); }
.x-stat__value--secondary { color: var(--text-2); }
.x-stat__value--accent { color: var(--accent); }
.x-stat__value--info { color: var(--info); }
</style>
