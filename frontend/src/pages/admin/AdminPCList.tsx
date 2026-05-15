import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { Link } from 'react-router-dom';
import { Plus, Pencil, Trash2, ToggleLeft, ToggleRight } from 'lucide-react';
import { catalogApi } from '../../api/catalog';
import toast from 'react-hot-toast';

const STATUS_COLORS = {
  available: 'bg-green-100 text-green-700',
  booked: 'bg-blue-100 text-blue-700',
  maintenance: 'bg-yellow-100 text-yellow-700',
};

const STATUS_LABELS = {
  available: 'Доступен',
  booked: 'Занят',
  maintenance: 'Обслуживание',
};

export function AdminPCList() {
  const queryClient = useQueryClient();

  const { data, isLoading } = useQuery({
    queryKey: ['admin-pcs'],
    queryFn: () => catalogApi.listPCs({ limit: 100 }),
  });

  const deleteMutation = useMutation({
    mutationFn: (id: string) => catalogApi.deletePC(id),
    onSuccess: () => {
      toast.success('ПК удалён');
      queryClient.invalidateQueries({ queryKey: ['admin-pcs'] });
    },
    onError: () => toast.error('Не удалось удалить'),
  });

  const statusMutation = useMutation({
    mutationFn: ({ id, status }: { id: string; status: string }) => catalogApi.updatePCStatus(id, status),
    onSuccess: () => {
      toast.success('Статус обновлён');
      queryClient.invalidateQueries({ queryKey: ['admin-pcs'] });
    },
  });

  const pcs = data?.data?.data || [];

  return (
    <div className="max-w-6xl mx-auto px-4 py-8 sm:px-6">
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold text-gray-900">Системные блоки</h1>
        <Link
          to="/admin/pcs/new"
          className="flex items-center gap-2 bg-primary-600 hover:bg-primary-700 text-white px-4 py-2.5 rounded-xl text-sm font-medium transition-colors"
        >
          <Plus className="w-4 h-4" /> Добавить ПК
        </Link>
      </div>

      {isLoading ? (
        <div className="space-y-3">
          {[...Array(4)].map((_, i) => (
            <div key={i} className="bg-white rounded-xl border border-gray-100 p-4 animate-pulse h-16" />
          ))}
        </div>
      ) : (
        <div className="bg-white rounded-2xl border border-gray-100 overflow-hidden">
          <table className="w-full text-sm">
            <thead className="bg-gray-50 border-b border-gray-100">
              <tr>
                <th className="text-left px-4 py-3 font-medium text-gray-600">Название</th>
                <th className="text-left px-4 py-3 font-medium text-gray-600 hidden md:table-cell">GPU</th>
                <th className="text-left px-4 py-3 font-medium text-gray-600">Цена/ч</th>
                <th className="text-left px-4 py-3 font-medium text-gray-600">Статус</th>
                <th className="px-4 py-3" />
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-50">
              {pcs.map((pc) => (
                <tr key={pc.id} className="hover:bg-gray-50">
                  <td className="px-4 py-3 font-medium text-gray-900">{pc.name}</td>
                  <td className="px-4 py-3 text-gray-500 hidden md:table-cell">{pc.specs.gpu}</td>
                  <td className="px-4 py-3 text-gray-700">{pc.price_per_hour.toLocaleString()} ₸</td>
                  <td className="px-4 py-3">
                    <span className={`text-xs font-medium px-2 py-1 rounded-full ${STATUS_COLORS[pc.status]}`}>
                      {STATUS_LABELS[pc.status]}
                    </span>
                  </td>
                  <td className="px-4 py-3">
                    <div className="flex items-center justify-end gap-1">
                      <button
                        onClick={() => statusMutation.mutate({
                          id: pc.id,
                          status: pc.status === 'available' ? 'maintenance' : 'available',
                        })}
                        className="p-1.5 text-gray-400 hover:text-gray-600"
                        title={pc.status === 'available' ? 'Поставить на обслуживание' : 'Сделать доступным'}
                      >
                        {pc.status === 'available' ? <ToggleRight className="w-4 h-4 text-green-500" /> : <ToggleLeft className="w-4 h-4" />}
                      </button>
                      <Link to={`/admin/pcs/${pc.id}/edit`} className="p-1.5 text-gray-400 hover:text-blue-600">
                        <Pencil className="w-4 h-4" />
                      </Link>
                      <button
                        onClick={() => { if (confirm('Удалить ПК?')) deleteMutation.mutate(pc.id); }}
                        className="p-1.5 text-gray-400 hover:text-red-600"
                      >
                        <Trash2 className="w-4 h-4" />
                      </button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
          {pcs.length === 0 && (
            <div className="text-center py-12 text-gray-400">
              <p>ПК не добавлены</p>
              <Link to="/admin/pcs/new" className="text-primary-600 hover:underline text-sm mt-1 block">Добавить первый ПК</Link>
            </div>
          )}
        </div>
      )}
    </div>
  );
}
