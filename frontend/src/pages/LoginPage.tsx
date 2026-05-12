import { useState } from 'react';
import { Link, useNavigate, useLocation } from 'react-router-dom';
import { Monitor } from 'lucide-react';
import toast from 'react-hot-toast';
import { authApi } from '../api/auth';
import { useAuthStore } from '../store/authStore';

export function LoginPage() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [loading, setLoading] = useState(false);
  const { setUser, setTokens } = useAuthStore();
  const navigate = useNavigate();
  const location = useLocation();
  const from = (location.state as { from?: { pathname: string } })?.from?.pathname || '/';

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    try {
      const { data } = await authApi.login({ email, password });
      setTokens(data.tokens.access_token, data.tokens.refresh_token);
      setUser(data.user);
      toast.success('Добро пожаловать!');
      if (data.user.role === 'worker') {
        navigate('/worker');
      } else if (data.user.role === 'admin') {
        navigate('/admin');
      } else {
        navigate(from);
      }
    } catch {
      toast.error('Неверный email или пароль');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-[80vh] flex items-center justify-center px-4">
      <div className="w-full max-w-md">
        <div className="text-center mb-8">
          <div className="flex justify-center mb-4">
            <Monitor className="w-10 h-10 text-primary-600" />
          </div>
          <h1 className="text-2xl font-bold text-gray-900">Войдите в PCDoma</h1>
          <p className="text-gray-500 mt-1 text-sm">Арендуйте ПК в Астане</p>
        </div>

        <form onSubmit={handleSubmit} className="bg-white rounded-2xl shadow-sm border border-gray-100 p-8 space-y-5">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1.5">Email</label>
            <input
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
              className="w-full px-4 py-3 border border-gray-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent"
              placeholder="your@email.com"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1.5">Пароль</label>
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
              className="w-full px-4 py-3 border border-gray-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent"
              placeholder="••••••••"
            />
          </div>

          <button
            type="submit"
            disabled={loading}
            className="w-full bg-primary-600 text-white py-3 rounded-xl font-medium text-sm hover:bg-primary-700 transition-colors disabled:opacity-50"
          >
            {loading ? 'Входим...' : 'Войти'}
          </button>

          <p className="text-center text-sm text-gray-500">
            Нет аккаунта?{' '}
            <Link to="/register" className="text-primary-600 hover:text-primary-700 font-medium">
              Зарегистрироваться
            </Link>
          </p>
        </form>
      </div>
    </div>
  );
}
