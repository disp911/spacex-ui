<script setup lang="ts">
import { computed, nextTick, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { api, ApiError } from '../api'
import {
  toast, XActionBar, XBusyOverlay, XButton, XCard, XConfirm, XEmpty, XNotice, XSettingsSkeleton, XStatusText, XTabs, type TabItem,
} from '../kit'
import AppPage from '../layout/AppPage.vue'
import { securityAlerts, tabs, type AllSetting, type SettingKey, type TabKey } from './settings/model'
import { addressAfterRestart, waitForPanel } from './settings/restart'

const { t } = useI18n()

// --- Loading ----------------------------------------------------------------

/** As the server has them, and as edited on the page. */
const saved = ref<AllSetting | null>(null)
const draft = ref<AllSetting | null>(null)
const loading = ref(true)
const failed = ref(false)

/** quiet keeps the page as it is, e.g. after a save. */
async function load(quiet = false) {
  if (!quiet) {
    loading.value = true
    failed.value = false
  }
  try {
    const s = await api.post<AllSetting>('panel/setting/all')
    saved.value = s
    draft.value = { ...s }
  } catch {
    if (!quiet) failed.value = true
  } finally {
    loading.value = false
  }
}
onMounted(() => load())

const ready = computed(() => !loading.value && !failed.value && !!draft.value)

// --- Changes and errors -------------------------------------------------------

const isDirty = (k: SettingKey) => !!draft.value && !!saved.value && draft.value[k] !== saved.value[k]
/** Messages of fields that do not pass their checks; the fields add them. */
const errors = computed<Partial<Record<SettingKey, string>>>(() => ({}))

const count = (keys: SettingKey[]) => ({
  dirty: keys.filter(isDirty).length,
  errors: keys.filter((k) => errors.value[k]).length,
})
const total = computed(() => count(tabs.flatMap((tab) => tab.keys)))

const status = computed(() => {
  const { dirty, errors: n } = total.value
  if (!ready.value || (!dirty && !n)) return null
  if (!n) return { tone: 'mod' as const, text: t('settings.unsaved', { n: dirty }) }
  return { tone: 'err' as const, text: dirty ? t('settings.fixErrorsUnsaved', { errors: n, n: dirty }) : t('settings.fixErrors', { n }) }
})

// --- Tabs -------------------------------------------------------------------

const tab = ref<TabKey>('panel')
/** Formats has something to set only while JSON or Clash subscriptions are on. */
const formats = computed(() => !!draft.value && (draft.value.subJsonEnable || draft.value.subClashEnable))
const tabItems = computed<TabItem<TabKey>[]>(() =>
  tabs
    .filter((item) => item.key !== 'formats' || formats.value)
    .map((item) => {
      const c = count(item.keys)
      return { key: item.key, label: t(item.label), dot: c.errors ? 'err' : c.dirty ? 'mod' : undefined }
    }),
)
watch(formats, (on) => {
  if (!on && tab.value === 'formats') tab.value = 'sub'
})
// A tab starts from its top.
watch(tab, () => nextTick(() => window.scrollTo({ top: 0 })))

// --- Security banner --------------------------------------------------------

const bannerClosed = ref(false)
const alerts = computed(() => (saved.value ? securityAlerts(saved.value).map((key) => t(key)) : []))
const showBanner = computed(() => ready.value && !bannerClosed.value && alerts.value.length > 0)

// --- Save -------------------------------------------------------------------

const saving = ref(false)
const canSave = computed(() => ready.value && total.value.dirty > 0 && !total.value.errors && !saving.value)

/** The server wraps the cause in parentheses after its own words; the cause is what helps. */
const cause = (e: unknown) => {
  const msg = e instanceof ApiError ? e.message : String(e)
  return /\((.+)\)$/s.exec(msg)?.[1] ?? msg
}

async function save() {
  if (!canSave.value || !draft.value) return
  saving.value = true
  try {
    await api.post('panel/setting/update', draft.value)
    saved.value = { ...draft.value }
    toast('ok', t('settings.saved'))
    // The server may tidy values up, e.g. the slashes of a path.
    await load(true)
  } catch (e) {
    toast('err', t('settings.saveFailed', { error: cause(e) }))
  } finally {
    saving.value = false
  }
}

// --- Restart ----------------------------------------------------------------

const canRestart = computed(() => ready.value && !total.value.dirty && !saving.value)
const confirming = ref(false)
const restarting = ref(false)
const restarted = ref(false)
const progress = ref(0)

async function restart() {
  if (!saved.value) return
  confirming.value = false
  restarted.value = false
  progress.value = 0
  restarting.value = true
  const { panel, page } = addressAfterRestart(saved.value)
  try {
    await api.post('panel/setting/restartPanel')
  } catch (e) {
    restarting.value = false
    toast('err', t('settings.restartFailed', { error: cause(e) }))
    return
  }
  // The bar runs most of the way while the panel is away, and fills when it answers.
  requestAnimationFrame(() => requestAnimationFrame(() => (progress.value = 72)))
  await waitForPanel(panel)
  progress.value = 100
  restarted.value = true
  setTimeout(() => {
    // At a new address the panel opens afresh; at the same one the page reloads the settings.
    if (page !== location.origin + location.pathname) {
      location.replace(page)
      return
    }
    restarting.value = false
    load()
  }, 1000)
}
</script>

<template>
  <AppPage :title="t('nav.settings')" fixed>
    <template v-if="showBanner" #banner>
      <XNotice tone="warn" :title="t('settings.securityTitle')" :items="alerts" closable @close="bannerClosed = true">
        {{ t('settings.securityText') }}
      </XNotice>
    </template>
    <template #subtitle>{{ t('settings.subtitle') }}</template>
    <template #actions>
      <XStatusText v-if="status" :tone="status.tone">{{ status.text }}</XStatusText>
      <XButton variant="primary" icon="check" :disabled="!canSave && !saving" :loading="saving" :title="canSave || saving ? undefined : t('settings.nothingToSave')" @click="save">
        {{ saving ? t('settings.saving') : t('settings.save') }}
      </XButton>
      <XButton variant="primary" icon="refresh" :disabled="!canRestart" :title="canRestart ? undefined : t('settings.saveFirst')" @click="confirming = true">
        {{ t('settings.restart') }}
      </XButton>
    </template>
    <template #tabs>
      <XTabs v-model="tab" :items="tabItems" :label="t('settings.sections')" />
    </template>

    <XSettingsSkeleton v-if="loading" />
    <XCard v-else-if="failed">
      <XEmpty tone="err" icon="alert" :title="t('settings.loadFailed')" />
    </XCard>
    <XCard v-else>
      <XEmpty :title="t('common.screenDraft')" />
    </XCard>

    <template #bar>
      <XActionBar>
        <template #state>
          <XStatusText v-if="status" :tone="status.tone">{{ status.text }}</XStatusText>
          <span v-else>{{ t('settings.subtitle') }}</span>
        </template>
        <XButton variant="primary" :disabled="!canSave && !saving" :loading="saving" @click="save">
          {{ saving ? t('settings.saving') : t('settings.save') }}
        </XButton>
        <XButton variant="primary" icon="refresh" :disabled="!canRestart" @click="confirming = true">{{ t('settings.restartShort') }}</XButton>
      </XActionBar>
    </template>
  </AppPage>

  <XConfirm v-model:open="confirming" icon="refresh" :title="t('settings.restart')" :cta="t('settings.restartYes')" @confirm="restart">
    {{ t('settings.restartText') }}
  </XConfirm>
  <XBusyOverlay :open="restarting" :title="t('settings.restart')" :text="t(restarted ? 'settings.restartDone' : 'settings.restartWait')" :progress="progress" />
</template>
