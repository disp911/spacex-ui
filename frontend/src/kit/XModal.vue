<script setup lang="ts">
import { nextTick, onBeforeUnmount, ref, useId, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import XButton from './XButton.vue'

/** Dialog on desktop, bottom sheet on phones. */
const props = withDefaults(
  defineProps<{
    title: string
    size?: 'sm' | 'md' | 'lg' | 'xl'
    /** Close on Escape and on a click outside. Off while a form is saving. */
    dismissible?: boolean
    /** The body draws its own sections edge to edge: no padding, no gap. */
    flush?: boolean
    /** On phones fills the screen instead of rising as a sheet: logs, editors. */
    full?: boolean
  }>(),
  { size: 'md', dismissible: true },
)

const open = defineModel<boolean>('open', { default: false })
const { t } = useI18n()
const id = useId()
const dialog = ref<HTMLElement | null>(null)
let returnFocus: HTMLElement | null = null

const FOCUSABLE = 'a[href], button:not(:disabled), input:not(:disabled), select:not(:disabled), textarea:not(:disabled), [tabindex]:not([tabindex="-1"])'

function close() {
  if (props.dismissible) open.value = false
}

function onKeydown(e: KeyboardEvent) {
  if (e.key === 'Escape') {
    e.stopPropagation()
    close()
  } else if (e.key === 'Tab' && dialog.value) {
    const list = [...dialog.value.querySelectorAll<HTMLElement>(FOCUSABLE)]
    const first = list[0]
    const last = list[list.length - 1]
    if (!first || !last) return
    if (e.shiftKey && document.activeElement === first) {
      e.preventDefault()
      last.focus()
    } else if (!e.shiftKey && document.activeElement === last) {
      e.preventDefault()
      first.focus()
    }
  }
}

function lockScroll(lock: boolean) {
  document.documentElement.style.overflow = lock ? 'hidden' : ''
}

watch(open, async (on) => {
  lockScroll(on)
  if (on) {
    returnFocus = document.activeElement as HTMLElement | null
    await nextTick()
    const target = dialog.value?.querySelector<HTMLElement>('[autofocus]') ?? dialog.value
    target?.focus()
  } else {
    returnFocus?.focus()
    returnFocus = null
  }
}, { immediate: true })

onBeforeUnmount(() => {
  if (open.value) lockScroll(false)
})
</script>

<template>
  <Teleport to="body">
    <Transition name="x-modal">
      <div v-if="open" :class="['x-modal', { 'x-modal--full': full }]" @keydown="onKeydown" @click.self="close">
        <div
          ref="dialog"
          role="dialog"
          aria-modal="true"
          :aria-labelledby="`${id}-title`"
          tabindex="-1"
          :class="['x-modal__dialog', `x-modal__dialog--${size}`, { 'x-modal__dialog--full': full }]"
        >
          <header :class="['x-modal__head', { 'x-modal__head--sub': $slots.subtitle }]">
            <div class="x-modal__heading">
              <h2 :id="`${id}-title`" class="x-modal__title">{{ title }}</h2>
              <p v-if="$slots.subtitle" class="x-modal__subtitle"><slot name="subtitle" /></p>
            </div>
            <slot name="head" />
            <XButton variant="ghost" size="sm" icon="close" :label="t('kit.close')" :disabled="!dismissible" @click="close" />
          </header>
          <div :class="['x-modal__body', { 'x-modal__body--flush': flush }]"><slot /></div>
          <footer v-if="$slots.footer" class="x-modal__foot"><slot name="footer" /></footer>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.x-modal {
  position: fixed;
  inset: 0;
  z-index: 900;
  display: flex;
  align-items: flex-start;
  justify-content: center;
  padding: 48px 16px;
  overflow-y: auto;
  background: var(--scrim);
  overscroll-behavior: contain;
}

.x-modal__dialog {
  display: flex;
  flex: none;
  flex-direction: column;
  width: 100%;
  border: 1px solid var(--line);
  border-radius: var(--r-16);
  background: var(--card);
  box-shadow: var(--shadow-modal);
  outline: none;
}

.x-modal__dialog--sm { max-width: 440px; }
.x-modal__dialog--md { max-width: 560px; }
.x-modal__dialog--lg { max-width: 720px; }
.x-modal__dialog--xl { max-width: 800px; }

.x-modal__head {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 12px 12px 16px;
  border-bottom: 1px solid var(--div);
}

.x-modal__head--sub {
  padding-block: 14px;
}

.x-modal__heading {
  flex: 1;
  min-width: 0;
}

.x-modal__title {
  margin: 0;
  overflow: hidden;
  color: var(--fg);
  font: 600 var(--fs-16) / var(--lh-tight) var(--font-sans);
  letter-spacing: -.3px;
  white-space: nowrap;
  text-overflow: ellipsis;
}

.x-modal__subtitle {
  margin: 2px 0 0;
  overflow: hidden;
  color: var(--dim);
  font: 400 var(--fs-12) / var(--lh-tight) var(--font-sans);
  white-space: nowrap;
  text-overflow: ellipsis;
}

.x-modal__body {
  display: flex;
  flex-direction: column;
  gap: 14px;
  padding: 16px;
}

.x-modal__body--flush {
  gap: 0;
  padding: 0;
}

.x-modal__foot {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  padding: 12px 16px;
  border-top: 1px solid var(--div);
  border-radius: 0 0 var(--r-16) var(--r-16);
  background: var(--foot);
}

.x-modal-enter-active,
.x-modal-leave-active {
  transition: opacity var(--t-med);
}

.x-modal-enter-active .x-modal__dialog,
.x-modal-leave-active .x-modal__dialog {
  transition: transform var(--t-med) var(--ease);
}

.x-modal-enter-from,
.x-modal-leave-to {
  opacity: 0;
}

.x-modal-enter-from .x-modal__dialog,
.x-modal-leave-to .x-modal__dialog {
  transform: translateY(8px) scale(.98);
}

@media (max-width: 767px) {
  .x-modal {
    align-items: flex-end;
    padding: 24px 0 0;
    overflow: hidden;
  }

  .x-modal__dialog {
    max-width: none;
    max-height: 100%;
    border-bottom: 0;
    border-radius: var(--r-22) var(--r-22) 0 0;
    box-shadow: var(--shadow-sheet);
  }

  .x-modal__body {
    overflow-y: auto;
    overscroll-behavior: contain;
  }

  .x-modal--full {
    padding: 0;
  }

  .x-modal__dialog--full {
    height: 100%;
    border: 0;
    border-radius: 0;
    box-shadow: none;
  }

  .x-modal__dialog--full .x-modal__head {
    padding-top: calc(12px + env(safe-area-inset-top));
  }

  .x-modal__dialog--full .x-modal__head--sub {
    padding-top: calc(14px + env(safe-area-inset-top));
  }

  .x-modal__dialog--full .x-modal__body {
    flex: 1;
    min-height: 0;
  }

  .x-modal__foot {
    padding-bottom: calc(12px + env(safe-area-inset-bottom));
    border-radius: 0;
  }

  .x-modal__foot > :deep(*) {
    flex: 1;
  }

  .x-modal-enter-from .x-modal__dialog,
  .x-modal-leave-to .x-modal__dialog {
    transform: translateY(100%);
  }
}
</style>
