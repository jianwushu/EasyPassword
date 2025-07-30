<template>
  <n-config-provider :theme="theme">
    <n-global-style />
    <n-dialog-provider>
      <n-message-provider>
        <router-view />
      </n-message-provider>
    </n-dialog-provider>
  </n-config-provider>
</template>

<script setup lang="ts">
import { ref, computed, provide } from 'vue';
import {
  NMessageProvider,
  NDialogProvider,
  NConfigProvider,
  NGlobalStyle,
  useOsTheme,
  darkTheme
} from 'naive-ui';

const osThemeRef = useOsTheme();
const themeMode = ref(osThemeRef.value || 'light');

const theme = computed(() => (themeMode.value === 'dark' ? darkTheme : null));

const toggleTheme = () => {
  themeMode.value = themeMode.value === 'light' ? 'dark' : 'light';
  document.documentElement.setAttribute('data-theme', themeMode.value);
};

provide('theme', {
  themeMode,
  toggleTheme,
});

// Set initial theme
document.documentElement.setAttribute('data-theme', themeMode.value);
</script>

<style>
/* Global styles can remain here, but theme-specific variables will be in style.css */
#app {
  font-family: system-ui, Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}
</style>
