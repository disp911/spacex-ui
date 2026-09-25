<script setup lang="ts">
import { onBeforeUnmount, onMounted } from 'vue'
import { isMobile } from '../lib/media'
import { t } from '../lib/i18n'
import Icon from './Icon.vue'

// Dialog shell: centred window on desktop, full-screen sheet on phones.
// Only the body scrolls; header and footer stay put.
const props = withDefaults(
  defineProps<{
    title: string
    subtitle?: string
    /** Short mono note on the title line, e.g. a file name. */
    note?: string
    width?: number
    closable?: boolean
    /** Fixed dialog height (log viewers); otherwise it sizes to content. */
    height?: number
    bodyClass?: string
  }>(),
  { width: 436, closable: true },
)
const emit = defineEmits<{ close: [] }>()

function onKey(e: KeyboardEvent) {
  if (e.key === 'Escape' && props.closable) emit('close')
}
onMounted(() => document.addEventListener('keydown', onKey))
onBeforeUnmount(() => document.removeEventListener('keydown', onKey))
</script>

<template>
  <Teleport to="body">
    <div class="wrap" :class="{ mobile: isMobile }" @click.self="closable && emit('close')">
      <div
        class="dialog"
        role="dialog"
        aria-modal="true"
        :style="
          isMobile
            ? undefined
            : {
                width: `min(${width}px, calc(100vw - 32px))`,
                height: height ? `min(${height}px, calc(100dvh - 32px))` : undefined,
              }
        "
      >
        <div class="head">
          <div class="titles">
            <div class="title">{{ title }}<span v-if="note" class="note">{{ note }}</span></div>
            <div v-if="subtitle || $slots.subtitle" class="subtitle"><slot name="subtitle">{{ subtitle }}</slot></div>
          </div>
          <slot name="head-extra" />
          <button v-if="closable" class="ctl icon close" :title="t('close')" @click="emit('close')">
            <Icon name="close" :size="isMobile ? 16 : 14" />
          </button>
        </div>
        <div v-if="$slots.toolbar" class="toolbar no-scrollbar"><slot name="toolbar" /></div>
        <div class="body" :class="bodyClass"><slot /></div>
        <div v-if="$slots.footer" class="foot"><slot name="footer" /></div>
      </div>
    </div>
  </Teleport>
</template>

<style scoped>
.wrap {
  position: fixed;
  inset: 0;
  z-index: 100;
  display: flex;
  padding: 16px;
  overflow: auto;
  background: rgba(8, 11, 16, 0.55);
  backdrop-filter: blur(3px) saturate(0.92);
  animation: scrimIn 0.22s ease-out;
}
:root[data-theme='light'] .wrap {
  background: rgba(16, 20, 28, 0.32);
}
.dialog {
  margin: auto;
  max-height: calc(100vh - 32px);
  max-height: calc(100dvh - 32px);
  display: flex;
  flex-direction: column;
  background: var(--dialog);
  border: 1px solid var(--line-card);
  border-radius: 16px;
  box-shadow: var(--shadow-card);
  backdrop-filter: blur(10px);
  overflow: hidden;
  animation: cDlgIn 0.26s cubic-bezier(0.22, 1, 0.36, 1);
}
.head {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 15px 17px 14px;
  border-bottom: 1px solid var(--line);
  flex: none;
}
.titles {
  flex: 1;
  min-width: 0;
}
.title {
  font: 600 15.5px var(--font-sans);
  letter-spacing: -0.3px;
}
.note {
  margin-left: 10px;
  font: 400 12px var(--font-mono);
  color: var(--muted);
  letter-spacing: 0;
}
.subtitle {
  font: 400 11.5px var(--font-sans);
  color: var(--muted);
  margin-top: 2px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.close {
  width: 30px;
  height: 30px;
  border-radius: 9px;
  background: var(--chip);
}
.toolbar {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 17px;
  border-bottom: 1px solid var(--line);
  flex: none;
}
.body {
  flex: 1 1 auto;
  min-height: 0;
  overflow: auto;
  padding: 14px 17px 16px;
}
.foot {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 9px;
  padding: 12px 17px;
  border-top: 1px solid var(--line);
  background: rgba(11, 14, 19, 0.4);
  flex: none;
}
:root[data-theme='light'] .foot {
  background: rgba(16, 20, 28, 0.035);
}

/* Phones: full-screen sheet. */
.wrap.mobile {
  padding: 0;
  background: var(--bg);
  backdrop-filter: none;
  animation: fadeIn 0.18s ease-out;
}
.mobile .dialog {
  width: 100%;
  height: 100%;
  max-height: none;
  margin: 0;
  border: 0;
  border-radius: 0;
  box-shadow: none;
  background: var(--bg);
  backdrop-filter: none;
  animation: sheetUp 0.2s ease-out;
}
.mobile .head {
  padding: max(16px, env(safe-area-inset-top)) 16px 14px;
}
.mobile .title {
  font-size: 17px;
}
.mobile .close {
  width: 34px;
  height: 34px;
  border-radius: 11px;
}
.mobile .toolbar {
  padding: 12px 16px;
  overflow-x: auto;
}
.mobile .body {
  padding: 14px 16px 16px;
}
.mobile .foot {
  padding: 12px 16px max(16px, env(safe-area-inset-bottom));
  background: transparent;
}
.mobile .foot :deep(.ctl) {
  flex: 1;
  height: 44px;
  border-radius: 12px;
  font-size: 13.5px;
}
</style>
