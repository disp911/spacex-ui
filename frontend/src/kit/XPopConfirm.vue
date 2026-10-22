<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import XButton from './XButton.vue'
import XIcon from './XIcon.vue'
import XPopover from './XPopover.vue'

/**
 * A short question next to the button that asks it, e.g. reset all traffic:
 * a warning sign, the question, Cancel and the action. The trigger binds the
 * slot props as with XPopover; confirm is emitted on the action. With run the
 * question stays open while the action runs, its button spinning, and closes
 * once it is done.
 */
const props = withDefaults(defineProps<{ text: string; cta: string; danger?: boolean; disabled?: boolean; run?: () => Promise<unknown> }>(), { danger: true })
const emit = defineEmits<{ confirm: [] }>()
const { t } = useI18n()

const running = ref(false)
async function confirm(close: () => void) {
  if (!props.run) {
    close()
    emit('confirm')
    return
  }
  running.value = true
  try {
    await props.run()
  } finally {
    running.value = false
    close()
  }
}
</script>

<template>
  <XPopover :label="text" :disabled="disabled">
    <template #default="slot"><slot v-bind="slot" /></template>
    <template #content="{ close }">
      <div class="x-pop-confirm">
        <p class="x-pop-confirm__text"><XIcon name="warning" :size="16" :stroke="1.5" class="x-pop-confirm__icon" />{{ text }}</p>
        <div class="x-pop-confirm__actions">
          <XButton size="sm" :disabled="running" @click="close">{{ t('kit.cancel') }}</XButton>
          <XButton size="sm" :variant="danger ? 'danger' : 'primary'" :loading="running" autofocus @click="confirm(close)">{{ cta }}</XButton>
        </div>
      </div>
    </template>
  </XPopover>
</template>

<style scoped>
.x-pop-confirm {
  display: flex;
  flex-direction: column;
  gap: var(--sp-12);
  width: 256px;
  max-width: 100%;
}

.x-pop-confirm__text {
  display: flex;
  align-items: flex-start;
  gap: var(--sp-8);
  margin: 0;
  color: var(--fg);
  font: 500 var(--fs-13) / var(--lh-text) var(--font-sans);
  text-wrap: pretty;
}

.x-pop-confirm__icon {
  flex: none;
  margin-top: 2px;
  color: var(--warn);
}

.x-pop-confirm__actions {
  display: flex;
  justify-content: flex-end;
  gap: var(--sp-8);
}
</style>
