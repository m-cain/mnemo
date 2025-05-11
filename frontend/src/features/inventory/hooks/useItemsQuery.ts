import { useQuery } from "@tanstack/react-query";
import { getItems } from "../../../lib/apiClient";
import type { ApiResponse, Item } from "../../../types/models";

/**
 * Hook to fetch all items for a specific home
 */
export const useItemsQuery = (homeId: string) => {
  return useQuery<ApiResponse<Item[]>>({
    queryKey: ["items", homeId],
    queryFn: () => getItems(homeId),
    enabled: !!homeId,
  });
};
