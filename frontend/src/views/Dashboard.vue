<template>
<div class="dashboard-layout">
  <header class="dashboard-header">
    <h1 class="dashboard-title">{{ t('dashboard.title') }}</h1>
    <div class="toolbar">
      <div class="search-wrapper">
        <SearchIcon :size="18" class="search-icon" />
        <input v-model="searchQuery" type="text" :placeholder="t('dashboard.search_placeholder')" class="search-input" />
      </div>

      <div class="toolbar-actions">
        <n-tooltip trigger="hover" :disabled="!isDesktop">
          <template #trigger>
            <button @click="promptForRelock" class="icon-button" :aria-label="t('dashboard.redecrypt_tooltip')">
              <KeyIcon :size="18" />
            </button>
          </template>
          <span>{{ t('dashboard.redecrypt_tooltip') }}</span>
        </n-tooltip>

        <div class="divider"></div>

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
    <div v-if="isLoading" class="items-grid">
      <PasswordItemCardSkeleton v-for="n in 6" :key="n" />
    </div>
    <div v-else-if="vaultStore.encryptedItems.length > 0 && vaultStore.decryptedItems.length === 0" class="decrypt-prompt">
      <n-input
        v-model:value="decryptionPassword"
        type="password"
        :placeholder="t('dashboard.decrypt_prompt.placeholder')"
        @keyup.enter="handleDecryptVault"
        size="large"
      />
      <n-button type="primary" @click="handleDecryptVault" :loading="isLoading" size="large">{{ t('dashboard.decrypt_prompt.button') }}</n-button>
    </div>
    <div v-else-if="filteredItems.length > 0" class="results-container">
      <h2 class="item-count-title">{{ t('dashboard.total_items', { count: itemCount }) }}</h2>
      <n-virtual-list v-if="useVirtualList" :items="filteredItems" :item-size="180" item-resizable>
        <template #default="{ item }">
          <div :key="item.id" style="padding: 0 8px 16px;">
            <PasswordItemCard
              :item="item"
              @copy="handleCopy"
              @edit="(itemId) => openModal('edit', itemId)"
              @delete="handleDeleteItem"
            />
          </div>
        </template>
      </n-virtual-list>
      <div v-else class="items-grid">
        <PasswordItemCard
          v-for="item in filteredItems"
          :key="item.ID"
          :item="item"
          @copy="handleCopy"
          @edit="(itemId) => openModal('edit', itemId)"
          @delete="handleDeleteItem"
        />
      </div>
    </div>
    <div v-else-if="searchQuery" class="empty-state">
      <p>{{ t('dashboard.empty_state.no_search_results', { query: searchQuery }) }}</p>
    </div>
    <div v-else class="empty-state">
      <h2>{{ t('dashboard.empty_state.title') }}</h2>
      <p>{{ t('dashboard.empty_state.subtitle') }}</p>
      <n-button type="primary" @click="openModal('create')" size="large">{{ t('dashboard.empty_state.add_first_item_button') }}</n-button>
    </div>
  </main>

  <button class="fab" @click="openModal('create')" :aria-label="t('dashboard.add_new_item_button')">
    <PlusIcon :size="24" />
  </button>

  <!-- 创建/编辑模态框暂时保持不变 -->
  <!-- 统一的创建/编辑模态框 -->
  <n-modal
    v-model:show="showItemModal"
    preset="card"
    style="width: 600px; max-width: 95vw;"
    :title="modalTitle"
    :bordered="false"
  >
    <n-form @submit.prevent="handleFormSubmit">
      <n-form-item :label="t('dashboard.form.name_label')">
        <n-input v-model:value="currentItem.name" :placeholder="t('dashboard.form.name_placeholder')" required />
      </n-form-item>
      <n-form-item :label="t('dashboard.form.website_label')">
        <n-input v-model:value="currentItem.website" :placeholder="t('dashboard.form.website_placeholder')" required />
      </n-form-item>
      <n-form-item :label="t('dashboard.form.login_label')">
        <n-input v-model:value="currentItem.account" :placeholder="t('dashboard.form.login_placeholder')" required />
      </n-form-item>
      <n-form-item :label="t('dashboard.form.password_label')">
        <n-input
          v-model:value="currentItem.password"
          type="password"
          show-password-on="click"
          :placeholder="t('dashboard.form.password_placeholder')"
          required
        >
          <template #suffix>
            <n-icon @click="showPasswordGenerator = !showPasswordGenerator" style="cursor: pointer;">
              <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24"><path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 1 1-18 0a9 9 0 0 1 18 0Z" clip-rule="evenodd"/><path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 7v5l3 3"/></svg>
            </n-icon>
          </template>
        </n-input>
      </n-form-item>
      <PasswordGenerator v-if="showPasswordGenerator" @password-generated="setFormPassword" />
      <n-form-item :label="t('dashboard.form.category_label')">
        <n-input v-model:value="currentItem.category" :placeholder="t('dashboard.form.category_placeholder')" />
      </n-form-item>
      <n-form-item :label="t('dashboard.form.master_password_label')">
        <n-input
          v-model:value="modalMasterPassword"
          type="password"
          :placeholder="modalMode === 'create' ? t('dashboard.create_modal.master_password_placeholder') : t('dashboard.edit_modal.master_password_placeholder')"
          required
        />
      </n-form-item>
      <n-space justify="end">
        <n-button @click="closeModal" :disabled="isLoading">{{ t('common.cancel') }}</n-button>
        <n-button type="primary" attr-type="submit" :loading="isLoading">{{ submitButtonText }}</n-button>
      </n-space>
    </n-form>
  </n-modal>
</div>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive, computed, inject, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import type { Ref } from 'vue';
import { useMessage, useDialog, NForm, NFormItem, NInput, NButton, NSpace, NModal, NVirtualList, NIcon, NDropdown, NAvatar, NTooltip } from 'naive-ui';
import {
  Search as SearchIcon,
  Key as KeyIcon,
  Plus as PlusIcon,
  MoreVertical as MoreVerticalIcon,
  Sun as SunIcon,
  Moon as MoonIcon,
  Languages as LanguagesIcon,
  LogOut as LogOutIcon,
  User as UserIcon,
} from 'lucide-vue-next';
import PasswordGenerator from '../components/PasswordGenerator.vue';
import PasswordItemCard from '../components/PasswordItemCard.vue';
import PasswordItemCardSkeleton from '../components/PasswordItemCardSkeleton.vue';
import { useVaultStore } from '../store/vault';
import { useAuthStore } from '../store/auth';
import { useRouter } from 'vue-router';

const vaultStore = useVaultStore();
const authStore = useAuthStore();
const router = useRouter();
const message = useMessage();
const dialog = useDialog();
const { t, locale } = useI18n();

const { themeMode, toggleTheme } = inject('theme') as { themeMode: Ref<'light' | 'dark'>, toggleTheme: () => void };

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

import { h } from 'vue';

const handleLogout = () => {
  authStore.clearAuthData();
  vaultStore.$reset();
  router.push('/login');
};

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

const isDesktop = ref(window.innerWidth > 768);

const decryptionPassword = ref(''); // 用于解密保险库
const searchQuery = ref('');
const isLoading = ref(false);
const showPasswordGenerator = ref(false);

// 创建/编辑模态框的统一状态
const showItemModal = ref(false);
const modalMode = ref<'create' | 'edit'>('create');
const currentItem = reactive({
  ID: '',
  name: '',
  website: '',
  account: '',
  password: '',
  notes: '',
  category: '',
});
const modalMasterPassword = ref('');
const clipboardClearTimer = ref<number | null>(null);

const filteredItems = computed(() => {
  if (!searchQuery.value) {
    return vaultStore.decryptedItems;
  }
 
  const lowerCaseQuery = searchQuery.value.toLowerCase();
  return vaultStore.decryptedItems.filter(item => {
   return (
     item.name?.toLowerCase().includes(lowerCaseQuery) ||
     item.website?.toLowerCase().includes(lowerCaseQuery) ||
     item.Category?.toLowerCase().includes(lowerCaseQuery)
   );
 });
});

const itemCount = computed(() => filteredItems.value.length);

const useVirtualList = computed(() => filteredItems.value.length > 50);

const handleCopy = async ({ text, type, callback }: { text: string, type: string, callback: (success: boolean) => void }) => {
  if (clipboardClearTimer.value) {
    clearTimeout(clipboardClearTimer.value);
    clipboardClearTimer.value = null;
  }

  try {
    await navigator.clipboard.writeText(text);
    message.success(t('dashboard.messages.copy_success', { type }));
    callback(true);

    clipboardClearTimer.value = window.setTimeout(() => {
      navigator.clipboard.writeText('').then(() => {
        message.info(t('dashboard.messages.clipboard_cleared'));
      });
      clipboardClearTimer.value = null;
    }, 30000); // 30 秒

  } catch (err) {
    message.error(t('dashboard.messages.copy_error', { type }));
    console.error('无法复制: ', err);
    callback(false);
  }
};

const setFormPassword = (password: string) => {
  currentItem.password = password;
};

const modalTitle = computed(() => {
  return modalMode.value === 'create'
    ? t('dashboard.create_modal.title')
    : t('dashboard.edit_modal.title');
});

const submitButtonText = computed(() => {
  return modalMode.value === 'create'
    ? t('dashboard.create_modal.create_button')
    : t('dashboard.edit_modal.save_button');
});

const handleFormSubmit = async () => {
  if (isLoading.value) return;
  isLoading.value = true;
  try {
    if (!modalMasterPassword.value) {
      message.error(t('dashboard.messages.master_password_required_update'));
      return;
    }

    if (modalMode.value === 'create') {
      const { ID, ...createData } = currentItem;
      await vaultStore.createVaultItem(createData, modalMasterPassword.value);
      message.success(t('dashboard.messages.item_created_success'));
      await handleFetchVault(); // 刷新列表
    } else {
      const { ID, ...updateData } = currentItem;
      await vaultStore.updateVaultItem(ID, updateData, modalMasterPassword.value);
      message.success(t('dashboard.messages.item_updated_success'));
    }
    closeModal();
  } catch (error: any) {
    const defaultMessage = modalMode.value === 'create'
      ? t('dashboard.messages.item_created_fail')
      : t('dashboard.messages.item_updated_fail');
    message.error(error.message || defaultMessage);
    console.error(error);
  } finally {
    isLoading.value = false;
  }
};

const handleFetchVault = async () => {
  if (isLoading.value) return;
  isLoading.value = true;
  try {
    await vaultStore.fetchVault();
    if (vaultStore.encryptedItems.length > 0 && vaultStore.decryptedItems.length === 0) {
        // message.info("保险库已获取。请输入主密码以解密。");
    } else if (vaultStore.encryptedItems.length === 0) {
        // message.success("您的保险库是空的。");
    }
  } catch (error: any) {
    console.error(error);
    message.error(error.message || t('dashboard.messages.vault_fetch_fail'));
  } finally {
    isLoading.value = false;
  }
};

const handleDecryptVault = () => {
  if (isLoading.value) return;
  isLoading.value = true;
  try {
    if (!decryptionPassword.value) {
      message.warning(t('dashboard.messages.master_password_required_decryption'));
      return;
    }
    vaultStore.decryptAllItems(decryptionPassword.value);
    message.success(t('dashboard.messages.vault_decrypted_success'));
    decryptionPassword.value = '';
  } catch (error: any) {
    console.error(error);
    message.error(error.message || t('dashboard.messages.vault_decrypted_fail'));
  } finally {
    isLoading.value = false;
  }
};

const promptForRelock = () => {
  dialog.warning({
    title: t('dashboard.relock_dialog.title'),
    content: t('dashboard.relock_dialog.content'),
    positiveText: t('dashboard.relock_dialog.confirm_button'),
    negativeText: t('common.cancel'),
    onPositiveClick: () => {
      vaultStore.relockVault();
      message.success(t('dashboard.messages.vault_locked_success'));
    },
  });
};

watch(showItemModal, (val) => {
  if (!val) {
    showPasswordGenerator.value = false;
  }
});

// 组件挂载时自动获取保险库
onMounted(() => {
  handleFetchVault();
  window.addEventListener('resize', () => {
    isDesktop.value = window.innerWidth > 768;
  });
});

const openModal = (mode: 'create' | 'edit', itemId?: string) => {
  modalMode.value = mode;
  if (mode === 'create') {
    Object.assign(currentItem, {
      ID: '',
      name: '',
      website: '',
      account: '',
      password: '',
      notes: '',
      category: '',
    });
    modalMasterPassword.value = '';
  } else if (itemId) {
    const item = vaultStore.decryptedItems.find(i => i.ID === itemId);
    if (item) {
      // 处理类别键中潜在的不一致（'Category' vs 'category'）
      const category = item.Category || item.Category || '';
      Object.assign(currentItem, { ...item, category });
      modalMasterPassword.value = '';
    } else {
      console.error(`未找到 ID 为 ${itemId} 的项目。`);
      return; // 如果未找到项目，则不打开模态框
    }
  }
  showItemModal.value = true;
};

const closeModal = () => {
  showItemModal.value = false;
};

const handleDeleteItem = (itemId: string) => {
  dialog.warning({
    title: t('dashboard.delete_dialog.title'),
    content: t('dashboard.delete_dialog.content'),
    positiveText: t('dashboard.delete_dialog.confirm_button'),
    negativeText: t('common.cancel'),
    onPositiveClick: async () => {
      if (isLoading.value) return;
      isLoading.value = true;
      try {
        await vaultStore.deleteVaultItem(itemId);
        message.success(t('dashboard.messages.item_deleted_success'));
      } catch (error: any) {
        console.error(error);
        message.error(error.message || t('dashboard.messages.item_deleted_fail'));
      } finally {
        isLoading.value = false;
      }
    },
  });
};
</script>

<style scoped>
.dashboard-layout {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background-color: var(--bg-color, #f7f9fc);
  color: var(--text-color, #333);
}

.dashboard-header {
  padding: 16px 24px;
  background-color: var(--header-bg-color, #fff);
  border-bottom: 1px solid var(--border-color, #e0e0e0);
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 16px;
}

.dashboard-title {
  font-size: 1.5em;
  font-weight: 600;
  margin: 0;
}

.toolbar {
  display: flex;
  align-items: center;
  gap: 16px;
  flex-grow: 1;
  justify-content: flex-end;
}

.toolbar-actions {
  display: flex;
  align-items: center;
  gap: 16px;
}

.user-profile {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  padding: 4px 8px;
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
}

.more-icon {
  color: #999;
}

.toolbar-icons-desktop {
  display: flex;
  align-items: center;
  gap: 12px;
}

.search-wrapper {
  position: relative;
  display: flex;
  align-items: center;
}

.search-icon {
  position: absolute;
  left: 12px;
  color: #999;
}

.search-input {
  background-color: var(--input-bg-color, #f0f2f5);
  border: 1px solid transparent;
  border-radius: 20px;
  padding: 8px 16px 8px 38px;
  width: 250px;
  transition: all 0.2s;
}
.search-input:focus {
  outline: none;
  border-color: var(--primary-color, #4A90E2);
  background-color: var(--bg-color, #fff);
  width: 300px;
}

.divider {
  width: 1px;
  height: 24px;
  background-color: var(--border-color, #e0e0e0);
}

.toolbar-info {
  display: flex;
  align-items: center;
  gap: 12px;
  color: #666;
  font-size: 0.9em;
}

.results-container {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.item-count-title {
  font-size: 1.2em;
  font-weight: 500;
  color: var(--text-color-secondary, #666);
  margin: 0 0 8px 8px;
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
}
.icon-button:hover {
  background-color: var(--button-hover-bg);
}

.rotating {
  animation: spin 1s linear infinite;
}
@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.fab {
  position: fixed;
  bottom: 24px;
  right: 24px;
  width: 56px;
  height: 56px;
  border-radius: 50%;
  background-color: var(--primary-color, #4A90E2);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;
  cursor: pointer;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  transition: all 0.3s ease;
  z-index: 1000;
}

.fab:hover {
  background-color: var(--primary-color-dark, #357ABD);
  transform: translateY(-2px);
}

.dashboard-main {
  flex-grow: 1;
  padding: 24px;
  overflow-y: auto;
}

.items-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 20px;
}

.decrypt-prompt, .empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;
  gap: 16px;
  height: 100%;
  color: #666;
}

.empty-state h2 {
  font-size: 1.8em;
  margin: 0;
  color: #333;
}

/* 响应式设计 */
@media (max-width: 1280px) {
  .items-grid {
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  }
}

@media (max-width: 768px) {
  .dashboard-header {
    padding: 12px;
    flex-wrap: nowrap;
    justify-content: space-between;
  }
  .toolbar {
    gap: 8px;
  }
  .search-wrapper {
    flex-grow: 1;
  }
  .search-input, .search-input:focus {
    width: 100%;
  }
  .toolbar-info, .toolbar-icons-desktop {
    display: none;
  }
  .divider {
    display: none;
  }
  .items-grid {
    grid-template-columns: 1fr;
  }
  .fab {
    bottom: 16px;
    right: 16px;
  }
}
</style>