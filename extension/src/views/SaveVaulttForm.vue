<template>
  <div class="form-container">
    <n-card title="保存密码" size="small">
      <n-form @submit.prevent="saveCredential" size="small">
        <n-form-item-row label="网站名称">
          <n-input v-model:value="model.name" placeholder="例如：Google" />
        </n-form-item-row>
        <n-form-item-row label="网址">
          <n-input v-model:value="model.website" placeholder="https://www.google.com" />
        </n-form-item-row>
        <n-form-item-row label="账号">
          <n-input v-model:value="model.account" placeholder="请输入账号" />
        </n-form-item-row>
        <n-form-item-row label="密码">
          <n-input
            v-model:value="model.password"
            type="password"
            show-password-on="mousedown"
            placeholder="请输入密码"
          />
        </n-form-item-row>
        <n-form-item-row label="类别">
          <n-select
            v-model:value="model.category"
            :options="categoryOptions"
            placeholder="请选择类别"
          />
        </n-form-item-row>
        <n-form-item-row label="主密码">
          <n-input
            v-model:value="masterPassword"
            type="password"
            show-password-on="mousedown"
            placeholder="请输入主密码以加密保存"
          />
        </n-form-item-row>
        <n-button type="primary" attr-type="submit" >保存</n-button>
        <n-button type="warning" attr-type="button"  @click="handleClear()" >取消</n-button>
      </n-form>
    </n-card>
  </div>

</template>

<script setup lang="ts">
import { ref, onMounted} from 'vue';
import { useVaultStore } from '../store/vault';
import { NCard, NForm, NFormItemRow, NInput, NButton, useMessage, NSelect, backTopDark } from 'naive-ui';

const model = ref({
  name: '',
  website: '',
  account: '',
  password: '',
  category: '默认',
});

const vaultStore = useVaultStore();
const message = useMessage();
const masterPassword = ref('');

const categoryOptions = [
  { label: '默认', value: '默认' },
  { label: '网站', value: '网站' },
  { label: '应用', value: '应用' },
  { label: '社交媒体', value: '社交媒体' },
];

onMounted(() => {
  // 检查是否有重定向请求
  chrome.storage.local.get('pendingCredential', (result) => {
    if (result.pendingCredential) {
        model.value = result.pendingCredential;
        // 清除 pendingCredential，以防再次触发 
        chrome.storage.local.remove('pendingCredential');
    }
  });
});

const saveCredential = async () => {
  if (!masterPassword.value) {
    message.error('请输入主密码以加密您的数据');
    return;
  }
  try {
    await vaultStore.createVaultItem(model.value, masterPassword.value);
    message.success('密码已成功保存！');
    handleClear();
  } catch (error) {
    console.error('Failed to save credential:', error);
    message.error('保存失败，请稍后再试。');
  }
};

const handleClear = () => {
  // 清除 pendingCredential，以防再次触发
  model.value = {
    name: '',
    website: '',
    account: '',
    password: '',
    category: '默认',
  };
  chrome.storage.local.remove('pendingCredential');
  // 返回到上一个页面或默认页面
  window.history.back();
  
};
</script>

<style scoped>
.form-container {
  min-height: 400px;
  min-width: 350px;
  margin: 20px auto;
  padding: 20px;
  width: 100%;
  box-sizing: border-box;
}
</style>