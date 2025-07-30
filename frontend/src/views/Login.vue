<template>
  <div class="login-container">
    <LanguageSwitcher />
    <h1 class="page-title">{{ t('app.title') }}</h1>
    <n-card :title="t('login_view.title')">
      <n-form @submit.prevent="handleLogin">
        <n-form-item :label="t('login_view.identifier_label')">
          <n-input v-model:value="model.identifier" :placeholder="t('login_view.identifier_placeholder')" />
        </n-form-item>
        <n-form-item :label="t('login_view.password_label')">
          <n-input
            type="password"
            v-model:value="model.masterPassword"
            :placeholder="t('login_view.password_placeholder')"
            show-password-on="click"
          />
        </n-form-item>
        <n-button type="primary" attr-type="submit" block>
          {{ t('login_view.login_button') }}
        </n-button>
      </n-form>
      <template #footer>
        <div style="display: flex; justify-content: space-between; align-items: center;">
          <p>
            {{ t('login_view.no_account') }}
            <router-link to="/register">{{ t('login_view.register_now') }}</router-link>
          </p>
          <router-link to="/forgot-password">{{ t('login_view.forgot_password') }}</router-link>
        </div>
      </template>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { NCard, NForm, NFormItem, NInput, NButton, useMessage } from 'naive-ui';
import { useAuthStore } from '../store/auth';
import { useRouter } from 'vue-router';
import LanguageSwitcher from '../components/LanguageSwitcher.vue';

const model = ref({
  identifier: '',
  masterPassword: '',
});

const authStore: any = useAuthStore();
const router = useRouter();
const message = useMessage();
const { t } = useI18n();

const handleLogin = async () => {
  try {
    await authStore.login(model.value.identifier, model.value.masterPassword);
    // 登录成功后，路由守卫会自动处理跳转
    router.push('/');
  } catch (error) {
    message.error(t('login_view.login_failed'));
    console.error('登录失败:', error);
  }
};
</script>

<style scoped>
.login-container {
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