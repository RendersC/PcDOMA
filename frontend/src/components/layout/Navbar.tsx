import { Link, useNavigate } from 'react-router-dom';
import { Monitor, Bell, User, Menu, X, LogOut } from 'lucide-react';
import { useState, useEffect } from 'react';
import { useAuthStore } from '../../store/authStore';
import { useNotificationStore } from '../../store/notificationStore';
import { notificationsApi } from '../../api/notifications';
import { authApi } from '../../api/auth';

export function Navbar() {
  const { user, isAuthenticated, logout } = useAuthStore();
  const { unreadCount, setNotifications } = useNotificationStore();
  const [isMenuOpen, setIsMenuOpen] = useState(false);
  const navigate = useNavigate();

  useEffect(() => {
    if (!isAuthenticated) return;
    const load = async () => {
      try {
        const res = await notificationsApi.list();
        setNotifications(res.data.data || []);
      } catch {}
    };
    load();
    const interval = setInterval(load, 10000);
    return () => clearInterval(interval);
  }, [isAuthenticated]);

  const handleLogout = async () => {
    const refreshToken = localStorage.getItem('refresh_token') || '';
    await authApi.logout(refreshToken).catch(() => {});
    logout();
    navigate('/');
  };

  return (
    <nav className="bg-white shadow-sm border-b border-gray-100 sticky top-0 z-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between h-16">
          <div className="flex items-center gap-8">
            <Link to="/" className="flex items-center gap-2 font-bold text-xl text-primary-700">
              <Monitor className="w-6 h-6" />
              <span>PCDoma</span>
            </Link>
            <div className="hidden md:flex items-center gap-6">
              <Link to="/catalog" className="text-gray-600 hover:text-primary-600 text-sm font-medium transition-colors">
                Каталог ПК
              </Link>
              <Link to="/setups" className="text-gray-600 hover:text-primary-600 text-sm font-medium transition-colors">
                Готовые сеапы
              </Link>
              <Link to="/configurator" className="text-gray-600 hover:text-primary-600 text-sm font-medium transition-colors">
                Конструктор
              </Link>
            </div>
          </div>

          <div className="hidden md:flex items-center gap-4">
            {isAuthenticated ? (
              <>
                {user?.role === 'worker' && (
                  <Link to="/worker" className="text-sm font-medium text-accent-600 hover:text-accent-700">
                    Мои заказы
                  </Link>
                )}
                {user?.role === 'admin' && (
                  <Link to="/admin" className="text-sm font-medium text-red-600 hover:text-red-700">
                    Админ
                  </Link>
                )}
                <Link to="/notifications" className="relative p-2 text-gray-500 hover:text-primary-600">
                  <Bell className="w-5 h-5" />
                  {unreadCount > 0 && (
                    <span className="absolute -top-1 -right-1 bg-red-500 text-white text-xs rounded-full w-4 h-4 flex items-center justify-center">
                      {unreadCount > 9 ? '9+' : unreadCount}
                    </span>
                  )}
                </Link>
                <Link to="/my-bookings" className="text-sm font-medium text-gray-600 hover:text-primary-600">
                  Мои аренды
                </Link>
                <div className="flex items-center gap-2">
                  <Link to="/profile" className="flex items-center gap-2 text-sm text-gray-700 hover:text-primary-600">
                    <User className="w-4 h-4" />
                    {user?.name}
                  </Link>
                  <button onClick={handleLogout} className="p-1 text-gray-400 hover:text-red-500">
                    <LogOut className="w-4 h-4" />
                  </button>
                </div>
              </>
            ) : (
              <>
                <Link to="/login" className="text-sm text-gray-600 hover:text-primary-600">
                  Войти
                </Link>
                <Link
                  to="/register"
                  className="bg-primary-600 text-white text-sm px-4 py-2 rounded-lg hover:bg-primary-700 transition-colors"
                >
                  Зарегистрироваться
                </Link>
              </>
            )}
          </div>

          <button className="md:hidden p-2" onClick={() => setIsMenuOpen(!isMenuOpen)}>
            {isMenuOpen ? <X className="w-5 h-5" /> : <Menu className="w-5 h-5" />}
          </button>
        </div>

        {isMenuOpen && (
          <div className="md:hidden py-3 border-t border-gray-100 space-y-2">
            <Link to="/catalog" className="block py-2 text-sm text-gray-600" onClick={() => setIsMenuOpen(false)}>Каталог ПК</Link>
            <Link to="/setups" className="block py-2 text-sm text-gray-600" onClick={() => setIsMenuOpen(false)}>Готовые сеапы</Link>
            <Link to="/configurator" className="block py-2 text-sm text-gray-600" onClick={() => setIsMenuOpen(false)}>Конструктор</Link>
            {isAuthenticated ? (
              <>
                <Link to="/my-bookings" className="block py-2 text-sm text-gray-600" onClick={() => setIsMenuOpen(false)}>Мои аренды</Link>
                <button onClick={handleLogout} className="block py-2 text-sm text-red-500 w-full text-left">Выйти</button>
              </>
            ) : (
              <>
                <Link to="/login" className="block py-2 text-sm text-gray-600" onClick={() => setIsMenuOpen(false)}>Войти</Link>
                <Link to="/register" className="block py-2 text-sm text-primary-600 font-medium" onClick={() => setIsMenuOpen(false)}>Зарегистрироваться</Link>
              </>
            )}
          </div>
        )}
      </div>
    </nav>
  );
}
