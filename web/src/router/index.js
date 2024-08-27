// import { createRouter, createWebHashHistory } from 'vue-router' 

import { createRouter, createWebHistory } from 'vue-router'

const routes = [{
  path: '/',
  redirect: '/login'
},
{
  path: '/checkin',
  name: 'Checkin',
  component: () => import('@/view/checkins/index.vue')
},
{
  path: '/init',
  name: 'Init',
  component: () => import('@/view/init/index.vue')
},
{
  path: '/login',
  name: 'Login',
  component: () => import('@/view/login/index.vue')
},

{
  path: '/:catchAll(.*)',
  meta: {
    closeTab: true,
  },
  component: () => import('@/view/error/index.vue')
}
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
