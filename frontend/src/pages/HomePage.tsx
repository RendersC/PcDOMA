import { Link } from 'react-router-dom';
import { Monitor, Zap, Shield, MapPin, ChevronRight } from 'lucide-react';
import { useQuery } from '@tanstack/react-query';
import { catalogApi } from '../api/catalog';
import type { PC } from '../types';
import { Hero } from '../components/home/Hero';

function PCCard({ pc }: { pc: PC }) {
  return (
    <Link to={`/catalog/${pc.id}`} className="group bg-white rounded-2xl shadow-sm border border-gray-100 overflow-hidden hover:shadow-md transition-all hover:-translate-y-1">
      <div className="h-40 bg-gradient-to-br from-gray-800 to-gray-900 flex items-center justify-center">
        {pc.images?.[0] ? (
          <img src={pc.images[0]} alt={pc.name} className="w-full h-full object-cover" />
        ) : (
          <Monitor className="w-12 h-12 text-gray-600" />
        )}
      </div>
      <div className="p-4">
        <h3 className="font-semibold text-gray-900 mb-1">{pc.name}</h3>
        <p className="text-xs text-gray-500 mb-2">{pc.specs.cpu} • {pc.specs.gpu}</p>
        <div className="flex items-center justify-between">
          <span className="text-primary-600 font-bold text-sm">от {pc.price_per_hour.toLocaleString()} ₸/ч</span>
          <span className={`text-xs px-2 py-0.5 rounded-full ${pc.status === 'available' ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'}`}>
            {pc.status === 'available' ? 'Свободен' : 'Занят'}
          </span>
        </div>
      </div>
    </Link>
  );
}

const features = [
  { icon: Zap, title: 'Мощные конфигурации', desc: 'RTX 4080, i9-13900K — топовое железо для любых задач' },
  { icon: Shield, title: 'Надёжно и безопасно', desc: 'Все ПК проверены и застрахованы. Работаем официально' },
  { icon: MapPin, title: 'Несколько точек', desc: 'Пункты выдачи во всех районах Астаны' },
];

export function HomePage() {
  const { data } = useQuery({
    queryKey: ['pcs', 'featured'],
    queryFn: () => catalogApi.listPCs({ limit: 6, status: 'available' }),
  });

  const pcs = data?.data?.data || [];

  return (
    <div>
      {/* Hero */}
      <Hero />

      {/* Features */}
      <section className="py-16 bg-gray-50">
        <div className="max-w-5xl mx-auto px-4">
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            {features.map(({ icon: Icon, title, desc }) => (
              <div key={title} className="bg-white rounded-2xl p-6 shadow-sm border border-gray-100">
                <div className="w-10 h-10 bg-primary-50 rounded-xl flex items-center justify-center mb-4">
                  <Icon className="w-5 h-5 text-primary-600" />
                </div>
                <h3 className="font-semibold text-gray-900 mb-2">{title}</h3>
                <p className="text-sm text-gray-500">{desc}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Featured PCs */}
      {pcs.length > 0 && (
        <section className="py-16">
          <div className="max-w-6xl mx-auto px-4">
            <div className="flex items-center justify-between mb-8">
              <h2 className="text-2xl font-bold text-gray-900">Доступные ПК</h2>
              <Link to="/catalog" className="text-primary-600 text-sm font-medium flex items-center gap-1 hover:gap-2 transition-all">
                Все ПК <ChevronRight className="w-4 h-4" />
              </Link>
            </div>
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-5">
              {pcs.map((pc) => <PCCard key={pc.id} pc={pc} />)}
            </div>
          </div>
        </section>
      )}
    </div>
  );
}
