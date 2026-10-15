<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { checkJson, XCard, XCodeEditor, XSegmented } from '../../kit'
import { useXrayForm } from './form'
import type { XrayConfig } from './model'

/**
 * The template as JSON, whole or one part of it. Text that parses goes into
 * the template as it is typed; text with a mistake stays in the editor only.
 */
const { t } = useI18n()
const f = useXrayForm()

type Part = 'all' | 'inbounds' | 'outbounds' | 'rules'
const part = ref<Part>('all')
const parts = computed(() => [
  { value: 'all' as const, label: t('xray.partAll') },
  { value: 'inbounds' as const, label: t('xray.partInbounds') },
  { value: 'outbounds' as const, label: t('xray.partOutbounds') },
  { value: 'rules' as const, label: t('xray.partRules') },
])

function read(c: XrayConfig, p: Part): unknown {
  if (p === 'all') return c
  if (p === 'rules') return c.routing?.rules ?? []
  return c[p] ?? []
}

function write(c: XrayConfig, p: Part, value: unknown) {
  if (p === 'all') {
    for (const key of Object.keys(c)) delete c[key]
    Object.assign(c, value)
  } else if (p === 'rules') (c.routing ??= {}).rules = value
  else c[p] = value
}

const same = (a: unknown, b: unknown) => JSON.stringify(a) === JSON.stringify(b)
const shown = computed(() => read(f.draft.config, part.value))
const text = ref(JSON.stringify(shown.value, null, 2))

// Typed text that parses goes into the template.
watch(text, (value) => {
  if (checkJson(value)) return
  const parsed = JSON.parse(value)
  if (!same(parsed, shown.value)) write(f.draft.config, part.value, parsed)
})

// The template changed elsewhere, or another part is chosen: the editor shows it,
// unless the text already says the same in its own layout.
watch([shown, part], ([value], [, before]) => {
  if (part.value === before && !checkJson(text.value) && same(JSON.parse(text.value), value)) return
  text.value = JSON.stringify(value, null, 2)
}, { deep: true })
</script>

<template>
  <XCard :title="t('xray.template')" :description="t('xray.templateDesc')" seamless>
    <template #actions>
      <XSegmented v-model="part" :options="parts" :label="t('xray.template')" scroll />
    </template>
    <XCodeEditor v-model="text" :label="t('xray.template')" />
  </XCard>
</template>
