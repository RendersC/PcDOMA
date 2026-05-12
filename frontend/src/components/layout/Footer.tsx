import { Monitor } from 'lucide-react';
import { Link } from 'react-router-dom';

export function Footer() {
  return (
    <footer className="bg-gray-900 text-gray-300 mt-auto">
      <div className="max-w-7xl mx-auto px-4 py-10 sm:px-6 lg:px-8">
        <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
          <div>
            <div className="flex items-center gap-2 text-white font-bold text-lg mb-3">
              <Monitor className="w-5 h-5 text-primary-400" />
              PCDoma
            </div>
            <p className="text-sm text-gray-400">
              Аренда мощных ПК в Астане. Для игр, работы и творчества.
            </p>
          </div>
          <div>
            <h4 className="text-white font-semibold mb-3 text-sm">Навигация</h4>
            <ul className="space-y-2 text-sm">
              <li><Link to="/catalog" className="hover:text-white transition-colors">Каталог ПК</Link></li>
              <li><Link to="/setups" className="hover:text-white transition-colors">Готовые сеапы</Link></li>
              <li><Link to="/configurator" className="hover:text-white transition-colors">Конструктор</Link></li>
            </ul>
          </div>
          <div>
            <h4 className="text-white font-semibold mb-3 text-sm">Контакты</h4>
            <ul className="space-y-2 text-sm">
              <li>г. Астана, пр. Туран 24</li>
              <li>+7 (700) 123-45-67</li>
              <li>support@pcdoma.kz</li>
            </ul>
          </div>
        </div>
        <div className="border-t border-gray-800 mt-8 pt-6 text-center text-xs text-gray-500">
          © {new Date().getFullYear()} PCDoma. Все права защищены.
        </div>
      </div>
    </footer>
  );
}
