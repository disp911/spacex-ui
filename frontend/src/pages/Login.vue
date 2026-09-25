<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, ref } from 'vue'
import { api } from '../lib/api'
import { config, panelUrl } from '../lib/config'
import { t } from '../lib/i18n'
import Icon from '../components/Icon.vue'
import Logo from '../components/Logo.vue'
import TopControls from '../components/TopControls.vue'
import Spinner from '../components/Spinner.vue'

const username = ref('')
const password = ref('')
const otp = ref('')
const showPassword = ref(false)
const status = ref<'idle' | 'loading' | 'success'>('idle')
const error = ref('')
const otpFocused = ref(false)
const otpInput = ref<HTMLInputElement>()
const userInput = ref<HTMLInputElement>()

const twoFactor = config.twoFactor
const otpCells = computed(() => Array.from({ length: 6 }, (_, i) => otp.value[i] ?? ''))
const activeCell = computed(() => (otpFocused.value ? Math.min(otp.value.length, 5) : -1))

const buttonLabel = computed(() =>
  status.value === 'loading' ? t('spx.checking') : status.value === 'success' ? t('spx.signedIn') : t('spx.signIn'),
)

function onOtpInput(e: Event) {
  const el = e.target as HTMLInputElement
  otp.value = el.value.replace(/\D/g, '').slice(0, 6)
  el.value = otp.value
  if (otp.value.length === 6 && username.value && password.value) submit()
}

async function submit() {
  if (status.value !== 'idle') return
  error.value = ''
  status.value = 'loading'
  const started = performance.now()
  const msg = await api.post('login', {
    username: username.value,
    password: password.value,
    twoFactorCode: otp.value,
  }, { notify: false })
  // Keep the spinner up long enough to read as a state, not a flicker.
  const wait = Math.max(0, 450 - (performance.now() - started))
  setTimeout(() => {
    if (msg.success) {
      status.value = 'success'
      setTimeout(() => location.assign(panelUrl('panel/')), 520)
    } else {
      status.value = 'idle'
      error.value = msg.msg || t('pages.login.toasts.wrongUsernameOrPassword')
      if (twoFactor) otp.value = ''
    }
  }, wait)
}

// Round trip to the server, measured with a tiny request every few seconds.
const ping = ref<number | null>(null)
let pingTimer: number | undefined
async function measure() {
  const t0 = performance.now()
  try {
    await fetch(panelUrl('getTwoFactorEnable'), {
      method: 'POST',
      headers: { 'X-Requested-With': 'XMLHttpRequest' },
      cache: 'no-store',
    })
    ping.value = Math.round(performance.now() - t0)
  } catch {
    ping.value = null
  }
}

onMounted(async () => {
  measure()
  pingTimer = window.setInterval(measure, 5000)
  await nextTick()
  userInput.value?.focus()
})
onBeforeUnmount(() => clearInterval(pingTimer))
</script>

<template>
  <div class="login">
    <div class="glow" />
    <div class="waves" aria-hidden="true">
      <div class="wave a">
        <svg viewBox="0 0 2880 340" preserveAspectRatio="none">
          <path d="M0 196 C 240 138 480 254 720 196 C 960 138 1200 254 1440 196 C 1680 138 1920 254 2160 196 C 2400 138 2640 254 2880 196 L2880 340 L0 340 Z" />
        </svg>
      </div>
      <div class="wave b">
        <svg viewBox="0 0 2880 340" preserveAspectRatio="none">
          <path d="M0 232 C 120 190 360 274 480 232 C 600 190 840 274 960 232 C 1080 190 1320 274 1440 232 C 1560 190 1800 274 1920 232 C 2040 190 2280 274 2400 232 C 2520 190 2760 274 2880 232 L2880 340 L0 340 Z" />
        </svg>
      </div>
      <div class="wave c">
        <svg viewBox="0 0 2880 340" preserveAspectRatio="none">
          <path d="M0 268 C 90 242 270 294 360 268 C 450 242 630 294 720 268 C 810 242 990 294 1080 268 C 1170 242 1350 294 1440 268 C 1530 242 1710 294 1800 268 C 1890 242 2070 294 2160 268 C 2250 242 2430 294 2520 268 C 2610 242 2790 294 2880 268 L2880 340 L0 340 Z" />
        </svg>
      </div>
    </div>

    <div class="frame">
      <header class="top">
        <Logo />
        <TopControls class="controls" />
      </header>

      <main class="center">
        <form class="card" novalidate @submit.prevent="submit">
          <h1 class="title">{{ t('spx.loginTitle') }}</h1>
          <div class="sub">{{ t('spx.server') }} <span class="host">{{ config.host }}</span></div>

          <div class="fields">
            <label class="field">
              <span class="caps">{{ t('spx.loginLabel') }}</span>
              <span class="input" :class="{ err: error }">
                <input
                  ref="userInput"
                  v-model.trim="username"
                  name="username"
                  autocomplete="username"
                  autocapitalize="off"
                  spellcheck="false"
                  required
                />
              </span>
            </label>
            <label class="field">
              <span class="caps">{{ t('password') }}</span>
              <span class="input pass" :class="{ err: error, ring: error }">
                <input
                  v-model="password"
                  name="password"
                  :type="showPassword ? 'text' : 'password'"
                  autocomplete="current-password"
                  required
                />
                <button
                  type="button"
                  class="eye"
                  :title="t(showPassword ? 'spx.hidePassword' : 'spx.showPassword')"
                  @click="showPassword = !showPassword"
                >
                  <Icon :name="showPassword ? 'eyeOff' : 'eye'" :size="17" />
                </button>
              </span>
            </label>
            <div v-if="twoFactor" class="field">
              <span class="caps">{{ t('spx.twoFactorCode') }}</span>
              <div class="otp" @click="otpInput?.focus()">
                <input
                  ref="otpInput"
                  class="otp-input"
                  :value="otp"
                  inputmode="numeric"
                  autocomplete="one-time-code"
                  maxlength="6"
                  :aria-label="t('spx.twoFactorCode')"
                  @input="onOtpInput"
                  @focus="otpFocused = true"
                  @blur="otpFocused = false"
                />
                <div
                  v-for="(d, i) in otpCells"
                  :key="i"
                  class="cell"
                  :class="{ filled: d, active: i === activeCell }"
                >
                  {{ d || '' }}
                </div>
              </div>
            </div>
          </div>

          <div v-if="error && status === 'idle'" class="banner" role="alert">
            <Icon name="alert" />
            <div class="banner-text">{{ error }}</div>
            <button type="button" class="banner-close" :title="t('close')" @click="error = ''">✕</button>
          </div>

          <button type="submit" class="submit" :class="status">
            <Spinner v-if="status === 'loading'" :size="18" on-accent />
            <svg v-if="status === 'success'" width="19" height="19" viewBox="0 0 20 20" fill="none">
              <path class="tick" d="M4.4 10.6 8.1 14.3 15.6 6.3" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round" />
            </svg>
            <span>{{ buttonLabel }}</span>
          </button>

          <div class="foot">
            <span class="pdot" :class="{ off: ping === null }" />
            {{ ping === null ? t('spx.pingUnavailable') : t('spx.ping', { ms: ping }) }}
          </div>
        </form>
      </main>
    </div>
  </div>
</template>

<style scoped>
.login {
  position: relative;
  min-height: 100%;
  overflow: hidden;
  display: flex;
  background: var(--bg);
}
.glow {
  position: absolute;
  inset: 0;
  background: radial-gradient(120% 90% at 50% 118%, var(--login-glow) 0%, transparent 62%);
  pointer-events: none;
}
.waves {
  position: absolute;
  left: 0;
  right: 0;
  bottom: 0;
  height: 440px;
  overflow: hidden;
  pointer-events: none;
}
.wave {
  position: absolute;
  left: 0;
  width: 200%;
  height: 100%;
  will-change: transform;
}
.wave svg {
  width: 100%;
  height: 100%;
}
.wave.a {
  bottom: -6px;
  opacity: var(--wave-a-opacity);
  animation: waveA 26s linear infinite;
  fill: var(--wave-a);
}
.wave.b {
  bottom: -14px;
  opacity: var(--wave-b-opacity);
  animation: waveB 17s linear infinite;
  fill: var(--wave-b);
}
.wave.c {
  bottom: -20px;
  opacity: var(--wave-c-opacity);
  animation: waveA 11s linear infinite;
  fill: var(--wave-c);
}
@keyframes waveA {
  from {
    transform: translate3d(0, 0, 0);
  }
  to {
    transform: translate3d(-50%, 0, 0);
  }
}
@keyframes waveB {
  from {
    transform: translate3d(-50%, 0, 0);
  }
  to {
    transform: translate3d(0, 0, 0);
  }
}

.frame {
  position: relative;
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  padding: 34px 44px;
  min-height: 100vh;
  min-height: 100dvh;
}
.top {
  display: flex;
  align-items: center;
  gap: 11px;
}
.controls {
  margin-left: auto;
}
.center {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px 0 40px;
}
.card {
  width: 100%;
  max-width: 396px;
  background: var(--card);
  border: 1px solid var(--line-card);
  border-radius: 16px;
  backdrop-filter: blur(10px);
  box-shadow: var(--shadow-card);
  padding: 30px 32px;
}
.title {
  margin: 0 0 5px;
  font: 600 19px var(--font-sans);
  letter-spacing: -0.4px;
}
.sub {
  font: 400 12.5px var(--font-sans);
  color: var(--muted);
  margin-bottom: 24px;
}
.host {
  font-family: var(--font-mono);
  color: var(--text-2);
}
.fields {
  display: flex;
  flex-direction: column;
  gap: 14px;
}
.field {
  display: block;
}
.field .caps {
  display: block;
  margin-bottom: 7px;
}
.input {
  height: 44px;
  border-radius: 12px;
  background: var(--fill);
  border: 1px solid var(--border);
  display: flex;
  align-items: center;
  padding: 0 13px;
  transition:
    box-shadow 0.22s ease,
    border-color 0.22s ease;
}
.input:hover,
.input:focus-within {
  border-color: var(--glow-border);
  box-shadow: var(--glow);
}
.input input {
  flex: 1;
  min-width: 0;
  height: 100%;
  border: 0;
  outline: 0;
  background: transparent;
  font: 450 13px var(--font-sans);
  color: var(--text-2);
}
.input.pass input {
  font: 500 13px var(--font-mono);
  color: var(--text);
}
.input.pass input[type='password'] {
  letter-spacing: 0.12em;
}
.input.err {
  border-color: var(--danger);
}
.input.ring {
  box-shadow: 0 0 0 3px rgba(var(--danger-rgb), 0.13);
}
.eye {
  display: flex;
  align-items: center;
  color: var(--text-3b);
  margin-left: 8px;
  transition: color 0.18s ease;
}
.eye:hover {
  color: var(--accent);
}
.otp {
  position: relative;
  display: flex;
  gap: 7px;
}
.otp-input {
  position: absolute;
  inset: 0;
  width: 100%;
  opacity: 0;
  border: 0;
  caret-color: transparent;
  font-size: 16px;
  cursor: text;
}
.cell {
  flex: 1;
  height: 44px;
  border-radius: 12px;
  background: var(--fill);
  border: 1px solid var(--border-soft);
  display: flex;
  align-items: center;
  justify-content: center;
  font: 500 15px var(--font-mono);
  color: var(--text);
  pointer-events: none;
  transition:
    border-color 0.18s ease,
    box-shadow 0.18s ease;
}
.cell.filled {
  border-color: var(--border);
}
.cell.active {
  border-color: var(--accent);
  box-shadow: var(--focus-ring);
}
.banner {
  margin-top: 14px;
  display: flex;
  gap: 10px;
  align-items: flex-start;
  background: rgba(var(--danger-rgb), 0.11);
  border: 1px solid rgba(var(--danger-rgb), 0.32);
  border-radius: 10px;
  padding: 11px 12px;
  color: var(--danger-text);
  animation: sheetUp 0.16s ease-out;
}
.banner svg {
  margin-top: 1px;
}
.banner-text {
  flex: 1;
  min-width: 0;
  font: 500 12.5px var(--font-sans);
}
.banner-close {
  flex: none;
  font: 400 12px var(--font-mono);
  color: var(--danger-text);
  padding: 0 2px;
}
.submit {
  width: 100%;
  height: 44px;
  margin-top: 16px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 9px;
  font: 600 13.5px var(--font-sans);
  background: var(--accent);
  color: var(--on-accent);
  transition:
    background 0.2s ease,
    box-shadow 0.22s ease;
}
.submit:hover {
  background: var(--accent-hover);
  box-shadow: var(--btn-glow);
}
.submit.loading {
  background: var(--accent-press);
  cursor: progress;
}
.submit.success {
  animation: btnPop 0.34s ease-out;
}
.tick {
  stroke-dasharray: 24;
  stroke-dashoffset: 24;
  animation: checkDraw 0.34s 0.05s ease-out forwards;
}
.foot {
  margin-top: 20px;
  padding-top: 16px;
  border-top: 1px solid var(--line);
  display: flex;
  align-items: center;
  gap: 9px;
  font: 450 12px var(--font-sans);
  color: var(--text-3);
}
.pdot {
  width: 6px;
  height: 6px;
  flex: none;
  border-radius: 3px;
  background: var(--accent);
  animation: pulseDot 2.4s infinite 0.8s;
}
.pdot.off {
  background: var(--unknown);
}

/* Phones: no card, the form sits on the background. */
@media (max-width: 760px) {
  .frame {
    padding: max(22px, env(safe-area-inset-top)) 22px 22px;
  }
  .top :deep(.logo svg) {
    width: 28px;
    height: 28px;
  }
  .top :deep(.logo .name) {
    font-size: 14.5px;
  }
  .center {
    align-items: flex-start;
    padding: 12vh 24px 180px;
  }
  .card {
    background: none;
    border: 0;
    box-shadow: none;
    backdrop-filter: none;
    padding: 0;
  }
  .title {
    font-size: 24px;
    letter-spacing: -0.6px;
  }
  .input,
  .cell {
    height: 46px;
  }
  .input {
    padding: 0 14px;
  }
  .input input {
    font-size: 13.5px;
  }
  .submit {
    height: 46px;
  }
  .foot {
    border-top: 0;
    padding-top: 0;
  }
  .waves {
    height: 200px;
  }
}
</style>
