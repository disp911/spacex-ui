<script setup lang="ts">
import { ref } from 'vue'
import XCard from '@/components/ui/XCard.vue'
import XButton from '@/components/ui/XButton.vue'
import type { ServerStatus } from '@/api/panel'
import { t } from '@/i18n'

defineProps<{ status: ServerStatus | null }>()
// Hidden by default so the addresses don't leak in screenshots.
const hidden = ref(true)
</script>

<template>
  <XCard :title="t('dash.ip')">
    <template #aside>
      <XButton icon="eye" icon-only size="xs" :label="t('dash.toggleIp')" @click="hidden = !hidden" />
    </template>
    <div class="ip">
      <div class="ip__tile">
        <div class="caps ip__label">IPv4</div>
        <div class="ip__value" :class="{ 'is-hidden': hidden }">{{ status?.publicIP.ipv4 || 'N/A' }}</div>
      </div>
      <div class="ip__tile">
        <div class="caps ip__label">IPv6</div>
        <div class="ip__value" :class="{ 'is-hidden': hidden }">{{ status?.publicIP.ipv6 || 'N/A' }}</div>
      </div>
    </div>
  </XCard>
</template>

<style scoped>
.ip {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}
.ip__tile {
  min-width: 0;
  padding: 10px 14px;
  background: var(--field);
  border: 1px solid var(--border-field);
  border-radius: var(--r-md);
}
.ip__label {
  margin-bottom: 6px;
  text-transform: none;
}
.ip__value {
  font: var(--fw-medium) var(--fs-lg) var(--font-mono);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  transition: filter var(--dur-base) ease;
}
.ip__value.is-hidden {
  filter: blur(5px);
  user-select: none;
}
@media (max-width: 767px) {
  .ip {
    gap: 10px;
  }
  .ip__tile {
    padding: 11px 13px;
  }
}
</style>
