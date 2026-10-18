<script setup lang="ts">
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import {
  XAddButton, XField, XInlineButton, XInput, XInputList, XListHead, XNest, XNestCard, XNumber, XSection, XSelect, XSwitch, XTagInput, type Option,
} from '../../kit'
import { wireguardKeys } from '../../wireguard'
import type { Outbound } from './model'
import {
  canFlow, canVisionSeed, defaultDnsRule, defaultFinalRule, defaultFragment, defaultNoise, defaultPeer, domainStrategies, flows, hasUser, server, ssMethods, user,
  items, vmessSecurities, wireguardStrategies,
} from './outbound'

/** The part of the outbound window that belongs to its protocol, in a card titled with it. */
const props = defineProps<{ o: Outbound; tags: string[] }>()
const { t } = useI18n()
const open = ref(true)

const s = computed(() => (props.o.settings ??= {}))
const srv = computed(() => server(props.o))
const usr = computed(() => user(props.o))

const list = (values: readonly string[], empty = false): Option[] => [
  ...(empty ? [{ value: '', label: 'Empty', plain: true }] : []),
  ...values.map((v) => ({ value: v, label: v })),
]

/** A comma list kept as an array, e.g. WireGuard addresses. */
const commaList = (key: string, numbers = false) =>
  computed({
    get: () => (Array.isArray(s.value[key]) ? s.value[key].join(',') : String(s.value[key] ?? '')),
    set: (v: string) => {
      const parts = v.split(',').map((x) => x.trim()).filter(Boolean)
      s.value[key] = numbers ? parts.map(Number) : parts
    },
  })

// --- VLESS -------------------------------------------------------------------------

const reverseTag = computed({
  get: () => s.value.reverse?.tag ?? '',
  set: (v: string) => (s.value.reverse = { ...(s.value.reverse ?? {}), tag: v }),
})
const seeds = computed(() => (usr.value.testseed ??= [900, 500, 900, 256]) as number[])

// --- Shadowsocks, SOCKS, HTTP ------------------------------------------------------

/** SOCKS and HTTP keep one optional user: without a name there is none. */
const account = (key: 'user' | 'pass') =>
  computed({
    get: () => srv.value.users?.[0]?.[key] ?? '',
    set: (v: string) => {
      const u = { ...(srv.value.users?.[0] ?? {}), [key]: v }
      srv.value.users = u.user ? [u] : []
    },
  })
const socksUser = account('user')
const socksPass = account('pass')

// --- WireGuard ---------------------------------------------------------------------

const wgAddress = commaList('address')
const wgReserved = commaList('reserved', true)
const publicKey = computed(() => {
  try {
    return s.value.secretKey ? wireguardKeys(s.value.secretKey).publicKey : ''
  } catch {
    return ''
  }
})
const newKey = () => (s.value.secretKey = wireguardKeys().privateKey)

// --- Freedom -----------------------------------------------------------------------

const fragmentOn = computed({
  get: () => !!s.value.fragment && Object.keys(s.value.fragment).length > 0,
  set: (on: boolean) => (s.value.fragment = on ? defaultFragment() : undefined),
})
const noisesOn = computed({
  get: () => (s.value.noises?.length ?? 0) > 0,
  set: (on: boolean) => (s.value.noises = on ? [defaultNoise()] : []),
})
/** A rule's network may come as a list; the field shows it as text. */
const networkText = (r: Record<string, any>) => (Array.isArray(r.network) ? r.network.join(',') : (r.network ?? ''))

// --- Blackhole ---------------------------------------------------------------------

const responseType = computed({
  get: () => s.value.response?.type ?? '',
  set: (v: string) => (s.value.response = v ? { type: v } : undefined),
})

const remove = (list: unknown[], i: number) => list.splice(i, 1)
</script>

<template>
  <XSection v-model:open="open" :title="o.protocol" card mono>
    <!-- VMess, VLESS -->
    <template v-if="o.protocol === 'vmess' || o.protocol === 'vless'">
      <div class="of-grid of-grid--2">
        <XField v-slot="a" label="ID">
          <XInput v-bind="a" v-model="usr.id" mono />
        </XField>
        <XField v-if="o.protocol === 'vmess'" v-slot="a" label="Security">
          <XSelect v-bind="a" v-model="usr.security" :options="list(vmessSecurities)" mono />
        </XField>
        <XField v-else v-slot="a" label="encryption">
          <XInput v-bind="a" v-model="usr.encryption" mono />
        </XField>
      </div>
      <div v-if="o.protocol === 'vless'" class="of-grid of-grid--2">
        <XField v-slot="a" :label="t('xray.reverseTag')" :help="t('xray.reverseTagDesc')">
          <XInput v-bind="a" v-model="reverseTag" mono />
        </XField>
        <XField v-if="canFlow(o)" v-slot="a" label="Flow">
          <XSelect v-bind="a" :model-value="usr.flow ?? ''" :options="list(flows, true)" @update:model-value="usr.flow = String($event)" mono />
        </XField>
      </div>
      <XNest v-if="canVisionSeed(o)" bare>
        <div class="of-stack">
          <div class="of-grid of-grid--3">
            <XField v-slot="a" label="Vision Pre-Connect">
              <XNumber v-bind="a" block v-model="usr.testpre" :min="0" :max="10" placeholder="0" />
            </XField>
          </div>
          <XListHead title="Vision Seed" />
          <div class="of-grid of-grid--seed">
            <XField v-for="(_, i) in 4" :key="i" v-slot="a" :label="`[${i}]`">
              <XNumber v-bind="a" block v-model="seeds[i]" :min="0" :max="9999" />
            </XField>
          </div>
        </div>
      </XNest>
    </template>

    <!-- Trojan, Shadowsocks -->
    <div v-if="o.protocol === 'trojan' || o.protocol === 'shadowsocks'" class="of-grid of-grid--2">
      <XField v-slot="a" :label="t('xray.password')">
        <XInput v-bind="a" v-model="srv.password" mono />
      </XField>
      <XField v-if="o.protocol === 'shadowsocks'" v-slot="a" :label="t('xray.ssMethod')">
        <XSelect v-bind="a" v-model="srv.method" :options="list(ssMethods)" mono />
      </XField>
    </div>
    <div v-if="o.protocol === 'shadowsocks'" class="of-grid of-grid--2">
      <XSwitch v-model="srv.uot" boxed label="UDP over TCP" />
      <XField v-slot="a" label="UoTVersion">
        <XNumber v-bind="a" block v-model="srv.UoTVersion" :min="1" :max="2" />
      </XField>
    </div>

    <!-- SOCKS, HTTP -->
    <div v-if="hasUser(o)" class="of-grid of-grid--2">
      <XField v-slot="a" :label="t('xray.username')">
        <XInput v-bind="a" v-model="socksUser" mono />
      </XField>
      <XField v-slot="a" :label="t('xray.password')">
        <XInput v-bind="a" v-model="socksPass" mono />
      </XField>
    </div>

    <!-- Hysteria -->
    <div v-if="o.protocol === 'hysteria'" class="of-grid of-grid--3">
      <XField v-slot="a" label="Version">
        <XNumber v-bind="a" block v-model="s.version" disabled />
      </XField>
    </div>

    <!-- DNS -->
    <template v-if="o.protocol === 'dns'">
      <div class="of-grid of-grid--2">
        <XField v-slot="a" :label="t('xray.ruleNetwork')">
          <XSelect v-bind="a" v-model="s.network" :options="list(['udp', 'tcp'])" mono />
        </XField>
      </div>
      <XListHead title="Rules" add @add="(s.rules ??= []).push(defaultDnsRule())" />
      <XNestCard v-for="(r, i) in items(s.rules)" :key="i" :title="String(i + 1)">
        <template #actions>
          <XInlineButton icon="delete" danger :label="t('xray.removeItem', { n: i + 1 })" @click="remove(s.rules, i)" />
        </template>
        <div class="of-grid of-grid--3">
          <XField v-slot="a" label="Action">
            <XSelect v-bind="a" v-model="r.action" :options="list(['direct', 'drop', 'return', 'hijack'])" mono />
          </XField>
          <XField v-slot="a" label="QType" help="Single qType (e.g. 28) or list/range (e.g. 1,3,23-24)">
            <XInput v-bind="a" v-model="r.qType" mono placeholder="1,3,23-24" />
          </XField>
          <XField v-slot="a" label="Domain" help="Comma-separated domain rules, e.g. domain:example.com,full:example.com">
            <XInput v-bind="a" v-model="r.domain" mono placeholder="domain:example.com" />
          </XField>
        </div>
      </XNestCard>
    </template>

    <!-- Blackhole -->
    <div v-if="o.protocol === 'blackhole'" class="of-grid of-grid--2">
      <XField v-slot="a" label="Response Type">
        <XSelect v-bind="a" v-model="responseType" :options="list(['none', 'http'], true)" mono />
      </XField>
    </div>

    <!-- Freedom -->
    <template v-if="o.protocol === 'freedom'">
      <div class="of-grid of-grid--2">
        <XField v-slot="a" label="Strategy">
          <XSelect v-bind="a" :model-value="s.domainStrategy ?? ''" :options="list(domainStrategies, true)" @update:model-value="s.domainStrategy = String($event)" mono />
        </XField>
        <XField v-slot="a" label="Redirect">
          <XInput v-bind="a" v-model="s.redirect" mono />
        </XField>
      </div>
      <XListHead title="Final Rules" add @add="(s.finalRules ??= []).push(defaultFinalRule())" />
      <XNestCard v-for="(r, i) in items(s.finalRules)" :key="i" :title="String(i + 1)">
        <template #actions>
          <XInlineButton icon="delete" danger :label="t('xray.removeItem', { n: i + 1 })" @click="remove(s.finalRules, i)" />
        </template>
        <div class="of-grid of-grid--3 of-grid--phone-2">
          <XField v-slot="a" label="Action">
            <XSelect v-bind="a" v-model="r.action" :options="list(['allow', 'block'])" mono />
          </XField>
          <XField v-slot="a" label="Network" help="Optional: tcp, udp, tcp,udp">
            <XInput v-bind="a" :model-value="networkText(r)" mono placeholder="tcp,udp" @update:model-value="r.network = $event" />
          </XField>
          <XField v-slot="a" label="Port" help="Optional: same format as routing port, e.g. 53,443,1000-2000">
            <XInput v-bind="a" v-model="r.port" mono placeholder="53,443" />
          </XField>
          <XField v-slot="a" label="IP" help="Optional IP/CIDR/geoip list, e.g. geoip:cn, 10.0.0.0/8">
            <XTagInput v-bind="a" v-model="r.ip" mono placeholder="geoip:cn,10.0.0.0/8" />
          </XField>
          <XField v-if="r.action === 'block'" v-slot="a" label="Block Delay">
            <XInput v-bind="a" v-model="r.blockDelay" mono placeholder="30-90" />
          </XField>
        </div>
      </XNestCard>

      <XSwitch v-model="fragmentOn" boxed label="Fragment" />
      <XNest v-if="fragmentOn" bare>
        <div class="of-grid of-grid--4">
          <XField v-slot="a" label="Packets">
            <XSelect v-bind="a" v-model="s.fragment.packets" :options="list(['1-3', 'tlshello'])" mono />
          </XField>
          <XField v-slot="a" label="Length">
            <XInput v-bind="a" v-model="s.fragment.length" mono />
          </XField>
          <XField v-slot="a" label="Interval">
            <XInput v-bind="a" v-model="s.fragment.interval" mono />
          </XField>
          <XField v-slot="a" label="Max Split">
            <XInput v-bind="a" v-model="s.fragment.maxSplit" mono />
          </XField>
        </div>
      </XNest>

      <XSwitch v-model="noisesOn" boxed label="Noises">
        <template v-if="noisesOn" #actions>
          <XAddButton :label="t('kit.addRow', { list: 'Noises' })" @click="s.noises.push(defaultNoise())" />
        </template>
      </XSwitch>
      <XNest v-if="noisesOn" bare>
        <div class="of-stack">
          <XNestCard v-for="(n, i) in items(s.noises)" :key="i" :title="String(i + 1)">
            <template #actions>
              <XInlineButton v-if="s.noises.length > 1" icon="delete" danger :label="t('xray.removeItem', { n: i + 1 })" @click="remove(s.noises, i)" />
            </template>
            <div class="of-grid of-grid--4">
              <XField v-slot="a" label="Type">
                <XSelect v-bind="a" v-model="n.type" :options="list(['rand', 'base64', 'str', 'hex'])" mono />
              </XField>
              <XField v-slot="a" label="Packet">
                <XInput v-bind="a" v-model="n.packet" mono />
              </XField>
              <XField v-slot="a" label="Delay">
                <XInput v-bind="a" v-model="n.delay" mono />
              </XField>
              <XField v-slot="a" label="Apply To">
                <XSelect v-bind="a" v-model="n.applyTo" :options="list(['ip', 'ipv4', 'ipv6'])" mono />
              </XField>
            </div>
          </XNestCard>
        </div>
      </XNest>
    </template>

    <!-- WireGuard -->
    <template v-if="o.protocol === 'wireguard'">
      <div class="of-grid of-grid--2">
        <XField v-slot="a" :label="t('xray.privateKey')">
          <XInput v-bind="a" v-model="s.secretKey" mono>
            <template #actions><XInlineButton icon="refresh" :label="t('xray.newKey')" @click="newKey" /></template>
          </XInput>
        </XField>
        <XField v-slot="a" :label="t('xray.publicKey')">
          <XInput v-bind="a" :model-value="publicKey" mono readonly />
        </XField>
      </div>
      <div class="of-grid of-grid--3">
        <XField v-slot="a" :label="t('xray.domainStrategy')">
          <XSelect v-bind="a" :model-value="s.domainStrategy ?? ''" :options="list(wireguardStrategies, true)" @update:model-value="s.domainStrategy = String($event)" mono />
        </XField>
        <XField v-slot="a" label="MTU">
          <XNumber v-bind="a" block v-model="s.mtu" :min="0" />
        </XField>
        <XField v-slot="a" label="Workers">
          <XNumber v-bind="a" block v-model="s.workers" :min="0" />
        </XField>
      </div>
      <div class="of-grid of-grid--2">
        <XField v-slot="a" :label="t('xray.outAddress')" :help="t('xray.useComma')">
          <XInput v-bind="a" v-model="wgAddress" mono />
        </XField>
        <XField v-slot="a" label="Reserved" :help="t('xray.useComma')">
          <XInput v-bind="a" v-model="wgReserved" mono />
        </XField>
        <XSwitch v-model="s.noKernelTun" boxed label="No Kernel Tun" />
      </div>
      <XListHead title="Peers" add @add="(s.peers ??= []).push(defaultPeer())" />
      <XNestCard v-for="(p, i) in items(s.peers)" :key="i" :title="t('xray.peer', { n: i + 1 })">
        <template #actions>
          <XInlineButton v-if="s.peers.length > 1" icon="delete" danger :label="t('xray.removeItem', { n: i + 1 })" @click="remove(s.peers, i)" />
        </template>
        <div class="of-stack">
          <XField v-slot="a" :label="t('xray.endpoint')">
            <XInput v-bind="a" v-model="p.endpoint" mono />
          </XField>
          <div class="of-grid of-grid--2">
            <XField v-slot="a" :label="t('xray.publicKey')">
              <XInput v-bind="a" v-model="p.publicKey" mono />
            </XField>
            <XField v-slot="a" :label="t('xray.psk')">
              <XInput v-bind="a" v-model="p.preSharedKey" mono />
            </XField>
          </div>
          <XInputList v-model="p.allowedIPs" title="Allowed IPs" mono />
          <div class="of-grid of-grid--3">
            <XField v-slot="a" label="Keep Alive">
              <XNumber v-bind="a" block v-model="p.keepAlive" :min="0" />
            </XField>
          </div>
        </div>
      </XNestCard>
    </template>
  </XSection>
</template>
