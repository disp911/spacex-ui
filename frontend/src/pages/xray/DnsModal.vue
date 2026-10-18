<script setup lang="ts">
import { computed, reactive, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { XButton, XField, XFieldGrid, XInput, XInputList, XModal, XNumber, XSelect, XSwitch } from '../../kit'
import { defaultDnsServer, dnsStrategies, listOf, type DnsServer } from './model'

/**
 * A DNS server to create, or the one at index to edit: where it is, how it
 * resolves and for which domains, and the IPs its answers must or must not
 * hold. A server given as a bare address is saved with the old panel's
 * defaults around it.
 */
const props = defineProps<{ server: DnsServer | null; index: number }>()
const open = defineModel<boolean>('open', { default: false })
const emit = defineEmits<{ save: [server: DnsServer] }>()
const { t } = useI18n()

const form = reactive({
  address: '',
  port: 53 as number | null,
  queryStrategy: 'UseIP',
  domains: [] as string[],
  expectIPs: [] as string[],
  unexpectedIPs: [] as string[],
  skipFallback: true,
  disableCache: false,
  finalQuery: false,
})

/** The server as an object: a bare address gets the defaults around it. */
const base = () => {
  const s = props.server
  if (s == null) return defaultDnsServer()
  return typeof s === 'string' ? { ...defaultDnsServer(), address: s } : s
}

type Form = typeof form
/** The form as it opened, to tell what was changed. */
let shown: Form | null = null

watch(open, (on) => {
  if (!on) return
  const s = base()
  Object.assign(form, {
    address: String(s.address ?? ''),
    port: typeof s.port === 'number' ? s.port : null,
    queryStrategy: String(s.queryStrategy ?? 'UseIP'),
    domains: listOf(s.domains),
    expectIPs: listOf(s.expectIPs),
    unexpectedIPs: listOf(s.unexpectedIPs),
    skipFallback: !!s.skipFallback,
    disableCache: !!s.disableCache,
    finalQuery: !!s.finalQuery,
  })
  shown = JSON.parse(JSON.stringify(form))
})

const list = (v: string[]) => v.map((x) => x.trim()).filter(Boolean)

function save() {
  const now: Form = { ...form, address: form.address.trim(), domains: list(form.domains), expectIPs: list(form.expectIPs), unexpectedIPs: list(form.unexpectedIPs) }
  const changed = (Object.keys(now) as (keyof Form)[]).filter((k) => JSON.stringify(now[k]) !== JSON.stringify(shown?.[k]))
  open.value = false
  if (props.server != null && !changed.length) return
  const s = props.server
  if (typeof s === 'string' && changed.every((k) => k === 'address')) return emit('save', now.address)
  // A new server or a bare address with more settings: the whole form, as the old panel wrote it.
  const out: Record<string, unknown> = s && typeof s === 'object' ? { ...s } : { ...base() }
  for (const k of s && typeof s === 'object' ? changed : (Object.keys(now) as (keyof Form)[])) out[k] = now[k]
  if (out.port == null) delete out.port
  emit('save', out)
}

const strategies = dnsStrategies.map((v) => ({ value: v, label: v }))
const title = computed(() => (props.server != null ? t('xray.dnsEditN', { n: props.index + 1 }) : t('xray.dnsAdd')))
</script>

<template>
  <XModal v-model:open="open" :title="title" cancel>
    <div class="dns">
      <XFieldGrid columns="minmax(0, 1fr) 96px 160px" :phone="2">
        <XField v-slot="a" :label="t('xray.outAddress')" class="x-wide-phone">
          <XInput v-bind="a" v-model="form.address" mono />
        </XField>
        <XField v-slot="a" :label="t('xray.port')">
          <XNumber v-bind="a" block v-model="form.port" :min="1" :max="65535" />
        </XField>
        <XField v-slot="a" :label="t('xray.dnsStrategy')">
          <XSelect v-bind="a" v-model="form.queryStrategy" :options="strategies" mono />
        </XField>
      </XFieldGrid>

      <XInputList v-model="form.domains" :title="t('xray.dnsDomains')" mono />
      <XInputList v-model="form.expectIPs" :title="t('xray.dnsExpectIPs')" mono />
      <XInputList v-model="form.unexpectedIPs" :title="t('xray.dnsUnexpectIPs')" mono />

      <XFieldGrid :cols="3">
        <XSwitch v-model="form.skipFallback" boxed label="Skip Fallback" />
        <XSwitch v-model="form.disableCache" boxed label="Disable Cache" />
        <XSwitch v-model="form.finalQuery" boxed label="Final Query" />
      </XFieldGrid>
    </div>

    <template #footer>
      <XButton variant="primary" @click="save">{{ t('xray.confirm') }}</XButton>
    </template>
  </XModal>
</template>

<style scoped>
.dns {
  display: flex;
  flex-direction: column;
  gap: 12px;
  min-width: 0;
}
</style>
