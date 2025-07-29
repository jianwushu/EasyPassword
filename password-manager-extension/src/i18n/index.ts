import { createI18n } from 'vue-i18n'
import messages from '@intlify/unplugin-vue-i18n/messages'

const i18n = createI18n({
  legacy: false, // 必须设置为 false 才能使用 Composition API
  locale: localStorage.getItem('user-locale') || navigator.language || 'en', // 默认语言
  fallbackLocale: 'en', // 回退语言
  messages,
})

export default i18n;