<script setup lang="ts">
import { computed } from 'vue'
import XCard from '@/components/ui/XCard.vue'
import XIcon from '@/components/ui/XIcon.vue'
import XStat from '@/components/ui/XStat.vue'
import type { ServerStatus } from '@/api/panel'
import { formatSpeed, sizeParts } from '@/utils/format'
import { t } from '@/i18n'

const props = defineProps<{ status: ServerStatus | null }>()
const sent = computed(() => sizeParts(props.status?.netTraffic.sent ?? 0))
const recv = computed(() => sizeParts(props.status?.netTraffic.recv ?? 0))
</script>

<template>
  <XCard :title="t('dash.network')">
    <div class="net__totals">
      <div class="tile">
        <div class="caps tile__label"><XIcon name="arrow-up" :size="13" :stroke-width="1.5" class="is-info" />
          <span class="d">{{ t('dash.sentTotal') }}</span><span class="m">{{ t('dash.sentShort') }}</span>
        </div>
        <div class="tile__value is-info">{{ sent[0] }} <small>{{ sent[1] }}</small></div>
      </div>
      <div class="tile">
        <div class="caps tile__label"><XIcon name="arrow-down" :size="13" :stroke-width="1.5" class="is-accent" />
          <span class="d">{{ t('dash.recvTotal') }}</span><span class="m">{{ t('dash.recvShort') }}</span>
        </div>
        <div class="tile__value is-accent">{{ recv[0] }} <small>{{ recv[1] }}</small></div>
      </div>
    </div>
    <div class="net__live">
      <XStat :label="t('dash.upload')" :value="formatSpeed(status?.netIO.up ?? 0)" tone="info" class="net__s" />
      <XStat :label="t('dash.downloadSpeed')" :value="formatSpeed(status?.netIO.down ?? 0)" tone="accent" class="net__s" />
      <XStat label="TCP" :value="status?.tcpCount ?? '—'" tone="secondary" class="net__s" />
      <XStat label="UDP" :value="status?.udpCount ?? '—'" tone="secondary" class="net__s" />
    </div>
  </XCard>
</template>

<style scoped>
.net__totals {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
  margin-bottom: 18px;
}
.tile {
  min-width: 0;
  padding: 12px 14px;
  background: var(--field);
  border: 1px solid var(--border-field);
  border-radius: var(--r-md);
}
.tile__label {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 7px;
}
.tile__label .m {
  display: none;
}
.tile__value {
  font: var(--fw-semibold) var(--fs-2xl) var(--font-mono);
  letter-spacing: -0.5px;
  white-space: nowrap;
}
.tile__value small {
  font-size: var(--fs-md);
  font-weight: var(--fw-medium);
  color: var(--text-3);
  letter-spacing: 0;
}
.is-info {
  color: var(--info);
}
.is-accent {
  color: var(--accent);
}
.net__live {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 12px;
  padding: 0 14px;
}
.net__s :deep(.x-stat__value) {
  font-size: var(--fs-lg);
}
@media (max-width: 767px) {
  .net__totals {
    gap: 10px;
    margin-bottom: 14px;
  }
  .tile {
    padding: 11px 13px;
  }
  .tile__label .d {
    display: none;
  }
  .tile__label .m {
    display: inline;
  }
  .net__live {
    padding: 0;
    gap: 10px;
    grid-template-columns: 1fr 1fr;
  }
}
</style>
