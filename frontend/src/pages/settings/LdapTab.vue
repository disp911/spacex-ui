<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { api } from '../../api'
import { XInput, XMultiSelect, XNumber, XSettingRow, XSettingsSection, XSwitch, XTextarea } from '../../kit'
import { useSettingsForm } from './form'
import type { SettingKey } from './model'

/** Clients from an LDAP directory: where it is, and how its users become clients. */
const { t } = useI18n()
const f = useSettingsForm()

const open = reactive({ connection: true, sync: true })
const sections = {
  connection: ['ldapEnable', 'ldapHost', 'ldapPort', 'ldapUseTLS', 'ldapBindDN', 'ldapPassword', 'ldapBaseDN', 'ldapUserFilter', 'ldapUserAttr'],
  sync: [
    'ldapVlessField', 'ldapFlagField', 'ldapTruthyValues', 'ldapInvertFlag', 'ldapSyncCron', 'ldapInboundTags',
    'ldapAutoCreate', 'ldapAutoDelete', 'ldapDefaultTotalGB', 'ldapDefaultExpiryDays', 'ldapDefaultLimitIP',
  ],
} satisfies Record<string, SettingKey[]>

// The inbounds to choose from; the setting keeps their tags comma-separated.
const inbounds = ref<{ value: string; label: string }[]>([])
onMounted(async () => {
  try {
    const list = await api.get<{ tag: string; protocol: string; port: number }[]>('panel/api/inbounds/list')
    inbounds.value = (list ?? []).map((i) => ({ value: i.tag, label: `${i.tag} (${i.protocol}@${i.port})` }))
  } catch {
    // nothing to choose from
  }
})
const tags = computed({
  get: () => f.draft.ldapInboundTags.split(',').map((s) => s.trim()).filter(Boolean),
  set: (list: string[]) => (f.draft.ldapInboundTags = list.join(',')),
})
</script>

<template>
  <div class="settings-tab">
    <XSettingsSection v-model:open="open.connection" :title="t('settings.secLdapConnection')" :count="9" :modified="f.count(sections.connection).dirty" :errors="f.count(sections.connection).errors">
      <XSettingRow :title="t('settings.ldapEnable')" inline :modified="f.isDirty('ldapEnable')">
        <XSwitch v-model="f.draft.ldapEnable" :label="t('settings.ldapEnable')" :modified="f.isDirty('ldapEnable')" />
      </XSettingRow>
      <XSettingRow :title="t('settings.ldapHost')" :modified="f.isDirty('ldapHost')">
        <XInput v-model="f.draft.ldapHost" mono :modified="f.isDirty('ldapHost')" />
      </XSettingRow>
      <XSettingRow :title="t('settings.ldapPort')" narrow :modified="f.isDirty('ldapPort')" :error="f.error('ldapPort')">
        <XNumber v-model="f.draft.ldapPort" :modified="f.isDirty('ldapPort')" :invalid="!!f.error('ldapPort')" />
      </XSettingRow>
      <XSettingRow :title="t('settings.ldapUseTls')" inline :modified="f.isDirty('ldapUseTLS')">
        <XSwitch v-model="f.draft.ldapUseTLS" :label="t('settings.ldapUseTls')" :modified="f.isDirty('ldapUseTLS')" />
      </XSettingRow>
      <XSettingRow :title="t('settings.ldapBindDn')" tall :modified="f.isDirty('ldapBindDN')">
        <XTextarea v-model="f.draft.ldapBindDN" mono :grow="3" :modified="f.isDirty('ldapBindDN')" />
      </XSettingRow>
      <XSettingRow :title="t('settings.ldapPassword')" :modified="f.isDirty('ldapPassword')">
        <XInput v-model="f.draft.ldapPassword" type="password" mono autocomplete="off" :modified="f.isDirty('ldapPassword')" />
      </XSettingRow>
      <XSettingRow :title="t('settings.ldapBaseDn')" tall :modified="f.isDirty('ldapBaseDN')">
        <XTextarea v-model="f.draft.ldapBaseDN" mono :grow="3" :modified="f.isDirty('ldapBaseDN')" />
      </XSettingRow>
      <XSettingRow :title="t('settings.ldapUserFilter')" tall :modified="f.isDirty('ldapUserFilter')">
        <XTextarea v-model="f.draft.ldapUserFilter" mono :grow="3" placeholder="(objectClass=person)" :modified="f.isDirty('ldapUserFilter')" />
      </XSettingRow>
      <XSettingRow :title="t('settings.ldapUserAttr')" :modified="f.isDirty('ldapUserAttr')">
        <XInput v-model="f.draft.ldapUserAttr" mono :modified="f.isDirty('ldapUserAttr')" />
      </XSettingRow>
    </XSettingsSection>

    <XSettingsSection v-model:open="open.sync" :title="t('settings.secLdapSync')" :count="11" :modified="f.count(sections.sync).dirty" :errors="f.count(sections.sync).errors">
      <XSettingRow :title="t('settings.ldapVlessField')" :modified="f.isDirty('ldapVlessField')">
        <XInput v-model="f.draft.ldapVlessField" mono :modified="f.isDirty('ldapVlessField')" />
      </XSettingRow>
      <XSettingRow :title="t('settings.ldapFlagField')" :description="t('settings.ldapFlagFieldDesc')" :modified="f.isDirty('ldapFlagField')">
        <XInput v-model="f.draft.ldapFlagField" mono :modified="f.isDirty('ldapFlagField')" />
      </XSettingRow>
      <XSettingRow :title="t('settings.ldapTruthy')" :description="t('settings.ldapTruthyDesc')" :modified="f.isDirty('ldapTruthyValues')">
        <XInput v-model="f.draft.ldapTruthyValues" mono placeholder="true,1,yes,on" :modified="f.isDirty('ldapTruthyValues')" />
      </XSettingRow>
      <XSettingRow :title="t('settings.ldapInvert')" :description="t('settings.ldapInvertDesc')" inline :modified="f.isDirty('ldapInvertFlag')">
        <XSwitch v-model="f.draft.ldapInvertFlag" :label="t('settings.ldapInvert')" :modified="f.isDirty('ldapInvertFlag')" />
      </XSettingRow>
      <XSettingRow :title="t('settings.ldapCron')" :description="t('settings.ldapCronDesc')" :modified="f.isDirty('ldapSyncCron')">
        <XInput v-model="f.draft.ldapSyncCron" mono placeholder="@every 1m" :modified="f.isDirty('ldapSyncCron')" />
      </XSettingRow>
      <XSettingRow :title="t('settings.ldapInbounds')" :description="t('settings.ldapInboundsDesc')" tall :modified="f.isDirty('ldapInboundTags')">
        <XMultiSelect
          v-model="tags"
          :options="inbounds"
          mono
          :disabled="!inbounds.length"
          :placeholder="inbounds.length ? t('settings.ldapInboundsPick') : t('settings.ldapInboundsNone')"
          :modified="f.isDirty('ldapInboundTags')"
        />
      </XSettingRow>
      <XSettingRow :title="t('settings.ldapAutoCreate')" inline :modified="f.isDirty('ldapAutoCreate')">
        <XSwitch v-model="f.draft.ldapAutoCreate" :label="t('settings.ldapAutoCreate')" :modified="f.isDirty('ldapAutoCreate')" />
      </XSettingRow>
      <XSettingRow :title="t('settings.ldapAutoDelete')" inline :modified="f.isDirty('ldapAutoDelete')">
        <XSwitch v-model="f.draft.ldapAutoDelete" :label="t('settings.ldapAutoDelete')" :modified="f.isDirty('ldapAutoDelete')" />
      </XSettingRow>
      <XSettingRow :title="t('settings.ldapTotalGb')" narrow :modified="f.isDirty('ldapDefaultTotalGB')" :error="f.error('ldapDefaultTotalGB')">
        <XNumber v-model="f.draft.ldapDefaultTotalGB" :unit="t('settings.unitGb')" :modified="f.isDirty('ldapDefaultTotalGB')" :invalid="!!f.error('ldapDefaultTotalGB')" />
      </XSettingRow>
      <XSettingRow :title="t('settings.ldapExpiryDays')" narrow :modified="f.isDirty('ldapDefaultExpiryDays')" :error="f.error('ldapDefaultExpiryDays')">
        <XNumber v-model="f.draft.ldapDefaultExpiryDays" :unit="t('settings.unitDays')" :modified="f.isDirty('ldapDefaultExpiryDays')" :invalid="!!f.error('ldapDefaultExpiryDays')" />
      </XSettingRow>
      <XSettingRow :title="t('settings.ldapLimitIp')" narrow :modified="f.isDirty('ldapDefaultLimitIP')" :error="f.error('ldapDefaultLimitIP')">
        <XNumber v-model="f.draft.ldapDefaultLimitIP" :modified="f.isDirty('ldapDefaultLimitIP')" :invalid="!!f.error('ldapDefaultLimitIP')" />
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
</style>
