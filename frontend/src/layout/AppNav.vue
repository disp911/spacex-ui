<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { XNav, XNavItem } from '../kit'
import { sections } from '../router'
import { useSession } from '../stores/session'

defineEmits<{ navigate: [] }>()
const { t } = useI18n()
const session = useSession()
</script>

<template>
  <XNav :label="t('nav.main')">
    <XNavItem
      v-for="s in sections"
      :key="s.name"
      :to="{ name: s.name }"
      :icon="s.icon"
      :label="t(s.titleKey)"
      @click="$emit('navigate')"
    />
  </XNav>
  <XNavItem icon="logout" :label="t('nav.logout')" danger class="nav__logout" @click="session.logout()" />
</template>

<style scoped>
.nav__logout {
  align-self: end;
}
</style>
