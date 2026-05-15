import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { ArrowLeft, Save } from 'lucide-react';
import { catalogApi } from '../../api/catalog';
import toast from 'react-hot-toast';

interface PCFormData {
  name: string;
  cpu: string;
  gpu: string;
  ram: string;
  storage: string;
  price_per_hour: string;
  price_per_day: string;
  district: string;
  address: string;
}

const EMPTY: PCFormData = {
  name: '', cpu: '', gpu: '', ram: '', storage: '',
  price_per_hour: '', price_per_day: '', district: '', address: '',
};

export function AdminPCForm() {
  const { id } = useParams<{ id?: string }>();
  const navigate = useNavigate();
  const isEdit = !!id && id !== 'new';
  const [form, setForm] = useState<PCFormData>(EMPTY);
  const [loading, setLoading] = useState(false);

  const { data: pcRes } = useQuery({
    queryKey: ['pc-edit', id],
    queryFn: () => catalogApi.getPC(id!),
    enabled: isEdit,
  });

  useEffect(() => {
    if (pcRes?.data) {
      const pc = pcRes.data;
      setForm({
        name: pc.name,
        cpu: pc.specs.cpu,
        gpu: pc.specs.gpu,
        ram: pc.specs.ram,
        storage: pc.specs.storage,
        price_per_hour: String(pc.price_per_hour),
        price_per_day: String(pc.price_per_day),
        district: pc.location?.district || '',
        address: pc.location?.address || '',
      });
    }
  }, [pcRes]);

  const set = (field: keyof PCFormData) => (e: React.ChangeEvent<HTMLInputElement>) =>
    setForm((prev) => ({ ...prev, [field]: e.target.value }));

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    try {
      const payload = {
        name: form.name,
        specs: { cpu: form.cpu, gpu: form.gpu, ram: form.ram, storage: form.storage },
        price_per_hour: Number(form.price_per_hour),
        price_per_day: Number(form.price_per_day),
        currency: 'KZT',
        location: { district: form.district, address: form.address, lat: 0, lng: 0 },
        status: 'available',
      };
      if (isEdit) {
        await catalogApi.updatePC(id!, payload);
        toast.success('ПК обновлён');
      } else {
        await catalogApi.createPC(payload);
        toast.success('ПК добавлен');
      }
      navigate('/admin/pcs');
    } catch {
      toast.error('Ошибка при сохранении');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="max-w-2xl mx-auto px-4 py-8 sm:px-6">
      <button onClick={() => navigate('/admin/pcs')} className="flex items-center gap-2 text-gray-500 hover:text-gray-700 mb-6 text-sm">
        <ArrowLeft className="w-4 h-4" /> К списку ПК
      </button>

      <h1 className="text-2xl font-bold text-gray-900 mb-6">{isEdit ? 'Редактировать ПК' : 'Добавить ПК'}</h1>

      <form onSubmit={handleSubmit} className="space-y-6">
        <div className="bg-white rounded-2xl border border-gray-100 p-6 space-y-4">
          <h2 className="font-bold text-gray-900">Основная информация</h2>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Название</label>
            <input type="text" required value={form.name} onChange={set('name')}
              placeholder="Gaming PC Pro" className="w-full px-4 py-3 border border-gray-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-primary-500" />
          </div>
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Цена за час (₸)</label>
              <input type="number" required min={0} value={form.price_per_hour} onChange={set('price_per_hour')}
                placeholder="1500" className="w-full px-4 py-3 border border-gray-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-primary-500" />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Цена за сутки (₸)</label>
              <input type="number" required min={0} value={form.price_per_day} onChange={set('price_per_day')}
                placeholder="8000" className="w-full px-4 py-3 border border-gray-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-primary-500" />
            </div>
          </div>
        </div>

        <div className="bg-white rounded-2xl border border-gray-100 p-6 space-y-4">
          <h2 className="font-bold text-gray-900">Характеристики</h2>
          {(['cpu', 'gpu', 'ram', 'storage'] as const).map((field) => (
            <div key={field}>
              <label className="block text-sm font-medium text-gray-700 mb-1 uppercase">{field}</label>
              <input type="text" required value={form[field]} onChange={set(field)}
                placeholder={field === 'cpu' ? 'Intel Core i9-13900K' : field === 'gpu' ? 'NVIDIA RTX 4090' : field === 'ram' ? '32GB DDR5' : '2TB NVMe SSD'}
                className="w-full px-4 py-3 border border-gray-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-primary-500" />
            </div>
          ))}
        </div>

        <div className="bg-white rounded-2xl border border-gray-100 p-6 space-y-4">
          <h2 className="font-bold text-gray-900">Расположение</h2>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Район</label>
            <input type="text" value={form.district} onChange={set('district')}
              placeholder="Есиль" className="w-full px-4 py-3 border border-gray-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-primary-500" />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Адрес</label>
            <input type="text" value={form.address} onChange={set('address')}
              placeholder="ул. Кунаева 14" className="w-full px-4 py-3 border border-gray-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-primary-500" />
          </div>
        </div>

        <button
          type="submit"
          disabled={loading}
          className="w-full flex items-center justify-center gap-2 bg-primary-600 hover:bg-primary-700 disabled:opacity-50 text-white py-3 rounded-xl font-semibold text-sm transition-colors"
        >
          <Save className="w-4 h-4" />
          {loading ? 'Сохраняем...' : isEdit ? 'Сохранить изменения' : 'Добавить ПК'}
        </button>
      </form>
    </div>
  );
}
