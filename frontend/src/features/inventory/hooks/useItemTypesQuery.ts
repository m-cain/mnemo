import { useQuery } from "@tanstack/react-query";
import { getItemTypes } from "../../../lib/apiClient";
import type { ApiResponse, ItemType } from "../../../types/models";

/**
 * Hook to fetch all item types for a specific home
 */
export const useItemTypesQuery = (homeId: string) => {
  return useQuery<ApiResponse<ItemType[]>>({
    queryKey: ["itemTypes", homeId],
    queryFn: () => getItemTypes(homeId),
    enabled: !!homeId,
  });
};
