import { createRouter, createWebHistory } from 'vue-router'
import { appBase } from './env'

export const router = createRouter({
  history: createWebHistory(appBase),
  routes: [
    { path: '/', name: 'dashboard', component: () => import('./views/DashboardView.vue') },
    { path: '/login', name: 'login', component: () => import('./views/LoginView.vue') },
    { path: '/:pathMatch(.*)*', redirect: '/' },
  ],
})
