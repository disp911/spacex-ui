<script setup lang="ts">
import { icons, type IconName } from './icons'

withDefaults(defineProps<{ name: IconName; size?: number; stroke?: number }>(), { size: 14, stroke: 1.5 })
</script>

<template>
  <!-- eslint-disable-next-line vue/no-v-html -- static icon markup from icons.ts -->
  <svg
    class="x-icon"
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
.x-icon {
  display: block;
  flex: none;
}
</style>
