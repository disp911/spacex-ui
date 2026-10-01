<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import XModal from '@/components/ui/XModal.vue'
import XButton from '@/components/ui/XButton.vue'
import XSpinner from '@/components/ui/XSpinner.vue'
import { server } from '@/api/panel'
import { saveBlob } from '@/api/http'
import { toastError } from '@/composables/useToast'
import { t } from '@/i18n'

// Read-only view of the generated Xray config.json (editable locally to
// copy/download a tweaked version; the panel builds the real file itself).
const open = defineModel<boolean>('open', { default: false })

const text = ref('')
const loading = ref(false)
const copied = ref(false)
const gutter = ref<HTMLDivElement>()

watch(open, async (v) => {
  if (!v) return
  loading.value = true
  try {
    text.value = JSON.stringify(await server.configJson(), null, 2)
  } catch (e) {
    toastError(e)
  } finally {
    loading.value = false
  }
})

const lines = computed(() => text.value.split('\n').length)
const parsed = computed(() => {
  try {
    JSON.parse(text.value)
    return { ok: true, line: 0 }
  } catch (e) {
    const msg = (e as Error).message
    const lm = /line (\d+)/.exec(msg)
    const pm = /position (\d+)/.exec(msg)
    const line = lm ? +lm[1] : pm ? text.value.slice(0, +pm[1]).split('\n').length : 0
    return { ok: false, line }
  }
})

function syncScroll(e: Event) {
  if (gutter.value) gutter.value.scrollTop = (e.target as HTMLTextAreaElement).scrollTop
}

async function copy() {
  try {
    await navigator.clipboard.writeText(text.value)
  } catch {
    /* clipboard unavailable over plain HTTP; fall through */
  }
  copied.value = true
  setTimeout(() => (copied.value = false), 1600)
}

function download() {
  saveBlob(new Blob([text.value], { type: 'application/json' }), 'config.json')
}
</script>

<template>
  <XModal v-model:open="open" :title="t('config.title')" :width="600">
    <template #tag>config.json</template>
    <div class="editor" :class="{ 'is-bad': !parsed.ok }">
      <div v-if="loading" class="editor__loading"><XSpinner :size="20" /></div>
      <template v-else>
        <div ref="gutter" class="editor__gutter" aria-hidden="true">
          <div v-for="n in lines" :key="n" :class="{ 'is-err': n === parsed.line }">{{ n }}</div>
        </div>
        <textarea v-model="text" class="editor__text" spellcheck="false" wrap="off" aria-label="config.json" @scroll="syncScroll" />
      </template>
    </div>
    <template #footer>
      <div class="status" :class="{ 'is-bad': !parsed.ok }">
        <span class="status__dot" />
        <span class="ellipsis">
          {{ parsed.ok ? t('config.valid', { n: lines }) : parsed.line ? t('config.invalidLine', { n: parsed.line }) : t('config.invalid') }}
        </span>
      </div>
      <XButton icon="download" @click="download"><span class="mono">config.json</span></XButton>
      <XButton variant="primary" :icon="copied ? 'check' : 'copy'" @click="copy">{{ copied ? t('config.copied') : t('config.copy') }}</XButton>
    </template>
  </XModal>
</template>

<style scoped>
.editor {
  position: relative;
  height: min(420px, 60vh);
  display: flex;
  padding: 5px 3px 5px 0;
  border-radius: var(--r-md);
  border: 1px solid var(--border-field);
  background: var(--field);
  overflow: hidden;
  transition: border-color 0.2s ease;
}
.editor.is-bad {
  border-color: var(--danger-border);
}
.editor__loading {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-4);
}
.editor__gutter {
  flex: none;
  min-width: 34px;
  height: 100%;
  overflow: hidden;
  padding: 7px 7px 24px 6px;
  border-right: 1px solid var(--divider);
  font: var(--fw-regular) var(--fs-sm) / 1.62 var(--font-mono);
  color: var(--text-5);
  text-align: right;
  user-select: none;
}
.editor__gutter .is-err {
  margin: 0 -7px 0 -6px;
  padding: 0 7px 0 6px;
  color: var(--danger);
  background: var(--danger-soft);
  font-weight: var(--fw-semibold);
}
.editor__text {
  flex: 1;
  min-width: 0;
  height: 100%;
  margin: 0;
  padding: 7px 11px 7px 12px;
  border: 0;
  outline: 0;
  resize: none;
  background: transparent;
  color: var(--text-2);
  font: var(--fw-regular) var(--fs-sm) / 1.62 var(--font-mono);
  tab-size: 2;
  white-space: pre;
  overflow: auto;
}
.status {
  flex: 1;
  min-width: 0;
  display: flex;
  align-items: center;
  gap: 7px;
  font: var(--fw-regular) var(--fs-xs) var(--font-mono);
  color: var(--text-3);
}
.status__dot {
  width: 6px;
  height: 6px;
  flex: none;
  border-radius: 50%;
  background: var(--accent);
}
.status.is-bad {
  color: var(--danger);
}
.status.is-bad .status__dot {
  background: var(--danger);
}
@media (max-width: 767px) {
  .editor {
    height: auto;
    flex: 1;
    min-height: 300px;
  }
  .status {
    display: none;
  }
}
</style>
