import apiClient from './index';

export const getVault = () => {
  return apiClient.get('/vault/items');
};

export const addVaultItem = (item: { encrypted_data: string; category: string }) => {
  return apiClient.post('/vault/items', item);
};

export const updateVaultItem = (id: string, item: { encrypted_data: string; category: string }) => {
  return apiClient.put(`/vault/items/${id}`, item);
};

export const deleteVaultItem = (id: string) => {
  return apiClient.delete(`/vault/${id}`);
};