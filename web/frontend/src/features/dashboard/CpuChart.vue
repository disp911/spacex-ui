<script setup lang="ts">
import { computed } from 'vue'
import type { CpuPoint } from '@/api/panel'

// Area chart of CPU %, y-axis auto-scaled to the next 25% above the peak.
const props = defineProps<{ points: CpuPoint[] }>()

const W = 560
const PLOT_W = 512 // leave room on the right for tick labels
const H = 48

const model = computed(() => {
  const vals = props.points.map((p) => Math.max(0, Math.min(100, p.cpu)))
  const peak = vals.length ? Math.max(...vals) : 0
  const top = Math.min(100, Math.max(25, Math.ceil((peak * 1.15) / 25) * 25))
  const y = (v: number) => H - 2 - (v / top) * 42

  const ticks = []
  for (let v = 0; v <= top; v += 25) ticks.push({ v, y: y(v) })

  // A single sample is drawn as a flat line across the plot.
  const pts = vals.length === 1 ? [vals[0], vals[0]] : vals
  const x = (i: number) => (PLOT_W / (pts.length - 1)) * i
  const line = pts.length ? pts.map((v, i) => `${i ? 'L' : 'M'}${x(i).toFixed(1)} ${y(v).toFixed(1)}`).join(' ') : ''
  const area = line ? `${line} L${PLOT_W} ${H} L0 ${H} Z` : ''
  return { ticks, line, area, peak }
})

defineExpose({ peak: computed(() => model.value.peak) })
</script>

<template>
  <div class="cpu-chart">
    <svg :viewBox="`0 0 ${W} ${H}`" preserveAspectRatio="none" role="img" aria-label="CPU">
      <line
        v-for="tk in model.ticks"
        :key="tk.v"
        x1="0"
        :x2="PLOT_W"
        :y1="tk.y"
        :y2="tk.y"
        class="cpu-chart__grid"
        vector-effect="non-scaling-stroke"
      />
      <path v-if="model.area" :d="model.area" class="cpu-chart__area" />
      <path v-if="model.line" :d="model.line" class="cpu-chart__line" vector-effect="non-scaling-stroke" />
    </svg>
    <span
      v-for="tk in model.ticks"
      :key="'l' + tk.v"
      class="cpu-chart__tick"
      :style="{ top: (tk.y / H) * 100 + '%' }"
    >{{ tk.v }}%</span>
  </div>
</template>

<style scoped>
.cpu-chart {
  position: relative;
  background: var(--field);
  border: 1px solid var(--border-field);
  border-radius: var(--r-md);
  overflow: hidden;
}
.cpu-chart svg {
  display: block;
  width: 100%;
  height: 88px;
}
.cpu-chart__grid {
  stroke: var(--chart-grid);
  stroke-width: 1;
}
.cpu-chart__area {
  fill: var(--accent);
  opacity: 0.14;
}
.cpu-chart__line {
  fill: none;
  stroke: var(--accent);
  stroke-width: 1.6;
  stroke-linejoin: round;
}
.cpu-chart__tick {
  position: absolute;
  right: 8px;
  transform: translateY(-50%);
  font: var(--fw-medium) 9px var(--font-mono);
  color: var(--text-5);
  pointer-events: none;
}
@media (max-width: 767px) {
  .cpu-chart svg {
    height: 96px;
  }
}
</style>
