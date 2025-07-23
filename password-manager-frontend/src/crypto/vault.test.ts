import { describe, it, expect } from 'vitest';
import { encryptVaultItem, decryptVaultData } from './vault';

// 辅助函数：将 Base64 转换为 ArrayBuffer
function base64ToArrBuf(base64: string): ArrayBuffer {
    const binary_string = window.atob(base64);
    const len = binary_string.length;
    const bytes = new Uint8Array(len);
    for (let i = 0; i < len; i++) {
        bytes[i] = binary_string.charCodeAt(i);
    }
    return bytes.buffer;
}

// 辅助函数：将 ArrayBuffer 转换为 Base64
function arrBufToBase64(buffer: ArrayBuffer): string {
    let binary = '';
    const bytes = new Uint8Array(buffer);
    for (let i = 0; i < bytes.length; i++) {
        binary += String.fromCharCode(bytes[i]);
    }
    return window.btoa(binary);
}


describe('Vault Encryption and Decryption with SubtleCrypto', () => {
  const masterPassword = 'my-super-secret-password';
  const salt = 'a-very-salty-salt';
  const vaultItem = {
    name: 'Test Item',
    website: 'example.com',
    login: 'testuser',
    password: 'testpassword',
  };

  it('should correctly encrypt and decrypt a vault item', async () => {
    // 加密项目
    const encryptedData = await encryptVaultItem(vaultItem, masterPassword, salt);

    // 确保加密后的数据是一个非空的字符串 (Base64)
    expect(typeof encryptedData).toBe('string');
    expect(encryptedData.length).toBeGreaterThan(0);

    // 解密项目
    const decryptedItem = await decryptVaultData(encryptedData, masterPassword, salt);

    // 确保解密后的项目与原始项目匹配
    expect(decryptedItem).toEqual(vaultItem);
  }, 20000);

  it('should throw an error when decrypting with an incorrect master password', async () => {
    const wrongPassword = 'wrong-password';

    // 使用正确的密码加密项目
    const encryptedData = await encryptVaultItem(vaultItem, masterPassword, salt);

    // 期望使用错误的密码解密会抛出错误
    await expect(
      decryptVaultData(encryptedData, wrongPassword, salt)
    ).rejects.toThrow('Decryption failed. Invalid master password or corrupted data.');
  }, 20000);

  it('should throw an error when decrypting corrupted data (integrity check)', async () => {
    // 加密项目
    const encryptedData = await encryptVaultItem(vaultItem, masterPassword, salt);

    // 通过翻转一个位来损坏数据
    const buffer = base64ToArrBuf(encryptedData);
    const corruptedBytes = new Uint8Array(buffer);
    // 在密文中间的某个位置篡改一个字节
    const tamperedIndex = Math.floor(corruptedBytes.length / 2);
    corruptedBytes[tamperedIndex] = corruptedBytes[tamperedIndex] ^ 1; // 翻转一个位
    const corruptedBase64 = arrBufToBase64(corruptedBytes.buffer);

    // 期望使用损坏的数据解密会抛出错误
    await expect(
      decryptVaultData(corruptedBase64, masterPassword, salt)
    ).rejects.toThrow('Decryption failed. Invalid master password or corrupted data.');
  }, 20000);


  it('should handle different vault item structures', async () => {
    const anotherItem = {
      service: 'Another Service',
      username: 'anotheruser',
      secret: 'another-secret-key-12345',
      notes: 'Some important notes here.',
    };

    const encryptedData = await encryptVaultItem(anotherItem, masterPassword, salt);
    const decryptedItem = await decryptVaultData(encryptedData, masterPassword, salt);

    expect(decryptedItem).toEqual(anotherItem);
  }, 20000);
});