import { useState } from 'react';
import { Link } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { Zap, Briefcase, Palette, Radio, Star, ChevronRight } from 'lucide-react';
import { catalogApi } from '../api/catalog';
import type { Setup } from '../types';

const CATEGORIES = [
  { value: '', label: 'Все сеапы' },
  { value: 'gaming', label: 'Игровые', icon: <Zap className="w-4 h-4" /> },
  { value: 'work', label: 'Рабочие', icon: <Briefcase className="w-4 h-4" /> },
  { value: 'design', label: 'Дизайн', icon: <Palette className="w-4 h-4" /> },
  { value: 'streaming', label: 'Стриминг', icon: <Radio className="w-4 h-4" /> },
];

const CATEGORY_COLORS: Record<string, string> = {
  gaming: 'bg-purple-100 text-purple-700',
  work: 'bg-blue-100 text-blue-700',
  design: 'bg-pink-100 text-pink-700',
  streaming: 'bg-red-100 text-red-700',
};

function SetupCard({ setup }: { setup: Setup }) {
  const discounted = setup.discount_percent > 0;
  const originalPerHour = discounted
    ? Math.round(setup.total_price_per_hour / (1 - setup.discount_percent / 100))
    : setup.total_price_per_hour;

  return (
    <Link
      to={`/setups/${setup.id}`}
      className="bg-white rounded-2xl border border-gray-100 overflow-hidden hover:shadow-md transition-shadow group"
    >
      <div className="h-40 bg-gradient-to-br from-gray-800 to-gray-900 flex items-center justify-center relative">
        <span className="text-5xl">🖥</span>
        {discounted && (
          <span className="absolute top-3 right-3 bg-green-500 text-white text-xs font-bold px-2 py-1 rounded-full">
            -{setup.discount_percent}%
          </span>
        )}
      </div>
      <div className="p-5">
        <div className="flex items-start justify-between gap-2 mb-2">
          <h3 className="font-bold text-gray-900 group-hover:text-primary-600 transition-colors text-sm leading-tight">
            {setup.name}
          </h3>
          <span className={`text-xs px-2 py-0.5 rounded-full font-medium flex-shrink-0 ${CATEGORY_COLORS[setup.category] || 'bg-gray-100 text-gray-600'}`}>
            {CATEGORIES.find((c) => c.value === setup.category)?.label || setup.category}
          </span>
        </div>

        <p className="text-xs text-gray-500 mb-3 line-clamp-2">{setup.description}</p>

        <div className="flex items-center gap-1 mb-3">
          <Star className="w-3.5 h-3.5 text-yellow-400 fill-yellow-400" />
          <span className="text-xs font-medium text-gray-700">{setup.rating_avg.toFixed(1)}</span>
          <span className="text-xs text-gray-400">({setup.rating_count})</span>
        </div>

        <div className="flex items-center justify-between">
          <div>
            {discounted && (
              <p className="text-xs text-gray-400 line-through">{originalPerHour.toLocaleString()} ₸/ч</p>
            )}
            <p className="font-bold text-primary-600">{setup.total_price_per_hour.toLocaleString()} ₸/ч</p>
          </div>
          <ChevronRight className="w-4 h-4 text-gray-400 group-hover:text-primary-600 transition-colors" />
        </div>
      </div>
    </Link>
  );
}

export function SetupsPage() {
  const [category, setCategory] = useState('');

  const { data, isLoading } = useQuery({
    queryKey: ['setups', category],
    queryFn: () => catalogApi.listSetups(category ? { category } : undefined),
  });

  const setups = data?.data?.data || [];

  return (
    <div className="max-w-7xl mx-auto px-4 py-8 sm:px-6 lg:px-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900 mb-2">Готовые сеапы</h1>
        <p className="text-gray-500">Преднастроенные конфигурации со скидкой</p>
      </div>

      {/* Category filter */}
      <div className="flex flex-wrap gap-2 mb-8">
        {CATEGORIES.map((cat) => (
          <button
            key={cat.value}
            onClick={() => setCategory(cat.value)}
            className={`flex items-center gap-1.5 px-4 py-2 rounded-full text-sm font-medium transition-colors ${
              category === cat.value
                ? 'bg-primary-600 text-white'
                : 'bg-white border border-gray-200 text-gray-600 hover:border-primary-300'
            }`}
          >
            {cat.icon}
            {cat.label}
          </button>
        ))}
      </div>

      {isLoading ? (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
          {[...Array(6)].map((_, i) => (
            <div key={i} className="bg-white rounded-2xl border border-gray-100 overflow-hidden animate-pulse">
              <div className="h-40 bg-gray-200" />
              <div className="p-5 space-y-3">
                <div className="h-4 bg-gray-200 rounded w-3/4" />
                <div className="h-3 bg-gray-200 rounded w-full" />
                <div className="h-3 bg-gray-200 rounded w-1/2" />
              </div>
            </div>
          ))}
        </div>
      ) : setups.length === 0 ? (
        <div className="text-center py-16">
          <p className="text-gray-400 text-lg mb-2">Сеапы не найдены</p>
          <p className="text-gray-400 text-sm">Попробуйте другую категорию</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
          {setups.map((setup) => (
            <SetupCard key={setup.id} setup={setup} />
          ))}
        </div>
      )}
    </div>
  );
}
