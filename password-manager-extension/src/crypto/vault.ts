// 此实现使用浏览器原生的 `SubtleCrypto` API，
// 这是现代安全 Web 应用程序的推荐标准。
// 它在单独的线程中运行，提供更好的性能，并使密钥
// 更难提取，从而提供优于纯 JS 库的安全性。

// --- 辅助函数 ---

/**
 * 使用 UTF-8 编码将字符串转换为 ArrayBuffer。
 * @param str 要转换的字符串。
 * @returns 生成的 ArrayBuffer。
 */
function strToArrBuf(str: string): ArrayBuffer {
    return new TextEncoder().encode(str).slice().buffer;
}

/**
 * 使用 UTF-8 解码将 ArrayBuffer 转换为字符串。
 * @param buffer 要转换的 ArrayBuffer。
 * @returns 生成的字符串。
 */
function arrBufToStr(buffer: ArrayBuffer): string {
    return new TextDecoder().decode(buffer);
}


/**
 * 将 ArrayBuffer 转换为 Base64 字符串。
 * @param buffer 要转换的 ArrayBuffer。
 * @returns 生成的 Base64 字符串。
 */
function arrBufToBase64(buffer: ArrayBuffer): string {
    let binary = '';
    const bytes = new Uint8Array(buffer);
    for (let i = 0; i < bytes.length; i++) {
        binary += String.fromCharCode(bytes[i]);
    }
    return window.btoa(binary);
}

/**
 * 将 Base64 字符串转换为 ArrayBuffer。
 * @param base64 要转换的 Base64 字符串。
 * @returns 生成的 ArrayBuffer。
 */
function base64ToArrBuf(base64: string): ArrayBuffer {
    const binary_string = window.atob(base64);
    const len = binary_string.length;
    const bytes = new Uint8Array(len);
    for (let i = 0; i < len; i++) {
        bytes[i] = binary_string.charCodeAt(i);
    }
    return bytes.buffer;
}

// --- 辅助函数 ---

/**
 * 将 ArrayBuffer 转换为十六进制字符串。
 * @param buffer 要转换的 ArrayBuffer。
 * @returns 生成的十六进制字符串。
 */
function arrBufToHex(buffer: ArrayBuffer): string {
    return Array.from(new Uint8Array(buffer))
        .map(b => b.toString(16).padStart(2, '0'))
        .join('');
}

/**
 * 将十六进制字符串转换为 ArrayBuffer。
 * @param hex 要转换的十六进制字符串。
 * @returns 生成的 ArrayBuffer。
 */
function hexToArrBuf(hex: string): ArrayBuffer {
    if (hex.length % 2 !== 0) {
        throw new Error("Invalid hex string");
    }
    const buffer = new ArrayBuffer(hex.length / 2);
    const view = new DataView(buffer);
    for (let i = 0; i < hex.length; i += 2) {
        view.setUint8(i / 2, parseInt(hex.substring(i, i + 2), 16));
    }
    return buffer;
}

// --- 核心加密函数 ---

/**
 * 使用 PBKDF2 从主密码和盐派生出安全的加密密钥。
 * @param masterPassword - 用户的master-password。
 * @param salt - 与用户关联的盐，从后端获取。
 * @returns {Promise<CryptoKey>} - 适用于 AES-GCM 加密和解密的 CryptoKey。
 */
export async function deriveKey(masterPassword: string, salt: string): Promise<CryptoKey> {
    const masterKey = await window.crypto.subtle.importKey(
        "raw",
        strToArrBuf(masterPassword),
        { name: "PBKDF2" },
        false,
        ["deriveKey"]
    );
    const derivedKey = await window.crypto.subtle.deriveKey(
        {
            name: "PBKDF2",
            salt: hexToArrBuf(salt), // 使用 hexToArrBuf 修复了盐的处理
            iterations: 100000,
            hash: "SHA-256",
        },
        masterKey,
        { name: "AES-GCM", length: 256 },
        true, // 使密钥可导出，以便可以对其进行哈希。
        ["encrypt", "decrypt"]
    );

    return derivedKey;
}

/**
 * 使用 AES-GCM 加密保险库项目对象。
 * @param itemObject - 要加密的明文对象。
 * @param masterPassword - 用户的master-password。
 * @param salt - 用户的盐。
 * @returns {Promise<string>} - 包含 IV 和加密数据的 Base64 编码字符串。
 */
export async function encryptVaultItem(itemObject: object, masterPassword: string, salt: string): Promise<string> {
    const key = await deriveKey(masterPassword, salt);
    const plaintext = JSON.stringify(itemObject);

    // AES-GCM 每次加密都需要一个唯一的初始化向量 (IV)。
    // 建议 GCM 使用 12 字节（96 位）的大小。
    const iv = window.crypto.getRandomValues(new Uint8Array(12));

    const encryptedData = await window.crypto.subtle.encrypt(
        {
            name: "AES-GCM",
            iv: iv,
        },
        key,
        strToArrBuf(plaintext)
    );

    // 将 IV 和加密数据（密文）合并到单个 ArrayBuffer 中。
    // 格式：[IV (12 字节)][密文]
    const combined = new Uint8Array(iv.length + encryptedData.byteLength);
    combined.set(iv, 0);
    combined.set(new Uint8Array(encryptedData), iv.length);

    // 以 Base64 字符串形式返回组合数据，以便于存储和传输。
    return arrBufToBase64(combined.buffer);
}

/**
 * 解密加密的保险库项目数据字符串。
 * @param encryptedDataB64 - 来自后端的 Base64 编码字符串。
 * @param masterPassword - 用户的master-password。
 * @param salt - 用户的盐。
 * @returns {Promise<object>} - 解密的明文对象。
 */
export async function decryptVaultData(encryptedDataB64: string, masterPassword: string, salt: string): Promise<object> {
    const key = await deriveKey(masterPassword, salt);
    
    const combinedData = base64ToArrBuf(encryptedDataB64);

    // IV 是组合数据的前 12 个字节。
    const iv = combinedData.slice(0, 12);
    const ciphertext = combinedData.slice(12);

    try {
        const decryptedBytes = await window.crypto.subtle.decrypt(
            {
                name: "AES-GCM",
                iv: iv,
            },
            key,
            ciphertext
        );

        const plaintext = arrBufToStr(decryptedBytes);
        return JSON.parse(plaintext);
    } catch (error) {
        console.error("Decryption failed:", error);
        // 这个错误是致命的。这意味着主密码错误、数据已损坏或代码中存在错误。
        throw new Error("Decryption failed. Invalid master password or corrupted data.");
    }
}

/**
 * 生成一个安全的随机盐。
 * @returns {string} - 16字节（128位）盐的十六进制编码字符串。
 */
export function generateSalt(): string {
    const salt = window.crypto.getRandomValues(new Uint8Array(16));
    return arrBufToHex(salt.buffer);
}

/**
 * 计算 CryptoKey 的 SHA-256 哈希值。
 * @param key - 要哈希的 CryptoKey。
 * @returns {Promise<string>} - 哈希值的十六进制编码字符串。
 */
export async function hashKey(key: CryptoKey): Promise<string> {
    // 请注意，AES-GCM 密钥默认情况下是不可导出的。
    // 为了对密钥进行哈希，我们需要在派生它时使其可导出。
    // `deriveKey` 函数需要更新以包含 `extractable: true`。
    const keyData = await window.crypto.subtle.exportKey("raw", key);
    const hashBuffer = await window.crypto.subtle.digest("SHA-256", keyData);
    return arrBufToHex(hashBuffer);
}