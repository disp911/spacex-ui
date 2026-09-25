<script setup lang="ts">
import { onBeforeUnmount, onMounted } from 'vue'
import { cancelConfirm, confirmState as c, runConfirm } from '../lib/confirm'
import { isMobile } from '../lib/media'
import { t } from '../lib/i18n'
import Icon from './Icon.vue'
import Spinner from './Spinner.vue'

function onKey(e: KeyboardEvent) {
  if (!c.open) return
  if (e.key === 'Escape') cancelConfirm()
  if (e.key === 'Enter') runConfirm()
}
onMounted(() => document.addEventListener('keydown', onKey))
onBeforeUnmount(() => document.removeEventListener('keydown', onKey))
</script>

<template>
  <Teleport to="body">
    <div v-if="c.open" class="scrim" :class="{ closing: c.closing, mobile: isMobile }" @click.self="cancelConfirm">
      <div class="box" :class="{ danger: c.tone === 'danger' }" role="alertdialog" aria-modal="true">
        <div class="grab" />
        <div class="main">
          <div class="icon">
            <Icon :name="c.tone === 'danger' ? 'trash' : 'alert'" :size="isMobile ? 18 : 17" />
          </div>
          <div class="texts">
            <div class="title">{{ c.title }}</div>
            <div v-if="c.text" class="text">{{ c.text }}</div>
          </div>
        </div>
        <div class="foot">
          <button class="ctl cancel" :disabled="c.running" @click="cancelConfirm">{{ t('cancel') }}</button>
          <button class="cta" @click="runConfirm">
            <template v-if="c.running">
              <Spinner :size="15" /><span>{{ c.busy || t('loading') }}</span>
            </template>
            <span v-else>{{ c.cta }}</span>
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<style scoped>
.scrim {
  position: fixed;
  inset: 0;
  z-index: 400;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 16px;
  background: rgba(8, 11, 16, 0.55);
  backdrop-filter: blur(3px) saturate(0.92);
  animation: scrimIn 0.22s ease-out;
}
:root[data-theme='light'] .scrim {
  background: rgba(16, 20, 28, 0.32);
}
.scrim.closing {
  animation: scrimOut 0.2s ease-in forwards;
}
.box {
  width: min(400px, 100%);
  background: var(--dialog);
  border: 1px solid var(--line-card);
  border-radius: 16px;
  box-shadow: var(--shadow-card);
  padding: 20px;
  animation: cDlgIn 0.26s cubic-bezier(0.22, 1, 0.36, 1);
}
.closing .box {
  animation: cDlgOut 0.18s ease-in forwards;
}
.grab {
  display: none;
}
.main {
  display: flex;
  gap: 12px;
  margin-bottom: 16px;
}
.icon {
  flex: none;
  width: 32px;
  height: 32px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(var(--warn-rgb), 0.12);
  border: 1px solid rgba(var(--warn-rgb), 0.3);
  color: var(--warn-state);
}
.danger .icon {
  background: rgba(var(--danger-rgb), 0.1);
  border-color: rgba(var(--danger-rgb), 0.32);
  color: var(--danger);
}
.texts {
  flex: 1;
  min-width: 0;
}
.title {
  font: 600 14.5px var(--font-sans);
  letter-spacing: -0.2px;
  margin-bottom: 5px;
}
.text {
  font: 400 12px/1.5 var(--font-sans);
  color: var(--text-3);
}
.foot {
  display: flex;
  justify-content: flex-end;
  gap: 9px;
  margin: 0 -20px -20px;
  padding: 12px 20px;
  border-radius: 0 0 16px 16px;
  border-top: 1px solid var(--line);
  background: rgba(11, 14, 19, 0.4);
}
:root[data-theme='light'] .foot {
  background: rgba(16, 20, 28, 0.035);
}
.cancel {
  background: var(--chip);
  padding: 0 15px;
}
.cta {
  height: 34px;
  padding: 0 15px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  font: 600 12.5px var(--font-sans);
  background: var(--accent);
  color: var(--on-accent);
  transition:
    background 0.2s ease,
    box-shadow 0.22s ease;
}
.cta:hover {
  background: var(--accent-hover);
  box-shadow: var(--btn-glow);
}
.danger .cta {
  background: var(--danger);
  color: #1a0a0b;
}
:root[data-theme='light'] .danger .cta {
  color: #fff;
}
.danger .cta:hover {
  box-shadow: 0 0 0 4px rgba(var(--danger-rgb), 0.16);
}

/* Phones: bottom sheet */
.scrim.mobile {
  align-items: flex-end;
  padding: 0;
}
.mobile .box {
  width: 100%;
  border-radius: 22px 22px 0 0;
  border-bottom: 0;
  padding: 8px 18px max(16px, env(safe-area-inset-bottom));
  animation: cSheetIn 0.34s cubic-bezier(0.22, 1, 0.36, 1);
}
.mobile.closing .box {
  animation: cSheetOut 0.22s cubic-bezier(0.4, 0, 1, 1) forwards;
}
.mobile .grab {
  display: block;
  width: 38px;
  height: 4px;
  border-radius: 2px;
  background: var(--scroll-thumb);
  margin: 2px auto 16px;
}
.mobile .main {
  margin-bottom: 20px;
}
.mobile .foot {
  margin: 0;
  padding: 0;
  border: 0;
  background: none;
}
.mobile .foot button {
  flex: 1;
  height: 44px;
  border-radius: 12px;
  font-size: 13.5px;
}
</style>
