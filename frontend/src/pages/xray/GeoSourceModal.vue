<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { api } from '../../api'
import { toast, XButton, XField, XInput, XModal, XSelect } from '../../kit'
import { cause } from '../settings/form'
import type { GeoSource } from './model'

/** A geo file of the user's own to add, or the one given to edit; the server takes it at once. */
const props = defineProps<{ source: GeoSource | null }>()
const open = defineModel<boolean>('open', { default: false })
const emit = defineEmits<{ saved: [] }>()
const { t } = useI18n()

const form = reactive({ type: 'geosite' as GeoSource['type'], alias: '', url: '' })
const errors = reactive({ alias: '', url: '' })
const saving = ref(false)

watch(open, (on) => {
  if (!on) return
  Object.assign(form, props.source ? { type: props.source.type, alias: props.source.alias, url: props.source.url } : { type: 'geosite', alias: '', url: '' })
  Object.assign(errors, { alias: '', url: '' })
})

// In the script: vue-tsc misreads the // of a URL inside the template.
const urlHint = 'https://'

const types = [
  { value: 'geosite' as const, label: 'geosite' },
  { value: 'geoip' as const, label: 'geoip' },
]

function check() {
  errors.alias = /^[a-z0-9_-]+$/.test(form.alias.trim()) ? '' : t('xray.geoAliasInvalid')
  errors.url = /^https?:\/\//i.test(form.url.trim()) ? '' : t('xray.geoUrlInvalid')
  return !errors.alias && !errors.url
}

async function save() {
  if (!check()) return
  saving.value = true
  const body = { type: form.type, alias: form.alias.trim(), url: form.url.trim() }
  try {
    await api.post(props.source ? `panel/api/custom-geo/update/${props.source.id}` : 'panel/api/custom-geo/add', body)
    emit('saved')
    open.value = false
  } catch (e) {
    // A taken alias belongs under its field; other refusals come as a toast.
    const why = cause(e)
    if (why === t('xray.geoAliasTaken')) errors.alias = why
    else toast('err', why)
  } finally {
    saving.value = false
  }
}

const title = computed(() => (props.source ? t('xray.geoSourceEdit') : t('xray.geoSourceAdd')))
</script>

<template>
  <XModal v-model:open="open" :title="title" size="sm" cancel :dismissible="!saving">
    <div class="geo-source">
      <XField v-slot="a" :label="t('xray.geoType')">
        <XSelect v-bind="a" v-model="form.type" :options="types" mono :disabled="!!source" />
      </XField>
      <XField v-slot="a" :label="t('xray.geoAlias')" :error="errors.alias">
        <XInput v-bind="a" v-model="form.alias" mono :placeholder="t('xray.geoAliasPlaceholder')" :disabled="!!source" :invalid="!!errors.alias" />
      </XField>
      <XField v-slot="u" :label="t('xray.geoUrl')" :error="errors.url">
        <XInput v-bind="u" v-model="form.url" mono :placeholder="urlHint" :invalid="!!errors.url" />
      </XField>
    </div>
    <template #footer>
      <XButton variant="primary" :loading="saving" @click="save">{{ t('xray.save') }}</XButton>
    </template>
  </XModal>
</template>

<style scoped>
.geo-source {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
</style>
