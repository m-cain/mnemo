import { Minus, Package, Plus } from "lucide-react";
import * as React from "react";
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

interface QuantityAdjustmentDialogProps {
  currentQuantity: number;
  onAdjust?: (data: {
    newQuantity: number;
    reason: string;
    notes: string;
  }) => void;
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
  const [quantity, setQuantity] = React.useState(currentQuantity);
  const [reason, setReason] = React.useState("");
  const [notes, setNotes] = React.useState("");
  const [open, setOpen] = React.useState(false);

  const handleIncrement = () => {
    setQuantity((prev) => prev + 1);
  };

  const handleDecrement = () => {
    setQuantity((prev) => (prev > 0 ? prev - 1 : 0));
  };

  const handleQuantityChange = (
    e: React.ChangeEvent<HTMLInputElement>
  ): void => {
    const value = parseInt(e.target.value);
    if (!isNaN(value) && value >= 0) {
      setQuantity(value);
    } else if (e.target.value === "") {
      setQuantity(0);
    }
  };

  const handleSubmit = () => {
    if (onAdjust) {
      onAdjust({
        newQuantity: quantity,
        reason,
        notes,
      });
    }
    setOpen(false);
  };

  const handleCancel = () => {
    setQuantity(currentQuantity);
    setReason("");
    setNotes("");
    setOpen(false);
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button variant="outline">Adjust Quantity</Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[425px]">
        <div className="flex flex-col gap-2 max-sm:items-center sm:flex-row sm:gap-4">
          <div
            className="flex size-9 shrink-0 items-center justify-center rounded-full border border-border"
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

        <div className="grid gap-4 py-4">
          <div className="space-y-2">
            <Label htmlFor="current-quantity">Current Quantity</Label>
            <div className="text-sm font-medium">{currentQuantity}</div>
          </div>

          <div className="space-y-2">
            <Label htmlFor="new-quantity">New Quantity</Label>
            <div className="flex items-center space-x-2">
              <Button
                type="button"
                variant="outline"
                size="icon"
                onClick={handleDecrement}
                className="h-8 w-8"
              >
                <Minus className="h-4 w-4" />
              </Button>
              <Input
                id="new-quantity"
                type="number"
                min="0"
                value={quantity}
                onChange={handleQuantityChange}
                className="w-20 text-center"
              />
              <Button
                type="button"
                variant="outline"
                size="icon"
                onClick={handleIncrement}
                className="h-8 w-8"
              >
                <Plus className="h-4 w-4" />
              </Button>
            </div>
          </div>

          <div className="space-y-2">
            <Label htmlFor="reason">Reason for Adjustment</Label>
            <Select value={reason} onValueChange={setReason}>
              <SelectTrigger id="reason">
                <SelectValue placeholder="Select a reason" />
              </SelectTrigger>
              <SelectContent>
                {reasonOptions.map((option) => (
                  <SelectItem key={option.value} value={option.value}>
                    {option.label}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          <div className="space-y-2">
            <Label htmlFor="notes">Notes</Label>
            <Textarea
              id="notes"
              placeholder="Add any additional details about this adjustment"
              value={notes}
              onChange={(e) => setNotes(e.target.value)}
              rows={3}
            />
          </div>
        </div>

        <DialogFooter>
          <Button variant="outline" onClick={handleCancel}>
            Cancel
          </Button>
          <Button onClick={handleSubmit} disabled={!reason}>
            Confirm Adjustment
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
