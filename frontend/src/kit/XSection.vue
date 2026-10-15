<script setup lang="ts">
import { useId } from 'vue'
import XIcon from './XIcon.vue'

/** Collapsible group of settings. Extra attributes go to the header button. */
defineOptions({ inheritAttrs: false })
defineProps<{ title: string; description?: string }>()

const open = defineModel<boolean>('open', { default: true })
const id = useId()
</script>

<template>
  <div :class="['x-section', { 'is-open': open }]">
    <div class="x-section__bar">
      <button
        v-bind="$attrs"
        type="button"
        class="x-section__head"
        :aria-expanded="open"
        :aria-controls="`${id}-body`"
        @click="open = !open"
      >
        <XIcon name="chevron" :size="12" :stroke="1.7" class="x-section__chevron" />
        <span class="x-section__title">{{ title }}</span>
      </button>
      <div v-if="$slots.meta" class="x-section__meta"><slot name="meta" /></div>
    </div>
    <p v-if="description && open" class="x-section__description">{{ description }}</p>
    <div v-show="open" :id="`${id}-body`" class="x-section__body">
      <slot />
    </div>
  </div>
</template>

<style scoped>
.x-section {
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.x-section__bar {
  display: flex;
  align-items: center;
  gap: 8px;
}

.x-section__head {
  display: flex;
  flex: 1;
  min-width: 0;
  align-items: center;
  gap: 8px;
  width: 100%;
  height: 40px;
  padding: 0 10px;
  border: 0;
  border-radius: var(--r-10);
  background: transparent;
  color: var(--fg);
  font: 600 var(--fs-13) / 1 var(--font-sans);
  text-align: left;
  outline: none;
  transition: background-color var(--t-fast);
}

.x-section__head:hover,
.x-section__head.is-hover {
  background: var(--hov);
}

.x-section__head:active,
.x-section__head.is-active {
  background: var(--press);
}

/* Keyboard focus marks the chevron, not the whole header. */
.x-section__head:focus-visible .x-section__chevron,
.x-section__head.is-focus .x-section__chevron {
  border-radius: 2px;
  background: var(--acc-soft);
  color: var(--acc);
  box-shadow: 0 0 0 4px var(--acc-soft), 0 0 0 5px var(--acc-line);
}

.x-section__chevron {
  color: var(--dim);
  transition: transform var(--t-fast) var(--ease), color var(--t-fast), box-shadow var(--t-fast);
}

.is-open > .x-section__bar .x-section__chevron {
  transform: rotate(90deg);
}

.x-section__title {
  flex: 1;
  min-width: 0;
}

.x-section__meta {
  display: flex;
  flex: none;
  align-items: center;
  gap: 8px;
  padding-right: 10px;
}

.x-section__description {
  /* Text lines up with the title: header padding, chevron, gap. */
  margin: 0 var(--sp-10) var(--sp-4) calc(var(--sp-10) + var(--icon-s) + var(--sp-8));
  color: var(--dim);
  font: 400 var(--fs-12) / var(--lh-text) var(--font-sans);
}

.x-section__body {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 8px 10px 12px calc(var(--sp-10) + var(--icon-s) + var(--sp-8));
}

@media (max-width: 767px) {
  .x-section__head {
    height: var(--ctl-h);
    font-size: var(--fs-14);
  }

  .x-section__body {
    padding-left: 10px;
  }
}
</style>
