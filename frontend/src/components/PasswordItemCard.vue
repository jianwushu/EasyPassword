<template>
  <div class="password-item-card" tabindex="0" @mouseenter="cardHovered = true" @mouseleave="cardHovered = false">
    <div class="card-header">
      <div class="website-info">
        <img 
          :src="logoUrl" 
          @error="onLogoError" 
          class="logo-image" 
          alt="" 
          loading="lazy"
        />
        <div v-if="!logoLoaded" class="logo-placeholder">
          {{ websiteInitial }}
        </div>
        <div class="website-details">
          <span class="website-name">{{ item.name }}</span>
          <a :href="websiteUrl" target="_blank" class="website-url">{{ item.website }}</a>
        </div>
      </div>
    </div>

    <div class="card-body">
      <div class="info-row">
        <span class="info-label">{{ $t('password_item_card.username') }}</span>
        <span class="info-value">{{ item.account }}</span>
        <button class="icon-button" @click.stop="copy('login')" :aria-label="$t('password_item_card.copy_username_aria')">
          <CheckIcon v-if="copiedState === 'login'" :size="16" color="#4CAF50" />
          <CopyIcon v-else :size="16" />
        </button>
      </div>
      <div class="info-row">
        <span class="info-label">{{ $t('password_item_card.password') }}</span>
        <span class="info-value">{{ showPassword ? item.password : '••••••••' }}</span>
        <div class="password-actions">
          <button class="icon-button" @click.stop="togglePasswordVisibility" :aria-label="showPassword ? $t('password_item_card.hide_password_aria') : $t('password_item_card.show_password_aria')">
            <EyeIcon v-if="!showPassword" :size="16" />
            <EyeOffIcon v-else :size="16" />
          </button>
          <button class="icon-button" @click.stop="copy('password')" :aria-label="$t('password_item_card.copy_password_aria')">
            <CheckIcon v-if="copiedState === 'password'" :size="16" color="#4CAF50" />
            <CopyIcon v-else :size="16" />
          </button>
        </div>
      </div>
    </div>

    <div class="card-footer">
      <div class="timestamps">
        <span class="timestamp">{{ $t('password_item_card.created_at') }}: {{ formatTimestamp(item.CreatedAt) }}</span>
        <span class="timestamp">{{ $t('password_item_card.updated_at') }}: {{ formatTimestamp(item.UpdatedAt) }}</span>
      </div>
      <div class="category-tag" :style="{ backgroundColor: categoryColor }">
        #{{ item.Category }}
      </div>
      <div class="menu-container" @mouseenter="handleMenuEnter" @mouseleave="handleMenuLeave">
        <button class="icon-button" :aria-label="$t('password_item_card.more_actions_aria')">
          <MoreHorizontalIcon :size="20" />
        </button>
        <div v-show="showMenu" class="menu-dropdown">
          <button class="menu-item" @click.stop="$emit('edit', item.ID)">{{ $t('password_item_card.edit') }}</button>
          <button class="menu-item delete" @click.stop="$emit('delete', item.ID)">{{ $t('password_item_card.delete') }}</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { useI18n } from 'vue-i18n';
import {
  Copy as CopyIcon,
  Check as CheckIcon,
  Eye as EyeIcon,
  EyeOff as EyeOffIcon,
  MoreHorizontal as MoreHorizontalIcon,
} from 'lucide-vue-next';
import type { DecryptedVaultItem } from '../types';

const props = withDefaults(defineProps<{ item: DecryptedVaultItem }>(), {});

const emit = defineEmits(['copy', 'edit', 'delete']);

const { t } = useI18n();

const showPassword = ref(false);
const cardHovered = ref(false);
const logoLoaded = ref(true);
const copiedState = ref<null | 'login' | 'password'>(null);
const showMenu = ref(false);
let menuTimeout: number | null = null;

const handleMenuEnter = () => {
  if (menuTimeout) {
    clearTimeout(menuTimeout);
    menuTimeout = null;
  }
  showMenu.value = true;
};

const handleMenuLeave = () => {
  menuTimeout = window.setTimeout(() => {
    showMenu.value = false;
  }, 200);
};

const logoUrl = computed(() => {
  if (!props.item.website) return '';
  // 假设名称可能是一个域名，用于获取徽标
  try {
    const url = new URL(props.item.website.startsWith('http') ? props.item.website : `https://${props.item.website}`);
    return `https://logo.clearbit.com/${url.hostname}`;
  } catch (e) {
    return ''; // 无效的 URL
  }
});

const websiteUrl = computed(() => {
  if (!props.item.website) return '#';
  return props.item.website.startsWith('http') ? props.item.website : `https://${props.item.website}`;
});

const websiteInitial = computed(() => {
  return props.item.name ? props.item.name.charAt(0).toUpperCase() : '?';
});

const onLogoError = () => {
  logoLoaded.value = false;
};

const togglePasswordVisibility = () => {
  showPassword.value = !showPassword.value;
};

const copy = (type: 'login' | 'password') => {
  if (copiedState.value) return; // 防止多次点击

  const textToCopy = type === 'login' ? props.item.account : props.item.password;
  if (textToCopy) {
    emit('copy', {
      text: textToCopy,
      type: type === 'login' ? t('password_item_card.copied_username_message') : t('password_item_card.copied_password_message'),
      callback: (success: boolean) => {
        if (success) {
          copiedState.value = type;
          setTimeout(() => {
            copiedState.value = null;
          }, 2000);
        }
      }
    });
  }
};

const categoryColor = computed(() => {
  const colors = ['#4A90E2', '#50E3C2', '#F5A623', '#F8E71C', '#BD10E0', '#9013FE', '#B8E986', '#7ED321'];
  let hash = 0;
  for (let i = 0; i < props.item.Category.length; i++) {
    hash = props.item.Category.charCodeAt(i) + ((hash << 5) - hash);
  }
  const index = Math.abs(hash % colors.length);
  return colors[index];
});

const formatTimestamp = (timestamp: string) => {
  if (!timestamp) return 'N/A';
  return new Date(timestamp).toLocaleString();
};

</script>

<style scoped>
.password-item-card {
  background-color: var(--card-bg-color, #fff);
  border: 1px solid var(--card-border-color, #e0e0e0);
  border-radius: 12px;
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
  transition: all 0.2s ease-in-out;
  box-shadow: 0 2px 4px rgba(0,0,0,0.05);
}

.password-item-card:hover,
.password-item-card:focus-within {
  transform: translateY(-4px);
  box-shadow: 0 8px 16px rgba(0,0,0,0.1);
  border-color: var(--primary-color, #4A90E2);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.website-info {
  flex-grow: 1;
  overflow: hidden;
  display: flex;
  align-items: center;
  gap: 12px;
  position: relative;
}

.logo-image {
  width: 40px;
  height: 40px;
  border-radius: 8px;
  object-fit: contain;
  background-color: #f0f0f0;
}

.logo-placeholder {
  width: 40px;
  height: 40px;
  border-radius: 8px;
  background-color: #e0e0e0;
  color: #757575;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  font-weight: bold;
  position: absolute;
  top: 0;
  left: 0;
}

.website-details {
  display: flex;
  flex-direction: column;
}

.website-name {
  font-weight: 600;
  color: var(--text-color-primary, #333);
}

.website-url {
  font-size: 0.85em;
  color: var(--text-color-secondary, #757575);
  text-decoration: none;
}
.website-url:hover {
  text-decoration: underline;
}

.card-body {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.info-row {
  display: grid;
  grid-template-columns: 70px 1fr auto;
  align-items: center;
  gap: 8px;
  font-size: 0.9em;
}

.info-label {
  color: var(--text-color-secondary, #757575);
}

.info-value {
  color: var(--text-color-primary, #333);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  font-family: 'Courier New', Courier, monospace;
}

.password-actions {
  display: flex;
  gap: 4px;
}

.icon-button {
  background: none;
  border: none;
  padding: 4px;
  cursor: pointer;
  color: var(--text-color-secondary, #757575);
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  transition: background-color 0.2s;
}

.icon-button:hover {
  background-color: var(--button-hover-bg, #f0f0f0);
}

.card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 8px;
  flex-wrap: wrap;
  gap: 8px;
  position: relative;
}

.timestamps {
  flex-grow: 1;
  display: flex;
  flex-direction: column;
  font-size: 0.75em;
  color: var(--text-color-secondary, #757575);
}

.timestamp {
  white-space: nowrap;
}

.category-tag {
  padding: 4px 8px;
  border-radius: 16px;
  font-size: 0.75em;
  color: white;
  font-weight: 500;
}

.menu-container {
  position: relative;
  flex-shrink: 0;
}

.menu-dropdown {
  position: absolute;
  bottom: calc(100% + 4px);
  right: 0;
  background-color: var(--card-bg-color, #fff);
  border: 1px solid var(--card-border-color, #e0e0e0);
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
  z-index: 100;
  display: flex;
  flex-direction: column;
  padding: 4px;
  min-width: 120px;
}

.menu-item {
  background: none;
  border: none;
  padding: 8px 12px;
  text-align: left;
  cursor: pointer;
  font-size: 0.9em;
  transition: background-color 0.2s;
  border-radius: 4px;
  color: var(--text-color-primary);
  white-space: nowrap;
}

.menu-item:hover {
  background-color: var(--button-hover-bg, #f0f0f0);
}

.menu-item.delete {
  color: #E53935;
}

.menu-item.delete:hover {
  background-color: #E53935;
  color: white;
}
</style>