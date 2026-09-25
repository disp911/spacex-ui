<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { config, panelUrl } from '../lib/config'
import { t } from '../lib/i18n'
import { isMobile } from '../lib/media'
import { stateColor, stateLabel, type ServiceState } from '../lib/status'
import Icon, { type IconName } from './Icon.vue'
import Logo from './Logo.vue'

// Page frame: sidebar on desktop; header with a burger and a slide-in
// drawer on phones. Only the content column scrolls.
const props = defineProps<{
  active: 'dashboard' | 'inbounds' | 'settings' | 'xray'
  title: string
  xrayState: ServiceState
}>()

const nav: { key: typeof props.active; icon: IconName; label: string; href: string }[] = [
  { key: 'dashboard', icon: 'dashboard', label: t('menu.dashboard'), href: panelUrl('panel/') },
  { key: 'inbounds', icon: 'connections', label: t('menu.inbounds'), href: panelUrl('panel/inbounds') },
  { key: 'settings', icon: 'settings', label: t('menu.settings'), href: panelUrl('panel/settings') },
  { key: 'xray', icon: 'xray', label: t('menu.xray'), href: panelUrl('panel/xray') },
]

const drawer = ref(false)
watch(isMobile, (m) => !m && (drawer.value = false))
const pill = computed(() => ({ color: stateColor(props.xrayState), label: stateLabel(props.xrayState) }))
</script>

<template>
  <div class="shell" :class="{ mobile: isMobile }">
    <aside v-if="!isMobile" class="sidebar">
      <Logo :size="26" class="brand" />
      <nav class="nav">
        <a v-for="n in nav" :key="n.key" :href="n.href" class="nav-item" :class="{ on: n.key === active }">
          <Icon :name="n.icon" />{{ n.label }}
        </a>
      </nav>
      <div class="bottom">
        <div class="pill" :style="{ color: pill.color }">
          <span class="dot" />
          <span class="pill-name">XRAY</span>
          <span class="pill-state">{{ pill.label }}</span>
        </div>
        <a :href="panelUrl('logout')" class="nav-item logout"><Icon name="logout" />{{ t('menu.logout') }}</a>
      </div>
    </aside>

    <header v-else class="mhead">
      <button class="ctl icon burger" :aria-label="t('spx.menu')" @click="drawer = true"><Icon name="menu" /></button>
      <div class="mtitles">
        <div class="mtitle">{{ title }}</div>
        <div class="mhost">{{ config.host }}</div>
      </div>
    </header>

    <Teleport to="body">
      <template v-if="drawer">
        <div class="drawer-scrim" @click="drawer = false" />
        <aside class="drawer">
          <div class="drawer-top">
            <Logo :size="26" class="brand" />
            <button class="ctl icon" :aria-label="t('close')" @click="drawer = false"><Icon name="close" :size="14" /></button>
          </div>
          <nav class="nav big">
            <a v-for="n in nav" :key="n.key" :href="n.href" class="nav-item" :class="{ on: n.key === active }">
              <Icon :name="n.icon" />{{ n.label }}
            </a>
          </nav>
          <div class="bottom">
            <div class="pill" :style="{ color: pill.color }">
              <span class="dot" />
              <span class="pill-name">XRAY</span>
              <span class="pill-state">{{ pill.label }}</span>
            </div>
            <a :href="panelUrl('logout')" class="nav-item logout"><Icon name="logout" />{{ t('menu.logout') }}</a>
          </div>
        </aside>
      </template>
    </Teleport>

    <main class="content" :class="{ 'no-scrollbar': isMobile }">
      <div v-if="!isMobile" class="head">
        <div class="titles">
          <h1 class="title">{{ title }}</h1>
          <div class="sub">{{ t('spx.server') }} <span class="host">{{ config.host }}</span></div>
        </div>
        <slot name="head" />
      </div>
      <slot />
    </main>
  </div>
</template>

<style scoped>
.shell {
  height: 100vh;
  height: 100dvh;
  display: flex;
  background: var(--bg);
}
.sidebar {
  --icon-bg: var(--sidebar);
  width: 236px;
  flex: none;
  position: relative;
  display: flex;
  flex-direction: column;
  padding: 20px 14px 16px;
  background: var(--sidebar);
  border-right: 1px solid var(--line-soft);
}
.brand {
  padding: 0 5px 0 4px;
  gap: 10px;
}
.brand :deep(.name) {
  font-size: 14.5px;
}
.nav {
  position: absolute;
  left: 14px;
  right: 14px;
  top: 50%;
  transform: translateY(-50%);
  display: flex;
  flex-direction: column;
  gap: 2px;
}
.nav-item {
  display: flex;
  align-items: center;
  gap: 11px;
  height: 38px;
  padding: 0 11px;
  border-radius: 10px;
  color: var(--text-3);
  font: 450 13px var(--font-sans);
  transition:
    background 0.18s ease,
    color 0.18s ease;
}
.nav-item:hover {
  background: var(--row-hover-soft);
  color: var(--text);
}
.nav-item.on {
  background: var(--accent-soft);
  color: var(--accent);
  font-weight: 500;
}
.nav-item.on :deep(svg) {
  --icon-bg: var(--accent-soft);
}
.bottom {
  margin-top: auto;
  display: flex;
  flex-direction: column;
  gap: 10px;
}
.pill {
  display: flex;
  align-items: center;
  gap: 9px;
  padding: 9px 11px;
  border-radius: 10px;
  background: var(--fill);
  border: 1px solid var(--line-soft);
}
.pill-name {
  font: 500 10.5px var(--font-mono);
  letter-spacing: 0.06em;
  color: var(--text-3b);
}
.pill-state {
  margin-left: auto;
  font: 500 11px var(--font-sans);
}
.logout {
  color: var(--text-3b);
}
.logout:hover {
  background: rgba(var(--danger-rgb), 0.1);
  color: var(--danger-text);
}

.content {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  padding: 22px 26px 24px;
  overflow-x: hidden;
  overflow-y: auto;
}
.head {
  display: flex;
  align-items: center;
  gap: 20px;
  margin-bottom: 16px;
  min-height: 45px;
  flex: none;
}
.title {
  margin: 0 0 3px;
  font: 600 21px var(--font-sans);
  letter-spacing: -0.5px;
  white-space: nowrap;
}
.sub {
  font: 400 12px var(--font-sans);
  color: var(--muted);
  white-space: nowrap;
}
.host {
  font-family: var(--font-mono);
  color: var(--text-2);
}

/* Phones */
.shell.mobile {
  flex-direction: column;
}
.mhead {
  flex: none;
  height: 72px;
  padding: 0 16px;
  padding-top: env(safe-area-inset-top);
  display: flex;
  align-items: center;
  gap: 13px;
  background: var(--sidebar);
  border-bottom: 1px solid var(--line-soft);
}
.burger {
  width: 36px;
  height: 36px;
  border-radius: 11px;
}
.mtitles {
  min-width: 0;
}
.mtitle {
  font: 600 17px var(--font-sans);
  letter-spacing: -0.4px;
}
.mhost {
  font: 400 11px var(--font-mono);
  color: var(--muted);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.mobile .content {
  padding: 13px 14px 16px;
  padding-bottom: max(16px, env(safe-area-inset-bottom));
}
.drawer-scrim {
  position: fixed;
  inset: 0;
  z-index: 200;
  background: var(--scrim-strong);
  backdrop-filter: blur(2px);
  animation: scrimIn 0.2s ease-out;
}
.drawer {
  --icon-bg: var(--sidebar);
  position: fixed;
  top: 0;
  left: 0;
  bottom: 0;
  z-index: 201;
  width: 282px;
  max-width: 86vw;
  display: flex;
  flex-direction: column;
  padding: max(16px, env(safe-area-inset-top)) 14px max(18px, env(safe-area-inset-bottom));
  background: var(--sidebar);
  border-right: 1px solid var(--line);
  box-shadow: var(--shadow-drawer);
  animation: drawerIn 0.24s cubic-bezier(0.22, 0.7, 0.3, 1);
}
.drawer-top {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 0 4px;
}
.drawer-top .brand {
  padding: 0;
}
.drawer-top .brand :deep(.name) {
  font-size: 15px;
}
.drawer-top .ctl {
  margin-left: auto;
  color: var(--text-3b);
}
.nav.big {
  gap: 3px;
}
.nav.big .nav-item,
.drawer .logout {
  height: 46px;
  gap: 12px;
  padding: 0 12px;
  border-radius: 11px;
  font-size: 14px;
}
</style>
