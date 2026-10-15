<script setup lang="ts">
import XBadge from './XBadge.vue'
import XStatusDot from './XStatusDot.vue'
import type { RowTone } from './XDataTable.vue'

/**
 * One log record as a card, for phones where a table does not fit: a dot and
 * the title with the time on top; a name, a detail and a badge below. Tinted
 * by tone like XDataTable rows. loading draws grey bars instead.
 */
defineProps<{
  title?: string
  loading?: boolean
  time?: string
  name?: string
  detail?: string
  badge?: string
  tone?: RowTone
}>()
</script>

<template>
  <div v-if="loading" class="x-record x-skeleton" aria-hidden="true">
    <div class="x-record__top"><span class="x-skeleton__bar" style="width: 55%" /></div>
    <div class="x-record__bottom"><span class="x-skeleton__bar" style="width: 35%" /></div>
  </div>
  <div v-else :class="['x-record', tone && `x-record--${tone}`]">
    <div class="x-record__top">
      <XStatusDot size="sm" :color="tone ? `var(--${tone === 'err' ? 'err' : 'upload'})` : 'var(--faint)'" />
      <span class="x-record__title">{{ title }}</span>
      <span v-if="time" class="x-record__time">{{ time }}</span>
    </div>
    <div class="x-record__bottom">
      <span v-if="name" class="x-record__name">{{ name }}</span>
      <span class="x-record__detail">{{ detail }}</span>
      <XBadge v-if="badge" kind="outline" mono :tone="tone">{{ badge }}</XBadge>
    </div>
  </div>
</template>

<style scoped>
.x-record {
  display: flex;
  flex-direction: column;
  gap: 8px;
  min-width: 0;
  padding: 12px;
  border: 1px solid var(--field-line);
  border-radius: var(--r-12);
  background: var(--field);
  color: var(--fg2);
}

.x-record--err {
  background: color-mix(in srgb, var(--err) 6%, var(--field));
  color: var(--err-t);
}

.x-record--upload {
  background: color-mix(in srgb, var(--upload) 6%, var(--field));
  color: var(--upload-t);
}

.x-record__top,
.x-record__bottom {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
  min-height: 16px;
}

.x-record__title {
  flex: 1;
  min-width: 0;
  overflow: hidden;
  font: 500 var(--fs-13) / var(--lh-tight) var(--font-mono);
  white-space: nowrap;
  text-overflow: ellipsis;
}

.x-record__time {
  flex: none;
  color: var(--dim);
  font: 400 var(--fs-11) / var(--lh-tight) var(--font-mono);
}

.x-record__name {
  flex: none;
  max-width: 120px;
  overflow: hidden;
  color: var(--fg2);
  font: 500 var(--fs-11) / var(--lh-tight) var(--font-sans);
  white-space: nowrap;
  text-overflow: ellipsis;
}

.x-record__detail {
  flex: 1;
  min-width: 0;
  overflow: hidden;
  color: var(--dim);
  font: 400 var(--fs-11) / var(--lh-tight) var(--font-mono);
  white-space: nowrap;
  text-overflow: ellipsis;
}
</style>
