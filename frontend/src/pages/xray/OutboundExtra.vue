<script setup lang="ts">
import { computed } from 'vue'
import { XField, XNest, XNumber, XSelect, XSwitch, XTagInput, type Option } from '../../kit'
import type { Outbound } from './model'
import { addressPortStrategies, canMux, defaultMux, defaultSockopt, forwardedForHeaders } from './outbound'

/** Socket options and Mux at the foot of the outbound window: a switch each, their fields under it when on. */
const props = defineProps<{ o: Outbound; tags: string[] }>()

const sockoptOn = computed({
  get: () => !!props.o.streamSettings?.sockopt,
  set: (on: boolean) => {
    if (on) (props.o.streamSettings ??= {}).sockopt = defaultSockopt()
    else if (props.o.streamSettings) delete props.o.streamSettings.sockopt
  },
})
const so = computed(() => props.o.streamSettings?.sockopt ?? {})

const muxOn = computed({
  get: () => !!props.o.mux?.enabled,
  set: (on: boolean) => (props.o.mux = on ? { ...defaultMux(), ...(props.o.mux ?? {}), enabled: true } : undefined),
})

const list = (values: readonly string[]): Option[] => values.map((v) => ({ value: v, label: v }))
const dialers = computed<Option[]>(() => [{ value: '', label: 'Empty', plain: true }, ...list(props.tags)])
const headers = [{ options: forwardedForHeaders.map((h) => ({ value: h, label: h })) }]
</script>

<template>
  <div class="of-grid of-grid--2">
    <XSwitch v-model="sockoptOn" boxed label="Sockopts" />
    <XSwitch v-if="canMux(o)" v-model="muxOn" boxed label="Mux" />
  </div>

  <XNest v-if="sockoptOn" bare>
    <div class="of-stack">
      <div class="of-grid of-grid--3">
        <XField v-slot="a" label="Dialer Proxy"><XSelect v-bind="a" :model-value="so.dialerProxy ?? ''" @update:model-value="so.dialerProxy = String($event)" :options="dialers" mono /></XField>
        <XField v-slot="a" label="Address Port Strategy"><XSelect v-bind="a" v-model="so.addressPortStrategy" :options="list(addressPortStrategies)" mono /></XField>
        <XField v-slot="a" label="Keep Alive Interval"><XNumber v-bind="a" v-model="so.tcpKeepAliveInterval" :min="0" /></XField>
        <XSwitch v-model="so.tcpFastOpen" boxed label="TCP Fast Open" />
        <XSwitch v-model="so.tcpMptcp" boxed label="Multipath TCP" />
        <XSwitch v-model="so.penetrate" boxed label="Penetrate" />
      </div>
      <XField v-slot="a" label="Trusted X-Forwarded-For">
        <XTagInput v-bind="a" v-model="so.trustedXForwardedFor" :groups="headers" mono />
      </XField>
    </div>
  </XNest>

  <XNest v-if="muxOn && canMux(o)" bare>
    <div class="of-grid of-grid--3">
      <XField v-slot="a" label="Concurrency"><XNumber v-bind="a" v-model="o.mux.concurrency" :min="-1" :max="1024" /></XField>
      <XField v-slot="a" label="xudp Concurrency"><XNumber v-bind="a" v-model="o.mux.xudpConcurrency" :min="-1" :max="1024" /></XField>
      <XField v-slot="a" label="xudp UDP 443"><XSelect v-bind="a" v-model="o.mux.xudpProxyUDP443" :options="list(['reject', 'allow', 'skip'])" mono /></XField>
    </div>
  </XNest>
</template>
