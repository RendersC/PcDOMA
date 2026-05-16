import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { Link } from 'react-router-dom';
import { Package, CheckCircle, XCircle, Clock, RefreshCw } from 'lucide-react';
import { bookingsApi } from '../../api/bookings';
import type { Booking } from '../../types';
import toast from 'react-hot-toast';

function BookingCard({ booking, onAccept, onReject }: {
  booking: Booking;
  onAccept: (id: string) => void;
  onReject: (id: string, reason: string) => void;
}) {
  const [showReject, setShowReject] = useState(false);
  const [reason, setReason] = useState('');

  return (
    <div className="bg-white rounded-2xl border-2 border-orange-100 p-5">
      <div className="flex items-start justify-between mb-3">
        <div>
          <p className="text-xs text-gray-400 mb-1">Заказ #{booking.id.slice(-8).toUpperCase()}</p>
          <div className="flex items-center gap-1.5 text-orange-600 text-sm font-medium">
            <Clock className="w-4 h-4" />
            Ожидает принятия
          </div>
        </div>
        <p className="font-bold text-primary-600 text-lg">{booking.total_price.toLocaleString()} ₸</p>
      </div>

      <div className="grid grid-cols-2 gap-2 mb-3 text-sm">
        <div>
          <p className="text-gray-400 text-xs">Тип</p>
          <p className="font-medium">{booking.booking_type === 'custom' ? 'Конструктор' : 'Готовый сеап'}</p>
        </div>
        <div>
          <p className="text-gray-400 text-xs">Тариф</p>
          <p className="font-medium">{booking.rental_type === 'hourly' ? 'Почасовой' : 'Посуточный'}</p>
        </div>
        <div>
          <p className="text-gray-400 text-xs">Начало</p>
          <p className="font-medium">{new Date(booking.start_time).toLocaleString('ru-RU', { day: 'numeric', month: 'short', hour: '2-digit', minute: '2-digit' })}</p>
        </div>
        <div>
          <p className="text-gray-400 text-xs">Конец</p>
          <p className="font-medium">{new Date(booking.end_time).toLocaleString('ru-RU', { day: 'numeric', month: 'short', hour: '2-digit', minute: '2-digit' })}</p>
        </div>
      </div>

      {booking.peripherals && booking.peripherals.length > 0 && (
        <div className="bg-gray-50 rounded-xl p-3 mb-3">
          <p className="text-xs text-gray-500 mb-1">Периферия: {booking.peripherals.length} шт.</p>
        </div>
      )}

      <div className="flex gap-2">
        <Link
          to={`/worker/order/${booking.id}`}
          className="flex-1 border border-gray-200 text-gray-600 hover:bg-gray-50 text-sm py-2.5 rounded-xl font-medium text-center transition-colors"
        >
          Детали
        </Link>
        {!showReject ? (
          <>
            <button
              onClick={() => onAccept(booking.id)}
              className="flex-1 bg-green-500 hover:bg-green-600 text-white text-sm py-2.5 rounded-xl font-medium flex items-center justify-center gap-1.5 transition-colors"
            >
              <CheckCircle className="w-4 h-4" /> Принять
            </button>
            <button
              onClick={() => setShowReject(true)}
              className="flex-1 bg-red-50 hover:bg-red-100 text-red-600 text-sm py-2.5 rounded-xl font-medium flex items-center justify-center gap-1.5 transition-colors"
            >
              <XCircle className="w-4 h-4" /> Отклонить
            </button>
          </>
        ) : (
          <div className="flex-1 space-y-2">
            <input
              type="text"
              value={reason}
              onChange={(e) => setReason(e.target.value)}
              placeholder="Причина отказа..."
              className="w-full px-3 py-2 border border-gray-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-red-300"
            />
            <div className="flex gap-2">
              <button
                onClick={() => { onReject(booking.id, reason); setShowReject(false); }}
                className="flex-1 bg-red-500 hover:bg-red-600 text-white text-sm py-2 rounded-xl font-medium transition-colors"
              >
                Подтвердить
              </button>
              <button
                onClick={() => { setShowReject(false); setReason(''); }}
                className="flex-1 border border-gray-200 text-gray-600 text-sm py-2 rounded-xl font-medium transition-colors"
              >
                Отмена
              </button>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}

export function WorkerDashboard() {
  const queryClient = useQueryClient();

  const { data, isLoading, refetch, isFetching } = useQuery({
    queryKey: ['worker-bookings'],
    queryFn: () => bookingsApi.workerList(),
    refetchInterval: 10000,
  });

  const acceptMutation = useMutation({
    mutationFn: (id: string) => bookingsApi.workerAccept(id),
    onSuccess: () => {
      toast.success('Заказ принят');
      queryClient.invalidateQueries({ queryKey: ['worker-bookings'] });
    },
    onError: () => toast.error('Не удалось принять'),
  });

  const rejectMutation = useMutation({
    mutationFn: ({ id, reason }: { id: string; reason: string }) => bookingsApi.workerReject(id, reason),
    onSuccess: () => {
      toast.success('Заказ отклонён');
      queryClient.invalidateQueries({ queryKey: ['worker-bookings'] });
    },
    onError: () => toast.error('Не удалось отклонить'),
  });

  const bookings = data?.data?.data || [];
  const pending = bookings.filter((b) => b.status === 'pending_worker');

  return (
    <div className="max-w-3xl mx-auto px-4 py-8 sm:px-6">
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Новые заказы</h1>
          <p className="text-sm text-gray-400 mt-0.5">Обновляется каждые 10 секунд</p>
        </div>
        <div className="flex items-center gap-3">
          <Link
            to="/worker/active"
            className="text-sm text-primary-600 hover:text-primary-700 font-medium"
          >
            Активные аренды →
          </Link>
          <button
            onClick={() => refetch()}
            className={`p-2 text-gray-400 hover:text-gray-600 ${isFetching ? 'animate-spin' : ''}`}
          >
            <RefreshCw className="w-4 h-4" />
          </button>
        </div>
      </div>

      {isLoading ? (
        <div className="space-y-4">
          {[...Array(3)].map((_, i) => (
            <div key={i} className="bg-white rounded-2xl border border-gray-100 p-5 animate-pulse">
              <div className="h-4 bg-gray-200 rounded w-1/3 mb-3" />
              <div className="h-3 bg-gray-200 rounded w-1/2 mb-2" />
              <div className="h-3 bg-gray-200 rounded w-2/3" />
            </div>
          ))}
        </div>
      ) : pending.length === 0 ? (
        <div className="text-center py-16">
          <Package className="w-12 h-12 text-gray-200 mx-auto mb-3" />
          <p className="text-gray-400 text-lg mb-1">Новых заказов нет</p>
          <p className="text-gray-400 text-sm">Страница обновляется автоматически</p>
        </div>
      ) : (
        <div className="space-y-4">
          {pending.map((b) => (
            <BookingCard
              key={b.id}
              booking={b}
              onAccept={(id) => acceptMutation.mutate(id)}
              onReject={(id, reason) => rejectMutation.mutate({ id, reason })}
            />
          ))}
        </div>
      )}
    </div>
  );
}
