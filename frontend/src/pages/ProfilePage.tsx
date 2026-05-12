import { useState } from 'react';
import { User, Lock, Save } from 'lucide-react';
import { authApi } from '../api/auth';
import { useAuthStore } from '../store/authStore';
import toast from 'react-hot-toast';

export function ProfilePage() {
  const { user, setUser } = useAuthStore();

  const [name, setName] = useState(user?.name || '');
  const [phone, setPhone] = useState(user?.phone || '');
  const [savingProfile, setSavingProfile] = useState(false);

  const [oldPassword, setOldPassword] = useState('');
  const [newPassword, setNewPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [savingPassword, setSavingPassword] = useState(false);

  const handleSaveProfile = async () => {
    if (!user) return;
    setSavingProfile(true);
    try {
      const res = await authApi.updateProfile(user.id, { name, phone });
      setUser(res.data);
      toast.success('Профиль обновлён');
    } catch {
      toast.error('Не удалось обновить профиль');
    } finally {
      setSavingProfile(false);
    }
  };

  const handleChangePassword = async () => {
    if (!user) return;
    if (newPassword !== confirmPassword) {
      toast.error('Пароли не совпадают');
      return;
    }
    if (newPassword.length < 8) {
      toast.error('Минимум 8 символов');
      return;
    }
    setSavingPassword(true);
    try {
      await authApi.changePassword(user.id, { old_password: oldPassword, new_password: newPassword });
      toast.success('Пароль изменён');
      setOldPassword('');
      setNewPassword('');
      setConfirmPassword('');
    } catch {
      toast.error('Неверный текущий пароль');
    } finally {
      setSavingPassword(false);
    }
  };

  return (
    <div className="max-w-2xl mx-auto px-4 py-8 sm:px-6">
      <h1 className="text-2xl font-bold text-gray-900 mb-6">Профиль</h1>

      {/* Profile info */}
      <div className="bg-white rounded-2xl border border-gray-100 p-6 mb-6">
        <h2 className="font-bold text-gray-900 mb-4 flex items-center gap-2">
          <User className="w-5 h-5 text-primary-600" /> Личные данные
        </h2>

        <div className="mb-4">
          <label className="block text-sm font-medium text-gray-700 mb-1">Email</label>
          <input
            type="email"
            value={user?.email || ''}
            disabled
            className="w-full px-4 py-3 border border-gray-100 bg-gray-50 rounded-xl text-sm text-gray-500 cursor-not-allowed"
          />
        </div>

        <div className="mb-4">
          <label className="block text-sm font-medium text-gray-700 mb-1">Имя</label>
          <input
            type="text"
            value={name}
            onChange={(e) => setName(e.target.value)}
            className="w-full px-4 py-3 border border-gray-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-primary-500"
          />
        </div>

        <div className="mb-5">
          <label className="block text-sm font-medium text-gray-700 mb-1">Телефон</label>
          <input
            type="tel"
            value={phone}
            onChange={(e) => setPhone(e.target.value)}
            placeholder="+7 (___) ___-__-__"
            className="w-full px-4 py-3 border border-gray-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-primary-500"
          />
        </div>

        <div className="mb-4 flex items-center gap-2">
          <span className="text-xs text-gray-400">Роль:</span>
          <span className={`text-xs font-medium px-2 py-0.5 rounded-full ${
            user?.role === 'admin' ? 'bg-red-100 text-red-700' :
            user?.role === 'worker' ? 'bg-blue-100 text-blue-700' :
            'bg-green-100 text-green-700'
          }`}>
            {user?.role === 'admin' ? 'Администратор' : user?.role === 'worker' ? 'Работник' : 'Клиент'}
          </span>
        </div>

        <button
          onClick={handleSaveProfile}
          disabled={savingProfile}
          className="flex items-center gap-2 bg-primary-600 hover:bg-primary-700 disabled:opacity-50 text-white px-5 py-2.5 rounded-xl text-sm font-medium transition-colors"
        >
          <Save className="w-4 h-4" />
          {savingProfile ? 'Сохраняем...' : 'Сохранить'}
        </button>
      </div>

      {/* Change password */}
      <div className="bg-white rounded-2xl border border-gray-100 p-6">
        <h2 className="font-bold text-gray-900 mb-4 flex items-center gap-2">
          <Lock className="w-5 h-5 text-primary-600" /> Изменить пароль
        </h2>

        <div className="space-y-4 mb-5">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Текущий пароль</label>
            <input
              type="password"
              value={oldPassword}
              onChange={(e) => setOldPassword(e.target.value)}
              className="w-full px-4 py-3 border border-gray-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-primary-500"
              placeholder="••••••••"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Новый пароль</label>
            <input
              type="password"
              value={newPassword}
              onChange={(e) => setNewPassword(e.target.value)}
              className="w-full px-4 py-3 border border-gray-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-primary-500"
              placeholder="Минимум 8 символов"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Повторите пароль</label>
            <input
              type="password"
              value={confirmPassword}
              onChange={(e) => setConfirmPassword(e.target.value)}
              className="w-full px-4 py-3 border border-gray-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-primary-500"
              placeholder="••••••••"
            />
          </div>
        </div>

        <button
          onClick={handleChangePassword}
          disabled={savingPassword || !oldPassword || !newPassword}
          className="flex items-center gap-2 bg-gray-800 hover:bg-gray-900 disabled:opacity-40 text-white px-5 py-2.5 rounded-xl text-sm font-medium transition-colors"
        >
          <Lock className="w-4 h-4" />
          {savingPassword ? 'Меняем...' : 'Изменить пароль'}
        </button>
      </div>
    </div>
  );
}
