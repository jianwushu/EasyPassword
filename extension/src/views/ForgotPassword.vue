<script setup lang="ts">
import { ref } from 'vue';
import { NCard, NForm, NFormItem, NInput, NButton, useMessage } from 'naive-ui';
import { requestPasswordReset } from '@/api/auth';
import { useI18n } from 'vue-i18n';
import LanguageSwitcher from '../components/LanguageSwitcher.vue';

const { t } = useI18n();
const email = ref('');
const isLoading = ref(false);
const message = useMessage();

const handleSubmit = async () => {
  isLoading.value = true;
  try {
    await requestPasswordReset(email.value);
    message.success(t('forgotPassword.successMessage'));
  } catch (error) {
    message.error(t('forgotPassword.errorMessage'));
  } finally {
    isLoading.value = false;
  }
};
</script>

<template>
  <div class="forgot-password-container">
    <LanguageSwitcher />
    <h1 class="page-title">{{ t('app.title') }}</h1>
    <n-card :title="t('forgotPassword.title')">
      <p>{{ t('forgotPassword.description') }}</p>
      <n-form @submit.prevent="handleSubmit">
        <n-form-item :label="t('common.email')">
          <n-input v-model:value="email" :placeholder="t('common.emailPlaceholder')" />
        </n-form-item>
        <n-button type="primary" attr-type="submit" :loading="isLoading" block>
          {{ t('forgotPassword.submitButton') }}
        </n-button>
      </n-form>
      <template #footer>
        <div style="text-align: center;">
          <router-link to="/login">{{ t('forgotPassword.backToLogin') }}</router-link>
        </div>
      </template>
    </n-card>
  </div>
</template>

<style scoped>
.forgot-password-container {
  min-height: 400px;
  min-width: 350px;
  margin: 20px auto;
  padding: 20px;
  width: 100%;
  box-sizing: border-box;
}
.page-title {
  text-align: center;
  font-size: 2em;
  margin-bottom: 20px;
  color: #333;
}
p {
  margin-bottom: 20px;
  color: #666;
  text-align: center;
}
</style>