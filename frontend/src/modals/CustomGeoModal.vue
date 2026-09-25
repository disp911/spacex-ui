<script lang="ts">
export interface CustomGeo {
  id: number
  type: 'geosite' | 'geoip'
  alias: string
  url: string
  lastUpdatedAt: number
}
</script>

<script setup lang="ts">
import { computed, reactive, ref } from 'vue'
import { api } from '../lib/api'
import { dateTime } from '../lib/format'
import { t } from '../lib/i18n'
import Icon from '../components/Icon.vue'
import Modal from '../components/Modal.vue'
import Spinner from '../components/Spinner.vue'

const props = defineProps<{ source: CustomGeo | null }>()
const emit = defineEmits<{ close: []; saved: [] }>()

const editing = computed(() => !!props.source)
const form = reactive({
  type: props.source?.type ?? ('geosite' as 'geosite' | 'geoip'),
  alias: props.source?.alias ?? '',
  url: props.source?.url ?? '',
})
const touched = reactive({ alias: false, url: false })
const saving = ref(false)

const aliasError = computed(() =>
  !editing.value && touched.alias && !/^[a-z0-9_-]+$/.test(form.alias) ? t('pages.index.customGeoValidationAlias') : '',
)
const urlError = computed(() => {
  if (!touched.url) return ''
  try {
    const u = new URL(form.url.trim())
    return u.protocol === 'http:' || u.protocol === 'https:' ? '' : t('spx.urlInvalid')
  } catch {
    return t('spx.urlInvalid')
  }
})
const preview = computed(() => `ext:${form.type}_${form.alias || 'alias'}.dat:${t('spx.tag')}`)

async function save() {
  touched.alias = true
  touched.url = true
  if (aliasError.value || urlError.value || (!editing.value && !form.alias)) return
  saving.value = true
  const body = { type: form.type, alias: form.alias, url: form.url.trim() }
  const msg = editing.value
    ? await api.post(`panel/api/custom-geo/update/${props.source!.id}`, body)
    : await api.post('panel/api/custom-geo/add', body)
  saving.value = false
  if (msg.success) emit('saved')
}
</script>

<template>
  <Modal
    :title="editing ? t('pages.index.customGeoModalEdit') : t('pages.index.customGeoModalAdd')"
    :subtitle="editing ? `${source!.type}_${source!.alias}.dat` : undefined"
    :width="464"
    :closable="!saving"
    @close="emit('close')"
  >
    <form class="form" @submit.prevent="save">
      <div class="field">
        <div class="caps">{{ t('pages.index.customGeoType') }}<Icon v-if="editing" name="lock" :size="11" /></div>
        <div v-if="!editing" class="seg">
          <button type="button" :class="{ on: form.type === 'geosite' }" @click="form.type = 'geosite'">geosite</button>
          <button type="button" :class="{ on: form.type === 'geoip' }" @click="form.type = 'geoip'">geoip</button>
        </div>
        <div v-else class="input locked">{{ form.type }}</div>
      </div>

      <div class="field">
        <div class="caps">{{ t('pages.index.customGeoAlias') }}<Icon v-if="editing" name="lock" :size="11" /></div>
        <div v-if="editing" class="input locked">{{ form.alias }}</div>
        <input
          v-else
          v-model.trim="form.alias"
          class="input"
          :class="{ err: aliasError }"
          placeholder="myads"
          autocapitalize="off"
          spellcheck="false"
          @blur="touched.alias = true"
        />
        <div v-if="aliasError" class="error"><Icon name="alert" :size="13" />{{ aliasError }}</div>
        <div v-else class="hint">
          <template v-if="editing">{{ t('spx.aliasLocked') }}</template>
          <template v-else>{{ t('spx.aliasAllowed') }} <span class="mono">a-z 0-9 _ -</span> · {{ t('spx.aliasInFileName') }}</template>
        </div>
      </div>

      <div class="field">
        <div class="caps">URL</div>
        <input
          v-model="form.url"
          class="input"
          :class="{ err: urlError }"
          :placeholder="`https://example.com/${form.type}_${form.alias || 'myads'}.dat`"
          inputmode="url"
          autocapitalize="off"
          spellcheck="false"
          @blur="touched.url = true"
        />
        <div v-if="urlError" class="error"><Icon name="alert" :size="13" />{{ urlError }}</div>
      </div>

      <div v-if="!editing" class="preview">
        <span class="caps-sm">{{ t('spx.inRules') }}</span>
        <span class="mono accent">{{ preview }}</span>
      </div>
      <button type="submit" hidden />
    </form>

    <template #footer>
      <span v-if="editing && source!.lastUpdatedAt" class="updated">{{ t('spx.updatedAt', { date: dateTime(source!.lastUpdatedAt) }) }}</span>
      <button class="ctl cancel" :disabled="saving" @click="emit('close')">{{ t('cancel') }}</button>
      <button class="ctl primary save" @click="save">
        <Spinner v-if="saving" :size="14" on-accent />{{ t('pages.index.customGeoModalSave') }}
      </button>
    </template>
  </Modal>
</template>

<style scoped>
.form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}
.caps {
  display: flex;
  align-items: center;
  gap: 5px;
  margin-bottom: 7px;
}
.seg {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 8px;
}
.seg button {
  height: 40px;
  border-radius: 11px;
  border: 1px solid var(--border);
  background: var(--fill);
  font: 500 12.5px var(--font-mono);
  color: var(--text-3);
  transition:
    border-color 0.18s ease,
    background 0.18s ease,
    box-shadow 0.18s ease,
    color 0.18s ease;
}
.seg button:hover {
  border-color: var(--glow-border);
}
.seg button.on {
  border-color: var(--accent);
  background: var(--accent-soft);
  color: var(--accent);
  box-shadow: var(--focus-ring);
}
.input {
  display: flex;
  align-items: center;
  width: 100%;
  height: 44px;
  padding: 0 13px;
  border-radius: 12px;
  background: var(--fill);
  border: 1px solid var(--border);
  outline: 0;
  font: 450 12.5px var(--font-mono);
  color: var(--text);
  transition:
    border-color 0.22s ease,
    box-shadow 0.22s ease;
}
input.input:hover,
input.input:focus {
  border-color: var(--glow-border);
  box-shadow: var(--glow);
}
.input::placeholder {
  color: var(--faint);
}
.input.locked {
  color: var(--text-3b);
  cursor: not-allowed;
}
.input.err,
.input.err:focus {
  border-color: var(--danger);
  box-shadow: 0 0 0 3px rgba(var(--danger-rgb), 0.13);
}
.hint {
  margin-top: 7px;
  font: 400 11.5px var(--font-sans);
  color: var(--text-3b);
}
.mono {
  font-family: var(--font-mono);
  color: var(--text-2);
}
.error {
  margin-top: 7px;
  display: flex;
  align-items: center;
  gap: 6px;
  font: 400 11.5px var(--font-sans);
  color: var(--danger-text);
}
.preview {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  border-radius: 11px;
  background: var(--fill);
  border: 1px solid var(--border);
  overflow: hidden;
}
.accent {
  color: var(--accent);
  font: 500 11.5px var(--font-mono);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.updated {
  margin-right: auto;
  font: 400 11.5px var(--font-sans);
  color: var(--text-3b);
}
.cancel {
  background: var(--chip);
  padding: 0 15px;
}
.save {
  padding: 0 15px;
}
</style>
