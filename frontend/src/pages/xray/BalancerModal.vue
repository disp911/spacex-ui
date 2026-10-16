<script setup lang="ts">
import { computed, reactive, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { XButton, XField, XInput, XModal, XSelect, XTagInput, type Option, type PresetGroup } from '../../kit'
import { useXrayForm } from './form'
import { listOf, strategies, strategyOf, type Balancer } from './model'

/**
 * A balancer to create, or the one at index to edit: its tag, how it picks an
 * outbound, the outbounds it picks from and where traffic goes when none is up.
 * An empty or taken tag and an empty selector list are marked from the start.
 */
const props = defineProps<{ balancer: Balancer | null; index: number }>()
const open = defineModel<boolean>('open', { default: false })
const emit = defineEmits<{ save: [balancer: Balancer] }>()
const { t } = useI18n()
const f = useXrayForm()

const form = reactive({ tag: '', strategy: 'random', selector: [] as string[], fallbackTag: '' })

watch(open, (on) => {
  if (!on) return
  const b = props.balancer ?? {}
  Object.assign(form, {
    tag: String(b.tag ?? ''),
    strategy: strategyOf(b),
    selector: listOf(b.selector),
    fallbackTag: String(b.fallbackTag ?? ''),
  })
})

const tags = (list: unknown) => [...new Set((Array.isArray(list) ? list : []).map((x) => x?.tag).filter((tag): tag is string => !!tag))]

const taken = computed(() => {
  const others = ((f.draft.config.routing?.balancers ?? []) as Balancer[]).filter((_, i) => i !== props.index)
  return tags(others).includes(form.tag.trim())
})
const tagWrong = computed(() => !form.tag.trim() || taken.value)
const selectorEmpty = computed(() => !form.selector.length)

/** The balancer as Xray takes it; a strategy kept as it was keeps its settings, fields the form does not know stay. */
const result = computed<Balancer>(() => {
  const { tag: _t, strategy, selector: _s, fallbackTag: _f, ...kept } = (props.balancer ?? {}) as Balancer
  const out: Balancer = { ...kept, tag: form.tag.trim(), selector: form.selector }
  if (form.fallbackTag) out.fallbackTag = form.fallbackTag
  if (strategy && strategyOf({ strategy }) === form.strategy) out.strategy = strategy
  else if (form.strategy !== 'random') out.strategy = { type: form.strategy }
  return out
})

function save() {
  if (tagWrong.value || selectorEmpty.value) return
  emit('save', result.value)
  open.value = false
}

const outboundTags = computed(() => tags(f.draft.config.outbounds))
const selectorGroups = computed<PresetGroup[]>(() => [{ options: outboundTags.value.map((v) => ({ value: v, label: v })) }])
// A fallback the balancer already names stays on offer even when no outbound carries it any more.
const fallbacks = computed<Option[]>(() => [
  { value: '', label: '—', plain: true },
  ...[...new Set([...outboundTags.value, ...(form.fallbackTag ? [form.fallbackTag] : [])])].map((v) => ({ value: v, label: v })),
])

const title = computed(() => (props.balancer ? t('xray.balEdit') : t('xray.balAdd')))
</script>

<template>
  <XModal v-model:open="open" :title="title" cancel>
    <div class="balancer">
      <XField v-slot="a" :label="t('xray.outTag')">
        <XInput v-bind="a" v-model="form.tag" mono :invalid="tagWrong" :placeholder="t('xray.balTagHint')" />
      </XField>
      <XField v-slot="a" :label="t('xray.balStrategy')">
        <XSelect v-bind="a" v-model="form.strategy" :options="strategies" />
      </XField>
      <XField v-slot="a" :label="t('xray.balSelectors')">
        <XTagInput v-bind="a" v-model="form.selector" :groups="selectorGroups" mono :invalid="selectorEmpty" />
      </XField>
      <XField v-slot="a" label="Fallback">
        <XSelect v-bind="a" v-model="form.fallbackTag" :options="fallbacks" mono />
      </XField>
    </div>

    <template #footer>
      <XButton variant="primary" :disabled="tagWrong || selectorEmpty" @click="save">{{ balancer ? t('xray.save') : t('xray.balAdd') }}</XButton>
    </template>
  </XModal>
</template>

<style scoped>
.balancer {
  display: flex;
  flex-direction: column;
  gap: 12px;
  min-width: 0;
}
</style>
