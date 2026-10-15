<script setup lang="ts">
import { reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { Secret } from 'otpauth'
import { api } from '../../api'
import { toast, XButton, XField, XFieldNote, XInput, XSettingBlock, XSettingRow, XSettingsSection, XSwitch } from '../../kit'
import { useSession } from '../../stores/session'
import { cause, useSettingsForm } from './form'
import TwoFactorWindow from './TwoFactorWindow.vue'

/**
 * The admin's credentials and two-factor authentication. Both apply on their
 * own, not with Save: the credentials with Confirm, 2FA once the code is right.
 */
const { t } = useI18n()
const router = useRouter()
const session = useSession()
const f = useSettingsForm()

const open = reactive({ admin: true, twoFactor: true })

// --- Credentials --------------------------------------------------------------

const creds = reactive({ oldUsername: '', oldPassword: '', newUsername: '', newPassword: '' })
/** What went wrong and which fields it points at. */
const failure = ref<{ text: string; current: boolean } | null>(null)
const sending = ref(false)
watch(creds, () => (failure.value = null))

const blank = (s: string) => !s.trim()
const marked = (field: keyof typeof creds) =>
  !!failure.value && (failure.value.current ? field.startsWith('old') : field.startsWith('new') && blank(creds[field]))

function submit() {
  if (sending.value) return
  if (blank(creds.newUsername) || blank(creds.newPassword)) {
    failure.value = { text: t('settings.credentialsEmpty'), current: false }
    return
  }
  // With 2FA on, the change waits for a code from the app.
  if (f.saved.twoFactorEnable) codeFor.value = 'credentials'
  else send()
}

async function send() {
  sending.value = true
  try {
    await api.post('panel/setting/updateUser', { ...creds })
  } catch (e) {
    failure.value = { text: cause(e), current: true }
    sending.value = false
    return
  }
  // The session belongs to the old credentials: the panel asks to log in again.
  await session.logout()
  router.replace({ name: 'login' })
}

// --- Two-factor authentication ------------------------------------------------

/** Which window asks for a code, if any. */
const codeFor = ref<'enable' | 'disable' | 'credentials' | null>(null)
const windowOpen = ref(false)
watch(codeFor, (mode) => (windowOpen.value = !!mode))
watch(windowOpen, (on) => !on && setTimeout(() => !windowOpen.value && (codeFor.value = null), 250))

/** A new token for enabling: kept while the window is open, so the QR code stays the same. */
const newToken = ref('')
const changing = ref(false)

function toggleTwoFactor(on: boolean) {
  if (on) newToken.value = new Secret().base32
  codeFor.value = on ? 'enable' : 'disable'
}

async function applyTwoFactor() {
  const on = codeFor.value === 'enable'
  changing.value = true
  try {
    await f.applyNow({ twoFactorEnable: on, twoFactorToken: on ? newToken.value : '' })
  } catch (e) {
    toast('err', t('settings.twoFactorFailed', { error: cause(e) }))
    throw e
  } finally {
    changing.value = false
  }
}
</script>

<template>
  <div class="settings-tab">
    <XSettingsSection v-model:open="open.admin" :title="t('settings.secAdmin')" :count="1" :note="t('settings.secAdminNote')">
      <XSettingBlock>
        <form class="creds" autocomplete="off" @submit.prevent="submit">
          <XField v-slot="a" :label="t('settings.currentUsername')">
            <XInput v-bind="a" v-model="creds.oldUsername" mono autocomplete="off" :invalid="marked('oldUsername')" />
          </XField>
          <XField v-slot="a" :label="t('settings.currentPassword')">
            <XInput v-bind="a" v-model="creds.oldPassword" type="password" mono autocomplete="off" :invalid="marked('oldPassword')" />
          </XField>
          <XField v-slot="a" :label="t('settings.newUsername')">
            <XInput v-bind="a" v-model="creds.newUsername" mono autocomplete="off" :invalid="marked('newUsername')" />
          </XField>
          <XField v-slot="a" :label="t('settings.newPassword')">
            <XInput v-bind="a" v-model="creds.newPassword" type="password" mono autocomplete="off" :invalid="marked('newPassword')" />
          </XField>
          <!-- Enter in a field confirms, as the button does. -->
          <button type="submit" hidden />
        </form>
        <template #note>
          <XFieldNote v-if="failure" tone="err">{{ failure.text }}</XFieldNote>
          <span>{{ t('settings.credentialsHint') }}</span>
        </template>
        <template #actions>
          <XButton variant="soft" icon="lock" :loading="sending" @click="submit">
            {{ sending ? t('settings.checking') : t('settings.confirm') }}
          </XButton>
        </template>
      </XSettingBlock>
    </XSettingsSection>

    <XSettingsSection v-model:open="open.twoFactor" :title="t('settings.secTwoFactor')" :count="1">
      <XSettingRow :title="t('settings.twoFactor')" :description="t('settings.twoFactorDesc')" inline>
        <XSwitch
          :model-value="f.saved.twoFactorEnable"
          :loading="changing"
          :label="t('settings.twoFactor')"
          :text="changing ? t('settings.twoFactorChanging') : f.saved.twoFactorEnable ? t('settings.twoFactorOn') : t('settings.twoFactorOff')"
          @update:model-value="toggleTwoFactor"
        />
      </XSettingRow>
    </XSettingsSection>
  </div>

  <TwoFactorWindow
    v-if="codeFor"
    v-model:open="windowOpen"
    :mode="codeFor"
    :token="codeFor === 'enable' ? newToken : f.saved.twoFactorToken"
    :apply="codeFor === 'credentials' ? undefined : applyTwoFactor"
    @verified="codeFor === 'credentials' && send()"
  />
</template>

<style scoped>
.settings-tab {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.creds {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 14px 16px;
}

@media (max-width: 767px) {
  .creds {
    grid-template-columns: minmax(0, 1fr);
  }
}
</style>
