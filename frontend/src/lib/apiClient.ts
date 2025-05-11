/**
 * API client for making requests to the backend
 */

// Base URL for API endpoints
const API_BASE_URL =
  import.meta.env.VITE_API_URL || "http://localhost:8080/api";

// Stored authentication token
let authToken: string | null = null;

/**
 * Set the authentication token for API requests
 */
export const setAuthToken = (token: string) => {
  authToken = token;
  // Optionally store in localStorage for persistence
  localStorage.setItem("authToken", token);
};

/**
 * Initialize the auth token from localStorage on app startup
 */
export const initAuthToken = () => {
  const storedToken = localStorage.getItem("authToken");
  if (storedToken) {
    authToken = storedToken;
  }
};

/**
 * Clear the authentication token (for logout)
 */
export const clearAuthToken = () => {
  authToken = null;
  localStorage.removeItem("authToken");
};

/**
 * Get the current authentication token
 */
export const getAuthToken = (): string | null => {
  return authToken;
};

/**
 * Check if user is authenticated
 */
export const isAuthenticated = (): boolean => {
  return !!authToken;
};

/**
 * Generic request function with proper error handling
 */
async function request<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> {
  const url = `${API_BASE_URL}${endpoint}`;

  // Add authorization header if token exists
  const headers = new Headers(options.headers || {});
  headers.set("Content-Type", "application/json");

  if (authToken) {
    headers.set("Authorization", `Bearer ${authToken}`);
  }

  const config: RequestInit = {
    ...options,
    headers,
  };

  try {
    const response = await fetch(url, config);

    // Handle 401 Unauthorized by clearing token and redirecting to login
    if (response.status === 401) {
      clearAuthToken();
      // Optionally redirect to login
      // window.location.href = '/login';
      throw new Error("Unauthorized - Please log in again");
    }

    // For all responses, attempt to parse as JSON
    const data = await response.json();

    // Handle unsuccessful responses
    if (!response.ok) {
      throw new Error(data.message || `API error: ${response.status}`);
    }

    return data as T;
  } catch (error) {
    // Enhance error with additional information
    if (error instanceof Error) {
      throw new Error(`API request failed: ${error.message}`);
    }
    throw new Error("Unknown API error occurred");
  }
}

/**
 * API client with methods for different request types
 */
export const apiClient = {
  get: <T>(endpoint: string, options?: RequestInit) =>
    request<T>(endpoint, { method: "GET", ...options }),

  post: <T>(endpoint: string, data?: unknown, options?: RequestInit) =>
    request<T>(endpoint, {
      method: "POST",
      body: data ? JSON.stringify(data) : undefined,
      ...options,
    }),

  put: <T>(endpoint: string, data?: unknown, options?: RequestInit) =>
    request<T>(endpoint, {
      method: "PUT",
      body: data ? JSON.stringify(data) : undefined,
      ...options,
    }),

  patch: <T>(endpoint: string, data?: unknown, options?: RequestInit) =>
    request<T>(endpoint, {
      method: "PATCH",
      body: data ? JSON.stringify(data) : undefined,
      ...options,
    }),

  delete: <T>(endpoint: string, options?: RequestInit) =>
    request<T>(endpoint, { method: "DELETE", ...options }),
};
