import { Link, useNavigate, useParams } from "@tanstack/react-router";
import {
  AlertTriangle,
  ArrowLeft,
  MapPin,
  Package,
  Pencil,
  Trash,
} from "lucide-react";
import { Badge } from "../../../components/ui/badge";
import { Button } from "../../../components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "../../../components/ui/card";
import { useDeleteItemMutation } from "../hooks/useDeleteItemMutation";
import { useItemQuery } from "../hooks/useItemQuery";
import { useUpdateItemQuantityMutation } from "../hooks/useUpdateItemQuantityMutation";
import type { QuantityAdjustmentFormValues } from "../schemas/itemValidationSchema";
import { QuantityAdjustmentDialog } from "./QuantityAdjustmentDialog.zod";

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

  // Quantity adjustment mutation
  const quantityMutation = useUpdateItemQuantityMutation(itemId);

  // Handle item deletion
  const handleDeleteItem = async () => {
    if (window.confirm("Are you sure you want to delete this item?")) {
      await deleteMutation.mutateAsync(itemId);
      // Navigate back to inventory list after deletion
      navigate({ to: "/inventory" });
    }
  };

  // Handle quantity adjustment from the dialog
  const handleQuantityAdjustment = async (
    data: QuantityAdjustmentFormValues
  ) => {
    await quantityMutation.mutateAsync(data.newQuantity);
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
      <Card className="mx-auto max-w-4xl">
        <CardHeader>
          <CardTitle className="flex items-center gap-2 text-red-600">
            <AlertTriangle className="h-5 w-5" />
            Error Loading Item
          </CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-red-600">
            {error?.message || "Failed to load item details"}
          </p>
        </CardContent>
        <CardFooter>
          <Button variant="outline" asChild>
            <Link to="/inventory">
              <ArrowLeft className="mr-2 h-4 w-4" />
              Back to Inventory
            </Link>
          </Button>
        </CardFooter>
      </Card>
    );
  }

  // If item is not found
  if (!itemResponse?.data) {
    return (
      <Card className="mx-auto max-w-4xl text-center">
        <CardHeader>
          <CardTitle>Item Not Found</CardTitle>
          <CardDescription>
            The requested item could not be found.
          </CardDescription>
        </CardHeader>
        <CardFooter className="justify-center">
          <Button asChild>
            <Link to="/inventory">
              <ArrowLeft className="mr-2 h-4 w-4" />
              Back to Inventory
            </Link>
          </Button>
        </CardFooter>
      </Card>
    );
  }

  const item = itemResponse.data;
  const isLowStock = item.quantity <= 5;

  return (
    <div className="container mx-auto px-4 py-6">
      {/* Header with actions */}
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold">{item.name}</h1>
        <div className="flex space-x-2">
          <Button variant="outline" asChild>
            <Link to="/inventory">
              <ArrowLeft className="mr-2 h-4 w-4" />
              Back to List
            </Link>
          </Button>
          <Button asChild>
            <Link to={`/inventory/${item.id}/edit`}>
              <Pencil className="mr-2 h-4 w-4" />
              Edit Item
            </Link>
          </Button>
        </div>
      </div>

      {/* Main content */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Basic Info Card */}
        <Card className="lg:col-span-2">
          <CardHeader>
            <CardTitle className="flex items-center">
              <Package className="mr-2 h-5 w-5" />
              Item Details
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div className="space-y-4">
                <div>
                  <h3 className="text-sm font-medium text-gray-500">Name</h3>
                  <p className="text-base font-semibold">{item.name}</p>
                </div>
                {item.description && (
                  <div>
                    <h3 className="text-sm font-medium text-gray-500">
                      Description
                    </h3>
                    <p className="text-base">{item.description}</p>
                  </div>
                )}
                <div>
                  <h3 className="text-sm font-medium text-gray-500">Type</h3>
                  <p className="text-base">
                    {item.item_type?.name || "Unknown"}
                  </p>
                </div>
                <div>
                  <h3 className="text-sm font-medium text-gray-500">Barcode</h3>
                  <p className="text-base">{item.barcode || "N/A"}</p>
                </div>
              </div>
              <div className="space-y-4">
                <div>
                  <h3 className="text-sm font-medium text-gray-500">Created</h3>
                  <p className="text-base">
                    {new Date(item.created_at).toLocaleDateString("en-US", {
                      year: "numeric",
                      month: "short",
                      day: "numeric",
                      hour: "2-digit",
                      minute: "2-digit",
                    })}
                  </p>
                </div>
                <div>
                  <h3 className="text-sm font-medium text-gray-500">
                    Last Updated
                  </h3>
                  <p className="text-base">
                    {new Date(item.updated_at).toLocaleDateString("en-US", {
                      year: "numeric",
                      month: "short",
                      day: "numeric",
                      hour: "2-digit",
                      minute: "2-digit",
                    })}
                  </p>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Quantity & Location Card */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center">
              <MapPin className="mr-2 h-5 w-5" />
              Quantity & Location
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-6">
              <div>
                <h3 className="text-sm font-medium text-gray-500">Quantity</h3>
                <div className="flex items-center mt-1">
                  <p className="text-xl font-semibold">
                    {item.quantity}
                    {item.quantity_unit && (
                      <span className="ml-1 text-base font-normal">
                        {item.quantity_unit}
                      </span>
                    )}
                  </p>
                  {isLowStock && (
                    <Badge className="ml-2 bg-red-100 text-red-800">
                      Low Stock
                    </Badge>
                  )}
                </div>
                <div className="mt-4">
                  <QuantityAdjustmentDialog
                    currentQuantity={item.quantity}
                    onAdjust={handleQuantityAdjustment}
                  />
                </div>
              </div>
              <div>
                <h3 className="text-sm font-medium text-gray-500">Location</h3>
                <p className="text-base font-medium mt-1">
                  {item.location?.name || "Not specified"}
                </p>
              </div>
            </div>
          </CardContent>
          <CardFooter className="border-t pt-6">
            <Button
              variant="destructive"
              onClick={handleDeleteItem}
              className="w-full"
            >
              <Trash className="mr-2 h-4 w-4" />
              Delete Item
            </Button>
          </CardFooter>
        </Card>
      </div>
    </div>
  );
}
