<script setup lang="ts">
import { computed, reactive } from 'vue'
import { useI18n } from 'vue-i18n'
import {
  XExample, XInput, XNumber, XOrderToggles, XSelect, XSettingDivider, XSettingRow, XSettingsSection, XSwitch, XTextarea,
} from '../../kit'
import { useSettingsForm } from './form'
import type { SettingKey } from './model'

/** The subscription service: where it listens, what clients see in their apps, its certificate. */
const { t } = useI18n()
const f = useSettingsForm()

const open = reactive({ general: true, info: true, certs: true, intervals: true })
const sections = {
  general: ['remarkModel', 'subEnable', 'subJsonEnable', 'subClashEnable', 'subListen', 'subDomain', 'subPort', 'subPath', 'subURI'],
  info: ['subEncrypt', 'subShowInfo', 'subTitle', 'subSupportUrl', 'subProfileUrl', 'subAnnounce', 'subEnableRouting', 'subRoutingRules'],
  certs: ['subCertFile', 'subKeyFile'],
  intervals: ['subUpdates'],
} satisfies Record<string, SettingKey[]>

// The remark model is the separator followed by the parts in order, e.g. "-ieo".
type Part = 'i' | 'e' | 'o'
const parts: { value: Part; label: string }[] = [
  { value: 'i', label: 'Inbound' },
  { value: 'e', label: 'Email' },
  { value: 'o', label: 'Other' },
]
const separators = computed(() =>
  [' ', '-', '_', '@', ':', '~', '|', ',', '.', '/'].map((s) => (s === ' ' ? { value: s, label: t('settings.space'), glyph: '␣' } : { value: s, label: s })),
)
const order = computed({
  get: () => (f.draft.remarkModel.length > 1 ? (f.draft.remarkModel.slice(1).split('') as Part[]) : []),
  set: (list: Part[]) => (f.draft.remarkModel = separator.value + list.join('')),
})
const separator = computed({
  get: () => (f.draft.remarkModel.length > 1 ? f.draft.remarkModel.charAt(0) : '-'),
  set: (s: string | null) => (f.draft.remarkModel = (s ?? '-') + f.draft.remarkModel.slice(1)),
})
const sample = computed(() => '#' + order.value.map((p) => parts.find((x) => x.value === p)?.label).join(separator.value))
</script>

<template>
  <div class="settings-tab">
    <XSettingsSection v-model:open="open.general" :title="t('settings.secGeneral')" :count="9" :modified="f.count(sections.general).dirty" :errors="f.count(sections.general).errors">
      <XSettingRow :title="t('settings.remarkModel')" :description="t('settings.remarkModelDesc')" tall :modified="f.isDirty('remarkModel')">
        <div class="remark">
          <XExample :label="t('settings.remarkSample')" :value="sample" />
          <div class="remark__line">
            <XOrderToggles v-model="order" :options="parts" :label="t('settings.remarkParts')" mono />
            <div class="remark__separator">
              <XSelect v-model="separator" :options="separators" :columns="5" mono :title="t('settings.remarkSeparator')" :modified="f.isDirty('remarkModel')" />
            </div>
          </div>
        </div>
      </XSettingRow>
      <XSettingRow :title="t('settings.subEnable')" :description="t('settings.subEnableDesc')" inline :modified="f.isDirty('subEnable')">
        <XSwitch v-model="f.draft.subEnable" :label="t('settings.subEnable')" :modified="f.isDirty('subEnable')" />
      </XSettingRow>
      <XSettingRow :title="t('settings.subJson')" :description="t('settings.subJsonDesc')" inline :modified="f.isDirty('subJsonEnable')">
        <XSwitch v-model="f.draft.subJsonEnable" :label="t('settings.subJson')" :modified="f.isDirty('subJsonEnable')" />
      </XSettingRow>
      <XSettingRow :title="t('settings.subClash')" :description="t('settings.subClashDesc')" inline :modified="f.isDirty('subClashEnable')">
        <XSwitch v-model="f.draft.subClashEnable" :label="t('settings.subClash')" :modified="f.isDirty('subClashEnable')" />
      </XSettingRow>
      <XSettingRow :title="t('settings.subListen')" :description="t('settings.subListenDesc')" :modified="f.isDirty('subListen')">
        <XInput v-model="f.draft.subListen" mono :modified="f.isDirty('subListen')" />
      </XSettingRow>
      <XSettingRow :title="t('settings.subDomain')" :description="t('settings.subDomainDesc')" :modified="f.isDirty('subDomain')">
        <XInput v-model="f.draft.subDomain" mono :modified="f.isDirty('subDomain')" />
      </XSettingRow>
      <XSettingRow :title="t('settings.subPort')" :description="t('settings.subPortDesc')" narrow :modified="f.isDirty('subPort')" :error="f.error('subPort')">
        <XNumber v-model="f.draft.subPort" :modified="f.isDirty('subPort')" :invalid="!!f.error('subPort')" />
      </XSettingRow>
      <XSettingRow :title="t('settings.subPath')" :description="t('settings.subPathDesc')" :modified="f.isDirty('subPath')" :error="f.error('subPath')">
        <XInput v-model="f.draft.subPath" mono placeholder="/sub/" :modified="f.isDirty('subPath')" :invalid="!!f.error('subPath')" />
      </XSettingRow>
      <XSettingRow :title="t('settings.subUri')" :description="t('settings.subUriDesc')" tall :modified="f.isDirty('subURI')">
        <XTextarea v-model="f.draft.subURI" mono :grow="3" placeholder="(http|https)://domain[:port]/path/" :modified="f.isDirty('subURI')" />
      </XSettingRow>
    </XSettingsSection>

    <XSettingsSection v-model:open="open.info" :title="t('settings.secInfo')" :count="8" :modified="f.count(sections.info).dirty" :errors="f.count(sections.info).errors">
      <XSettingRow :title="t('settings.subEncrypt')" :description="t('settings.subEncryptDesc')" inline :modified="f.isDirty('subEncrypt')">
        <XSwitch v-model="f.draft.subEncrypt" :label="t('settings.subEncrypt')" :modified="f.isDirty('subEncrypt')" />
      </XSettingRow>
      <XSettingRow :title="t('settings.subShowInfo')" :description="t('settings.subShowInfoDesc')" inline :modified="f.isDirty('subShowInfo')">
        <XSwitch v-model="f.draft.subShowInfo" :label="t('settings.subShowInfo')" :modified="f.isDirty('subShowInfo')" />
      </XSettingRow>
      <XSettingDivider :label="t('settings.groupMain')" />
      <XSettingRow :title="t('settings.subTitle')" :description="t('settings.subTitleDesc')" :modified="f.isDirty('subTitle')">
        <XInput v-model="f.draft.subTitle" :modified="f.isDirty('subTitle')" />
      </XSettingRow>
      <XSettingRow :title="t('settings.subSupportUrl')" :description="t('settings.subSupportUrlDesc')" :modified="f.isDirty('subSupportUrl')">
        <XInput v-model="f.draft.subSupportUrl" mono :modified="f.isDirty('subSupportUrl')" />
      </XSettingRow>
      <XSettingRow :title="t('settings.subProfileUrl')" :description="t('settings.subProfileUrlDesc')" :modified="f.isDirty('subProfileUrl')">
        <XInput v-model="f.draft.subProfileUrl" mono :modified="f.isDirty('subProfileUrl')" />
      </XSettingRow>
      <XSettingRow :title="t('settings.subAnnounce')" :description="t('settings.subAnnounceDesc')" tall :modified="f.isDirty('subAnnounce')">
        <XTextarea v-model="f.draft.subAnnounce" :rows="3" :modified="f.isDirty('subAnnounce')" />
      </XSettingRow>
      <XSettingDivider :label="t('settings.groupHapp')" />
      <XSettingRow :title="t('settings.subRouting')" :description="t('settings.subRoutingDesc')" inline :modified="f.isDirty('subEnableRouting')">
        <XSwitch v-model="f.draft.subEnableRouting" :label="t('settings.subRouting')" :modified="f.isDirty('subEnableRouting')" />
      </XSettingRow>
      <XSettingRow :title="t('settings.subRoutingRules')" :description="t('settings.subRoutingRulesDesc')" tall :modified="f.isDirty('subRoutingRules')">
        <XTextarea v-model="f.draft.subRoutingRules" :rows="3" mono placeholder="happ://routing/add/..." :modified="f.isDirty('subRoutingRules')" />
      </XSettingRow>
    </XSettingsSection>

    <XSettingsSection v-model:open="open.certs" :title="t('settings.secCerts')" :count="2" :modified="f.count(sections.certs).dirty" :errors="f.count(sections.certs).errors">
      <XSettingRow :title="t('settings.subCertFile')" tall :modified="f.isDirty('subCertFile')" :error="f.error('subCertFile')">
        <XTextarea v-model="f.draft.subCertFile" mono :grow="3" :modified="f.isDirty('subCertFile')" :invalid="!!f.error('subCertFile')" />
      </XSettingRow>
      <XSettingRow :title="t('settings.subKeyFile')" tall :modified="f.isDirty('subKeyFile')" :error="f.error('subKeyFile')">
        <XTextarea v-model="f.draft.subKeyFile" mono :grow="3" :modified="f.isDirty('subKeyFile')" :invalid="!!f.error('subKeyFile')" />
      </XSettingRow>
    </XSettingsSection>

    <XSettingsSection v-model:open="open.intervals" :title="t('settings.secIntervals')" :count="1" :modified="f.count(sections.intervals).dirty" :errors="f.count(sections.intervals).errors">
      <XSettingRow :title="t('settings.subUpdates')" :description="t('settings.subUpdatesDesc')" narrow :modified="f.isDirty('subUpdates')" :error="f.error('subUpdates')">
        <XNumber v-model="f.draft.subUpdates" :unit="t('settings.unitHours')" :modified="f.isDirty('subUpdates')" :invalid="!!f.error('subUpdates')" />
      </XSettingRow>
    </XSettingsSection>
  </div>
</template>

<style scoped>
.settings-tab {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.remark {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.remark__line {
  display: flex;
  gap: 8px;
}

.remark__separator {
  flex: none;
  width: 98px;
}

/* Phones: the parts take the whole line, the separator goes under them. */
@media (max-width: 767px) {
  .remark__line {
    flex-direction: column;
  }
}
</style>
