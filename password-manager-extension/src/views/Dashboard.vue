<template>
<div class="dashboard-layout">
  <header class="dashboard-header">
    <h1 class="dashboard-title">{{ t('dashboard.title') }}</h1>
    <div class="toolbar">

      <div class="toolbar-actions">
        <div class="toolbar-icons-desktop">
          <n-tooltip trigger="hover" :disabled="!isDesktop">
            <template #trigger>
              <button @click="toggleTheme" class="icon-button" :aria-label="t('dashboard.toggle_theme_tooltip')">
                <SunIcon v-if="themeMode === 'light'" :size="20" />
                <MoonIcon v-else :size="20" />
              </button>
            </template>
            <span>{{ t('dashboard.toggle_theme_tooltip') }}</span>
          </n-tooltip>
          <n-dropdown :options="languageOptions" @select="handleLanguageSelect" trigger="hover">
            <button class="icon-button" :aria-label="t('dashboard.toggle_language_tooltip')">
              <LanguagesIcon :size="20" />
            </button>
          </n-dropdown>
        </div>

        <div class="divider"></div>

        <n-dropdown :options="dropdownOptions" @select="handleDropdownSelect">
          <div class="user-profile">
            <n-avatar round size="small" class="user-avatar">
              <UserIcon :size="18" />
            </n-avatar>
            <span class="username">{{ authStore.username }}</span>
            <MoreVerticalIcon :size="20" class="more-icon" />
          </div>
        </n-dropdown>
      </div>
    </div>
  </header>

  <main class="dashboard-main">

  </main>

  <div class="dashboard-footer">
      <n-button @click="openPage('vault')">
        <!-- Vault Icon -->
        <span>My Vault</span>
      </n-button>
      <n-button @click="openPage('settings')">
        <!-- Settings Icon -->
        <span>Settings</span>
      </n-button>
  </div>

</div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, inject, h } from 'vue';
import type { Ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { NButton, NDropdown, NAvatar, NTooltip } from 'naive-ui';
import {
  Sun as SunIcon,
  Moon as MoonIcon,
  Languages as LanguagesIcon,
  LogOut as LogOutIcon,
  User as UserIcon,
  MoreVertical as MoreVerticalIcon,
} from 'lucide-vue-next';
import { useVaultStore } from '../store/vault';
import { useAuthStore } from '../store/auth';
import { useRouter } from 'vue-router';

const vaultStore = useVaultStore();
const authStore = useAuthStore();
const router = useRouter();
const { t, locale } = useI18n();

// Injected from App.vue, will be implemented later
const { themeMode, toggleTheme } = inject('theme', { themeMode: ref('light'), toggleTheme: () => {} }) as { themeMode: Ref<'light' | 'dark'>, toggleTheme: () => void };

const languages = [
  { label: 'English', key: 'en' },
  { label: '简体中文', key: 'zh-CN' },
  { label: '繁體中文', key: 'zh-TW' },
  { label: '日本語', key: 'ja' },
  { label: '한국어', key: 'ko' },
];

const languageOptions = computed(() => {
  return languages.map(lang => ({
    label: lang.label,
    key: lang.key,
    disabled: locale.value === lang.key
  }));
});

const handleLanguageSelect = (key: string) => {
  locale.value = key;
  localStorage.setItem('user-locale', key);
};

const handleLogout = () => {
  authStore.clearAuthData();
  vaultStore.$reset();
  router.push('/login');
};

const isDesktop = ref(window.innerWidth > 768); // Simplified for popup

const dropdownOptions = computed(() => {
  const options = [];
  if (!isDesktop.value) {
    options.push(
      {
        label: t('dashboard.toggle_theme_tooltip'),
        key: 'toggle-theme',
        icon: () => themeMode.value === 'light' ? h(SunIcon) : h(MoonIcon),
      },
      {
        label: t('dashboard.toggle_language_tooltip'),
        key: 'language',
        icon: () => h(LanguagesIcon),
        children: languageOptions.value
      }
    );
  }
  options.push({
    label: t('dashboard.logout_tooltip'),
    key: 'logout',
    icon: () => h(LogOutIcon),
  });
  return options;
});

const handleDropdownSelect = (key: string) => {
  if (languages.some(lang => lang.key === key)) {
    handleLanguageSelect(key);
    return;
  }

  switch (key) {
    case 'toggle-theme':
      toggleTheme();
      break;
    case 'logout':
      handleLogout();
      break;
  }
};

const openPage = (page: 'vault' | 'settings') => {
  let url = chrome.runtime.getURL('index.html');
  chrome.tabs.create({ url: `${url}#/${page}` });
};

onMounted(() => {

  isDesktop.value = window.innerWidth > 600;
  window.addEventListener('resize', () => {
    isDesktop.value = window.innerWidth > 600;
  });
});
</script>

<style scoped>
.dashboard-layout {
  display: flex;
  flex-direction: column;
  min-height: 400px;
  min-width: 350px;
  width: 100%;
  background-color: var(--bg-color, #f7f9fc);
  color: var(--text-color, #333);
}

.dashboard-header {
  padding: 12px;
  background-color: var(--header-bg-color, #fff);
  border-bottom: 1px solid var(--border-color, #e0e0e0);
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-shrink: 0;
}

.dashboard-title {
  font-size: 1.2em;
  font-weight: 600;
  margin: 0;
  white-space: nowrap;
}

.toolbar {
  display: flex;
  align-items: center;
  gap: 8px;
}

.toolbar-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.user-profile {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  padding: 4px;
  border-radius: 18px;
  transition: background-color 0.2s;
}

.user-profile:hover {
  background-color: var(--button-hover-bg);
}

.user-avatar {
  background-color: var(--primary-color);
}

.username {
  font-weight: 500;
  font-size: 0.9em;
  display: none; /* Hide username in popup by default */
}

.more-icon {
  color: #999;
}

.toolbar-icons-desktop {
  display: flex;
  align-items: center;
  gap: 8px;
}


.divider {
  width: 1px;
  height: 20px;
  background-color: var(--border-color, #e0e0e0);
}

.icon-button {
  background: none;
  border: none;
  cursor: pointer;
  color: #666;
  padding: 4px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}
.icon-button:hover {
  background-color: var(--button-hover-bg);
}

.dashboard-main {
  flex-grow: 1;
  padding: 16px;
  overflow-y: auto;
}


.dashboard-footer {
    display: flex;
    justify-content: space-around;
    padding: 8px 0;
    border-top: 1px solid var(--border-color, #e0e0e0);
    flex-shrink: 0;
}

.dashboard-footer button {
    background: none;
    border: none;
    cursor: pointer;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 4px;
    font-size: 0.8em;
    color: var(--text-color-secondary);
}

/* Hide desktop-only elements in the smaller popup view */
@media (max-width: 600px) {
  .toolbar-icons-desktop, .divider {
    display: none;
  }
  .username {
      display: inline; /* Show username on smaller screens if space allows */
  }
}
</style>