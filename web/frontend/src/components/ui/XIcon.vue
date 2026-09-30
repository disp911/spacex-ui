<script lang="ts">
// Stroke icons taken from the design mockups, all on a 16×16 grid and
// colored with currentColor. Size defaults to 16px.
export type IconName =
  | 'dashboard' | 'inbounds' | 'xray' | 'settings' | 'logout' | 'menu' | 'close'
  | 'globe' | 'moon' | 'sun' | 'eye' | 'restart' | 'stop' | 'log' | 'download'
  | 'upload' | 'code' | 'database' | 'arrow-up' | 'arrow-down' | 'alert' | 'info'
  | 'check' | 'search' | 'edit' | 'trash' | 'lock' | 'copy' | 'plus' | 'doc'
  | 'chevron-right' | 'chevron-left' | 'chevron-down' | 'filter'
</script>

<script setup lang="ts">
defineProps<{ name: IconName; size?: number | string; strokeWidth?: number }>()
</script>

<template>
  <svg
    class="x-icon"
    :width="size ?? 16"
    :height="size ?? 16"
    viewBox="0 0 16 16"
    fill="none"
    stroke="currentColor"
    :stroke-width="strokeWidth ?? 1.4"
    stroke-linecap="round"
    stroke-linejoin="round"
    aria-hidden="true"
  >
    <template v-if="name === 'dashboard'">
      <rect x="1.8" y="1.8" width="5.1" height="5.1" rx="1.4" /><rect x="9.1" y="1.8" width="5.1" height="5.1" rx="1.4" />
      <rect x="1.8" y="9.1" width="5.1" height="5.1" rx="1.4" /><rect x="9.1" y="9.1" width="5.1" height="5.1" rx="1.4" />
    </template>
    <path v-else-if="name === 'inbounds'" d="M2.6 5.4h10.4M10.4 2.8 13 5.4 10.4 8M13.4 10.6H3M5.6 8 3 10.6l2.6 2.6" />
    <template v-else-if="name === 'xray'">
      <path d="M8 1.8l5.4 3.1v6.2L8 14.2l-5.4-3.1V4.9z" /><circle cx="8" cy="8" r="2" />
    </template>
    <template v-else-if="name === 'settings'">
      <path d="M2 4.6h12M2 11.4h12" />
      <circle cx="10.2" cy="4.6" r="2.1" fill="var(--x-icon-knob, var(--bg-nav))" />
      <circle cx="5.8" cy="11.4" r="2.1" fill="var(--x-icon-knob, var(--bg-nav))" />
    </template>
    <template v-else-if="name === 'logout'">
      <path d="M6.2 2.4H3.4A1.4 1.4 0 0 0 2 3.8v8.4a1.4 1.4 0 0 0 1.4 1.4h2.8" /><path d="M10.4 11.2 13.6 8l-3.2-3.2M13.2 8H6.4" />
    </template>
    <path v-else-if="name === 'menu'" d="M2.4 4h11.2M2.4 8h11.2M2.4 12h11.2" />
    <path v-else-if="name === 'close'" d="M4 4l8 8M12 4l-8 8" />
    <template v-else-if="name === 'globe'">
      <circle cx="8" cy="8" r="6.4" /><path d="M1.6 8h12.8M8 1.6c3.4 3.6 3.4 9.2 0 12.8M8 1.6C4.6 5.2 4.6 10.8 8 14.4" stroke-width="1.2" />
    </template>
    <path v-else-if="name === 'moon'" d="M13.4 10.2A5.8 5.8 0 0 1 5.8 2.6a5.9 5.9 0 1 0 7.6 7.6Z" />
    <template v-else-if="name === 'sun'">
      <circle cx="8" cy="8" r="3.4" />
      <path d="M8 1v1.8M8 13.2V15M1 8h1.8M13.2 8H15M3.1 3.1l1.3 1.3M11.6 11.6l1.3 1.3M12.9 3.1l-1.3 1.3M4.4 11.6l-1.3 1.3" stroke-width="1.3" />
    </template>
    <template v-else-if="name === 'eye'">
      <path d="M1.4 8S3.9 3.7 8 3.7 14.6 8 14.6 8 12.1 12.3 8 12.3 1.4 8 1.4 8Z" /><circle cx="8" cy="8" r="2.1" />
    </template>
    <template v-else-if="name === 'restart'">
      <path d="M13.4 8a5.4 5.4 0 1 1-1.6-3.8" /><path d="M13.6 2.2v3.2h-3.2" />
    </template>
    <rect v-else-if="name === 'stop'" x="4.2" y="4.2" width="7.6" height="7.6" rx="1.6" />
    <template v-else-if="name === 'log'">
      <path d="M3.4 2.6h9.2v10.8H3.4z" /><path d="M5.8 5.8h4.4M5.8 8h4.4M5.8 10.2h2.6" stroke-width="1.3" />
    </template>
    <path v-else-if="name === 'download'" d="M8 2.8v7.6M4.8 7.4 8 10.6l3.2-3.2M3 13.2h10" />
    <path v-else-if="name === 'upload'" d="M8 10.6V3M4.8 6.2 8 3l3.2 3.2M3 13.2h10" />
    <path v-else-if="name === 'code'" d="M6 3.4 2.8 8 6 12.6M10 3.4 13.2 8 10 12.6" />
    <template v-else-if="name === 'database'">
      <ellipse cx="8" cy="4" rx="5" ry="2" /><path d="M3 4v8c0 1.1 2.2 2 5 2s5-.9 5-2V4M3 8c0 1.1 2.2 2 5 2s5-.9 5-2" />
    </template>
    <path v-else-if="name === 'arrow-up'" d="M8 13.4V3M8 2.6l4.2 4.2M8 2.6 3.8 6.8" />
    <path v-else-if="name === 'arrow-down'" d="M8 2.6V13M8 13.4l4.2-4.2M8 13.4 3.8 9.2" />
    <template v-else-if="name === 'alert'">
      <circle cx="8" cy="8" r="6.9" /><path d="M8 4.6v4.1" stroke-width="1.6" /><circle cx="8" cy="11.2" r=".95" fill="currentColor" stroke="none" />
    </template>
    <template v-else-if="name === 'info'">
      <circle cx="8" cy="8" r="6.9" /><path d="M8 7.3v4" stroke-width="1.6" /><circle cx="8" cy="4.8" r=".95" fill="currentColor" stroke="none" />
    </template>
    <path v-else-if="name === 'check'" d="M3.4 8.4 6.5 11.5 12.6 4.9" />
    <template v-else-if="name === 'search'">
      <circle cx="7.2" cy="7.2" r="4.6" /><path d="M10.6 10.6 13.4 13.4" />
    </template>
    <path v-else-if="name === 'edit'" d="M10.6 2.9 13.1 5.4 5.8 12.7l-3.1.6.6-3.1z" />
    <path v-else-if="name === 'trash'" d="M3.2 4.4h9.6M6.4 4.4V3.1h3.2v1.3M4.6 4.4l.6 8.5h5.6l.6-8.5" />
    <template v-else-if="name === 'lock'">
      <rect x="3.2" y="7" width="9.6" height="6.4" rx="1.8" /><path d="M5.6 7V5.2a2.4 2.4 0 0 1 4.8 0V7" />
    </template>
    <template v-else-if="name === 'copy'">
      <rect x="5.3" y="5.3" width="8" height="8" rx="1.8" /><path d="M10.7 3.6V3.5a1.3 1.3 0 0 0-1.3-1.3H4a1.8 1.8 0 0 0-1.8 1.8v5.4c0 .7.6 1.3 1.3 1.3" />
    </template>
    <path v-else-if="name === 'plus'" d="M8 3v10M3 8h10" />
    <template v-else-if="name === 'doc'">
      <rect x="2.2" y="2.6" width="11.6" height="10.8" rx="2.4" /><path d="M5 6.4h6M5 9.2h3.6" />
    </template>
    <path v-else-if="name === 'chevron-right'" d="M6 3.5 10.5 8 6 12.5" />
    <path v-else-if="name === 'chevron-left'" d="M10 3.5 5.5 8l4.5 4.5" />
    <path v-else-if="name === 'chevron-down'" d="M3.5 6 8 10.5 12.5 6" />
    <path v-else-if="name === 'filter'" d="M2.4 3.6h11.2L9.4 8.6v4.2l-2.8-1.4V8.6z" />
  </svg>
</template>

<style scoped>
.x-icon {
  display: block;
  flex: none;
}
</style>
