<script setup lang="ts">
import { computed } from 'vue'
import { expiryParts } from './model'
import { t } from '@/i18n'

// Expiry date plus a days-left hint colored by urgency.
const props = defineProps<{ expiryTime: number; warnDays?: number }>()

const view = computed(() => {
  const { date, days } = expiryParts(props.expiryTime)
  if (!props.expiryTime) return { main: t('ib.never'), hint: '', tone: 'muted' }
  if (props.expiryTime < 0) return { main: t('ib.daysLeft', { n: days ?? 0 }), hint: t('ib.afterFirstUse'), tone: '' }
  if (days !== null && days <= 0) return { main: date, hint: t('ib.expired'), tone: 'danger' }
  const warn = days !== null && days <= (props.warnDays || 3)
  return { main: date, hint: t('ib.daysLeft', { n: days ?? 0 }), tone: warn ? 'warn' : '' }
})
</script>

<template>
  <div class="exp" :class="view.tone && `is-${view.tone}`">
    <span class="exp__date">{{ view.main }}</span>
    <span v-if="view.hint" class="exp__hint">{{ view.hint }}</span>
  </div>
</template>

<style scoped>
.exp {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
  white-space: nowrap;
}
.exp__date {
  font: var(--fw-regular) var(--fs-sm) var(--font-mono);
  color: var(--text-2);
}
.exp__hint {
  font: var(--fw-medium) var(--fs-2xs) var(--font-mono);
  color: var(--text-4);
}
.exp.is-muted .exp__date { color: var(--text-4); font-family: var(--font-sans); }
.exp.is-warn .exp__hint { color: var(--warning); }
.exp.is-danger .exp__date,
.exp.is-danger .exp__hint { color: var(--danger-text); }
</style>
