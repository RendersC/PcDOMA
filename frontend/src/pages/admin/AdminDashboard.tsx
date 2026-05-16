import { useQuery } from '@tanstack/react-query';
import { Link } from 'react-router-dom';
import { Monitor, Mouse, Package, ShoppingBag, ChevronRight } from 'lucide-react';
import { catalogApi } from '../../api/catalog';
import { bookingsApi } from '../../api/bookings';

export function AdminDashboard() {
  const { data: pcsData } = useQuery({ queryKey: ['admin-pcs'], queryFn: () => catalogApi.listPCs({ limit: 100 }) });
  const { data: peripheralsData } = useQuery({ queryKey: ['admin-peripherals'], queryFn: () => catalogApi.listPeripherals({ limit: '100' }) });
  const { data: setupsData } = useQuery({ queryKey: ['admin-setups'], queryFn: () => catalogApi.listSetups() });
  const { data: bookingsData } = useQuery({ queryKey: ['admin-bookings-all'], queryFn: () => bookingsApi.adminList({ limit: '100' }) });

  const pcs = pcsData?.data?.data || [];
  const peripherals = peripheralsData?.data?.data || [];
  const setups = setupsData?.data?.data || [];
  const allBookings = bookingsData?.data?.data || [];
  const activeBookings = allBookings.filter((b) => ['pending_worker', 'accepted', 'delivering', 'active'].includes(b.status));
  const totalRevenue = allBookings.filter((b) => b.status !== 'cancelled' && b.status !== 'rejected').reduce((s, b) => s + b.total_price, 0);

  const stats = [
    { label: 'Системных блоков', value: pcs.length, icon: <Monitor className="w-6 h-6" />, color: 'bg-blue-50 text-blue-600', to: '/admin/pcs' },
    { label: 'Периферии', value: peripherals.length, icon: <Mouse className="w-6 h-6" />, color: 'bg-purple-50 text-purple-600', to: '/admin/peripherals' },
    { label: 'Готовых сеапов', value: setups.length, icon: <Package className="w-6 h-6" />, color: 'bg-pink-50 text-pink-600', to: '/admin/setups' },
    { label: 'Активных аренд', value: activeBookings.length, icon: <ShoppingBag className="w-6 h-6" />, color: 'bg-green-50 text-green-600', to: '/admin/bookings' },
  ];

  const quickLinks = [
    { label: 'Управление ПК', desc: 'Добавить, изменить, удалить', to: '/admin/pcs' },
    { label: 'Управление периферией', desc: 'Мышки, наушники, мониторы', to: '/admin/peripherals' },
    { label: 'Управление сеапами', desc: 'Готовые конфигурации', to: '/admin/setups' },
    { label: 'Все бронирования', desc: 'История и статистика', to: '/admin/bookings' },
  ];

  return (
    <div className="max-w-6xl mx-auto px-4 py-8 sm:px-6">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900">Панель администратора</h1>
        <p className="text-gray-500 mt-1">PCDoma — управление сервисом</p>
      </div>

      {/* Stats */}
      <div className="grid grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
        {stats.map((s) => (
          <Link key={s.to} to={s.to} className="bg-white rounded-2xl border border-gray-100 p-5 hover:shadow-md transition-shadow">
            <div className={`w-12 h-12 rounded-xl flex items-center justify-center mb-3 ${s.color}`}>
              {s.icon}
            </div>
            <p className="text-2xl font-bold text-gray-900">{s.value}</p>
            <p className="text-sm text-gray-500 mt-0.5">{s.label}</p>
          </Link>
        ))}
      </div>

      {/* Revenue */}
      <div className="bg-gradient-to-r from-primary-600 to-primary-700 rounded-2xl p-6 mb-8 text-white">
        <p className="text-primary-100 text-sm mb-1">Общая выручка</p>
        <p className="text-4xl font-bold">{totalRevenue.toLocaleString()} ₸</p>
        <p className="text-primary-200 text-sm mt-1">{allBookings.length} бронирований всего</p>
      </div>

      {/* Quick links */}
      <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
        {quickLinks.map((l) => (
          <Link
            key={l.to}
            to={l.to}
            className="bg-white rounded-2xl border border-gray-100 p-5 flex items-center justify-between hover:border-primary-200 hover:shadow-sm transition-all group"
          >
            <div>
              <p className="font-semibold text-gray-900 group-hover:text-primary-600 transition-colors">{l.label}</p>
              <p className="text-sm text-gray-400 mt-0.5">{l.desc}</p>
            </div>
            <ChevronRight className="w-5 h-5 text-gray-300 group-hover:text-primary-400 transition-colors" />
          </Link>
        ))}
      </div>
    </div>
  );
}
