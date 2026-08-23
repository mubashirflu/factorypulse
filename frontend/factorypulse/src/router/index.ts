import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: () => import('@/view/login.vue'),
    },
    {
      path: '/register',
      name: 'register',
      component: () => import('@/view/Register.vue'),
    },
    {
      path: '/dashboard',
      name: 'dashboard',
      component: () => import('@/view/Dashboard.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/machines',
      name: 'machines',
      component: () => import('@/view/Machines.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/machines/:id',
      name: 'machine-detail',
      component: () => import('@/view/Machinedetaile.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/',
      redirect: '/dashboard',
    },
    {
  path: '/maintenance',
  name: 'maintenance',
  component: () => import('@/view/Maintenance.vue'),
  meta: { requiresAuth: true },
},  
{
  path: '/analytics',
  name: 'analytics',
  component: () => import('@/view/Analytics.vue'),
  meta: { requiresAuth: true },
},
  ],
})

router.beforeEach((to) => {
  const authStore = useAuthStore()

  if (to.meta.requiresAuth && !authStore.isLoggedIn) {
    return '/login'
  }
})

export default router