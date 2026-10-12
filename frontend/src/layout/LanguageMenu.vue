<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { locales, setLocale, type Locale } from '../i18n'
import { XButton, XIcon, XMenu, type MenuItem } from '../kit'

/** Language picker; the page behind the open menu is blurred. */
const { t, locale } = useI18n()

const code = (l: string) => l.slice(0, 2).toUpperCase()
const items = computed<MenuItem[]>(() =>
  locales.map((l) => ({ key: l.value, label: l.label, hint: code(l.value), checked: l.value === locale.value })),
)
</script>

<template>
  <XMenu v-slot="{ props, open }" :items="items" :title="t('common.language')" @select="setLocale($event as Locale)">
    <XButton v-bind="props" icon="globe" :label="t('common.language')" class="lang">
      {{ code(locale) }}
      <XIcon name="caret" :size="10" :stroke="1.6" class="lang__caret" />
    </XButton>
    <Teleport to="body">
      <Transition name="lang-scrim">
        <div v-if="open" class="lang-scrim" aria-hidden="true" />
      </Transition>
    </Teleport>
  </XMenu>
</template>

<style scoped>
.x-btn.lang {
  padding: 0 10px 0 12px;
  font-size: var(--fs-12);
}

.lang :deep(.x-btn__text) {
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.lang__caret {
  color: var(--dim);
}

.lang-scrim {
  position: fixed;
  inset: 0;
  z-index: 999;
  background: rgba(8, 11, 16, .14);
  backdrop-filter: blur(3px) saturate(.92);
}

.lang-scrim-enter-active,
.lang-scrim-leave-active {
  transition: opacity var(--t-fast);
}

.lang-scrim-enter-from,
.lang-scrim-leave-to {
  opacity: 0;
}
</style>
