import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { Toaster } from 'react-hot-toast';

import { Navbar } from './components/layout/Navbar';
import { Footer } from './components/layout/Footer';
import { ProtectedRoute } from './components/auth/ProtectedRoute';

import { HomePage } from './pages/HomePage';
import { CatalogPage } from './pages/CatalogPage';
import { PCDetailPage } from './pages/PCDetailPage';
import { SetupsPage } from './pages/SetupsPage';
import { SetupDetailPage } from './pages/SetupDetailPage';
import { ConfiguratorPage } from './pages/ConfiguratorPage';
import { BookingPage } from './pages/BookingPage';
import { PaymentPage } from './pages/PaymentPage';
import { MyBookingsPage } from './pages/MyBookingsPage';
import { ProfilePage } from './pages/ProfilePage';
import { NotificationsPage } from './pages/NotificationsPage';
import { LoginPage } from './pages/LoginPage';
import { RegisterPage } from './pages/RegisterPage';

import { WorkerDashboard } from './pages/worker/WorkerDashboard';
import { WorkerActiveOrders } from './pages/worker/WorkerActiveOrders';
import { WorkerOrderDetail } from './pages/worker/WorkerOrderDetail';

import { AdminDashboard } from './pages/admin/AdminDashboard';
import { AdminPCList } from './pages/admin/AdminPCList';
import { AdminPCForm } from './pages/admin/AdminPCForm';
import { AdminPeripheralList } from './pages/admin/AdminPeripheralList';
import { AdminPeripheralForm } from './pages/admin/AdminPeripheralForm';
import { AdminSetupList } from './pages/admin/AdminSetupList';
import { AdminSetupForm } from './pages/admin/AdminSetupForm';
import { AdminBookings } from './pages/admin/AdminBookings';

const queryClient = new QueryClient({
  defaultOptions: {
    queries: { retry: 1, staleTime: 30000 },
  },
});

function AppLayout({ children }: { children: React.ReactNode }) {
  return (
    <div className="min-h-screen flex flex-col bg-gray-50">
      <Navbar />
      <main className="flex-1 pt-16">{children}</main>
      <Footer />
    </div>
  );
}

export default function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <Toaster position="top-right" toastOptions={{ duration: 3000 }} />
        <Routes>
          {/* Public */}
          <Route path="/" element={<AppLayout><HomePage /></AppLayout>} />
          <Route path="/catalog" element={<AppLayout><CatalogPage /></AppLayout>} />
          <Route path="/catalog/:id" element={<AppLayout><PCDetailPage /></AppLayout>} />
          <Route path="/setups" element={<AppLayout><SetupsPage /></AppLayout>} />
          <Route path="/setups/:id" element={<AppLayout><SetupDetailPage /></AppLayout>} />
          <Route path="/configurator" element={<AppLayout><ConfiguratorPage /></AppLayout>} />
          <Route path="/login" element={<AppLayout><LoginPage /></AppLayout>} />
          <Route path="/register" element={<AppLayout><RegisterPage /></AppLayout>} />

          {/* Protected — user */}
          <Route path="/booking/:id" element={
            <AppLayout>
              <ProtectedRoute><BookingPage /></ProtectedRoute>
            </AppLayout>
          } />
          <Route path="/payment/:id" element={
            <AppLayout>
              <ProtectedRoute><PaymentPage /></ProtectedRoute>
            </AppLayout>
          } />
          <Route path="/my-bookings" element={
            <AppLayout>
              <ProtectedRoute><MyBookingsPage /></ProtectedRoute>
            </AppLayout>
          } />
          <Route path="/profile" element={
            <AppLayout>
              <ProtectedRoute><ProfilePage /></ProtectedRoute>
            </AppLayout>
          } />
          <Route path="/notifications" element={
            <AppLayout>
              <ProtectedRoute><NotificationsPage /></ProtectedRoute>
            </AppLayout>
          } />

          {/* Worker */}
          <Route path="/worker" element={
            <AppLayout>
              <ProtectedRoute roles={['worker', 'admin']}><WorkerDashboard /></ProtectedRoute>
            </AppLayout>
          } />
          <Route path="/worker/active" element={
            <AppLayout>
              <ProtectedRoute roles={['worker', 'admin']}><WorkerActiveOrders /></ProtectedRoute>
            </AppLayout>
          } />
          <Route path="/worker/order/:id" element={
            <AppLayout>
              <ProtectedRoute roles={['worker', 'admin']}><WorkerOrderDetail /></ProtectedRoute>
            </AppLayout>
          } />

          {/* Admin */}
          <Route path="/admin" element={
            <AppLayout>
              <ProtectedRoute roles={['admin']}><AdminDashboard /></ProtectedRoute>
            </AppLayout>
          } />
          <Route path="/admin/pcs" element={
            <AppLayout>
              <ProtectedRoute roles={['admin']}><AdminPCList /></ProtectedRoute>
            </AppLayout>
          } />
          <Route path="/admin/pcs/new" element={
            <AppLayout>
              <ProtectedRoute roles={['admin']}><AdminPCForm /></ProtectedRoute>
            </AppLayout>
          } />
          <Route path="/admin/pcs/:id/edit" element={
            <AppLayout>
              <ProtectedRoute roles={['admin']}><AdminPCForm /></ProtectedRoute>
            </AppLayout>
          } />
          <Route path="/admin/peripherals" element={
            <AppLayout>
              <ProtectedRoute roles={['admin']}><AdminPeripheralList /></ProtectedRoute>
            </AppLayout>
          } />
          <Route path="/admin/peripherals/new" element={
            <AppLayout>
              <ProtectedRoute roles={['admin']}><AdminPeripheralForm /></ProtectedRoute>
            </AppLayout>
          } />
          <Route path="/admin/setups" element={
            <AppLayout>
              <ProtectedRoute roles={['admin']}><AdminSetupList /></ProtectedRoute>
            </AppLayout>
          } />
          <Route path="/admin/setups/new" element={
            <AppLayout>
              <ProtectedRoute roles={['admin']}><AdminSetupForm /></ProtectedRoute>
            </AppLayout>
          } />
          <Route path="/admin/bookings" element={
            <AppLayout>
              <ProtectedRoute roles={['admin']}><AdminBookings /></ProtectedRoute>
            </AppLayout>
          } />

          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
      </BrowserRouter>
    </QueryClientProvider>
  );
}
