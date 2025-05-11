/**
 * Type definitions for backend models
 */

// User
export interface User {
  id: string;
  email: string;
  name: string;
  created_at: string;
  updated_at: string;
}

// Home
export interface Home {
  id: string;
  name: string;
  owner_id: string;
  created_at: string;
  updated_at: string;
}

// Home User
export interface HomeUser {
  home_id: string;
  user_id: string;
  role: string; // 'owner', 'admin', 'member'
  created_at: string;
  updated_at: string;
  // Include user information when listing home users
  user?: User;
}

// Location
export interface Location {
  id: string;
  home_id: string;
  name: string;
  description: string | null;
  parent_id: string | null;
  created_at: string;
  updated_at: string;
  // Include child locations when retrieving a location
  children?: Location[];
}

// ItemType
export interface ItemType {
  id: string;
  home_id: string;
  name: string;
  description: string | null;
  default_unit: string | null;
  created_at: string;
  updated_at: string;
}

// Item
export interface Item {
  id: string;
  home_id: string;
  type_id: string;
  name: string;
  description: string | null;
  barcode: string | null;
  location_id: string | null;
  quantity: number;
  quantity_unit: string | null;
  created_at: string;
  updated_at: string;
  // Include related entities when retrieving an item
  item_type?: ItemType;
  location?: Location;
}

// API Key
export interface ApiKey {
  id: string;
  name: string;
  key: string; // Only included when creating a new key
  user_id: string;
  expires_at: string | null;
  created_at: string;
  updated_at: string;
}

// Auth related types
export interface LoginCredentials {
  email: string;
  password: string;
}

export interface RegistrationData {
  name: string;
  email: string;
  password: string;
}

export interface AuthResponse {
  token: string;
  user: User;
}

// API Response types
export interface ApiResponse<T> {
  data: T;
  message?: string;
}

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page: number;
  per_page: number;
  total_pages: number;
}

// Error response
export interface ApiError {
  message: string;
  errors?: Record<string, string[]>;
  status?: number;
}
