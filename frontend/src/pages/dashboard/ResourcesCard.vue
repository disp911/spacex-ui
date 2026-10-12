<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { formatBytes, formatGigahertz, formatNumber } from '../../i18n/format'
import type { ServerStatus } from '../../stores/server'
import CpuChart from './CpuChart.vue'
import DashCard from './DashCard.vue'

const props = defineProps<{ status: ServerStatus | null }>()
const { t } = useI18n()

const share = (used: number, total: number) => (total > 0 ? (used / total) * 100 : 0)
const usedOf = (m: { current: number; total: number }) => t('dashboard.usedOf', { used: formatBytes(m.current), total: formatBytes(m.total) })

const rows = computed(() => {
  const s = props.status
  if (!s) return []
  return [
    { key: 'cpu', label: t('dashboard.cpu'), detail: t('dashboard.cores', { n: s.cpuCores, speed: formatGigahertz(s.cpuSpeedMhz) }, s.cpuCores), value: s.cpu },
    { key: 'ram', label: t('dashboard.ram'), detail: usedOf(s.mem), value: share(s.mem.current, s.mem.total) },
    { key: 'swap', label: t('dashboard.swap'), detail: usedOf(s.swap), value: share(s.swap.current, s.swap.total) },
    { key: 'disk', label: t('dashboard.disk'), detail: usedOf(s.disk), value: share(s.disk.current, s.disk.total) },
  ]
})

/** Green while there is room, amber from 80 %, red from 90 %. */
const tone = (v: number) => (v >= 90 ? 'err' : v >= 80 ? 'warn' : 'ok')
</script>

<template>
  <DashCard :title="t('dashboard.resources')">
    <div class="rows">
      <div v-for="row in rows" :key="row.key" :class="['res', `res--${tone(row.value)}`]">
        <div class="res__head">
          <span class="res__label">{{ row.label }}</span>
          <span class="res__detail">{{ row.detail }}</span>
          <span class="res__value">{{ formatNumber(row.value, row.value < 10 ? 2 : 1) }}<small>%</small></span>
        </div>
        <div class="res__bar" role="meter" :aria-label="row.label" :aria-valuenow="Math.round(row.value)" aria-valuemin="0" aria-valuemax="100">
          <span :style="{ width: `${Math.min(100, row.value)}%` }" />
        </div>
      </div>
    </div>
    <CpuChart class="chart" />
  </DashCard>
</template>

<style scoped>
.rows {
  display: flex;
  flex: 1;
  flex-direction: column;
  justify-content: center;
  gap: 16px;
}

.res {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.res--ok { --tone: var(--ok); }
.res--warn { --tone: var(--warn); }
.res--err { --tone: var(--err); }

.res__head {
  display: flex;
  align-items: baseline;
  gap: 8px;
  min-width: 0;
}

.res__label {
  color: var(--fg2);
  font: 500 var(--fs-12) / var(--lh-tight) var(--font-mono);
  letter-spacing: .06em;
  text-transform: uppercase;
  white-space: nowrap;
}

.res__detail {
  overflow: hidden;
  color: var(--dim);
  font: 400 var(--fs-11) / var(--lh-tight) var(--font-mono);
  white-space: nowrap;
  text-overflow: ellipsis;
}

.res__value {
  margin-left: auto;
  color: var(--tone);
  font: 600 var(--fs-16) / var(--lh-tight) var(--font-mono);
  white-space: nowrap;
}

.res__value small {
  margin-left: 1px;
  color: var(--dim);
  font-size: var(--fs-11);
  font-weight: 500;
}

.res__bar {
  height: 6px;
  overflow: hidden;
  border-radius: var(--r-pill);
  background: var(--bar-track);
}

.res__bar span {
  display: block;
  height: 100%;
  border-radius: var(--r-pill);
  background: var(--tone);
  transition: width var(--t-med) var(--ease);
}

.chart {
  margin-top: 14px;
  padding-top: 12px;
  border-top: 1px solid var(--div);
}
</style>
