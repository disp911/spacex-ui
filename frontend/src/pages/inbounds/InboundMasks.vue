<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { XField, XFieldGrid, XInlineButton, XInput, XListHead, XNestCard, XNumber, XSelect, XTagInput, type Option } from '../../kit'
import OutboundPackets from '../xray/OutboundPackets.vue'
import { canUdpMask, defaultTcpMask, defaultUdpMask, network, type InboundForm } from './inbound'

/**
 * TCP and UDP masks of an inbound's transport, as the old panel's form had
 * them on the server side: TCP masks for raw TCP, WebSocket, gRPC,
 * HTTPUpgrade and XHTTP; UDP masks for mKCP and Hysteria.
 */
const props = defineProps<{ f: InboundForm }>()
const { t } = useI18n()

type Json = Record<string, any>
const st = computed(() => {
  props.f.streamSettings ??= {}
  // Read back through the reactive object: ??= hands back the plain new value.
  return props.f.streamSettings
})
const fm = computed<Json>(() => {
  st.value.finalmask ??= {}
  // Read back through the reactive object: ??= hands back the plain new value.
  return st.value.finalmask
})
const items = (v: unknown): Json[] => (Array.isArray(v) ? v : [])
const list = (values: string[]): Option[] => values.map((v) => ({ value: v, label: v }))

const tcpMasks = computed(() => ['tcp', 'httpupgrade', 'ws', 'grpc', 'xhttp'].includes(network(props.f)))
const udpMasks = computed(() => canUdpMask(props.f))
const hysteria = computed(() => props.f.protocol === 'hysteria')

const packetItem = () => ({ delay: 0, rand: 0, randRange: '0-255', type: 'array', packet: [] })

const tcpTypes: Option[] = [
  { value: 'fragment', label: 'Fragment' },
  { value: 'header-custom', label: 'Header Custom' },
  { value: 'sudoku', label: 'Sudoku' },
]
const udpTypes = computed<Option[]>(() =>
  hysteria.value
    ? [{ value: 'salamander', label: 'Salamander (Hysteria2)' }]
    : [
        { value: 'mkcp-legacy', label: 'mKCP' },
        { value: 'xdns', label: 'xDNS' },
        { value: 'xicmp', label: 'xICMP' },
        { value: 'header-custom', label: 'Header Custom' },
        { value: 'noise', label: 'Noise' },
      ],
)
const kcpHeaders: Option[] = [
  { value: '', label: 'None' },
  { value: 'dns', label: 'DNS' },
  { value: 'dtls', label: 'DTLS 1.2' },
  { value: 'srtp', label: 'SRTP' },
  { value: 'utp', label: 'uTP' },
  { value: 'wechat', label: 'WeChat Video' },
  { value: 'wireguard', label: 'WireGuard' },
]

function setTcpType(m: Json, type: string) {
  Object.assign(m, defaultTcpMask(type))
}

function setUdpType(m: Json, type: string) {
  Object.assign(m, defaultUdpMask(type))
  // xDNS rides in small packets: mKCP's MTU follows it, as on the old panel.
  if (network(props.f) === 'kcp') st.value.kcpSettings.mtu = type === 'xdns' ? 900 : 1350
}

const addUdpMask = () => (fm.value.udp ??= []).push(defaultUdpMask(hysteria.value ? 'salamander' : 'mkcp-legacy'))

/** Sudoku's custom tables are one text in the field, a list in the JSON. */
const tablesText = (m: Json) => (Array.isArray(m.settings.customTables) ? m.settings.customTables.join(',') : '')
const setTables = (m: Json, v: string) => (m.settings.customTables = v.split(',').map((x) => x.trim()).filter(Boolean))
</script>

<template>
  <!-- TCP masks -->
  <template v-if="tcpMasks">
    <XListHead title="TCP Masks" add @add="(fm.tcp ??= []).push(defaultTcpMask())">
      <XNestCard v-for="(m, i) in items(fm.tcp)" :key="i" :title="String(i + 1)">
        <template #actions>
          <XInlineButton icon="delete" danger :label="t('xray.removeItem', { n: i + 1 })" @click="fm.tcp.splice(i, 1)" />
        </template>
        <XFieldGrid :cols="1">
          <XFieldGrid :cols="2">
            <XField v-slot="a" :label="t('inbounds.type')">
              <XSelect v-bind="a" :model-value="m.type" :options="tcpTypes" @update:model-value="setTcpType(m, String($event))" />
            </XField>
          </XFieldGrid>
          <XFieldGrid v-if="m.type === 'fragment'" :cols="4">
            <XField v-slot="a" label="Packets"><XSelect v-bind="a" v-model="m.settings.packets" :options="list(['tlshello', '1-3', '1-5'])" mono /></XField>
            <XField v-slot="a" label="Length"><XInput v-bind="a" v-model="m.settings.length" mono :placeholder="t('xray.hintExample', { v: '100-200' })" /></XField>
            <XField v-slot="a" label="Delay"><XInput v-bind="a" v-model="m.settings.delay" mono :placeholder="t('xray.hintExample', { v: '10-20' })" /></XField>
            <XField v-slot="a" label="Max Split"><XInput v-bind="a" v-model="m.settings.maxSplit" mono :placeholder="t('xray.hintExample', { v: '3-6' })" /></XField>
          </XFieldGrid>
          <XFieldGrid v-if="m.type === 'sudoku'" :cols="3">
            <XField v-slot="a" label="Password"><XInput v-bind="a" v-model="m.settings.password" mono :placeholder="t('xray.hintObfs')" /></XField>
            <XField v-slot="a" label="ASCII"><XInput v-bind="a" v-model="m.settings.ascii" mono /></XField>
            <XField v-slot="a" label="Custom Table"><XInput v-bind="a" v-model="m.settings.customTable" mono /></XField>
            <XField v-slot="a" label="Custom Tables">
              <XInput v-bind="a" :model-value="tablesText(m)" mono @update:model-value="setTables(m, $event)" />
            </XField>
            <XField v-slot="a" label="Padding Min"><XNumber v-bind="a" block v-model="m.settings.paddingMin" :min="0" /></XField>
            <XField v-slot="a" label="Padding Max"><XNumber v-bind="a" block v-model="m.settings.paddingMax" :min="0" /></XField>
          </XFieldGrid>
          <template v-if="m.type === 'header-custom'">
            <template v-for="side in (['clients', 'servers'] as const)" :key="side">
              <XListHead :title="side === 'clients' ? 'Clients' : 'Servers'" add @add="(m.settings[side] ??= []).push([packetItem()])">
                <XNestCard v-for="(group, gi) in items(m.settings[side])" :key="gi" :title="`${side === 'clients' ? 'Clients' : 'Servers'} Group ${gi + 1}`">
                  <template #actions>
                    <XInlineButton icon="delete" danger :label="t('xray.removeItem', { n: gi + 1 })" @click="m.settings[side].splice(gi, 1)" />
                  </template>
                  <OutboundPackets :items="items(group)" with-delay />
                </XNestCard>
              </XListHead>
            </template>
          </template>
        </XFieldGrid>
      </XNestCard>
    </XListHead>
  </template>

  <!-- UDP masks -->
  <template v-if="udpMasks">
    <XListHead title="UDP Masks" add @add="addUdpMask">
      <XNestCard v-for="(m, i) in items(fm.udp)" :key="i" :title="String(i + 1)">
        <template #actions>
          <XInlineButton icon="delete" danger :label="t('xray.removeItem', { n: i + 1 })" @click="fm.udp.splice(i, 1)" />
        </template>
        <XFieldGrid :cols="1">
          <XFieldGrid :cols="2">
            <XField v-slot="a" :label="t('inbounds.type')">
              <XSelect v-bind="a" :model-value="m.type" :options="udpTypes" @update:model-value="setUdpType(m, String($event))" />
            </XField>
            <XField v-if="m.type === 'salamander'" v-slot="a" label="Password">
              <XInput v-bind="a" v-model="m.settings.password" mono :placeholder="t('xray.hintObfs')" />
            </XField>
            <template v-if="m.type === 'mkcp-legacy'">
              <XField v-slot="a" label="Header"><XSelect v-bind="a" v-model="m.settings.header" :options="kcpHeaders" /></XField>
              <XField v-if="!m.settings.header" v-slot="a" label="Password">
                <XInput v-bind="a" v-model="m.settings.value" mono :placeholder="t('xray.hintMkcp')" />
              </XField>
              <XField v-if="m.settings.header === 'dns'" v-slot="a" label="Domain">
                <XInput v-bind="a" v-model="m.settings.value" mono :placeholder="t('xray.hintExample', { v: 'www.example.com' })" />
              </XField>
            </template>
            <XField v-if="m.type === 'xdns'" v-slot="a" label="Domains" class="x-wide">
              <XTagInput v-bind="a" v-model="m.settings.domains[0].names" mono :placeholder="t('xray.hintExample', { v: 'www.example.com' })" />
            </XField>
            <XField v-if="m.type === 'noise'" v-slot="a" label="Reset"><XNumber v-bind="a" block v-model="m.settings.reset" :min="0" /></XField>
          </XFieldGrid>
          <template v-if="m.type === 'noise'">
            <XListHead title="Noise" add @add="m.settings.noise.push({ rand: '1-8192', randRange: '0-255', type: 'array', packet: [], delay: '10-20' })">
              <OutboundPackets :items="m.settings.noise" with-delay removable />
            </XListHead>
          </template>
          <template v-if="m.type === 'header-custom'">
            <XListHead title="Client" add @add="m.settings.client.push(packetItem())">
              <OutboundPackets :items="m.settings.client" removable />
            </XListHead>
            <XListHead title="Server" add @add="m.settings.server.push(packetItem())">
              <OutboundPackets :items="m.settings.server" removable />
            </XListHead>
          </template>
        </XFieldGrid>
      </XNestCard>
    </XListHead>
  </template>
</template>
