import { useState } from 'react';
import { Link } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { Monitor, Search, SlidersHorizontal, Star } from 'lucide-react';
import { catalogApi } from '../api/catalog';
import type { PC } from '../types';

function PCCard({ pc }: { pc: PC }) {
  return (
    <Link to={`/catalog/${pc.id}`} className="group bg-white rounded-2xl border border-gray-100 shadow-sm hover:shadow-md transition-all hover:-translate-y-0.5 overflow-hidden">
      <div className="h-44 bg-gradient-to-br from-gray-800 to-gray-900 flex items-center justify-center relative">
        {pc.images?.[0] ? (
          <img src={pc.images[0]} alt={pc.name} className="w-full h-full object-cover" />
        ) : (
          <Monitor className="w-16 h-16 text-gray-700" />
        )}
        <div className="absolute top-3 right-3">
          <span className={`text-xs px-2.5 py-1 rounded-full font-medium ${
            pc.status === 'available' ? 'bg-green-500 text-white' : 'bg-red-500 text-white'
          }`}>
            {pc.status === 'available' ? 'Свободен' : pc.status === 'booked' ? 'Занят' : 'Тех. работы'}
          </span>
        </div>
      </div>
      <div className="p-5">
        <h3 className="font-bold text-gray-900 mb-1">{pc.name}</h3>
        <div className="text-xs text-gray-500 space-y-0.5 mb-3">
          <p>{pc.specs.cpu}</p>
          <p>{pc.specs.gpu} • {pc.specs.ram}</p>
        </div>
        <div className="flex items-center gap-1 mb-3">
          <Star className="w-3.5 h-3.5 text-yellow-400 fill-current" />
          <span className="text-xs text-gray-600">{pc.rating_avg.toFixed(1)} ({pc.rating_count})</span>
        </div>
        <div className="flex items-center justify-between pt-3 border-t border-gray-50">
          <div>
            <span className="text-primary-600 font-bold">{pc.price_per_hour.toLocaleString()} ₸</span>
            <span className="text-gray-400 text-xs">/час</span>
          </div>
          <span className="text-xs text-gray-400">{pc.location.district}</span>
        </div>
      </div>
    </Link>
  );
}

const STATUSES = [
  { value: '', label: 'Все' },
  { value: 'available', label: 'Свободные' },
  { value: 'booked', label: 'Занятые' },
];

const SORTS = [
  { value: '', label: 'По умолчанию' },
  { value: 'price_asc', label: 'Цена ↑' },
  { value: 'price_desc', label: 'Цена ↓' },
  { value: 'rating', label: 'По рейтингу' },
];

export function CatalogPage() {
  const [search, setSearch] = useState('');
  const [status, setStatus] = useState('');
  const [sort, setSort] = useState('');
  const [showFilters, setShowFilters] = useState(false);

  const { data, isLoading } = useQuery({
    queryKey: ['pcs', { status, sort }],
    queryFn: () => catalogApi.listPCs({ status, sort, limit: 50 }),
  });

  const pcs = (data?.data?.data || []).filter((pc) =>
    search === '' || pc.name.toLowerCase().includes(search.toLowerCase())
  );

  return (
    <div className="max-w-7xl mx-auto px-4 py-8 sm:px-6 lg:px-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900 mb-2">Каталог ПК</h1>
        <p className="text-gray-500">Выберите подходящую конфигурацию для аренды в Астане</p>
      </div>

      {/* Search & filters */}
      <div className="flex flex-col sm:flex-row gap-3 mb-6">
        <div className="relative flex-1">
          <Search className="absolute left-3 top-3 w-4 h-4 text-gray-400" />
          <input
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="Поиск по названию..."
            className="w-full pl-10 pr-4 py-2.5 border border-gray-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-primary-500"
          />
        </div>
        <button
          onClick={() => setShowFilters(!showFilters)}
          className="flex items-center gap-2 px-4 py-2.5 border border-gray-200 rounded-xl text-sm text-gray-600 hover:border-primary-300"
        >
          <SlidersHorizontal className="w-4 h-4" />
          Фильтры
        </button>
      </div>

      {showFilters && (
        <div className="bg-gray-50 rounded-xl p-4 mb-6 flex flex-wrap gap-4">
          <div>
            <p className="text-xs font-medium text-gray-700 mb-2">Статус</p>
            <div className="flex gap-2">
              {STATUSES.map((s) => (
                <button
                  key={s.value}
                  onClick={() => setStatus(s.value)}
                  className={`text-xs px-3 py-1.5 rounded-lg transition-colors ${
                    status === s.value
                      ? 'bg-primary-600 text-white'
                      : 'bg-white text-gray-600 border border-gray-200 hover:border-primary-300'
                  }`}
                >
                  {s.label}
                </button>
              ))}
            </div>
          </div>
          <div>
            <p className="text-xs font-medium text-gray-700 mb-2">Сортировка</p>
            <div className="flex gap-2 flex-wrap">
              {SORTS.map((s) => (
                <button
                  key={s.value}
                  onClick={() => setSort(s.value)}
                  className={`text-xs px-3 py-1.5 rounded-lg transition-colors ${
                    sort === s.value
                      ? 'bg-primary-600 text-white'
                      : 'bg-white text-gray-600 border border-gray-200 hover:border-primary-300'
                  }`}
                >
                  {s.label}
                </button>
              ))}
            </div>
          </div>
        </div>
      )}

      {isLoading ? (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-5">
          {[...Array(6)].map((_, i) => (
            <div key={i} className="bg-gray-100 rounded-2xl h-64 animate-pulse" />
          ))}
        </div>
      ) : pcs.length === 0 ? (
        <div className="text-center py-20 text-gray-400">
          <Monitor className="w-12 h-12 mx-auto mb-3 opacity-30" />
          <p>ПК не найдены</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-5">
          {pcs.map((pc) => <PCCard key={pc.id} pc={pc} />)}
        </div>
      )}
    </div>
  );
}
