<script setup lang="ts">
import { computed } from 'vue'
import type { Health } from './model'
import { formatSize } from '@/utils/format'
import { t } from '@/i18n'

// "12.40 GB / 200.00 GB" with a thin quota bar; ∞ when unlimited.
const props = defineProps<{ up: number; down: number; total: number; health?: Health }>()

const used = computed(() => props.up + props.down)
const pct = computed(() => (props.total > 0 ? Math.min(100, (used.value / props.total) * 100) : 0))
const tone = computed(() => (props.health === 'depleted' ? 'is-danger' : props.health === 'warn' ? 'is-warn' : ''))
</script>

<template>
  <div class="usage" :class="tone">
    <div class="usage__text">
      <span class="usage__used">{{ formatSize(used) }}</span>
      <span class="usage__sep">/</span>
      <span class="usage__total">{{ total > 0 ? formatSize(total) : t('ib.unlimited') }}</span>
    </div>
    <div v-if="total > 0" class="usage__track"><div class="usage__bar" :style="{ width: pct + '%' }" /></div>
  </div>
</template>

<style scoped>
.usage {
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 5px;
}
.usage__text {
  display: flex;
  align-items: baseline;
  gap: 5px;
  font: var(--fw-regular) var(--fs-sm) var(--font-mono);
  white-space: nowrap;
}
.usage__used {
  color: var(--text);
  font-weight: var(--fw-medium);
}
.usage__sep,
.usage__total {
  color: var(--text-4);
}
.usage__track {
  height: 3px;
  border-radius: 2px;
  background: var(--track);
  overflow: hidden;
}
.usage__bar {
  height: 100%;
  border-radius: 2px;
  background: var(--accent);
}
.usage.is-warn .usage__bar { background: var(--warning); }
.usage.is-danger .usage__bar { background: var(--danger); }
.usage.is-danger .usage__used { color: var(--danger-text); }
</style>
