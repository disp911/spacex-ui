<script setup lang="ts">
import { computed } from 'vue'

// Area chart of CPU load. The y axis tops out at the next 25% step above
// 115% of the peak (at least 25%), with a grid line every 25%.
const props = defineProps<{ values: number[] }>()

const W = 512
const H = 48
const chart = computed(() => {
  const v = props.values.length ? props.values : [0]
  const n = v.length
  const max = Math.max(...v)
  const top = Math.min(100, Math.max(25, Math.ceil((max * 1.15) / 25) * 25))
  const x = (i: number) => (n > 1 ? (W / (n - 1)) * i : W).toFixed(1)
  const y = (p: number) => (46 - (Math.min(p, top) / top) * 42).toFixed(1)
  const line = v.map((p, i) => `${i ? 'L' : 'M'}${x(i)} ${y(p)}`).join(' ')
  const ticks = []
  for (let p = 0; p <= top; p += 25) {
    const yy = 46 - (p / top) * 42
    ticks.push({ p, y: yy.toFixed(1), pct: ((yy / H) * 100).toFixed(2) })
  }
  return { line, area: `${line} L${W} ${H} L0 ${H} Z`, ticks }
})
</script>

<template>
  <div class="chart">
    <svg viewBox="0 0 560 48" preserveAspectRatio="none">
      <line
        v-for="tk in chart.ticks"
        :key="tk.p"
        x1="0"
        :x2="W"
        :y1="tk.y"
        :y2="tk.y"
        stroke="var(--grid-line)"
        stroke-width="1"
        vector-effect="non-scaling-stroke"
      />
      <path :d="chart.area" fill="rgba(var(--accent-rgb), .14)" />
      <path :d="chart.line" fill="none" stroke="var(--accent)" stroke-width="1.6" stroke-linejoin="round" vector-effect="non-scaling-stroke" />
    </svg>
    <span v-for="tk in chart.ticks" :key="tk.p" class="tick" :style="{ top: `${tk.pct}%` }">{{ tk.p }}%</span>
  </div>
</template>

<style scoped>
.chart {
  position: relative;
  background: var(--fill);
  border: 1px solid var(--border);
  border-radius: 12px;
  overflow: hidden;
}
svg {
  width: 100%;
  height: 88px;
}
.tick {
  position: absolute;
  right: 8px;
  transform: translateY(-50%);
  font: 500 8.5px var(--font-mono);
  color: var(--faint);
  pointer-events: none;
}
</style>
