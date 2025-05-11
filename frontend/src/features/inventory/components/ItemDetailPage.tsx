import { Link, useNavigate, useParams } from "@tanstack/react-router";
import { useDeleteItemMutation } from "../hooks/useDeleteItemMutation";
import { useItemQuery } from "../hooks/useItemQuery";

/**
 * Page to display details of a single inventory item
 */
export default function ItemDetailPage() {
  const { itemId } = useParams({ from: "/inventory/$itemId" });
  const navigate = useNavigate();

  // Fetch the item details
  const {
    data: itemResponse,
    isLoading,
    isError,
    error,
  } = useItemQuery(itemId);

  // Delete item mutation
  const deleteMutation = useDeleteItemMutation();

  // Handle item deletion
  const handleDeleteItem = async () => {
    if (window.confirm("Are you sure you want to delete this item?")) {
      await deleteMutation.mutateAsync(itemId);
      // Navigate back to inventory list after deletion
      navigate({ to: "/inventory" });
    }
  };

  // Render loading state
  if (isLoading) {
    return (
      <div className="flex justify-center items-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-gray-900"></div>
      </div>
    );
  }

  // Render error state
  if (isError) {
    return (
      <div className="bg-red-50 border border-red-200 text-red-800 px-4 py-3 rounded relative">
        <strong className="font-bold">Error:</strong>
        <span className="block sm:inline">
          {" "}
          {error?.message || "Failed to load item details"}
        </span>
      </div>
    );
  }

  // If item is not found
  if (!itemResponse?.data) {
    return (
      <div className="text-center py-10">
        <h2 className="text-2xl font-bold mb-4">Item Not Found</h2>
        <p className="mb-4">The requested item could not be found.</p>
        <Link
          to="/inventory"
          className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
        >
          Back to Inventory
        </Link>
      </div>
    );
  }

  const item = itemResponse.data;

  return (
    <div className="container mx-auto px-4 py-6">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold">{item.name}</h1>
        <div className="space-x-2">
          <Link
            to="/inventory"
            className="px-4 py-2 border border-gray-300 rounded text-gray-700 hover:bg-gray-50"
          >
            Back to List
          </Link>
          <Link
            to={`/inventory/${item.id}/edit`}
            className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
          >
            Edit Item
          </Link>
        </div>
      </div>

      <div className="bg-white shadow-lg rounded-lg overflow-hidden">
        <div className="p-6">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
              <h2 className="text-lg font-semibold mb-4">Basic Information</h2>
              <div className="space-y-3">
                <div>
                  <span className="text-gray-500">Name:</span>
                  <span className="ml-2 font-medium">{item.name}</span>
                </div>
                {item.description && (
                  <div>
                    <span className="text-gray-500">Description:</span>
                    <p className="mt-1 text-gray-800">{item.description}</p>
                  </div>
                )}
                <div>
                  <span className="text-gray-500">Type:</span>
                  <span className="ml-2">
                    {item.item_type?.name || "Unknown"}
                  </span>
                </div>
                <div>
                  <span className="text-gray-500">Barcode:</span>
                  <span className="ml-2">{item.barcode || "N/A"}</span>
                </div>
              </div>
            </div>

            <div>
              <h2 className="text-lg font-semibold mb-4">
                Quantity & Location
              </h2>
              <div className="space-y-3">
                <div>
                  <span className="text-gray-500">Quantity:</span>
                  <span className="ml-2 font-medium">
                    {item.quantity} {item.quantity_unit || ""}
                  </span>
                </div>
                <div>
                  <span className="text-gray-500">Location:</span>
                  <span className="ml-2">
                    {item.location?.name || "Not specified"}
                  </span>
                </div>
              </div>
            </div>
          </div>

          <div className="mt-8">
            <h2 className="text-lg font-semibold mb-4">
              Additional Information
            </h2>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <span className="text-gray-500">Created:</span>
                <span className="ml-2">
                  {new Date(item.created_at).toLocaleString()}
                </span>
              </div>
              <div>
                <span className="text-gray-500">Last Updated:</span>
                <span className="ml-2">
                  {new Date(item.updated_at).toLocaleString()}
                </span>
              </div>
            </div>
          </div>

          <div className="mt-10 pt-6 border-t border-gray-200">
            <button
              onClick={handleDeleteItem}
              className="px-4 py-2 bg-red-600 text-white rounded hover:bg-red-700"
            >
              Delete Item
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
