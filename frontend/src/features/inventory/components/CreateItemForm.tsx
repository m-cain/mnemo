import { Link, useNavigate } from "@tanstack/react-router";
import {
  ArrowLeft,
  CheckCircle2,
  HelpCircle,
  Package,
  PlusCircle,
} from "lucide-react";
import { useEffect, useState } from "react";
import { Button } from "../../../components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "../../../components/ui/card";
import { Input } from "../../../components/ui/input";
import { Label } from "../../../components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "../../../components/ui/select";
import { Textarea } from "../../../components/ui/textarea";
import { useCreateItemMutation } from "../hooks/useCreateItemMutation";
import { useItemTypesQuery } from "../hooks/useItemTypesQuery";

/**
 * Form to create a new inventory item
 */
export default function CreateItemForm() {
  // TODO: This should be retrieved from context or user selection
  const [selectedHomeId, setSelectedHomeId] = useState<string | null>(null);
  const navigate = useNavigate();
  const [isSubmitting, setIsSubmitting] = useState(false);

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

  // Validation state - track touched fields to show errors only after interaction
  const [errors, setErrors] = useState<Record<string, string>>({});
  const [touchedFields, setTouchedFields] = useState<Record<string, boolean>>(
    {}
  );

  // Item creation mutation
  const createMutation = useCreateItemMutation(selectedHomeId || "");

  // Fetch item types for dropdown
  const { data: itemTypesResponse, isLoading: itemTypesLoading } =
    useItemTypesQuery(selectedHomeId || "");

  // Handle form field changes - for regular Input components
  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>
  ) => {
    const { name, value } = e.target;

    // Mark field as touched
    setTouchedFields({
      ...touchedFields,
      [name]: true,
    });

    // Convert quantity to number if it's the quantity field
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

    // Validate field on change
    validateField(
      name,
      name === "quantity" ? (value === "" ? 0 : parseFloat(value)) : value
    );
  };

  // Handle select changes - for Shadcn Select component
  const handleSelectChange = (name: string, value: string) => {
    setTouchedFields({
      ...touchedFields,
      [name]: true,
    });

    setFormData({
      ...formData,
      [name]: value,
    });

    // Validate field on change
    validateField(name, value);
  };

  // Validate a single field
  const validateField = (name: string, value: string | number) => {
    const newErrors = { ...errors };

    switch (name) {
      case "name":
        if (typeof value === "string" && !value.trim()) {
          newErrors.name = "Name is required";
        } else {
          delete newErrors.name;
        }
        break;

      case "type_id":
        if (value === "") {
          newErrors.type_id = "Item type is required";
        } else {
          delete newErrors.type_id;
        }
        break;

      case "quantity":
        if (typeof value === "number" && value < 0) {
          newErrors.quantity = "Quantity cannot be negative";
        } else {
          delete newErrors.quantity;
        }
        break;

      default:
        break;
    }

    setErrors(newErrors);
  };

  // Validate the entire form
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

    // Mark all fields as touched to show all validation errors
    const allFields = [
      "name",
      "description",
      "type_id",
      "barcode",
      "quantity",
      "quantity_unit",
      "location_id",
    ];
    const allTouched = allFields.reduce((acc, field) => {
      acc[field] = true;
      return acc;
    }, {} as Record<string, boolean>);

    setTouchedFields(allTouched);

    if (!selectedHomeId) {
      setErrors({ homeId: "Please select a home" });
      return;
    }

    if (!validateForm()) {
      return;
    }

    setIsSubmitting(true);

    try {
      await createMutation.mutateAsync({
        ...formData,
        home_id: selectedHomeId,
      });
      navigate({ to: "/inventory" });
    } catch (error) {
      console.error("Failed to create item:", error);
      setErrors({ submit: "Failed to create item. Please try again." });
      setIsSubmitting(false);
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
      <Card className="mx-auto max-w-md text-center">
        <CardHeader>
          <CardTitle>No Home Selected</CardTitle>
          <CardDescription>
            Please select a home to create an inventory item.
          </CardDescription>
        </CardHeader>
        <CardFooter className="justify-center">
          <Button onClick={() => setSelectedHomeId("sample-home-id")}>
            Select Sample Home
          </Button>
        </CardFooter>
      </Card>
    );
  }

  return (
    <div className="container mx-auto px-4 py-6">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold">Create New Item</h1>
        <Button variant="outline" asChild>
          <Link to="/inventory">
            <ArrowLeft className="mr-2 h-4 w-4" />
            Back to List
          </Link>
        </Button>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Main Form Card */}
        <Card className="lg:col-span-2">
          <CardHeader>
            <CardTitle className="flex items-center">
              <Package className="mr-2 h-5 w-5" />
              Item Details
            </CardTitle>
            <CardDescription>Add a new item to your inventory</CardDescription>
          </CardHeader>
          <CardContent>
            <form
              id="create-item-form"
              onSubmit={handleSubmit}
              className="space-y-6"
            >
              {/* Basic Information Section */}
              <div className="space-y-4">
                <h3 className="text-lg font-medium">Basic Information</h3>

                {/* Item Name */}
                <div className="space-y-2">
                  <Label htmlFor="name" className="flex items-center">
                    Name <span className="text-red-500 ml-1">*</span>
                  </Label>
                  <Input
                    id="name"
                    name="name"
                    value={formData.name}
                    onChange={handleChange}
                    placeholder="Enter item name"
                    className={
                      errors.name && touchedFields.name ? "border-red-500" : ""
                    }
                  />
                  {errors.name && touchedFields.name && (
                    <p className="text-sm text-red-600">{errors.name}</p>
                  )}
                </div>

                {/* Item Description */}
                <div className="space-y-2">
                  <Label htmlFor="description">Description</Label>
                  <Textarea
                    id="description"
                    name="description"
                    value={formData.description}
                    onChange={handleChange}
                    placeholder="Optional description of the item"
                    rows={3}
                  />
                </div>

                {/* Item Type */}
                <div className="space-y-2">
                  <Label htmlFor="type_id" className="flex items-center">
                    Item Type <span className="text-red-500 ml-1">*</span>
                  </Label>
                  <Select
                    value={formData.type_id}
                    onValueChange={(value) =>
                      handleSelectChange("type_id", value)
                    }
                  >
                    <SelectTrigger
                      id="type_id"
                      className={
                        errors.type_id && touchedFields.type_id
                          ? "border-red-500"
                          : ""
                      }
                    >
                      <SelectValue placeholder="Select Item Type" />
                    </SelectTrigger>
                    <SelectContent>
                      {itemTypesLoading ? (
                        <SelectItem value="loading" disabled>
                          Loading...
                        </SelectItem>
                      ) : (
                        itemTypesResponse?.data?.map((type) => (
                          <SelectItem key={type.id} value={type.id}>
                            {type.name}
                          </SelectItem>
                        ))
                      )}
                    </SelectContent>
                  </Select>
                  {errors.type_id && touchedFields.type_id && (
                    <p className="text-sm text-red-600">{errors.type_id}</p>
                  )}
                </div>

                {/* Barcode */}
                <div className="space-y-2">
                  <Label htmlFor="barcode">Barcode</Label>
                  <Input
                    id="barcode"
                    name="barcode"
                    value={formData.barcode}
                    onChange={handleChange}
                    placeholder="Optional barcode"
                  />
                </div>
              </div>

              {/* Quantity and Location Section */}
              <div className="space-y-4 pt-4 border-t">
                <h3 className="text-lg font-medium">Quantity & Location</h3>

                {/* Quantity and Unit */}
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div className="space-y-2">
                    <Label htmlFor="quantity">Quantity</Label>
                    <Input
                      type="number"
                      id="quantity"
                      name="quantity"
                      value={formData.quantity.toString()}
                      onChange={handleChange}
                      min="0"
                      step="0.01"
                      className={
                        errors.quantity && touchedFields.quantity
                          ? "border-red-500"
                          : ""
                      }
                    />
                    {errors.quantity && touchedFields.quantity && (
                      <p className="text-sm text-red-600">{errors.quantity}</p>
                    )}
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="quantity_unit">Unit (optional)</Label>
                    <Input
                      id="quantity_unit"
                      name="quantity_unit"
                      value={formData.quantity_unit}
                      onChange={handleChange}
                      placeholder="e.g., kg, pcs, etc."
                    />
                  </div>
                </div>

                {/* Location - This would be populated from location data in a real app */}
                <div className="space-y-2">
                  <Label htmlFor="location_id">Location</Label>
                  <Select
                    value={formData.location_id}
                    onValueChange={(value) =>
                      handleSelectChange("location_id", value)
                    }
                  >
                    <SelectTrigger id="location_id">
                      <SelectValue placeholder="Select Location (Optional)" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="">No Location</SelectItem>
                      {/* This would be populated from location data */}
                      <SelectItem value="sample-location-1">
                        Living Room
                      </SelectItem>
                      <SelectItem value="sample-location-2">Kitchen</SelectItem>
                      <SelectItem value="sample-location-3">Bedroom</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
              </div>

              {/* Form error */}
              {errors.submit && (
                <div className="p-3 bg-red-50 border border-red-200 text-red-800 rounded-md">
                  {errors.submit}
                </div>
              )}
            </form>
          </CardContent>
          <CardFooter className="flex justify-between border-t pt-6">
            <Button
              variant="outline"
              onClick={() => navigate({ to: "/inventory" })}
            >
              Cancel
            </Button>
            <Button
              type="submit"
              form="create-item-form"
              disabled={isSubmitting || createMutation.isPending}
            >
              <PlusCircle className="mr-2 h-4 w-4" />
              {isSubmitting || createMutation.isPending
                ? "Creating..."
                : "Create Item"}
            </Button>
          </CardFooter>
        </Card>

        {/* Help and Tips Card */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center">
              <HelpCircle className="mr-2 h-5 w-5" />
              Tips
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="space-y-2">
              <h3 className="text-sm font-medium flex items-center">
                <CheckCircle2 className="mr-2 h-4 w-4 text-green-600" />
                Required Fields
              </h3>
              <p className="text-sm text-gray-600">
                Fields marked with an asterisk (*) are required. Make sure to
                fill them in.
              </p>
            </div>
            <div className="space-y-2">
              <h3 className="text-sm font-medium flex items-center">
                <CheckCircle2 className="mr-2 h-4 w-4 text-green-600" />
                Item Types
              </h3>
              <p className="text-sm text-gray-600">
                Categorize your items by type to make them easier to find and
                organize.
              </p>
            </div>
            <div className="space-y-2">
              <h3 className="text-sm font-medium flex items-center">
                <CheckCircle2 className="mr-2 h-4 w-4 text-green-600" />
                Units
              </h3>
              <p className="text-sm text-gray-600">
                Specify units (like kg, lb, pcs) to track your inventory more
                accurately.
              </p>
            </div>
            <div className="space-y-2">
              <h3 className="text-sm font-medium flex items-center">
                <CheckCircle2 className="mr-2 h-4 w-4 text-green-600" />
                Locations
              </h3>
              <p className="text-sm text-gray-600">
                Assigning a location helps you quickly find items in your home.
              </p>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
