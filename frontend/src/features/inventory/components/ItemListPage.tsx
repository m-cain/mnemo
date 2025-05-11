import { Link, useNavigate } from "@tanstack/react-router";
import {
  ArrowDown,
  ArrowUp,
  CircleX,
  Filter,
  MapPin,
  MoveDownLeft,
  MoveUpRight,
  Package,
  Search,
} from "lucide-react";
import { useMemo, useRef, useState } from "react";
import { Badge } from "../../../components/ui/badge";
import { Button } from "../../../components/ui/button";
import { Card } from "../../../components/ui/card";
import { Checkbox } from "../../../components/ui/checkbox";
import { Input } from "../../../components/ui/input";
import { Label } from "../../../components/ui/label";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "../../../components/ui/popover";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "../../../components/ui/table";
import { cn } from "../../../lib/utils";
import { useDeleteItemMutation } from "../hooks/useDeleteItemMutation";
import { useItemsQuery } from "../hooks/useItemsQuery";
import { useItemTypesQuery } from "../hooks/useItemTypesQuery";

/**
 * Page to display a list of inventory items
 */
export default function ItemListPage() {
  // TODO: This should be retrieved from context or user selection
  const [selectedHomeId, setSelectedHomeId] = useState<string | null>(null);
  const navigate = useNavigate();
  const [searchTerm, setSearchTerm] = useState("");
  const [selectedTypes, setSelectedTypes] = useState<string[]>([]);
  const inputRef = useRef<HTMLInputElement>(null);

  // Sorting state
  const [sortField, setSortField] = useState<string>("name");
  const [sortOrder, setSortOrder] = useState<"asc" | "desc">("asc");

  // Fetch items for the selected home
  const {
    data: itemsResponse,
    isLoading,
    isError,
    error,
  } = useItemsQuery(selectedHomeId || "");

  // Fetch item types for filtering
  const { data: itemTypesResponse } = useItemTypesQuery(selectedHomeId || "");

  // Delete item mutation
  const deleteMutation = useDeleteItemMutation(selectedHomeId || undefined);

  // Get unique item types
  const itemTypes = useMemo(() => {
    if (!itemTypesResponse?.data) return [];
    return itemTypesResponse.data;
  }, [itemTypesResponse]);

  // Calculate stats
  const stats = useMemo(() => {
    if (!itemsResponse?.data) {
      return {
        totalItems: 0,
        totalQuantity: 0,
        itemTypes: 0,
        locations: 0,
        lowStockItems: 0,
      };
    }

    const items = itemsResponse.data;
    const totalQuantity = items.reduce((sum, item) => sum + item.quantity, 0);
    const uniqueTypeIds = new Set(items.map((item) => item.type_id));
    const uniqueLocationIds = new Set(
      items.map((item) => item.location_id).filter(Boolean)
    );
    const lowStockItems = items.filter((item) => item.quantity <= 5).length;

    return {
      totalItems: items.length,
      totalQuantity,
      itemTypes: uniqueTypeIds.size,
      locations: uniqueLocationIds.size,
      lowStockItems,
    };
  }, [itemsResponse]);

  // Filter data based on search term and selected types
  const filteredItems = useMemo(() => {
    if (!itemsResponse?.data) return [];

    return itemsResponse.data.filter((item) => {
      const matchesSearch =
        item.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
        (item.description &&
          item.description.toLowerCase().includes(searchTerm.toLowerCase()));

      const matchesType =
        selectedTypes.length === 0 ||
        (item.type_id && selectedTypes.includes(item.type_id));

      return matchesSearch && matchesType;
    });
  }, [itemsResponse, searchTerm, selectedTypes]);

  // Sort the filtered items
  const sortedItems = useMemo(() => {
    if (!filteredItems.length) return [];

    return [...filteredItems].sort((a, b) => {
      let aValue: string | number | null;
      let bValue: string | number | null;

      // Get the values based on sort field
      switch (sortField) {
        case "name":
          aValue = a.name.toLowerCase();
          bValue = b.name.toLowerCase();
          break;
        case "type":
          aValue =
            itemTypes.find((t) => t.id === a.type_id)?.name?.toLowerCase() ||
            "";
          bValue =
            itemTypes.find((t) => t.id === b.type_id)?.name?.toLowerCase() ||
            "";
          break;
        case "location":
          aValue = a.location?.name?.toLowerCase() || "";
          bValue = b.location?.name?.toLowerCase() || "";
          break;
        case "quantity":
          aValue = a.quantity;
          bValue = b.quantity;
          break;
        default:
          aValue = a.name.toLowerCase();
          bValue = b.name.toLowerCase();
      }

      // Handle the sort order
      if (sortOrder === "asc") {
        return aValue > bValue ? 1 : aValue < bValue ? -1 : 0;
      } else {
        return aValue < bValue ? 1 : aValue > bValue ? -1 : 0;
      }
    });
  }, [filteredItems, sortField, sortOrder, itemTypes]);

  // Handle item deletion
  const handleDeleteItem = async (itemId: string) => {
    if (window.confirm("Are you sure you want to delete this item?")) {
      await deleteMutation.mutateAsync(itemId);
    }
  };

  // Handle column sort
  const handleSort = (field: string) => {
    if (sortField === field) {
      // If the same field is clicked, toggle the sort order
      setSortOrder(sortOrder === "asc" ? "desc" : "asc");
    } else {
      // If a different field is clicked, set it as the new sort field and default to ascending
      setSortField(field);
      setSortOrder("asc");
    }
  };

  // Helper function to render sort indicator
  const renderSortIndicator = (field: string) => {
    if (sortField !== field) return null;

    return sortOrder === "asc" ? (
      <ArrowUp className="inline h-4 w-4 ml-1" />
    ) : (
      <ArrowDown className="inline h-4 w-4 ml-1" />
    );
  };

  // Handle type selection
  const handleTypeChange = (checked: boolean, typeId: string) => {
    if (checked) {
      setSelectedTypes((prev) => [...prev, typeId]);
    } else {
      setSelectedTypes((prev) => prev.filter((id) => id !== typeId));
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
          {error?.message || "Failed to load items"}
        </span>
      </div>
    );
  }

  // Render no home selected state
  if (!selectedHomeId) {
    return (
      <div className="text-center py-10">
        <h2 className="text-2xl font-bold mb-4">No Home Selected</h2>
        <p className="mb-4">Please select a home to view inventory items.</p>
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

  // Render the inventory dashboard
  return (
    <div className="container mx-auto px-4 py-6">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold">Inventory Management</h1>
        <Link
          to="/inventory/new"
          className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
        >
          Add New Item
        </Link>
      </div>

      {/* Stats Section */}
      <div className="grid text-left grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 w-full gap-4 lg:gap-8 mb-6">
        <div className="flex gap-0 flex-col justify-between p-6 border rounded-md">
          <MoveUpRight className="w-4 h-4 mb-10 text-blue-500" />
          <h2 className="text-4xl tracking-tighter max-w-xl text-left font-regular flex flex-row gap-4 items-end">
            {stats.totalQuantity}
          </h2>
          <p className="text-base leading-relaxed tracking-tight text-gray-500 max-w-xl text-left">
            Total Items Quantity
          </p>
        </div>
        <div className="flex gap-0 flex-col justify-between p-6 border rounded-md">
          <Package className="w-4 h-4 mb-10 text-blue-500" />
          <h2 className="text-4xl tracking-tighter max-w-xl text-left font-regular flex flex-row gap-4 items-end">
            {stats.itemTypes}
          </h2>
          <p className="text-base leading-relaxed tracking-tight text-gray-500 max-w-xl text-left">
            Item Types
          </p>
        </div>
        <div className="flex gap-0 flex-col justify-between p-6 border rounded-md">
          <MapPin className="w-4 h-4 mb-10 text-green-500" />
          <h2 className="text-4xl tracking-tighter max-w-xl text-left font-regular flex flex-row gap-4 items-end">
            {stats.locations}
          </h2>
          <p className="text-base leading-relaxed tracking-tight text-gray-500 max-w-xl text-left">
            Locations
          </p>
        </div>
        <div className="flex gap-0 flex-col justify-between p-6 border rounded-md">
          <MoveDownLeft className="w-4 h-4 mb-10 text-red-500" />
          <h2 className="text-4xl tracking-tighter max-w-xl text-left font-regular flex flex-row gap-4 items-end">
            {stats.lowStockItems}
          </h2>
          <p className="text-base leading-relaxed tracking-tight text-gray-500 max-w-xl text-left">
            Low Stock Items
          </p>
        </div>
      </div>

      {/* Search and Filters */}
      <Card className="p-6 mb-6">
        <div className="flex flex-wrap items-center justify-between gap-4">
          <div className="flex items-center gap-4 flex-grow">
            <div className="relative flex-grow max-w-md">
              <Input
                ref={inputRef}
                className={cn("ps-9", Boolean(searchTerm) && "pe-9")}
                value={searchTerm}
                onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                  setSearchTerm(e.target.value)
                }
                placeholder="Search inventory..."
                type="text"
              />
              <div className="pointer-events-none absolute inset-y-0 start-0 flex items-center justify-center ps-3 text-muted-foreground/80">
                <Search size={16} strokeWidth={2} aria-hidden="true" />
              </div>
              {Boolean(searchTerm) && (
                <button
                  className="absolute inset-y-0 end-0 flex h-full w-9 items-center justify-center rounded-e-lg text-muted-foreground/80 hover:text-foreground"
                  aria-label="Clear search"
                  onClick={() => {
                    setSearchTerm("");
                    if (inputRef.current) {
                      inputRef.current.focus();
                    }
                  }}
                >
                  <CircleX size={16} strokeWidth={2} aria-hidden="true" />
                </button>
              )}
            </div>

            {/* Type Filter */}
            {itemTypes.length > 0 && (
              <Popover>
                <PopoverTrigger asChild>
                  <Button variant="outline">
                    <Filter
                      className="-ms-1 me-2 opacity-60"
                      size={16}
                      strokeWidth={2}
                      aria-hidden="true"
                    />
                    Item Type
                    {selectedTypes.length > 0 && (
                      <span className="-me-1 ms-3 inline-flex h-5 max-h-full items-center rounded border border-gray-200 bg-white px-1 font-[inherit] text-[0.625rem] font-medium text-gray-500">
                        {selectedTypes.length}
                      </span>
                    )}
                  </Button>
                </PopoverTrigger>
                <PopoverContent className="min-w-36 p-3" align="start">
                  <div className="space-y-3">
                    <div className="text-xs font-medium text-gray-500">
                      Filter by Type
                    </div>
                    <div className="space-y-3">
                      {itemTypes.map((type) => (
                        <div key={type.id} className="flex items-center gap-2">
                          <Checkbox
                            id={`type-${type.id}`}
                            checked={selectedTypes.includes(type.id)}
                            onCheckedChange={(checked: boolean) =>
                              handleTypeChange(checked, type.id)
                            }
                          />
                          <Label
                            htmlFor={`type-${type.id}`}
                            className="flex grow justify-between gap-2 font-normal"
                          >
                            {type.name}{" "}
                            <span className="ms-2 text-xs text-gray-500">
                              {itemsResponse?.data?.filter(
                                (item) => item.type_id === type.id
                              ).length || 0}
                            </span>
                          </Label>
                        </div>
                      ))}
                    </div>
                  </div>
                </PopoverContent>
              </Popover>
            )}
          </div>
        </div>
      </Card>

      {/* Render empty items state */}
      {(!itemsResponse?.data || filteredItems.length === 0) && (
        <div className="text-center py-10">
          <h2 className="text-2xl font-bold mb-4">No Items Found</h2>
          <p className="mb-4">
            {filteredItems.length === 0 && itemsResponse?.data?.length
              ? "No items match your search criteria."
              : "Your inventory is empty. Start by adding some items."}
          </p>
          {filteredItems.length === 0 && itemsResponse?.data?.length ? (
            <Button
              onClick={() => {
                setSearchTerm("");
                setSelectedTypes([]);
              }}
            >
              Reset Filters
            </Button>
          ) : (
            <Link
              to="/inventory/new"
              className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
            >
              Add First Item
            </Link>
          )}
        </div>
      )}

      {/* Inventory Table */}
      {filteredItems.length > 0 && (
        <div className="overflow-hidden rounded-lg border border-gray-200 bg-white">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead
                  onClick={() => handleSort("name")}
                  className="cursor-pointer hover:bg-gray-50"
                >
                  Name {renderSortIndicator("name")}
                </TableHead>
                <TableHead
                  onClick={() => handleSort("type")}
                  className="cursor-pointer hover:bg-gray-50"
                >
                  Type {renderSortIndicator("type")}
                </TableHead>
                <TableHead
                  onClick={() => handleSort("location")}
                  className="cursor-pointer hover:bg-gray-50"
                >
                  Location {renderSortIndicator("location")}
                </TableHead>
                <TableHead
                  onClick={() => handleSort("quantity")}
                  className="cursor-pointer hover:bg-gray-50 text-right"
                >
                  Quantity {renderSortIndicator("quantity")}
                </TableHead>
                <TableHead className="text-right">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {sortedItems.map((item) => (
                <TableRow key={item.id}>
                  <TableCell>
                    <div className="font-medium">{item.name}</div>
                    {item.description && (
                      <div className="text-sm text-gray-500">
                        {item.description}
                      </div>
                    )}
                  </TableCell>
                  <TableCell>
                    {itemTypes.find((t) => t.id === item.type_id)?.name ||
                      "Unknown"}
                  </TableCell>
                  <TableCell>
                    {item.location?.name || "Not specified"}
                  </TableCell>
                  <TableCell className="text-right">
                    <div className="inline-flex">
                      {item.quantity}
                      {item.quantity_unit && (
                        <span className="ml-1">{item.quantity_unit}</span>
                      )}
                      {item.quantity <= 5 && (
                        <Badge className="ml-2 bg-red-100 text-red-800">
                          Low Stock
                        </Badge>
                      )}
                    </div>
                  </TableCell>
                  <TableCell className="text-right">
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => navigate({ to: `/inventory/${item.id}` })}
                      className="text-blue-600 hover:text-blue-900 mr-1"
                    >
                      View
                    </Button>
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() =>
                        navigate({ to: `/inventory/${item.id}/edit` })
                      }
                      className="text-indigo-600 hover:text-indigo-900 mr-1"
                    >
                      Edit
                    </Button>
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => handleDeleteItem(item.id)}
                      className="text-red-600 hover:text-red-900"
                    >
                      Delete
                    </Button>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </div>
      )}

      {/* Pagination info */}
      {filteredItems.length > 0 && (
        <div className="flex items-center justify-between mt-4">
          <div className="text-sm text-gray-500">
            {sortField !== "name" && (
              <span>
                Sorted by: <strong>{sortField}</strong> (
                {sortOrder === "asc" ? "ascending" : "descending"})
              </span>
            )}
          </div>
          <div className="text-sm text-gray-500">
            Showing <span className="font-medium">{sortedItems.length}</span> of{" "}
            <span className="font-medium">
              {itemsResponse?.data?.length || 0}
            </span>{" "}
            items
          </div>
        </div>
      )}
    </div>
  );
}
