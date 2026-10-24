<script setup lang="ts">
import { icons, type IconName } from './icons'

withDefaults(defineProps<{ name: IconName; size?: number; stroke?: number }>(), { size: 14, stroke: 1.5 })
</script>

<template>
  <!-- eslint-disable-next-line vue/no-v-html -- static icon markup from icons.ts -->
  <svg
    class="x-icon"
    :style="{ '--x-icon': `${size}px` }"
    :width="size"
    :height="size"
    viewBox="0 0 16 16"
    fill="none"
    stroke="currentColor"
    :stroke-width="stroke"
    stroke-linecap="round"
    stroke-linejoin="round"
    aria-hidden="true"
    v-html="icons[name]"
  />
</template>

<style scoped>
/* Its size grows by --icon-k on phones, to whole pixels; a part that sizes its icons itself (buttons) still wins. */
.x-icon {
  display: block;
  flex: none;
  width: round(calc(var(--x-icon) * var(--icon-k)), 1px);
  height: round(calc(var(--x-icon) * var(--icon-k)), 1px);
}
</style>
