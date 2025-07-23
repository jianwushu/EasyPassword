<template>
  <div class="login-container">
    <h1 class="page-title">{{ $t('app.title') }}</h1>
    <n-card :title="$t('login_view.title')">
      <n-form @submit.prevent="handleLogin">
        <n-form-item :label="$t('login_view.username_label')">
          <n-input v-model:value="model.username" :placeholder="$t('login_view.username_placeholder')" />
        </n-form-item>
        <n-form-item :label="$t('login_view.password_label')">
          <n-input
            type="password"
            v-model:value="model.masterPassword"
            :placeholder="$t('login_view.password_placeholder')"
            show-password-on="click"
          />
        </n-form-item>
        <n-button type="primary" attr-type="submit" block>
          {{ $t('login_view.login_button') }}
        </n-button>
      </n-form>
      <template #footer>
        <p>
          {{ $t('login_view.no_account') }}
          <router-link to="/register">{{ $t('login_view.register_now') }}</router-link>
        </p>
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

const model = ref({
  username: '',
  masterPassword: '',
});

const authStore: any = useAuthStore();
const router = useRouter();
const message = useMessage();
const { t } = useI18n();

const handleLogin = async () => {
  try {
    console.log('登录数据:', model.value);
    await authStore.login(model.value.username, model.value.masterPassword);
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