<script setup lang="ts">
import XCard from '@/components/ui/XCard.vue'
import XListLink from '@/components/ui/XListLink.vue'
import XButton from '@/components/ui/XButton.vue'
import { t } from '@/i18n'
import { withV } from '@/utils/format'

defineProps<{ version: string; latest?: string; updateAvailable?: boolean }>()
defineEmits<{ open: ['panelLog' | 'config' | 'backup' | 'panelUpdate'] }>()
</script>

<template>
  <XCard :title="t('dash.manage')" fill>
    <div class="manage__links">
      <XListLink icon="log" :label="t('dash.panelLog')" @click="$emit('open', 'panelLog')" />
      <XListLink icon="code" :label="t('dash.config')" @click="$emit('open', 'config')" />
      <XListLink icon="database" :label="t('dash.backup')" @click="$emit('open', 'backup')" />
    </div>
    <template #footer>
      <div class="manage__foot">
        <span class="manage__ver-label">{{ t('dash.panelVersion') }}</span>
        <button type="button" class="manage__ver mono" @click="$emit('open', 'panelUpdate')">{{ withV(version) || '—' }}</button>
        <XButton v-if="updateAvailable" variant="soft" size="sm" class="manage__update" @click="$emit('open', 'panelUpdate')">
          <span class="manage__pulse" />
          {{ t('dash.updatePanelTo') }} <span class="mono">{{ withV(latest ?? '') }}</span>
        </XButton>
      </div>
    </template>
  </XCard>
</template>

<style scoped>
.manage__links {
  display: flex;
  flex-direction: column;
  gap: 2px;
}
.manage__foot {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 8px;
  min-height: 28px;
}
.manage__ver-label {
  font: var(--fw-regular) var(--fs-sm) var(--font-sans);
  color: var(--text-4);
  white-space: nowrap;
}
.manage__ver {
  padding: 0;
  border: 0;
  background: none;
  font-size: var(--fs-sm);
  font-weight: var(--fw-medium);
  color: var(--text-2);
}
.manage__ver:hover {
  color: var(--accent);
}
.manage__update {
  margin-left: auto;
}
.manage__pulse {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--accent);
  animation: spx-pulse 2.4s infinite;
}
@media (max-width: 767px) {
  .manage__update {
    margin: 4px 0 0;
    width: 100%;
    height: 42px;
    font-size: var(--fs-md);
  }
}
</style>
