<script setup lang="ts">
import { ref } from 'vue'
import XIcon from '@/components/ui/XIcon.vue'
import { LOCALES, locale, t } from '@/i18n'
import { useEscape } from '@/composables/useMedia'

const open = ref(false)
useEscape(() => (open.value = false))
</script>

<template>
  <div class="lang">
    <button type="button" class="lang__trigger" :class="{ 'is-open': open }" :aria-expanded="open" :aria-label="t('common.language')" @click="open = !open">
      <XIcon name="globe" :size="14" />
      <span class="lang__code">{{ locale.toUpperCase() }}</span>
      <XIcon name="chevron-down" :size="10" :stroke-width="1.6" class="lang__caret" />
    </button>
    <div v-if="open" class="lang__scrim" @click="open = false" />
    <Transition name="x-pop">
      <div v-if="open" class="lang__menu" role="listbox">
        <div class="caps lang__title">{{ t('common.language') }}</div>
        <button
          v-for="l in LOCALES"
          :key="l.code"
          type="button"
          role="option"
          :aria-selected="locale === l.code"
          class="lang__opt"
          :class="{ 'is-on': locale === l.code }"
          @click="(locale = l.code), (open = false)"
        >
          <span class="lang__label">{{ l.label }}</span>
          <span class="lang__opt-code">{{ l.code.toUpperCase() }}</span>
          <XIcon name="check" :size="13" :stroke-width="1.9" class="lang__check" />
        </button>
      </div>
    </Transition>
  </div>
</template>

<style scoped>
.lang {
  position: relative;
}
.lang__trigger {
  position: relative;
  z-index: 6;
  display: flex;
  align-items: center;
  gap: 7px;
  height: var(--h-md);
  padding: 0 10px 0 11px;
  border-radius: var(--r-md);
  border: 1px solid var(--border-control);
  background: var(--field);
  color: var(--text-2);
  transition: border-color var(--dur-base) ease, background var(--dur-base) ease, box-shadow var(--dur-base) ease;
}
.lang__trigger:hover,
.lang__trigger.is-open {
  border-color: var(--accent-focus);
  background: var(--accent-soft-hover);
  box-shadow: 0 0 0 3px var(--accent-ring), 0 0 14px var(--accent-glow);
  color: var(--text);
}
.lang__code {
  font: var(--fw-medium) var(--fs-sm) var(--font-sans);
}
.lang__caret {
  color: var(--text-4);
  transition: transform var(--dur-fast) ease;
}
.lang__trigger.is-open .lang__caret {
  transform: rotate(180deg);
}
.lang__scrim {
  position: fixed;
  inset: 0;
  z-index: 4;
  background: var(--scrim-soft);
  backdrop-filter: blur(3px) saturate(0.92);
  animation: spx-fade-in 0.16s ease-out;
}
.lang__menu {
  position: absolute;
  top: calc(100% + 8px);
  right: 0;
  z-index: 7;
  width: 220px;
  padding: 6px;
  display: flex;
  flex-direction: column;
  gap: 1px;
  background: var(--surface-raised);
  border: 1px solid var(--border);
  border-radius: var(--r-md);
  box-shadow: var(--shadow-pop);
  backdrop-filter: blur(10px);
}
.lang__title {
  padding: 7px 9px 6px;
  font-size: var(--fs-2xs);
  color: var(--text-3);
}
.lang__opt {
  display: flex;
  align-items: center;
  gap: 9px;
  padding: 8px 10px;
  border: 0;
  border-radius: var(--r-sm);
  background: transparent;
  color: var(--text-2);
  text-align: left;
}
.lang__opt:hover {
  background: var(--control-hover);
}
.lang__opt.is-on {
  background: var(--accent-soft);
  color: var(--accent);
}
.lang__label {
  flex: 1;
  font: var(--fw-regular) var(--fs-md) var(--font-sans);
}
.lang__opt-code {
  font: var(--fw-medium) var(--fs-2xs) var(--font-mono);
  color: var(--text-3);
}
.lang__check {
  color: transparent;
}
.lang__opt.is-on .lang__check {
  color: var(--accent);
}
.x-pop-enter-active {
  transition: opacity 0.14s ease-out, transform 0.14s ease-out;
}
.x-pop-enter-from {
  opacity: 0;
  transform: translateY(6px);
}
</style>
