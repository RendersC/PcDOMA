import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { ArrowLeft, Save } from 'lucide-react';
import { catalogApi } from '../../api/catalog';
import toast from 'react-hot-toast';

const PERIPHERAL_TYPES = ['mouse', 'headphones', 'headset', 'monitor', 'keyboard'];

export function AdminPeripheralForm() {
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [form, setForm] = useState({
    name: '',
    type: 'mouse',
    description: '',
    price_per_hour: '',
    price_per_day: '',
  });

  const set = (field: string) => (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement>) =>
    setForm((prev) => ({ ...prev, [field]: e.target.value }));

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    try {
      await catalogApi.createPeripheral({
        name: form.name,
        type: form.type,
        description: form.description,
        specs: {},
        price_per_hour: Number(form.price_per_hour),
        price_per_day: Number(form.price_per_day),
        currency: 'KZT',
        status: 'available',
      });
      toast.success('Периферия добавлена');
      navigate('/admin/peripherals');
    } catch {
      toast.error('Ошибка при сохранении');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="max-w-2xl mx-auto px-4 py-8 sm:px-6">
      <button onClick={() => navigate('/admin/peripherals')} className="flex items-center gap-2 text-gray-500 hover:text-gray-700 mb-6 text-sm">
        <ArrowLeft className="w-4 h-4" /> К списку
      </button>

      <h1 className="text-2xl font-bold text-gray-900 mb-6">Добавить периферию</h1>

      <form onSubmit={handleSubmit} className="bg-white rounded-2xl border border-gray-100 p-6 space-y-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Название</label>
          <input type="text" required value={form.name} onChange={set('name')}
            placeholder="Razer DeathAdder V3"
            className="w-full px-4 py-3 border border-gray-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-primary-500" />
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Тип</label>
          <select value={form.type} onChange={set('type')}
            className="w-full px-4 py-3 border border-gray-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 bg-white">
            {PERIPHERAL_TYPES.map((t) => (
              <option key={t} value={t}>
                {t === 'mouse' ? 'Мышь' : t === 'headphones' ? 'Наушники' : t === 'headset' ? 'Гарнитура' : t === 'monitor' ? 'Монитор' : 'Клавиатура'}
              </option>
            ))}
          </select>
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Описание</label>
          <textarea value={form.description} onChange={set('description')} rows={3}
            placeholder="Описание устройства..."
            className="w-full px-4 py-3 border border-gray-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 resize-none" />
        </div>

        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Цена за час (₸)</label>
            <input type="number" required min={0} value={form.price_per_hour} onChange={set('price_per_hour')}
              placeholder="200"
              className="w-full px-4 py-3 border border-gray-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-primary-500" />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Цена за сутки (₸)</label>
            <input type="number" required min={0} value={form.price_per_day} onChange={set('price_per_day')}
              placeholder="800"
              className="w-full px-4 py-3 border border-gray-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-primary-500" />
          </div>
        </div>

        <button
          type="submit"
          disabled={loading}
          className="w-full flex items-center justify-center gap-2 bg-primary-600 hover:bg-primary-700 disabled:opacity-50 text-white py-3 rounded-xl font-semibold text-sm transition-colors"
        >
          <Save className="w-4 h-4" />
          {loading ? 'Сохраняем...' : 'Добавить периферию'}
        </button>
      </form>
    </div>
  );
}
