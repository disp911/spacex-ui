<script setup lang="ts">
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { XField, XFieldGrid, XInput, XMultiSelect, XSection, XSegmented, XSelect, type Option } from '../../kit'
import type { Outbound } from './model'
import { alpns, canReality, defaultReality, defaultTls, fingerprints } from './outbound'

/** The security of the outbound window: none, TLS or Reality in the card's head, and its settings. */
const props = defineProps<{ o: Outbound }>()
const { t } = useI18n()
const open = ref(true)

const st = computed(() => (props.o.streamSettings ??= { network: 'tcp', security: 'none' }))
const sec = computed({
  get: () => st.value.security ?? 'none',
  set: (v: string) => {
    st.value.security = v
    if (v === 'tls') st.value.tlsSettings ??= defaultTls()
    if (v === 'reality') st.value.realitySettings ??= defaultReality()
  },
})
const tls = computed(() => (st.value.tlsSettings ??= defaultTls()))
const reality = computed(() => (st.value.realitySettings ??= defaultReality()))

const kinds = computed(() => [
  { value: 'none', label: 'none' },
  { value: 'tls', label: 'TLS' },
  ...(canReality(props.o) ? [{ value: 'reality', label: 'Reality' }] : []),
])
const fingerprintOptions = (empty: boolean): Option[] => [
  ...(empty ? [{ value: '', label: 'Empty', plain: true }] : []),
  ...fingerprints.map((f) => ({ value: f, label: f })),
]
const alpnOptions = alpns.map((a) => ({ value: a, label: a }))
</script>

<template>
  <XSection v-model:open="open" :title="t('xray.security')" card>
    <template #meta>
      <div class="security__kinds">
        <XSegmented v-model="sec" :options="kinds" :label="t('xray.security')" block />
      </div>
    </template>

    <template v-if="sec === 'tls'" #default>
      <XFieldGrid :cols="2">
        <XField v-slot="a" label="SNI"><XInput v-bind="a" v-model="tls.serverName" mono /></XField>
        <XField v-slot="a" label="uTLS"><XSelect v-bind="a" :model-value="tls.fingerprint ?? ''" @update:model-value="tls.fingerprint = String($event)" :options="fingerprintOptions(true)" mono /></XField>
      </XFieldGrid>
      <XField v-slot="a" label="ALPN"><XMultiSelect v-bind="a" v-model="tls.alpn" :options="alpnOptions" mono /></XField>
      <XField v-slot="a" label="ECH Config List"><XInput v-bind="a" v-model="tls.echConfigList" mono /></XField>
      <XFieldGrid :cols="2">
        <XField v-slot="a" label="verify Peer Cert By Name"><XInput v-bind="a" v-model="tls.verifyPeerCertByName" mono placeholder="cloudflare-dns.com" /></XField>
        <XField v-slot="a" label="pinned Peer Cert Sha256">
          <XInput v-bind="a" v-model="tls.pinnedPeerCertSha256" mono placeholder="Enter SHA256 fingerprints (base64)" />
        </XField>
      </XFieldGrid>
    </template>

    <template v-else-if="sec === 'reality'" #default>
      <XFieldGrid :cols="2">
        <XField v-slot="a" label="SNI"><XInput v-bind="a" v-model="reality.serverName" mono /></XField>
        <XField v-slot="a" label="uTLS"><XSelect v-bind="a" v-model="reality.fingerprint" :options="fingerprintOptions(false)" mono /></XField>
        <XField v-slot="a" label="Short ID"><XInput v-bind="a" v-model="reality.shortId" mono /></XField>
        <XField v-slot="a" label="SpiderX"><XInput v-bind="a" v-model="reality.spiderX" mono /></XField>
      </XFieldGrid>
      <XField v-slot="a" label="Public Key"><XInput v-bind="a" v-model="reality.publicKey" mono /></XField>
      <XField v-slot="a" label="mldsa65 Verify"><XInput v-bind="a" v-model="reality.mldsa65Verify" mono /></XField>
    </template>
  </XSection>
</template>

<style scoped>
.security__kinds {
  width: 250px;
}

@media (max-width: 767px) {
  .security__kinds {
    width: 190px;
  }
}
</style>
