import { Link, useNavigate, useLocation } from 'react-router-dom';
import { Monitor, Bell, User, LogOut } from 'lucide-react';
import { useState, useEffect } from 'react';
import { useAuthStore } from '../../store/authStore';
import { useNotificationStore } from '../../store/notificationStore';
import { notificationsApi } from '../../api/notifications';
import { authApi } from '../../api/auth';
import { useScroll } from '../ui/useScroll';
import { MenuToggleIcon } from '../ui/MenuToggleIcon';
import { cn } from '../../lib/utils';

const navLinks = [
  { to: '/catalog', label: 'Каталог ПК' },
  { to: '/setups', label: 'Готовые сетапы' },
  { to: '/configurator', label: 'Конструктор' },
];

export function Navbar() {
  const { user, isAuthenticated, logout } = useAuthStore();
  const { unreadCount, setNotifications } = useNotificationStore();
  const [open, setOpen] = useState(false);
  const scrolled = useScroll(10);
  const navigate = useNavigate();
  const { pathname } = useLocation();

  // Over the homepage hero (at the top, menu closed) the bar is transparent
  // with light text so the video shows through. Everywhere else it's a
  // translucent blurred bar with dark text.
  const transparent = pathname === '/' && !scrolled && !open;

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

  // Lock body scroll while the mobile menu is open.
  useEffect(() => {
    document.body.style.overflow = open ? 'hidden' : '';
    return () => {
      document.body.style.overflow = '';
    };
  }, [open]);

  const handleLogout = async () => {
    const refreshToken = localStorage.getItem('refresh_token') || '';
    await authApi.logout(refreshToken).catch(() => {});
    logout();
    setOpen(false);
    navigate('/');
  };

  const linkCls = transparent
    ? 'text-white/80 hover:text-white'
    : 'text-gray-600 hover:text-primary-600';
  const iconCls = transparent
    ? 'text-white/80 hover:text-white'
    : 'text-gray-500 hover:text-primary-600';

  return (
    <header
      className={cn(
        'fixed inset-x-0 top-0 z-50 w-full transition-all duration-300 ease-out',
        transparent ? 'bg-transparent' : 'border-b border-gray-200/50 bg-white/70 backdrop-blur-lg',
      )}
    >
      <div className="mx-auto max-w-7xl px-4">
        <nav
          className={cn(
            'flex w-full items-center justify-between transition-all duration-300 ease-out',
            scrolled ? 'h-14' : 'h-16',
          )}
        >
          {/* Brand + primary links */}
          <div className="flex items-center gap-8">
            <Link
              to="/"
              className={cn('flex items-center gap-2 text-xl font-bold', transparent ? 'text-white' : 'text-primary-700')}
              onClick={() => setOpen(false)}
            >
              <Monitor className="h-6 w-6" />
              <span>PCDoma</span>
            </Link>
            <div className="hidden items-center gap-6 md:flex">
              {navLinks.map((link) => (
                <Link key={link.to} to={link.to} className={cn('text-sm font-medium transition-colors', linkCls)}>
                  {link.label}
                </Link>
              ))}
            </div>
          </div>

          {/* Desktop auth area */}
          <div className="hidden items-center gap-4 md:flex">
            {isAuthenticated ? (
              <>
                {user?.role === 'worker' && (
                  <Link to="/worker" className={cn('text-sm font-medium', transparent ? 'text-amber-300 hover:text-amber-200' : 'text-accent-600 hover:text-accent-700')}>
                    Мои заказы
                  </Link>
                )}
                {user?.role === 'admin' && (
                  <Link to="/admin" className={cn('text-sm font-medium', transparent ? 'text-red-300 hover:text-red-200' : 'text-red-600 hover:text-red-700')}>
                    Админ
                  </Link>
                )}
                <Link to="/notifications" className={cn('relative p-2', iconCls)}>
                  <Bell className="h-5 w-5" />
                  {unreadCount > 0 && (
                    <span className="absolute -right-1 -top-1 flex h-4 w-4 items-center justify-center rounded-full bg-red-500 text-xs text-white">
                      {unreadCount > 9 ? '9+' : unreadCount}
                    </span>
                  )}
                </Link>
                <Link to="/my-bookings" className={cn('text-sm font-medium', linkCls)}>
                  Мои аренды
                </Link>
                <div className="flex items-center gap-2">
                  <Link to="/profile" className={cn('flex items-center gap-2 text-sm', transparent ? 'text-white/90 hover:text-white' : 'text-gray-700 hover:text-primary-600')}>
                    <User className="h-4 w-4" />
                    {user?.name}
                  </Link>
                  <button onClick={handleLogout} className={cn('p-1', transparent ? 'text-white/70 hover:text-red-300' : 'text-gray-400 hover:text-red-500')}>
                    <LogOut className="h-4 w-4" />
                  </button>
                </div>
              </>
            ) : (
              <>
                <Link to="/login" className={cn('text-sm', linkCls)}>
                  Войти
                </Link>
                <Link
                  to="/register"
                  className="rounded-lg bg-primary-600 px-4 py-2 text-sm text-white transition-colors hover:bg-primary-700"
                >
                  Зарегистрироваться
                </Link>
              </>
            )}
          </div>

          {/* Mobile toggle */}
          <button
            className={cn('rounded-lg border p-2 md:hidden', transparent ? 'border-white/30 text-white' : 'border-gray-200 text-gray-700')}
            onClick={() => setOpen(!open)}
            aria-label="Меню"
          >
            <MenuToggleIcon open={open} className="size-5" duration={300} />
          </button>
        </nav>
      </div>

      {/* Mobile menu */}
      <div
        className={cn(
          'fixed inset-x-0 bottom-0 top-16 z-50 flex flex-col overflow-y-auto border-t bg-white/95 backdrop-blur-lg md:hidden',
          open ? 'flex' : 'hidden',
        )}
      >
        <div className="flex h-full flex-col justify-between gap-y-4 p-5">
          <div className="grid gap-y-1">
            {navLinks.map((link) => (
              <Link
                key={link.to}
                to={link.to}
                className="rounded-lg px-3 py-3 text-base font-medium text-gray-700 hover:bg-gray-50"
                onClick={() => setOpen(false)}
              >
                {link.label}
              </Link>
            ))}
            {isAuthenticated && (
              <>
                {user?.role === 'worker' && (
                  <Link to="/worker" className="rounded-lg px-3 py-3 text-base font-medium text-accent-600 hover:bg-gray-50" onClick={() => setOpen(false)}>
                    Мои заказы
                  </Link>
                )}
                {user?.role === 'admin' && (
                  <Link to="/admin" className="rounded-lg px-3 py-3 text-base font-medium text-red-600 hover:bg-gray-50" onClick={() => setOpen(false)}>
                    Админ
                  </Link>
                )}
                <Link to="/my-bookings" className="rounded-lg px-3 py-3 text-base font-medium text-gray-700 hover:bg-gray-50" onClick={() => setOpen(false)}>
                  Мои аренды
                </Link>
                <Link to="/notifications" className="rounded-lg px-3 py-3 text-base font-medium text-gray-700 hover:bg-gray-50" onClick={() => setOpen(false)}>
                  Уведомления{unreadCount > 0 ? ` (${unreadCount})` : ''}
                </Link>
                <Link to="/profile" className="rounded-lg px-3 py-3 text-base font-medium text-gray-700 hover:bg-gray-50" onClick={() => setOpen(false)}>
                  Профиль
                </Link>
              </>
            )}
          </div>

          <div className="flex flex-col gap-2">
            {isAuthenticated ? (
              <button
                onClick={handleLogout}
                className="w-full rounded-lg border border-gray-200 px-4 py-3 text-base font-medium text-red-500 hover:bg-gray-50"
              >
                Выйти
              </button>
            ) : (
              <>
                <Link
                  to="/login"
                  className="w-full rounded-lg border border-gray-200 px-4 py-3 text-center text-base font-medium text-gray-700 hover:bg-gray-50"
                  onClick={() => setOpen(false)}
                >
                  Войти
                </Link>
                <Link
                  to="/register"
                  className="w-full rounded-lg bg-primary-600 px-4 py-3 text-center text-base font-medium text-white hover:bg-primary-700"
                  onClick={() => setOpen(false)}
                >
                  Зарегистрироваться
                </Link>
              </>
            )}
          </div>
        </div>
      </div>
    </header>
  );
}
