// Define the storage type based on the expected interface for pinia-plugin-persistedstate
type Storage = {
  getItem(key: string): Promise<string | null>;
  setItem(key: string, value: string): Promise<void>;
  removeItem(key: string): Promise<void>;
};

/**
 * Creates a storage adapter for pinia-plugin-persistedstate that uses the
 * chrome.storage.local API. This allows Pinia state to be shared across
 * all parts of the extension (background, popup, content scripts).
 *
 * @returns A Storage-like object for Pinia persistence.
 */
export const createChromeStorage = (): Storage => {
  return {
    /**
     * Asynchronously gets an item from chrome.storage.local.
     * @param key - The key of the item to get.
     * @returns A Promise that resolves to the stored value (stringified JSON).
     */
    getItem: async (key: string): Promise<string | null> => {
      const data = await chrome.storage.local.get(key);
      return data[key] || null;
    },

    /**
     * Asynchronously sets an item in chrome.storage.local.
     * @param key - The key of the item to set.
     * @param value - The value of the item to set (stringified JSON).
     * @returns A Promise that resolves when the item is set.
     */
    setItem: async (key: string, value: string): Promise<void> => {
      await chrome.storage.local.set({ [key]: value });
    },

    /**
     * Asynchronously removes an item from chrome.storage.local.
     * @param key - The key of the item to remove.
     * @returns A Promise that resolves when the item is removed.
     */
    removeItem: async (key: string): Promise<void> => {
      await chrome.storage.local.remove(key);
    },
  };
};