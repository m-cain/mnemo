import { Outlet } from "@tanstack/react-router";

/**
 * Layout for authentication pages (login, register)
 */
export default function AuthLayout() {
  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <div className="w-full max-w-md p-8 space-y-8 bg-white rounded-lg shadow">
        <div className="text-center">
          <h1 className="text-3xl font-bold text-gray-900 mb-2">Mnemo</h1>
          <p className="text-gray-500">Home Inventory Tracking System</p>
        </div>
        <div>
          <Outlet />
        </div>
      </div>
    </div>
  );
}
