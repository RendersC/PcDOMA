import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { Link } from 'react-router-dom';
import { Plus, Trash2 } from 'lucide-react';
import { catalogApi } from '../../api/catalog';
import toast from 'react-hot-toast';

const TYPE_LABELS: Record<string, string> = {
  mouse: 'Мышь',
  headphones: 'Наушники',
  headset: 'Гарнитура',
  monitor: 'Монитор',
  keyboard: 'Клавиатура',
};

export function AdminPeripheralList() {
  const queryClient = useQueryClient();

  const { data, isLoading } = useQuery({
    queryKey: ['admin-peripherals'],
    queryFn: () => catalogApi.listPeripherals({ limit: '100' }),
  });

  const deleteMutation = useMutation({
    mutationFn: (id: string) => catalogApi.deletePeripheral(id),
    onSuccess: () => {
      toast.success('Периферия удалена');
      queryClient.invalidateQueries({ queryKey: ['admin-peripherals'] });
    },
    onError: () => toast.error('Не удалось удалить'),
  });

  const peripherals = data?.data?.data || [];

  return (
    <div className="max-w-6xl mx-auto px-4 py-8 sm:px-6">
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold text-gray-900">Периферия</h1>
        <Link
          to="/admin/peripherals/new"
          className="flex items-center gap-2 bg-primary-600 hover:bg-primary-700 text-white px-4 py-2.5 rounded-xl text-sm font-medium transition-colors"
        >
          <Plus className="w-4 h-4" /> Добавить периферию
        </Link>
      </div>

      {isLoading ? (
        <div className="space-y-3">
          {[...Array(4)].map((_, i) => <div key={i} className="bg-white rounded-xl border border-gray-100 h-14 animate-pulse" />)}
        </div>
      ) : (
        <div className="bg-white rounded-2xl border border-gray-100 overflow-hidden">
          <table className="w-full text-sm">
            <thead className="bg-gray-50 border-b border-gray-100">
              <tr>
                <th className="text-left px-4 py-3 font-medium text-gray-600">Название</th>
                <th className="text-left px-4 py-3 font-medium text-gray-600">Тип</th>
                <th className="text-left px-4 py-3 font-medium text-gray-600">Цена/ч</th>
                <th className="text-left px-4 py-3 font-medium text-gray-600">Статус</th>
                <th className="px-4 py-3" />
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-50">
              {peripherals.map((p) => (
                <tr key={p.id} className="hover:bg-gray-50">
                  <td className="px-4 py-3 font-medium text-gray-900">{p.name}</td>
                  <td className="px-4 py-3 text-gray-500">{TYPE_LABELS[p.type] || p.type}</td>
                  <td className="px-4 py-3 text-gray-700">{p.price_per_hour.toLocaleString()} ₸</td>
                  <td className="px-4 py-3">
                    <span className={`text-xs font-medium px-2 py-1 rounded-full ${
                      p.status === 'available' ? 'bg-green-100 text-green-700' : 'bg-yellow-100 text-yellow-700'
                    }`}>
                      {p.status === 'available' ? 'Доступна' : p.status === 'booked' ? 'Занята' : 'Обслуживание'}
                    </span>
                  </td>
                  <td className="px-4 py-3 text-right">
                    <button
                      onClick={() => { if (confirm('Удалить периферию?')) deleteMutation.mutate(p.id); }}
                      className="p-1.5 text-gray-400 hover:text-red-600"
                    >
                      <Trash2 className="w-4 h-4" />
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
          {peripherals.length === 0 && (
            <div className="text-center py-12 text-gray-400">
              <p>Периферия не добавлена</p>
              <Link to="/admin/peripherals/new" className="text-primary-600 hover:underline text-sm mt-1 block">Добавить первую</Link>
            </div>
          )}
        </div>
      )}
    </div>
  );
}
