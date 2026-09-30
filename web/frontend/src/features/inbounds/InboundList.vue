<script setup lang="ts">
import XSwitch from '@/components/ui/XSwitch.vue'
import XButton from '@/components/ui/XButton.vue'
import XMenu from '@/components/ui/XMenu.vue'
import XTag from '@/components/ui/XTag.vue'
import XIcon from '@/components/ui/XIcon.vue'
import type { MenuItem } from '@/components/ui/types'
import UsageCell from './UsageCell.vue'
import ExpiryCell from './ExpiryCell.vue'
import ClientList from './ClientList.vue'
import type { ClientRow, InboundRow } from './model'
import { panelUrl } from '@/env'
import { t } from '@/i18n'

// Inbound rows with expandable client lists. Desktop: table-like grid.
// Mobile: the same rows reflow into cards (grid areas).
const props = defineProps<{
  rows: InboundRow[]
  expanded: Set<number>
  busy: Set<string>
  showIp?: boolean
  filter?: string
}>()
const emit = defineEmits<{
  expand: [number]
  toggle: [InboundRow, boolean]
  qr: [InboundRow, ClientRow | null]
  action: [string, InboundRow]
  clientToggle: [InboundRow, ClientRow, boolean]
  clientAction: [string, InboundRow, ClientRow]
}>()

const legacy = panelUrl('panel/inbounds')

function menuFor(r: InboundRow): MenuItem[] {
  const items: MenuItem[] = [{ key: 'edit', label: t('ib.edit'), icon: 'edit', href: legacy }]
  if (r.multiUser) {
    items.push(
      { key: 'resetClients', label: t('ib.resetClients'), icon: 'restart' },
      { key: 'delDepleted', label: t('ib.delDepleted'), icon: 'trash', disabled: !r.depleted },
    )
  }
  items.push({ key: 'delete', label: t('ib.delete'), icon: 'trash', danger: true, divided: true })
  return items
}

const canExpand = (r: InboundRow) => r.multiUser
const onRowClick = (r: InboundRow) => canExpand(r) && emit('expand', r.id)
const isOpen = (r: InboundRow) => props.expanded.has(r.id)
</script>

<template>
  <div class="list">
    <div class="list__head caps">
      <span class="i-chev" />
      <span class="i-sw" />
      <span class="i-main">{{ t('ib.colRemark') }}</span>
      <span class="i-port">{{ t('ib.colPort') }}</span>
      <span class="i-cl">{{ t('ib.colClients') }}</span>
      <span class="i-usage">{{ t('ib.colTraffic') }}</span>
      <span class="i-exp">{{ t('ib.colExpiry') }}</span>
      <span class="i-act" />
    </div>

    <div v-for="r in rows" :key="r.id" class="item" :class="[`health-${r.health}`, { 'is-open': isOpen(r) }]">
      <div class="row" :class="{ 'is-clickable': canExpand(r) }" @click="onRowClick(r)">
        <span class="i-chev">
          <XIcon v-if="canExpand(r)" name="chevron-right" :size="13" :stroke-width="1.6" class="chev" />
        </span>
        <span class="i-sw" @click.stop>
          <XSwitch :model-value="r.enable" :loading="busy.has('i' + r.id)" :label="t('ib.enable')" @update:model-value="emit('toggle', r, $event)" />
        </span>
        <div class="i-main">
          <div class="i-remark ellipsis">{{ r.remark || '#' + r.id }}</div>
          <div class="i-tags">
            <XTag tone="accent">{{ r.protocolLabel }}</XTag>
            <XTag v-if="r.network">{{ r.network }}</XTag>
            <XTag v-if="r.security" tone="info">{{ r.security }}</XTag>
          </div>
        </div>
        <span class="i-port mono">{{ r.port }}</span>
        <span class="i-cl">
          <template v-if="r.multiUser">
            <span class="dot" :class="{ 'is-online': r.online > 0 }" />
            <span class="mono"><b>{{ r.online }}</b> / {{ r.clients.length }}</span>
            <span v-if="r.depleted" class="i-badge is-danger" :title="t('ib.attention')">{{ r.depleted }}</span>
            <span v-else-if="r.expiring" class="i-badge is-warn" :title="t('ib.attention')">{{ r.expiring }}</span>
          </template>
          <span v-else class="muted">—</span>
        </span>
        <UsageCell class="i-usage" :up="r.up" :down="r.down" :total="r.total" :health="r.health === 'warn' ? 'ok' : r.health" />
        <ExpiryCell class="i-exp" :expiry-time="r.expiryTime" />
        <span class="i-act" @click.stop>
          <XButton v-if="!r.multiUser && (r.protocol === 'wireguard')" icon="qr" icon-only size="sm" :label="t('ib.qr')" @click="emit('qr', r, null)" />
          <XMenu :items="menuFor(r)" :label="t('ib.more')" :title="r.remark" @select="emit('action', $event, r)" />
        </span>
      </div>

      <div v-if="isOpen(r)" class="item__clients">
        <ClientList
          v-if="r.clients.length"
          :inbound="r"
          :busy="busy"
          :show-ip="showIp"
          :filter="filter"
          @toggle="(c, v) => emit('clientToggle', r, c, v)"
          @qr="(c) => emit('qr', r, c)"
          @action="(k, c) => emit('clientAction', k, r, c)"
        />
        <div v-else class="item__none">{{ t('ib.noClients') }}</div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.list {
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--r-lg);
  overflow: hidden;
}
.list__head,
.row {
  display: grid;
  grid-template-columns: 16px 34px minmax(200px, 2fr) 70px 120px minmax(150px, 1.3fr) 110px 76px;
  grid-template-areas: 'chev sw main port cl usage exp act';
  align-items: center;
  gap: 14px;
  padding: 0 16px;
}
.list__head {
  height: 38px;
  font-size: var(--fs-2xs);
  background: var(--bar-muted);
  border-bottom: 1px solid var(--divider);
}
.row {
  min-height: 60px;
  padding-top: 9px;
  padding-bottom: 9px;
  transition: background var(--dur-fast) ease;
}
.row.is-clickable {
  cursor: pointer;
}
.row:hover {
  background: var(--hover);
}
.item + .item {
  border-top: 1px solid var(--divider);
}
.item.is-open > .row {
  background: var(--control);
}
.item.health-off .i-main,
.item.health-off .i-port,
.item.health-off .i-cl,
.item.health-off .i-usage,
.item.health-off .i-exp {
  opacity: 0.5;
}
.item__clients {
  padding: 4px 0 6px 64px;
  background: var(--bg);
  border-top: 1px solid var(--divider);
  animation: spx-pop-in 0.18s var(--ease-out);
}
.item__none {
  padding: 14px;
  font: var(--fw-regular) var(--fs-sm) var(--font-sans);
  color: var(--text-4);
}

.i-chev { grid-area: chev; display: flex; }
.i-sw { grid-area: sw; display: flex; }
.i-main { grid-area: main; min-width: 0; display: flex; flex-direction: column; gap: 5px; }
.i-port { grid-area: port; }
.i-cl { grid-area: cl; }
.row .i-port { font-size: var(--fs-md); color: var(--text-2); }
.row .i-cl { display: flex; align-items: center; gap: 7px; font-size: var(--fs-sm); color: var(--text-3); white-space: nowrap; }
.i-cl b { color: var(--text); font-weight: var(--fw-medium); }
.i-usage { grid-area: usage; }
.i-exp { grid-area: exp; }
.i-act { grid-area: act; display: flex; justify-content: flex-end; gap: 8px; }

.chev {
  color: var(--text-4);
  transition: transform var(--dur-fast) ease, color var(--dur-fast) ease;
}
.is-open .chev {
  transform: rotate(90deg);
  color: var(--accent);
}
.i-remark {
  font: var(--fw-medium) var(--fs-md) var(--font-sans);
  color: var(--text);
}
.i-tags {
  display: flex;
  gap: 5px;
  flex-wrap: wrap;
}
.dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: var(--muted-dot);
  flex: none;
}
.dot.is-online {
  background: var(--accent);
}
.i-badge {
  min-width: 18px;
  height: 18px;
  padding: 0 5px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--r-pill);
  font: var(--fw-semibold) var(--fs-2xs) var(--font-mono);
}
.i-badge.is-danger { background: var(--danger-soft); color: var(--danger-text); border: 1px solid var(--danger-border); }
.i-badge.is-warn { background: var(--warning-soft); color: var(--warning-title); border: 1px solid var(--warning-border); }
.muted { color: var(--text-5); }

/* Narrower desktops: drop the expiry column into the traffic cell's row. */
@media (max-width: 1180px) and (min-width: 768px) {
  .list__head,
  .row {
    grid-template-columns: 16px 34px minmax(180px, 2fr) 60px 110px minmax(140px, 1.3fr) 76px;
    grid-template-areas: 'chev sw main port cl usage act';
  }
  .i-exp { display: none; }
}

@media (max-width: 767px) {
  .list {
    background: none;
    border: 0;
    border-radius: 0;
    display: flex;
    flex-direction: column;
    gap: 10px;
    overflow: visible;
  }
  .list__head { display: none; }
  .item {
    background: var(--surface);
    border: 1px solid var(--border);
    border-radius: var(--r-lg);
    overflow: hidden;
  }
  .item + .item { border-top: 1px solid var(--border); }
  .row {
    grid-template-columns: 1fr auto auto;
    grid-template-areas:
      'main sw act'
      'cl cl port'
      'usage usage exp';
    gap: 12px 10px;
    padding: 13px 14px 14px;
  }
  .i-chev { display: none; }
  .row .i-port { justify-self: end; font-size: var(--fs-sm); color: var(--text-3); }
  .row .i-port::before { content: ':'; }
  .i-exp { align-items: flex-end; }
  .item.is-open > .row { background: var(--surface); }
  .item__clients {
    padding: 12px;
    border-top: 1px solid var(--divider);
  }
}
</style>
