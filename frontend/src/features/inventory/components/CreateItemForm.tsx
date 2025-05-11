import { useNavigate } from "@tanstack/react-router";
import { useEffect, useState } from "react";
import { useCreateItemMutation } from "../hooks/useCreateItemMutation";
import { useItemTypesQuery } from "../hooks/useItemTypesQuery";

/**
 * Form to create a new inventory item
 */
export default function CreateItemForm() {
  // TODO: This should be retrieved from context or user selection
  const [selectedHomeId, setSelectedHomeId] = useState<string | null>(null);
  const navigate = useNavigate();

  // Form state
  const [formData, setFormData] = useState({
    name: "",
    description: "",
    type_id: "",
    location_id: "",
    barcode: "",
    quantity: 1,
    quantity_unit: "",
  });

  // Validation state
  const [errors, setErrors] = useState<Record<string, string>>({});

  // Item creation mutation
  const createMutation = useCreateItemMutation(selectedHomeId || "");

  // Fetch item types for dropdown
  const { data: itemTypesResponse, isLoading: itemTypesLoading } =
    useItemTypesQuery(selectedHomeId || "");

  // Handle form field changes
  const handleChange = (
    e: React.ChangeEvent<
      HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement
    >
  ) => {
    const { name, value } = e.target;

    // Convert quantity to number
    if (name === "quantity") {
      setFormData({
        ...formData,
        [name]: value === "" ? 0 : parseFloat(value),
      });
    } else {
      setFormData({
        ...formData,
        [name]: value,
      });
    }
  };

  // Validate form
  const validateForm = () => {
    const newErrors: Record<string, string> = {};

    if (!formData.name.trim()) {
      newErrors.name = "Name is required";
    }

    if (!formData.type_id) {
      newErrors.type_id = "Item type is required";
    }

    if (formData.quantity < 0) {
      newErrors.quantity = "Quantity cannot be negative";
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  // Handle form submission
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!selectedHomeId) {
      setErrors({ homeId: "Please select a home" });
      return;
    }

    if (!validateForm()) {
      return;
    }

    try {
      await createMutation.mutateAsync({
        ...formData,
        home_id: selectedHomeId,
      });
      navigate({ to: "/inventory" });
    } catch (error) {
      console.error("Failed to create item:", error);
      setErrors({ submit: "Failed to create item. Please try again." });
    }
  };

  // For demo purposes, set a sample home ID if none is selected
  useEffect(() => {
    if (!selectedHomeId) {
      // This is just for demo purposes
      // In a real app, this would be set via context or user selection
      setSelectedHomeId("sample-home-id");
    }
  }, [selectedHomeId]);

  // Render no home selected state
  if (!selectedHomeId) {
    return (
      <div className="text-center py-10">
        <h2 className="text-2xl font-bold mb-4">No Home Selected</h2>
        <p className="mb-4">
          Please select a home to create an inventory item.
        </p>
        {/* This would be connected to home selection UI in the future */}
        <button
          className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
          onClick={() => setSelectedHomeId("sample-home-id")} // Temporary for demo
        >
          Select Sample Home
        </button>
      </div>
    );
  }

  return (
    <div className="container mx-auto px-4 py-6">
      <div className="mb-6">
        <h1 className="text-2xl font-bold">Create New Item</h1>
        <p className="text-gray-600">Add a new item to your inventory</p>
      </div>

      <div className="bg-white rounded-lg shadow-md p-6">
        <form onSubmit={handleSubmit}>
          {/* Item Name */}
          <div className="mb-4">
            <label
              htmlFor="name"
              className="block text-sm font-medium text-gray-700 mb-1"
            >
              Name *
            </label>
            <input
              type="text"
              id="name"
              name="name"
              value={formData.name}
              onChange={handleChange}
              className={`w-full p-2 border rounded ${
                errors.name ? "border-red-500" : "border-gray-300"
              }`}
              placeholder="Item name"
            />
            {errors.name && (
              <p className="mt-1 text-sm text-red-600">{errors.name}</p>
            )}
          </div>

          {/* Item Description */}
          <div className="mb-4">
            <label
              htmlFor="description"
              className="block text-sm font-medium text-gray-700 mb-1"
            >
              Description
            </label>
            <textarea
              id="description"
              name="description"
              value={formData.description}
              onChange={handleChange}
              rows={3}
              className="w-full p-2 border border-gray-300 rounded"
              placeholder="Optional description"
            />
          </div>

          {/* Item Type */}
          <div className="mb-4">
            <label
              htmlFor="type_id"
              className="block text-sm font-medium text-gray-700 mb-1"
            >
              Item Type *
            </label>
            <select
              id="type_id"
              name="type_id"
              value={formData.type_id}
              onChange={handleChange}
              className={`w-full p-2 border rounded ${
                errors.type_id ? "border-red-500" : "border-gray-300"
              }`}
            >
              <option value="">Select Item Type</option>
              {itemTypesLoading ? (
                <option disabled>Loading...</option>
              ) : (
                itemTypesResponse?.data?.map((type) => (
                  <option key={type.id} value={type.id}>
                    {type.name}
                  </option>
                ))
              )}
            </select>
            {errors.type_id && (
              <p className="mt-1 text-sm text-red-600">{errors.type_id}</p>
            )}
          </div>

          {/* Barcode */}
          <div className="mb-4">
            <label
              htmlFor="barcode"
              className="block text-sm font-medium text-gray-700 mb-1"
            >
              Barcode
            </label>
            <input
              type="text"
              id="barcode"
              name="barcode"
              value={formData.barcode}
              onChange={handleChange}
              className="w-full p-2 border border-gray-300 rounded"
              placeholder="Optional barcode"
            />
          </div>

          {/* Quantity and Unit */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
            <div>
              <label
                htmlFor="quantity"
                className="block text-sm font-medium text-gray-700 mb-1"
              >
                Quantity
              </label>
              <input
                type="number"
                id="quantity"
                name="quantity"
                value={formData.quantity.toString()}
                onChange={handleChange}
                min="0"
                step="0.01"
                className={`w-full p-2 border rounded ${
                  errors.quantity ? "border-red-500" : "border-gray-300"
                }`}
              />
              {errors.quantity && (
                <p className="mt-1 text-sm text-red-600">{errors.quantity}</p>
              )}
            </div>
            <div>
              <label
                htmlFor="quantity_unit"
                className="block text-sm font-medium text-gray-700 mb-1"
              >
                Unit (optional)
              </label>
              <input
                type="text"
                id="quantity_unit"
                name="quantity_unit"
                value={formData.quantity_unit}
                onChange={handleChange}
                className="w-full p-2 border border-gray-300 rounded"
                placeholder="e.g., kg, pcs, etc."
              />
            </div>
          </div>

          {/* Location - This would be populated from location data in a real app */}
          <div className="mb-6">
            <label
              htmlFor="location_id"
              className="block text-sm font-medium text-gray-700 mb-1"
            >
              Location
            </label>
            <select
              id="location_id"
              name="location_id"
              value={formData.location_id}
              onChange={handleChange}
              className="w-full p-2 border border-gray-300 rounded"
            >
              <option value="">No Location</option>
              {/* This would be populated from location data */}
              <option value="sample-location-1">Living Room</option>
              <option value="sample-location-2">Kitchen</option>
              <option value="sample-location-3">Bedroom</option>
            </select>
          </div>

          {/* Form error */}
          {errors.submit && (
            <div className="mb-4 p-2 bg-red-50 border border-red-200 text-red-800 rounded">
              {errors.submit}
            </div>
          )}

          {/* Submit and Cancel buttons */}
          <div className="flex justify-end space-x-2">
            <button
              type="button"
              onClick={() => navigate({ to: "/inventory" })}
              className="px-4 py-2 border border-gray-300 rounded text-gray-700 hover:bg-gray-50"
            >
              Cancel
            </button>
            <button
              type="submit"
              className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
              disabled={createMutation.isPending}
            >
              {createMutation.isPending ? "Creating..." : "Create Item"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
