<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, useId, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { api, ApiError } from '../api'
import { formatMilliseconds } from '../i18n/format'
import { XButton, XIcon, XInput, XNotice, XSpinner, XStatusDot } from '../kit'
import AppLogo from '../layout/AppLogo.vue'
import LanguageMenu from '../layout/LanguageMenu.vue'
import LoginBackdrop from '../layout/LoginBackdrop.vue'
import ThemeToggle from '../layout/ThemeToggle.vue'
import { useSession } from '../stores/session'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const session = useSession()
const id = useId()
const host = location.hostname

const username = ref('')
const password = ref('')
const code = ref('')
const twoFactor = ref(false)

type Problem = 'username' | 'password' | 'code' | 'credentials' | 'connection'
const problem = ref<Problem | null>(null)
const phase = ref<'idle' | 'checking' | 'success'>('idle')

const message = computed(() => {
  switch (problem.value) {
    case 'username': return t('auth.enterUsername')
    case 'password': return t('auth.enterPassword')
    case 'code': return t('auth.enterCode')
    case 'credentials': return t('auth.wrongCredentials')
    case 'connection': return t('auth.noConnection')
    default: return ''
  }
})

// Editing a field takes the error back.
watch([username, password, code], () => (problem.value = null))

const usernameInput = ref<InstanceType<typeof XInput> | null>(null)
const passwordInput = ref<InstanceType<typeof XInput> | null>(null)
const codeInput = ref<InstanceType<typeof XInput> | null>(null)

async function submit() {
  if (phase.value !== 'idle') return
  const empty = !username.value ? 'username' : !password.value ? 'password' : twoFactor.value && !code.value ? 'code' : null
  if (empty) {
    problem.value = empty
    ;({ username: usernameInput, password: passwordInput, code: codeInput })[empty].value?.focus()
    return
  }

  phase.value = 'checking'
  try {
    await session.login(username.value, password.value, code.value)
  } catch (e) {
    phase.value = 'idle'
    // A refusal comes as a normal answer; anything else means the server was not reached.
    problem.value = e instanceof ApiError && e.status === 200 ? 'credentials' : 'connection'
    passwordInput.value?.focus()
    return
  }
  phase.value = 'success'
  const next = typeof route.query.next === 'string' && route.query.next.startsWith('/') ? route.query.next : '/'
  setTimeout(() => router.replace(next), 600)
}

// The 2FA question doubles as the ping: the server answers it without a session.
const ping = ref<number | null>(null)
const pingFailed = ref(false)
let timer = 0

async function check() {
  const started = performance.now()
  try {
    const enabled = await api.post<boolean>('getTwoFactorEnable')
    ping.value = performance.now() - started
    pingFailed.value = false
    twoFactor.value = enabled === true
  } catch (e) {
    // An empty answer means the setting could not be read, not that the server is gone.
    pingFailed.value = !(e instanceof SyntaxError)
    if (!pingFailed.value) ping.value = performance.now() - started
  }
}

const pingTone = computed(() => (pingFailed.value ? 'err' : ping.value == null || ping.value < 150 ? 'ok' : ping.value < 400 ? 'warn' : 'err'))

onMounted(() => {
  // On phones the keyboard would cover the form before the user asked for it.
  if (matchMedia('(min-width: 768px)').matches) usernameInput.value?.focus()
  check()
  timer = window.setInterval(check, 10_000)
})
onBeforeUnmount(() => clearInterval(timer))
</script>

<template>
  <div class="login">
    <LoginBackdrop />

    <header class="login__top">
      <AppLogo large />
      <div class="login__tools">
        <LanguageMenu />
        <ThemeToggle />
      </div>
    </header>

    <main class="login__main">
      <form class="login__card" novalidate @submit.prevent="submit">
        <h1 class="login__title">{{ t('auth.title') }}</h1>
        <i18n-t keypath="auth.server" tag="p" class="login__server">
          <template #host><span class="login__host">{{ host }}</span></template>
        </i18n-t>

        <div class="login__fields">
          <div class="login__field">
            <label class="login__label" :for="`${id}-user`">{{ t('auth.username') }}</label>
            <XInput
              :id="`${id}-user`"
              ref="usernameInput"
              v-model="username"
              name="username"
              autocomplete="username"
              autocapitalize="none"
              spellcheck="false"
              :invalid="problem === 'username' || problem === 'credentials'"
              :aria-describedby="message ? `${id}-error` : undefined"
            />
          </div>
          <div class="login__field">
            <label class="login__label" :for="`${id}-pass`">{{ t('auth.password') }}</label>
            <XInput
              :id="`${id}-pass`"
              ref="passwordInput"
              v-model="password"
              type="password"
              name="password"
              autocomplete="current-password"
              mono
              :invalid="problem === 'password' || problem === 'credentials'"
              :aria-describedby="message ? `${id}-error` : undefined"
            />
          </div>
          <div v-if="twoFactor" class="login__field">
            <label class="login__label" :for="`${id}-code`">{{ t('auth.code') }}</label>
            <XInput
              :id="`${id}-code`"
              ref="codeInput"
              v-model="code"
              name="code"
              :placeholder="t('auth.code')"
              autocomplete="one-time-code"
              inputmode="numeric"
              maxlength="6"
              mono
              :invalid="problem === 'code' || problem === 'credentials'"
              :aria-describedby="message ? `${id}-error` : undefined"
            />
          </div>
        </div>

        <Transition name="login-error">
          <XNotice v-if="message" :id="`${id}-error`" tone="err" class="login__error">
            {{ message }}
            <template #actions>
              <XButton variant="ghost-danger" size="sm" icon="close" :label="t('kit.close')" class="login__dismiss" @click="problem = null" />
            </template>
          </XNotice>
        </Transition>

        <XButton
          type="submit"
          variant="primary"
          block
          :class="['login__submit', { 'is-success': phase === 'success' }]"
          :aria-busy="phase === 'checking' || undefined"
        >
          <template v-if="phase !== 'idle'" #icon>
            <XSpinner v-if="phase === 'checking'" :size="16" />
            <XIcon v-else name="check" :size="16" :stroke="2" class="login__check" />
          </template>
          {{ phase === 'checking' ? t('auth.checking') : phase === 'success' ? t('auth.success') : t('auth.submit') }}
        </XButton>

        <footer class="login__foot" aria-live="polite">
          <XStatusDot :tone="pingTone" pulse />
          <span v-if="pingFailed">{{ t('auth.noConnection') }}</span>
          <span v-else-if="ping != null">{{ t('auth.ping', { value: formatMilliseconds(ping) }) }}</span>
        </footer>
      </form>
    </main>
  </div>
</template>

<style scoped>
.login {
  /* The login form uses the large controls of the mockup on every screen. */
  --ctl-h: 44px;
  --ctl-r: var(--r-12);
  position: relative;
  display: flex;
  flex-direction: column;
  min-height: 100vh;
  min-height: 100dvh;
  padding: 32px 40px;
}

.login__top {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.login__tools {
  display: flex;
  gap: 8px;
  /* The page's controls stay at the regular size. */
  --ctl-h: 34px;
  --ctl-r: var(--r-10);
}

.login__main {
  position: relative;
  display: flex;
  flex: 1;
  align-items: center;
  justify-content: center;
  padding: 24px 0 40px;
}

.login__card {
  display: flex;
  flex-direction: column;
  width: 100%;
  max-width: 396px;
  padding: 32px;
  border: 1px solid var(--line);
  border-radius: var(--r-16);
  background: color-mix(in srgb, var(--card) 86%, transparent);
  box-shadow: var(--shadow-modal);
  backdrop-filter: blur(10px);
}

.login__title {
  margin: 0 0 4px;
  font: 600 var(--fs-20) / var(--lh-tight) var(--font-sans);
  letter-spacing: -.4px;
}

.login__server {
  margin: 0 0 24px;
  color: var(--dim);
  font-size: var(--fs-13);
}

.login__host {
  color: var(--fg2);
  font-family: var(--font-mono);
}

.login__fields {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.login__field {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.login__label {
  color: var(--dim);
  font: 500 var(--fs-11) / 1 var(--font-mono);
  letter-spacing: .07em;
  text-transform: uppercase;
}

.login__error {
  margin-top: 14px;
  align-items: center;
}

.login__dismiss {
  margin: -4px -6px -4px 0;
}

.login__submit {
  margin-top: 16px;
  font-size: var(--fs-14);
}

.login__submit.is-success {
  animation: login-pop .34s ease-out;
}

.login__check {
  stroke-dasharray: 16;
  stroke-dashoffset: 16;
  animation: login-check .3s ease-out forwards;
}

.login__foot {
  display: flex;
  align-items: center;
  gap: 4px;
  min-height: 37px;
  margin-top: 20px;
  padding-top: 16px;
  border-top: 1px solid var(--div);
  color: var(--mut);
  font-size: var(--fs-12);
}

.login__foot :deep(.x-dot) {
  margin-left: -5px;
}

.login-error-enter-active,
.login-error-leave-active {
  transition: opacity var(--t-fast), transform var(--t-fast) var(--ease);
}

.login-error-enter-from,
.login-error-leave-to {
  opacity: 0;
  transform: translateY(-4px);
}

@keyframes login-pop {
  40% {
    transform: scale(.975);
  }
}

@keyframes login-check {
  to {
    stroke-dashoffset: 0;
  }
}

@media (max-width: 767px) {
  .login {
    padding: 16px 20px 24px;
  }

  .login__tools {
    gap: 6px;
  }

  .login__main {
    padding: 24px 0 0;
  }

  .login__card {
    padding: 24px 20px;
  }

  .login__title {
    font-size: var(--fs-24);
    letter-spacing: -.7px;
  }
}
</style>
