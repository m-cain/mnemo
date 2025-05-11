import { z } from "zod";

/**
 * Validation schema for creating or updating inventory items
 */
export const itemSchema = z.object({
  name: z
    .string()
    .min(1, { message: "Name is required" })
    .max(100, { message: "Name must be less than 100 characters" }),
  description: z.string().optional(),
  type_id: z.string({ required_error: "Item type is required" }).min(1, {
    message: "Item type is required",
  }),
  barcode: z.string().optional(),
  quantity: z
    .number({ invalid_type_error: "Quantity must be a number" })
    .nonnegative({ message: "Quantity cannot be negative" })
    .default(1),
  quantity_unit: z.string().optional(),
  location_id: z.string().optional(),
  home_id: z.string({ required_error: "Home ID is required" }).min(1, {
    message: "Home ID is required",
  }),
});

/**
 * Type for the item form data
 */
export type ItemFormValues = z.infer<typeof itemSchema>;

/**
 * Validation schema for quantity adjustments
 */
export const quantityAdjustmentSchema = z.object({
  newQuantity: z
    .number({ invalid_type_error: "Quantity must be a number" })
    .nonnegative({ message: "Quantity cannot be negative" }),
  reason: z
    .string({ required_error: "Reason is required" })
    .min(1, { message: "Reason is required" }),
  notes: z.string().optional(),
});

/**
 * Type for the quantity adjustment form data
 */
export type QuantityAdjustmentFormValues = z.infer<
  typeof quantityAdjustmentSchema
>;
