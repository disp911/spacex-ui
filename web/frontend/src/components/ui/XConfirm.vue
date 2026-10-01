<script setup lang="ts">
import { toRef } from 'vue'
import XIcon from './XIcon.vue'
import XButton from './XButton.vue'
import { useEscape } from '@/composables/useMedia'
import { useScrollLock } from '@/composables/useScrollLock'
import { t } from '@/i18n'

// Confirmation: small dialog on desktop, bottom sheet on mobile.
const open = defineModel<boolean>('open', { default: false })
const props = withDefaults(
  defineProps<{ title: string; text?: string; cta: string; tone?: 'warning' | 'danger'; busy?: boolean; busyLabel?: string }>(),
  { tone: 'warning' },
)
const emit = defineEmits<{ confirm: [] }>()

useScrollLock(toRef(open))
useEscape(() => {
  if (!props.busy) open.value = false
})
</script>

<template>
  <Teleport to="body">
    <Transition name="x-confirm">
      <div v-if="open" class="x-confirm" @mousedown.self="!busy && (open = false)">
        <div class="x-confirm__box" role="alertdialog" aria-modal="true" :aria-label="title">
          <div class="x-confirm__handle" />
          <div class="x-confirm__main">
            <div class="x-confirm__icon" :class="`is-${tone}`">
              <XIcon :name="tone === 'danger' ? 'trash' : 'alert'" :size="17" />
            </div>
            <div class="x-confirm__texts">
              <div class="x-confirm__title">{{ title }}</div>
              <div v-if="text" class="x-confirm__text">{{ text }}</div>
            </div>
          </div>
          <div class="x-confirm__foot">
            <XButton :disabled="busy" @click="open = false">{{ t('common.cancel') }}</XButton>
            <XButton :variant="tone === 'danger' ? 'danger-solid' : 'primary'" :loading="busy" @click="emit('confirm')">
              {{ busy ? busyLabel || t('common.updating') : cta }}
            </XButton>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.x-confirm {
  position: fixed;
  inset: 0;
  z-index: 70;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px;
  background: var(--scrim);
  backdrop-filter: blur(3px) saturate(0.92);
}
.x-confirm__box {
  width: min(400px, 100%);
  display: flex;
  flex-direction: column;
  background: var(--surface-raised);
  border: 1px solid var(--border);
  border-radius: var(--r-lg);
  box-shadow: var(--shadow-dialog);
  color: var(--text);
  overflow: hidden;
}
.x-confirm__handle {
  display: none;
}
.x-confirm__main {
  display: flex;
  gap: 12px;
  padding: 20px;
}
.x-confirm__icon {
  flex: none;
  width: 32px;
  height: 32px;
  border-radius: var(--r-md);
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px solid;
}
.x-confirm__icon.is-warning {
  background: var(--warning-soft);
  border-color: var(--warning-border);
  color: var(--warning);
}
.x-confirm__icon.is-danger {
  background: var(--danger-soft);
  border-color: var(--danger-border);
  color: var(--danger);
}
.x-confirm__texts {
  flex: 1;
  min-width: 0;
}
.x-confirm__title {
  font: var(--fw-semibold) var(--fs-lg) / 1.35 var(--font-sans);
  letter-spacing: -0.2px;
  margin-bottom: 5px;
  text-wrap: pretty;
}
.x-confirm__text {
  font: var(--fw-regular) var(--fs-sm) / 1.5 var(--font-sans);
  color: var(--text-3);
  text-wrap: pretty;
}
.x-confirm__foot {
  display: flex;
  justify-content: flex-end;
  gap: 9px;
  padding: 12px 20px;
  border-top: 1px solid var(--divider);
  background: var(--bar-muted);
}

@media (max-width: 767px) {
  .x-confirm {
    padding: 0;
    align-items: flex-end;
  }
  .x-confirm__box {
    width: 100%;
    border-radius: var(--r-sheet) var(--r-sheet) 0 0;
    border-width: 1px 0 0;
    box-shadow: var(--shadow-sheet);
  }
  .x-confirm__handle {
    display: block;
    align-self: center;
    width: 38px;
    height: 4px;
    margin-top: 10px;
    border-radius: 2px;
    background: var(--scrollbar);
  }
  .x-confirm__main {
    padding: 14px 16px 16px;
  }
  .x-confirm__title {
    font-size: var(--fs-xl);
  }
  .x-confirm__foot {
    padding: 12px 16px calc(18px + env(safe-area-inset-bottom));
  }
  .x-confirm__foot :deep(.x-btn) {
    flex: 1;
    height: var(--h-lg);
    font-size: var(--fs-lg);
  }
}

.x-confirm-enter-active,
.x-confirm-leave-active {
  transition: opacity 0.2s ease;
}
.x-confirm-enter-active .x-confirm__box {
  transition: transform 0.26s var(--ease-out), opacity 0.26s var(--ease-out);
}
.x-confirm-leave-active .x-confirm__box {
  transition: transform 0.18s ease-in, opacity 0.18s ease-in;
}
.x-confirm-enter-from,
.x-confirm-leave-to {
  opacity: 0;
}
.x-confirm-enter-from .x-confirm__box,
.x-confirm-leave-to .x-confirm__box {
  opacity: 0;
  transform: translateY(12px) scale(0.96);
}
@media (max-width: 767px) {
  .x-confirm-enter-from .x-confirm__box,
  .x-confirm-leave-to .x-confirm__box {
    opacity: 1;
    transform: translateY(100%);
  }
}
</style>
