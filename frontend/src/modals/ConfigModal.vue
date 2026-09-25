<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { api, saveBlob } from '../lib/api'
import { t } from '../lib/i18n'
import Icon from '../components/Icon.vue'
import Modal from '../components/Modal.vue'
import Spinner from '../components/Spinner.vue'

// Read-only viewer for the generated Xray config.json.
const emit = defineEmits<{ close: [] }>()
const text = ref<string | null>(null)
const copied = ref(false)
const gutter = ref<HTMLElement>()

const lines = computed(() => (text.value ?? '').split('\n'))
const check = computed(() => {
  if (text.value === null) return { ok: true, line: 0 }
  try {
    JSON.parse(text.value)
    return { ok: true, line: 0 }
  } catch (e) {
    const m = /line (\d+)/.exec(String(e)) ?? /position (\d+)/.exec(String(e))
    const line = m ? (String(e).includes('line') ? +m[1] : text.value.slice(0, +m[1]).split('\n').length) : 0
    return { ok: false, line }
  }
})

onMounted(async () => {
  const msg = await api.get<unknown>('panel/api/server/getConfigJson')
  text.value = msg.success ? JSON.stringify(msg.obj, null, 2) : ''
})

async function copy() {
  if (!text.value) return
  try {
    await navigator.clipboard.writeText(text.value)
  } catch {
    const ta = document.createElement('textarea')
    ta.value = text.value
    document.body.appendChild(ta)
    ta.select()
    document.execCommand('copy')
    ta.remove()
  }
  copied.value = true
  setTimeout(() => (copied.value = false), 1600)
}

function download() {
  if (text.value) saveBlob(new Blob([text.value], { type: 'application/json' }), 'config.json')
}

function syncGutter(e: Event) {
  if (gutter.value) gutter.value.scrollTop = (e.target as HTMLElement).scrollTop
}
</script>

<template>
  <Modal :title="t('pages.index.config')" note="config.json" :width="600" :height="660" body-class="cfg-body" @close="emit('close')">
    <div class="editor" :class="{ bad: !check.ok }">
      <div v-if="text === null" class="loading"><Spinner :size="20" /></div>
      <template v-else>
        <div ref="gutter" class="gutter" aria-hidden="true">
          <div v-for="(_, i) in lines" :key="i" :class="{ errline: i + 1 === check.line }">{{ i + 1 }}</div>
        </div>
        <pre class="code" @scroll="syncGutter">{{ text }}</pre>
      </template>
    </div>
    <template #footer>
      <span class="status" :class="{ bad: !check.ok }">
        <span class="sdot" />
        {{ check.ok ? t('spx.jsonValid', { n: lines.length }) : t('spx.jsonError', { line: check.line }) }}
      </span>
      <button class="ctl dl" @click="download"><Icon name="download" :size="14" /><span class="mono">config.json</span></button>
      <button class="ctl primary" @click="copy">
        <Icon :name="copied ? 'check' : 'copy'" :size="15" />{{ copied ? t('copied') : t('copy') }}
      </button>
    </template>
  </Modal>
</template>

<style scoped>
:deep(.cfg-body) {
  display: flex;
  flex-direction: column;
}
.editor {
  flex: 1;
  min-height: 0;
  display: flex;
  background: var(--fill);
  border: 1px solid var(--border);
  border-radius: 12px;
  overflow: hidden;
}
.editor.bad {
  border-color: rgba(var(--danger-rgb), 0.5);
}
.loading {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--accent);
}
.gutter {
  flex: none;
  padding: 12px 7px 12px 6px;
  min-width: 36px;
  text-align: right;
  font: 400 11.5px/1.65 var(--font-mono);
  color: var(--faint);
  border-right: 1px solid var(--border-soft);
  overflow: hidden;
  user-select: none;
}
.errline {
  color: var(--danger);
  background: rgba(var(--danger-rgb), 0.14);
  margin: 0 -7px 0 -6px;
  padding: 0 7px 0 6px;
  font-weight: 600;
}
.code {
  flex: 1;
  min-width: 0;
  margin: 0;
  padding: 12px 12px 12px 12px;
  overflow: auto;
  font: 400 11.5px/1.65 var(--font-mono);
  color: var(--text-2);
  white-space: pre;
  tab-size: 2;
}
.status {
  margin-right: auto;
  display: flex;
  align-items: center;
  gap: 8px;
  font: 500 11px var(--font-mono);
  color: var(--text-3b);
}
.sdot {
  width: 6px;
  height: 6px;
  border-radius: 3px;
  background: var(--accent);
}
.status.bad {
  color: var(--danger-text);
}
.status.bad .sdot {
  background: var(--danger);
}
.dl {
  background: var(--chip);
}
.dl .mono {
  font: 500 12px var(--font-mono);
}
@media (max-width: 760px) {
  .status {
    display: none;
  }
}
</style>
