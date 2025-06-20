package router

import (
	"encoding/json"
	"errors" // Import the errors package
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/m-cain/mnemo/backend/apperrors" // Import the apperrors package
	"github.com/m-cain/mnemo/backend/auth"
	"github.com/m-cain/mnemo/backend/contextkey"
	"github.com/m-cain/mnemo/backend/home"
	"github.com/m-cain/mnemo/backend/inventory"
	"github.com/m-cain/mnemo/backend/models"
)

// RegisterInventoryItemRoutes registers the inventory item related routes.
func RegisterInventoryItemRoutes(r chi.Router, inventoryService *inventory.InventoryService, authService *auth.AuthService, homeService *home.HomeService) {
	r.Route("/items", func(r chi.Router) {
		r.Use(authService.AuthMiddleware)    // Protect inventory routes
		r.Use(homeIDMiddleware(homeService)) // Middleware to check home membership and set homeID in context

		r.Get("/", listItemsHandler(inventoryService))
		r.Post("/", createItemHandler(inventoryService))
		r.Get("/{itemID}", getItemByIDHandler(inventoryService))
		r.Put("/{itemID}", updateItemHandler(inventoryService))
		r.Delete("/{itemID}", deleteItemHandler(inventoryService))
		r.Put("/{itemID}/quantity", updateItemQuantityHandler(inventoryService))
	})
}

// updateItemQuantityHandler returns a http.HandlerFunc that updates the quantity of an existing item.
func updateItemQuantityHandler(inventoryService *inventory.InventoryService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		itemIDStr := chi.URLParam(r, "itemID")
		itemID, err := uuid.Parse(itemIDStr)
		if err != nil {
			http.Error(w, "Invalid item ID format", http.StatusBadRequest)
			return
		}

		var req struct {
			Quantity int `json:"quantity"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		err = inventoryService.UpdateItemQuantity(r.Context(), itemID, req.Quantity)
		if err != nil {
			if err == pgx.ErrNoRows {
				http.Error(w, "Item not found", http.StatusNotFound)
				return
			}
			http.Error(w, "Failed to update item quantity", http.StatusInternalServerError)
			log.Printf("Error updating item quantity: %v", err)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Item quantity updated successfully"})
	}
}

// listItemsHandler returns a http.HandlerFunc that lists items for a given home.
func listItemsHandler(inventoryService *inventory.InventoryService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		homeIDStr, ok := r.Context().Value(contextkey.HomeIDKey).(string)
		if !ok {
			http.Error(w, "Home ID not found in context", http.StatusInternalServerError)
			return
		}

		homeID, err := uuid.Parse(homeIDStr) // Correctly parse and assign homeID
		if err != nil {
			http.Error(w, "Invalid Home ID format", http.StatusBadRequest)
			return
		}

		items, err := inventoryService.ListItems(r.Context(), homeID)
		if err != nil {
			http.Error(w, "Failed to list items", http.StatusInternalServerError)
			log.Printf("Error listing items: %v", err)
			return
		}

		json.NewEncoder(w).Encode(items)
	}
}

// createItemHandler returns a http.HandlerFunc that creates a new item.
func createItemHandler(inventoryService *inventory.InventoryService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// homeIDStr, ok := r.Context().Value(contextkey.HomeIDKey).(string) // homeIDStr is not used in this handler currently
		// if !ok {
		// 	http.Error(w, "Home ID not found in context", http.StatusInternalServerError)
		// 	return
		// }

		// homeID, err := uuid.Parse(homeIDStr) // homeID is not used directly in this handler currently
		// if err != nil {
		// 	http.Error(w, "Invalid Home ID format", http.StatusBadRequest)
		// 	return
		// }

		var req models.Item
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Ensure the item is associated with the home from the context
		// req.HomeID = homeID // Items table does not have home_id directly, it's linked via location

		// For now, we'll assume location_id is provided in the request body
		// TODO: Implement proper location handling and validation

		createdItem, err := inventoryService.CreateItem(r.Context(), req)
		if err != nil {
			http.Error(w, "Failed to create item", http.StatusInternalServerError)
			log.Printf("Error creating item: %v", err)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(createdItem)
	}
}

// getItemByIDHandler returns a http.HandlerFunc that retrieves an item by its ID.
func getItemByIDHandler(inventoryService *inventory.InventoryService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		itemIDStr := chi.URLParam(r, "itemID")
		itemID, err := uuid.Parse(itemIDStr)
		if err != nil {
			http.Error(w, "Invalid item ID format", http.StatusBadRequest)
			return
		}

		item, err := inventoryService.GetItemByID(r.Context(), itemID)
		if err != nil {
			if errors.Is(err, apperrors.ErrNotFound) {
				http.Error(w, "Item not found", http.StatusNotFound)
				return
			}
			http.Error(w, "Failed to get item", http.StatusInternalServerError)
			log.Printf("Error getting item by ID: %v", err)
			return
		}

		json.NewEncoder(w).Encode(item)
	}
}

// updateItemHandler returns a http.HandlerFunc that updates an existing item.
func updateItemHandler(inventoryService *inventory.InventoryService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		itemIDStr := chi.URLParam(r, "itemID")
		itemID, err := uuid.Parse(itemIDStr)
		if err != nil {
			http.Error(w, "Invalid item ID format", http.StatusBadRequest)
			return
		}

		var req models.Item
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		updatedItem, err := inventoryService.UpdateItem(r.Context(), itemID, req)
		if err != nil {
			http.Error(w, "Failed to update item", http.StatusInternalServerError)
			log.Printf("Error updating item: %v", err)
			return
		}

		if updatedItem == nil {
			http.Error(w, "Item not found", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(updatedItem)
	}
}

// deleteItemHandler returns a http.HandlerFunc that deletes an item by its ID.
func deleteItemHandler(inventoryService *inventory.InventoryService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		itemIDStr := chi.URLParam(r, "itemID")
		itemID, err := uuid.Parse(itemIDStr)
		if err != nil {
			http.Error(w, "Invalid item ID format", http.StatusBadRequest)
			return
		}

		err = inventoryService.DeleteItem(r.Context(), itemID)
		if err != nil {
			if err == pgx.ErrNoRows {
				http.Error(w, "Item not found", http.StatusNotFound)
				return
			}
			http.Error(w, "Failed to delete item", http.StatusInternalServerError)
			log.Printf("Error deleting item: %v", err)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Item deleted successfully"})
	}
}
