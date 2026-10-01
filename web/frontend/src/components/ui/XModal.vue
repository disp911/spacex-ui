<script setup lang="ts">
import { toRef } from 'vue'
import XIcon from './XIcon.vue'
import { useEscape } from '@/composables/useMedia'
import { useScrollLock } from '@/composables/useScrollLock'
import { t } from '@/i18n'

// Dialog: centered card on desktop, full-screen page on mobile.
// Slots: default (body), toolbar (strip under the header), footer, subtitle.
const open = defineModel<boolean>('open', { default: false })
const props = withDefaults(
  defineProps<{ title: string; subtitle?: string; width?: number; closable?: boolean; flush?: boolean; tall?: boolean }>(),
  { width: 464, closable: true },
)

useScrollLock(toRef(open))
useEscape(() => {
  if (open.value && props.closable) open.value = false
})
</script>

<template>
  <Teleport to="body">
    <Transition name="x-modal">
      <div v-if="open" class="x-modal" @mousedown.self="closable && (open = false)">
        <div
          class="x-modal__dialog"
          :class="{ 'is-tall': tall }"
          role="dialog"
          aria-modal="true"
          :aria-label="title"
          :style="{ '--modal-w': width + 'px' }"
        >
          <header class="x-modal__head">
            <div class="x-modal__titles">
              <h2 class="x-modal__title">
                {{ title }}
                <span v-if="$slots.tag" class="x-modal__tag"><slot name="tag" /></span>
              </h2>
              <div v-if="subtitle || $slots.subtitle" class="x-modal__sub">
                <slot name="subtitle">{{ subtitle }}</slot>
              </div>
            </div>
            <button v-if="closable" type="button" class="x-modal__close" :aria-label="t('common.close')" @click="open = false">
              <XIcon name="close" :size="13" :stroke-width="1.6" />
            </button>
          </header>
          <div v-if="$slots.toolbar" class="x-modal__toolbar"><slot name="toolbar" /></div>
          <div class="x-modal__body" :class="{ 'is-flush': flush }"><slot /></div>
          <footer v-if="$slots.footer" class="x-modal__foot"><slot name="footer" /></footer>
          <slot name="overlay" />
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.x-modal {
  position: fixed;
  inset: 0;
  z-index: 50;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 32px 24px;
  background: var(--scrim);
  backdrop-filter: blur(3px) saturate(0.92);
}
.x-modal__dialog {
  position: relative;
  width: min(var(--modal-w), 100%);
  max-height: 100%;
  display: flex;
  flex-direction: column;
  background: var(--surface-raised);
  border: 1px solid var(--border);
  border-radius: var(--r-lg);
  box-shadow: var(--shadow-dialog);
  backdrop-filter: blur(10px);
  overflow: hidden;
  color: var(--text);
}
.x-modal__dialog.is-tall {
  height: 100%;
}
.x-modal__head {
  flex: none;
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 15px 17px 14px;
  border-bottom: 1px solid var(--divider);
}
.x-modal__titles {
  flex: 1;
  min-width: 0;
}
.x-modal__title {
  margin: 0;
  display: flex;
  align-items: baseline;
  gap: 9px;
  font: var(--fw-semibold) var(--fs-xl) / 1.25 var(--font-sans);
  letter-spacing: var(--tracking-tight);
}
.x-modal__tag {
  font: var(--fw-regular) var(--fs-sm) var(--font-mono);
  color: var(--text-3);
  letter-spacing: 0;
}
.x-modal__sub {
  margin-top: 2px;
  font: var(--fw-regular) var(--fs-sm) var(--font-sans);
  color: var(--text-4);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.x-modal__close {
  flex: none;
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--r-sm);
  border: 1px solid var(--border-control);
  background: var(--control);
  color: var(--text-3);
  transition: border-color var(--dur-base) ease, background var(--dur-base) ease, color var(--dur-base) ease, box-shadow var(--dur-base) ease;
}
.x-modal__close:hover {
  border-color: var(--accent-focus);
  background: var(--accent-soft-hover);
  box-shadow: 0 0 0 3px var(--accent-ring);
  color: var(--text);
}
.x-modal__toolbar {
  flex: none;
  display: flex;
  align-items: center;
  gap: 9px;
  padding: 11px 17px;
  border-bottom: 1px solid var(--divider);
  background: var(--bar-muted);
}
.x-modal__body {
  flex: 1;
  min-height: 0;
  overflow: auto;
  padding: 16px 17px 18px;
  display: flex;
  flex-direction: column;
  gap: 14px;
}
.x-modal__body.is-flush {
  padding: 12px 13px 14px;
  background: var(--bg);
  gap: 0;
}
.x-modal__foot {
  flex: none;
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 9px;
  padding: 12px 17px;
  border-top: 1px solid var(--divider);
  background: var(--bar-muted);
}

/* Mobile: full-screen page */
@media (max-width: 767px) {
  .x-modal {
    padding: 0;
    align-items: stretch;
    background: var(--bg);
    backdrop-filter: none;
  }
  .x-modal__dialog {
    width: 100%;
    height: 100%;
    max-height: none;
    border: 0;
    border-radius: 0;
    background: var(--bg);
    box-shadow: none;
    backdrop-filter: none;
  }
  .x-modal__head {
    padding: calc(16px + env(safe-area-inset-top)) 15px 14px;
  }
  .x-modal__title {
    font-size: var(--fs-xl);
  }
  .x-modal__close {
    width: var(--h-md);
    height: var(--h-md);
    border-radius: var(--r-md);
  }
  .x-modal__toolbar {
    flex-direction: column;
    align-items: stretch;
    gap: 7px;
    padding: 0 15px 12px;
    background: none;
  }
  .x-modal__body {
    padding: 14px 15px 18px;
    gap: 12px;
  }
  .x-modal__body.is-flush {
    padding: 12px 14px 14px;
  }
  .x-modal__foot {
    padding: 12px 16px calc(18px + env(safe-area-inset-bottom));
  }
  .x-modal__foot :deep(.x-btn) {
    flex: 1;
    height: var(--h-lg);
    font-size: var(--fs-lg);
  }
}

/* Transitions */
.x-modal-enter-active,
.x-modal-leave-active {
  transition: opacity 0.2s ease;
}
.x-modal-enter-active .x-modal__dialog {
  transition: transform 0.26s var(--ease-out), opacity 0.26s var(--ease-out);
}
.x-modal-leave-active .x-modal__dialog {
  transition: transform 0.18s ease-in, opacity 0.18s ease-in;
}
.x-modal-enter-from,
.x-modal-leave-to {
  opacity: 0;
}
.x-modal-enter-from .x-modal__dialog,
.x-modal-leave-to .x-modal__dialog {
  opacity: 0;
  transform: translateY(12px) scale(0.97);
}
@media (max-width: 767px) {
  .x-modal-enter-from .x-modal__dialog,
  .x-modal-leave-to .x-modal__dialog {
    transform: translateY(24px);
  }
}
</style>
