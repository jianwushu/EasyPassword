<template>
  <n-card title="我的密码库">
    <div v-if="vaultStore.decryptedItems.length > 0">
      <n-list bordered>
        <n-list-item v-for="item in vaultStore.decryptedItems" :key="item.ID">
          <template #prefix>
            <n-tag type="info">{{ item.name }}</n-tag>
          </template>
          {{ item.account }}
        </n-list-item>
      </n-list>
    </div>
    <div v-else>
      <n-empty description="还没有保存任何密码">
        <template #extra>
          <n-input
            v-model:value="masterPassword"
            type="password"
            show-password-on="mousedown"
            placeholder="输入主密码以解密"
            @keyup.enter="decrypt"
          />
          <n-button type="primary" @click="decrypt" style="margin-top: 10px;">解密密码库</n-button>
        </template>
      </n-empty>
    </div>
  </n-card>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useVaultStore } from '../store/vault';
import { NCard, NList, NListItem, NTag, NEmpty, NInput, NButton, useMessage } from 'naive-ui';

const vaultStore = useVaultStore();
const message = useMessage();
const masterPassword = ref('');

onMounted(() => {
  vaultStore.fetchVault();
});

const decrypt = async () => {
  if (!masterPassword.value) {
    message.error('请输入主密码');
    return;
  }
  try {
    await vaultStore.decryptAllItems(masterPassword.value);
    if (vaultStore.decryptedItems.length === 0) {
        message.info('密码库为空或解密失败');
    }
  } catch (error) {
    message.error('解密失败，请检查您的主密码');
  }
};
</script>