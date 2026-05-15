import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { bookingsApi } from '../../api/bookings';
import type { BookingStatus } from '../../types';

const STATUS_COLORS: Record<BookingStatus, string> = {
  pending_payment: 'bg-yellow-100 text-yellow-700',
  paid: 'bg-blue-100 text-blue-700',
  pending_worker: 'bg-orange-100 text-orange-700',
  accepted: 'bg-green-100 text-green-700',
  rejected: 'bg-red-100 text-red-700',
  delivering: 'bg-purple-100 text-purple-700',
  active: 'bg-teal-100 text-teal-700',
  completed: 'bg-gray-100 text-gray-600',
  cancelled: 'bg-gray-100 text-gray-400',
};

const STATUS_LABELS: Record<BookingStatus, string> = {
  pending_payment: 'Ожидает оплаты',
  paid: 'Оплачено',
  pending_worker: 'Ожидает работника',
  accepted: 'Принято',
  rejected: 'Отклонено',
  delivering: 'Доставляется',
  active: 'Активно',
  completed: 'Завершено',
  cancelled: 'Отменено',
};

const ALL_STATUSES: BookingStatus[] = [
  'pending_payment', 'paid', 'pending_worker', 'accepted',
  'rejected', 'delivering', 'active', 'completed', 'cancelled',
];

export function AdminBookings() {
  const [statusFilter, setStatusFilter] = useState('');

  const { data, isLoading } = useQuery({
    queryKey: ['admin-bookings', statusFilter],
    queryFn: () => bookingsApi.adminList(statusFilter ? { status: statusFilter, limit: '100' } : { limit: '100' }),
  });

  const bookings = data?.data?.data || [];

  return (
    <div className="max-w-7xl mx-auto px-4 py-8 sm:px-6">
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold text-gray-900">Все бронирования</h1>
        <span className="text-sm text-gray-400">{bookings.length} записей</span>
      </div>

      {/* Status filter */}
      <div className="flex flex-wrap gap-2 mb-6">
        <button
          onClick={() => setStatusFilter('')}
          className={`px-3 py-1.5 rounded-full text-xs font-medium transition-colors ${
            statusFilter === '' ? 'bg-gray-800 text-white' : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
          }`}
        >
          Все
        </button>
        {ALL_STATUSES.map((s) => (
          <button
            key={s}
            onClick={() => setStatusFilter(s)}
            className={`px-3 py-1.5 rounded-full text-xs font-medium transition-colors ${
              statusFilter === s ? 'bg-gray-800 text-white' : `${STATUS_COLORS[s]} hover:opacity-80`
            }`}
          >
            {STATUS_LABELS[s]}
          </button>
        ))}
      </div>

      {isLoading ? (
        <div className="space-y-3">
          {[...Array(5)].map((_, i) => <div key={i} className="bg-white rounded-xl border border-gray-100 h-16 animate-pulse" />)}
        </div>
      ) : (
        <div className="bg-white rounded-2xl border border-gray-100 overflow-hidden">
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead className="bg-gray-50 border-b border-gray-100">
                <tr>
                  <th className="text-left px-4 py-3 font-medium text-gray-600">ID</th>
                  <th className="text-left px-4 py-3 font-medium text-gray-600">Тип</th>
                  <th className="text-left px-4 py-3 font-medium text-gray-600">Начало</th>
                  <th className="text-left px-4 py-3 font-medium text-gray-600">Конец</th>
                  <th className="text-left px-4 py-3 font-medium text-gray-600">Сумма</th>
                  <th className="text-left px-4 py-3 font-medium text-gray-600">Статус</th>
                  <th className="text-left px-4 py-3 font-medium text-gray-600">Создано</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-50">
                {bookings.map((b) => (
                  <tr key={b.id} className="hover:bg-gray-50">
                    <td className="px-4 py-3 font-mono text-xs text-gray-500">#{b.id.slice(-8).toUpperCase()}</td>
                    <td className="px-4 py-3 text-gray-700">{b.booking_type === 'custom' ? 'Конструктор' : 'Сеап'}</td>
                    <td className="px-4 py-3 text-gray-600">{new Date(b.start_time).toLocaleString('ru-RU', { day: 'numeric', month: 'short', hour: '2-digit', minute: '2-digit' })}</td>
                    <td className="px-4 py-3 text-gray-600">{new Date(b.end_time).toLocaleString('ru-RU', { day: 'numeric', month: 'short', hour: '2-digit', minute: '2-digit' })}</td>
                    <td className="px-4 py-3 font-semibold text-gray-900">{b.total_price.toLocaleString()} ₸</td>
                    <td className="px-4 py-3">
                      <span className={`text-xs font-medium px-2 py-1 rounded-full ${STATUS_COLORS[b.status]}`}>
                        {STATUS_LABELS[b.status]}
                      </span>
                    </td>
                    <td className="px-4 py-3 text-gray-400 text-xs">{new Date(b.created_at).toLocaleDateString('ru-RU')}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
          {bookings.length === 0 && (
            <div className="text-center py-12 text-gray-400">Бронирований нет</div>
          )}
        </div>
      )}
    </div>
  );
}
