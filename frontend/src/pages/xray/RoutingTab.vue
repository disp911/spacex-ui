<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { XButton, XCard, XEmpty, XGridTable, XRouteTarget, XValueChip, type GridColumn, type MenuItem } from '../../kit'
import { useXrayForm } from './form'
import { ruleGroups, ruleTarget, type Rule } from './model'

/** The routing rules in their order: Xray takes the first that matches. */
const { t } = useI18n()
const f = useXrayForm()

const rules = computed<Rule[]>(() => f.draft.config.routing?.rules ?? [])

const columns = computed<GridColumn[]>(() => [
  ...ruleGroups.map((g) => ({ key: g.key, label: t(g.label) })),
  { key: 'target', label: t('xray.ruleOutbound'), size: '.85fr', phone: 'head' as const },
])

/** The parts of a group a rule fills in, each with its values. */
const chips = (r: Rule, key: string) =>
  (ruleGroups.find((g) => g.key === key)?.parts ?? []).map(([label, read]) => ({ label, values: read(r) })).filter((c) => c.values.length)

const blank = (r: Rule, key: string) => (key === 'target' ? !ruleTarget(r) : !chips(r, key).length)

const menu = (): MenuItem[] => [
  { key: 'edit', label: t('xray.edit'), icon: 'edit' },
  { key: 'delete', label: t('xray.delete'), icon: 'delete', danger: true },
]

function move(from: number, to: number) {
  const list = f.draft.config.routing.rules as Rule[]
  list.splice(to, 0, list.splice(from, 1)[0]!)
}

function select(key: string, _row: Rule, index: number) {
  if (key === 'delete') (f.draft.config.routing.rules as Rule[]).splice(index, 1)
}
</script>

<template>
  <div class="xray-tab">
    <div class="xray-tab__tools">
      <XButton variant="primary" icon="plus" class="xray-tab__main">{{ t('xray.ruleAdd') }}</XButton>
      <XButton icon="globe">{{ t('xray.geoFiles') }}</XButton>
    </div>

    <XGridTable
      v-if="rules.length"
      :columns="columns"
      :rows="rules"
      :row-key="(_, i) => i"
      :blank="blank"
      :menu="menu"
      :row-label="(i) => t('xray.ruleN', { n: i + 1 })"
      sortable
      numbered
      @move="move"
      @select="select"
    >
      <template v-for="g in ruleGroups" :key="g.key" #[`cell-${g.key}`]="{ row }">
        <XValueChip v-for="c in chips(row, g.key)" :key="c.label" :label="c.label" :values="c.values" />
      </template>
      <template #cell-target="{ row }">
        <XRouteTarget :tag="ruleTarget(row)!.tag" :balancer="ruleTarget(row)!.kind === 'balancer'" />
      </template>
    </XGridTable>
    <XCard v-else>
      <XEmpty icon="rules" />
    </XCard>
  </div>
</template>

<style scoped>
.xray-tab__tools {
  display: flex;
  gap: 8px;
  margin-bottom: 12px;
}

@media (max-width: 767px) {
  .xray-tab__main {
    flex: 1;
  }
}
</style>
