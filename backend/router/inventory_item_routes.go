package router

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
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
		// TODO: Add other item routes (GET by ID, PUT, DELETE, PUT quantity)
	})
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
