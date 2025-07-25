<template>
  <div class="register-container">
    <h1 class="page-title">{{ $t('app.title') }}</h1>
    <n-card :title="$t('register_view.title')">
      <n-form ref="formRef" :model="model" :rules="rules" @submit.prevent="handleRegister">
        <n-form-item path="username" :label="$t('register_view.username_label')">
          <n-input v-model:value="model.username" :placeholder="$t('register_view.username_placeholder')" />
        </n-form-item>
        <n-form-item path="masterPassword" :label="$t('register_view.password_label')">
          <n-input
            type="password"
            v-model:value="model.masterPassword"
            :placeholder="$t('register_view.password_placeholder')"
            show-password-on="click"
          />
        </n-form-item>
        <n-button type="primary" attr-type="submit" block :loading="loading">
          {{ $t('register_view.register_button') }}
        </n-button>
      </n-form>
      <template #footer>
        <p>
          {{ $t('register_view.has_account') }}
          <router-link to="/login">{{ $t('register_view.login_now') }}</router-link>
        </p>
      </template>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { useI18n } from 'vue-i18n';
import {
  NCard,
  NForm,
  NFormItem,
  NInput,
  NButton,
  useMessage,
  type FormInst,
  type FormRules,
} from 'naive-ui';
import { useAuthStore } from '../store/auth';
import { useRouter } from 'vue-router';

const { t } = useI18n();
const formRef = ref<FormInst | null>(null);
const loading = ref(false);
const model = ref({
  username: '',
  masterPassword: '',
});

const rules = computed<FormRules>(() => ({
  username: [
    {
      required: true,
      message: t('register_view.validation_username_required'),
      trigger: ['input', 'blur'],
    },
  ],
  masterPassword: [
    {
      required: true,
      message: t('register_view.validation_password_required'),
      trigger: ['input', 'blur'],
    },
    {
      min: 8,
      message: t('register_view.validation_password_min_length'),
      trigger: ['input', 'blur'],
    },
  ],
}));

const authStore: any = useAuthStore();
const router = useRouter();
const message = useMessage();

const handleRegister = (e: Event) => {
  e.preventDefault();
  formRef.value?.validate(async (errors) => {
    if (!errors) {
      console.log('注册数据:', model.value);
      loading.value = true;
      try {
        await authStore.register(model.value.username, model.value.masterPassword);
        message.success(t('register_view.registration_success'));
        router.push('/login');
      } catch (error) {
        message.error(t('register_view.registration_failed'));
        console.error('注册失败:', error);
      } finally {
        loading.value = false;
      }
    } else {
      console.log(errors);
      message.error(t('register_view.check_input'));
    }
  });
};
</script>

<style scoped>
.register-container {
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