import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { Bell, CheckCheck } from 'lucide-react';
import { notificationsApi } from '../api/notifications';
import { useNotificationStore } from '../store/notificationStore';
import type { Notification } from '../types';

function NotificationItem({ n, onRead }: { n: Notification; onRead: (id: string) => void }) {
  return (
    <div
      className={`p-4 rounded-xl border transition-all cursor-pointer ${
        n.is_read ? 'bg-white border-gray-100' : 'bg-primary-50 border-primary-100'
      }`}
      onClick={() => !n.is_read && onRead(n.id)}
    >
      <div className="flex items-start justify-between gap-3">
        <div className="flex-1">
          <p className={`text-sm font-medium ${n.is_read ? 'text-gray-700' : 'text-gray-900'}`}>{n.title}</p>
          <p className="text-sm text-gray-500 mt-0.5">{n.body}</p>
          <p className="text-xs text-gray-400 mt-1">
            {new Date(n.created_at).toLocaleString('ru-RU', { day: 'numeric', month: 'short', hour: '2-digit', minute: '2-digit' })}
          </p>
        </div>
        {!n.is_read && <div className="w-2 h-2 bg-primary-500 rounded-full flex-shrink-0 mt-1.5" />}
      </div>
    </div>
  );
}

export function NotificationsPage() {
  const { setNotifications, markRead, markAllRead } = useNotificationStore();
  const queryClient = useQueryClient();

  const { data, isLoading } = useQuery({
    queryKey: ['notifications'],
    queryFn: async () => {
      const res = await notificationsApi.list();
      setNotifications(res.data.data || []);
      return res;
    },
  });

  const readMutation = useMutation({
    mutationFn: (id: string) => notificationsApi.markRead(id),
    onSuccess: (_data, id) => {
      markRead(id);
      queryClient.invalidateQueries({ queryKey: ['notifications'] });
    },
  });

  const readAllMutation = useMutation({
    mutationFn: () => notificationsApi.markAllRead(),
    onSuccess: () => {
      markAllRead();
      queryClient.invalidateQueries({ queryKey: ['notifications'] });
    },
  });

  const notifications = data?.data?.data || [];
  const unread = notifications.filter((n) => !n.is_read).length;

  return (
    <div className="max-w-2xl mx-auto px-4 py-8 sm:px-6">
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Уведомления</h1>
          {unread > 0 && <p className="text-sm text-gray-400 mt-0.5">{unread} непрочитанных</p>}
        </div>
        {unread > 0 && (
          <button
            onClick={() => readAllMutation.mutate()}
            className="flex items-center gap-1.5 text-sm text-primary-600 hover:text-primary-700 font-medium"
          >
            <CheckCheck className="w-4 h-4" /> Прочитать все
          </button>
        )}
      </div>

      {isLoading ? (
        <div className="space-y-3">
          {[...Array(4)].map((_, i) => (
            <div key={i} className="bg-white rounded-xl border border-gray-100 p-4 animate-pulse">
              <div className="h-4 bg-gray-200 rounded w-3/4 mb-2" />
              <div className="h-3 bg-gray-200 rounded w-full mb-2" />
              <div className="h-3 bg-gray-200 rounded w-1/4" />
            </div>
          ))}
        </div>
      ) : notifications.length === 0 ? (
        <div className="text-center py-16">
          <Bell className="w-12 h-12 text-gray-200 mx-auto mb-3" />
          <p className="text-gray-400">Уведомлений нет</p>
        </div>
      ) : (
        <div className="space-y-2">
          {notifications.map((n) => (
            <NotificationItem key={n.id} n={n} onRead={(id) => readMutation.mutate(id)} />
          ))}
        </div>
      )}
    </div>
  );
}
