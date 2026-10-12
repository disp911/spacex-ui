import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import { config } from './config'
import type { IconName } from './kit'
import AppLayout from './layout/AppLayout.vue'
import { useSession } from './stores/session'

type TitleKey = 'nav.dashboard' | 'nav.inbounds' | 'nav.xray' | 'nav.settings' | 'common.notFound' | 'auth.title'

declare module 'vue-router' {
  interface RouteMeta {
    titleKey?: TitleKey
    /** Open without a session. */
    public?: boolean
  }
}

/** Sections of the main menu, in order. */
export const sections: { name: string; titleKey: TitleKey; icon: IconName }[] = [
  { name: 'dashboard', titleKey: 'nav.dashboard', icon: 'dashboard' },
  { name: 'inbounds', titleKey: 'nav.inbounds', icon: 'swap' },
  { name: 'xray', titleKey: 'nav.xray', icon: 'xray' },
  { name: 'settings', titleKey: 'nav.settings', icon: 'sliders' },
]

// Screens that are not built yet.
const DraftPage = () => import('./pages/DraftPage.vue')

const routes: RouteRecordRaw[] = [
  { path: '/login', name: 'login', component: () => import('./pages/LoginPage.vue'), meta: { titleKey: 'auth.title', public: true } },
  {
    path: '/',
    component: AppLayout,
    children: [
      { path: '', name: 'dashboard', component: () => import('./pages/DashboardPage.vue'), meta: { titleKey: 'nav.dashboard' } },
      { path: 'inbounds', name: 'inbounds', component: DraftPage, meta: { titleKey: 'nav.inbounds' } },
      { path: 'xray', name: 'xray', component: DraftPage, meta: { titleKey: 'nav.xray' } },
      { path: 'settings', name: 'settings', component: DraftPage, meta: { titleKey: 'nav.settings' } },
      { path: ':path(.*)*', name: 'not-found', component: () => import('./pages/NotFoundPage.vue'), meta: { titleKey: 'common.notFound' } },
    ],
  },
]

// The component kit's showcase, for development only.
if (import.meta.env.DEV) {
  routes.unshift({ path: '/kit', component: () => import('./showcase/KitShowcase.vue'), meta: { public: true } })
}

export const router = createRouter({
  history: createWebHistory(config.mount),
  routes,
  scrollBehavior: (_to, _from, saved) => saved ?? { top: 0 },
})

// Without a session every screen but the public ones leads to the login screen,
// which returns to the requested address after login.
router.beforeEach((to) => {
  if (to.meta.public || useSession().loggedIn) return true
  return { name: 'login', query: to.fullPath === '/' ? {} : { next: to.fullPath } }
})
