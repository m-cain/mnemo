import { useMutation, useQueryClient } from "@tanstack/react-query";
import { updateItem } from "../../../lib/apiClient";
import type { Item } from "../../../types/models";

/**
 * Hook to update an existing inventory item
 */
export const useUpdateItemMutation = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      itemId,
      itemData,
    }: {
      itemId: string;
      itemData: Partial<Item>;
    }) => updateItem(itemId, itemData),
    onSuccess: (data) => {
      // Invalidate the specific item query
      queryClient.invalidateQueries({ queryKey: ["item", data.data.id] });

      // Invalidate the items list for the home
      if (data.data.home_id) {
        queryClient.invalidateQueries({
          queryKey: ["items", data.data.home_id],
        });
      }
    },
  });
};
