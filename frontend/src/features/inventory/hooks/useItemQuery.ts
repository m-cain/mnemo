import { useQuery } from "@tanstack/react-query";
import { getItem } from "../../../lib/apiClient";
import type { ApiResponse, Item } from "../../../types/models";

/**
 * Hook to fetch a single item by ID
 */
export const useItemQuery = (itemId: string) => {
  return useQuery<ApiResponse<Item>>({
    queryKey: ["item", itemId],
    queryFn: () => getItem(itemId),
    enabled: !!itemId,
  });
};
