<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { api, ApiError } from '../../api'
import { toast, XButton, XConfirm, XEmpty, XModal, XNotice, XPickList, type PickItem } from '../../kit'
import { useServerStatus } from '../../stores/server'

/** Xray releases from GitHub; picking one asks, then installs it and restarts the core. */
const open = defineModel<boolean>('open', { default: false })
const { t } = useI18n()
const server = useServerStatus()

const versions = ref<string[]>([])
const loading = ref(false)
const failed = ref('')

const message = (e: unknown) => (e instanceof ApiError ? e.message : String(e))

// The status reports 26.4.25, releases are tagged v26.4.25.
const installed = computed(() => (server.status?.xray.version ? `v${server.status.xray.version}` : ''))
const items = computed<PickItem[]>(() => versions.value.map((v) => ({ key: v, label: v, current: v === installed.value })))

let request = 0
async function load() {
  const id = ++request
  loading.value = true
  failed.value = ''
  try {
    const reply = await api.get<string[] | null>('panel/api/server/getXrayVersion')
    if (id === request) versions.value = reply ?? []
  } catch (e) {
    if (id !== request) return
    versions.value = []
    failed.value = message(e)
  } finally {
    if (id === request) loading.value = false
  }
}

watch(open, (on) => on && load(), { immediate: true })

const target = ref('')
const confirming = ref(false)
const installing = ref(false)

function pick(version: string) {
  target.value = version
  confirming.value = true
}

async function install() {
  installing.value = true
  try {
    await api.post(`panel/api/server/installXray/${target.value}`)
    await server.refresh().catch(() => {})
  } catch (e) {
    toast('err', t('dashboard.actionFailed', { error: message(e) }))
  } finally {
    installing.value = false
    confirming.value = false
  }
}
</script>

<template>
  <XModal v-model:open="open" :title="t('dashboard.xrayUpdate')">
    <template #subtitle>{{ t('dashboard.xrayVersion') }}</template>

    <XNotice tone="warn" :title="t('dashboard.xrayVersionWarnTitle')">{{ t('dashboard.xrayVersionWarn') }}</XNotice>

    <XPickList :items="items" :loading="loading" :current-label="t('dashboard.xrayInstalled')" @pick="pick">
      <XEmpty :title="failed ? t('dashboard.xrayVersionsFailed', { error: failed }) : t('dashboard.xrayVersionsEmpty')">
        <XButton v-if="failed" size="sm" icon="refresh" @click="load">{{ t('dashboard.retry') }}</XButton>
      </XEmpty>
    </XPickList>

    <template #footer>
      <XButton @click="open = false">{{ t('kit.cancel') }}</XButton>
    </template>
  </XModal>

  <XConfirm
    v-model:open="confirming"
    :title="t('dashboard.xrayInstallTitle')"
    :cta="t('dashboard.xrayInstall')"
    :busy="installing"
    :busy-label="t('dashboard.xrayInstalling')"
    @confirm="install"
  >
    {{ installed ? t('dashboard.xrayInstallText', { version: target, current: installed }) : t('dashboard.xrayInstallTextNew', { version: target }) }}
  </XConfirm>
</template>
