<script setup lang="ts">
import { nextTick, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'
import { XButton } from '../kit'
import AppLogo from './AppLogo.vue'
import AppNav from './AppNav.vue'

const { t } = useI18n()
const route = useRoute()
const host = location.hostname

const drawerOpen = ref(false)
const drawer = ref<HTMLElement | null>(null)
const menuButton = ref<InstanceType<typeof XButton> | null>(null)

function openDrawer() {
  drawerOpen.value = true
  nextTick(() => drawer.value?.querySelector<HTMLElement>('a, button')?.focus())
}

function closeDrawer() {
  if (!drawerOpen.value) return
  drawerOpen.value = false
  nextTick(() => (menuButton.value?.$el as HTMLElement | undefined)?.focus())
}

watch(() => route.fullPath, () => (drawerOpen.value = false))
</script>

<template>
  <div class="layout">
    <aside class="sidebar">
      <AppLogo class="sidebar__logo" />
      <AppNav />
    </aside>

    <header class="topbar">
      <XButton ref="menuButton" icon="menu" :label="t('nav.openMenu')" :aria-expanded="drawerOpen" @click="openDrawer" />
      <div class="topbar__text">
        <div class="topbar__title">{{ route.meta.titleKey ? t(route.meta.titleKey) : '' }}</div>
        <div class="topbar__host">{{ host }}</div>
      </div>
    </header>

    <Transition name="drawer">
      <div v-if="drawerOpen" class="drawer" @keydown.esc="closeDrawer">
        <div class="drawer__scrim" @click="closeDrawer" />
        <div ref="drawer" class="drawer__panel" role="dialog" aria-modal="true" :aria-label="t('nav.main')">
          <div class="drawer__head">
            <AppLogo />
            <XButton icon="close" :label="t('nav.closeMenu')" @click="closeDrawer" />
          </div>
          <AppNav @navigate="closeDrawer" />
        </div>
      </div>
    </Transition>

    <main class="content">
      <RouterView />
    </main>
  </div>
</template>

<style scoped>
.layout {
  display: grid;
  grid-template-columns: 236px minmax(0, 1fr);
  min-height: 100vh;
  min-height: 100dvh;
}

.sidebar {
  position: sticky;
  top: 0;
  display: flex;
  flex-direction: column;
  height: 100vh;
  height: 100dvh;
  padding: 20px 14px 16px;
  border-right: 1px solid var(--line);
  background: var(--side);
}

.sidebar__logo {
  padding: 0 5px 0 4px;
}

.topbar,
.drawer {
  display: none;
}

.content {
  min-width: 0;
  padding: 24px 28px 40px;
}

@media (max-width: 767px) {
  .layout {
    display: flex;
    flex-direction: column;
  }

  .sidebar {
    display: none;
  }

  .topbar {
    position: sticky;
    top: 0;
    z-index: 10;
    display: flex;
    align-items: center;
    gap: 12px;
    height: 64px;
    padding: 0 14px;
    border-bottom: 1px solid var(--div);
    background: var(--side);
  }

  .topbar__text {
    min-width: 0;
  }

  .topbar__title {
    font: 600 var(--fs-18) / var(--lh-tight) var(--font-sans);
    letter-spacing: -.4px;
  }

  .topbar__host {
    overflow: hidden;
    color: var(--dim);
    font: 400 var(--fs-11) / var(--lh-text) var(--font-mono);
    white-space: nowrap;
    text-overflow: ellipsis;
  }

  .drawer {
    position: fixed;
    inset: 0;
    z-index: 800;
    display: block;
  }

  .drawer__scrim {
    position: absolute;
    inset: 0;
    background: var(--scrim);
    backdrop-filter: blur(2px);
  }

  .drawer__panel {
    position: relative;
    display: flex;
    flex-direction: column;
    width: min(282px, 85vw);
    height: 100%;
    padding: 16px 14px calc(16px + env(safe-area-inset-bottom));
    border-right: 1px solid var(--div);
    background: var(--side);
    box-shadow: var(--shadow-pop);
  }

  .drawer__head {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 10px;
    padding-left: 4px;
  }

  .drawer-enter-active,
  .drawer-leave-active {
    transition: opacity var(--t-med);
  }

  .drawer-enter-active .drawer__panel,
  .drawer-leave-active .drawer__panel {
    transition: transform var(--t-med) var(--ease);
  }

  .drawer-enter-from,
  .drawer-leave-to {
    opacity: 0;
  }

  .drawer-enter-from .drawer__panel,
  .drawer-leave-to .drawer__panel {
    transform: translateX(-100%);
  }

  .content {
    padding: 14px 14px 32px;
  }
}
</style>
