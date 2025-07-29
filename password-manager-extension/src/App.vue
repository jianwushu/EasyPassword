<template>
  <n-config-provider :theme="theme">
    <n-dialog-provider>
      <n-message-provider>
        <router-view />
      </n-message-provider>
    </n-dialog-provider>
  </n-config-provider>
</template>

<script setup lang="ts">
import { ref, provide, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { NConfigProvider, NMessageProvider, NDialogProvider, darkTheme } from 'naive-ui';
import type { GlobalTheme } from 'naive-ui';

const router = useRouter();
const themeMode = ref<'light' | 'dark'>('light');

const toggleTheme = () => {
  themeMode.value = themeMode.value === 'light' ? 'dark' : 'light';
  localStorage.setItem('theme-mode', themeMode.value);
  updateBodyClass();
};

const theme = computed<GlobalTheme | null>(() => {
  return themeMode.value === 'dark' ? darkTheme : null;
});

const updateBodyClass = () => {
  if (themeMode.value === 'dark') {
    document.body.classList.add('dark');
  } else {
    document.body.classList.remove('dark');
  }
};

onMounted(() => {
  const savedTheme = localStorage.getItem('theme-mode');
  if (savedTheme && (savedTheme === 'light' || savedTheme === 'dark')) {
    themeMode.value = savedTheme;
  }
  updateBodyClass();

  // 检查是否有重定向请求
  chrome.storage.local.get('redirectTo', (result) => {
    if (result.redirectTo) {
      router.push(result.redirectTo);
      // 清除重定向标志，以防再次触发
      chrome.storage.local.remove('redirectTo');
    }
  });
});

provide('theme', {
  themeMode,
  toggleTheme,
});
</script>

<style>
/* Global styles for theme */
body.dark {
  --bg-color: #1a1a1a;
  --header-bg-color: #242424;
  --text-color: #ffffffde;
  --text-color-secondary: #ffffffb3;
  --border-color: #3a3a3a;
  --button-hover-bg: #3a3a3a;
  --input-bg-color: #242424;
  --item-bg-color: #242424;
  --primary-color: #63a3e4;
  --primary-color-dark: #4a90e2;
}
</style>