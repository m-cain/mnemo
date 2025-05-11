import { zodResolver } from "@hookform/resolvers/zod";
import { MinusCircle, Package, PlusCircle } from "lucide-react";
import { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { Button } from "../../../components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "../../../components/ui/dialog";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "../../../components/ui/form";
import { Input } from "../../../components/ui/input";
import { Label } from "../../../components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "../../../components/ui/select";
import { Textarea } from "../../../components/ui/textarea";
import type { QuantityAdjustmentFormValues } from "../schemas/itemValidationSchema";
import { quantityAdjustmentSchema } from "../schemas/itemValidationSchema";

interface QuantityAdjustmentDialogProps {
  currentQuantity: number;
  onAdjust?: (data: QuantityAdjustmentFormValues) => void;
}

const reasonOptions = [
  { value: "stock-count", label: "Stock Count" },
  { value: "damage", label: "Damage" },
  { value: "sale", label: "Sale" },
  { value: "purchase", label: "Purchase" },
  { value: "other", label: "Other" },
];

export function QuantityAdjustmentDialog({
  currentQuantity,
  onAdjust,
}: QuantityAdjustmentDialogProps) {
  const [open, setOpen] = useState(false);
  const [quantity, setQuantity] = useState(currentQuantity);

  // Initialize form with React Hook Form and Zod validation
  const form = useForm<QuantityAdjustmentFormValues>({
    resolver: zodResolver(quantityAdjustmentSchema),
    defaultValues: {
      newQuantity: currentQuantity,
      reason: "",
      notes: "",
    },
  });

  // Handle quantity increment
  const handleIncrement = () => {
    const newQuantity = quantity + 1;
    setQuantity(newQuantity);
    form.setValue("newQuantity", newQuantity, { shouldValidate: true });
  };

  // Handle quantity decrement
  const handleDecrement = () => {
    if (quantity > 0) {
      const newQuantity = quantity - 1;
      setQuantity(newQuantity);
      form.setValue("newQuantity", newQuantity, { shouldValidate: true });
    }
  };

  // Handle form submission
  const onSubmit = (data: QuantityAdjustmentFormValues) => {
    if (onAdjust) {
      onAdjust(data);
    }
    setOpen(false);
    form.reset();
  };

  // Reset form when dialog opens or current quantity changes
  useEffect(() => {
    if (open) {
      setQuantity(currentQuantity);
      form.reset({
        newQuantity: currentQuantity,
        reason: "",
        notes: "",
      });
    }
  }, [open, currentQuantity, form]);

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button variant="outline">Adjust Quantity</Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[425px]">
        <div className="flex flex-col gap-2 max-sm:items-center sm:flex-row sm:gap-4">
          <div
            className="flex size-11 shrink-0 items-center justify-center rounded-full border border-border"
            aria-hidden="true"
          >
            <Package className="opacity-80" size={16} strokeWidth={2} />
          </div>
          <DialogHeader>
            <DialogTitle>Adjust Inventory Quantity</DialogTitle>
            <DialogDescription>
              Update the quantity and provide a reason for this adjustment.
            </DialogDescription>
          </DialogHeader>
        </div>

        <Form {...form}>
          <form
            onSubmit={form.handleSubmit(onSubmit)}
            className="grid gap-4 py-4"
          >
            <div className="space-y-2">
              <Label>Current Quantity</Label>
              <div className="flex items-center justify-between rounded-md border border-input bg-background px-3 py-2 text-sm">
                <span className="text-muted-foreground">{currentQuantity}</span>
              </div>
            </div>

            <FormField
              control={form.control}
              name="newQuantity"
              render={({ field }) => (
                <FormItem className="space-y-2">
                  <FormLabel>New Quantity</FormLabel>
                  <div className="flex items-center space-x-2">
                    <Button
                      type="button"
                      variant="outline"
                      size="icon"
                      onClick={handleDecrement}
                      className="h-8 w-8"
                    >
                      <MinusCircle className="h-4 w-4" />
                    </Button>
                    <FormControl>
                      <Input
                        type="number"
                        min="0"
                        className="w-20 text-center"
                        value={quantity}
                        onChange={(e) => {
                          const value = parseInt(e.target.value);
                          if (!isNaN(value)) {
                            setQuantity(value);
                            field.onChange(value);
                          }
                        }}
                      />
                    </FormControl>
                    <Button
                      type="button"
                      variant="outline"
                      size="icon"
                      onClick={handleIncrement}
                      className="h-8 w-8"
                    >
                      <PlusCircle className="h-4 w-4" />
                    </Button>
                  </div>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="reason"
              render={({ field }) => (
                <FormItem className="space-y-2">
                  <FormLabel className="flex items-center">
                    Reason for Adjustment{" "}
                    <span className="text-red-500 ml-1">*</span>
                  </FormLabel>
                  <Select
                    onValueChange={field.onChange}
                    defaultValue={field.value}
                  >
                    <FormControl>
                      <SelectTrigger>
                        <SelectValue placeholder="Select a reason" />
                      </SelectTrigger>
                    </FormControl>
                    <SelectContent>
                      {reasonOptions.map((option) => (
                        <SelectItem key={option.value} value={option.value}>
                          {option.label}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="notes"
              render={({ field }) => (
                <FormItem className="space-y-2">
                  <FormLabel>Notes</FormLabel>
                  <FormControl>
                    <Textarea
                      placeholder="Add any additional details about this adjustment"
                      className="resize-none"
                      {...field}
                      rows={3}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <DialogFooter className="mt-4">
              <Button
                variant="outline"
                type="button"
                onClick={() => setOpen(false)}
              >
                Cancel
              </Button>
              <Button type="submit" disabled={form.formState.isSubmitting}>
                Confirm Adjustment
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}
