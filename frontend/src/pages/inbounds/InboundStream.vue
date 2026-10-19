<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import {
  XDivider, XField, XFieldGrid, XInput, XInputList, XNest, XNumber, XPairList, XPartHead, XSelect, XSwitch, type Option,
} from '../../kit'
import { canTls, defaultMasquerade, defaultTcpRequest, defaultTcpResponse, network, networks, setNetwork, xhttpModes, type InboundForm } from './inbound'
import InboundMasks from './InboundMasks.vue'
import InboundSockopt from './InboundSockopt.vue'
import { headerRows } from './rows'

/**
 * The transport part of window 11, as the old panel's stream forms had it:
 * the transport and its settings (TCP with HTTP masking, mKCP, WebSocket,
 * gRPC, HTTPUpgrade, XHTTP, Hysteria's own), then Sockopt and the masks.
 */
const props = defineProps<{ f: InboundForm }>()
const { t } = useI18n()

type Json = Record<string, any>
const st = computed<Json>(() => props.f.streamSettings ?? {})
const net = computed(() => network(props.f))
/** The settings of the transport in use, e.g. xhttpSettings. */
const ns = computed<Json>(() => st.value[`${net.value}Settings`] ?? {})
const hysteria = computed(() => props.f.protocol === 'hysteria')

const list = (values: readonly string[]): Option[] => values.map((v) => ({ value: v, label: v }))
const netModel = computed({ get: () => net.value, set: (v: string) => setNetwork(props.f, v) })

// --- TCP ---------------------------------------------------------------------------------

const httpMask = computed({
  get: () => ns.value.header?.type === 'http',
  set: (on: boolean) => {
    const h = (ns.value.header ??= {})
    h.type = on ? 'http' : 'none'
    if (on) {
      h.request = { ...defaultTcpRequest(), ...h.request }
      h.response = { ...defaultTcpResponse(), ...h.response }
    }
  },
})
const req = computed<Json>(() => ns.value.header?.request ?? {})
const res = computed<Json>(() => ns.value.header?.response ?? {})
const reqHeaders = headerRows(() => (httpMask.value ? req.value : undefined), 'headers', true)
const resHeaders = headerRows(() => (httpMask.value ? res.value : undefined), 'headers', true)

// --- WebSocket, HTTPUpgrade, XHTTP ---------------------------------------------------------

const headers = headerRows(() => (['ws', 'httpupgrade', 'xhttp'].includes(net.value) ? ns.value : undefined), 'headers')

const placements = (dflt: string, values: string[]): Option[] => [{ value: '', label: `Default (${dflt})` }, ...list(values)]
const paddingPlacements = placements('queryInHeader', ['queryInHeader', 'header', 'cookie', 'query'])
const paddingMethods = placements('repeat-x', ['repeat-x', 'tokenish'])
const uplinkMethods: Option[] = [
  { value: '', label: 'Default (POST)' },
  { value: 'POST', label: 'POST' },
  { value: 'PUT', label: 'PUT' },
  { value: 'GET', label: 'GET (packet-up only)' },
]
const idPlacements = placements('path', ['path', 'header', 'cookie', 'query'])
const dataPlacements = placements('body', ['body', 'header', 'cookie', 'query'])

// --- Hysteria ------------------------------------------------------------------------------

const masqueradeOn = computed({
  get: () => !!ns.value.masquerade,
  set: (on: boolean) => {
    if (on) ns.value.masquerade = defaultMasquerade()
    else delete ns.value.masquerade
  },
})
const mq = computed<Json>(() => ns.value.masquerade ?? {})
const mqHeaders = headerRows(() => ns.value.masquerade, 'headers')
const masqueradeTypes: Option[] = [
  { value: 'file', label: 'File' },
  { value: 'proxy', label: 'Proxy' },
  { value: 'string', label: 'String' },
]
</script>

<template>
  <XPartHead :title="t('inbounds.transport')">
    <XSelect v-if="!hysteria" v-model="netModel" :options="networks" :aria-label="t('inbounds.transport')" />
  </XPartHead>

  <!-- TCP (RAW) -->
  <template v-if="net === 'tcp'">
    <XFieldGrid>
      <XSwitch v-if="canTls(f)" v-model="ns.acceptProxyProtocol" boxed label="Proxy Protocol" />
      <XSwitch v-model="httpMask" boxed :label="`HTTP ${t('inbounds.camouflage')}`" />
    </XFieldGrid>
    <XNest v-if="httpMask" bare>
      <XFieldGrid :cols="1">
        <XDivider :label="t('inbounds.request')" />
        <XFieldGrid>
          <XField v-slot="a" :label="t('inbounds.version')"><XInput v-bind="a" v-model.trim="req.version" mono /></XField>
          <XField v-slot="a" :label="t('inbounds.method')"><XInput v-bind="a" v-model.trim="req.method" mono /></XField>
        </XFieldGrid>
        <XInputList v-model="req.path" :title="t('inbounds.path')" mono />
        <XPairList v-model="reqHeaders" :title="t('inbounds.requestHeader')" :name-label="t('inbounds.headerName')" :value-label="t('inbounds.headerValue')" mono />
        <XDivider :label="t('inbounds.response')" />
        <XFieldGrid :cols="3">
          <XField v-slot="a" :label="t('inbounds.version')"><XInput v-bind="a" v-model.trim="res.version" mono /></XField>
          <XField v-slot="a" :label="t('inbounds.status')"><XInput v-bind="a" v-model.trim="res.status" mono /></XField>
          <XField v-slot="a" :label="t('inbounds.statusDescription')"><XInput v-bind="a" v-model.trim="res.reason" /></XField>
        </XFieldGrid>
        <XPairList v-model="resHeaders" :title="t('inbounds.responseHeader')" :name-label="t('inbounds.headerName')" :value-label="t('inbounds.headerValue')" mono />
      </XFieldGrid>
    </XNest>
  </template>

  <!-- mKCP -->
  <XFieldGrid v-if="net === 'kcp'" :cols="3" :phone="2">
    <XField v-slot="a" label="MTU"><XNumber v-bind="a" v-model="ns.mtu" block :min="576" :max="1460" /></XField>
    <XField v-slot="a" label="TTI (ms)"><XNumber v-bind="a" v-model="ns.tti" block :min="10" :max="100" /></XField>
    <XField v-slot="a" label="Uplink (MB/s)"><XNumber v-bind="a" v-model="ns.uplinkCapacity" block :min="0" /></XField>
    <XField v-slot="a" label="Downlink (MB/s)"><XNumber v-bind="a" v-model="ns.downlinkCapacity" block :min="0" /></XField>
    <XField v-slot="a" label="CWND Multiplier"><XNumber v-bind="a" v-model="ns.cwndMultiplier" block :min="1" /></XField>
    <XField v-slot="a" label="Max Sending Window"><XNumber v-bind="a" v-model="ns.maxSendingWindow" block :min="0" /></XField>
  </XFieldGrid>

  <!-- WebSocket and HTTPUpgrade -->
  <template v-if="net === 'ws' || net === 'httpupgrade'">
    <XSwitch v-model="ns.acceptProxyProtocol" boxed label="Proxy Protocol" />
    <XFieldGrid>
      <XField v-slot="a" :label="t('inbounds.host')"><XInput v-bind="a" v-model.trim="ns.host" mono /></XField>
      <XField v-slot="a" :label="t('inbounds.path')"><XInput v-bind="a" v-model.trim="ns.path" mono /></XField>
      <XField v-if="net === 'ws'" v-slot="a" label="Heartbeat Period"><XNumber v-bind="a" v-model="ns.heartbeatPeriod" block :min="0" /></XField>
    </XFieldGrid>
    <XPairList v-model="headers" :title="t('inbounds.requestHeader')" :name-label="t('inbounds.headerName')" :value-label="t('inbounds.headerValue')" mono />
  </template>

  <!-- gRPC -->
  <template v-if="net === 'grpc'">
    <XFieldGrid>
      <XField v-slot="a" label="Service Name"><XInput v-bind="a" v-model.trim="ns.serviceName" mono /></XField>
      <XField v-slot="a" label="Authority"><XInput v-bind="a" v-model.trim="ns.authority" mono /></XField>
    </XFieldGrid>
    <XSwitch v-model="ns.multiMode" boxed label="Multi Mode" />
  </template>

  <!-- XHTTP -->
  <template v-if="net === 'xhttp'">
    <XFieldGrid>
      <XField v-slot="a" :label="t('inbounds.host')"><XInput v-bind="a" v-model.trim="ns.host" mono /></XField>
      <XField v-slot="a" :label="t('inbounds.path')"><XInput v-bind="a" v-model.trim="ns.path" mono /></XField>
    </XFieldGrid>
    <XPairList v-model="headers" :title="t('inbounds.requestHeader')" :name-label="t('inbounds.headerName')" :value-label="t('inbounds.headerValue')" mono />
    <XFieldGrid>
      <XField v-slot="a" label="Mode"><XSelect v-bind="a" v-model="ns.mode" :options="list(xhttpModes)" mono /></XField>
    </XFieldGrid>
    <XFieldGrid v-if="ns.mode === 'packet-up' || ns.mode === 'stream-up'">
      <template v-if="ns.mode === 'packet-up'">
        <XField v-slot="a" label="Max Buffered Upload"><XNumber v-bind="a" v-model="ns.scMaxBufferedPosts" block /></XField>
        <XField v-slot="a" label="Max Upload Size (Byte)"><XInput v-bind="a" v-model.trim="ns.scMaxEachPostBytes" mono /></XField>
      </template>
      <XField v-else v-slot="a" label="Stream-Up Server"><XInput v-bind="a" v-model.trim="ns.scStreamUpServerSecs" mono /></XField>
    </XFieldGrid>
    <XFieldGrid>
      <XField v-slot="a" label="Padding Bytes"><XInput v-bind="a" v-model.trim="ns.xPaddingBytes" mono /></XField>
    </XFieldGrid>
    <XSwitch v-model="ns.xPaddingObfsMode" boxed label="Padding Obfs Mode" />
    <XFieldGrid v-if="ns.xPaddingObfsMode">
      <XField v-slot="a" label="Padding Key"><XInput v-bind="a" v-model.trim="ns.xPaddingKey" mono placeholder="x_padding" /></XField>
      <XField v-slot="a" label="Padding Header"><XInput v-bind="a" v-model.trim="ns.xPaddingHeader" mono placeholder="X-Padding" /></XField>
      <XField v-slot="a" label="Padding Placement"><XSelect v-bind="a" v-model="ns.xPaddingPlacement" :options="paddingPlacements" mono /></XField>
      <XField v-slot="a" label="Padding Method"><XSelect v-bind="a" v-model="ns.xPaddingMethod" :options="paddingMethods" mono /></XField>
    </XFieldGrid>
    <XFieldGrid>
      <XField v-slot="a" label="Uplink HTTP Method"><XSelect v-bind="a" v-model="ns.uplinkHTTPMethod" :options="uplinkMethods" mono /></XField>
    </XFieldGrid>
    <XFieldGrid>
      <XField v-slot="a" label="Session ID Placement"><XSelect v-bind="a" v-model="ns.sessionIDPlacement" :options="idPlacements" mono /></XField>
      <XField v-if="ns.sessionIDPlacement && ns.sessionIDPlacement !== 'path'" v-slot="a" label="Session ID Key">
        <XInput v-bind="a" v-model.trim="ns.sessionIDKey" mono placeholder="x_session" />
      </XField>
      <XField v-slot="a" label="Sequence Placement"><XSelect v-bind="a" v-model="ns.seqPlacement" :options="idPlacements" mono /></XField>
      <XField v-if="ns.seqPlacement && ns.seqPlacement !== 'path'" v-slot="a" label="Sequence Key">
        <XInput v-bind="a" v-model.trim="ns.seqKey" mono placeholder="x_seq" />
      </XField>
    </XFieldGrid>
    <XFieldGrid v-if="ns.mode === 'packet-up'" :cols="3">
      <XField v-slot="a" label="Uplink Data Placement"><XSelect v-bind="a" v-model="ns.uplinkDataPlacement" :options="dataPlacements" mono /></XField>
      <template v-if="ns.uplinkDataPlacement && ns.uplinkDataPlacement !== 'body'">
        <XField v-slot="a" label="Uplink Data Key"><XInput v-bind="a" v-model.trim="ns.uplinkDataKey" mono placeholder="x_data" /></XField>
        <XField v-slot="a" label="Uplink Chunk Size"><XNumber v-bind="a" v-model="ns.uplinkChunkSize" block :min="0" placeholder="0 (unlimited)" /></XField>
      </template>
    </XFieldGrid>
    <XSwitch v-model="ns.noSSEHeader" boxed label="No SSE Header" />
  </template>

  <!-- Hysteria -->
  <template v-if="net === 'hysteria'">
    <XFieldGrid>
      <XField v-slot="a" label="UDP Idle Timeout"><XNumber v-bind="a" v-model="ns.udpIdleTimeout" block :min="0" /></XField>
    </XFieldGrid>
    <XSwitch v-model="masqueradeOn" boxed label="Masquerade" />
    <XNest v-if="masqueradeOn" bare>
      <XFieldGrid :cols="1">
        <XFieldGrid>
          <XField v-slot="a" label="Type"><XSelect v-bind="a" v-model="mq.type" :options="masqueradeTypes" /></XField>
          <XField v-if="mq.type === 'file'" v-slot="a" label="Dir"><XInput v-bind="a" v-model.trim="mq.dir" mono /></XField>
          <XField v-if="mq.type === 'proxy'" v-slot="a" label="URL"><XInput v-bind="a" v-model.trim="mq.url" mono /></XField>
          <template v-if="mq.type === 'string'">
            <XField v-slot="a" label="Content"><XInput v-bind="a" v-model.trim="mq.content" /></XField>
            <XField v-slot="a" label="Status Code"><XNumber v-bind="a" v-model="mq.statusCode" block /></XField>
          </template>
        </XFieldGrid>
        <XFieldGrid v-if="mq.type === 'proxy'">
          <XSwitch v-model="mq.rewriteHost" boxed label="Rewrite Host" />
          <XSwitch v-model="mq.insecure" boxed label="Insecure" />
        </XFieldGrid>
        <XPairList v-if="mq.type === 'string'" v-model="mqHeaders" :title="t('inbounds.requestHeader')" :name-label="t('inbounds.headerName')" :value-label="t('inbounds.headerValue')" mono />
      </XFieldGrid>
    </XNest>
  </template>

  <InboundSockopt :f="f" />
  <InboundMasks :f="f" />
</template>
