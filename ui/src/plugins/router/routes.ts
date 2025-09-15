export const routes = [
  { path: '/', redirect: '/dashboard' },
  {
    path: '/',
    meta: { requiresAuth: true },
    component: () => import('@/layouts/default.vue'),
    children: [
      {
        path: 'dashboard',
        name: 'dashboard',
        component: () => import('@/pages/dashboard/index.vue'),
      },
      {
        path: 'account-settings',
        component: () => import('@/pages/account/settings.vue'),
      },
    ],
  },
  {
    path: '/',
    component: () => import('@/layouts/blank.vue'),
    children: [
      {
        path: 'login',
        name: 'login',
        meta: { guestOnly: true },
        component: () => import('@/pages/auth/login.vue'),
      },
      {
        path: 'register',
        name: 'register',
        meta: { guestOnly: true },
        component: () => import('@/pages/auth/register.vue'),
      },
      {
        path: 'forgot-password',
        name: 'forgot-password',
        meta: { guestOnly: true },
        component: () => import('@/pages/auth/forgot-password.vue'),
      },
      {
        path: '/:pathMatch(.*)*',
        component: () => import('@/pages/[...error].vue'),
      },
    ],
  },
]
