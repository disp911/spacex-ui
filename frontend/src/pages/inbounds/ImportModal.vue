<script setup lang="ts">
import { ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { api } from '../../api'
import { toast, XButton, XCodeArea, XModal } from '../../kit'
import { useInbounds } from '../../stores/inbounds'
import { cause } from '../settings/form'

/**
 * Window 18 of the mockup: an inbound pasted as the JSON its data window
 * gives, added as a new one (POST panel/api/inbounds/import), as the old
 * panel's import. Ctrl+S sends it, as there.
 */
const open = defineModel<boolean>('open', { default: false })
const { t } = useI18n()
const store = useInbounds()

const text = ref('')
const saving = ref(false)
watch(open, (on) => on && (text.value = ''))

async function submit() {
  if (saving.value) return
  saving.value = true
  const form = new FormData()
  form.set('data', text.value)
  try {
    await api.post('panel/api/inbounds/import', form)
    open.value = false
    await store.load().catch(() => {})
  } catch (e) {
    toast('err', cause(e))
  } finally {
    saving.value = false
  }
}

function onKeydown(e: KeyboardEvent) {
  if ((e.ctrlKey || e.metaKey) && e.key.toLowerCase() === 's') {
    e.preventDefault()
    void submit()
  }
}
</script>

<template>
  <XModal v-model:open="open" :title="t('inbounds.importInbound')" size="lg" cancel :dismissible="!saving">
    <XCodeArea v-model="text" :height="308" :label="t('inbounds.importInbound')" plain autofocus @keydown="onKeydown" />
    <template #footer>
      <XButton variant="primary" :loading="saving" @click="submit">{{ t('inbounds.import') }}</XButton>
    </template>
  </XModal>
</template>
