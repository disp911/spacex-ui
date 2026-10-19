<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { XCheckGroup, XField, XFieldGrid, XNest, XSection, XSwitch, XTagInput } from '../../kit'
import { sniffingTargets, type InboundForm } from './inbound'

/**
 * Sniffing in window 11, as the old panel's form had it: on or off, then
 * what to sniff, Metadata Only, Route Only and the addresses left alone.
 * The card starts closed, as on the old panel; Sniffing itself starts off.
 */
const props = defineProps<{ f: InboundForm }>()
const open = defineModel<boolean>('open', { default: false })
const { t } = useI18n()

const targets = sniffingTargets.map((v) => ({ value: v, label: v }))
const list = (key: 'ipsExcluded' | 'domainsExcluded') => (Array.isArray(props.f.sniffing[key]) ? props.f.sniffing[key] : [])
</script>

<template>
  <XSection v-model:open="open" title="Sniffing" card>
    <XSwitch v-model="f.sniffing.enabled" boxed :label="t('inbounds.enabled')" :help="t('inbounds.sniffingHelp')" />
    <XNest v-if="f.sniffing.enabled" bare>
      <XFieldGrid :cols="1">
        <XCheckGroup v-model="f.sniffing.destOverride" :options="targets" mono label="Sniffing" />
        <XFieldGrid>
          <XSwitch v-model="f.sniffing.metadataOnly" boxed label="Metadata Only" />
          <XSwitch v-model="f.sniffing.routeOnly" boxed label="Route Only" />
        </XFieldGrid>
        <XFieldGrid>
          <XField v-slot="a" label="IPs Excluded">
            <XTagInput v-bind="a" :model-value="list('ipsExcluded')" mono placeholder="IP/CIDR/geoip:*/ext:*" @update:model-value="f.sniffing.ipsExcluded = $event" />
          </XField>
          <XField v-slot="a" label="Domains Excluded">
            <XTagInput v-bind="a" :model-value="list('domainsExcluded')" mono placeholder="domain:*/ext:*" @update:model-value="f.sniffing.domainsExcluded = $event" />
          </XField>
        </XFieldGrid>
      </XFieldGrid>
    </XNest>
  </XSection>
</template>
