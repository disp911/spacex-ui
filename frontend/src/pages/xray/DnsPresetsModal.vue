<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { XButton, XChip, XModal, XPlainTable } from '../../kit'
import { dnsPresets } from './model'

/** Ready lists of DNS servers; installing one replaces the servers and closes the window. */
const open = defineModel<boolean>('open', { default: false })
const emit = defineEmits<{ install: [servers: string[]] }>()
const { t } = useI18n()

const columns = [{ key: 'name' }, { key: 'install', fit: true }]

function install(servers: string[]) {
  emit('install', [...servers])
  open.value = false
}
</script>

<template>
  <XModal v-model:open="open" :title="t('xray.dnsPresets')">
    <XPlainTable :columns="columns" :rows="dnsPresets" :row-key="(p) => p.name">
      <template #cell-name="{ row }">
        <span class="presets__name">
          {{ row.name }}
          <XChip :tone="row.family ? 'unl' : 'neutral'">{{ row.family ? t('xray.dnsFamily') : 'DNS' }}</XChip>
        </span>
      </template>
      <template #cell-install="{ row }">
        <XButton size="sm" @click="install(row.servers)">{{ t('xray.dnsInstall') }}</XButton>
      </template>
    </XPlainTable>
  </XModal>
</template>

<style scoped>
.presets__name {
  display: flex;
  align-items: center;
  gap: 10px;
  color: var(--fg);
  font: 600 var(--fs-13) / var(--lh-tight) var(--font-sans);
}
</style>
