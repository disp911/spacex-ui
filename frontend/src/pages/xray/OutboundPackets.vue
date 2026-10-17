<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { XField, XInlineButton, XInput, XNestCard, XNumber, XSelect, type Option } from '../../kit'

/**
 * Packets of a mask, e.g. Header Custom or Noise: random bytes of a length
 * (array) or bytes given as text, hex or base64, each with an optional delay.
 */
const props = defineProps<{ items: Record<string, any>[]; withDelay?: boolean; removable?: boolean }>()
const { t } = useI18n()

const types: Option[] = [
  { value: 'array', label: 'Array' },
  { value: 'str', label: 'String' },
  { value: 'hex', label: 'Hex' },
  { value: 'base64', label: 'Base64' },
]

function setType(item: Record<string, any>, type: string) {
  item.type = type
  if (type === 'array') {
    item.rand = 0
    item.packet = []
  } else item.packet = ''
}
</script>

<template>
  <XNestCard v-for="(item, i) in props.items" :key="i" :title="String(i + 1)">
    <template v-if="removable" #actions>
      <XInlineButton icon="delete" danger :label="t('xray.removeItem', { n: i + 1 })" @click="props.items.splice(i, 1)" />
    </template>
    <div class="of-grid of-grid--4">
      <XField v-slot="a" label="Type">
        <XSelect v-bind="a" :model-value="item.type" :options="types" @update:model-value="setType(item, String($event))" />
      </XField>
      <!-- Header Custom delays are numbers, Noise delays ranges such as 10-20. -->
      <XField v-if="withDelay" v-slot="a" label="Delay (ms)">
        <XNumber v-if="typeof item.delay === 'number'" v-bind="a" block v-model="item.delay" :min="0" />
        <XInput v-else v-bind="a" v-model="item.delay" mono placeholder="10-20" />
      </XField>
      <template v-if="item.type === 'array'">
        <XField v-slot="a" label="Rand"><XNumber v-bind="a" block v-model="item.rand" :min="0" /></XField>
        <XField v-slot="a" label="Rand Range"><XInput v-bind="a" v-model="item.randRange" mono placeholder="0-255" /></XField>
      </template>
      <XField v-else v-slot="a" label="Packet">
        <XInput v-bind="a" v-model="item.packet" mono placeholder="binary data" />
      </XField>
    </div>
  </XNestCard>
</template>
