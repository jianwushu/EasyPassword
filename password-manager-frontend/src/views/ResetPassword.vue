<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { NCard, NForm, NFormItem, NInput, NButton, useMessage } from 'naive-ui';
import { resetPassword } from '@/api/auth';
import { deriveKey, hashKey, generateSalt } from '@/crypto/vault';
import { useI18n } from 'vue-i18n';
import LanguageSwitcher from '../components/LanguageSwitcher.vue';

const { t } = useI18n();
const route = useRoute();
const router = useRouter();
const message = useMessage();

const token = ref('');
const newPassword = ref('');
const confirmPassword = ref('');
const isLoading = ref(false);

onMounted(() => {
  token.value = route.params.token as string;
});

const handleSubmit = async () => {
  if (newPassword.value !== confirmPassword.value) {
    message.error(t('resetPassword.passwordMismatch'));
    return;
  }
  isLoading.value = true;
  try {
    // 生成一个新的盐用于重置密码
    const newMasterSalt = generateSalt();
    const newMasterKey = await deriveKey(newPassword.value, newMasterSalt);
    const newMasterKeyHash = await hashKey(newMasterKey);
    await resetPassword(token.value, newMasterKeyHash, newMasterSalt);
    message.success(t('resetPassword.successMessage'));
    router.push('/login');
  } catch (error) {
    message.error(t('resetPassword.errorMessage'));
  } finally {
    isLoading.value = false;
  }
};
</script>

<template>
  <div class="reset-password-container">
    <LanguageSwitcher />
    <h1 class="page-title">{{ t('app.title') }}</h1>
    <n-card :title="t('resetPassword.title')">
      <n-form @submit.prevent="handleSubmit">
        <n-form-item :label="t('resetPassword.newPassword')">
          <n-input
            v-model:value="newPassword"
            type="password"
            :placeholder="t('resetPassword.newPasswordPlaceholder')"
            show-password-on="mousedown"
          />
        </n-form-item>
        <n-form-item :label="t('resetPassword.confirmPassword')">
          <n-input
            v-model:value="confirmPassword"
            type="password"
            :placeholder="t('resetPassword.confirmPasswordPlaceholder')"
            show-password-on="mousedown"
          />
        </n-form-item>
        <n-button type="primary" attr-type="submit" :loading="isLoading" block>
          {{ t('resetPassword.submitButton') }}
        </n-button>
      </n-form>
      <template #footer>
        <div style="text-align: center;">
          <router-link to="/login">{{ t('resetPassword.backToLogin') }}</router-link>
        </div>
      </template>
    </n-card>
  </div>
</template>

<style scoped>
.reset-password-container {
  max-width: 400px;
  margin: 100px auto;
  padding: 20px;
}
.page-title {
  text-align: center;
  font-size: 2em;
  margin-bottom: 20px;
  color: #333;
}
</style>