import apiClient from './index';
import type { UserCredentials } from '@/types';

export const register = (credentials: UserCredentials) => {
  return apiClient.post('/auth/register', credentials);
};

export const login = (credentials: UserCredentials) => {
  return apiClient.post('/auth/login', credentials);
};

export const getSalt = (username: string) => {
  return apiClient.get(`/auth/salt/${username}`);
};