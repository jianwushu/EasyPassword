import axios from 'axios';
import { useAuthStore } from '@/store/auth';

const apiClient = axios.create({
  baseURL: 'http://localhost:8081/api/v1',
  headers: {
    'Content-Type': 'application/json',
  },
});

apiClient.interceptors.request.use(
  (config) => {
    const authStore = useAuthStore();
    const token = authStore.token;
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

apiClient.interceptors.response.use(
  (response) => {
    return response;
  },
  (error) => {
    // 在这里你可以处理全局错误，例如，显示一个通知
    // 例如，如果错误是 401 未授权，你可能想要重定向到登录页面
    if (error.response && error.response.status === 401) {
      const authStore = useAuthStore();
      authStore.clearAuthData();
      // 可选地重定向到登录页面
      // window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export default apiClient;