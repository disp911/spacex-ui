<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { XInput, XNumber, XSettingRow, XSettingsSection, XSwitch, XTextarea } from '../../kit'
import AddressNote from './AddressNote.vue'
import { useSettingsForm } from './form'
import { panelAddress, type SettingKey } from './model'
import { loadSite } from './site'

/** The panel's address, its certificate and what it does with clients. */
const { t } = useI18n()
const f = useSettingsForm()

const open = reactive({ general: true, clients: true })
const sections = {
  general: ['webListen', 'webDomain', 'webPort', 'webBasePath', 'sessionMaxAge'],
  clients: ['restartXrayOnClientDisable', 'externalTrafficInformEnable', 'externalTrafficInformURI'],
} satisfies Record<string, SettingKey[]>

// A changed part of the address shows where the panel will be after a restart.
const address = computed(() => panelAddress(f.draft))
const moved = (k: SettingKey) => f.isDirty(k) && !f.error(k)

// Behind REALITY the panel answers on 443 with its website and opens no port of its own.
const reality = ref(false)
onMounted(async () => {
  try {
    reality.value = (await loadSite()).mode === 'reality'
  } catch {
    // the port shows as it is stored
  }
})
</script>

<template>
  <div class="settings-tab">
    <XSettingsSection v-model:open="open.general" :title="t('settings.secGeneral')" :count="5" :modified="f.count(sections.general).dirty" :errors="f.count(sections.general).errors">
      <XSettingRow :title="t('settings.webListen')" :description="t('settings.webListenDesc')" :modified="f.isDirty('webListen')">
        <XInput v-model="f.draft.webListen" mono :modified="f.isDirty('webListen')" />
        <template v-if="moved('webListen')" #note><AddressNote :address="address" /></template>
      </XSettingRow>
      <XSettingRow :title="t('settings.webDomain')" :description="t('settings.webDomainDesc')" :modified="f.isDirty('webDomain')">
        <XInput v-model="f.draft.webDomain" mono :modified="f.isDirty('webDomain')" />
        <template v-if="moved('webDomain')" #note><AddressNote :address="address" /></template>
      </XSettingRow>
      <XSettingRow v-if="reality" :title="t('settings.webPort')" :description="t('settings.webPortDesc')" narrow>
        <XNumber :model-value="443" readonly />
      </XSettingRow>
      <XSettingRow v-else :title="t('settings.webPort')" :description="t('settings.webPortDesc')" narrow :modified="f.isDirty('webPort')" :error="f.error('webPort')">
        <XNumber v-model="f.draft.webPort" :modified="f.isDirty('webPort')" :invalid="!!f.error('webPort')" />
        <template v-if="moved('webPort')" #note><AddressNote :address="address" /></template>
      </XSettingRow>
      <XSettingRow :title="t('settings.webBasePath')" :description="t('settings.pathBothDesc')" :modified="f.isDirty('webBasePath')" :error="f.error('webBasePath')">
        <XInput v-model="f.draft.webBasePath" mono :modified="f.isDirty('webBasePath')" :invalid="!!f.error('webBasePath')" />
        <template v-if="moved('webBasePath')" #note><AddressNote :address="address" /></template>
      </XSettingRow>
      <XSettingRow :title="t('settings.sessionMaxAge')" :description="t('settings.sessionMaxAgeDesc')" narrow :modified="f.isDirty('sessionMaxAge')" :error="f.error('sessionMaxAge')">
        <XNumber v-model="f.draft.sessionMaxAge" :unit="t('settings.unitMinutes')" :modified="f.isDirty('sessionMaxAge')" :invalid="!!f.error('sessionMaxAge')" />
      </XSettingRow>
    </XSettingsSection>

    <XSettingsSection v-model:open="open.clients" :title="t('settings.secClients')" :count="3" :modified="f.count(sections.clients).dirty" :errors="f.count(sections.clients).errors">
      <XSettingRow :title="t('settings.restartXray')" :description="t('settings.restartXrayDesc')" inline :modified="f.isDirty('restartXrayOnClientDisable')">
        <XSwitch v-model="f.draft.restartXrayOnClientDisable" :label="t('settings.restartXray')" :modified="f.isDirty('restartXrayOnClientDisable')" />
      </XSettingRow>
      <XSettingRow :title="t('settings.externalTraffic')" :description="t('settings.externalTrafficDesc')" inline :modified="f.isDirty('externalTrafficInformEnable')">
        <XSwitch v-model="f.draft.externalTrafficInformEnable" :label="t('settings.externalTraffic')" :modified="f.isDirty('externalTrafficInformEnable')" />
      </XSettingRow>
      <XSettingRow :title="t('settings.externalTrafficUri')" :description="t('settings.externalTrafficUriDesc')" tall :modified="f.isDirty('externalTrafficInformURI')">
        <XTextarea v-model="f.draft.externalTrafficInformURI" mono :grow="3" placeholder="(http|https)://domain[:port]/path/" :modified="f.isDirty('externalTrafficInformURI')" />
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
