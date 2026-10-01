<script setup lang="ts">
import { computed } from 'vue'
import XSwitch from '@/components/ui/XSwitch.vue'
import XButton from '@/components/ui/XButton.vue'
import XMenu from '@/components/ui/XMenu.vue'
import type { MenuItem } from '@/components/ui/types'
import UsageCell from './UsageCell.vue'
import ExpiryCell from './ExpiryCell.vue'
import type { ClientRow, InboundRow } from './model'
import { t } from '@/i18n'

// Clients of one inbound: a dense grid on desktop, cards on mobile
// (same markup, the layout switches via grid areas).
const props = defineProps<{ inbound: InboundRow; busy: Set<string>; showIp?: boolean; filter?: string }>()
const emit = defineEmits<{
  toggle: [ClientRow, boolean]
  qr: [ClientRow]
  action: [string, ClientRow]
}>()

const rows = computed(() => {
  const q = props.filter?.trim().toLowerCase()
  const list = props.inbound.clients
  // When the page search matches a client email, show matches first.
  if (!q) return list
  return [...list].sort((a, b) => Number(b.email.toLowerCase().includes(q)) - Number(a.email.toLowerCase().includes(q)))
})

const menu: MenuItem[] = [
  { key: 'reset', label: t('ib.resetClient'), icon: 'restart' },
  { key: 'delete', label: t('ib.deleteClient'), icon: 'trash', danger: true, divided: true },
]
const isMatch = (c: ClientRow) => !!props.filter?.trim() && c.email.toLowerCase().includes(props.filter.trim().toLowerCase())
</script>

<template>
  <div class="clients" :class="{ 'with-ip': showIp }">
    <div class="clients__head caps">
      <span class="c-email">Email</span>
      <span class="c-usage">{{ t('ib.colTraffic') }}</span>
      <span class="c-exp">{{ t('ib.colExpiry') }}</span>
      <span v-if="showIp" class="c-ip">{{ t('ib.colIp') }}</span>
      <span class="c-act" />
    </div>
    <div
      v-for="c in rows"
      :key="c.email"
      class="client"
      :class="[`health-${c.health}`, { 'is-match': isMatch(c) }]"
    >
      <div class="c-email">
        <span class="dot" :class="{ 'is-online': c.online }" :title="c.online ? t('ib.online') : ''" />
        <div class="c-email__text">
          <span class="c-email__name ellipsis">{{ c.email }}</span>
          <span v-if="c.comment" class="c-email__comment ellipsis">{{ c.comment }}</span>
        </div>
      </div>
      <UsageCell class="c-usage" :up="c.up" :down="c.down" :total="c.total" :health="c.health" />
      <ExpiryCell class="c-exp" :expiry-time="c.expiryTime" />
      <span v-if="showIp" class="c-ip mono">{{ c.limitIp || '∞' }}</span>
      <div class="c-act">
        <XSwitch :model-value="c.enable" :loading="busy.has(c.email)" :label="t('ib.enable')" @update:model-value="emit('toggle', c, $event)" />
        <XButton icon="qr" icon-only size="xs" :label="t('ib.qr')" @click="emit('qr', c)" />
        <XMenu :items="menu" :label="t('ib.more')" :title="c.email" size="xs" @select="emit('action', $event, c)" />
      </div>
    </div>
  </div>
</template>

<style scoped>
.clients__head,
.client {
  display: grid;
  grid-template-columns: minmax(180px, 1.6fr) minmax(150px, 1.2fr) 118px 110px;
  grid-template-areas: 'email usage exp act';
  align-items: center;
  gap: 14px;
  padding: 0 14px;
}
.with-ip .clients__head,
.with-ip .client {
  grid-template-columns: minmax(180px, 1.6fr) minmax(150px, 1.2fr) 118px 48px 110px;
  grid-template-areas: 'email usage exp ip act';
}
.clients__head {
  height: 30px;
  font-size: var(--fs-2xs);
  border-bottom: 1px solid var(--divider);
}
.client {
  min-height: 48px;
  padding-top: 7px;
  padding-bottom: 7px;
  border-bottom: 1px solid var(--divider);
  transition: background var(--dur-fast) ease;
}
.client:last-child {
  border-bottom: 0;
}
.client:hover {
  background: var(--hover);
}
.client.is-match {
  background: var(--accent-soft);
}
.client.health-off {
  opacity: 0.55;
}
.c-email { grid-area: email; display: flex; align-items: center; gap: 9px; min-width: 0; }
.c-usage { grid-area: usage; }
.c-exp { grid-area: exp; }
.c-ip { grid-area: ip; font-size: var(--fs-sm); color: var(--text-3); }
.c-act { grid-area: act; display: flex; align-items: center; justify-content: flex-end; gap: 8px; }
.c-email__text {
  min-width: 0;
  display: flex;
  flex-direction: column;
}
.c-email__name {
  font: var(--fw-medium) var(--fs-md) var(--font-mono);
  color: var(--text);
}
.c-email__comment {
  font: var(--fw-regular) var(--fs-xs) var(--font-sans);
  color: var(--text-4);
}
.dot {
  width: 7px;
  height: 7px;
  flex: none;
  border-radius: 50%;
  background: var(--muted-dot);
}
.dot.is-online {
  background: var(--accent);
  box-shadow: 0 0 0 3px var(--accent-ring);
  animation: spx-pulse 2.4s infinite;
}
.health-depleted .c-email__name { color: var(--danger-text); }

@media (max-width: 767px) {
  .clients__head {
    display: none;
  }
  .clients {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }
  .client,
  .with-ip .client {
    grid-template-columns: 1fr auto;
    grid-template-areas:
      'email act'
      'usage exp';
    gap: 10px 12px;
    padding: 11px 12px;
    border: 1px solid var(--border-field);
    border-radius: var(--r-md);
    background: var(--field);
  }
  .client:last-child {
    border-bottom: 1px solid var(--border-field);
  }
  .c-ip {
    display: none;
  }
  .c-exp {
    align-items: flex-end;
  }
}
</style>
