<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { api, saveBlob } from '../lib/api'
import { t, tn } from '../lib/i18n'
import Dropdown from '../components/Dropdown.vue'
import Icon from '../components/Icon.vue'
import Modal from '../components/Modal.vue'
import Spinner from '../components/Spinner.vue'

// The panel's own log (journalctl -u x-ui, or syslog).
const emit = defineEmits<{ close: [] }>()

const LEVEL_COLORS: Record<string, string> = {
  DEBUG: 'var(--info)',
  INFO: 'var(--accent)',
  NOTICE: 'var(--accent)',
  WARNING: 'var(--warn-state)',
  ERROR: 'var(--danger)',
}
const rowOptions = [10, 20, 50, 100, 500].map((n) => ({ value: n, label: String(n) }))
const levelOptions = ['debug', 'info', 'notice', 'warning', 'error'].map((l) => ({
  value: l,
  label: l[0].toUpperCase() + l.slice(1),
  dot: LEVEL_COLORS[l.toUpperCase()],
}))

const rows = ref(100)
const level = ref('info')
const syslog = ref(false)
const loading = ref(false)
const raw = ref<string[] | null>(null)

interface Line {
  ts: string
  level: string
  source: string
  message: string
}

// "2026/09/18 04:13:30 INFO - XRAY: message" → parts.
function parse(line: string): Line {
  const [head, ...rest] = line.split(' - ')
  const message = rest.join(' - ')
  const parts = head.trim().split(' ')
  if (parts.length === 3 && message) {
    const isXray = message.startsWith('XRAY:')
    return {
      ts: `${parts[0]} ${parts[1]}`,
      level: parts[2],
      source: isXray ? 'XRAY:' : 'X-UI:',
      message: isXray ? message.slice(5).trimStart() : message,
    }
  }
  return { ts: '', level: '', source: '', message: line }
}

const lines = ref<Line[]>([])
async function load() {
  loading.value = true
  const msg = await api.post<string[]>(`panel/api/server/logs/${rows.value}`, { level: level.value, syslog: syslog.value }, { notify: false })
  raw.value = msg.success && Array.isArray(msg.obj) ? msg.obj.filter((l) => l.trim()) : []
  lines.value = raw.value.map(parse)
  loading.value = false
}

function download() {
  saveBlob(new Blob([(raw.value ?? []).join('\n')], { type: 'text/plain' }), 'x-ui.log')
}

watch([rows, level, syslog], load)
onMounted(load)
</script>

<template>
  <Modal :title="t('spx.panelLog')" :width="800" :height="560" body-class="log-body" @close="emit('close')">
    <template #subtitle>
      {{ tn('spx.lastLines', rows) }} · <span class="mono">{{ syslog ? 'syslog' : 'journalctl -u x-ui' }}</span>
    </template>
    <template #toolbar>
      <button class="ctl icon" :title="t('spx.refresh')" @click="load">
        <Icon name="restart" :size="15" :class="{ spinning: loading }" />
      </button>
      <Dropdown v-model="rows" :options="rowOptions" :label="t('spx.rows')" :title="t('spx.rows')" />
      <Dropdown v-model="level" :options="levelOptions" :label="t('spx.level')" :title="t('spx.level')" />
      <button class="ctl check" :class="{ on: syslog }" @click="syslog = !syslog">
        <span class="box"><Icon v-if="syslog" name="check" :size="11" /></span>SysLog
      </button>
      <button class="ctl dl" :title="t('download')" @click="download">
        <Icon name="download" :size="14" /><span class="dl-label">{{ t('download') }}</span>
      </button>
    </template>

    <div class="log">
      <div v-if="raw === null" class="state"><Spinner :size="20" /></div>
      <div v-else-if="!lines.length" class="state empty">{{ t('spx.noRecords') }}</div>
      <div v-else class="lines">
        <div v-for="(l, i) in lines" :key="i" class="line">
          <template v-if="l.ts">
            <span class="ts">{{ l.ts }}</span>
            <span class="lvl" :style="{ color: LEVEL_COLORS[l.level] ?? 'var(--text-3)' }">{{ l.level }}</span>
            <span class="dash">-</span>
            <span class="src">{{ l.source }}</span>
          </template>
          <span class="msg">{{ l.message }}</span>
        </div>
      </div>
    </div>
  </Modal>
</template>

<style scoped>
:deep(.log-body) {
  display: flex;
  flex-direction: column;
}
.mono {
  font-family: var(--font-mono);
  color: var(--text-2);
}
.spinning {
  animation: spin 0.7s linear infinite;
}
.check {
  background: var(--chip);
  gap: 8px;
}
.box {
  width: 15px;
  height: 15px;
  border-radius: 5px;
  border: 1px solid var(--scroll-thumb);
  display: flex;
  align-items: center;
  justify-content: center;
}
.check.on .box {
  border-color: transparent;
  background: var(--accent);
  color: var(--on-accent);
}
.dl {
  margin-left: auto;
  background: var(--chip);
}
.log {
  flex: 1;
  min-height: 0;
  display: flex;
  background: var(--fill);
  border: 1px solid var(--border);
  border-radius: 12px;
  overflow: hidden;
}
.state {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--accent);
}
.state.empty {
  color: var(--muted);
  font: 400 12.5px var(--font-sans);
}
.lines {
  flex: 1;
  overflow: auto;
  padding: 10px 12px;
}
.line {
  font: 400 11.5px/1.75 var(--font-mono);
  white-space: nowrap;
  color: var(--text-2);
}
.ts {
  color: var(--info);
}
.lvl {
  margin-left: 0.6em;
  font-weight: 500;
}
.dash {
  margin: 0 0.5em;
  color: var(--faint);
}
.src {
  font-weight: 600;
  color: var(--text-2);
  margin-right: 0.5em;
}
@media (max-width: 760px) {
  .dl-label {
    display: none;
  }
  .dl {
    width: 34px;
    padding: 0;
  }
}
</style>
