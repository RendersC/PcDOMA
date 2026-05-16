import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { Link } from 'react-router-dom';
import { Truck, CheckCircle, Monitor } from 'lucide-react';
import { bookingsApi } from '../../api/bookings';
import type { Booking } from '../../types';
import toast from 'react-hot-toast';

const STATUS_CONFIG = {
  accepted: { label: 'Принято — подготовьте оборудование', color: 'bg-blue-100 text-blue-700', icon: <Monitor className="w-3.5 h-3.5" /> },
  delivering: { label: 'Доставляется клиенту', color: 'bg-purple-100 text-purple-700', icon: <Truck className="w-3.5 h-3.5" /> },
  active: { label: 'Аренда активна', color: 'bg-green-100 text-green-700', icon: <CheckCircle className="w-3.5 h-3.5" /> },
};

function ActiveCard({ booking, onDeliver, onComplete }: {
  booking: Booking;
  onDeliver: (id: string) => void;
  onComplete: (id: string) => void;
}) {
  const cfg = STATUS_CONFIG[booking.status as keyof typeof STATUS_CONFIG];

  return (
    <div className="bg-white rounded-2xl border border-gray-100 p-5">
      <div className="flex items-start justify-between mb-3">
        <div>
          <p className="text-xs text-gray-400 mb-1">#{booking.id.slice(-8).toUpperCase()}</p>
          {cfg && (
            <span className={`inline-flex items-center gap-1 text-xs font-medium px-2.5 py-1 rounded-full ${cfg.color}`}>
              {cfg.icon}
              {cfg.label}
            </span>
          )}
        </div>
        <p className="font-bold text-primary-600">{booking.total_price.toLocaleString()} ₸</p>
      </div>

      <div className="grid grid-cols-2 gap-2 mb-4 text-sm">
        <div>
          <p className="text-gray-400 text-xs">Начало</p>
          <p className="font-medium">{new Date(booking.start_time).toLocaleString('ru-RU', { day: 'numeric', month: 'short', hour: '2-digit', minute: '2-digit' })}</p>
        </div>
        <div>
          <p className="text-gray-400 text-xs">Конец</p>
          <p className="font-medium">{new Date(booking.end_time).toLocaleString('ru-RU', { day: 'numeric', month: 'short', hour: '2-digit', minute: '2-digit' })}</p>
        </div>
      </div>

      <div className="flex gap-2">
        <Link
          to={`/worker/order/${booking.id}`}
          className="flex-1 border border-gray-200 text-gray-600 hover:bg-gray-50 text-sm py-2.5 rounded-xl font-medium text-center transition-colors"
        >
          Детали заказа
        </Link>
        {booking.status === 'accepted' && (
          <button
            onClick={() => onDeliver(booking.id)}
            className="flex-1 bg-purple-500 hover:bg-purple-600 text-white text-sm py-2.5 rounded-xl font-medium flex items-center justify-center gap-1.5 transition-colors"
          >
            <Truck className="w-4 h-4" /> Выдаю клиенту
          </button>
        )}
        {booking.status === 'delivering' && (
          <button
            onClick={() => onComplete(booking.id)}
            className="flex-1 bg-green-500 hover:bg-green-600 text-white text-sm py-2.5 rounded-xl font-medium flex items-center justify-center gap-1.5 transition-colors"
          >
            <CheckCircle className="w-4 h-4" /> Завершить
          </button>
        )}
      </div>
    </div>
  );
}

export function WorkerActiveOrders() {
  const queryClient = useQueryClient();

  const { data, isLoading } = useQuery({
    queryKey: ['worker-active'],
    queryFn: () => bookingsApi.workerActiveList(),
    refetchInterval: 15000,
  });

  const deliverMutation = useMutation({
    mutationFn: (id: string) => bookingsApi.workerDeliver(id),
    onSuccess: () => {
      toast.success('Статус обновлён: Доставляется');
      queryClient.invalidateQueries({ queryKey: ['worker-active'] });
    },
  });

  const completeMutation = useMutation({
    mutationFn: (id: string) => bookingsApi.workerComplete(id),
    onSuccess: () => {
      toast.success('Аренда завершена');
      queryClient.invalidateQueries({ queryKey: ['worker-active'] });
    },
  });

  const bookings = data?.data?.data || [];

  return (
    <div className="max-w-3xl mx-auto px-4 py-8 sm:px-6">
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold text-gray-900">Активные аренды</h1>
        <Link to="/worker" className="text-sm text-primary-600 hover:text-primary-700 font-medium">
          ← Новые заказы
        </Link>
      </div>

      {isLoading ? (
        <div className="space-y-4">
          {[...Array(2)].map((_, i) => (
            <div key={i} className="bg-white rounded-2xl border border-gray-100 p-5 animate-pulse">
              <div className="h-4 bg-gray-200 rounded w-1/3 mb-3" />
              <div className="h-3 bg-gray-200 rounded w-1/2" />
            </div>
          ))}
        </div>
      ) : bookings.length === 0 ? (
        <div className="text-center py-16">
          <p className="text-gray-400 text-lg mb-1">Активных аренд нет</p>
          <Link to="/worker" className="text-primary-600 text-sm font-medium hover:underline">
            Смотреть новые заказы
          </Link>
        </div>
      ) : (
        <div className="space-y-4">
          {bookings.map((b) => (
            <ActiveCard
              key={b.id}
              booking={b}
              onDeliver={(id) => deliverMutation.mutate(id)}
              onComplete={(id) => completeMutation.mutate(id)}
            />
          ))}
        </div>
      )}
    </div>
  );
}
