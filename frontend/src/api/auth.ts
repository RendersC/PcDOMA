import { api } from './client';
import type { AuthResponse, User } from '../types';

export const authApi = {
  register: (data: { name: string; email: string; password: string; phone?: string }) =>
    api.post<AuthResponse>('/auth/register', data),

  login: (data: { email: string; password: string }) =>
    api.post<AuthResponse>('/auth/login', data),

  logout: (refreshToken: string) =>
    api.post('/auth/logout', { refresh_token: refreshToken }),

  me: () => api.get<User>('/auth/me'),

  changePassword: (id: string, data: { old_password: string; new_password: string }) =>
    api.put(`/users/${id}/password`, data),

  updateProfile: (id: string, data: { name?: string; phone?: string }) =>
    api.put<User>(`/users/${id}`, data),
};
