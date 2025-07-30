<template>
  <div class="language-switcher">
    <select v-model="currentLanguage" @change="switchLanguage">
      <option v-for="lang in languages" :key="lang.code" :value="lang.code">
        {{ lang.name }}
      </option>
    </select>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

const { locale } = useI18n()

const languages = [
  { code: 'en', name: 'English' },
  { code: 'zh-CN', name: '简体中文' },
  { code: 'zh-TW', name: '繁體中文' },
  { code: 'ja', name: '日本語' },
  { code: 'ko', name: '한국어' },
]

const currentLanguage = ref(locale.value)

const switchLanguage = () => {
  locale.value = currentLanguage.value
  localStorage.setItem('user-locale', currentLanguage.value)
}

watch(locale, (newLocale) => {
  currentLanguage.value = newLocale
})
</script>

<style scoped>
.language-switcher {
  position: absolute;
  top: 1rem;
  right: 1rem;
}
</style>