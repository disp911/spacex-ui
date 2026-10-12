<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { api, ApiError } from '../../api'
import { XButton, XCheckbox, XIcon, XMenu, XModal, type MenuItem } from '../../kit'

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
        <XButton v-bind="props" :label="`${t('dashboard.logRows')}: ${count}`" class="log-filter log-filter--count">
          <span class="log-filter__label log-filter__label--wide">{{ t('dashboard.logRows') }}</span>
          <span class="mono log-filter__value">{{ count }}</span>
          <XIcon name="caret" :size="10" :stroke="1.6" class="log-filter__caret" />
        </XButton>
      </XMenu>

      <XMenu v-slot="{ props }" :items="levelItems" :title="t('dashboard.logLevel')" placement="bottom-start" sheet @select="level = $event">
        <XButton v-bind="props" :label="`${t('dashboard.logLevel')}: ${current.label}`" class="log-filter">
          <span class="log-filter__label log-filter__label--wide">{{ t('dashboard.logLevel') }}</span>
          <span :class="['log-dot', `log-tone--${current.tone}`]" />
          <span class="log-filter__value">{{ current.label }}</span>
          <XIcon name="caret" :size="10" :stroke="1.6" class="log-filter__caret" />
        </XButton>
      </XMenu>

      <XCheckbox v-model="syslog" class="log-syslog">SysLog</XCheckbox>

      <XButton icon="download" :label="t('dashboard.logDownload')" :disabled="!lines.length" class="log-download" @click="download">
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
          <div v-else-if="!rows.length" class="log-empty">
            <span class="log-empty__icon"><XIcon name="doc" :size="16" :stroke="1.3" /></span>
            <span class="log-empty__title">{{ failed ? t('dashboard.logFailed', { error: failed }) : t('dashboard.logEmpty') }}</span>
            <span v-if="!failed" class="log-empty__hint">{{ t('dashboard.logEmptyHint') }}</span>
          </div>
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

.x-btn.log-filter {
  padding: 0 10px 0 12px;
  font-size: var(--fs-12);
  font-weight: 500;
}

.log-filter :deep(.x-btn__text) {
  display: inline-flex;
  align-items: center;
  gap: 8px;
}

.log-filter__label {
  flex: none;
  color: var(--dim);
  font: 500 var(--fs-10) / 1 var(--font-mono);
  letter-spacing: .07em;
  text-transform: uppercase;
}

.log-filter__value {
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
}

.log-filter__caret {
  flex: none;
  color: var(--dim);
}

.log-dot {
  flex: none;
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--tone);
}

.log-syslog {
  flex: none;
  height: var(--ctl-h);
  padding: 0 12px 0 10px;
  border: 1px solid var(--field-line);
  border-radius: var(--ctl-r);
  background: var(--ctrl);
  font-size: var(--fs-12);
  font-weight: 500;
  transition: border-color var(--t-fast);
}

.log-syslog:hover {
  border-color: var(--field-hover);
}

.log-syslog:has(:focus-visible) {
  border-color: var(--acc-line);
  box-shadow: 0 0 0 3px var(--acc-ring);
}

.x-btn.log-download {
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

.log-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  height: 100%;
  min-height: 200px;
  padding: 0 16px;
  text-align: center;
  font-family: var(--font-sans);
}

.log-empty__icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  border: 1px dashed var(--track);
  border-radius: var(--r-12);
  color: var(--faint);
}

.log-empty__title {
  color: var(--mut);
  font: 500 var(--fs-13) / var(--lh-tight) var(--font-sans);
}

.log-empty__hint {
  color: var(--dim);
  font: 400 var(--fs-12) / var(--lh-text) var(--font-sans);
  text-wrap: pretty;
}

@media (max-width: 767px) {
  /* 44px controls leave no room for the filter labels: the sheet titles name them.
     Lines and level share what is left; narrow phones scroll the row. */
  .log-bar {
    gap: 6px;
    padding: 10px 14px;
    overflow-x: auto;
    scrollbar-width: none;
  }

  .log-bar > * {
    flex: none;
  }

  .log-bar > .x-btn.log-filter {
    flex: 1 1 auto;
    min-width: 0;
    justify-content: space-between;
    padding: 0 8px 0 10px;
  }

  .log-bar > .x-btn.log-filter--count {
    flex-shrink: 0;
  }

  .log-filter :deep(.x-btn__text) {
    flex: 1;
    gap: 6px;
  }

  .log-filter__caret {
    margin-left: auto;
  }

  .log-syslog {
    padding: 0 10px 0 8px;
  }

  .log-filter__label--wide {
    display: none;
  }

  .x-btn.log-download {
    width: var(--ctl-h);
    padding: 0;
  }

  .log-download :deep(.x-btn__text) {
    display: none;
  }

  .log-box {
    height: auto;
  }
}
</style>
