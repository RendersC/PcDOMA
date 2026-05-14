import { useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { CreditCard, CheckCircle, XCircle, ArrowLeft } from 'lucide-react';
import { bookingsApi } from '../api/bookings';
import { paymentsApi } from '../api/payments';
import toast from 'react-hot-toast';

const STATUS_LABELS: Record<string, string> = {
  pending_payment: 'Ожидает оплаты',
  paid: 'Оплачено',
  pending_worker: 'Ожидает работника',
  accepted: 'Принято',
  rejected: 'Отклонено',
  delivering: 'Доставляется',
  active: 'Активно',
  completed: 'Завершено',
  cancelled: 'Отменено',
};

export function PaymentPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [method, setMethod] = useState<'card' | 'kaspi'>('card');
  const [cardNumber, setCardNumber] = useState('');
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState<'success' | 'failed' | null>(null);

  const { data: bookingRes, isLoading } = useQuery({
    queryKey: ['booking', id],
    queryFn: () => bookingsApi.getById(id!),
    enabled: !!id,
  });

  const booking = bookingRes?.data;

  const formatCardNumber = (val: string) => {
    const digits = val.replace(/\D/g, '').slice(0, 16);
    return digits.replace(/(.{4})/g, '$1 ').trim();
  };

  const handlePay = async () => {
    if (!booking) return;
    if (method === 'card' && cardNumber.replace(/\s/g, '').length < 16) {
      toast.error('Введите номер карты');
      return;
    }

    setLoading(true);
    try {
      const paymentId = booking.payment_id;
      if (!paymentId) {
        toast.error('Платёж не найден');
        return;
      }
      await paymentsApi.process(paymentId, method);
      setResult('success');
    } catch {
      setResult('failed');
    } finally {
      setLoading(false);
    }
  };

  if (isLoading) {
    return (
      <div className="max-w-lg mx-auto px-4 py-16">
        <div className="animate-pulse space-y-4">
          <div className="h-8 bg-gray-200 rounded w-1/2" />
          <div className="h-48 bg-gray-200 rounded-2xl" />
        </div>
      </div>
    );
  }

  if (result === 'success') {
    return (
      <div className="max-w-lg mx-auto px-4 py-16 text-center">
        <div className="bg-white rounded-2xl border border-gray-100 p-10 shadow-sm">
          <CheckCircle className="w-16 h-16 text-green-500 mx-auto mb-4" />
          <h1 className="text-2xl font-bold text-gray-900 mb-2">Оплата прошла!</h1>
          <p className="text-gray-500 mb-2">Ваш заказ передан работникам.</p>
          <p className="text-gray-400 text-sm mb-8">Они свяжутся с вами и доставят оборудование.</p>
          <div className="space-y-3">
            <button
              onClick={() => navigate('/my-bookings')}
              className="w-full bg-primary-600 hover:bg-primary-700 text-white py-3 rounded-xl font-semibold text-sm transition-colors"
            >
              Мои аренды
            </button>
            <button
              onClick={() => navigate('/')}
              className="w-full border border-gray-200 text-gray-600 hover:text-gray-800 py-3 rounded-xl font-medium text-sm transition-colors"
            >
              На главную
            </button>
          </div>
        </div>
      </div>
    );
  }

  if (result === 'failed') {
    return (
      <div className="max-w-lg mx-auto px-4 py-16 text-center">
        <div className="bg-white rounded-2xl border border-gray-100 p-10 shadow-sm">
          <XCircle className="w-16 h-16 text-red-500 mx-auto mb-4" />
          <h1 className="text-2xl font-bold text-gray-900 mb-2">Ошибка оплаты</h1>
          <p className="text-gray-500 mb-8">Платёж не прошёл. Попробуйте ещё раз.</p>
          <div className="space-y-3">
            <button
              onClick={() => setResult(null)}
              className="w-full bg-primary-600 hover:bg-primary-700 text-white py-3 rounded-xl font-semibold text-sm"
            >
              Попробовать снова
            </button>
            <button
              onClick={() => navigate('/my-bookings')}
              className="w-full border border-gray-200 text-gray-600 py-3 rounded-xl font-medium text-sm"
            >
              Мои аренды
            </button>
          </div>
        </div>
      </div>
    );
  }

  if (!booking) {
    return (
      <div className="max-w-lg mx-auto px-4 py-16 text-center">
        <p className="text-gray-500">Бронирование не найдено</p>
      </div>
    );
  }

  return (
    <div className="max-w-lg mx-auto px-4 py-8">
      <button onClick={() => navigate(-1)} className="flex items-center gap-2 text-gray-500 hover:text-gray-700 mb-6 text-sm">
        <ArrowLeft className="w-4 h-4" /> Назад
      </button>

      <h1 className="text-2xl font-bold text-gray-900 mb-6">Оплата</h1>

      {/* Booking summary */}
      <div className="bg-white rounded-2xl border border-gray-100 p-6 mb-6">
        <h2 className="font-bold text-gray-900 mb-3">Детали заказа</h2>
        <div className="space-y-2 text-sm">
          <div className="flex justify-between">
            <span className="text-gray-500">Тип</span>
            <span className="font-medium">{booking.booking_type === 'custom' ? 'Конструктор' : 'Готовый сеап'}</span>
          </div>
          <div className="flex justify-between">
            <span className="text-gray-500">Тариф</span>
            <span className="font-medium">{booking.rental_type === 'hourly' ? 'Почасовой' : 'Посуточный'}</span>
          </div>
          <div className="flex justify-between">
            <span className="text-gray-500">Начало</span>
            <span className="font-medium">{new Date(booking.start_time).toLocaleString('ru-RU')}</span>
          </div>
          <div className="flex justify-between">
            <span className="text-gray-500">Конец</span>
            <span className="font-medium">{new Date(booking.end_time).toLocaleString('ru-RU')}</span>
          </div>
          <div className="flex justify-between">
            <span className="text-gray-500">Статус</span>
            <span className="font-medium">{STATUS_LABELS[booking.status] || booking.status}</span>
          </div>
        </div>
        <div className="border-t mt-4 pt-4 flex justify-between font-bold text-lg">
          <span>К оплате</span>
          <span className="text-primary-600">{booking.total_price.toLocaleString()} ₸</span>
        </div>
      </div>

      {/* Payment method */}
      <div className="bg-white rounded-2xl border border-gray-100 p-6 mb-6">
        <h2 className="font-bold text-gray-900 mb-4">Способ оплаты</h2>
        <div className="grid grid-cols-2 gap-3 mb-5">
          {(['card', 'kaspi'] as const).map((m) => (
            <button
              key={m}
              onClick={() => setMethod(m)}
              className={`p-4 rounded-xl border-2 text-left transition-all ${
                method === m ? 'border-primary-500 bg-primary-50' : 'border-gray-100 hover:border-gray-200'
              }`}
            >
              <CreditCard className={`w-5 h-5 mb-1 ${method === m ? 'text-primary-600' : 'text-gray-400'}`} />
              <p className="font-semibold text-sm text-gray-900">
                {m === 'card' ? 'Банковская карта' : 'Kaspi Pay'}
              </p>
            </button>
          ))}
        </div>

        {method === 'card' && (
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1.5">Номер карты</label>
              <input
                type="text"
                value={cardNumber}
                onChange={(e) => setCardNumber(formatCardNumber(e.target.value))}
                placeholder="0000 0000 0000 0000"
                className="w-full px-4 py-3 border border-gray-200 rounded-xl text-sm font-mono focus:outline-none focus:ring-2 focus:ring-primary-500"
              />
            </div>
            <div className="grid grid-cols-2 gap-3">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1.5">Срок</label>
                <input
                  type="text"
                  placeholder="MM/YY"
                  maxLength={5}
                  className="w-full px-4 py-3 border border-gray-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-primary-500"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1.5">CVV</label>
                <input
                  type="password"
                  placeholder="•••"
                  maxLength={3}
                  className="w-full px-4 py-3 border border-gray-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-primary-500"
                />
              </div>
            </div>
          </div>
        )}

        {method === 'kaspi' && (
          <div className="bg-yellow-50 rounded-xl p-4 text-center">
            <p className="text-yellow-700 text-sm font-medium">Оплата через Kaspi Pay</p>
            <p className="text-yellow-600 text-xs mt-1">Нажмите кнопку — откроется Kaspi приложение</p>
          </div>
        )}
      </div>

      <button
        onClick={handlePay}
        disabled={loading}
        className="w-full bg-primary-600 hover:bg-primary-700 disabled:opacity-50 text-white py-3.5 rounded-xl font-bold text-base transition-colors"
      >
        {loading ? 'Обработка...' : `Оплатить ${booking.total_price.toLocaleString()} ₸`}
      </button>

      <p className="text-xs text-gray-400 text-center mt-3">
        Демо-режим: платёж обрабатывается мгновенно
      </p>
    </div>
  );
}
