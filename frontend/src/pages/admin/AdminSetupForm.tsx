import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { ArrowLeft, Save, Plus, X } from 'lucide-react';
import { catalogApi } from '../../api/catalog';
import toast from 'react-hot-toast';

export function AdminSetupForm() {
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [form, setForm] = useState({
    name: '',
    category: 'gaming',
    description: '',
    pc_id: '',
    discount_percent: '10',
  });
  const [selectedPeripherals, setSelectedPeripherals] = useState<string[]>([]);

  const { data: pcsData } = useQuery({ queryKey: ['pcs-all'], queryFn: () => catalogApi.listPCs({ limit: 100 }) });
  const { data: peripheralsData } = useQuery({ queryKey: ['peripherals-all'], queryFn: () => catalogApi.listPeripherals({ limit: '100' }) });

  const pcs = pcsData?.data?.data || [];
  const peripherals = peripheralsData?.data?.data || [];

  const set = (field: string) => (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement>) =>
    setForm((prev) => ({ ...prev, [field]: e.target.value }));

  const togglePeripheral = (id: string) =>
    setSelectedPeripherals((prev) => prev.includes(id) ? prev.filter((x) => x !== id) : [...prev, id]);

  const selectedPC = pcs.find((p) => p.id === form.pc_id);
  const selectedPeripheralObjects = peripherals.filter((p) => selectedPeripherals.includes(p.id));

  const basePerHour = selectedPC?.price_per_hour || 0;
  const extrasPerHour = selectedPeripheralObjects.reduce((s, p) => s + p.price_per_hour, 0);
  const rawPerHour = basePerHour + extrasPerHour;
  const discount = Number(form.discount_percent) || 0;
  const totalPerHour = Math.round(rawPerHour * (1 - discount / 100));
  const totalPerDay = Math.round(((selectedPC?.price_per_day || 0) + selectedPeripheralObjects.reduce((s, p) => s + p.price_per_day, 0)) * (1 - discount / 100));

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!form.pc_id) { toast.error('Выберите ПК'); return; }
    setLoading(true);
    try {
      await catalogApi.createSetup({
        name: form.name,
        category: form.category,
        description: form.description,
        pc_id: form.pc_id,
        peripheral_ids: selectedPeripherals,
        discount_percent: Number(form.discount_percent),
        total_price_per_hour: totalPerHour,
        total_price_per_day: totalPerDay,
        currency: 'KZT',
      });
      toast.success('Сеап создан');
      navigate('/admin/setups');
    } catch {
      toast.error('Ошибка при создании');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="max-w-3xl mx-auto px-4 py-8 sm:px-6">
      <button onClick={() => navigate('/admin/setups')} className="flex items-center gap-2 text-gray-500 hover:text-gray-700 mb-6 text-sm">
        <ArrowLeft className="w-4 h-4" /> К списку сеапов
      </button>

      <h1 className="text-2xl font-bold text-gray-900 mb-6">Создать сеап</h1>

      <form onSubmit={handleSubmit} className="space-y-6">
        <div className="bg-white rounded-2xl border border-gray-100 p-6 space-y-4">
          <h2 className="font-bold text-gray-900">Основная информация</h2>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Название сеапа</label>
            <input type="text" required value={form.name} onChange={set('name')} placeholder="Ultimate Gaming Setup"
              className="w-full px-4 py-3 border border-gray-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-primary-500" />
          </div>
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Категория</label>
              <select value={form.category} onChange={set('category')}
                className="w-full px-4 py-3 border border-gray-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 bg-white">
                <option value="gaming">Игровой</option>
                <option value="work">Рабочий</option>
                <option value="design">Дизайн</option>
                <option value="streaming">Стриминг</option>
              </select>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Скидка (%)</label>
              <input type="number" min={0} max={50} value={form.discount_percent} onChange={set('discount_percent')}
                className="w-full px-4 py-3 border border-gray-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-primary-500" />
            </div>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Описание</label>
            <textarea value={form.description} onChange={set('description')} rows={3}
              placeholder="Описание сеапа..."
              className="w-full px-4 py-3 border border-gray-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 resize-none" />
          </div>
        </div>

        <div className="bg-white rounded-2xl border border-gray-100 p-6">
          <h2 className="font-bold text-gray-900 mb-4">Системный блок</h2>
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-2">
            {pcs.map((pc) => (
              <button
                key={pc.id}
                type="button"
                onClick={() => setForm((prev) => ({ ...prev, pc_id: pc.id }))}
                className={`text-left p-3 rounded-xl border-2 transition-all ${
                  form.pc_id === pc.id ? 'border-primary-500 bg-primary-50' : 'border-gray-100 hover:border-gray-200'
                }`}
              >
                <p className="font-medium text-sm text-gray-900">{pc.name}</p>
                <p className="text-xs text-gray-400">{pc.price_per_hour.toLocaleString()} ₸/ч</p>
              </button>
            ))}
          </div>
        </div>

        <div className="bg-white rounded-2xl border border-gray-100 p-6">
          <h2 className="font-bold text-gray-900 mb-4">Периферия</h2>
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-2">
            {peripherals.map((p) => {
              const selected = selectedPeripherals.includes(p.id);
              return (
                <button
                  key={p.id}
                  type="button"
                  onClick={() => togglePeripheral(p.id)}
                  className={`text-left p-3 rounded-xl border transition-all flex items-center justify-between ${
                    selected ? 'border-primary-400 bg-primary-50' : 'border-gray-100 hover:border-gray-200'
                  }`}
                >
                  <div>
                    <p className="text-sm font-medium text-gray-900">{p.name}</p>
                    <p className="text-xs text-gray-400">+{p.price_per_hour.toLocaleString()} ₸/ч</p>
                  </div>
                  {selected ? <X className="w-4 h-4 text-primary-600" /> : <Plus className="w-4 h-4 text-gray-300" />}
                </button>
              );
            })}
          </div>
        </div>

        {selectedPC && (
          <div className="bg-gray-50 rounded-2xl border border-gray-200 p-5">
            <h2 className="font-bold text-gray-700 mb-2">Итоговая цена</h2>
            <p className="text-sm text-gray-500">Базовая: {rawPerHour.toLocaleString()} ₸/ч</p>
            <p className="text-sm text-gray-500">Скидка: {discount}%</p>
            <p className="font-bold text-primary-600 text-lg mt-1">{totalPerHour.toLocaleString()} ₸/ч · {totalPerDay.toLocaleString()} ₸/сут</p>
          </div>
        )}

        <button
          type="submit"
          disabled={loading}
          className="w-full flex items-center justify-center gap-2 bg-primary-600 hover:bg-primary-700 disabled:opacity-50 text-white py-3 rounded-xl font-semibold text-sm transition-colors"
        >
          <Save className="w-4 h-4" />
          {loading ? 'Создаём...' : 'Создать сеап'}
        </button>
      </form>
    </div>
  );
}
