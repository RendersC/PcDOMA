// User & Auth
export interface User {
  id: string;
  name: string;
  email: string;
  role: 'user' | 'worker' | 'admin';
  phone?: string;
  location_id?: string;
  created_at: string;
}

export interface TokenPair {
  access_token: string;
  refresh_token: string;
  expires_in: number;
}

export interface AuthResponse {
  user: User;
  tokens: TokenPair;
}

// PC Catalog
export interface PCSpecs {
  cpu: string;
  gpu: string;
  ram: string;
  storage: string;
}

export interface Location {
  district: string;
  address: string;
  lat: number;
  lng: number;
}

export interface PC {
  id: string;
  name: string;
  slug: string;
  specs: PCSpecs;
  price_per_hour: number;
  price_per_day: number;
  currency: string;
  location: Location;
  status: 'available' | 'booked' | 'maintenance';
  images: string[];
  rating_avg: number;
  rating_count: number;
  created_at: string;
}

export interface Peripheral {
  id: string;
  name: string;
  type: 'mouse' | 'headphones' | 'headset' | 'monitor' | 'keyboard';
  description: string;
  specs: Record<string, string>;
  price_per_hour: number;
  price_per_day: number;
  currency: string;
  location_id: string;
  status: 'available' | 'booked' | 'maintenance';
  images: string[];
}

export interface Setup {
  id: string;
  name: string;
  category: 'gaming' | 'work' | 'design' | 'streaming';
  pc_id: string;
  peripheral_ids: string[];
  total_price_per_hour: number;
  total_price_per_day: number;
  discount_percent: number;
  description: string;
  images: string[];
  rating_avg: number;
  rating_count: number;
}

export interface Review {
  id: string;
  pc_id: string;
  user_id: string;
  user_name: string;
  rating: number;
  comment: string;
  created_at: string;
}

// Booking
export type BookingStatus =
  | 'pending_payment'
  | 'paid'
  | 'pending_worker'
  | 'accepted'
  | 'rejected'
  | 'delivering'
  | 'active'
  | 'completed'
  | 'cancelled';

export interface BookingPeripheral {
  id: string;
  booking_id: string;
  peripheral_id: string;
  price_per_unit: number;
}

export interface Booking {
  id: string;
  user_id: string;
  booking_type: 'custom' | 'setup';
  pc_id?: string;
  setup_id?: string;
  location_id: string;
  worker_id?: string;
  rental_type: 'hourly' | 'daily';
  start_time: string;
  end_time: string;
  base_price: number;
  extras_price: number;
  discount: number;
  total_price: number;
  currency: string;
  status: BookingStatus;
  rejection_reason?: string;
  payment_id?: string;
  peripherals?: BookingPeripheral[];
  created_at: string;
}

// Payment
export type PaymentStatus = 'pending' | 'processing' | 'completed' | 'failed' | 'refunded';

export interface Payment {
  id: string;
  booking_id: string;
  user_id: string;
  amount: number;
  currency: string;
  method: 'card' | 'kaspi';
  status: PaymentStatus;
  processed_at?: string;
  created_at: string;
}

// Notification
export interface Notification {
  id: string;
  user_id: string;
  type: string;
  title: string;
  body: string;
  data: Record<string, string>;
  is_read: boolean;
  created_at: string;
}

// API Response wrappers
export interface ListResponse<T> {
  data: T[];
  total: number;
  limit?: number;
  offset?: number;
}
