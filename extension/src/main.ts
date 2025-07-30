import { createApp } from 'vue';
import { createPinia } from 'pinia';
import App from './App.vue';
import router from './router';
import { useAuthStore } from './store/auth';
import { createChromeStorage } from './store/storage';
import './style.css';
import naive from 'naive-ui';
import i18n from './i18n';

const pinia = createPinia();
const app = createApp(App);
app.use(pinia);

// Manually restore state on startup
const storage = createChromeStorage();
storage.getItem('auth').then((persistedState) => {
  if (persistedState) {
    const authStore = useAuthStore();
    authStore.$patch(JSON.parse(persistedState));
  }
});
app.use(router);
app.use(i18n);
app.use(naive);
app.mount('#app');