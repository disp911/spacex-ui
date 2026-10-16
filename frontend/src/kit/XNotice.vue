<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import XIcon from './XIcon.vue'
import XInlineButton from './XInlineButton.vue'
import type { IconName } from './icons'

/** Inline message inside a page or a dialog. */
const props = withDefaults(
  defineProps<{
    tone?: 'info' | 'ok' | 'warn' | 'err'
    title?: string
    /** A slim one-line notice beside a page title: smaller text and padding. */
    compact?: boolean
    /** Points listed under the text, e.g. what makes the panel unsafe. */
    items?: string[]
    /** A cross at the top right that emits close. */
    closable?: boolean
    /** In place of the tone's own icon. */
    icon?: IconName
    /** The icon on a solid square of the tone and heavier text, e.g. how a window behaves. */
    badge?: boolean
  }>(),
  { tone: 'info' },
)

defineEmits<{ close: [] }>()
const { t } = useI18n()

const glyph = computed<IconName>(() => props.icon ?? (props.tone === 'ok' ? 'check' : props.tone === 'info' ? 'info' : 'alert'))
</script>

<template>
  <div :class="['x-notice', `x-notice--${tone}`, { 'x-notice--compact': compact, 'x-notice--plain': !title, 'x-notice--badge': badge }]" :role="tone === 'err' ? 'alert' : 'status'">
    <span v-if="badge" class="x-notice__badge"><XIcon :name="glyph" :size="14" :stroke="1.6" /></span>
    <XIcon v-else :name="glyph" :size="16" :stroke="1.5" class="x-notice__icon" />
    <div class="x-notice__body">
      <p v-if="title" class="x-notice__title">{{ title }}</p>
      <div class="x-notice__text"><slot /></div>
      <ul v-if="items?.length" class="x-notice__items">
        <li v-for="item in items" :key="item">{{ item }}</li>
      </ul>
    </div>
    <XInlineButton v-if="closable" icon="close" :label="t('kit.close')" class="x-notice__close" @click="$emit('close')" />
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
  color: var(--tone-text, var(--fg2));
}

.x-notice__title {
  margin: 0 0 2px;
  color: var(--tone-t);
  font-weight: 600;
}

.x-notice__items {
  display: flex;
  flex-direction: column;
  gap: 4px;
  margin: 8px 0 0;
  padding: 0;
  list-style: none;
  font-size: var(--fs-12);
}

.x-notice__items li {
  position: relative;
  padding-left: 12px;
  text-wrap: pretty;
}

.x-notice__items li::before {
  content: '';
  position: absolute;
  top: 7px;
  left: 0;
  width: 4px;
  height: 4px;
  border-radius: 50%;
  background: var(--tone);
}

.x-notice--compact {
  align-items: center;
  padding: 6px 6px 6px 12px;
  font-size: var(--fs-12);
  line-height: var(--lh-tight);
}

.x-notice--compact .x-notice__title {
  margin-bottom: 1px;
}

/* Without a title the text is one short message: the icon and actions centre on it. */
.x-notice--plain {
  align-items: center;
}

.x-notice--compact .x-notice__icon,
.x-notice--plain .x-notice__icon {
  margin-top: 0;
}

.x-notice__badge {
  display: flex;
  flex: none;
  align-items: center;
  justify-content: center;
  width: 22px;
  height: 22px;
  border-radius: var(--r-6);
  background: var(--tone);
  color: var(--tone-ink, var(--card));
}

.x-notice--badge .x-notice__body {
  color: var(--fg);
  font-weight: 500;
}

.x-notice__close {
  margin: -4px -4px 0 0;
  color: var(--tone-text, var(--fg2));
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
    /* Under the text: past the icon and the gap. */
    padding-left: calc(var(--icon-l) + var(--sp-10));
  }
}

.x-notice--info {
  --tone: var(--on);
  --tone-t: var(--on-t);
  --tone-bg: var(--on-bg);
  --tone-line: var(--on-line);
}
.x-notice--ok {
  --tone: var(--ok);
  --tone-ink: var(--acc-ink);
  --tone-t: var(--ok-t);
  --tone-bg: var(--ok-bg);
  --tone-line: var(--ok-line);
}
.x-notice--warn {
  --tone: var(--warn);
  --tone-t: var(--warn-t);
  --tone-text: var(--warn-x);
  --tone-bg: var(--warn-bg);
  --tone-line: var(--warn-line);
}
.x-notice--err {
  --tone: var(--err);
  --tone-t: var(--err-t);
  --tone-text: var(--err-t);
  --tone-bg: var(--err-bg);
  --tone-line: var(--err-line);
}
</style>
