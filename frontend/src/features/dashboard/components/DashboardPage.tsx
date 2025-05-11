import { useEffect, useState } from "react";
import { apiClient } from "../../../lib/apiClient";
import type {
  Home,
  Item,
  ItemType,
  Location,
  User,
} from "../../../types/models";
import { InventoryCharts } from "../../inventory/components/charts/InventoryCharts";

export default function DashboardPage() {
  const [user, setUser] = useState<User | null>(null);
  const [homes, setHomes] = useState<Home[]>([]);
  const [selectedHomeId, setSelectedHomeId] = useState<string | null>(null);
  const [items, setItems] = useState<Item[]>([]);
  const [itemTypes, setItemTypes] = useState<ItemType[]>([]);
  const [locations, setLocations] = useState<Location[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isLoadingCharts, setIsLoadingCharts] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchUserData = async () => {
      try {
        // Fetch current user information
        const userData = await apiClient.get<User>("/auth/me");
        setUser(userData);

        // Fetch user's homes
        const homesData = await apiClient.get<Home[]>("/homes");
        setHomes(homesData);

        // Set the first home as selected if available
        if (homesData.length > 0) {
          setSelectedHomeId(homesData[0].id);
        }
      } catch (err) {
        setError(
          err instanceof Error
            ? err.message
            : "Failed to load dashboard data. Please try again."
        );
      } finally {
        setIsLoading(false);
      }
    };

    fetchUserData();
  }, []);

  // Fetch inventory data when a home is selected
  useEffect(() => {
    const fetchInventoryData = async () => {
      if (!selectedHomeId) return;

      setIsLoadingCharts(true);
      try {
        // Fetch items, item types, and locations in parallel
        const [itemsData, itemTypesData, locationsData] = await Promise.all([
          apiClient.get<Item[]>(`/inventory/items?home_id=${selectedHomeId}`),
          apiClient.get<ItemType[]>(
            `/inventory/item-types?home_id=${selectedHomeId}`
          ),
          apiClient.get<Location[]>(`/locations?home_id=${selectedHomeId}`),
        ]);

        setItems(itemsData);
        setItemTypes(itemTypesData);
        setLocations(locationsData);
      } catch (err) {
        console.error("Failed to fetch inventory data:", err);
        // Don't set error state here to keep the dashboard usable
      } finally {
        setIsLoadingCharts(false);
      }
    };

    fetchInventoryData();
  }, [selectedHomeId]);

  // Handle home selection
  const handleSelectHome = (homeId: string) => {
    setSelectedHomeId(homeId);
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-[60vh]">
        <div className="text-center">
          <div className="w-12 h-12 border-4 border-t-blue-500 border-blue-200 rounded-full animate-spin mx-auto mb-4"></div>
          <p className="text-gray-600">Loading dashboard...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="p-4 mx-auto mt-8 max-w-2xl">
        <div className="p-4 bg-red-50 border border-red-200 rounded-md">
          <h2 className="text-lg font-medium text-red-600 mb-2">Error</h2>
          <p className="text-red-600">{error}</p>
          <button
            onClick={() => window.location.reload()}
            className="mt-4 px-4 py-2 bg-red-100 hover:bg-red-200 text-red-600 rounded-md"
          >
            Retry
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-2xl font-bold text-gray-900">Dashboard</h1>
        <p className="text-gray-600">Welcome, {user?.name || "User"}!</p>
      </div>

      <div className="bg-white shadow rounded-lg overflow-hidden">
        <div className="px-6 py-5 border-b">
          <h3 className="text-lg font-medium text-gray-900">Your Homes</h3>
        </div>

        {homes.length === 0 ? (
          <div className="p-6 text-center">
            <p className="text-gray-600 mb-4">You don't have any homes yet.</p>
            <button className="px-4 py-2 text-sm font-medium text-white bg-blue-600 rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500">
              Create a Home
            </button>
          </div>
        ) : (
          <ul className="divide-y divide-gray-200">
            {homes.map((home) => (
              <li key={home.id} className="p-6 hover:bg-gray-50">
                <div className="flex items-center justify-between">
                  <div>
                    <h4 className="text-lg font-medium text-gray-900">
                      {home.name}
                    </h4>
                    <p className="text-sm text-gray-500">
                      Created on{" "}
                      {new Date(home.created_at).toLocaleDateString()}
                    </p>
                  </div>
                  <button
                    className={`px-4 py-2 text-sm font-medium rounded-md focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 ${
                      selectedHomeId === home.id
                        ? "bg-blue-100 text-blue-700 hover:bg-blue-200 border border-blue-200"
                        : "text-white bg-blue-600 hover:bg-blue-700"
                    }`}
                    onClick={() => handleSelectHome(home.id)}
                  >
                    {selectedHomeId === home.id ? "Selected" : "View"}
                  </button>
                </div>
              </li>
            ))}
          </ul>
        )}
      </div>

      {selectedHomeId && (
        <>
          <div className="mb-8">
            <h2 className="text-xl font-bold mb-6">Inventory Insights</h2>
            <InventoryCharts
              items={items}
              itemTypes={itemTypes}
              locations={locations}
              isLoading={isLoadingCharts}
            />
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div className="bg-white shadow rounded-lg p-6">
              <h3 className="text-lg font-medium text-gray-900 mb-4">
                Recent Items
              </h3>
              {items.length === 0 ? (
                <p className="text-gray-600">
                  No items added yet. Add items to your inventory to see them
                  here.
                </p>
              ) : (
                <ul className="divide-y divide-gray-200">
                  {items.slice(0, 5).map((item) => (
                    <li key={item.id} className="py-3 flex justify-between">
                      <span className="text-gray-700">{item.name}</span>
                      <span className="text-gray-500">
                        {item.quantity} {item.quantity_unit || ""}
                      </span>
                    </li>
                  ))}
                </ul>
              )}
            </div>

            <div className="bg-white shadow rounded-lg p-6">
              <h3 className="text-lg font-medium text-gray-900 mb-4">
                Quick Actions
              </h3>
              <div className="space-y-2">
                <button className="flex items-center px-4 py-2 w-full text-left text-gray-700 rounded-md hover:bg-gray-100">
                  <svg
                    className="w-5 h-5 mr-3 text-gray-500"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth="2"
                      d="M12 4v16m8-8H4"
                    />
                  </svg>
                  Add New Item
                </button>
                <button className="flex items-center px-4 py-2 w-full text-left text-gray-700 rounded-md hover:bg-gray-100">
                  <svg
                    className="w-5 h-5 mr-3 text-gray-500"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth="2"
                      d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z"
                    />
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth="2"
                      d="M15 11a3 3 0 11-6 0 3 3 0 016 0z"
                    />
                  </svg>
                  Manage Locations
                </button>
                <button className="flex items-center px-4 py-2 w-full text-left text-gray-700 rounded-md hover:bg-gray-100">
                  <svg
                    className="w-5 h-5 mr-3 text-gray-500"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth="2"
                      d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z"
                    />
                  </svg>
                  Manage API Keys
                </button>
              </div>
            </div>
          </div>
        </>
      )}

      {!selectedHomeId && homes.length > 0 && (
        <div className="text-center p-8 bg-gray-50 rounded-lg">
          <p className="text-gray-600 mb-4">
            Please select a home to view inventory insights.
          </p>
        </div>
      )}
    </div>
  );
}
