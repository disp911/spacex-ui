<script setup lang="ts">
import { computed } from 'vue'
import XCard from '@/components/ui/XCard.vue'
import XProgress from '@/components/ui/XProgress.vue'
import XSelect from '@/components/ui/XSelect.vue'
import CpuChart from './CpuChart.vue'
import type { CpuPoint, ServerStatus } from '@/api/panel'
import { CPU_RANGES, type CpuRange } from '@/composables/useServerStatus'
import { formatSize, percent } from '@/utils/format'
import { t } from '@/i18n'

const props = defineProps<{ status: ServerStatus | null; series: CpuPoint[] }>()
const range = defineModel<CpuRange>('range', { required: true })

const s = computed(() => props.status)
const cpuMeta = computed(() => {
  if (!s.value) return ''
  const parts = [t('dash.cores', { n: s.value.cpuCores || s.value.logicalPro })]
  if (s.value.cpuSpeedMhz) parts.push((s.value.cpuSpeedMhz / 1000).toFixed(2) + ' GHz')
  return parts.join(' · ')
})
const pair = (c = 0, total = 0) => `${formatSize(c)} / ${formatSize(total)}`

const rangeOptions = computed(() =>
  (Object.keys(CPU_RANGES) as CpuRange[]).map((r) => ({ value: r, label: t('dash.range.' + r) })),
)
const peak = computed(() => (props.series.length ? Math.max(...props.series.map((p) => p.cpu)).toFixed(1) : '0.0'))
</script>

<template>
  <XCard :title="t('dash.resources')" fill class="res">
    <div class="res__bars">
      <XProgress :label="t('dash.cpu')" :meta="cpuMeta" :value="s?.cpu ?? 0" />
      <XProgress :label="t('dash.ram')" :meta="pair(s?.mem.current, s?.mem.total)" :value="percent(s?.mem.current ?? 0, s?.mem.total ?? 0)" />
      <XProgress :label="t('dash.swap')" :meta="pair(s?.swap.current, s?.swap.total)" :value="percent(s?.swap.current ?? 0, s?.swap.total ?? 0)" />
      <XProgress :label="t('dash.disk')" :meta="pair(s?.disk.current, s?.disk.total)" :value="percent(s?.disk.current ?? 0, s?.disk.total ?? 0)" />
    </div>

    <div class="res__chart">
      <div class="res__chart-head">
        <span class="caps">{{ t('dash.cpuLoad') }}</span>
        <XSelect v-model="range" :options="rangeOptions" size="xs" mono inline :menu-width="132">
          <template #value="{ option }">{{ t('dash.rangePrefix', { r: option?.label ?? '' }) }}</template>
        </XSelect>
        <span class="res__peak">{{ t('dash.peak', { v: peak }) }}</span>
      </div>
      <CpuChart :points="series" />
    </div>
  </XCard>
</template>

<style scoped>
.res__bars {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 16px;
}
.res__chart {
  flex: none;
  margin-top: 14px;
  padding-top: 13px;
  border-top: 1px solid var(--divider);
}
.res__chart-head {
  display: flex;
  align-items: center;
  gap: 9px;
  margin-bottom: 8px;
}
.res__peak {
  margin-left: auto;
  font: var(--fw-medium) var(--fs-2xs) var(--font-mono);
  color: var(--text-3);
  white-space: nowrap;
}
@media (max-width: 767px) {
  .res__bars {
    gap: 14px;
  }
}
</style>
