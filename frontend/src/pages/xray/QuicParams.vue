<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { XField, XFieldGrid, XInput, XNumber, XSelect, XSwitch, type Option } from '../../kit'

/**
 * The fields of finalmask.quicParams, for Hysteria and XHTTP over HTTP/3:
 * the same on an inbound and an outbound, each of which turns them on its own way.
 */
defineProps<{ qp: Record<string, any> }>()
const { t } = useI18n()

const congestions: Option[] = [
  { value: 'reno', label: 'Reno' },
  { value: 'bbr', label: 'BBR' },
  { value: 'brutal', label: 'Brutal' },
  { value: 'force-brutal', label: 'Force Brutal' },
]
</script>

<template>
  <XFieldGrid :cols="2">
    <XField v-slot="a" label="Congestion"><XSelect v-bind="a" v-model="qp.congestion" :options="congestions" /></XField>
    <XSwitch v-model="qp.debug" boxed label="Debug" />
    <template v-if="qp.congestion === 'brutal' || qp.congestion === 'force-brutal'">
      <XField v-slot="a" label="Brutal Up"><XInput v-bind="a" v-model="qp.brutalUp" mono placeholder="65537" /></XField>
      <XField v-slot="a" label="Brutal Down"><XInput v-bind="a" v-model="qp.brutalDown" mono placeholder="65537" /></XField>
    </template>
  </XFieldGrid>
  <XFieldGrid :cols="2">
    <XField v-slot="a" label="Max Idle Timeout (s)"><XNumber v-bind="a" block v-model="qp.maxIdleTimeout" :min="4" :max="120" /></XField>
    <XField v-slot="a" label="Keep Alive Period (s)"><XNumber v-bind="a" block v-model="qp.keepAlivePeriod" :min="0" :max="60" /></XField>
    <XSwitch v-model="qp.disablePathMTUDiscovery" boxed label="Disable Path MTU Dis" />
    <XField v-slot="a" label="Max Incoming Streams"><XNumber v-bind="a" block v-model="qp.maxIncomingStreams" :min="8" :placeholder="t('xray.hintDefault', { v: '8' })" /></XField>
    <XField v-slot="a" label="Init Stream Window"><XNumber v-bind="a" block v-model="qp.initStreamReceiveWindow" :min="16384" :placeholder="t('xray.hintDefault', { v: '8388608' })" /></XField>
    <XField v-slot="a" label="Max Stream Window"><XNumber v-bind="a" block v-model="qp.maxStreamReceiveWindow" :min="16384" :placeholder="t('xray.hintDefault', { v: '8388608' })" /></XField>
    <XField v-slot="a" label="Init Conn Window"><XNumber v-bind="a" block v-model="qp.initConnectionReceiveWindow" :min="16384" :placeholder="t('xray.hintDefault', { v: '20971520' })" /></XField>
    <XField v-slot="a" label="Max Conn Window"><XNumber v-bind="a" block v-model="qp.maxConnectionReceiveWindow" :min="16384" :placeholder="t('xray.hintDefault', { v: '20971520' })" /></XField>
  </XFieldGrid>
</template>
