<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import XIcon, { type IconName } from '@/components/ui/XIcon.vue'
import XLogo from '@/components/ui/XLogo.vue'
import ThemeToggle from './ThemeToggle.vue'
import LangSwitch from './LangSwitch.vue'
import { panelUrl } from '@/env'
import { t } from '@/i18n'
import { useEscape } from '@/composables/useMedia'
import { useScrollLock } from '@/composables/useScrollLock'

// Page frame: fixed sidebar on desktop; top bar + slide-in drawer on mobile.
defineProps<{ title: string; subtitle?: string }>()

interface NavItem {
  key: string
  icon: IconName
  label: string
  to?: string
  href?: string
}

// Pages that are not redesigned yet link to the legacy panel screens.
const nav = computed<NavItem[]>(() => [
  { key: 'dashboard', icon: 'dashboard', label: t('nav.dashboard'), to: '/' },
  { key: 'inbounds', icon: 'inbounds', label: t('nav.inbounds'), to: '/inbounds' },
  { key: 'xray', icon: 'xray', label: t('nav.xray'), href: panelUrl('panel/xray') },
  { key: 'settings', icon: 'settings', label: t('nav.settings'), href: panelUrl('panel/settings') },
])

const route = useRoute()
const drawer = ref(false)
watch(() => route.fullPath, () => (drawer.value = false))
useEscape(() => (drawer.value = false))
useScrollLock(drawer)

const logoutUrl = panelUrl('logout')
</script>

<template>
  <div class="shell">
    <!-- Desktop sidebar -->
    <aside class="shell__side">
      <XLogo :size="26" class="shell__logo" />
      <nav class="shell__nav" :aria-label="t('nav.menu')">
        <template v-for="n in nav" :key="n.key">
          <RouterLink v-if="n.to" :to="n.to" class="nav-item" active-class="is-active" exact-active-class="is-active">
            <XIcon :name="n.icon" />{{ n.label }}
          </RouterLink>
          <a v-else :href="n.href" class="nav-item"><XIcon :name="n.icon" />{{ n.label }}</a>
        </template>
      </nav>
      <div class="shell__side-foot">
        <div class="shell__prefs">
          <LangSwitch />
          <ThemeToggle />
        </div>
        <a :href="logoutUrl" class="nav-item nav-item--danger"><XIcon name="logout" />{{ t('nav.logout') }}</a>
      </div>
    </aside>

    <!-- Mobile top bar -->
    <header class="shell__top">
      <button type="button" class="shell__burger" :aria-label="t('nav.menu')" @click="drawer = true">
        <XIcon name="menu" :size="17" :stroke-width="1.5" />
      </button>
      <div class="shell__top-titles">
        <div class="shell__top-title">{{ title }}</div>
        <div v-if="subtitle" class="shell__top-sub">{{ subtitle }}</div>
      </div>
    </header>

    <!-- Mobile drawer -->
    <Transition name="drawer">
      <div v-if="drawer" class="shell__drawer" @mousedown.self="drawer = false">
        <div class="shell__drawer-panel" role="dialog" aria-modal="true" :aria-label="t('nav.menu')">
          <div class="shell__drawer-head">
            <XLogo :size="26" />
            <button type="button" class="shell__burger" :aria-label="t('common.close')" @click="drawer = false">
              <XIcon name="close" :size="15" :stroke-width="1.6" />
            </button>
          </div>
          <nav class="shell__nav">
            <template v-for="n in nav" :key="n.key">
              <RouterLink v-if="n.to" :to="n.to" class="nav-item nav-item--lg" exact-active-class="is-active">
                <XIcon :name="n.icon" :size="17" />{{ n.label }}
              </RouterLink>
              <a v-else :href="n.href" class="nav-item nav-item--lg"><XIcon :name="n.icon" :size="17" />{{ n.label }}</a>
            </template>
          </nav>
          <div class="shell__side-foot">
            <div class="shell__prefs">
              <LangSwitch />
              <ThemeToggle />
            </div>
            <a :href="logoutUrl" class="nav-item nav-item--lg nav-item--danger"><XIcon name="logout" :size="17" />{{ t('nav.logout') }}</a>
          </div>
        </div>
      </div>
    </Transition>

    <main class="shell__main">
      <header class="shell__head">
        <div class="shell__titles">
          <h1 class="shell__title">{{ title }}</h1>
          <div v-if="subtitle" class="shell__sub">
            {{ t('common.server') }} <span class="mono">{{ subtitle }}</span>
          </div>
        </div>
        <div v-if="$slots.aside" class="shell__head-aside"><slot name="aside" /></div>
      </header>
      <div v-if="$slots.aside" class="shell__aside-mobile"><slot name="aside" /></div>
      <slot />
    </main>
  </div>
</template>

<style scoped>
.shell {
  min-height: 100vh;
  min-height: 100dvh;
  display: flex;
}

/* ---- Sidebar -------------------------------------------------------- */
.shell__side {
  position: sticky;
  top: 0;
  flex: none;
  width: var(--sidebar-w);
  height: 100vh;
  height: 100dvh;
  display: flex;
  flex-direction: column;
  padding: 20px 14px 16px;
  background: var(--bg-nav);
  border-right: 1px solid var(--border);
  --x-icon-knob: var(--bg-nav);
}
.shell__logo {
  padding: 0 4px;
  margin-bottom: 22px;
}
.shell__nav {
  margin: auto 0;
  display: flex;
  flex-direction: column;
  gap: 2px;
}
.shell__side-foot {
  display: flex;
  flex-direction: column;
  gap: 10px;
}
.shell__prefs {
  display: flex;
  gap: 8px;
  padding: 0 2px;
}
.shell__prefs :deep(.lang) {
  flex: 1;
}
.shell__prefs :deep(.lang__trigger) {
  width: 100%;
}
.shell__prefs :deep(.lang__menu) {
  top: auto;
  bottom: calc(100% + 8px);
  left: 0;
  right: auto;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 11px;
  height: var(--h-md);
  padding: 0 11px;
  border-radius: var(--r-md);
  color: var(--text-3);
  font: var(--fw-regular) var(--fs-md) var(--font-sans);
  text-decoration: none;
  transition: background var(--dur-fast) ease, color var(--dur-fast) ease;
}
.nav-item:hover {
  background: var(--hover);
  color: var(--text);
}
.nav-item.is-active {
  background: var(--accent-soft);
  color: var(--accent);
  font-weight: var(--fw-medium);
}
.nav-item--danger:hover {
  background: var(--danger-soft);
  color: var(--danger-text);
}
.nav-item--lg {
  height: var(--h-lg);
  padding: 0 12px;
  gap: 12px;
  font-size: var(--fs-lg);
}

/* ---- Main ----------------------------------------------------------- */
.shell__main {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  padding: 22px 26px 24px;
}
.shell__head {
  display: flex;
  align-items: center;
  gap: 20px;
  min-height: 45px;
  margin-bottom: 16px;
}
.shell__title {
  margin: 0 0 3px;
  font: var(--fw-semibold) var(--fs-2xl) / 1.2 var(--font-sans);
  letter-spacing: -0.5px;
  white-space: nowrap;
}
.shell__sub {
  font: var(--fw-regular) var(--fs-sm) var(--font-sans);
  color: var(--text-4);
  white-space: nowrap;
}
.shell__sub .mono {
  color: var(--text-2);
}
.shell__head-aside {
  margin-left: auto;
  width: calc(50% - 6px);
  min-width: 0;
}
.shell__aside-mobile {
  display: none;
}

/* ---- Mobile --------------------------------------------------------- */
.shell__top,
.shell__drawer {
  display: none;
}
.shell__burger {
  flex: none;
  width: var(--h-md);
  height: var(--h-md);
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--r-md);
  border: 1px solid var(--border-control);
  background: var(--field);
  color: var(--text-2);
  transition: border-color var(--dur-base) ease, background var(--dur-base) ease, box-shadow var(--dur-base) ease;
}
.shell__burger:hover {
  border-color: var(--accent-focus);
  background: var(--accent-soft-hover);
  box-shadow: 0 0 0 3px var(--accent-ring);
  color: var(--text);
}

@media (max-width: 767px) {
  .shell {
    flex-direction: column;
  }
  .shell__side,
  .shell__head {
    display: none;
  }
  .shell__top {
    position: sticky;
    top: 0;
    z-index: 20;
    display: flex;
    align-items: center;
    gap: 11px;
    height: calc(var(--topbar-h) + env(safe-area-inset-top));
    padding: env(safe-area-inset-top) 15px 0;
    background: var(--bg-nav);
    border-bottom: 1px solid var(--divider);
  }
  .shell__top-titles {
    min-width: 0;
  }
  .shell__top-title {
    font: var(--fw-semibold) var(--fs-xl) / 1.15 var(--font-sans);
    letter-spacing: -0.4px;
  }
  .shell__top-sub {
    font: var(--fw-regular) var(--fs-xs) var(--font-mono);
    color: var(--text-4);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
  .shell__main {
    padding: 13px 14px calc(16px + env(safe-area-inset-bottom));
  }
  .shell__aside-mobile {
    display: block;
    margin-bottom: 12px;
  }

  .shell__drawer {
    position: fixed;
    inset: 0;
    z-index: 40;
    display: block;
    background: var(--scrim);
    backdrop-filter: blur(2px);
  }
  .shell__drawer-panel {
    position: relative;
    width: var(--drawer-w);
    max-width: 86vw;
    height: 100%;
    display: flex;
    flex-direction: column;
    padding: calc(16px + env(safe-area-inset-top)) 14px calc(18px + env(safe-area-inset-bottom));
    background: var(--bg-nav);
    border-right: 1px solid var(--divider);
    box-shadow: 26px 0 70px rgba(0, 0, 0, 0.35);
    --x-icon-knob: var(--bg-nav);
  }
  .shell__drawer-head {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 0 0 4px;
    margin-bottom: 20px;
  }
  .shell__drawer .shell__nav {
    gap: 3px;
  }
}

.drawer-enter-active,
.drawer-leave-active {
  transition: opacity 0.2s ease;
}
.drawer-enter-active .shell__drawer-panel,
.drawer-leave-active .shell__drawer-panel {
  transition: transform 0.24s cubic-bezier(0.22, 0.7, 0.3, 1);
}
.drawer-enter-from,
.drawer-leave-to {
  opacity: 0;
}
.drawer-enter-from .shell__drawer-panel,
.drawer-leave-to .shell__drawer-panel {
  transform: translateX(-100%);
}
</style>
