import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'comparisons',
      component: () => import('@/views/Comparisons.vue')
    },
    {
      path: '/comparisons/:id',
      name: 'comparison_objects',
      component: () => import('@/views/Objects.vue')
    },
    {
      path: '/comparisons/:id/manage_custom_options',
      name: 'manage_custom_options',
      component: () => import('@/views/ManageCustomOptions.vue')
    },
    {
      path: '/:catchall(.*)*',
      name: 'not_found',
      component: () => import('@/views/NotFound.vue')
    }
  ]
})

export default router
