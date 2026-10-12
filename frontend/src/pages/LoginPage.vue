<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { XNotice } from '../kit'
import AppLogo from '../layout/AppLogo.vue'

/** Login screen; the form is built next. */
const { t } = useI18n()
</script>

<template>
  <main class="login">
    <section class="login__card">
      <AppLogo />
      <h1 class="login__title">{{ t('auth.title') }}</h1>
      <XNotice tone="info">{{ t('common.screenDraft') }}</XNotice>
    </section>
  </main>
</template>

<style scoped>
.login {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  min-height: 100dvh;
  padding: 24px 16px;
}

.login__card {
  display: flex;
  flex-direction: column;
  gap: 16px;
  width: 100%;
  max-width: 400px;
  padding: 24px;
  border: 1px solid var(--line);
  border-radius: var(--r-16);
  background: var(--card);
}

.login__title {
  margin: 0;
  font: 600 var(--fs-20) / var(--lh-tight) var(--font-sans);
  letter-spacing: -.5px;
}
</style>
