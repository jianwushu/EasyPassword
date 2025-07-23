import { defineStore } from 'pinia';
import { useAuthStore } from './auth';
import * as api from '../api/vault';
import { encryptVaultItem, decryptVaultData } from '../crypto/vault';
import type { VaultItem, DecryptedVaultItem } from '../types';

export const useVaultStore = defineStore('vault', {
  state: () => ({
    encryptedItems: [] as VaultItem[],
    decryptedItems: [] as DecryptedVaultItem[],
  }),
  actions: {
    setEncryptedItems(items: VaultItem[]) {
      this.encryptedItems = items;
    },
    setDecryptedItems(items: DecryptedVaultItem[]) {
      this.decryptedItems = items;
    },
    clearVault() {
      this.encryptedItems = [];
      this.decryptedItems = [];
    },
    async createVaultItem(itemData: { [key: string]: any }, masterPassword: string) {
      const authStore = useAuthStore();
      if (!authStore.masterSalt) {
        throw new Error('Master salt is not available. Please log in again.');
      }

      const { category, ...dataToEncrypt } = itemData;
      const encryptedData = await encryptVaultItem(dataToEncrypt, masterPassword, authStore.masterSalt);

      await api.addVaultItem({ encrypted_data: encryptedData, category: category || '' });

      // Optionally, refresh the vault after adding
      this.fetchVault();
    },

    async fetchVault() {
      try {
        const response = await api.getVault();
        this.setEncryptedItems(response.data || []);
        console.log('Fetched encrypted items:', this.encryptedItems);
        // 获取新数据时清除已解密的数据
        this.setDecryptedItems([]);
      } catch (error) {
        console.error('Failed to fetch vault:', error);
        // 在UI中适当地处理错误
      }
    },

    async decryptAllItems(masterPassword: string) {
      const authStore = useAuthStore();
      if (!authStore.masterSalt) {
        throw new Error('Master salt is not available. Please log in again.');
      }
      if (this.encryptedItems.length === 0) {
        return;
      }
      // 解密所有加密的密码项
      const decryptionPromises = this.encryptedItems.map(async item => {
        try {
          const decryptedData = await decryptVaultData(
            item.EncryptedData,
            masterPassword,
            authStore.masterSalt!
          );
          const { EncryptedData, ...meta } = item;
          return {
            ...(decryptedData as object),
            ...meta,
          } as DecryptedVaultItem;
        } catch (e) {
          console.error(`Failed to decrypt item ${item.ID}:`, e);
          // 解密失败时返回 null
          return null;
        }
      });

      const decrypted = (await Promise.all(decryptionPromises)).filter(
        (item): item is DecryptedVaultItem => item !== null
      );

      this.setDecryptedItems(decrypted);
    },

    async updateVaultItem(itemId: string, itemData: { [key: string]: any }, masterPassword: string) {
      const authStore = useAuthStore();
      if (!authStore.masterSalt) {
        throw new Error('Master salt is not available. Please log in again.');
      }

      const { category, ...dataToEncrypt } = itemData;
      const encryptedData = await encryptVaultItem(
        dataToEncrypt,
        masterPassword,
        authStore.masterSalt
      );

      const updatedItem = await api.updateVaultItem(itemId, { encrypted_data: encryptedData, category: category || '' });

      // 更新本地存储中的加密项
      const index = this.encryptedItems.findIndex(item => item.ID === itemId);
      if (index !== -1) {
        this.encryptedItems[index] = updatedItem.data;
      }

      // 更新本地存储中的解密项
      const decryptedIndex = this.decryptedItems.findIndex(item => item.ID === itemId);
      if (decryptedIndex !== -1) {
        const decryptedData = await decryptVaultData(
          updatedItem.data.EncryptedData,
          masterPassword,
          authStore.masterSalt
        );

        const { ID, UserID, CreatedAt } = this.decryptedItems[decryptedIndex];
        this.decryptedItems[decryptedIndex] = {
          ID,
          UserID,
          CreatedAt,
          ...(decryptedData as object),
          Category: updatedItem.data.Category,
          UpdatedAt: updatedItem.data.UpdatedAt,
        } as DecryptedVaultItem;
      }
    },

    async deleteVaultItem(itemId: string) {
      await api.deleteVaultItem(itemId);

      // 从本地存储中移除
      this.encryptedItems = this.encryptedItems.filter(item => item.ID !== itemId);
      this.decryptedItems = this.decryptedItems.filter(item => item.ID !== itemId);
    },
    relockVault() {
      this.decryptedItems = [];
    },
  },
});