<script setup lang="ts">
import { onBeforeUnmount, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { saveFile } from '../../api'
import { copyText, toast, XButton, XCodeArea, XModal } from '../../kit'

/**
 * Window 17 of the mockup: text to take away, as the old panel's text window
 * showed exports: share links, subscription links, an inbound's data as
 * JSON. It can be changed before it is copied; with a file name it can be
 * saved as that file too.
 */
const props = defineProps<{ title: string; text: string; fileName?: string }>()
const open = defineModel<boolean>('open', { default: false })
const { t } = useI18n()

const content = ref('')
watch(open, (on) => on && (content.value = props.text))

const copied = ref(false)
let reset = 0

async function copy() {
  if (!(await copyText(content.value))) return toast('err', t('kit.copyFailed'))
  copied.value = true
  clearTimeout(reset)
  reset = window.setTimeout(() => (copied.value = false), 1600)
}
onBeforeUnmount(() => clearTimeout(reset))
</script>

<template>
  <XModal v-model:open="open" :title="title" size="lg">
    <XCodeArea v-model="content" :height="308" :label="title" plain />
    <template #footer>
      <XButton v-if="fileName" icon="download" @click="saveFile(content, fileName)">{{ fileName }}</XButton>
      <XButton variant="primary" icon="copy" :success="copied" @click="copy">{{ copied ? t('kit.copied') : t('kit.copy') }}</XButton>
    </template>
  </XModal>
</template>
