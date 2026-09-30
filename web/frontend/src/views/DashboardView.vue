<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import AppShell from '@/components/shell/AppShell.vue'
import XAlert from '@/components/ui/XAlert.vue'
import XListLink from '@/components/ui/XListLink.vue'
import ServiceCard from '@/features/dashboard/ServiceCard.vue'
import ResourcesCard from '@/features/dashboard/ResourcesCard.vue'
import NetworkCard from '@/features/dashboard/NetworkCard.vue'
import SystemCard from '@/features/dashboard/SystemCard.vue'
import IpCard from '@/features/dashboard/IpCard.vue'
import ManageCard from '@/features/dashboard/ManageCard.vue'
import ConnectionsLogModal from '@/features/modals/ConnectionsLogModal.vue'
import XrayUpdatesModal from '@/features/modals/XrayUpdatesModal.vue'
import PanelLogModal from '@/features/modals/PanelLogModal.vue'
import ConfigModal from '@/features/modals/ConfigModal.vue'
import BackupModal from '@/features/modals/BackupModal.vue'
import PanelUpdateModal from '@/features/modals/PanelUpdateModal.vue'
import { server, type PanelUpdateInfo } from '@/api/panel'
import { useServerStatus } from '@/composables/useServerStatus'
import { toastError, toastMsg } from '@/composables/useToast'
import { env } from '@/env'
import { formatSize, formatUptime, withV } from '@/utils/format'
import { t } from '@/i18n'

const { status, online, cpuRange, cpuSeries, refresh } = useServerStatus()

// ---- HTTP warning ----------------------------------------------------
const HTTP_KEY = 'spx-http-warning-dismissed'
const httpWarning = ref(location.protocol !== 'https:' && sessionStorage.getItem(HTTP_KEY) !== '1')
function dismissHttp() {
  httpWarning.value = false
  sessionStorage.setItem(HTTP_KEY, '1')
}

// ---- Xray controls ---------------------------------------------------
const xrayBusy = ref<'restart' | 'stop' | null>(null)
async function xrayAction(kind: 'restart' | 'stop') {
  xrayBusy.value = kind
  try {
    toastMsg(await (kind === 'restart' ? server.restartXray() : server.stopXray()))
    await refresh()
  } catch (e) {
    toastError(e)
  } finally {
    xrayBusy.value = null
  }
}

// ---- Panel update info -----------------------------------------------
const updateInfo = ref<PanelUpdateInfo | null>(null)
onMounted(async () => {
  try {
    updateInfo.value = await server.panelUpdateInfo()
  } catch {
    /* offline / rate limited: the version still shows from boot data */
  }
})
const panelVersion = computed(() => updateInfo.value?.currentVersion || env.version)

// ---- Modals ----------------------------------------------------------
const modals = reactive({
  xrayLog: false,
  updates: false,
  panelLog: false,
  config: false,
  backup: false,
  panelUpdate: false,
})
type ModalKey = keyof typeof modals
const openModal = (k: ModalKey) => (modals[k] = true)

const xray = computed(() => status.value?.xray)
const app = computed(() => status.value?.appStats)
const tg = computed(() => status.value?.telemt)
</script>

<template>
  <AppShell :title="t('dash.title')" :subtitle="env.host">
    <template v-if="httpWarning" #aside>
      <XAlert tone="warning" :title="t('http.title')" dismissible compact class="http-warn" @dismiss="dismissHttp">
        {{ t('http.text') }}
      </XAlert>
    </template>

    <div class="grid">
      <div class="col">
        <ServiceCard
          class="o-xray"
          title="Xray"
          :version="withV(xray?.version ?? '')"
          :state="xray?.state"
          :error-msg="xray?.errorMsg"
          :online="online === null ? undefined : String(online)"
          :uptime="app ? formatUptime(app.uptime) : undefined"
          :memory="app ? formatSize(app.mem) : undefined"
          :busy="xrayBusy"
          @restart="xrayAction('restart')"
          @stop="xrayAction('stop')"
        >
          <XListLink icon="log" :label="t('dash.xrayLog')" @click="openModal('xrayLog')" />
          <XListLink icon="download" :label="t('dash.versionPick')" @click="openModal('updates')" />
        </ServiceCard>

        <ServiceCard
          v-if="tg"
          class="o-tg"
          :title="t('dash.telegram')"
          :version="tg.version"
          :state="tg.state"
          :error-msg="tg.errorMsg"
          :online="String(tg.online)"
          :uptime="formatUptime(tg.uptime)"
          :memory="formatSize(tg.mem)"
          :controls="false"
        />

        <ManageCard
          class="o-manage"
          :version="panelVersion"
          :latest="updateInfo?.latestVersion"
          :update-available="updateInfo?.updateAvailable"
          @open="openModal"
        />

        <IpCard class="o-ip" :status="status" />
      </div>

      <div class="col">
        <ResourcesCard v-model:range="cpuRange" class="o-res" :status="status" :series="cpuSeries" />
        <NetworkCard class="o-net" :status="status" />
        <SystemCard class="o-sys" :status="status" />
      </div>
    </div>

    <ConnectionsLogModal v-model:open="modals.xrayLog" />
    <XrayUpdatesModal v-model:open="modals.updates" :installed="xray?.version" @changed="refresh" />
    <PanelLogModal v-model:open="modals.panelLog" />
    <ConfigModal v-model:open="modals.config" />
    <BackupModal v-model:open="modals.backup" />
    <PanelUpdateModal v-model:open="modals.panelUpdate" :info="updateInfo" />
  </AppShell>
</template>

<style scoped>
.http-warn {
  min-height: 45px;
}
.grid {
  flex: 1 0 auto;
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
  align-items: stretch;
  max-width: var(--content-max);
  width: 100%;
}
.col {
  display: flex;
  flex-direction: column;
  gap: 12px;
  min-width: 0;
}

/* Narrow desktops / tablets: one column in reading order. */
@media (max-width: 1180px) {
  .grid {
    grid-template-columns: 1fr;
  }
}

/* Mobile: flatten the two columns and reorder cards as in the mockup. */
@media (max-width: 767px) {
  .grid {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }
  .col {
    display: contents;
  }
  .o-xray { order: 1; }
  .o-tg { order: 2; }
  .o-res { order: 3; }
  .o-net { order: 4; }
  .o-ip { order: 5; }
  .o-manage { order: 6; }
  .o-sys { order: 7; }
}
</style>
