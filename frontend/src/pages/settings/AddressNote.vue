<script setup lang="ts">
import { XFieldNote } from '../../kit'

/** Under a changed part of the panel's address: where the panel will be after a restart. */
defineProps<{ address: string }>()
</script>

<template>
  <XFieldNote>
    <i18n-t keypath="settings.afterRestart"><template #address><code>{{ address }}</code></template></i18n-t>
  </XFieldNote>
</template>
