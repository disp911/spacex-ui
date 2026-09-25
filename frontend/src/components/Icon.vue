<script setup lang="ts">
// Stroked 16×16 icons copied from the design files; they draw in
// currentColor so they follow the surrounding text colour.
const ICONS = {
  dashboard:
    '<rect x="1.8" y="1.8" width="5.1" height="5.1" rx="1.4" stroke="currentColor" stroke-width="1.4"/><rect x="9.1" y="1.8" width="5.1" height="5.1" rx="1.4" stroke="currentColor" stroke-width="1.4"/><rect x="1.8" y="9.1" width="5.1" height="5.1" rx="1.4" stroke="currentColor" stroke-width="1.4"/><rect x="9.1" y="9.1" width="5.1" height="5.1" rx="1.4" stroke="currentColor" stroke-width="1.4"/>',
  connections:
    '<circle cx="4" cy="4.2" r="2.4" stroke="currentColor" stroke-width="1.4"/><circle cx="12" cy="11.8" r="2.4" stroke="currentColor" stroke-width="1.4"/><path d="M6 5.6c3.2 1 3.9 2.6 4.3 4.3" stroke="currentColor" stroke-width="1.4" stroke-linecap="round"/>',
  settings:
    '<path d="M2 4.6h12M2 11.4h12" stroke="currentColor" stroke-width="1.4" stroke-linecap="round"/><circle cx="10.2" cy="4.6" r="2.1" fill="var(--icon-bg, var(--sidebar))" stroke="currentColor" stroke-width="1.4"/><circle cx="5.8" cy="11.4" r="2.1" fill="var(--icon-bg, var(--sidebar))" stroke="currentColor" stroke-width="1.4"/>',
  xray: '<rect x="2" y="2" width="12" height="12" rx="3.4" stroke="currentColor" stroke-width="1.4"/><circle cx="8" cy="8" r="2.2" stroke="currentColor" stroke-width="1.4"/>',
  logout:
    '<path d="M6.2 2.4H3.4A1.4 1.4 0 0 0 2 3.8v8.4a1.4 1.4 0 0 0 1.4 1.4h2.8" stroke="currentColor" stroke-width="1.4" stroke-linecap="round"/><path d="M10.4 11.2 13.6 8l-3.2-3.2M13.2 8H6.4" stroke="currentColor" stroke-width="1.4" stroke-linecap="round" stroke-linejoin="round"/>',
  globe:
    '<circle cx="8" cy="8" r="6.4" stroke="currentColor" stroke-width="1.4"/><path d="M1.6 8h12.8M8 1.6c3.4 3.6 3.4 9.2 0 12.8M8 1.6C4.6 5.2 4.6 10.8 8 14.4" stroke="currentColor" stroke-width="1.2"/>',
  moon: '<path d="M13.4 10.2A5.8 5.8 0 0 1 5.8 2.6a5.9 5.9 0 1 0 7.6 7.6Z" stroke="currentColor" stroke-width="1.4" stroke-linejoin="round"/>',
  sun: '<circle cx="8" cy="8" r="3.4" stroke="currentColor" stroke-width="1.4"/><path d="M8 1v1.8M8 13.2V15M1 8h1.8M13.2 8H15M3.1 3.1l1.3 1.3M11.6 11.6l1.3 1.3M12.9 3.1l-1.3 1.3M4.4 11.6l-1.3 1.3" stroke="currentColor" stroke-width="1.3" stroke-linecap="round"/>',
  eye: '<path d="M1.5 8S4 3.7 8 3.7 14.5 8 14.5 8 12 12.3 8 12.3 1.5 8 1.5 8Z" stroke="currentColor" stroke-width="1.25" stroke-linejoin="round"/><circle cx="8" cy="8" r="2.1" stroke="currentColor" stroke-width="1.25"/>',
  eyeOff:
    '<path d="M1.5 8S4 3.7 8 3.7 14.5 8 14.5 8 12 12.3 8 12.3 1.5 8 1.5 8Z" stroke="currentColor" stroke-width="1.25" stroke-linejoin="round"/><circle cx="8" cy="8" r="2.1" stroke="currentColor" stroke-width="1.25"/><path d="M2.6 13.4 13.4 2.6" stroke="currentColor" stroke-width="1.3" stroke-linecap="round"/>',
  check: '<path d="M3.4 8.4 6.5 11.5 12.6 4.9" stroke="currentColor" stroke-width="1.9" stroke-linecap="round" stroke-linejoin="round"/>',
  alert:
    '<circle cx="8" cy="8" r="6.9" stroke="currentColor" stroke-width="1.4"/><path d="M8 4.6v4.1" stroke="currentColor" stroke-width="1.6" stroke-linecap="round"/><circle cx="8" cy="11.2" r=".95" fill="currentColor"/>',
  info: '<circle cx="8" cy="8" r="6.9" stroke="currentColor" stroke-width="1.4"/><path d="M8 7.3v4" stroke="currentColor" stroke-width="1.6" stroke-linecap="round"/><circle cx="8" cy="4.8" r=".95" fill="currentColor"/>',
  restart:
    '<path d="M13.4 8a5.4 5.4 0 1 1-1.6-3.8" stroke="currentColor" stroke-width="1.4" stroke-linecap="round"/><path d="M13.6 2.2v3.2h-3.2" stroke="currentColor" stroke-width="1.4" stroke-linecap="round" stroke-linejoin="round"/>',
  stop: '<rect x="4.2" y="4.2" width="7.6" height="7.6" rx="1.6" stroke="currentColor" stroke-width="1.4"/>',
  log: '<path d="M3.4 2.6h9.2v10.8H3.4z" stroke="currentColor" stroke-width="1.4" stroke-linejoin="round"/><path d="M5.8 5.8h4.4M5.8 8h4.4M5.8 10.2h2.6" stroke="currentColor" stroke-width="1.3" stroke-linecap="round"/>',
  version:
    '<path d="M8 2.4v7.2M5.2 7l2.8 2.8L10.8 7" stroke="currentColor" stroke-width="1.4" stroke-linecap="round" stroke-linejoin="round"/><path d="M2.8 12.4h10.4" stroke="currentColor" stroke-width="1.4" stroke-linecap="round"/>',
  code: '<path d="M6 3.4 2.8 8 6 12.6M10 3.4 13.2 8 10 12.6" stroke="currentColor" stroke-width="1.4" stroke-linecap="round" stroke-linejoin="round"/>',
  database:
    '<ellipse cx="8" cy="4" rx="5" ry="2" stroke="currentColor" stroke-width="1.4"/><path d="M3 4v8c0 1.1 2.2 2 5 2s5-.9 5-2V4M3 8c0 1.1 2.2 2 5 2s5-.9 5-2" stroke="currentColor" stroke-width="1.4"/>',
  up: '<path d="M8 13.4V3M8 2.6l4.2 4.2M8 2.6 3.8 6.8" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>',
  down: '<path d="M8 2.6V13M8 13.4l4.2-4.2M8 13.4 3.8 9.2" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>',
  download:
    '<path d="M8 2.8v7.6M4.8 7.4 8 10.6l3.2-3.2M3 13.2h10" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round"/>',
  upload:
    '<path d="M8 10.6V3M4.8 6.2 8 3l3.2 3.2M3 13.2h10" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round"/>',
  arrowUp: '<path d="M8 12.6V3.2M4.6 6.6 8 3.2l3.4 3.4" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round"/>',
  close: '<path d="M4 4l8 8M12 4l-8 8" stroke="currentColor" stroke-width="1.6" stroke-linecap="round"/>',
  search:
    '<circle cx="7.2" cy="7.2" r="4.6" stroke="currentColor" stroke-width="1.4"/><path d="M10.6 10.6 13.4 13.4" stroke="currentColor" stroke-width="1.4" stroke-linecap="round"/>',
  menu: '<path d="M2.5 4.3h11M2.5 8h11M2.5 11.7h11" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>',
  copy: '<rect x="5.3" y="5.3" width="8" height="8" rx="1.8" stroke="currentColor" stroke-width="1.5"/><path d="M10.7 3.6V3.5a1.3 1.3 0 0 0-1.3-1.3H4a1.8 1.8 0 0 0-1.8 1.8v5.4c0 .7.6 1.3 1.3 1.3" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>',
  edit: '<path d="M10.6 2.9 13.1 5.4 5.8 12.7l-3.1.6.6-3.1z" stroke="currentColor" stroke-width="1.4" stroke-linejoin="round"/>',
  trash:
    '<path d="M3.2 4.4h9.6M6.4 4.4V3.1h3.2v1.3M4.6 4.4l.6 8.5h5.6l.6-8.5" stroke="currentColor" stroke-width="1.35" stroke-linecap="round" stroke-linejoin="round"/>',
  plus: '<path d="M8 3.2v9.6M3.2 8h9.6" stroke="currentColor" stroke-width="1.6" stroke-linecap="round"/>',
  lock: '<rect x="3.6" y="7" width="8.8" height="6.4" rx="1.6" stroke="currentColor" stroke-width="1.4"/><path d="M5.6 7V5.2a2.4 2.4 0 0 1 4.8 0V7" stroke="currentColor" stroke-width="1.4"/>',
  chevronLeft: '<path d="M10 3.4 5.4 8l4.6 4.6" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>',
  chevronRight: '<path d="M6 3.4 10.6 8 6 12.6" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>',
  file: '<path d="M4 2h5.4L12.6 5.2V14H4z" stroke="currentColor" stroke-width="1.35" stroke-linejoin="round"/><path d="M9.2 2v3.4h3.4" stroke="currentColor" stroke-width="1.35" stroke-linejoin="round"/>',
  filter: '<path d="M2.4 3.4h11.2L9.4 8.4v4.2l-2.8 1.2V8.4z" stroke="currentColor" stroke-width="1.35" stroke-linejoin="round"/>',
} as const

export type IconName = keyof typeof ICONS

withDefaults(defineProps<{ name: IconName; size?: number }>(), { size: 16 })
</script>

<template>
  <svg :width="size" :height="size" viewBox="0 0 16 16" fill="none" aria-hidden="true" v-html="ICONS[name]" />
</template>
