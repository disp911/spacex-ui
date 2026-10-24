<script setup lang="ts">
import { computed, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { XButton, XInput, XSelect, XSettingRow, XSettingsSection, XSwitch, type Option } from '../../kit'
import { useXrayForm } from './form'
import { defaultTestUrl, fields, options, type FieldKey } from './model'

/** Strategies, statistics and logs of Xray, and the way back to the default template. */
const { t } = useI18n()
const f = useXrayForm()

const open = reactive({ general: true, stats: true, logs: true, reset: true })
const sections = {
  general: ['freedom', 'routing', 'testUrl'],
  stats: ['statsInboundUplink', 'statsInboundDownlink', 'statsOutboundUplink', 'statsOutboundDownlink'],
  logs: ['logLevel', 'access', 'error', 'mask', 'dnsLog'],
} satisfies Record<string, FieldKey[]>

/** A field of the template to bind a control to. */
function field<K extends keyof typeof fields>(key: K) {
  type T = ReturnType<(typeof fields)[K]['get']>
  return computed<T>({
    get: () => fields[key].get(f.draft.config) as T,
    set: (v) => (fields[key].set as (c: object, v: T) => void)(f.draft.config, v),
  })
}

const v = reactive({
  freedom: field('freedom'),
  routing: field('routing'),
  statsInboundUplink: field('statsInboundUplink'),
  statsInboundDownlink: field('statsInboundDownlink'),
  statsOutboundUplink: field('statsOutboundUplink'),
  statsOutboundDownlink: field('statsOutboundDownlink'),
  logLevel: field('logLevel'),
  access: field('access'),
  error: field('error'),
  mask: field('mask'),
  dnsLog: field('dnsLog'),
})

// An empty value leaves the setting to Xray; the list shows it as a word, not a code.
const list = (values: string[]): Option[] => values.map((value) => (value ? { value, label: value } : { value, label: 'Empty', plain: true }))

const stats = computed(() => [
  { key: 'statsInboundUplink' as const, title: t('xray.statsInboundUplink'), description: t('xray.statsInboundUplinkDesc') },
  { key: 'statsInboundDownlink' as const, title: t('xray.statsInboundDownlink'), description: t('xray.statsInboundDownlinkDesc') },
  { key: 'statsOutboundUplink' as const, title: t('xray.statsOutboundUplink'), description: t('xray.statsOutboundUplinkDesc') },
  { key: 'statsOutboundDownlink' as const, title: t('xray.statsOutboundDownlink'), description: t('xray.statsOutboundDownlinkDesc') },
])

const resetting = ref(false)
async function reset() {
  resetting.value = true
  try {
    await f.resetToDefault()
  } finally {
    resetting.value = false
  }
}
</script>

<template>
  <div class="xray-tab">
    <XSettingsSection v-model:open="open.general" :title="t('xray.secGeneral')" :count="3" :modified="f.count(sections.general)">
      <XSettingRow :title="t('xray.freedom')" :description="t('xray.freedomDesc')" :modified="f.isDirty('freedom')">
        <XSelect v-model="v.freedom" :options="list(options.freedom)" mono :modified="f.isDirty('freedom')" />
      </XSettingRow>
      <XSettingRow :title="t('xray.routing')" :description="t('xray.routingDesc')" :modified="f.isDirty('routing')">
        <XSelect v-model="v.routing" :options="list(options.routing)" mono :modified="f.isDirty('routing')" />
      </XSettingRow>
      <XSettingRow :title="t('xray.testUrl')" :description="t('xray.testUrlDesc')" :modified="f.isDirty('testUrl')">
        <XInput v-model="f.draft.testUrl" mono :placeholder="defaultTestUrl" :modified="f.isDirty('testUrl')" />
      </XSettingRow>
    </XSettingsSection>

    <XSettingsSection v-model:open="open.stats" :title="t('xray.secStats')" :count="4" :modified="f.count(sections.stats)">
      <XSettingRow v-for="s in stats" :key="s.key" :title="s.title" :description="s.description" inline :modified="f.isDirty(s.key)">
        <XSwitch v-model="v[s.key]" :label="s.title" :modified="f.isDirty(s.key)" />
      </XSettingRow>
    </XSettingsSection>

    <XSettingsSection v-model:open="open.logs" :title="t('xray.secLogs')" :count="5" :modified="f.count(sections.logs)">
      <XSettingRow :title="t('xray.logLevel')" :description="t('xray.logLevelDesc')" :modified="f.isDirty('logLevel')">
        <XSelect v-model="v.logLevel" :options="list(options.logLevel)" mono :modified="f.isDirty('logLevel')" />
      </XSettingRow>
      <XSettingRow :title="t('xray.accessLog')" :description="t('xray.accessLogDesc')" :modified="f.isDirty('access')">
        <XSelect v-model="v.access" :options="list(options.access)" mono :modified="f.isDirty('access')" />
      </XSettingRow>
      <XSettingRow :title="t('xray.errorLog')" :description="t('xray.errorLogDesc')" :modified="f.isDirty('error')">
        <XSelect v-model="v.error" :options="list(options.error)" mono :modified="f.isDirty('error')" />
      </XSettingRow>
      <XSettingRow :title="t('xray.maskAddress')" :description="t('xray.maskAddressDesc')" :modified="f.isDirty('mask')">
        <XSelect v-model="v.mask" :options="list(options.mask)" mono :modified="f.isDirty('mask')" />
      </XSettingRow>
      <XSettingRow :title="t('xray.dnsLog')" :description="t('xray.dnsLogDesc')" inline :modified="f.isDirty('dnsLog')">
        <XSwitch v-model="v.dnsLog" :label="t('xray.dnsLog')" :modified="f.isDirty('dnsLog')" />
      </XSettingRow>
    </XSettingsSection>

    <XSettingsSection v-model:open="open.reset" :title="t('xray.secReset')" padded>
      <XButton variant="danger" icon="reset" :loading="resetting" @click="reset">{{ t('xray.secResetCta') }}</XButton>
    </XSettingsSection>
  </div>
</template>

<style scoped>
.xray-tab {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
</style>
