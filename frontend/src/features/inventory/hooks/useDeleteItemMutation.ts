import { useMutation, useQueryClient } from "@tanstack/react-query";
import { deleteItem } from "../../../lib/apiClient";

/**
 * Hook to delete an inventory item
 */
export const useDeleteItemMutation = (homeId?: string) => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (itemId: string) => deleteItem(itemId),
    onSuccess: (_, itemId) => {
      // Invalidate the specific item query
      queryClient.invalidateQueries({ queryKey: ["item", itemId] });

      // Invalidate the items list for the home if homeId is provided
      if (homeId) {
        queryClient.invalidateQueries({ queryKey: ["items", homeId] });
      }
    },
  });
};
