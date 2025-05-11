import { useMutation, useQueryClient } from "@tanstack/react-query";
import { createItem } from "../../../lib/apiClient";
import type { Item } from "../../../types/models";

/**
 * Hook to create a new inventory item
 */
export const useCreateItemMutation = (homeId: string) => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (newItem: Omit<Item, "id" | "created_at" | "updated_at">) =>
      createItem(homeId, newItem),
    onSuccess: () => {
      // Invalidate the items query to refetch the updated list
      queryClient.invalidateQueries({ queryKey: ["items", homeId] });
    },
  });
};
