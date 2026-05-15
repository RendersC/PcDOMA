import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { Link } from 'react-router-dom';
import { Clock, CheckCircle, XCircle, Truck, Monitor, AlertCircle } from 'lucide-react';
import { bookingsApi } from '../api/bookings';
import type { Booking, BookingStatus } from '../types';
import toast from 'react-hot-toast';

const STATUS_CONFIG: Record<BookingStatus, { label: string; color: string; icon: React.ReactNode }> = {
  pending_payment: { label: 'Ожидает оплаты', color: 'bg-yellow-100 text-yellow-700', icon: <Clock className="w-3.5 h-3.5" /> },
  paid: { label: 'Оплачено', color: 'bg-blue-100 text-blue-700', icon: <CheckCircle className="w-3.5 h-3.5" /> },
  pending_worker: { label: 'Ожидает работника', color: 'bg-orange-100 text-orange-700', icon: <Clock className="w-3.5 h-3.5" /> },
  accepted: { label: 'Принято', color: 'bg-green-100 text-green-700', icon: <CheckCircle className="w-3.5 h-3.5" /> },
  rejected: { label: 'Отклонено', color: 'bg-red-100 text-red-700', icon: <XCircle className="w-3.5 h-3.5" /> },
  delivering: { label: 'Доставляется', color: 'bg-purple-100 text-purple-700', icon: <Truck className="w-3.5 h-3.5" /> },
  active: { label: 'Активно', color: 'bg-green-100 text-green-700', icon: <Monitor className="w-3.5 h-3.5" /> },
  completed: { label: 'Завершено', color: 'bg-gray-100 text-gray-600', icon: <CheckCircle className="w-3.5 h-3.5" /> },
  cancelled: { label: 'Отменено', color: 'bg-gray-100 text-gray-500', icon: <XCircle className="w-3.5 h-3.5" /> },
};

const ACTIVE_STATUSES: BookingStatus[] = ['pending_payment', 'paid', 'pending_worker', 'accepted', 'delivering', 'active'];
const PAST_STATUSES: BookingStatus[] = ['completed', 'rejected', 'cancelled'];

function BookingCard({ booking, onCancel }: { booking: Booking; onCancel: (id: string) => void }) {
  const cfg = STATUS_CONFIG[booking.status] || { label: booking.status, color: 'bg-gray-100 text-gray-600', icon: null };
  const canCancel = ['pending_payment', 'paid', 'pending_worker'].includes(booking.status);

  return (
    <div className="bg-white rounded-2xl border border-gray-100 p-5">
      <div className="flex items-start justify-between mb-3">
        <div>
          <p className="text-xs text-gray-400 mb-1">#{booking.id.slice(-8).toUpperCase()}</p>
          <span className={`inline-flex items-center gap-1 text-xs font-medium px-2.5 py-1 rounded-full ${cfg.color}`}>
            {cfg.icon}
            {cfg.label}
          </span>
        </div>
        <div className="text-right">
          <p className="font-bold text-primary-600">{booking.total_price.toLocaleString()} ₸</p>
          <p className="text-xs text-gray-400">{booking.rental_type === 'hourly' ? 'Почасовой' : 'Посуточный'}</p>
        </div>
      </div>

      <div className="grid grid-cols-2 gap-2 mb-3 text-sm">
        <div>
          <p className="text-gray-400 text-xs">Начало</p>
          <p className="font-medium text-gray-800">{new Date(booking.start_time).toLocaleString('ru-RU', { day: 'numeric', month: 'short', hour: '2-digit', minute: '2-digit' })}</p>
        </div>
        <div>
          <p className="text-gray-400 text-xs">Конец</p>
          <p className="font-medium text-gray-800">{new Date(booking.end_time).toLocaleString('ru-RU', { day: 'numeric', month: 'short', hour: '2-digit', minute: '2-digit' })}</p>
        </div>
      </div>

      {booking.rejection_reason && (
        <div className="flex items-start gap-2 bg-red-50 rounded-xl p-3 mb-3">
          <AlertCircle className="w-4 h-4 text-red-500 flex-shrink-0 mt-0.5" />
          <p className="text-sm text-red-700">{booking.rejection_reason}</p>
        </div>
      )}

      <div className="flex gap-2">
        {booking.status === 'pending_payment' && booking.payment_id && (
          <Link
            to={`/payment/${booking.id}`}
            className="flex-1 bg-primary-600 hover:bg-primary-700 text-white text-sm py-2 rounded-xl font-medium text-center transition-colors"
          >
            Оплатить
          </Link>
        )}
        {canCancel && (
          <button
            onClick={() => onCancel(booking.id)}
            className="flex-1 border border-red-200 text-red-600 hover:bg-red-50 text-sm py-2 rounded-xl font-medium transition-colors"
          >
            Отменить
          </button>
        )}
      </div>
    </div>
  );
}

export function MyBookingsPage() {
  const [tab, setTab] = useState<'active' | 'past' | 'all'>('active');
  const queryClient = useQueryClient();

  const { data, isLoading } = useQuery({
    queryKey: ['my-bookings'],
    queryFn: () => bookingsApi.list(),
  });

  const cancelMutation = useMutation({
    mutationFn: (id: string) => bookingsApi.cancel(id),
    onSuccess: () => {
      toast.success('Бронирование отменено');
      queryClient.invalidateQueries({ queryKey: ['my-bookings'] });
    },
    onError: () => toast.error('Не удалось отменить'),
  });

  const allBookings = data?.data?.data || [];
  const bookings =
    tab === 'active'
      ? allBookings.filter((b) => ACTIVE_STATUSES.includes(b.status))
      : tab === 'past'
      ? allBookings.filter((b) => PAST_STATUSES.includes(b.status))
      : allBookings;

  const tabs = [
    { key: 'active' as const, label: 'Активные', count: allBookings.filter((b) => ACTIVE_STATUSES.includes(b.status)).length },
    { key: 'past' as const, label: 'История', count: allBookings.filter((b) => PAST_STATUSES.includes(b.status)).length },
    { key: 'all' as const, label: 'Все', count: allBookings.length },
  ];

  return (
    <div className="max-w-3xl mx-auto px-4 py-8 sm:px-6">
      <h1 className="text-2xl font-bold text-gray-900 mb-6">Мои аренды</h1>

      {/* Tabs */}
      <div className="flex gap-1 bg-gray-100 p-1 rounded-xl mb-6">
        {tabs.map((t) => (
          <button
            key={t.key}
            onClick={() => setTab(t.key)}
            className={`flex-1 py-2 rounded-lg text-sm font-medium transition-all ${
              tab === t.key ? 'bg-white text-gray-900 shadow-sm' : 'text-gray-500 hover:text-gray-700'
            }`}
          >
            {t.label}
            {t.count > 0 && (
              <span className={`ml-1.5 text-xs px-1.5 py-0.5 rounded-full ${
                tab === t.key ? 'bg-primary-100 text-primary-700' : 'bg-gray-200 text-gray-500'
              }`}>{t.count}</span>
            )}
          </button>
        ))}
      </div>

      {isLoading ? (
        <div className="space-y-4">
          {[...Array(3)].map((_, i) => (
            <div key={i} className="bg-white rounded-2xl border border-gray-100 p-5 animate-pulse">
              <div className="h-4 bg-gray-200 rounded w-1/4 mb-3" />
              <div className="h-3 bg-gray-200 rounded w-1/2 mb-2" />
              <div className="h-3 bg-gray-200 rounded w-1/3" />
            </div>
          ))}
        </div>
      ) : bookings.length === 0 ? (
        <div className="text-center py-16">
          <p className="text-gray-400 text-lg mb-2">Аренд нет</p>
          <Link to="/catalog" className="text-primary-600 hover:text-primary-700 text-sm font-medium">
            Перейти в каталог →
          </Link>
        </div>
      ) : (
        <div className="space-y-4">
          {bookings.map((b) => (
            <BookingCard key={b.id} booking={b} onCancel={(id) => cancelMutation.mutate(id)} />
          ))}
        </div>
      )}
    </div>
  );
}
