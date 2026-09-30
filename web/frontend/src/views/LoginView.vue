<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import XLogo from '@/components/ui/XLogo.vue'
import XField from '@/components/ui/XField.vue'
import XInput from '@/components/ui/XInput.vue'
import XOtpInput from '@/components/ui/XOtpInput.vue'
import XAlert from '@/components/ui/XAlert.vue'
import XIcon from '@/components/ui/XIcon.vue'
import XSpinner from '@/components/ui/XSpinner.vue'
import LangSwitch from '@/components/shell/LangSwitch.vue'
import ThemeToggle from '@/components/shell/ThemeToggle.vue'
import WaveBackdrop from '@/components/shell/WaveBackdrop.vue'
import { auth } from '@/api/panel'
import { env } from '@/env'
import { t } from '@/i18n'

const router = useRouter()
const route = useRoute()

const username = ref('')
const password = ref('')
const otp = ref('')
const twoFactor = ref(false)
const error = ref('')
const phase = ref<'idle' | 'loading' | 'success'>('idle')
const ping = ref<number | null>(null)

const userInput = ref<InstanceType<typeof XInput>>()

onMounted(async () => {
  userInput.value?.focus()
  const started = performance.now()
  try {
    twoFactor.value = await auth.twoFactorEnabled()
    ping.value = Math.round(performance.now() - started)
  } catch {
    ping.value = null
  }
})

const label = computed(() =>
  phase.value === 'loading' ? t('login.checking') : phase.value === 'success' ? t('login.success') : t('login.submit'),
)

async function submit() {
  if (phase.value !== 'idle') return
  if (!username.value || !password.value) {
    error.value = t('login.fillFields')
    return
  }
  error.value = ''
  phase.value = 'loading'
  try {
    const msg = await auth.login(username.value, password.value, twoFactor.value ? otp.value : undefined)
    if (!msg.success) {
      phase.value = 'idle'
      error.value = msg.msg
      return
    }
    phase.value = 'success'
    const next = typeof route.query.next === 'string' && route.query.next.startsWith('/') ? route.query.next : '/'
    setTimeout(() => router.replace(next), 450)
  } catch (e) {
    phase.value = 'idle'
    error.value = (e as Error).message
  }
}
</script>

<template>
  <div class="login">
    <WaveBackdrop />

    <header class="login__bar">
      <XLogo :size="32" class="login__logo" />
      <div class="login__prefs">
        <LangSwitch />
        <ThemeToggle />
      </div>
    </header>

    <main class="login__center">
      <form class="login__card" novalidate @submit.prevent="submit">
        <h1 class="login__title">{{ t('login.title') }}</h1>
        <div class="login__sub">
          {{ t('common.server') }} <span class="mono">{{ env.host }}</span>
        </div>

        <div class="login__fields">
          <XField :label="t('login.username')" for="login-user">
            <XInput
              id="login-user"
              ref="userInput"
              v-model="username"
              autocomplete="username"
              :invalid="!!error"
            />
          </XField>
          <XField :label="t('login.password')" for="login-pass">
            <XInput
              id="login-pass"
              v-model="password"
              type="password"
              autocomplete="current-password"
              :invalid="!!error"
            />
          </XField>
          <XField v-if="twoFactor" :label="t('login.otp')" for="login-otp">
            <XOtpInput id="login-otp" v-model="otp" @complete="submit" />
          </XField>
        </div>

        <XAlert v-if="error && phase === 'idle'" tone="danger" dismissible class="login__error" @dismiss="error = ''">
          {{ error }}
        </XAlert>

        <button type="submit" class="login__submit" :class="`is-${phase}`" :aria-busy="phase === 'loading'">
          <XSpinner v-if="phase === 'loading'" :size="18" />
          <svg v-else-if="phase === 'success'" width="19" height="19" viewBox="0 0 20 20" fill="none" aria-hidden="true">
            <path d="M4.4 10.6 8.1 14.3 15.6 6.3" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round" class="login__check" />
          </svg>
          <span>{{ label }}</span>
        </button>

        <div class="login__foot">
          <span class="login__dot" :class="{ 'is-off': ping === null }" />
          <template v-if="ping !== null">{{ t('login.ping', { ms: ping }) }}</template>
          <template v-else>—</template>
          <span v-if="env.version" class="login__ver mono">
            <XIcon name="xray" :size="12" /> v{{ env.version }}
          </span>
        </div>
      </form>
    </main>
  </div>
</template>

<style scoped>
.login {
  position: relative;
  min-height: 100vh;
  min-height: 100dvh;
  display: flex;
  flex-direction: column;
  padding: 34px 44px;
}
.login__bar {
  position: relative;
  display: flex;
  align-items: center;
  gap: 11px;
  z-index: 2;
}
.login__prefs {
  margin-left: auto;
  display: flex;
  align-items: center;
  gap: 8px;
}
.login__center {
  position: relative;
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px 0 40px;
}
.login__card {
  width: 396px;
  max-width: 100%;
  padding: 30px 32px;
  background: var(--surface-raised);
  border: 1px solid var(--border);
  border-radius: var(--r-lg);
  box-shadow: var(--shadow-dialog);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  animation: spx-pop-in 0.3s var(--ease-out);
}
.login__title {
  margin: 0 0 5px;
  font: var(--fw-semibold) var(--fs-3xl) / 1.15 var(--font-sans);
  letter-spacing: -0.6px;
}
.login__sub {
  margin-bottom: 24px;
  font: var(--fw-regular) var(--fs-md) var(--font-sans);
  color: var(--text-4);
}
.login__sub .mono {
  color: var(--text-2);
}
.login__fields {
  display: flex;
  flex-direction: column;
  gap: 14px;
}
.login__error {
  margin-top: 14px;
}
.login__submit {
  width: 100%;
  height: var(--h-lg);
  margin-top: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 9px;
  border: 0;
  border-radius: var(--r-md);
  background: var(--accent);
  color: var(--on-accent);
  font: var(--fw-semibold) var(--fs-lg) var(--font-sans);
  transition: background 0.2s ease, box-shadow var(--dur-base) ease;
}
.login__submit.is-idle:hover {
  background: var(--accent-hover);
  box-shadow: 0 0 0 4px var(--accent-ring), 0 0 22px var(--accent-glow);
}
.login__submit.is-loading {
  background: var(--accent-press);
  cursor: progress;
}
.login__submit.is-success {
  animation: btn-pop 0.34s ease-out;
}
.login__check {
  stroke-dasharray: 24;
  stroke-dashoffset: 24;
  animation: spx-check-draw 0.34s 0.05s ease-out forwards;
}
@keyframes btn-pop {
  0% { transform: scale(1); }
  40% { transform: scale(0.975); }
  100% { transform: scale(1); }
}
.login__foot {
  margin-top: 20px;
  padding-top: 16px;
  border-top: 1px solid var(--divider);
  display: flex;
  align-items: center;
  gap: 9px;
  font: var(--fw-regular) var(--fs-sm) var(--font-sans);
  color: var(--text-3);
}
.login__dot {
  width: 6px;
  height: 6px;
  flex: none;
  border-radius: 50%;
  background: var(--accent);
  animation: spx-pulse 2.4s infinite 0.8s;
}
.login__dot.is-off {
  background: var(--text-5);
  animation: none;
}
.login__ver {
  margin-left: auto;
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: var(--fs-xs);
  color: var(--text-4);
}

@media (max-width: 767px) {
  .login {
    padding: calc(16px + env(safe-area-inset-top)) 22px calc(30px + env(safe-area-inset-bottom));
  }
  .login__logo :deep(svg) {
    width: 28px;
    height: 28px;
  }
  .login__prefs {
    gap: 7px;
  }
  .login__center {
    align-items: center;
    padding: 20px 0 0;
  }
  .login__card {
    padding: 26px 22px;
    border-radius: var(--r-sheet);
  }
  .login__title {
    font-size: var(--fs-3xl);
  }
  .login__foot {
    border-top: 0;
    padding-top: 0;
    margin-top: 22px;
  }
}
</style>
