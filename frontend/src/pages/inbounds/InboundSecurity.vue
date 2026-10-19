<script setup lang="ts">
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import {
  toast, XAddButton, XButton, XButtonRow, XField, XFieldGrid, XInlineButton, XInput, XMultiSelect, XNestCard, XNumber, XPartHead, XSegmented, XSelect, XSwitch,
  XTextarea, type Option,
} from '../../kit'
import { useInbounds } from '../../stores/inbounds'
import { cause } from '../settings/form'
import {
  alpns, canReality, canTls, certUsages, cipherSuites, defaultCert, fingerprints, newEch, newMldsa65, newX25519, randomRealityTarget, randomShortIds,
  security, setSecurity, tlsVersions, type InboundForm,
} from './inbound'

/**
 * The security part of window 11, as the old panel's TLS and Reality forms
 * had it: none, Reality or TLS in the head (Hysteria is always TLS); TLS with
 * its certificates by path or content and ECH; Reality with its target, short
 * IDs and keys. Choosing Reality fetches a key pair at once (mockup 11a).
 */
const props = defineProps<{ f: InboundForm }>()
const { t } = useI18n()
const store = useInbounds()

type Json = Record<string, any>
const st = computed<Json>(() => props.f.streamSettings ?? {})
const tls = computed<Json>(() => st.value.tlsSettings ?? {})
const reality = computed<Json>(() => st.value.realitySettings ?? {})
const hysteria = computed(() => props.f.protocol === 'hysteria')

/** Whether the transport takes any security: none of it for mKCP, as on the old panel. */
const shown = computed(() => canTls(props.f))

const sec = computed({
  get: () => security(props.f),
  set: (v: string) => {
    setSecurity(props.f, v)
    if (v === 'reality' && !reality.value.privateKey) void getX25519()
  },
})
const kinds = computed(() => [
  { value: 'none', label: t('inbounds.none') },
  ...(canReality(props.f) ? [{ value: 'reality', label: 'Reality' }] : []),
  { value: 'tls', label: 'TLS' },
])

const list = (values: readonly string[]): Option[] => values.map((v) => ({ value: v, label: v }))
const cipherOptions: Option[] = [{ value: '', label: 'Auto' }, ...list(cipherSuites)]
const tlsFingerprints: Option[] = [{ value: '', label: 'None' }, ...list(fingerprints)]
const alpnOptions = list(alpns)

// --- Keys from the server ---------------------------------------------------------------

const busy = ref('')

async function fetching(what: string, run: () => Promise<void>) {
  busy.value = what
  try {
    await run()
  } catch (e) {
    toast('err', cause(e))
  } finally {
    busy.value = ''
  }
}

const getX25519 = () =>
  fetching('x25519', async () => {
    const keys = await newX25519()
    reality.value.privateKey = keys.privateKey
    reality.value.settings.publicKey = keys.publicKey
  })
const getMldsa65 = () =>
  fetching('mldsa65', async () => {
    const keys = await newMldsa65()
    reality.value.mldsa65Seed = keys.seed
    reality.value.settings.mldsa65Verify = keys.verify
  })
const getEch = () =>
  fetching('ech', async () => {
    const keys = await newEch(tls.value.serverName)
    tls.value.echServerKeys = keys.echServerKeys
    tls.value.settings.echConfigList = keys.echConfigList
  })

function clearX25519() {
  reality.value.privateKey = ''
  reality.value.settings.publicKey = ''
}
function clearMldsa65() {
  reality.value.mldsa65Seed = ''
  reality.value.settings.mldsa65Verify = ''
}
function clearEch() {
  tls.value.echServerKeys = ''
  tls.value.settings.echConfigList = ''
}

// --- Reality ------------------------------------------------------------------------------

function newTarget() {
  const { target, sni } = randomRealityTarget()
  reality.value.target = target
  reality.value.serverNames = [sni]
}

// Typed as text and kept as typed; the lists are cleaned up on save.
const serverNames = computed({
  get: () => (Array.isArray(reality.value.serverNames) ? reality.value.serverNames.join(',') : String(reality.value.serverNames ?? '')),
  set: (v: string) => (reality.value.serverNames = v.split(',')),
})
const shortIds = computed({
  get: () => (Array.isArray(reality.value.shortIds) ? reality.value.shortIds.join('\n') : String(reality.value.shortIds ?? '')),
  set: (v: string) => (reality.value.shortIds = v.split('\n')),
})

// --- TLS certificates ------------------------------------------------------------------------

const certs = computed<Json[]>(() => (Array.isArray(tls.value.certificates) ? tls.value.certificates : []))
const certSources = computed(() => [
  { value: true, label: t('inbounds.certificatePath') },
  { value: false, label: t('inbounds.certificateContent') },
])

function panelCert(c: Json) {
  c.certificateFile = store.defaults?.defaultCert ?? ''
  c.keyFile = store.defaults?.defaultKey ?? ''
}
</script>

<template>
  <template v-if="shown">
    <XPartHead v-if="hysteria" title="TLS" />
    <XPartHead v-else :title="t('inbounds.security')" wide>
      <XSegmented v-model="sec" :options="kinds" :label="t('inbounds.security')" block />
    </XPartHead>

    <!-- TLS -->
    <template v-if="sec === 'tls'">
      <XFieldGrid>
        <XField v-slot="a" label="SNI">
          <XInput v-bind="a" v-model.trim="tls.serverName" mono />
        </XField>
        <XField v-slot="a" label="Cipher Suites">
          <XSelect v-bind="a" v-model="tls.cipherSuites" :options="cipherOptions" />
        </XField>
      </XFieldGrid>
      <XFieldGrid :cols="3" :phone="2">
        <XField v-slot="a" label="Min Version">
          <XSelect v-bind="a" v-model="tls.minVersion" :options="list(tlsVersions)" />
        </XField>
        <XField v-slot="a" label="Max Version">
          <XSelect v-bind="a" v-model="tls.maxVersion" :options="list(tlsVersions)" />
        </XField>
        <XField v-slot="a" label="uTLS">
          <XSelect v-bind="a" v-model="tls.settings.fingerprint" :options="tlsFingerprints" />
        </XField>
      </XFieldGrid>
      <XField v-slot="a" label="ALPN">
        <XMultiSelect v-bind="a" v-model="tls.alpn" :options="alpnOptions" mono />
      </XField>
      <XFieldGrid :cols="3">
        <XSwitch v-model="tls.rejectUnknownSni" boxed label="Reject Unknown SNI" />
        <XSwitch v-model="tls.disableSystemRoot" boxed label="Disable System Root" />
        <XSwitch v-model="tls.enableSessionResumption" boxed label="Session Resumption" />
      </XFieldGrid>

      <XNestCard v-for="(c, i) in certs" :key="i" :title="`#${i + 1}`">
        <template #actions>
          <XAddButton v-if="i === 0" :label="t('kit.addRow', { list: t('inbounds.certificate') })" @click="certs.push(defaultCert())" />
          <XInlineButton v-if="certs.length > 1" icon="minus" :label="t('kit.removeRow', { n: i + 1 })" @click="certs.splice(i, 1)" />
        </template>
        <XFieldGrid :cols="1">
          <XField v-slot="a" :label="t('inbounds.certificate')">
            <XSegmented v-bind="a" v-model="c.useFile" :options="certSources" :label="t('inbounds.certificate')" block />
          </XField>
          <XFieldGrid v-if="c.useFile">
            <XField v-slot="a" :label="t('inbounds.publicKey')">
              <XInput v-bind="a" v-model.trim="c.certificateFile" mono />
            </XField>
            <XField v-slot="a" :label="t('inbounds.privateKey')">
              <XInput v-bind="a" v-model.trim="c.keyFile" mono />
            </XField>
          </XFieldGrid>
          <XFieldGrid v-else>
            <XField v-slot="a" :label="t('inbounds.publicKey')">
              <XTextarea v-bind="a" v-model="c.certificate" mono :rows="4" />
            </XField>
            <XField v-slot="a" :label="t('inbounds.privateKey')">
              <XTextarea v-bind="a" v-model="c.key" mono :rows="4" />
            </XField>
          </XFieldGrid>
          <XButtonRow v-if="c.useFile">
            <XButton size="sm" @click="panelCert(c)">{{ t('inbounds.setDefaultCert') }}</XButton>
          </XButtonRow>
          <XFieldGrid :cols="3">
            <XSwitch v-model="c.oneTimeLoading" boxed label="One Time Loading" />
            <XField v-slot="a" label="Usage Option">
              <XSelect v-bind="a" v-model="c.usage" :options="list(certUsages)" />
            </XField>
            <XSwitch v-if="c.usage === 'issue'" v-model="c.buildChain" boxed label="Build Chain" />
          </XFieldGrid>
        </XFieldGrid>
      </XNestCard>

      <XFieldGrid>
        <XField v-slot="a" label="ECH key">
          <XInput v-bind="a" v-model="tls.echServerKeys" mono />
        </XField>
        <XField v-slot="a" label="ECH config">
          <XInput v-bind="a" v-model="tls.settings.echConfigList" mono />
        </XField>
      </XFieldGrid>
      <XButtonRow>
        <XButton size="sm" variant="primary" :loading="busy === 'ech'" @click="getEch">Get New ECH Cert</XButton>
        <XButton size="sm" variant="danger" @click="clearEch">Clear</XButton>
      </XButtonRow>
    </template>

    <!-- Reality -->
    <template v-if="sec === 'reality'">
      <XFieldGrid :cols="3">
        <XSwitch v-model="reality.show" boxed label="Show" />
      </XFieldGrid>
      <XFieldGrid>
        <XField v-slot="a" label="Xver">
          <XNumber v-bind="a" v-model="reality.xver" block :min="0" />
        </XField>
        <XField v-slot="a" label="uTLS">
          <XSelect v-bind="a" v-model="reality.settings.fingerprint" :options="list(fingerprints)" />
        </XField>
        <XField v-slot="a" label="Target" :action="t('inbounds.reset')" @action="newTarget">
          <XInput v-bind="a" v-model.trim="reality.target" mono />
        </XField>
        <XField v-slot="a" label="SNI" :action="t('inbounds.reset')" @action="newTarget">
          <XInput v-bind="a" v-model.trim="serverNames" mono />
        </XField>
      </XFieldGrid>
      <XFieldGrid :cols="3">
        <XField v-slot="a" label="Max Time Diff (ms)">
          <XNumber v-bind="a" v-model="reality.maxTimediff" block :min="0" />
        </XField>
        <XField v-slot="a" label="Min Client Ver">
          <XInput v-bind="a" v-model.trim="reality.minClientVer" placeholder="25.9.11" />
        </XField>
        <XField v-slot="a" label="Max Client Ver">
          <XInput v-bind="a" v-model.trim="reality.maxClientVer" placeholder="25.9.11" />
        </XField>
      </XFieldGrid>
      <XFieldGrid top>
        <XField v-slot="a" label="Short IDs" :action="t('inbounds.reset')" @action="reality.shortIds = randomShortIds()">
          <XTextarea v-bind="a" v-model="shortIds" mono :rows="4" />
        </XField>
        <XField v-slot="a" label="SpiderX">
          <XInput v-bind="a" v-model.trim="reality.settings.spiderX" mono />
        </XField>
      </XFieldGrid>
      <XFieldGrid>
        <XField v-slot="a" :label="t('inbounds.publicKey')">
          <XInput v-bind="a" v-model.trim="reality.settings.publicKey" mono />
        </XField>
        <XField v-slot="a" :label="t('inbounds.privateKey')">
          <XInput v-bind="a" v-model.trim="reality.privateKey" mono />
        </XField>
      </XFieldGrid>
      <XButtonRow>
        <XButton size="sm" variant="primary" :loading="busy === 'x25519'" @click="getX25519">Get New Cert</XButton>
        <XButton size="sm" variant="danger" @click="clearX25519">Clear</XButton>
      </XButtonRow>
      <XFieldGrid>
        <XField v-slot="a" label="mldsa65 Seed">
          <XInput v-bind="a" v-model.trim="reality.mldsa65Seed" mono />
        </XField>
        <XField v-slot="a" label="mldsa65 Verify">
          <XInput v-bind="a" v-model.trim="reality.settings.mldsa65Verify" mono />
        </XField>
      </XFieldGrid>
      <XButtonRow>
        <XButton size="sm" variant="primary" :loading="busy === 'mldsa65'" @click="getMldsa65">Get New Seed</XButton>
        <XButton size="sm" variant="danger" @click="clearMldsa65">Clear</XButton>
      </XButtonRow>
    </template>
  </template>
</template>
