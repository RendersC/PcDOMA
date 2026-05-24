import { useParams, useNavigate } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { Star, Zap, Briefcase, Palette, Radio, ArrowLeft, Monitor, Mouse, Headphones, Tv } from 'lucide-react';
import { catalogApi } from '../api/catalog';
import { useAuthStore } from '../store/authStore';

const CATEGORY_LABELS: Record<string, string> = {
  gaming: 'Игровой',
  work: 'Рабочий',
  design: 'Дизайн',
  streaming: 'Стриминг',
};

const CATEGORY_ICONS: Record<string, React.ReactNode> = {
  gaming: <Zap className="w-4 h-4" />,
  work: <Briefcase className="w-4 h-4" />,
  design: <Palette className="w-4 h-4" />,
  streaming: <Radio className="w-4 h-4" />,
};

const PERIPHERAL_ICONS: Record<string, React.ReactNode> = {
  mouse: <Mouse className="w-4 h-4" />,
  headphones: <Headphones className="w-4 h-4" />,
  headset: <Headphones className="w-4 h-4" />,
  monitor: <Tv className="w-4 h-4" />,
  keyboard: <Monitor className="w-4 h-4" />,
};

export function SetupDetailPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { isAuthenticated } = useAuthStore();

  const { data: setupRes, isLoading } = useQuery({
    queryKey: ['setup', id],
    queryFn: () => catalogApi.getSetup(id!),
    enabled: !!id,
  });

  const setup = setupRes?.data;

  const { data: pcRes } = useQuery({
    queryKey: ['pc', setup?.pc_id],
    queryFn: () => catalogApi.getPC(setup!.pc_id),
    enabled: !!setup?.pc_id,
  });

  const { data: peripheralsRes } = useQuery({
    queryKey: ['peripherals', 'all'],
    queryFn: () => catalogApi.listPeripherals({ limit: '100' }),
  });

  const pc = pcRes?.data;
  const allPeripherals = peripheralsRes?.data?.data || [];
  const setupPeripherals = allPeripherals.filter((p) => setup?.peripheral_ids?.includes(p.id));

  const handleBook = () => {
    if (!isAuthenticated) {
      navigate('/login', { state: { from: { pathname: `/setups/${id}` } } });
      return;
    }
    const params = new URLSearchParams({
      type: 'setup',
      setup_id: setup!.id,
    });
    navigate(`/booking/setup?${params}`);
  };

  if (isLoading) {
    return (
      <div className="max-w-5xl mx-auto px-4 py-8">
        <div className="animate-pulse space-y-6">
          <div className="h-8 bg-gray-200 rounded w-1/2" />
          <div className="h-64 bg-gray-200 rounded-2xl" />
          <div className="h-40 bg-gray-200 rounded-2xl" />
        </div>
      </div>
    );
  }

  if (!setup) {
    return (
      <div className="max-w-5xl mx-auto px-4 py-16 text-center">
        <p className="text-gray-500 text-lg">Сеап не найден</p>
      </div>
    );
  }

  const originalPerHour = setup.discount_percent > 0
    ? Math.round(setup.total_price_per_hour / (1 - setup.discount_percent / 100))
    : null;

  return (
    <div className="max-w-5xl mx-auto px-4 py-8 sm:px-6">
      <button onClick={() => navigate(-1)} className="flex items-center gap-2 text-gray-500 hover:text-gray-700 mb-6 text-sm">
        <ArrowLeft className="w-4 h-4" /> Назад
      </button>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        {/* Main content */}
        <div className="lg:col-span-2 space-y-6">
          {/* Cover image */}
          {setup.images?.[0] && (
            <div className="h-64 w-full overflow-hidden rounded-2xl border border-gray-100 bg-gray-100">
              <img src={setup.images[0]} alt={setup.name} className="h-full w-full object-cover" />
            </div>
          )}

          {/* Header */}
          <div className="bg-white rounded-2xl border border-gray-100 p-6">
            <div className="flex items-start justify-between gap-4 mb-4">
              <div>
                <div className="flex items-center gap-2 mb-2">
                  <span className="flex items-center gap-1 text-xs font-medium px-2 py-1 bg-primary-50 text-primary-700 rounded-full">
                    {CATEGORY_ICONS[setup.category]}
                    {CATEGORY_LABELS[setup.category] || setup.category}
                  </span>
                  {setup.discount_percent > 0 && (
                    <span className="text-xs font-bold px-2 py-1 bg-green-100 text-green-700 rounded-full">
                      Скидка {setup.discount_percent}%
                    </span>
                  )}
                </div>
                <h1 className="text-2xl font-bold text-gray-900">{setup.name}</h1>
              </div>
              <div className="flex items-center gap-1 flex-shrink-0">
                <Star className="w-4 h-4 text-yellow-400 fill-yellow-400" />
                <span className="font-medium text-gray-800">{setup.rating_avg.toFixed(1)}</span>
                <span className="text-gray-400 text-sm">({setup.rating_count})</span>
              </div>
            </div>
            <p className="text-gray-600 text-sm leading-relaxed">{setup.description}</p>
          </div>

          {/* PC Block */}
          {pc && (
            <div className="bg-white rounded-2xl border border-gray-100 p-6">
              <h2 className="font-bold text-gray-900 mb-4 flex items-center gap-2">
                <Monitor className="w-5 h-5 text-primary-600" /> Системный блок
              </h2>
              <div className="flex items-start justify-between">
                <div>
                  <p className="font-semibold text-gray-900">{pc.name}</p>
                  <div className="mt-2 grid grid-cols-2 gap-x-8 gap-y-1">
                    {Object.entries(pc.specs).map(([key, val]) => (
                      <div key={key} className="flex gap-2 text-sm">
                        <span className="text-gray-400 capitalize">{key}:</span>
                        <span className="text-gray-700 font-medium">{val}</span>
                      </div>
                    ))}
                  </div>
                </div>
                <span className="text-sm text-gray-500">{pc.price_per_hour.toLocaleString()} ₸/ч</span>
              </div>
            </div>
          )}

          {/* Peripherals */}
          {setupPeripherals.length > 0 && (
            <div className="bg-white rounded-2xl border border-gray-100 p-6">
              <h2 className="font-bold text-gray-900 mb-4">Периферия в сеапе</h2>
              <div className="space-y-3">
                {setupPeripherals.map((p) => (
                  <div key={p.id} className="flex items-center justify-between py-2 border-b border-gray-50 last:border-0">
                    <div className="flex items-center gap-3">
                      <span className="text-gray-500">{PERIPHERAL_ICONS[p.type] || <Monitor className="w-4 h-4" />}</span>
                      <div>
                        <p className="text-sm font-medium text-gray-900">{p.name}</p>
                        <p className="text-xs text-gray-400 capitalize">{p.type}</p>
                      </div>
                    </div>
                    <span className="text-sm text-gray-500">+{p.price_per_hour.toLocaleString()} ₸/ч</span>
                  </div>
                ))}
              </div>
              {setup.peripheral_ids.length > setupPeripherals.length && (
                <p className="text-xs text-gray-400 mt-2">
                  + ещё {setup.peripheral_ids.length - setupPeripherals.length} позиций
                </p>
              )}
            </div>
          )}
        </div>

        {/* Sidebar */}
        <div>
          <div className="bg-white rounded-2xl border border-gray-100 p-6 sticky top-24 shadow-sm">
            <h2 className="font-bold text-gray-900 mb-4">Стоимость</h2>

            <div className="space-y-2 mb-5">
              {originalPerHour && (
                <div className="flex justify-between text-sm text-gray-400">
                  <span>Без скидки</span>
                  <span className="line-through">{originalPerHour.toLocaleString()} ₸/ч</span>
                </div>
              )}
              <div className="flex justify-between font-bold text-lg">
                <span className="text-gray-900">За час</span>
                <span className="text-primary-600">{setup.total_price_per_hour.toLocaleString()} ₸</span>
              </div>
              <div className="flex justify-between text-sm text-gray-500">
                <span>За сутки</span>
                <span>{setup.total_price_per_day.toLocaleString()} ₸</span>
              </div>
              {setup.discount_percent > 0 && (
                <div className="bg-green-50 rounded-xl p-3 text-center">
                  <p className="text-green-700 text-sm font-medium">
                    Экономия {setup.discount_percent}% при выборе сеапа
                  </p>
                </div>
              )}
            </div>

            <button
              onClick={handleBook}
              className="w-full bg-primary-600 hover:bg-primary-700 text-white py-3 rounded-xl font-semibold text-sm transition-colors"
            >
              Забронировать сеап
            </button>

            <p className="text-xs text-gray-400 text-center mt-3">
              Включает ПК + {setup.peripheral_ids.length} периферийных устройств
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}
