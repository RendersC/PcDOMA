import { api } from './client';
import type { PC, Peripheral, Setup, Review, ListResponse } from '../types';

export const catalogApi = {
  // PCs
  listPCs: (params?: Record<string, string | number>) =>
    api.get<ListResponse<PC>>('/catalog/pcs', { params }),

  getPC: (id: string) =>
    api.get<PC>(`/catalog/pcs/${id}`),

  createPC: (data: unknown) =>
    api.post<PC>('/catalog/pcs', data),

  updatePC: (id: string, data: unknown) =>
    api.put<PC>(`/catalog/pcs/${id}`, data),

  deletePC: (id: string) =>
    api.delete(`/catalog/pcs/${id}`),

  updatePCStatus: (id: string, status: string) =>
    api.patch(`/catalog/pcs/${id}/status`, { status }),

  getLocations: () =>
    api.get<{ data: string[] }>('/catalog/locations'),

  // Reviews
  listReviews: (pcId: string) =>
    api.get<ListResponse<Review>>(`/catalog/pcs/${pcId}/reviews`),

  createReview: (pcId: string, data: { rating: number; comment: string }) =>
    api.post<Review>(`/catalog/pcs/${pcId}/reviews`, data),

  // Peripherals
  listPeripherals: (params?: Record<string, string>) =>
    api.get<ListResponse<Peripheral>>('/catalog/peripherals', { params }),

  getPeripheral: (id: string) =>
    api.get<Peripheral>(`/catalog/peripherals/${id}`),

  createPeripheral: (data: unknown) =>
    api.post<Peripheral>('/catalog/peripherals', data),

  deletePeripheral: (id: string) =>
    api.delete(`/catalog/peripherals/${id}`),

  // Setups
  listSetups: (params?: Record<string, string>) =>
    api.get<ListResponse<Setup>>('/catalog/setups', { params }),

  getSetup: (id: string) =>
    api.get<Setup>(`/catalog/setups/${id}`),

  createSetup: (data: unknown) =>
    api.post<Setup>('/catalog/setups', data),

  deleteSetup: (id: string) =>
    api.delete(`/catalog/setups/${id}`),
};
