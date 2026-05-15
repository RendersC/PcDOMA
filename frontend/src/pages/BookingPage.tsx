import { useState, useEffect } from 'react';
import { useParams, useSearchParams, useNavigate } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { Calendar, Clock, ChevronRight, ArrowLeft } from 'lucide-react';
import { catalogApi } from '../api/catalog';
import { bookingsApi } from '../api/bookings';
import toast from 'react-hot-toast';

function addHours(date: Date, hours: number) {
  return new Date(date.getTime() + hours * 60 * 60 * 1000);
}

function formatDateTimeLocal(date: Date) {
  const pad = (n: number) => String(n).padStart(2, '0');
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())}T${pad(date.getHours())}:${pad(date.getMinutes())}`;
}

function diffHours(start: Date, end: Date) {
  return Math.max(0, (end.getTime() - start.getTime()) / (1000 * 60 * 60));
}

export function BookingPage() {
  const { id } = useParams<{ id: string }>();
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();

  const bookingType = (searchParams.get('type') || 'custom') as 'custom' | 'setup';
  const pcId = searchParams.get('pc_id') || id || '';
  const setupId = searchParams.get('setup_id') || '';
  const peripheralIdsParam = searchParams.get('peripheral_ids') || '';
  const peripheralIds = peripheralIdsParam ? peripheralIdsParam.split(',').filter(Boolean) : [];

  const now = new Date();
  const defaultStart = addHours(now, 1);
  const defaultEnd = addHours(now, 3);

  const [rentalType, setRentalType] = useState<'hourly' | 'daily'>('hourly');
  const [startTime, setStartTime] = useState(formatDateTimeLocal(defaultStart));
  const [endTime, setEndTime] = useState(formatDateTimeLocal(defaultEnd));
  const [loading, setLoading] = useState(false);

  const { data: pcRes } = useQuery({
    queryKey: ['pc', pcId],
    queryFn: () => catalogApi.getPC(pcId),
    enabled: bookingType === 'custom' && !!pcId,
  });

  const { data: setupRes } = useQuery({
    queryKey: ['setup', setupId],
    queryFn: () => catalogApi.getSetup(setupId),
    enabled: bookingType === 'setup' && !!setupId,
  });

  const { data: peripheralsRes } = useQuery({
    queryKey: ['peripherals', 'all'],
    queryFn: () => catalogApi.listPeripherals({ limit: '100' }),
    enabled: bookingType === 'custom' && peripheralIds.length > 0,
  });

  const pc = pcRes?.data;
  const setup = setupRes?.data;
  const allPeripherals = peripheralsRes?.data?.data || [];
  const selectedPeripherals = allPeripherals.filter((p) => peripheralIds.includes(p.id));

  const start = new Date(startTime);
  const end = new Date(endTime);
  const hours = diffHours(start, end);
  const days = hours / 24;

  let totalPrice = 0;
  let priceBreakdown: { label: string; amount: number }[] = [];

  if (bookingType === 'custom' && pc) {
    if (rentalType === 'hourly') {
      const pcCost = pc.price_per_hour * hours;
      priceBreakdown.push({ label: pc.name, amount: pcCost });
      selectedPeripherals.forEach((p) => {
        priceBreakdown.push({ label: p.name, amount: p.price_per_hour * hours });
      });
    } else {
      const pcCost = pc.price_per_day * Math.ceil(days);
      priceBreakdown.push({ label: pc.name, amount: pcCost });
      selectedPeripherals.forEach((p) => {
        priceBreakdown.push({ label: p.name, amount: p.price_per_day * Math.ceil(days) });
      });
    }
    totalPrice = priceBreakdown.reduce((s, x) => s + x.amount, 0);
  } else if (bookingType === 'setup' && setup) {
    if (rentalType === 'hourly') {
      totalPrice = setup.total_price_per_hour * hours;
    } else {
      totalPrice = setup.total_price_per_day * Math.ceil(days);
    }
    priceBreakdown = [{ label: setup.name, amount: totalPrice }];
  }

  const durationLabel = rentalType === 'hourly'
    ? `${hours.toFixed(1)} ч`
    : `${Math.ceil(days)} сут`;

  const handleSubmit = async () => {
    if (hours <= 0) {
      toast.error('Время окончания должно быть позже начала');
      return;
    }
    if (hours < 1) {
      toast.error('Минимальная аренда — 1 час');
      return;
    }

    setLoading(true);
    try {
      const payload: Parameters<typeof bookingsApi.create>[0] = {
        booking_type: bookingType,
        start_time: new Date(startTime).toISOString(),
        end_time: new Date(endTime).toISOString(),
        rental_type: rentalType,
      };
      if (bookingType === 'custom') {
        payload.pc_id = pcId;
        if (peripheralIds.length > 0) payload.peripheral_ids = peripheralIds;
      } else {
        payload.setup_id = setupId;
      }

      const res = await bookingsApi.create(payload);
      const booking = res.data;
      toast.success('Бронирование создано!');
      navigate(`/payment/${booking.id}`);
    } catch (err: unknown) {
      const msg = (err as { response?: { data?: { error?: string } } })?.response?.data?.error;
      toast.error(msg || 'Ошибка при создании бронирования');
    } finally {
      setLoading(false);
    }
  };

  const isReady = bookingType === 'custom' ? !!pc : !!setup;

  return (
    <div className="max-w-4xl mx-auto px-4 py-8 sm:px-6">
      <button onClick={() => navigate(-1)} className="flex items-center gap-2 text-gray-500 hover:text-gray-700 mb-6 text-sm">
        <ArrowLeft className="w-4 h-4" /> Назад
      </button>

      <h1 className="text-2xl font-bold text-gray-900 mb-6">Оформление аренды</h1>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <div className="lg:col-span-2 space-y-6">
          {/* What's being booked */}
          <div className="bg-white rounded-2xl border border-gray-100 p-6">
            <h2 className="font-bold text-gray-900 mb-4">Что арендуем</h2>
            {!isReady ? (
              <div className="animate-pulse space-y-2">
                <div className="h-4 bg-gray-200 rounded w-1/2" />
                <div className="h-3 bg-gray-200 rounded w-1/3" />
              </div>
            ) : bookingType === 'custom' && pc ? (
              <div className="space-y-2">
                <div className="flex justify-between text-sm">
                  <span className="font-medium text-gray-900">🖥 {pc.name}</span>
                  <span className="text-gray-500">{pc.price_per_hour.toLocaleString()} ₸/ч</span>
                </div>
                <p className="text-xs text-gray-400">{pc.specs.cpu} · {pc.specs.gpu}</p>
                {selectedPeripherals.map((p) => (
                  <div key={p.id} className="flex justify-between text-sm text-gray-600">
                    <span>+ {p.name}</span>
                    <span className="text-gray-400">+{p.price_per_hour.toLocaleString()} ₸/ч</span>
                  </div>
                ))}
              </div>
            ) : setup ? (
              <div className="space-y-1">
                <p className="font-medium text-gray-900">{setup.name}</p>
                <p className="text-xs text-gray-400">
                  {setup.peripheral_ids.length} периферийных устройств · скидка {setup.discount_percent}%
                </p>
                <p className="text-sm text-primary-600 font-semibold">{setup.total_price_per_hour.toLocaleString()} ₸/ч</p>
              </div>
            ) : null}
          </div>

          {/* Rental type */}
          <div className="bg-white rounded-2xl border border-gray-100 p-6">
            <h2 className="font-bold text-gray-900 mb-4">Тариф</h2>
            <div className="grid grid-cols-2 gap-3">
              {(['hourly', 'daily'] as const).map((rt) => (
                <button
                  key={rt}
                  onClick={() => setRentalType(rt)}
                  className={`p-4 rounded-xl border-2 text-left transition-all ${
                    rentalType === rt ? 'border-primary-500 bg-primary-50' : 'border-gray-100 hover:border-gray-200'
                  }`}
                >
                  <p className="font-semibold text-sm text-gray-900">
                    {rt === 'hourly' ? 'Почасовой' : 'Посуточный'}
                  </p>
                  <p className="text-xs text-gray-400 mt-1">
                    {rt === 'hourly' ? 'Минимум 1 час' : 'Минимум 1 сутки'}
                  </p>
                </button>
              ))}
            </div>
          </div>

          {/* Date & time */}
          <div className="bg-white rounded-2xl border border-gray-100 p-6">
            <h2 className="font-bold text-gray-900 mb-4 flex items-center gap-2">
              <Calendar className="w-5 h-5 text-primary-600" /> Время аренды
            </h2>
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1.5">Начало</label>
                <input
                  type="datetime-local"
                  value={startTime}
                  onChange={(e) => setStartTime(e.target.value)}
                  className="w-full px-3 py-2.5 border border-gray-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-primary-500"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1.5">Конец</label>
                <input
                  type="datetime-local"
                  value={endTime}
                  onChange={(e) => setEndTime(e.target.value)}
                  className="w-full px-3 py-2.5 border border-gray-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-primary-500"
                />
              </div>
            </div>
            {hours > 0 && (
              <div className="mt-3 flex items-center gap-2 text-sm text-gray-500">
                <Clock className="w-4 h-4" />
                <span>Продолжительность: <strong className="text-gray-900">{durationLabel}</strong></span>
              </div>
            )}
          </div>
        </div>

        {/* Summary sidebar */}
        <div>
          <div className="bg-white rounded-2xl border border-gray-100 p-6 sticky top-24 shadow-sm">
            <h2 className="font-bold text-gray-900 mb-4">Итого</h2>

            {priceBreakdown.length > 0 && hours > 0 ? (
              <>
                <div className="space-y-2 mb-4">
                  {priceBreakdown.map((item, i) => (
                    <div key={i} className="flex justify-between text-sm">
                      <span className="text-gray-600 truncate mr-2">{item.label}</span>
                      <span className="text-gray-900 font-medium flex-shrink-0">{item.amount.toLocaleString()} ₸</span>
                    </div>
                  ))}
                </div>
                <div className="border-t pt-3 mb-5">
                  <div className="flex justify-between font-bold text-lg">
                    <span className="text-gray-900">Итого</span>
                    <span className="text-primary-600">{Math.round(totalPrice).toLocaleString()} ₸</span>
                  </div>
                  <p className="text-xs text-gray-400 mt-1">За {durationLabel}</p>
                </div>
              </>
            ) : (
              <p className="text-sm text-gray-400 mb-5">Выберите время для расчёта</p>
            )}

            <button
              onClick={handleSubmit}
              disabled={loading || !isReady || hours <= 0}
              className="w-full bg-primary-600 hover:bg-primary-700 disabled:opacity-40 disabled:cursor-not-allowed text-white py-3 rounded-xl font-semibold text-sm transition-colors flex items-center justify-center gap-2"
            >
              {loading ? 'Создаём бронирование...' : 'Перейти к оплате'}
              {!loading && <ChevronRight className="w-4 h-4" />}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
