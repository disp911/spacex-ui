<script setup lang="ts">
import { computed, nextTick, ref, watch } from 'vue'
import XModal from '@/components/ui/XModal.vue'
import XButton from '@/components/ui/XButton.vue'
import XSpinner from '@/components/ui/XSpinner.vue'
import XEmpty from '@/components/ui/XEmpty.vue'
import type { PanelDefaults } from '@/api/panel'
import { drawQr, shareLinks } from '@/api/legacyLinks'
import { toast } from '@/composables/useToast'
import { useTheme } from '@/composables/useTheme'
import type { ClientRow, InboundRow } from './model'
import { t } from '@/i18n'

// QR code + copyable links for a client (or a WireGuard inbound):
// the protocol links, then subscription URLs when enabled.
const open = defineModel<boolean>('open', { default: false })
const props = defineProps<{ inbound: InboundRow | null; client: ClientRow | null; defaults: PanelDefaults | null }>()

interface Entry {
  label: string
  link: string
}
const entries = ref<Entry[]>([])
const current = ref(0)
const loading = ref(false)
const error = ref('')
const canvas = ref<HTMLCanvasElement>()
const { theme } = useTheme()

watch(open, async (v) => {
  if (!v || !props.inbound) return
  entries.value = []
  current.value = 0
  error.value = ''
  loading.value = true
  try {
    const links = await shareLinks(props.inbound.db, props.client?.raw ?? null, props.defaults?.remarkModel ?? '-ieo')
    const list: Entry[] = links.map((l) => ({ label: l.remark || props.inbound!.protocolLabel, link: l.link }))
    const d = props.defaults
    const subId = props.client?.subId
    if (d && subId) {
      if (d.subEnable && d.subURI) list.push({ label: t('ib.subscription'), link: d.subURI + subId })
      if (d.subJsonEnable && d.subJsonURI) list.push({ label: t('ib.subscriptionJson'), link: d.subJsonURI + subId })
    }
    entries.value = list
  } catch (e) {
    error.value = (e as Error).message
  } finally {
    loading.value = false
  }
})

const entry = computed(() => entries.value[current.value])

watch([entry, theme, canvas], async () => {
  await nextTick()
  if (entry.value && canvas.value) drawQr(canvas.value, entry.value.link, 480, theme.value === 'dark')
})

async function copy() {
  if (!entry.value) return
  try {
    await navigator.clipboard.writeText(entry.value.link)
  } catch {
    // Plain-HTTP panels have no async clipboard; fall back to a selection copy.
    const ta = document.createElement('textarea')
    ta.value = entry.value.link
    document.body.appendChild(ta)
    ta.select()
    document.execCommand('copy')
    ta.remove()
  }
  toast(t('ib.copied'), 'success', 1800)
}

const title = computed(() => props.client?.email || props.inbound?.remark || t('ib.qrTitle'))
</script>

<template>
  <XModal v-model:open="open" :title="title" :subtitle="t('ib.qrTitle')" :width="440">
    <div v-if="loading" class="qr__loading"><XSpinner :size="22" /></div>
    <XEmpty v-else-if="error" icon="alert" :title="t('ib.linksError')" :text="error" />
    <XEmpty v-else-if="!entries.length" icon="qr" :title="t('ib.noLinks')" />
    <template v-else>
      <div v-if="entries.length > 1" class="qr__tabs no-scrollbar" role="tablist">
        <button
          v-for="(e, i) in entries"
          :key="i"
          type="button"
          role="tab"
          class="qr__tab"
          :class="{ 'is-on': i === current }"
          :aria-selected="i === current"
          @click="current = i"
        >
          {{ e.label }}
        </button>
      </div>
      <div class="qr__code"><canvas ref="canvas" /></div>
      <div class="qr__link">
        <span class="qr__text">{{ entry?.link }}</span>
      </div>
    </template>
    <template v-if="entries.length" #footer>
      <XButton variant="primary" icon="copy" block @click="copy">{{ t('ib.copy') }}</XButton>
    </template>
  </XModal>
</template>

<style scoped>
.qr__loading {
  display: flex;
  justify-content: center;
  padding: 60px 0;
  color: var(--text-4);
}
.qr__tabs {
  display: flex;
  gap: 6px;
  overflow-x: auto;
}
.qr__tab {
  flex: none;
  height: var(--h-sm);
  padding: 0 12px;
  border-radius: var(--r-sm);
  border: 1px solid var(--border-control);
  background: var(--control);
  color: var(--text-3);
  font: var(--fw-medium) var(--fs-sm) var(--font-sans);
  white-space: nowrap;
}
.qr__tab.is-on {
  border-color: var(--accent-border);
  background: var(--accent-soft);
  color: var(--accent);
}
.qr__code {
  align-self: center;
  width: min(240px, 100%);
  aspect-ratio: 1;
  padding: 0;
  border-radius: var(--r-md);
  border: 1px solid var(--border-field);
  overflow: hidden;
  background: var(--field);
}
.qr__code canvas {
  display: block;
  width: 100%;
  height: 100%;
}
.qr__link {
  padding: 10px 12px;
  max-height: 96px;
  overflow: auto;
  background: var(--field);
  border: 1px solid var(--border-field);
  border-radius: var(--r-md);
}
.qr__text {
  font: var(--fw-regular) var(--fs-xs) / 1.5 var(--font-mono);
  color: var(--text-2);
  word-break: break-all;
}
@media (max-width: 767px) {
  .qr__code {
    width: min(300px, 100%);
  }
}
</style>
