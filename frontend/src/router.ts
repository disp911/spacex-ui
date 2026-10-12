import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import { config } from './config'
import type { IconName } from './kit'
import AppLayout from './layout/AppLayout.vue'

type TitleKey = 'nav.dashboard' | 'nav.inbounds' | 'nav.xray' | 'nav.settings' | 'common.notFound'

declare module 'vue-router' {
  interface RouteMeta {
    titleKey?: TitleKey
    /** The same screen in the old interface, until this one is ready. */
    oldPath?: string
  }
}

/** Sections of the main menu, in order. */
export const sections: { name: string; titleKey: TitleKey; icon: IconName }[] = [
  { name: 'dashboard', titleKey: 'nav.dashboard', icon: 'dashboard' },
  { name: 'inbounds', titleKey: 'nav.inbounds', icon: 'swap' },
  { name: 'xray', titleKey: 'nav.xray', icon: 'xray' },
  { name: 'settings', titleKey: 'nav.settings', icon: 'sliders' },
]

const MovingPage = () => import('./pages/MovingPage.vue')

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    component: AppLayout,
    children: [
      { path: '', name: 'dashboard', component: MovingPage, meta: { titleKey: 'nav.dashboard', oldPath: 'panel/' } },
      { path: 'inbounds', name: 'inbounds', component: MovingPage, meta: { titleKey: 'nav.inbounds', oldPath: 'panel/inbounds' } },
      { path: 'xray', name: 'xray', component: MovingPage, meta: { titleKey: 'nav.xray', oldPath: 'panel/xray' } },
      { path: 'settings', name: 'settings', component: MovingPage, meta: { titleKey: 'nav.settings', oldPath: 'panel/settings' } },
      { path: ':path(.*)*', name: 'not-found', component: () => import('./pages/NotFoundPage.vue'), meta: { titleKey: 'common.notFound' } },
    ],
  },
]

// The component kit's showcase, for development only.
if (import.meta.env.DEV) {
  routes.unshift({ path: '/kit', component: () => import('./showcase/KitShowcase.vue') })
}

export const router = createRouter({
  history: createWebHistory(config.mount),
  routes,
  scrollBehavior: (_to, _from, saved) => saved ?? { top: 0 },
})
