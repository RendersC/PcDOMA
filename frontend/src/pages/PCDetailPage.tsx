import { useParams, Link } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { Monitor, Star, MapPin, Cpu, ChevronRight } from 'lucide-react';
import { catalogApi } from '../api/catalog';
import { useAuthStore } from '../store/authStore';

export function PCDetailPage() {
  const { id } = useParams<{ id: string }>();
  const { isAuthenticated } = useAuthStore();

  const { data: pcData, isLoading } = useQuery({
    queryKey: ['pc', id],
    queryFn: () => catalogApi.getPC(id!),
    enabled: !!id,
  });

  const { data: reviewsData } = useQuery({
    queryKey: ['reviews', id],
    queryFn: () => catalogApi.listReviews(id!),
    enabled: !!id,
  });

  const pc = pcData?.data;
  const reviews = reviewsData?.data?.data || [];

  if (isLoading) {
    return (
      <div className="max-w-5xl mx-auto px-4 py-10">
        <div className="animate-pulse space-y-4">
          <div className="h-64 bg-gray-200 rounded-2xl" />
          <div className="h-8 bg-gray-200 rounded-lg w-1/2" />
          <div className="h-4 bg-gray-100 rounded-lg w-3/4" />
        </div>
      </div>
    );
  }

  if (!pc) {
    return <div className="text-center py-20 text-gray-400">ПК не найден</div>;
  }

  return (
    <div className="max-w-5xl mx-auto px-4 py-8 sm:px-6 lg:px-8">
      {/* Breadcrumb */}
      <nav className="flex items-center gap-2 text-sm text-gray-500 mb-6">
        <Link to="/catalog" className="hover:text-primary-600">Каталог</Link>
        <ChevronRight className="w-4 h-4" />
        <span className="text-gray-900">{pc.name}</span>
      </nav>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        {/* Left: image + specs */}
        <div className="lg:col-span-2 space-y-6">
          <div className="h-72 bg-gradient-to-br from-gray-800 to-gray-900 rounded-2xl flex items-center justify-center overflow-hidden">
            {pc.images?.[0] ? (
              <img src={pc.images[0]} alt={pc.name} className="w-full h-full object-cover rounded-2xl" />
            ) : (
              <Monitor className="w-20 h-20 text-gray-700" />
            )}
          </div>

          <div className="bg-white rounded-2xl border border-gray-100 p-6">
            <h2 className="font-bold text-gray-900 mb-4 text-sm uppercase tracking-wide">Характеристики</h2>
            <div className="grid grid-cols-2 gap-3 text-sm">
              {[
                { label: 'Процессор', value: pc.specs.cpu },
                { label: 'Видеокарта', value: pc.specs.gpu },
                { label: 'Оперативная память', value: pc.specs.ram },
                { label: 'Накопитель', value: pc.specs.storage },
              ].map(({ label, value }) => (
                <div key={label} className="bg-gray-50 rounded-xl p-3">
                  <p className="text-gray-500 text-xs mb-0.5">{label}</p>
                  <p className="font-medium text-gray-900">{value}</p>
                </div>
              ))}
            </div>
          </div>

          {/* Reviews */}
          <div className="bg-white rounded-2xl border border-gray-100 p-6">
            <div className="flex items-center justify-between mb-4">
              <h2 className="font-bold text-gray-900">Отзывы ({reviews.length})</h2>
              <div className="flex items-center gap-1">
                <Star className="w-4 h-4 text-yellow-400 fill-current" />
                <span className="font-semibold text-sm">{pc.rating_avg.toFixed(1)}</span>
              </div>
            </div>
            {reviews.length === 0 ? (
              <p className="text-sm text-gray-400">Отзывов пока нет. Будьте первым!</p>
            ) : (
              <div className="space-y-3">
                {reviews.slice(0, 5).map((r) => (
                  <div key={r.id} className="flex gap-3">
                    <div className="w-8 h-8 bg-primary-100 text-primary-700 rounded-full flex items-center justify-center font-semibold text-xs flex-shrink-0">
                      {r.user_name[0]}
                    </div>
                    <div>
                      <div className="flex items-center gap-2 mb-1">
                        <span className="text-sm font-medium text-gray-900">{r.user_name}</span>
                        <div className="flex">
                          {[...Array(r.rating)].map((_, i) => (
                            <Star key={i} className="w-3 h-3 text-yellow-400 fill-current" />
                          ))}
                        </div>
                      </div>
                      <p className="text-sm text-gray-600">{r.comment}</p>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>

        {/* Right: booking card */}
        <div className="lg:col-span-1">
          <div className="bg-white rounded-2xl border border-gray-100 p-6 sticky top-24 shadow-sm">
            <h1 className="text-xl font-bold text-gray-900 mb-1">{pc.name}</h1>
            <div className="flex items-center gap-1 mb-4">
              <Star className="w-4 h-4 text-yellow-400 fill-current" />
              <span className="text-sm text-gray-600">{pc.rating_avg.toFixed(1)} ({pc.rating_count} отзывов)</span>
            </div>

            <div className="flex items-center gap-2 text-sm text-gray-500 mb-5">
              <MapPin className="w-4 h-4" />
              {pc.location.district}, {pc.location.address}
            </div>

            <div className="border-t border-gray-100 pt-4 mb-5 space-y-2">
              <div className="flex justify-between text-sm">
                <span className="text-gray-500">Почасовая аренда</span>
                <span className="font-bold text-gray-900">{pc.price_per_hour.toLocaleString()} ₸/ч</span>
              </div>
              <div className="flex justify-between text-sm">
                <span className="text-gray-500">Посуточная аренда</span>
                <span className="font-bold text-gray-900">{pc.price_per_day.toLocaleString()} ₸/сут</span>
              </div>
            </div>

            <div className="mb-4">
              <span className={`inline-flex items-center gap-1.5 text-sm px-3 py-1.5 rounded-full font-medium ${
                pc.status === 'available'
                  ? 'bg-green-50 text-green-700'
                  : 'bg-red-50 text-red-700'
              }`}>
                <span className={`w-1.5 h-1.5 rounded-full ${pc.status === 'available' ? 'bg-green-500' : 'bg-red-500'}`} />
                {pc.status === 'available' ? 'Доступен для аренды' : 'Недоступен'}
              </span>
            </div>

            {pc.status === 'available' ? (
              isAuthenticated ? (
                <Link
                  to={`/booking/${pc.id}`}
                  className="w-full bg-primary-600 hover:bg-primary-700 text-white py-3 rounded-xl font-semibold text-sm transition-colors block text-center"
                >
                  Забронировать
                </Link>
              ) : (
                <Link
                  to="/login"
                  className="w-full bg-gray-900 hover:bg-gray-800 text-white py-3 rounded-xl font-semibold text-sm transition-colors block text-center"
                >
                  Войдите для бронирования
                </Link>
              )
            ) : (
              <button disabled className="w-full bg-gray-100 text-gray-400 py-3 rounded-xl font-semibold text-sm cursor-not-allowed">
                Недоступен
              </button>
            )}

            <Link
              to={`/configurator?pc_id=${pc.id}`}
              className="mt-3 w-full border border-primary-200 text-primary-600 hover:bg-primary-50 py-2.5 rounded-xl font-medium text-sm transition-colors block text-center"
            >
              <Cpu className="w-4 h-4 inline mr-1" />
              Настроить периферию
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
}
