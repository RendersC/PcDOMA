import { api } from './client';
import type { Payment, ListResponse } from '../types';

export const paymentsApi = {
  list: () => api.get<ListResponse<Payment>>('/payments'),

  getById: (id: string) => api.get<Payment>(`/payments/${id}`),

  process: (id: string, method: 'card' | 'kaspi') =>
    api.post<Payment>(`/payments/${id}/process`, { method }),
};
