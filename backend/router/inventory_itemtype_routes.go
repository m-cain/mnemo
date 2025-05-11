package router

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/m-cain/mnemo/backend/auth"
	"github.com/m-cain/mnemo/backend/inventory"
)

// RegisterInventoryItemTypeRoutes registers the inventory item type related routes.
func RegisterInventoryItemTypeRoutes(r chi.Router, inventoryService *inventory.InventoryService, authService *auth.AuthService) {
	r.Route("/item-types", func(r chi.Router) {
		r.Use(authService.AuthMiddleware) // Protect item type routes
		r.Get("/", listItemTypesHandler(inventoryService))
		r.Post("/", createItemTypeHandler(inventoryService))
		r.Get("/{itemTypeID}", getItemTypeHandler(inventoryService))
		r.Put("/{itemTypeID}", updateItemTypeHandler(inventoryService))
		r.Delete("/{itemTypeID}", deleteItemTypeHandler(inventoryService))
	})
}

// listItemTypesHandler returns a http.HandlerFunc that lists all item types.
func listItemTypesHandler(inventoryService *inventory.InventoryService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		itemTypes, err := inventoryService.ListItemTypes(r.Context())
		if err != nil {
			http.Error(w, "Failed to list item types", http.StatusInternalServerError)
			log.Printf("Error listing item types: %v", err)
			return
		}

		json.NewEncoder(w).Encode(itemTypes)
	}
}

// createItemTypeHandler returns a http.HandlerFunc that creates a new item type.
func createItemTypeHandler(inventoryService *inventory.InventoryService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Name string `json:"name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		itemType, err := inventoryService.CreateItemType(r.Context(), req.Name)
		if err != nil {
			http.Error(w, "Failed to create item type", http.StatusInternalServerError)
			log.Printf("Error creating item type: %v", err)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(itemType)
	}
}

// getItemTypeHandler returns a http.HandlerFunc that retrieves an item type by ID.
func getItemTypeHandler(inventoryService *inventory.InventoryService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		itemTypeIDStr := chi.URLParam(r, "itemTypeID")
		itemTypeID, err := uuid.Parse(itemTypeIDStr)
		if err != nil {
			http.Error(w, "Invalid Item Type ID format", http.StatusBadRequest)
			return
		}

		itemType, err := inventoryService.GetItemTypeByID(r.Context(), itemTypeID)
		if err != nil {
			if err == pgx.ErrNoRows {
				http.Error(w, "Item type not found", http.StatusNotFound)
			} else {
				http.Error(w, "Failed to get item type", http.StatusInternalServerError)
				log.Printf("Error getting item type: %v", err)
			}
			return
		}

		if itemType == nil {
			http.Error(w, "Item type not found", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(itemType)
	}
}

// updateItemTypeHandler returns a http.HandlerFunc that updates an existing item type.
func updateItemTypeHandler(inventoryService *inventory.InventoryService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		itemTypeIDStr := chi.URLParam(r, "itemTypeID")
		itemTypeID, err := uuid.Parse(itemTypeIDStr)
		if err != nil {
			http.Error(w, "Invalid Item Type ID format", http.StatusBadRequest)
			return
		}

		var req struct {
			Name string `json:"name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		itemType, err := inventoryService.UpdateItemType(r.Context(), itemTypeID, req.Name)
		if err != nil {
			if err == pgx.ErrNoRows {
				http.Error(w, "Item type not found", http.StatusNotFound)
			} else {
				http.Error(w, "Failed to update item type", http.StatusInternalServerError)
				log.Printf("Error updating item type: %v", err)
			}
			return
		}

		if itemType == nil {
			http.Error(w, "Item type not found", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(itemType)
	}
}

// deleteItemTypeHandler returns a http.HandlerFunc that deletes an item type by ID.
func deleteItemTypeHandler(inventoryService *inventory.InventoryService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		itemTypeIDStr := chi.URLParam(r, "itemTypeID")
		itemTypeID, err := uuid.Parse(itemTypeIDStr)
		if err != nil {
			http.Error(w, "Invalid Item Type ID format", http.StatusBadRequest)
			return
		}

		err = inventoryService.DeleteItemType(r.Context(), itemTypeID)
		if err != nil {
			if err == pgx.ErrNoRows {
				http.Error(w, "Item type not found", http.StatusNotFound)
			} else {
				http.Error(w, "Failed to delete item type", http.StatusInternalServerError)
				log.Printf("Error deleting item type: %v", err)
			}
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
