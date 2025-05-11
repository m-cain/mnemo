import { useMemo } from "react";
import { Card } from "../../../../components/ui/card";
import type { Item, ItemType, Location } from "../../../../types/models";
import { ItemTypeDistributionChart } from "./ItemTypeDistributionChart";
import { LocationDistributionChart } from "./LocationDistributionChart";
import { LowStockChart } from "./LowStockChart";
import { QuantityTrendsChart } from "./QuantityTrendsChart";

interface InventoryChartsProps {
  items: Item[];
  itemTypes: ItemType[];
  locations: Location[];
  className?: string;
  isLoading?: boolean;
}

/**
 * Component to display various charts and visualizations for inventory data
 */
export function InventoryCharts({
  items,
  itemTypes,
  locations,
  className = "",
  isLoading = false,
}: InventoryChartsProps) {
  // Prepare data for charts
  const itemTypeData = useMemo(() => {
    const typeCounts: { [key: string]: number } = {};

    // Count items by type
    items.forEach((item) => {
      const typeId = item.type_id;
      typeCounts[typeId] = (typeCounts[typeId] || 0) + 1;
    });

    // Convert to chart-friendly format
    return Object.entries(typeCounts)
      .map(([typeId, count]) => {
        const typeName =
          itemTypes.find((t) => t.id === typeId)?.name || "Unknown";
        return {
          name: typeName,
          value: count,
          typeId,
        };
      })
      .sort((a, b) => b.value - a.value); // Sort by count, descending
  }, [items, itemTypes]);

  const locationData = useMemo(() => {
    const locationCounts: { [key: string]: number } = {};

    // Count items by location
    items.forEach((item) => {
      if (item.location_id) {
        locationCounts[item.location_id] =
          (locationCounts[item.location_id] || 0) + 1;
      } else {
        locationCounts["unassigned"] = (locationCounts["unassigned"] || 0) + 1;
      }
    });

    // Convert to chart-friendly format
    return Object.entries(locationCounts)
      .map(([locationId, count]) => {
        let locationName = "Unassigned";

        if (locationId !== "unassigned") {
          locationName =
            locations.find((l) => l.id === locationId)?.name || "Unknown";
        }

        return {
          name: locationName,
          value: count,
          locationId,
        };
      })
      .sort((a, b) => b.value - a.value); // Sort by count, descending
  }, [items, locations]);

  const quantityData = useMemo(() => {
    // Get top items by quantity
    return [...items]
      .sort((a, b) => b.quantity - a.quantity)
      .slice(0, 10)
      .map((item) => ({
        name: item.name,
        value: item.quantity,
        unit: item.quantity_unit || "",
        id: item.id,
      }));
  }, [items]);

  const lowStockData = useMemo(() => {
    // Get low stock items (quantity <= 5)
    return items
      .filter((item) => item.quantity <= 5)
      .map((item) => ({
        name: item.name,
        value: item.quantity,
        unit: item.quantity_unit || "",
        id: item.id,
      }))
      .sort((a, b) => a.value - b.value); // Sort by quantity, ascending
  }, [items]);

  if (isLoading) {
    return (
      <div className={`w-full ${className}`}>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
          {[1, 2, 3, 4].map((i) => (
            <Card key={i} className="p-6 h-64 animate-pulse flex flex-col">
              <div className="h-6 bg-gray-200 rounded mb-4 w-1/3"></div>
              <div className="flex-1 bg-gray-100 rounded"></div>
            </Card>
          ))}
        </div>
      </div>
    );
  }

  if (items.length === 0) {
    return (
      <div className={`w-full ${className}`}>
        <Card className="p-6 text-center">
          <p className="text-gray-500">
            No inventory items found. Add items to see visualizations.
          </p>
        </Card>
      </div>
    );
  }

  return (
    <div className={`w-full ${className}`}>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
        <Card className="p-6">
          <h3 className="text-lg font-medium text-gray-900 mb-4">
            Item Type Distribution
          </h3>
          <div className="h-56">
            <ItemTypeDistributionChart data={itemTypeData} />
          </div>
        </Card>

        <Card className="p-6">
          <h3 className="text-lg font-medium text-gray-900 mb-4">
            Location Distribution
          </h3>
          <div className="h-56">
            <LocationDistributionChart data={locationData} />
          </div>
        </Card>

        <Card className="p-6">
          <h3 className="text-lg font-medium text-gray-900 mb-4">
            Top 10 Items by Quantity
          </h3>
          <div className="h-56">
            <QuantityTrendsChart data={quantityData} />
          </div>
        </Card>

        <Card className="p-6">
          <h3 className="text-lg font-medium text-gray-900 mb-4">
            Low Stock Items
          </h3>
          <div className="h-56">
            <LowStockChart data={lowStockData} />
          </div>
        </Card>
      </div>
    </div>
  );
}
