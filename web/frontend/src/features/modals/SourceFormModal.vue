<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import XModal from '@/components/ui/XModal.vue'
import XField from '@/components/ui/XField.vue'
import XInput from '@/components/ui/XInput.vue'
import XSegmented from '@/components/ui/XSegmented.vue'
import XButton from '@/components/ui/XButton.vue'
import { customGeo, type CustomGeo } from '@/api/panel'
import { toastError, toastMsg } from '@/composables/useToast'
import { formatDateTime } from '@/utils/format'
import { t } from '@/i18n'

// Add or edit a custom geosite/geoip source. When editing, type and alias
// are locked because they form the file name.
const open = defineModel<boolean>('open', { default: false })
const props = defineProps<{ source?: CustomGeo | null }>()
const emit = defineEmits<{ saved: [] }>()

const TYPES = ['geosite', 'geoip'] as const
const type = ref<(typeof TYPES)[number]>('geosite')
const alias = ref('')
const url = ref('')
const touched = ref(false)
const saving = ref(false)

const editing = computed(() => !!props.source)

watch(open, (v) => {
  if (!v) return
  touched.value = false
  type.value = props.source?.type ?? 'geosite'
  alias.value = props.source?.alias ?? ''
  url.value = props.source?.url ?? ''
})

const aliasError = computed(() => (touched.value && !editing.value && !/^[a-z0-9_-]+$/.test(alias.value) ? t('source.aliasError') : ''))
const urlError = computed(() => {
  if (!touched.value) return ''
  try {
    const u = new URL(url.value.trim())
    return u.protocol === 'http:' || u.protocol === 'https:' ? '' : t('source.urlError')
  } catch {
    return t('source.urlError')
  }
})
const fileName = computed(() => `${type.value}_${alias.value || 'alias'}.dat`)
const aliasHint = computed(() => t('source.aliasHint').split('{chars}'))

async function save() {
  touched.value = true
  if (aliasError.value || urlError.value) return
  saving.value = true
  try {
    const f = { type: type.value, alias: alias.value, url: url.value.trim() }
    const m = props.source ? await customGeo.update(props.source.id, f) : await customGeo.add(f)
    toastMsg(m)
    if (m.success) {
      open.value = false
      emit('saved')
    }
  } catch (e) {
    toastError(e)
  } finally {
    saving.value = false
  }
}
</script>

<template>
  <XModal v-model:open="open" :title="editing ? t('source.editTitle') : t('source.addTitle')" :width="464">
    <template v-if="editing" #subtitle><span class="mono">{{ fileName }}</span></template>

    <XField :label="t('source.type')" :locked="editing">
      <XInput v-if="editing" :model-value="type" locked mono />
      <XSegmented v-else v-model="type" :options="TYPES" mono />
    </XField>

    <XField :label="t('source.alias')" :locked="editing" :error="aliasError" :hint="editing ? t('source.aliasLocked') : undefined" for="src-alias">
      <XInput id="src-alias" v-model="alias" mono :locked="editing" :invalid="!!aliasError" placeholder="myads" />
      <template v-if="!editing" #hint>
        {{ aliasHint[0] }}<span class="mono">a-z 0-9 _ -</span>{{ aliasHint[1] }}
      </template>
    </XField>

    <XField :label="t('source.url')" :error="urlError" for="src-url">
      <XInput id="src-url" v-model="url" type="url" mono :invalid="!!urlError" placeholder="https://example.com/geosite_myads.dat" />
    </XField>

    <div v-if="!editing" class="rule">
      <span class="caps">{{ t('source.inRules') }}</span>
      <span class="rule__value ellipsis">ext:{{ fileName }}:tag</span>
    </div>

    <template #footer>
      <span v-if="source?.lastUpdatedAt" class="updated">{{ t('source.updatedAt', { at: formatDateTime(new Date(source.lastUpdatedAt * 1000)) }) }}</span>
      <XButton @click="open = false">{{ t('common.cancel') }}</XButton>
      <XButton variant="primary" :loading="saving" @click="save">{{ t('common.save') }}</XButton>
    </template>
  </XModal>
</template>

<style scoped>
.rule {
  display: flex;
  align-items: center;
  gap: 9px;
  padding: 9px 12px;
  background: var(--field);
  border: 1px solid var(--border-field);
  border-radius: var(--r-md);
}
.rule__value {
  font: var(--fw-regular) var(--fs-sm) var(--font-mono);
  color: var(--accent);
}
.updated {
  flex: 1;
  min-width: 0;
  font: var(--fw-regular) var(--fs-sm) var(--font-sans);
  color: var(--text-3);
}
@media (max-width: 767px) {
  .updated {
    display: none;
  }
}
</style>
