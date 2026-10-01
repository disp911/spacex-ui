<script setup lang="ts">
import XCard from '@/components/ui/XCard.vue'
import XTag from '@/components/ui/XTag.vue'
import XStat from '@/components/ui/XStat.vue'
import XButton from '@/components/ui/XButton.vue'
import StateBadge from './StateBadge.vue'
import { useIsMobile } from '@/composables/useMedia'
import { t } from '@/i18n'

// A managed process (Xray core, Telegram proxy): state, stats, controls.
withDefaults(defineProps<{
  title: string
  version?: string
  state?: string
  errorMsg?: string
  online?: string
  uptime?: string
  memory?: string
  busy?: 'restart' | 'stop' | null
  /** Show restart/stop buttons. */
  controls?: boolean
}>(), { controls: true })
defineEmits<{ restart: []; stop: [] }>()
const isMobile = useIsMobile()
</script>

<template>
  <XCard :title="title" class="svc">
    <template v-if="version" #aside>
      <XTag tone="version">{{ version }}</XTag>
    </template>

    <div class="svc__state">
      <StateBadge :state="state" />
      <div v-if="controls" class="svc__actions svc__actions--mobile">
        <XButton icon="restart" icon-only :label="t('dash.restart')" :loading="busy === 'restart'" @click="$emit('restart')" />
        <XButton icon="stop" icon-only variant="danger" :label="t('dash.stop')" :loading="busy === 'stop'" @click="$emit('stop')" />
      </div>
    </div>

    <pre v-if="state === 'error' && errorMsg" class="svc__error">{{ errorMsg }}</pre>

    <div class="svc__row">
      <XStat :label="isMobile ? t('dash.onlineShort') : t('dash.online')" :value="online ?? '—'" class="svc__stat" />
      <XStat :label="isMobile ? t('dash.uptimeShort') : t('dash.uptime')" :value="uptime ?? '—'" class="svc__stat" />
      <XStat :label="t('dash.memory')" :value="memory ?? '—'" class="svc__stat" />
      <div v-if="controls" class="svc__actions svc__actions--desktop">
        <XButton icon="restart" icon-only :label="t('dash.restart')" :loading="busy === 'restart'" @click="$emit('restart')" />
        <XButton icon="stop" icon-only variant="danger" :label="t('dash.stop')" :loading="busy === 'stop'" @click="$emit('stop')" />
      </div>
    </div>

    <div v-if="$slots.default" class="svc__links"><slot /></div>
  </XCard>
</template>

<style scoped>
.svc__state {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 12px;
}
.svc__error {
  margin: 0 0 14px;
  padding: 10px 12px;
  max-height: 92px;
  overflow: auto;
  white-space: pre-wrap;
  word-break: break-word;
  background: var(--danger-soft);
  border: 1px solid var(--danger-border);
  border-radius: var(--r-md);
  font: var(--fw-regular) var(--fs-xs) / 1.5 var(--font-mono);
  color: var(--danger-text);
}
.svc__row {
  display: flex;
  align-items: flex-end;
  gap: 12px;
}
.svc__stat {
  flex: 1;
}
.svc__actions {
  flex: none;
  display: flex;
  gap: 8px;
}
.svc__actions--mobile {
  display: none;
  margin-left: auto;
}
.svc__links {
  display: flex;
  flex-direction: column;
  gap: 2px;
  margin-top: 16px;
  padding-top: 11px;
  border-top: 1px solid var(--divider);
}

@media (max-width: 767px) {
  .svc__state {
    margin-bottom: 14px;
  }
  .svc__actions--mobile {
    display: flex;
  }
  .svc__actions--desktop {
    display: none;
  }
  .svc__row {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 10px;
  }
  .svc__links {
    margin-top: 12px;
    padding-top: 10px;
  }
}
</style>
