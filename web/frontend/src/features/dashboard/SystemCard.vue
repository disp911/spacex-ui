<script setup lang="ts">
import { computed } from 'vue'
import XCard from '@/components/ui/XCard.vue'
import XStat from '@/components/ui/XStat.vue'
import type { ServerStatus } from '@/api/panel'
import { formatUptime } from '@/utils/format'
import { t } from '@/i18n'

const props = defineProps<{ status: ServerStatus | null }>()
const loads = computed(() => (props.status?.loads?.length ? props.status.loads.slice(0, 3).map((l) => l.toFixed(2)) : null))
</script>

<template>
  <XCard :title="t('dash.system')">
    <div class="sys">
      <XStat :label="t('dash.osUptime')" :value="status ? formatUptime(status.uptime) : '—'" class="sys__up" />
      <div class="sys__load">
        <XStat :label="t('dash.load')">
          <template v-if="loads">
            <template v-for="(l, i) in loads" :key="i"><span v-if="i" class="sys__sep"> | </span>{{ l }}</template>
          </template>
          <template v-else>—</template>
        </XStat>
        <div class="sys__hint">{{ t('dash.loadHint') }}</div>
      </div>
    </div>
  </XCard>
</template>

<style scoped>
.sys {
  display: flex;
  align-items: flex-start;
  gap: 24px;
}
.sys__up {
  flex: none;
}
.sys__load {
  flex: 1;
  min-width: 0;
}
.sys__sep {
  color: var(--text-disabled);
}
.sys__hint {
  margin-top: 5px;
  font: var(--fw-regular) var(--fs-xs) var(--font-sans);
  color: var(--text-4);
}
@media (max-width: 767px) {
  .sys {
    gap: 20px;
  }
  .sys__hint {
    display: none;
  }
}
</style>
