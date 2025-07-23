<template>
  <n-card size="small" :title="$t('password_generator.card_title')" style="margin-top: 1rem;">
    <n-space vertical>
      <n-space align="center">
        <n-form-item :label="$t('password_generator.length')" style="margin-bottom: 0;">
          <n-input-number v-model:value="options.length" :min="8" :max="128" style="width: 100px;" />
        </n-form-item>
        <n-checkbox v-model:checked="options.includeNumbers">
          {{ $t('password_generator.include_numbers') }}
        </n-checkbox>
        <n-checkbox v-model:checked="options.includeSymbols">
          {{ $t('password_generator.include_symbols') }}
        </n-checkbox>
      </n-space>
      <n-input-group>
        <n-input :value="generatedPassword" readonly :placeholder="$t('password_generator.generated_password')" />
        <n-button @click="generatePassword">{{ $t('password_generator.generate') }}</n-button>
        <n-button type="primary" @click="emitPassword" :disabled="!generatedPassword">{{ $t('password_generator.use_password') }}</n-button>
      </n-input-group>
    </n-space>
  </n-card>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';
import { useI18n } from 'vue-i18n';
import {
  NCard,
  NSpace,
  NInputNumber,
  NCheckbox,
  NInputGroup,
  NInput,
  NButton,
  NFormItem
} from 'naive-ui';

useI18n();

const emit = defineEmits(['password-generated']);

const options = reactive({
  length: 16,
  includeNumbers: true,
  includeSymbols: true,
});

const generatedPassword = ref('');

const generatePassword = () => {
  const lower = 'abcdefghijklmnopqrstuvwxyz';
  const upper = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ';
  const numbers = '0123456789';
  const symbols = '!@#$%^&*()_+-=[]{}|;:,.<>?';

  let charset = lower + upper;
  if (options.includeNumbers) charset += numbers;
  if (options.includeSymbols) charset += symbols;

  let password = '';
  for (let i = 0; i < options.length; i++) {
    const randomIndex = Math.floor(Math.random() * charset.length);
    password += charset[randomIndex];
  }
  generatedPassword.value = password;
};

const emitPassword = () => {
  if (generatedPassword.value) {
    emit('password-generated', generatedPassword.value);
  }
};
</script>

<style scoped>
/* 由于 Naive UI 处理组件样式，因此不再需要作用域样式。 */
</style>