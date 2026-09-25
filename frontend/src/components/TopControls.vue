<script setup lang="ts">
import { computed, ref } from 'vue'
import { config } from '../lib/config'
import { setLanguage, t } from '../lib/i18n'
import { isDark, toggleTheme } from '../lib/theme'
import Icon from './Icon.vue'

// Language picker and theme toggle from the login screen's top bar.
const open = ref(false)
const code = (c: string) => c.slice(0, 2).toUpperCase()
const current = computed(() => code(config.lang))

function pick(c: string) {
  open.value = false
  if (c !== config.lang) setLanguage(c)
}
</script>

<template>
  <div class="controls">
    <div class="lang" :class="{ open }">
      <button class="ctl lang-btn" :class="{ 'is-open': open }" @click="open = !open">
        <Icon name="globe" :size="14" />
        <span class="code">{{ current }}</span>
        <span class="caret">{{ open ? '▲' : '▼' }}</span>
      </button>
      <template v-if="open">
        <div class="scrim" @click="open = false" />
        <div class="menu">
          <div class="menu-title">{{ t('spx.interfaceLanguage') }}</div>
          <button
            v-for="l in config.languages"
            :key="l.code"
            class="row"
            :class="{ on: l.code === config.lang }"
            @click="pick(l.code)"
          >
            <span class="rl">{{ l.name }}</span>
            <span class="rc">{{ code(l.code) }}</span>
            <Icon name="check" :size="13" class="chk" />
          </button>
        </div>
      </template>
    </div>
    <button class="ctl icon theme" :title="t('menu.theme')" @click="toggleTheme">
      <span class="swap">
        <Icon name="moon" :class="{ shown: isDark }" />
        <Icon name="sun" :class="{ shown: !isDark }" class="sun" />
      </span>
    </button>
  </div>
</template>

<style scoped>
.controls {
  display: flex;
  align-items: center;
  gap: 8px;
}
.lang {
  position: relative;
}
.lang.open {
  z-index: 50;
}
.lang-btn {
  position: relative;
  z-index: 3;
  padding: 0 10px 0 11px;
}
.code {
  font: 500 11.5px var(--font-sans);
}
.caret {
  font-size: 9px;
  color: var(--muted);
}
.scrim {
  position: fixed;
  inset: 0;
  z-index: 1;
  backdrop-filter: blur(3px) saturate(0.92);
  background: var(--scrim);
  animation: scrimIn 0.16s ease-out;
}
.menu {
  position: absolute;
  top: 38px;
  right: 0;
  width: 232px;
  z-index: 2;
  background: var(--card);
  border: 1px solid var(--line-card);
  border-radius: 12px;
  backdrop-filter: blur(10px);
  box-shadow: var(--shadow-pop);
  padding: 6px;
  display: flex;
  flex-direction: column;
  gap: 1px;
  animation: sheetUp 0.14s ease-out;
}
:root[data-theme='light'] .menu {
  background: rgba(255, 255, 255, 0.94);
}
.menu-title {
  font: 500 9.5px var(--font-mono);
  color: var(--text-3b);
  letter-spacing: 0.09em;
  text-transform: uppercase;
  padding: 7px 9px 6px;
}
.row {
  display: flex;
  align-items: center;
  gap: 9px;
  padding: 8px 10px;
  border-radius: 8px;
  color: var(--text-2);
  text-align: left;
}
.row:hover {
  background: var(--row-hover);
}
.row.on {
  background: var(--accent-soft);
  color: var(--accent);
}
.rl {
  flex: 1;
  min-width: 0;
  font: 450 12.5px var(--font-sans);
}
.rc {
  font: 500 10px var(--font-mono);
  color: var(--text-3b);
}
.chk {
  color: transparent;
}
.row.on .chk {
  color: var(--accent);
}
.swap {
  position: relative;
  width: 16px;
  height: 16px;
}
.swap svg {
  position: absolute;
  inset: 0;
  opacity: 0;
  transform: rotate(-90deg) scale(0.55);
  transition:
    transform 0.42s cubic-bezier(0.34, 1.32, 0.5, 1),
    opacity 0.26s ease;
}
.swap svg.sun {
  transform: rotate(90deg) scale(0.55);
}
.swap svg.shown {
  opacity: 1;
  transform: rotate(0) scale(1);
}
</style>
