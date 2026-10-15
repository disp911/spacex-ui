<script setup lang="ts">
import { nextTick, onBeforeUnmount, ref, useId, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import XButton from './XButton.vue'
import { lockScroll } from './scrollLock'

/**
 * Dialog on desktop, bottom sheet on phones. The sheet has a handle instead
 * of the close button, so its footer always starts with Cancel, before the
 * window's own buttons. On desktop Cancel shows only with the cancel prop.
 */
const props = withDefaults(
  defineProps<{
    title: string
    /** wide: up to 1368, for tables such as the connection log. */
    size?: 'xs' | 'sm' | 'md' | 'lg' | 'xl' | 'wide'
    /** Close on Escape and on a click outside. Off while a form is saving. */
    dismissible?: boolean
    /** The body draws its own sections edge to edge: no padding, no gap. */
    flush?: boolean
    /** On phones fills the screen instead of rising as a sheet: logs, editors. */
    full?: boolean
    /** No title bar: the body carries the title, as XConfirm does. */
    headless?: boolean
    /** Cancel in the footer on desktop too, before the window's buttons. */
    cancel?: boolean
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

watch(open, async (on, was) => {
  if (on !== !!was) lockScroll(on)
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
          :aria-labelledby="headless ? undefined : `${id}-title`"
          :aria-label="headless ? title : undefined"
          tabindex="-1"
          :class="['x-modal__dialog', `x-modal__dialog--${size}`, { 'x-modal__dialog--full': full, 'x-modal__dialog--headless': headless }]"
        >
          <header v-if="!headless" :class="['x-modal__head', { 'x-modal__head--sub': $slots.subtitle }]">
            <div class="x-modal__heading">
              <h2 :id="`${id}-title`" class="x-modal__title">{{ title }}</h2>
              <p v-if="$slots.subtitle" class="x-modal__subtitle"><slot name="subtitle" /></p>
            </div>
            <slot name="head" />
            <XButton icon="close" :label="t('kit.close')" :disabled="!dismissible" class="x-modal__close" @click="close" />
          </header>
          <div :class="['x-modal__body', { 'x-modal__body--flush': flush }]"><slot /></div>
          <footer v-if="$slots.footer || cancel || !full" :class="['x-modal__foot', { 'x-modal__foot--phone': !$slots.footer && !cancel }]">
            <XButton v-if="cancel || !full" :disabled="!dismissible" :class="{ 'x-modal__cancel--phone': !cancel }" @click="close">
              {{ t('kit.cancel') }}
            </XButton>
            <slot name="footer" />
          </footer>
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
  /* Room for a scrollbar on both sides: a window taller than the screen does not move when it appears. */
  scrollbar-gutter: stable both-edges;
  background: var(--scrim);
  overscroll-behavior: contain;
}

.x-modal__dialog {
  display: flex;
  flex: none;
  flex-direction: column;
  width: 100%;
  /* In the middle of the screen; a window taller than it starts at the top and scrolls. */
  margin: auto 0;
  border: 1px solid var(--line);
  border-radius: var(--r-16);
  background: var(--card);
  box-shadow: var(--shadow-modal);
  outline: none;
  /* Bars and flush bodies with their own background follow the rounded corners; popups are teleported, so nothing is cut. */
  overflow: clip;
}

.x-modal__dialog--xs { max-width: 400px; }
.x-modal__dialog--sm { max-width: 440px; }
.x-modal__dialog--md { max-width: 560px; }
.x-modal__dialog--lg { max-width: 720px; }
.x-modal__dialog--xl { max-width: 800px; }
.x-modal__dialog--wide { max-width: 1368px; }

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
  background: var(--foot);
}

.x-modal__foot--phone,
.x-modal__cancel--phone {
  display: none;
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
    margin: 0;
    border-bottom: 0;
    border-radius: var(--r-22) var(--r-22) 0 0;
    box-shadow: var(--shadow-sheet);
  }

  .x-modal__body {
    overflow-y: auto;
    overscroll-behavior: contain;
  }

  /* Sheet: a handle on top, the title under it with no line, no close button. */
  .x-modal__dialog:not(.x-modal__dialog--full) {
    padding-top: 10px;
  }

  .x-modal__dialog:not(.x-modal__dialog--full)::before {
    content: '';
    flex: none;
    align-self: center;
    width: 38px;
    height: 4px;
    margin-bottom: 14px;
    border-radius: var(--r-pill);
    background: var(--track);
  }

  .x-modal__dialog:not(.x-modal__dialog--full) .x-modal__head {
    padding: 0 16px;
    border-bottom: 0;
  }

  .x-modal__dialog:not(.x-modal__dialog--full) .x-modal__close {
    display: none;
  }

  .x-modal__dialog:not(.x-modal__dialog--full) .x-modal__body:not(.x-modal__body--flush) {
    padding: 14px 16px 16px;
  }

  /* Without a title bar the body sits right under the handle. */
  .x-modal__dialog--headless:not(.x-modal__dialog--full) .x-modal__body:not(.x-modal__body--flush) {
    padding-top: 0;
  }

  .x-modal__foot--phone,
  .x-modal__cancel--phone {
    display: flex;
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
