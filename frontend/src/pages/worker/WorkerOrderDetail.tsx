import { useParams, useNavigate } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { ArrowLeft, Monitor, Mouse, Headphones, Tv } from 'lucide-react';
import { bookingsApi } from '../../api/bookings';
import { catalogApi } from '../../api/catalog';

const PERIPHERAL_ICONS: Record<string, React.ReactNode> = {
  mouse: <Mouse className="w-4 h-4" />,
  headphones: <Headphones className="w-4 h-4" />,
  headset: <Headphones className="w-4 h-4" />,
  monitor: <Tv className="w-4 h-4" />,
  keyboard: <Monitor className="w-4 h-4" />,
};

const STATUS_LABELS: Record<string, string> = {
  pending_payment: 'Ожидает оплаты',
  paid: 'Оплачено',
  pending_worker: 'Ожидает принятия',
  accepted: 'Принято',
  rejected: 'Отклонено',
  delivering: 'Доставляется',
  active: 'Активно',
  completed: 'Завершено',
  cancelled: 'Отменено',
};

export function WorkerOrderDetail() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();

  const { data: bookingRes, isLoading } = useQuery({
    queryKey: ['booking-detail', id],
    queryFn: () => bookingsApi.getById(id!),
    enabled: !!id,
  });

  const booking = bookingRes?.data;

  const { data: pcRes } = useQuery({
    queryKey: ['pc', booking?.pc_id],
    queryFn: () => catalogApi.getPC(booking!.pc_id!),
    enabled: !!booking?.pc_id,
  });

  const { data: peripheralsRes } = useQuery({
    queryKey: ['peripherals', 'all'],
    queryFn: () => catalogApi.listPeripherals({ limit: '100' }),
    enabled: !!booking?.peripherals?.length,
  });

  const pc = pcRes?.data;
  const allPeripherals = peripheralsRes?.data?.data || [];
  const peripheralIds = booking?.peripherals?.map((p) => p.peripheral_id) || [];
  const peripherals = allPeripherals.filter((p) => peripheralIds.includes(p.id));

  if (isLoading) {
    return (
      <div className="max-w-2xl mx-auto px-4 py-8">
        <div className="animate-pulse space-y-4">
          <div className="h-8 bg-gray-200 rounded w-1/2" />
          <div className="h-40 bg-gray-200 rounded-2xl" />
        </div>
      </div>
    );
  }

  if (!booking) {
    return (
      <div className="max-w-2xl mx-auto px-4 py-16 text-center">
        <p className="text-gray-500">Заказ не найден</p>
      </div>
    );
  }

  return (
    <div className="max-w-2xl mx-auto px-4 py-8 sm:px-6">
      <button onClick={() => navigate(-1)} className="flex items-center gap-2 text-gray-500 hover:text-gray-700 mb-6 text-sm">
        <ArrowLeft className="w-4 h-4" /> Назад
      </button>

      <h1 className="text-2xl font-bold text-gray-900 mb-2">Заказ #{booking.id.slice(-8).toUpperCase()}</h1>
      <p className="text-sm text-gray-500 mb-6">Статус: <strong>{STATUS_LABELS[booking.status] || booking.status}</strong></p>

      {/* What to prepare */}
      <div className="bg-white rounded-2xl border-2 border-primary-100 p-6 mb-5">
        <h2 className="font-bold text-gray-900 mb-4 text-primary-700">Что подготовить клиенту</h2>

        {pc && (
          <div className="flex items-start gap-3 p-3 bg-gray-50 rounded-xl mb-3">
            <Monitor className="w-5 h-5 text-gray-500 flex-shrink-0 mt-0.5" />
            <div>
              <p className="font-semibold text-gray-900">{pc.name}</p>
              <p className="text-xs text-gray-500 mt-0.5">{pc.specs.cpu} · {pc.specs.gpu}</p>
              <p className="text-xs text-gray-400">RAM: {pc.specs.ram} · Storage: {pc.specs.storage}</p>
            </div>
          </div>
        )}

        {peripherals.length > 0 && (
          <div className="space-y-2">
            {peripherals.map((p) => (
              <div key={p.id} className="flex items-center gap-3 p-3 bg-gray-50 rounded-xl">
                <span className="text-gray-500">{PERIPHERAL_ICONS[p.type] || <Monitor className="w-4 h-4" />}</span>
                <div>
                  <p className="font-medium text-sm text-gray-900">{p.name}</p>
                  <p className="text-xs text-gray-400 capitalize">{p.type}</p>
                </div>
              </div>
            ))}
          </div>
        )}

        {booking.booking_type === 'setup' && peripherals.length === 0 && (
          <p className="text-sm text-gray-500">Готовый сеап — уточните состав у администратора</p>
        )}
      </div>

      {/* Booking details */}
      <div className="bg-white rounded-2xl border border-gray-100 p-6">
        <h2 className="font-bold text-gray-900 mb-4">Детали заказа</h2>
        <div className="space-y-2 text-sm">
          <div className="flex justify-between">
            <span className="text-gray-500">Тип</span>
            <span className="font-medium">{booking.booking_type === 'custom' ? 'Конструктор' : 'Готовый сеап'}</span>
          </div>
          <div className="flex justify-between">
            <span className="text-gray-500">Тариф</span>
            <span className="font-medium">{booking.rental_type === 'hourly' ? 'Почасовой' : 'Посуточный'}</span>
          </div>
          <div className="flex justify-between">
            <span className="text-gray-500">Начало</span>
            <span className="font-medium">{new Date(booking.start_time).toLocaleString('ru-RU')}</span>
          </div>
          <div className="flex justify-between">
            <span className="text-gray-500">Конец</span>
            <span className="font-medium">{new Date(booking.end_time).toLocaleString('ru-RU')}</span>
          </div>
          <div className="flex justify-between border-t pt-2 font-bold">
            <span>Сумма</span>
            <span className="text-primary-600">{booking.total_price.toLocaleString()} ₸</span>
          </div>
        </div>
      </div>
    </div>
  );
}
