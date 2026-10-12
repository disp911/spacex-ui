<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { api, ApiError } from '../../api'
import { XButton, XCheckbox, XEmpty, XMenu, XModal, type MenuItem } from '../../kit'

/** The panel's own log, or the x-ui unit's journal with SysLog on. */
const open = defineModel<boolean>('open', { default: false })
const { t } = useI18n()

const counts = [10, 20, 50, 100, 500]
// Values are the syslog priorities; the panel's own log understands them too.
const levels = [
  { value: 'debug', label: 'Debug', tone: 'debug', dot: 'var(--upload)' },
  { value: 'info', label: 'Info', tone: 'info', dot: 'var(--acc)' },
  { value: 'notice', label: 'Notice', tone: 'info', dot: 'var(--acc)' },
  { value: 'warning', label: 'Warning', tone: 'warn', dot: 'var(--warn)' },
  { value: 'err', label: 'Error', tone: 'err', dot: 'var(--err)' },
]
const tones: Record<string, string> = { DEBUG: 'debug', INFO: 'info', NOTICE: 'info', WARNING: 'warn', ERROR: 'err', CRITICAL: 'err' }

const count = ref(100)
const level = ref('info')
const syslog = ref(false)
const lines = ref<string[]>([])
const loading = ref(false)
const failed = ref('')

const current = computed(() => levels.find((l) => l.value === level.value) ?? levels[1]!)
const countItems = computed<MenuItem[]>(() => counts.map((n) => ({ key: String(n), label: String(n), checked: n === count.value })))
const levelItems = computed<MenuItem[]>(() => levels.map((l) => ({ key: l.value, label: l.label, dot: l.dot, checked: l.value === level.value })))

// Panel lines read "2026/09/18 04:13:30 INFO - XRAY: message"; journal lines are shown as they are.
const LINE = /^(\d{4}\/\d{2}\/\d{2} \d{2}:\d{2}:\d{2}) ([A-Z]+) - (.*)$/
const rows = computed(() =>
  lines.value.map((raw) => {
    const m = LINE.exec(raw)
    if (!m) return { raw }
    const [, time, lvl, rest] = m as unknown as [string, string, string, string]
    const xray = rest.startsWith('XRAY:')
    return { time, lvl, tone: tones[lvl] ?? 'plain', src: xray ? 'XRAY:' : 'X-UI:', msg: xray ? rest.slice(5).trimStart() : rest }
  }),
)

let request = 0
async function load() {
  const id = ++request
  loading.value = true
  failed.value = ''
  const body = new FormData()
  body.set('level', level.value)
  body.set('syslog', String(syslog.value))
  try {
    const reply = await api.post<string[] | null>(`panel/api/server/logs/${count.value}`, body)
    if (id === request) lines.value = (reply ?? []).filter((l) => l.trim())
  } catch (e) {
    if (id !== request) return
    lines.value = []
    failed.value = e instanceof ApiError ? e.message : String(e)
  } finally {
    if (id === request) loading.value = false
  }
}

watch(open, (on) => on && load(), { immediate: true })
watch([count, level, syslog], () => open.value && load())

function download() {
  const url = URL.createObjectURL(new Blob([lines.value.join('\n')], { type: 'text/plain' }))
  const a = document.createElement('a')
  a.href = url
  a.download = 'x-ui.log'
  a.click()
  URL.revokeObjectURL(url)
}

// Skeleton bars: message widths vary like real lines.
const skeleton = [62, 48, 71, 55, 38, 66, 52, 44, 69, 58, 41, 63]
</script>

<template>
  <XModal v-model:open="open" :title="t('dashboard.panelLog')" size="xl" flush full>
    <template #subtitle>
      {{ t('dashboard.logLines', { n: count }) }}<template v-if="syslog"> · <span class="mono">journalctl -u x-ui</span></template>
    </template>

    <div class="log-bar">
      <XButton icon="refresh" :label="t('dashboard.logRefresh')" :loading="loading" @click="load" />

      <XMenu v-slot="{ props }" :items="countItems" :title="t('dashboard.logRows')" placement="bottom-start" sheet @select="count = Number($event)">
        <XButton v-bind="props" :label="`${t('dashboard.logRows')}: ${count}`" caret>
          <span class="x-caps log-filter__label">{{ t('dashboard.logRows') }}</span> <span class="mono">{{ count }}</span>
        </XButton>
      </XMenu>

      <XMenu v-slot="{ props }" :items="levelItems" :title="t('dashboard.logLevel')" placement="bottom-start" sheet @select="level = $event">
        <XButton v-bind="props" :label="`${t('dashboard.logLevel')}: ${current.label}`" caret>
          <span class="x-caps log-filter__label">{{ t('dashboard.logLevel') }}</span> <span :class="['log-dot', `log-tone--${current.tone}`]" /> {{ current.label }}
        </XButton>
      </XMenu>

      <XCheckbox v-model="syslog" boxed>SysLog</XCheckbox>

      <XButton icon="download" :label="t('dashboard.logDownload')" :disabled="!lines.length" phone-icon-only class="log-download" @click="download">
        {{ t('dashboard.logDownload') }}
      </XButton>
    </div>

    <div class="log-area">
      <div class="log-box">
        <div class="log-scroll" tabindex="0" :aria-busy="loading || undefined">
          <div v-if="loading" class="log-skeleton" aria-hidden="true">
            <div v-for="w in skeleton" :key="w" class="log-skeleton__row">
              <span class="log-skeleton__bar" style="width: 30%; max-width: 132px" />
              <span class="log-skeleton__bar" :style="{ width: `${w}%` }" />
            </div>
          </div>
          <XEmpty
            v-else-if="!rows.length"
            :title="failed ? t('dashboard.logFailed', { error: failed }) : t('dashboard.logEmpty')"
            :hint="failed ? undefined : t('dashboard.logEmptyHint')"
          />
          <template v-else>
            <div v-for="(row, i) in rows" :key="i" class="log-line">
              <!-- One line, so copied text keeps the spaces between the fields. -->
              <template v-if="row.time">
                <span class="log-time">{{ row.time }}</span> <span :class="['log-lvl', `log-tone--${row.tone}`]">{{ row.lvl }}</span> <span class="log-sep">-</span> <span class="log-src">{{ row.src }}</span> <span class="log-msg">{{ row.msg }}</span>
              </template>
              <span v-else class="log-msg">{{ row.raw }}</span>
            </div>
          </template>
        </div>
      </div>
    </div>
  </XModal>
</template>

<style scoped>
.log-bar {
  display: flex;
  flex: none;
  align-items: center;
  gap: 8px;
  padding: 10px 16px;
  border-bottom: 1px solid var(--div);
  background: var(--foot);
}

.log-filter__label {
  margin-right: 2px;
}

.log-dot {
  display: inline-block;
  width: 6px;
  height: 6px;
  margin: 0 2px;
  border-radius: 50%;
  background: var(--tone);
  vertical-align: middle;
}

.log-download {
  margin-left: auto;
}

.log-area {
  display: flex;
  flex: 1;
  min-height: 0;
  padding: 12px 12px 14px;
  background: var(--bg);
}

.log-box {
  display: flex;
  flex: 1;
  min-width: 0;
  height: 400px;
  padding: 4px 2px 4px 0;
  border: 1px solid var(--line);
  border-radius: var(--r-12);
  background: var(--bg);
  clip-path: inset(0 round var(--r-12));
}

.log-scroll {
  flex: 1;
  min-width: 0;
  padding: 6px 10px 6px 12px;
  overflow: auto;
  outline: none;
  font: 400 var(--fs-12) / var(--lh-code) var(--font-mono);
}

.log-line {
  white-space: nowrap;
}

.log-time { color: var(--upload); }
.log-lvl { font-weight: 500; color: var(--tone); }
.log-sep { color: var(--faint); }
.log-src { font-weight: 600; color: var(--fg2); }
.log-msg { color: var(--mut); }

.log-tone--debug { --tone: var(--upload); }
.log-tone--info { --tone: var(--acc); }
.log-tone--warn { --tone: var(--warn); }
.log-tone--err { --tone: var(--err); }
.log-tone--plain { --tone: var(--mut); }

/* On light backgrounds the solid tones are too pale for text. */
:root[data-theme='light'] .log-lvl.log-tone--warn { --tone: var(--warn-t); }
:root[data-theme='light'] .log-lvl.log-tone--err { --tone: var(--err-t); }

.log-skeleton {
  display: flex;
  flex-direction: column;
  animation: x-pulse 1.4s ease-in-out infinite;
}

.log-skeleton__row {
  display: flex;
  align-items: center;
  gap: 10px;
  height: 20px;
}

.log-skeleton__bar {
  flex: none;
  height: 8px;
  border-radius: var(--r-4);
  background: var(--skel);
}

@media (max-width: 767px) {
  /* 44px controls leave no room for the filter labels: the sheet titles name them.
     The row scrolls sideways when it does not fit, as in the mockup. */
  .log-bar {
    gap: 6px;
    padding: 10px 14px;
    overflow-x: auto;
    scrollbar-width: none;
  }

  .log-bar > * {
    flex: none;
  }

  .log-filter__label {
    display: none;
  }

  .log-box {
    height: auto;
  }
}
</style>
