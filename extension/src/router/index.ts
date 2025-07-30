import { createRouter, createWebHashHistory, RouteRecordRaw } from 'vue-router';
import { useAuthStore } from '../store/auth';

import Dashboard from '../views/Dashboard.vue';
import Login from '../views/Login.vue';
import Vault from '../views/Vault.vue';
import AddVault from '../views/SaveVaulttForm.vue';
import Settings from '../views/Settings.vue';
import Register from '../views/Register.vue';
import ForgotPassword from '../views/ForgotPassword.vue';
import ResetPassword from '../views/ResetPassword.vue';

const routes: Array<RouteRecordRaw> = [
  {
    path: '/',
    name: 'Dashboard',
    component: Dashboard,
  },
  {
    path: '/login',
    name: 'Login',
    component: Login,
  },
  {
    path: '/vault',
    name: 'Vault',
    component: Vault,
  },
  {
    path: '/vault/add',
    name: 'AddVault',
    component: AddVault,
  },
  {
    path: '/settings',
    name: 'Settings',
    component: Settings,
  },
  {
    path: '/register',
    name: 'Register',
    component: Register,
  },
  {
    path: '/forgot-password',
    name: 'ForgotPassword',
    component: ForgotPassword,
  },
  {
    path: '/reset-password',
    name: 'ResetPassword',
    component: ResetPassword,
  },
];

const router = createRouter({
  history: createWebHashHistory(),
  routes,
});

const whiteList = ['/login', '/register', '/forgot-password', '/reset-password']; // 无需重定向的白名单

router.beforeEach((to, from, next) => {
  const authStore = useAuthStore();

  if (authStore.isAuthenticated) {
    // 如果用户已登录
    if (to.path === '/login') {
      // 如果已登录，则重定向到主页
      next({ path: '/' });
    } else {
      next();
    }
  } else {
    // 如果用户未登录
    if (whiteList.indexOf(to.path) !== -1) {
      // 在免登录白名单中，直接进入
      next();
    } else {
      // 否则全部重定向到登录页
      next('/login');
    }
  }
});

export default router;