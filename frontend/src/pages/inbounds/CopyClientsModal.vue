<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { api } from '../../api'
import {
  toast, XButton, XButtonRow, XCheckbox, XChip, XChipBox, XField, XModal, XPlainTable, XSelect, type Option, type PlainColumn,
} from '../../kit'
import { formatBytes, formatDate } from '../../i18n/format'
import { useInbounds } from '../../stores/inbounds'
import { cause } from '../settings/form'
import { canFlow, flows, fromDb } from './inbound'
import { clientsOf, delayedDays, statsMap, usedOf, type DbInbound } from './model'

/**
 * Window 14 of the mockup: clients of another inbound copied into this one,
 * as the old panel's copy form did it. The source is picked first, then its
 * clients; the server adds them under new emails (email_<inbound id>) and
 * answers with what it added and what failed. VLESS targets may set one flow
 * for all of them.
 */
const props = defineProps<{ inbound: DbInbound | null }>()
const open = defineModel<boolean>('open', { default: false })
const { t } = useI18n()
const store = useInbounds()

const sourceId = ref<number | null>(null)
const selected = ref<string[]>([])
const flow = ref('')
const saving = ref(false)

watch(open, (on) => {
  if (!on) return
  sourceId.value = null
  selected.value = []
  flow.value = ''
})
watch(sourceId, () => (selected.value = []))

/** Other inbounds with clients of their own. */
const sources = computed<Option<number>[]>(() =>
  store.list
    .filter((ib) => ib.id !== props.inbound?.id && clientsOf(ib).length > 0)
    .map((ib) => ({ value: ib.id, label: `${ib.remark} (${ib.protocol}, ${clientsOf(ib).length})` })),
)
const source = computed(() => store.list.find((ib) => ib.id === sourceId.value) ?? null)

interface Row {
  email: string
  traffic: string
  expiry: string
}
const rows = computed<Row[]>(() => {
  const ib = source.value
  if (!ib) return []
  const stats = statsMap(ib)
  return clientsOf(ib).map((c) => {
    const days = delayedDays(c.expiryTime ?? 0)
    const expiry = (c.expiryTime ?? 0) > 0 ? formatDate(c.expiryTime!) : days !== null ? t('inbounds.periodDays', { n: days }) : t('inbounds.unlimited')
    return { email: c.email, traffic: formatBytes(usedOf(stats.get(c.email))), expiry }
  })
})
const columns = computed<PlainColumn[]>(() => [
  { key: 'email', label: 'Email' },
  { key: 'traffic', label: t('inbounds.colTraffic') },
  { key: 'expiry', label: t('inbounds.colExpiry') },
])

const allOn = computed(() => rows.value.length > 0 && selected.value.length === rows.value.length)
const someOn = computed(() => selected.value.length > 0 && !allOn.value)
const selectAll = () => (selected.value = rows.value.map((r) => r.email))
const clearAll = () => (selected.value = [])
function toggle(email: string, on: boolean) {
  selected.value = on ? [...selected.value, email] : selected.value.filter((e) => e !== email)
}

const showFlow = computed(() => !!props.inbound && canFlow(fromDb(props.inbound)))
const flowOptions = computed<Option[]>(() => [{ value: '', label: t('inbounds.none') }, ...flows.map((f) => ({ value: f, label: f }))])
const preview = computed(() => selected.value.map((e) => `${e}_${props.inbound?.id}`))

async function copy() {
  if (!props.inbound || saving.value) return
  if (!sourceId.value) return toast('err', t('inbounds.copySelectSourceFirst'))
  saving.value = true
  try {
    const res = await api.post<{ added?: string[]; errors?: string[] }>(`panel/api/inbounds/${props.inbound.id}/copyClients`, {
      sourceInboundId: sourceId.value,
      clientEmails: selected.value,
      ...(showFlow.value && flow.value ? { flow: flow.value } : {}),
    })
    const added = res?.added?.length ?? 0
    if (added > 0) toast('ok', `${t('inbounds.copyResultSuccess')}: ${added}`)
    else toast('warn', t('inbounds.copyResultNone'))
    if (res?.errors?.length) toast('err', `${t('inbounds.copyResultErrors')}: ${res.errors.join('; ')}`)
    open.value = false
    await store.load().catch(() => {})
  } catch (e) {
    toast('err', cause(e))
  } finally {
    saving.value = false
  }
}
</script>

<template>
  <XModal v-model:open="open" :title="`${t('inbounds.copyTo')} ${inbound?.remark ?? ''}`" size="table" cancel scroll-body :dismissible="!saving">
    <div class="copy">
      <XField v-slot="a" :label="t('inbounds.copySource')">
        <XSelect v-bind="a" v-model="sourceId" :options="sources" />
      </XField>

      <template v-if="source">
        <XButtonRow>
          <XButton size="sm" @click="selectAll">{{ t('inbounds.selectAll') }}</XButton>
          <XButton size="sm" @click="clearAll">{{ t('inbounds.clearAll') }}</XButton>
        </XButtonRow>
        <XPlainTable :columns="columns" :rows="rows" :row-key="(r) => r.email" :height="282">
          <template #head-email>
            <XCheckbox :model-value="allOn" :indeterminate="someOn" @update:model-value="$event ? selectAll() : clearAll()">Email</XCheckbox>
          </template>
          <template #cell-email="{ row }">
            <XCheckbox :model-value="selected.includes(row.email)" @update:model-value="toggle(row.email, $event)">{{ row.email }}</XCheckbox>
          </template>
          <template #cell-traffic="{ row }">{{ row.traffic }}</template>
          <template #cell-expiry="{ row }">{{ row.expiry }}</template>
        </XPlainTable>
      </template>

      <XField v-if="showFlow" v-slot="a" :label="t('inbounds.copyFlowLabel')" :hint="t('inbounds.copyFlowHint')">
        <XSelect v-bind="a" v-model="flow" :options="flowOptions" mono />
      </XField>

      <XField v-if="preview.length" :label="t('inbounds.copyEmailPreview')">
        <XChipBox :label="t('inbounds.copyEmailPreview')">
          <XChip v-for="e in preview" :key="e" mono>{{ e }}</XChip>
        </XChipBox>
      </XField>
    </div>

    <template #footer>
      <XButton variant="primary" :loading="saving" @click="copy">{{ t('inbounds.copySelected') }}</XButton>
    </template>
  </XModal>
</template>

<style scoped>
.copy {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
</style>
