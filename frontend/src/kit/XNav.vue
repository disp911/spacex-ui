<script setup lang="ts">
import { ref } from 'vue'
import { useRoute } from 'vue-router'
import { useSlidingThumb } from './slidingThumb'

/**
 * A column of XNavItem links, e.g. the side menu. The current page's tint is
 * one thumb that travels between the items like liquid instead of jumping.
 */
defineProps<{ label?: string }>()
const route = useRoute()
const root = ref<HTMLElement | null>(null)
const thumb = useSlidingThumb(root, '.is-current', [() => route.fullPath], 'y')
</script>

<template>
  <nav ref="root" class="x-nav" :aria-label="label">
    <span
      v-if="thumb"
      :class="['x-nav__thumb', thumb.ahead && `x-slide--${thumb.ahead}`]"
      :style="{ top: `${thumb.top}px`, bottom: `${thumb.bottom}px`, left: `${thumb.left}px`, right: `${thumb.right}px` }"
      aria-hidden="true"
    />
    <slot />
  </nav>
</template>

<style scoped>
.x-nav {
  position: relative;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.x-nav__thumb {
  position: absolute;
  border-radius: var(--ctl-r);
  background: var(--acc-soft);
  pointer-events: none;
}
</style>
