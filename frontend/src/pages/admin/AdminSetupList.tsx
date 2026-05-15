import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { Link } from 'react-router-dom';
import { Plus, Trash2, Star } from 'lucide-react';
import { catalogApi } from '../../api/catalog';
import toast from 'react-hot-toast';

const CATEGORY_LABELS: Record<string, string> = {
  gaming: 'Игровой',
  work: 'Рабочий',
  design: 'Дизайн',
  streaming: 'Стриминг',
};

export function AdminSetupList() {
  const queryClient = useQueryClient();

  const { data, isLoading } = useQuery({
    queryKey: ['admin-setups'],
    queryFn: () => catalogApi.listSetups(),
  });

  const deleteMutation = useMutation({
    mutationFn: (id: string) => catalogApi.deleteSetup(id),
    onSuccess: () => {
      toast.success('Сеап удалён');
      queryClient.invalidateQueries({ queryKey: ['admin-setups'] });
    },
    onError: () => toast.error('Не удалось удалить'),
  });

  const setups = data?.data?.data || [];

  return (
    <div className="max-w-6xl mx-auto px-4 py-8 sm:px-6">
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold text-gray-900">Готовые сеапы</h1>
        <Link
          to="/admin/setups/new"
          className="flex items-center gap-2 bg-primary-600 hover:bg-primary-700 text-white px-4 py-2.5 rounded-xl text-sm font-medium transition-colors"
        >
          <Plus className="w-4 h-4" /> Добавить сеап
        </Link>
      </div>

      {isLoading ? (
        <div className="space-y-3">
          {[...Array(3)].map((_, i) => <div key={i} className="bg-white rounded-xl border border-gray-100 h-14 animate-pulse" />)}
        </div>
      ) : (
        <div className="bg-white rounded-2xl border border-gray-100 overflow-hidden">
          <table className="w-full text-sm">
            <thead className="bg-gray-50 border-b border-gray-100">
              <tr>
                <th className="text-left px-4 py-3 font-medium text-gray-600">Название</th>
                <th className="text-left px-4 py-3 font-medium text-gray-600">Категория</th>
                <th className="text-left px-4 py-3 font-medium text-gray-600">Цена/ч</th>
                <th className="text-left px-4 py-3 font-medium text-gray-600">Скидка</th>
                <th className="text-left px-4 py-3 font-medium text-gray-600">Рейтинг</th>
                <th className="px-4 py-3" />
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-50">
              {setups.map((s) => (
                <tr key={s.id} className="hover:bg-gray-50">
                  <td className="px-4 py-3 font-medium text-gray-900">{s.name}</td>
                  <td className="px-4 py-3 text-gray-500">{CATEGORY_LABELS[s.category] || s.category}</td>
                  <td className="px-4 py-3 text-gray-700">{s.total_price_per_hour.toLocaleString()} ₸</td>
                  <td className="px-4 py-3">
                    {s.discount_percent > 0 ? (
                      <span className="text-xs font-medium px-2 py-1 rounded-full bg-green-100 text-green-700">
                        -{s.discount_percent}%
                      </span>
                    ) : (
                      <span className="text-gray-400">—</span>
                    )}
                  </td>
                  <td className="px-4 py-3">
                    <div className="flex items-center gap-1">
                      <Star className="w-3.5 h-3.5 text-yellow-400 fill-yellow-400" />
                      <span className="text-gray-700">{s.rating_avg.toFixed(1)}</span>
                    </div>
                  </td>
                  <td className="px-4 py-3 text-right">
                    <button
                      onClick={() => { if (confirm('Удалить сеап?')) deleteMutation.mutate(s.id); }}
                      className="p-1.5 text-gray-400 hover:text-red-600"
                    >
                      <Trash2 className="w-4 h-4" />
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
          {setups.length === 0 && (
            <div className="text-center py-12 text-gray-400">
              <p>Сеапы не добавлены</p>
              <Link to="/admin/setups/new" className="text-primary-600 hover:underline text-sm mt-1 block">Создать первый сеап</Link>
            </div>
          )}
        </div>
      )}
    </div>
  );
}
