import { createI18n } from 'vue-i18n';
import en from './locales/en.json';
import zhCN from './locales/zh-CN.json';

const i18n = createI18n({
  legacy: false, // 必须设置为 false 才能使用 Composition API
  locale: localStorage.getItem('user-locale') || navigator.language || 'en', // 默认语言
  fallbackLocale: 'en', // 回退语言
  messages: {
    en,
    'zh-CN': zhCN,
  },
});

export default i18n;