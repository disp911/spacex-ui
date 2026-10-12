<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { formatBytes, formatDuration, formatNumber } from '../../i18n/format'
import { XButton, XIcon, XPopover } from '../../kit'
import type { ServiceStatus } from '../../stores/server'
import DashCard from './DashCard.vue'
import DashStat from './DashStat.vue'

/** Xray or the Telegram proxy: state, figures, restart and stop. */
const props = defineProps<{
  title: string
  /** Undefined until the first status arrives. */
  service?: ServiceStatus
  errorTitle: string
  busy?: 'restart' | 'stop' | null
}>()
defineEmits<{ restart: []; stop: []; log: [] }>()
const { t } = useI18n()

const state = computed(() => {
  switch (props.service?.state) {
    case 'running': return { tone: 'ok', label: t('dashboard.running') }
    case 'stop': return { tone: 'warn', label: t('dashboard.stopped') }
    case 'error': return { tone: 'err', label: t('dashboard.error') }
    default: return { tone: 'off', label: t('dashboard.unknown') }
  }
})
const running = computed(() => props.service?.state === 'running')
const errorLines = computed(() => (props.service?.errorMsg ?? '').split('\n').filter(Boolean))
const dash = '—'
</script>

<template>
  <DashCard :title="title">
    <template v-if="service?.version" #head>
      <span class="service__version">{{ service.version }}</span>
    </template>

    <div class="service">
      <div class="service__state">
        <XPopover :disabled="!errorLines.length" :label="errorTitle">
          <template #default="{ props: trigger }">
          <component
            :is="errorLines.length ? 'button' : 'div'"
            v-bind="errorLines.length ? trigger : {}"
            :type="errorLines.length ? 'button' : undefined"
            :class="['service__status', `service__status--${state.tone}`, { 'is-clickable': errorLines.length }]"
            :title="errorLines.length ? t('dashboard.showError') : undefined"
          >
            <span class="service__dot" />
            <span>{{ state.label }}</span>
            <XIcon v-if="errorLines.length" name="caret" :size="14" :stroke="1.6" />
          </component>
          </template>
          <template #content>
            <div class="service__error-head">
              <span>{{ errorTitle }}</span>
              <slot name="error-actions" />
            </div>
            <div class="service__error">
              <div v-for="(line, i) in errorLines" :key="i">{{ line }}</div>
            </div>
          </template>
        </XPopover>
      </div>

      <div class="service__stats">
        <DashStat :label="t('dashboard.online')" :short="t('dashboard.onlineShort')" :value="running ? formatNumber(service!.online) : dash" />
        <DashStat :label="t('dashboard.uptime')" :short="t('dashboard.uptimeShort')" :value="running ? formatDuration(service!.uptime) : dash" />
        <DashStat :label="t('dashboard.memory')" :value="running ? formatBytes(service!.mem) : dash" />
      </div>

      <div class="service__actions">
        <XButton icon="refresh" :label="t('dashboard.restart')" :loading="busy === 'restart'" :disabled="!service || !!busy" @click="$emit('restart')" />
        <XButton
          variant="danger"
          icon="stop"
          :label="t('dashboard.stop')"
          :loading="busy === 'stop'"
          :disabled="!running || !!busy"
          @click="$emit('stop')"
        />
      </div>
    </div>

    <div v-if="$slots.default" class="service__links"><slot /></div>
  </DashCard>
</template>

<style scoped>
.service {
  display: grid;
  grid-template-areas:
    'state state'
    'stats actions';
  grid-template-columns: minmax(0, 1fr) auto;
  align-items: end;
  gap: 12px;
}

.service__state {
  grid-area: state;
  display: flex;
}

.service__stats {
  grid-area: stats;
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 12px;
}

.service__actions {
  grid-area: actions;
  display: flex;
  gap: 8px;
}

.service__version {
  display: inline-flex;
  align-items: center;
  height: 24px;
  padding: 0 10px;
  border: 1px solid var(--line);
  border-radius: var(--r-8);
  background: var(--ctrl);
  color: var(--fg);
  font: 500 var(--fs-13) / 1 var(--font-mono);
}

.service__status {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  padding: 0;
  border: 0;
  border-radius: var(--r-8);
  background: none;
  color: var(--tone);
  font: 600 var(--fs-18) / 1.2 var(--font-sans);
  letter-spacing: -.4px;
  outline: none;
}

.service__status.is-clickable {
  cursor: pointer;
}

.service__status:focus-visible {
  box-shadow: 0 0 0 2px var(--acc-line);
}

.service__status--ok { --tone: var(--ok); }
.service__status--warn { --tone: var(--warn); }
.service__status--err { --tone: var(--err); }
.service__status--off { --tone: var(--dim); }

.service__dot {
  width: 9px;
  height: 9px;
  border-radius: 50%;
  background: var(--tone);
  animation: service-pulse 2.4s infinite;
}

@keyframes service-pulse {
  50% {
    opacity: .35;
  }
}

.service__error-head {
  display: flex;
  align-items: center;
  gap: 4px;
  min-height: 28px;
  margin-bottom: 8px;
  color: var(--fg);
  font: 600 var(--fs-13) / 1.2 var(--font-sans);
}

.service__error-head > span {
  flex: 1;
}

.service__error {
  display: flex;
  flex-direction: column;
  gap: 2px;
  max-height: 200px;
  padding: 9px 11px;
  overflow: auto;
  border: 1px solid var(--err-line);
  border-radius: var(--r-8);
  background: var(--err-bg);
  color: var(--err-t);
  font: 400 var(--fs-11) / 1.55 var(--font-mono);
  word-break: break-word;
}

.service__links {
  display: flex;
  flex-direction: column;
  gap: 2px;
  margin-top: 16px;
  padding-top: 10px;
  border-top: 1px solid var(--div);
}

@media (max-width: 767px) {
  .service {
    grid-template-areas:
      'state actions'
      'stats stats';
    align-items: center;
  }

  .service__version {
    font-size: var(--fs-12);
  }
}
</style>
