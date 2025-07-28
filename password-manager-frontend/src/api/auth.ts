import apiClient from './index';
import type { RegisterRequestPayload, LoginRequestPayload } from '@/types';

export const register = (payload: RegisterRequestPayload) => {
  return apiClient.post('/auth/register', payload);
};

export const login = (payload: LoginRequestPayload) => {
  return apiClient.post('/auth/login', payload);
};

export const getSalt = (identifier: string) => {
  return apiClient.post('/auth/salt', { identifier });
};

export const sendVerificationCode = (data: { email: string }) => {
  return apiClient.post('/auth/send-verification-code', data);
};

export const requestPasswordReset = (email: string) => {
  return apiClient.post('/auth/request-password-reset', { email });
};

export const resetPassword = (token: string, new_master_key_hash: string, new_master_salt: string) => {
  return apiClient.post('/auth/reset-password', { token, new_master_key_hash, new_master_salt });
};