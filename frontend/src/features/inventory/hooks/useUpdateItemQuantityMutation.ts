import { useMutation, useQueryClient } from "@tanstack/react-query";
import { updateItemQuantity } from "../../../lib/apiClient";

/**
 * Hook to update just the quantity of an inventory item
 */
export const useUpdateItemQuantityMutation = (
  itemId: string,
  homeId?: string
) => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (quantity: number) => updateItemQuantity(itemId, quantity),
    onSuccess: () => {
      // Invalidate the specific item query
      queryClient.invalidateQueries({ queryKey: ["item", itemId] });

      // Invalidate the items list for the home if homeId is provided
      if (homeId) {
        queryClient.invalidateQueries({ queryKey: ["items", homeId] });
      }
    },
  });
};
