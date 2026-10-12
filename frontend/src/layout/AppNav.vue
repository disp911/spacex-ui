<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { panelUrl } from '../config'
import { XIcon } from '../kit'
import { sections } from '../router'

defineEmits<{ navigate: [] }>()
const { t } = useI18n()
</script>

<template>
  <nav class="nav" :aria-label="t('nav.main')">
    <RouterLink
      v-for="s in sections"
      :key="s.name"
      :to="{ name: s.name }"
      class="nav__item"
      exact-active-class="is-current"
      @click="$emit('navigate')"
    >
      <XIcon :name="s.icon" :size="16" :stroke="1.4" />
      <span>{{ t(s.titleKey) }}</span>
    </RouterLink>
  </nav>
  <a class="nav__item nav__item--logout" :href="panelUrl('logout')">
    <XIcon name="logout" :size="16" :stroke="1.4" />
    <span>{{ t('nav.logout') }}</span>
  </a>
</template>

<style scoped>
.nav {
  display: flex;
  flex-direction: column;
  gap: 2px;
  margin: auto 0;
}

.nav__item {
  display: flex;
  align-items: center;
  gap: 11px;
  height: var(--sidebar-item-h);
  padding: 0 11px;
  border-radius: var(--r-10);
  color: var(--mut);
  font: 400 var(--fs-13) / 1 var(--font-sans);
  white-space: nowrap;
  outline: none;
  transition: background-color var(--t-fast), color var(--t-fast);
}

.nav__item:hover {
  background: var(--hov);
  color: var(--fg);
}

.nav__item:focus-visible {
  box-shadow: inset 0 0 0 2px var(--acc-line);
  color: var(--fg);
}

.nav__item.is-current {
  background: var(--acc-soft);
  color: var(--acc);
  font-weight: 500;
}

.nav__item--logout {
  margin-top: auto;
  color: var(--dim);
}

.nav__item--logout:hover {
  background: var(--err-bg);
  color: var(--err-t);
}

@media (max-width: 767px) {
  .nav__item {
    gap: 12px;
    padding: 0 12px;
    border-radius: var(--r-12);
    font-size: var(--fs-14);
  }
}
</style>
