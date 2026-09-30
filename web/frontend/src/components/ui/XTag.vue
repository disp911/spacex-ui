<script setup lang="ts">
// Small outlined label: versions, types, statuses.
withDefaults(defineProps<{ tone?: 'neutral' | 'accent' | 'info' | 'danger' | 'version'; mono?: boolean }>(), {
  tone: 'neutral',
  mono: true,
})
</script>

<template>
  <span class="x-tag" :class="[`x-tag--${tone}`, { mono }]"><slot /></span>
</template>

<style scoped>
.x-tag {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  flex: none;
  gap: 6px;
  padding: 2px 6px;
  border: 1px solid;
  border-radius: var(--r-xs);
  font-size: var(--fs-2xs);
  font-weight: var(--fw-medium);
  letter-spacing: 0.04em;
  white-space: nowrap;
}
.x-tag--neutral { color: var(--text-2); border-color: var(--border-control); }
.x-tag--accent { color: var(--accent); border-color: var(--accent-border); }
.x-tag--info { color: var(--info); border-color: var(--info-border); }
.x-tag--danger { color: var(--danger-text); border-color: var(--danger-border); }
/* Version pill shown in card headers (e.g. "v26.4.25"). */
.x-tag--version {
  height: var(--h-xs);
  padding: 0 11px;
  border-radius: var(--r-sm);
  background: var(--control);
  border-color: var(--border);
  color: var(--text);
  font-size: var(--fs-md);
  letter-spacing: 0;
}
</style>
