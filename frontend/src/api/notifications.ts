import { api } from './client';
import type { Notification, ListResponse } from '../types';

export const notificationsApi = {
  list: (unread?: boolean) =>
    api.get<ListResponse<Notification>>('/notifications', { params: unread ? { unread: true } : {} }),

  markRead: (id: string) =>
    api.patch(`/notifications/${id}/read`),

  markAllRead: () =>
    api.patch('/notifications/read-all'),

  delete: (id: string) =>
    api.delete(`/notifications/${id}`),
};
