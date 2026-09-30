<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import XModal from '@/components/ui/XModal.vue'
import XButton from '@/components/ui/XButton.vue'
import XSelect from '@/components/ui/XSelect.vue'
import XCheckChip from '@/components/ui/XCheckChip.vue'
import XEmpty from '@/components/ui/XEmpty.vue'
import { server } from '@/api/panel'
import { saveBlob } from '@/api/http'
import { toastError } from '@/composables/useToast'
import { t } from '@/i18n'

const open = defineModel<boolean>('open', { default: false })

type Level = 'debug' | 'info' | 'notice' | 'warning' | 'error'
const LEVEL_COLOR: Record<string, string> = {
  DEBUG: 'var(--info)',
  INFO: 'var(--accent)',
  NOTICE: 'var(--accent)',
  WARNING: 'var(--warning)',
  ERROR: 'var(--danger)',
}

const count = ref(100)
const level = ref<Level>('info')
const syslog = ref(false)
const raw = ref<string[]>([])
const loading = ref(false)

const countOptions = [10, 20, 50, 100, 500].map((n) => ({ value: n, label: String(n) }))
const levelOptions = (['debug', 'info', 'notice', 'warning', 'error'] as Level[]).map((l) => ({
  value: l,
  label: l[0].toUpperCase() + l.slice(1),
  dot: LEVEL_COLOR[l.toUpperCase()],
}))

async function load() {
  loading.value = true
  try {
    raw.value = (await server.logs(count.value, level.value, syslog.value)) ?? []
  } catch (e) {
    toastError(e)
  } finally {
    loading.value = false
  }
}

watch(open, (v) => v && load())
watch([count, level, syslog], () => open.value && load())

// "2026/09/18 04:13:30 WARNING - X-UI: message"
const rows = computed(() =>
  raw.value.map((line) => {
    const [head, ...rest] = line.split(' - ')
    const msg = rest.join(' - ')
    const parts = head.split(' ')
    if (parts.length === 3) {
      const src = /^([A-Z-]+:)\s?(.*)$/.exec(msg)
      return { ts: `${parts[0]} ${parts[1]}`, lvl: parts[2], src: src?.[1] ?? '', msg: src ? src[2] : msg }
    }
    return { ts: '', lvl: '', src: '', msg: line }
  }),
)

function download() {
  saveBlob(new Blob([raw.value.join('\n')], { type: 'text/plain' }), 'x-ui.log')
}
</script>

<template>
  <XModal v-model:open="open" :title="t('plog.title')" :width="800" flush>
    <template #subtitle>
      {{ t('plog.subtitle', { n: count }) }} · <span class="mono">journalctl -u x-ui</span>
    </template>
    <template #toolbar>
      <div class="tb">
        <XButton icon="restart" icon-only size="sm" :label="t('common.refresh')" :loading="loading" @click="load" />
        <XSelect v-model="count" :options="countOptions" :label="t('plog.lines')" mono :menu-width="112" />
        <XSelect v-model="level" :options="levelOptions" :label="t('plog.level')" :menu-width="150" />
        <XCheckChip v-model="syslog" label="SysLog" />
        <XButton icon="download" size="sm" class="tb__end" @click="download">{{ t('common.download') }}</XButton>
      </div>
    </template>

    <div class="log">
      <div class="log__scroll">
        <div v-if="loading && !rows.length" class="log__skeleton">
          <div v-for="w in [62, 48, 71, 55, 38, 66, 52, 44, 69, 58]" :key="w" class="sk"><span /><span :style="{ width: w + '%' }" /></div>
        </div>
        <XEmpty v-else-if="!rows.length" icon="log" :title="t('plog.emptyTitle')" :text="t('plog.emptyText')" />
        <div v-else class="log__lines">
          <div v-for="(r, i) in rows" :key="i" class="line">
            <span class="line__ts">{{ r.ts }}</span>
            <span v-if="r.lvl" class="line__lvl" :style="{ color: LEVEL_COLOR[r.lvl] ?? 'var(--text-3)' }">{{ ' ' + r.lvl }}</span>
            <span v-if="r.lvl" class="line__sep">{{ ' - ' }}</span>
            <span class="line__src">{{ r.src }}</span>
            <span class="line__msg">{{ (r.src ? ' ' : '') + r.msg }}</span>
          </div>
        </div>
      </div>
    </div>
  </XModal>
</template>

<style scoped>
.tb {
  display: flex;
  align-items: center;
  gap: 9px;
  width: 100%;
  flex-wrap: wrap;
}
.tb__end {
  margin-left: auto;
}
.log {
  height: min(406px, 60vh);
  display: flex;
  padding: 5px 3px 5px 0;
  background: var(--bg);
  border: 1px solid var(--border-subtle);
  border-radius: var(--r-md);
  overflow: hidden;
}
.log__scroll {
  flex: 1;
  min-width: 0;
  overflow: auto;
  padding: 5px 9px 5px 12px;
}
.line {
  font: var(--fw-regular) var(--fs-sm) / 1.75 var(--font-mono);
  white-space: pre;
}
.line__ts { color: var(--info); }
.line__lvl { font-weight: var(--fw-medium); }
.line__sep { color: var(--text-5); }
.line__src { font-weight: var(--fw-semibold); color: var(--text-2); }
.line__msg { color: var(--text-3); }
.log__skeleton {
  animation: spx-pulse 1.4s ease-in-out infinite;
}
.sk {
  display: flex;
  align-items: center;
  gap: 10px;
  height: 20px;
}
.sk span {
  height: 8px;
  border-radius: 4px;
  background: var(--control);
}
.sk span:first-child {
  width: 30%;
  max-width: 132px;
  flex: none;
}
@media (max-width: 767px) {
  .tb {
    flex-wrap: nowrap;
    overflow-x: auto;
    scrollbar-width: none;
    margin: 0 -15px;
    padding: 0 15px;
  }
  .tb::-webkit-scrollbar {
    display: none;
  }
  .log {
    height: auto;
    flex: 1;
  }
}
</style>
