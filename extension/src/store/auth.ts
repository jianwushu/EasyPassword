import { defineStore } from 'pinia';
import * as api from '../api/auth';
import { deriveKey, generateSalt, hashKey } from '../crypto/vault';
import { createChromeStorage } from './storage';

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: null as string | null,
    username: null as string | null,
    masterSalt: null as string | null,
    isAuthenticated: false,
  }),
  actions: {
    async register(username: string, email: string, masterPassword: string, code: string): Promise<void> {
      const salt = generateSalt();
      const masterKey = await deriveKey(masterPassword, salt);
      const masterKeyHash = await hashKey(masterKey);

      await api.register({
        username,
        email,
        master_key_hash: masterKeyHash,
        master_salt: salt,
        code,
      });
    },
    async sendVerificationCode(email: string): Promise<void> {
      await api.sendVerificationCode({ email });
    },
    async login(identifier: string, masterPassword: string): Promise<void> {
      // 步骤 1：从服务器获取盐。
      const saltResponse = await api.getSalt(identifier);
      const salt = saltResponse.data.master_salt;

      // 步骤 2：派生主密钥并进行哈希。
      const masterKey = await deriveKey(masterPassword, salt);
      const masterKeyHash = await hashKey(masterKey);

      // 步骤 3：使用标识符和主密钥哈希调用登录 API。
      const loginResponse = await api.login({
        identifier,
        master_key_hash: masterKeyHash,
      });

      // 步骤 4：在 store 中设置认证数据。
      // 注意：这里我们将标识符用作用户名。如果需要显示确切的用户名，
      // 后端应在登录响应中返回它。
      this.setAuthData(loginResponse.data.token, loginResponse.data.username, loginResponse.data.master_salt);
    },
    setAuthData(token: string, username: string, masterSalt: string) {
      this.token = token;
      this.username = username;
      this.masterSalt = masterSalt;
      this.isAuthenticated = true;
      // Manually persist state
      const storage = createChromeStorage();
      storage.setItem('auth', JSON.stringify(this.$state));
    },
    clearAuthData() {
      this.token = null;
      this.username = null;
      this.masterSalt = null;
      this.isAuthenticated = false;
      // Manually clear persisted state
      const storage = createChromeStorage();
      storage.removeItem('auth');
    },
    async initializeAuth(): Promise<void> {
      const storage = createChromeStorage();
      const authDataString = await storage.getItem('auth');
      if (authDataString) {
        const authData = JSON.parse(authDataString);
        if (authData.token) {
          this.token = authData.token;
          this.username = authData.username;
          this.masterSalt = authData.masterSalt;
          this.isAuthenticated = true;
        }
      }
    },
  },
});