import { useState } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { Monitor, Mouse, Headphones, Tv, Plus, Minus, ChevronRight } from 'lucide-react';
import { catalogApi } from '../api/catalog';
import type { PC, Peripheral } from '../types';
import toast from 'react-hot-toast';

const PERIPHERAL_ICONS: Record<string, React.ReactNode> = {
  mouse: <Mouse className="w-5 h-5" />,
  headphones: <Headphones className="w-5 h-5" />,
  headset: <Headphones className="w-5 h-5" />,
  monitor: <Tv className="w-5 h-5" />,
  keyboard: <Monitor className="w-5 h-5" />,
};

const PERIPHERAL_LABELS: Record<string, string> = {
  mouse: 'Мышь',
  headphones: 'Наушники',
  headset: 'Гарнитура',
  monitor: 'Монитор',
  keyboard: 'Клавиатура',
};

export function ConfiguratorPage() {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const defaultPCId = searchParams.get('pc_id') || '';

  const [selectedPC, setSelectedPC] = useState<PC | null>(null);
  const [selectedPeripherals, setSelectedPeripherals] = useState<Peripheral[]>([]);

  const { data: pcsData } = useQuery({
    queryKey: ['pcs', 'available'],
    queryFn: () => catalogApi.listPCs({ status: 'available', limit: 50 }),
  });

  const { data: peripheralsData } = useQuery({
    queryKey: ['peripherals', 'all'],
    queryFn: () => catalogApi.listPeripherals({ limit: '100' }),
  });

  const pcs = pcsData?.data?.data || [];
  const peripherals = peripheralsData?.data?.data || [];

  const totalPerHour = (selectedPC?.price_per_hour || 0) +
    selectedPeripherals.reduce((sum, p) => sum + p.price_per_hour, 0);
  const totalPerDay = (selectedPC?.price_per_day || 0) +
    selectedPeripherals.reduce((sum, p) => sum + p.price_per_day, 0);

  const togglePeripheral = (p: Peripheral) => {
    setSelectedPeripherals((prev) =>
      prev.find((x) => x.id === p.id) ? prev.filter((x) => x.id !== p.id) : [...prev, p]
    );
  };

  const handleBook = () => {
    if (!selectedPC) {
      toast.error('Выберите системный блок');
      return;
    }
    const params = new URLSearchParams({
      type: 'custom',
      pc_id: selectedPC.id,
      peripheral_ids: selectedPeripherals.map((p) => p.id).join(','),
    });
    navigate(`/booking/${selectedPC.id}?${params}`);
  };

  const groupedPeripherals = peripherals.reduce<Record<string, Peripheral[]>>((acc, p) => {
    if (!acc[p.type]) acc[p.type] = [];
    acc[p.type].push(p);
    return acc;
  }, {});

  return (
    <div className="max-w-7xl mx-auto px-4 py-8 sm:px-6 lg:px-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900 mb-2">Конструктор сеапа</h1>
        <p className="text-gray-500">Выберите системный блок и добавьте нужную периферию</p>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <div className="lg:col-span-2 space-y-6">
          {/* Step 1: PC */}
          <div className="bg-white rounded-2xl border border-gray-100 p-6">
            <h2 className="font-bold text-gray-900 mb-1 flex items-center gap-2">
              <span className="w-6 h-6 bg-primary-600 text-white rounded-full text-xs flex items-center justify-center font-bold">1</span>
              Выберите системный блок
            </h2>
            <p className="text-sm text-gray-500 mb-4 ml-8">Основа вашего сеапа</p>

            <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
              {pcs.map((pc) => (
                <button
                  key={pc.id}
                  onClick={() => setSelectedPC(selectedPC?.id === pc.id ? null : pc)}
                  className={`text-left p-4 rounded-xl border-2 transition-all ${
                    selectedPC?.id === pc.id
                      ? 'border-primary-500 bg-primary-50'
                      : 'border-gray-100 hover:border-gray-200'
                  }`}
                >
                  <div className="flex items-start justify-between">
                    <div>
                      <p className="font-semibold text-sm text-gray-900">{pc.name}</p>
                      <p className="text-xs text-gray-500 mt-0.5">{pc.specs.gpu}</p>
                    </div>
                    {selectedPC?.id === pc.id && (
                      <div className="w-5 h-5 bg-primary-600 rounded-full flex items-center justify-center flex-shrink-0">
                        <span className="text-white text-xs">✓</span>
                      </div>
                    )}
                  </div>
                  <p className="text-primary-600 font-bold text-sm mt-2">
                    {pc.price_per_hour.toLocaleString()} ₸/ч
                  </p>
                </button>
              ))}
            </div>
          </div>

          {/* Step 2: Peripherals */}
          <div className="bg-white rounded-2xl border border-gray-100 p-6">
            <h2 className="font-bold text-gray-900 mb-1 flex items-center gap-2">
              <span className="w-6 h-6 bg-primary-600 text-white rounded-full text-xs flex items-center justify-center font-bold">2</span>
              Добавьте периферию
            </h2>
            <p className="text-sm text-gray-500 mb-4 ml-8">Необязательно — можно пропустить</p>

            <div className="space-y-6">
              {Object.entries(groupedPeripherals).map(([type, items]) => (
                <div key={type}>
                  <div className="flex items-center gap-2 mb-3">
                    <span className="text-gray-500">{PERIPHERAL_ICONS[type]}</span>
                    <h3 className="font-medium text-sm text-gray-700">{PERIPHERAL_LABELS[type] || type}</h3>
                  </div>
                  <div className="grid grid-cols-1 sm:grid-cols-2 gap-2">
                    {items.map((p) => {
                      const isSelected = selectedPeripherals.some((x) => x.id === p.id);
                      return (
                        <button
                          key={p.id}
                          onClick={() => togglePeripheral(p)}
                          className={`text-left p-3 rounded-xl border transition-all flex items-center justify-between ${
                            isSelected ? 'border-primary-400 bg-primary-50' : 'border-gray-100 hover:border-gray-200'
                          }`}
                        >
                          <div>
                            <p className="text-sm font-medium text-gray-900">{p.name}</p>
                            <p className="text-xs text-gray-500 mt-0.5">+{p.price_per_hour.toLocaleString()} ₸/ч</p>
                          </div>
                          <div className={`w-5 h-5 rounded-full border-2 flex items-center justify-center flex-shrink-0 ${
                            isSelected ? 'border-primary-600 bg-primary-600' : 'border-gray-300'
                          }`}>
                            {isSelected ? <span className="text-white text-xs">✓</span> : <Plus className="w-3 h-3 text-gray-400" />}
                          </div>
                        </button>
                      );
                    })}
                  </div>
                </div>
              ))}
              {Object.keys(groupedPeripherals).length === 0 && (
                <p className="text-sm text-gray-400 text-center py-4">Периферия не добавлена в каталог</p>
              )}
            </div>
          </div>
        </div>

        {/* Right: Summary */}
        <div>
          <div className="bg-white rounded-2xl border border-gray-100 p-6 sticky top-24 shadow-sm">
            <h2 className="font-bold text-gray-900 mb-4">Ваш сеап</h2>

            {!selectedPC ? (
              <p className="text-sm text-gray-400 text-center py-4">Выберите системный блок</p>
            ) : (
              <div className="space-y-3 mb-5">
                <div className="flex items-center justify-between text-sm">
                  <span className="font-medium text-gray-900">🖥 {selectedPC.name}</span>
                  <span className="text-gray-600">{selectedPC.price_per_hour.toLocaleString()} ₸/ч</span>
                </div>
                {selectedPeripherals.map((p) => (
                  <div key={p.id} className="flex items-center justify-between text-sm">
                    <div className="flex items-center gap-2">
                      <button onClick={() => togglePeripheral(p)} className="text-red-400 hover:text-red-600">
                        <Minus className="w-3 h-3" />
                      </button>
                      <span className="text-gray-700">{p.name}</span>
                    </div>
                    <span className="text-gray-600">+{p.price_per_hour.toLocaleString()} ₸/ч</span>
                  </div>
                ))}

                <div className="border-t pt-3 space-y-1">
                  <div className="flex justify-between font-bold text-gray-900 text-sm">
                    <span>Итого (час)</span>
                    <span className="text-primary-600">{totalPerHour.toLocaleString()} ₸</span>
                  </div>
                  <div className="flex justify-between text-xs text-gray-500">
                    <span>Итого (сутки)</span>
                    <span>{totalPerDay.toLocaleString()} ₸</span>
                  </div>
                </div>
              </div>
            )}

            <button
              onClick={handleBook}
              disabled={!selectedPC}
              className="w-full bg-primary-600 hover:bg-primary-700 disabled:opacity-40 disabled:cursor-not-allowed text-white py-3 rounded-xl font-semibold text-sm transition-colors flex items-center justify-center gap-2"
            >
              Перейти к бронированию <ChevronRight className="w-4 h-4" />
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
