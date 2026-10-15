<script setup lang="ts">
import { onBeforeUnmount, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { api, ApiError, saveFile } from '../../api'
import { copyText, toast, XButton, XCodeArea, XEmpty, XModal } from '../../kit'

/** The running Xray config.json: read, edit here for a copy or a download; nothing is saved to the server. */
const open = defineModel<boolean>('open', { default: false })
const { t } = useI18n()

const text = ref('')
const loading = ref(false)
const failed = ref('')

async function load() {
  loading.value = true
  failed.value = ''
  try {
    text.value = JSON.stringify(await api.get<unknown>('panel/api/server/getConfigJson'), null, 2)
  } catch (e) {
    failed.value = e instanceof ApiError ? e.message : String(e)
  } finally {
    loading.value = false
  }
}

watch(open, (on) => on && load(), { immediate: true })

const copied = ref(false)
let reset: ReturnType<typeof setTimeout> | undefined

async function copy() {
  if (!(await copyText(text.value))) {
    toast('err', t('dashboard.configCopyFailed'))
    return
  }
  copied.value = true
  clearTimeout(reset)
  reset = setTimeout(() => (copied.value = false), 1600)
}
onBeforeUnmount(() => clearTimeout(reset))

function download() {
  saveFile(text.value, 'config.json')
}
</script>

<template>
  <XModal v-model:open="open" :title="t('dashboard.config')" full>
    <template #subtitle><span class="mono">config.json</span></template>

    <XEmpty v-if="failed" :title="t('dashboard.configFailed', { error: failed })">
      <XButton size="sm" icon="refresh" @click="load">{{ t('dashboard.retry') }}</XButton>
    </XEmpty>
    <XCodeArea v-else v-model="text" :loading="loading" :height="420" phone-fill label="config.json" />

    <template #footer>
      <XButton icon="download" mono :disabled="loading || !!failed" @click="download">config.json</XButton>
      <XButton variant="primary" icon="copy" :success="copied" :disabled="loading || !!failed" @click="copy">
        {{ copied ? t('dashboard.configCopied') : t('dashboard.configCopy') }}
      </XButton>
    </template>
  </XModal>
</template>
