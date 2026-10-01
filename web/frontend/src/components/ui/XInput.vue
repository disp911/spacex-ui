<script setup lang="ts">
import { ref } from 'vue'
import XIcon, { type IconName } from './XIcon.vue'
import { t } from '@/i18n'

const model = defineModel<string>({ default: '' })
const props = withDefaults(
  defineProps<{
    id?: string
    type?: 'text' | 'password' | 'url' | 'search'
    placeholder?: string
    size?: 'sm' | 'md' | 'lg'
    icon?: IconName
    mono?: boolean
    invalid?: boolean
    locked?: boolean
    clearable?: boolean
    autocomplete?: string
    autofocus?: boolean
    inputmode?: 'text' | 'numeric' | 'url' | 'search'
  }>(),
  { type: 'text', size: 'lg' },
)

const reveal = ref(false)
const el = ref<HTMLInputElement>()
defineExpose({ focus: () => el.value?.focus() })
</script>

<template>
  <div class="x-input" :class="[`x-input--${size}`, { 'is-invalid': invalid, 'is-locked': locked, mono }]">
    <XIcon v-if="icon" :name="icon" :size="13" class="x-input__icon" />
    <input
      :id="id"
      ref="el"
      v-model="model"
      :type="type === 'password' && reveal ? 'text' : type"
      :placeholder="placeholder"
      :readonly="locked"
      :aria-invalid="invalid || undefined"
      :autocomplete="autocomplete"
      :autofocus="autofocus"
      :inputmode="inputmode"
      spellcheck="false"
    />
    <button v-if="clearable && model" type="button" class="x-input__action" :aria-label="t('common.close')" @click="model = ''">
      <XIcon name="close" :size="11" :stroke-width="1.6" />
    </button>
    <button
      v-if="props.type === 'password'"
      type="button"
      class="x-input__action"
      :class="{ 'is-on': reveal }"
      :title="reveal ? t('login.hidePassword') : t('login.showPassword')"
      :aria-label="reveal ? t('login.hidePassword') : t('login.showPassword')"
      @click="reveal = !reveal"
    >
      <XIcon name="eye" :size="17" :stroke-width="1.3" />
    </button>
  </div>
</template>

<style scoped>
.x-input {
  --in-h: var(--h-lg);
  display: flex;
  align-items: center;
  gap: 8px;
  height: var(--in-h);
  padding: 0 13px;
  border-radius: var(--r-md);
  background: var(--field);
  border: 1px solid var(--border-field);
  transition: border-color var(--dur-base) ease, box-shadow var(--dur-base) ease;
  min-width: 0;
}
.x-input--md { --in-h: var(--h-md); padding: 0 11px; }
.x-input--sm { --in-h: var(--h-sm); padding: 0 11px; border-radius: var(--r-sm); background: var(--control); border-color: var(--border-control); }
.x-input:hover {
  border-color: var(--accent-focus);
}
.x-input:focus-within {
  border-color: var(--accent);
  box-shadow: 0 0 0 3px var(--accent-ring);
}
.x-input.is-invalid {
  border-color: var(--danger);
}
.x-input.is-invalid:focus-within {
  box-shadow: 0 0 0 3px var(--danger-ring);
}
.x-input.is-locked {
  background: var(--surface-sunken);
  border-color: var(--border-subtle);
  box-shadow: none;
  cursor: not-allowed;
}
.x-input.is-locked input {
  color: var(--text-4);
  cursor: not-allowed;
}
.x-input input {
  flex: 1;
  min-width: 0;
  height: 100%;
  border: 0;
  outline: 0;
  background: transparent;
  color: var(--text);
  font: var(--fw-regular) var(--fs-md) var(--font-sans);
}
.x-input--sm input { font-size: var(--fs-sm); }
.x-input.mono input,
.x-input input[type='password'] {
  font-family: var(--font-mono);
  font-weight: var(--fw-medium);
}
.x-input.mono.x-input--sm input { font-weight: var(--fw-regular); }
.x-input input::placeholder {
  color: var(--text-5);
}
.x-input__icon {
  color: var(--text-4);
}
.x-input__action {
  flex: none;
  display: flex;
  align-items: center;
  padding: 2px;
  margin-right: -2px;
  border: 0;
  background: none;
  color: var(--text-3);
  transition: color var(--dur-fast) ease;
}
.x-input__action:hover,
.x-input__action.is-on {
  color: var(--accent);
}
@media (max-width: 767px) {
  .x-input input { font-size: 16px; } /* prevents iOS zoom on focus */
  .x-input--sm input { font-size: var(--fs-md); }
}
</style>
