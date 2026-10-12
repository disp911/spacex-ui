<script setup lang="ts">
import { computed } from 'vue'
import XIcon from './XIcon.vue'
import type { IconName } from './icons'

/** Inline message inside a page or a dialog. */
const props = withDefaults(defineProps<{ tone?: 'info' | 'ok' | 'warn' | 'err'; title?: string }>(), { tone: 'info' })

const icon = computed<IconName>(() => (props.tone === 'ok' ? 'check' : props.tone === 'info' ? 'info' : 'warning'))
</script>

<template>
  <div :class="['x-notice', `x-notice--${tone}`]" :role="tone === 'err' ? 'alert' : 'status'">
    <XIcon :name="icon" :size="16" :stroke="1.5" class="x-notice__icon" />
    <div class="x-notice__body">
      <p v-if="title" class="x-notice__title">{{ title }}</p>
      <div class="x-notice__text"><slot /></div>
    </div>
    <div v-if="$slots.actions" class="x-notice__actions"><slot name="actions" /></div>
  </div>
</template>

<style scoped>
.x-notice {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  padding: 10px 12px;
  border: 1px solid var(--tone-line);
  border-radius: var(--r-12);
  background: var(--tone-bg);
  color: var(--fg2);
  font: 400 var(--fs-13) / var(--lh-text) var(--font-sans);
}

.x-notice__icon {
  margin-top: 1px;
  color: var(--tone-t);
}

.x-notice__body {
  flex: 1;
  min-width: 0;
}

.x-notice__title {
  margin: 0 0 2px;
  color: var(--tone-t);
  font-weight: 600;
}

.x-notice__actions {
  display: flex;
  flex: none;
  gap: 8px;
}

@media (max-width: 767px) {
  .x-notice {
    flex-wrap: wrap;
  }

  .x-notice__actions {
    flex-basis: 100%;
    padding-left: 26px;
  }
}

.x-notice--info {
  --tone-t: var(--on-t);
  --tone-bg: var(--on-bg);
  --tone-line: var(--on-line);
}
.x-notice--ok {
  --tone-t: var(--ok-t);
  --tone-bg: var(--ok-bg);
  --tone-line: var(--ok-line);
}
.x-notice--warn {
  --tone-t: var(--warn-t);
  --tone-bg: var(--warn-bg);
  --tone-line: var(--warn-line);
}
.x-notice--err {
  --tone-t: var(--err-t);
  --tone-bg: var(--err-bg);
  --tone-line: var(--err-line);
}
</style>
