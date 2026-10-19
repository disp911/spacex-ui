<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { api } from '../../api'
import { toast, XButton, XDateTime, XField, XFieldGrid, XInput, XModal, XNumber, XSelect, XSwitch, type Option } from '../../kit'
import { formatDate } from '../../i18n/format'
import { useInbounds } from '../../stores/inbounds'
import { cause } from '../settings/form'
import { fromDb, newInbound, setProtocol, toForm, trafficResets, type InboundForm } from './inbound'
import { protocols, type DbInbound, type Protocol } from './model'

/**
 * Window 11 of the mockup: a new inbound, or the stored one to change. It
 * opens centred and its body scrolls, as every window. The protocol of a
 * stored inbound stays; the panel serves one MTProto inbound only.
 */
const props = defineProps<{ inbound: DbInbound | null }>()
const open = defineModel<boolean>('open', { default: false })
const { t } = useI18n()
const store = useInbounds()

const GB = 1024 ** 3
const f = ref<InboundForm>(newInbound())
const saving = ref(false)

watch(open, (on) => {
  if (on) f.value = props.inbound ? fromDb(props.inbound) : newInbound()
})

const editing = computed(() => !!props.inbound)
const mtprotoTaken = computed(() => !editing.value && store.list.some((ib) => ib.protocol === 'mtproto'))

const protocolOptions = computed<Option[]>(() =>
  protocols.map((p) => (p === 'mtproto' && mtprotoTaken.value ? { value: p, label: p, disabled: true, note: t('inbounds.mtprotoSingle') } : { value: p, label: p })),
)
const protocol = computed({
  get: () => f.value.protocol as string,
  set: (p: string) => setProtocol(f.value, p as Protocol),
})

const resetOptions = computed<Option[]>(() => {
  const labels: Record<string, string> = {
    never: t('inbounds.resetNever'),
    hourly: t('inbounds.resetHourly'),
    daily: t('inbounds.resetDaily'),
    weekly: t('inbounds.resetWeekly'),
    monthly: t('inbounds.resetMonthly'),
  }
  return trafficResets.map((r) => ({ value: r, label: labels[r] ?? r }))
})
const resetHelp = computed(() =>
  f.value.lastTrafficResetTime > 0
    ? `${t('inbounds.trafficResetHelp')}. ${t('inbounds.lastReset')}: ${formatDate(f.value.lastTrafficResetTime)}`
    : t('inbounds.trafficResetHelp'),
)

/** The limit in gigabytes, as the old panel showed it, to two places. */
const totalGb = computed({
  get: () => Math.round((f.value.total / GB) * 100) / 100,
  set: (gb: number | null) => (f.value.total = Math.round((gb ?? 0) * GB)),
})

const portValid = computed(() => !!f.value.port && f.value.port >= 1 && f.value.port <= 65535)
const valid = computed(() => portValid.value && !(f.value.protocol === 'mtproto' && mtprotoTaken.value))

async function save() {
  if (!valid.value || saving.value) return
  saving.value = true
  try {
    await api.post(editing.value ? `panel/api/inbounds/update/${props.inbound!.id}` : 'panel/api/inbounds/add', toForm(f.value))
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
  <XModal v-model:open="open" :title="editing ? t('inbounds.editTitle') : t('inbounds.addTitle')" size="lg" cancel :dismissible="!saving">
    <div class="inbound">
      <XSwitch v-model="f.enable" boxed :label="t('inbounds.enable')" />
      <XFieldGrid>
        <XField v-slot="a" :label="t('inbounds.colRemark')">
          <XInput v-bind="a" v-model.trim="f.remark" />
        </XField>
        <XField v-slot="a" :label="t('inbounds.colProtocol')">
          <XSelect v-bind="a" v-model="protocol" :options="protocolOptions" :disabled="editing" mono />
        </XField>
      </XFieldGrid>
      <XFieldGrid>
        <XField v-slot="a" :label="t('inbounds.listen')" :help="t('inbounds.listenHelp')">
          <XInput v-bind="a" v-model.trim="f.listen" mono />
        </XField>
        <XField v-slot="a" :label="t('inbounds.colPort')">
          <XNumber v-bind="a" v-model="f.port" block :min="1" :max="65535" :invalid="!portValid" />
        </XField>
      </XFieldGrid>
      <XFieldGrid :cols="3">
        <XField v-slot="a" :label="t('inbounds.totalFlow')" :help="t('inbounds.totalHelp')">
          <XNumber v-bind="a" v-model="totalGb" block decimal :min="0" :unit="t('inbounds.gb')" />
        </XField>
        <XField v-slot="a" :label="t('inbounds.trafficReset')" :help="resetHelp">
          <XSelect v-bind="a" v-model="f.trafficReset" :options="resetOptions" />
        </XField>
        <XField v-slot="a" :label="t('inbounds.colExpiry')" :help="t('inbounds.expiryHelp')">
          <XDateTime v-bind="a" v-model="f.expiryTime" />
        </XField>
      </XFieldGrid>
    </div>

    <template #footer>
      <XButton variant="primary" :loading="saving" :disabled="!valid" @click="save">{{ editing ? t('inbounds.update') : t('inbounds.createShort') }}</XButton>
    </template>
  </XModal>
</template>

<style scoped>
.inbound {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
</style>
