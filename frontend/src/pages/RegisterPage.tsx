import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import toast from 'react-hot-toast';
import { authApi } from '../api/auth';
import { useAuthStore } from '../store/authStore';

export function RegisterPage() {
  const [form, setForm] = useState({ name: '', email: '', password: '', phone: '' });
  const [loading, setLoading] = useState(false);
  const { setUser, setTokens } = useAuthStore();
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    try {
      const { data } = await authApi.register(form);
      setTokens(data.tokens.access_token, data.tokens.refresh_token);
      setUser(data.user);
      toast.success('Аккаунт создан!');
      navigate('/catalog');
    } catch (err: unknown) {
      const msg = (err as { response?: { data?: { error?: string } } })?.response?.data?.error;
      toast.error(msg || 'Ошибка регистрации');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-[80vh] flex items-center justify-center px-4 py-8">
      <div className="w-full max-w-md">
        <div className="text-center mb-8">
          <h1 className="text-2xl font-bold text-gray-900">Создать аккаунт</h1>
          <p className="text-gray-500 mt-1 text-sm">Начните арендовать ПК уже сегодня</p>
        </div>

        <form onSubmit={handleSubmit} className="bg-white rounded-2xl shadow-sm border border-gray-100 p-8 space-y-5">
          {(['name', 'email', 'phone'] as const).map((field) => (
            <div key={field}>
              <label className="block text-sm font-medium text-gray-700 mb-1.5 capitalize">
                {field === 'name' ? 'Имя' : field === 'email' ? 'Email' : 'Телефон'}
                {field !== 'phone' && <span className="text-red-500 ml-1">*</span>}
              </label>
              <input
                type={field === 'email' ? 'email' : 'text'}
                value={form[field]}
                onChange={(e) => setForm((f) => ({ ...f, [field]: e.target.value }))}
                required={field !== 'phone'}
                className="w-full px-4 py-3 border border-gray-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-primary-500"
                placeholder={field === 'phone' ? '+7 (700) 123-45-67' : ''}
              />
            </div>
          ))}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1.5">
              Пароль <span className="text-red-500">*</span>
            </label>
            <input
              type="password"
              value={form.password}
              onChange={(e) => setForm((f) => ({ ...f, password: e.target.value }))}
              required
              minLength={6}
              className="w-full px-4 py-3 border border-gray-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-primary-500"
              placeholder="Минимум 6 символов"
            />
          </div>

          <button
            type="submit"
            disabled={loading}
            className="w-full bg-primary-600 text-white py-3 rounded-xl font-medium text-sm hover:bg-primary-700 transition-colors disabled:opacity-50"
          >
            {loading ? 'Создаём...' : 'Зарегистрироваться'}
          </button>

          <p className="text-center text-sm text-gray-500">
            Уже есть аккаунт?{' '}
            <Link to="/login" className="text-primary-600 hover:text-primary-700 font-medium">
              Войти
            </Link>
          </p>
        </form>
      </div>
    </div>
  );
}
