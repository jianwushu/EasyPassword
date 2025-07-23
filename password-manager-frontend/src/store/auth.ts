import { defineStore } from 'pinia';
import * as api from '../api/auth';
import { deriveKey } from '../crypto/vault';
import * as CryptoJS from 'crypto-js';

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: null as string | null,
    username: null as string | null,
    masterSalt: null as string | null,
    isAuthenticated: false,
  }),
  actions: {
    async register(username: string, masterPassword: string): Promise<void> {
      const salt = CryptoJS.lib.WordArray.random(16).toString(CryptoJS.enc.Hex);
      const masterKey = await deriveKey(masterPassword, salt);
      const masterKeyHash = CryptoJS.SHA256(masterKey.toString()).toString(CryptoJS.enc.Hex);

      await api.register({
        username,
        master_key_hash: masterKeyHash,
        master_salt: salt,
      });
    },
    async login(username: string, masterPassword: string): Promise<void> {
      // 步骤 1：从服务器获取盐。
      const saltResponse = await api.getSalt(username);
      const salt = saltResponse.data.master_salt;

      // 步骤 2：派生主密钥并进行哈希。
      const masterKey = await deriveKey(masterPassword, salt);
      const masterKeyHash = CryptoJS.SHA256(masterKey.toString()).toString(CryptoJS.enc.Hex);

      // 步骤 3：使用用户名和主密钥哈希调用登录 API。
      const loginResponse = await api.login({
        username,
        master_key_hash: masterKeyHash,
      });

      // 步骤 4：在 store 中设置认证数据。
      this.setAuthData(loginResponse.data.token, username, loginResponse.data.master_salt);
    },
    setAuthData(token: string, username: string, masterSalt: string) {
      this.token = token;
      this.username = username;
      this.masterSalt = masterSalt;
      this.isAuthenticated = true;
    },
    clearAuthData() {
      this.token = null;
      this.username = null;
      this.masterSalt = null;
      this.isAuthenticated = false;
    },
  },
  persist: true,
});