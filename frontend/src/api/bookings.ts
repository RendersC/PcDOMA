import { api } from './client';
import type { Booking, ListResponse } from '../types';

export const bookingsApi = {
  create: (data: {
    booking_type: 'custom' | 'setup';
    pc_id?: string;
    setup_id?: string;
    peripheral_ids?: string[];
    start_time: string;
    end_time: string;
    rental_type: 'hourly' | 'daily';
  }) => api.post<Booking>('/bookings', data),

  list: (params?: { status?: string }) =>
    api.get<ListResponse<Booking>>('/bookings', { params }),

  getById: (id: string) =>
    api.get<Booking>(`/bookings/${id}`),

  cancel: (id: string) =>
    api.delete(`/bookings/${id}`),

  // Worker
  workerList: (locationId?: string) =>
    api.get<ListResponse<Booking>>('/worker/bookings', { params: { location_id: locationId } }),

  workerActiveList: () =>
    api.get<ListResponse<Booking>>('/worker/bookings/active'),

  workerAccept: (id: string) =>
    api.patch<Booking>(`/worker/bookings/${id}/accept`),

  workerReject: (id: string, reason?: string) =>
    api.patch<Booking>(`/worker/bookings/${id}/reject`, { rejection_reason: reason }),

  workerDeliver: (id: string) =>
    api.patch<Booking>(`/worker/bookings/${id}/delivering`),

  workerComplete: (id: string) =>
    api.patch<Booking>(`/worker/bookings/${id}/complete`),

  // Admin
  adminList: (params?: Record<string, string>) =>
    api.get<ListResponse<Booking>>('/admin/bookings', { params }),
};
