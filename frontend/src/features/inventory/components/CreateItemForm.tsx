import { zodResolver } from "@hookform/resolvers/zod";
import { Link, useNavigate } from "@tanstack/react-router";
import {
  ArrowLeft,
  CheckCircle2,
  HelpCircle,
  Package,
  PlusCircle,
} from "lucide-react";
import { useEffect, useState } from "react";
import { useForm, type SubmitHandler } from "react-hook-form";
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
import type { ItemFormValues } from "../schemas/itemValidationSchema";
import { itemSchema } from "../schemas/itemValidationSchema";

/**
 * Form to create a new inventory item
 */
export default function CreateItemForm() {
  // TODO: This should be retrieved from context or user selection
  const navigate = useNavigate();
  const [selectedHomeId, setSelectedHomeId] = useState<string | null>(null);

  // Set up react-hook-form with Zod validation
  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
    setValue,
    watch,
  } = useForm<ItemFormValues>({
    resolver: zodResolver(itemSchema) as any,
    defaultValues: {
      name: "",
      description: "",
      type_id: "",
      location_id: "",
      barcode: "",
      quantity: 1,
      quantity_unit: "",
      home_id: "",
    },
  });

  // Item creation mutation
  const createMutation = useCreateItemMutation(selectedHomeId || "");

  // Fetch item types for dropdown
  const { data: itemTypesResponse, isLoading: itemTypesLoading } =
    useItemTypesQuery(selectedHomeId || "");

  // Custom handler for Select component changes
  const handleSelectChange = (name: keyof ItemFormValues, value: string) => {
    setValue(name, value, {
      shouldValidate: true,
      shouldDirty: true,
      shouldTouch: true,
    });
  };

  // Handle form submission
  const onSubmit: SubmitHandler<ItemFormValues> = async (data) => {
    if (!selectedHomeId) {
      return;
    }

    try {
      await createMutation.mutateAsync({
        ...data,
        // Convert undefined to null for backend compatibility
        description: data.description || null,
        barcode: data.barcode || null,
        quantity_unit: data.quantity_unit || null,
        location_id: data.location_id || null,
        home_id: selectedHomeId,
      });
      navigate({ to: "/inventory" });
    } catch (error) {
      console.error("Failed to create item:", error);
    }
  };

  // For demo purposes, set a sample home ID if none is selected
  useEffect(() => {
    if (!selectedHomeId) {
      // This is just for demo purposes
      // In a real app, this would be set via context or user selection
      const sampleHomeId = "sample-home-id";
      setSelectedHomeId(sampleHomeId);
      setValue("home_id", sampleHomeId);
    }
  }, [selectedHomeId, setValue]);

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
              onSubmit={handleSubmit(onSubmit)}
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
                    placeholder="Enter item name"
                    {...register("name")}
                    className={errors.name ? "border-red-500" : ""}
                  />
                  {errors.name && (
                    <p className="text-sm text-red-600">
                      {errors.name.message}
                    </p>
                  )}
                </div>

                {/* Item Description */}
                <div className="space-y-2">
                  <Label htmlFor="description">Description</Label>
                  <Textarea
                    id="description"
                    placeholder="Optional description of the item"
                    rows={3}
                    {...register("description")}
                  />
                </div>

                {/* Item Type */}
                <div className="space-y-2">
                  <Label htmlFor="type_id" className="flex items-center">
                    Item Type <span className="text-red-500 ml-1">*</span>
                  </Label>
                  <Select
                    value={watch("type_id")}
                    onValueChange={(value) =>
                      handleSelectChange("type_id", value)
                    }
                  >
                    <SelectTrigger
                      id="type_id"
                      className={errors.type_id ? "border-red-500" : ""}
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
                  {errors.type_id && (
                    <p className="text-sm text-red-600">
                      {errors.type_id.message}
                    </p>
                  )}
                </div>

                {/* Barcode */}
                <div className="space-y-2">
                  <Label htmlFor="barcode">Barcode</Label>
                  <Input
                    id="barcode"
                    placeholder="Optional barcode"
                    {...register("barcode")}
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
                      min="0"
                      step="0.01"
                      {...register("quantity", { valueAsNumber: true })}
                      className={errors.quantity ? "border-red-500" : ""}
                    />
                    {errors.quantity && (
                      <p className="text-sm text-red-600">
                        {errors.quantity.message}
                      </p>
                    )}
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="quantity_unit">Unit (optional)</Label>
                    <Input
                      id="quantity_unit"
                      placeholder="e.g., kg, pcs, etc."
                      {...register("quantity_unit")}
                    />
                  </div>
                </div>

                {/* Location - This would be populated from location data in a real app */}
                <div className="space-y-2">
                  <Label htmlFor="location_id">Location</Label>
                  <Select
                    value={watch("location_id") || ""}
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
              {createMutation.isError && (
                <div className="p-3 bg-red-50 border border-red-200 text-red-800 rounded-md">
                  Failed to create item. Please try again.
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
