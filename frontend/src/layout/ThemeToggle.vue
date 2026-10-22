<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { XButton, XIcon } from '../kit'
import { useTheme } from '../stores/theme'

/** Switches between the dark and the light theme; the icons turn into each other. */
const { t } = useI18n()
const store = useTheme()
</script>

<template>
  <XButton
    :label="store.theme === 'dark' ? t('common.lightTheme') : t('common.darkTheme')"
    @click="store.set(store.theme === 'dark' ? 'light' : 'dark')"
  >
    <template #icon>
      <span class="theme-icons">
        <XIcon name="moon" :size="16" :class="['theme-icons__icon', { 'is-hidden': store.theme !== 'dark' }]" />
        <XIcon name="sun" :size="16" :class="['theme-icons__icon', 'theme-icons__icon--sun', { 'is-hidden': store.theme !== 'light' }]" />
      </span>
    </template>
  </XButton>
</template>

<style scoped>
.theme-icons {
  position: relative;
  width: var(--btn-icon);
  height: var(--btn-icon);
}

.theme-icons__icon {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  transition: transform .42s cubic-bezier(.34, 1.32, .5, 1), opacity .26s ease;
}

.theme-icons__icon.is-hidden {
  opacity: 0;
  transform: rotate(-90deg) scale(.55);
}

.theme-icons__icon--sun.is-hidden {
  transform: rotate(90deg) scale(.55);
}
</style>
