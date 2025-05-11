package router

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/m-cain/mnemo/backend/auth"
	"github.com/m-cain/mnemo/backend/contextkey"
	"github.com/m-cain/mnemo/backend/inventory"
)

// RegisterAPIKeyRoutes registers the API key related routes.
func RegisterAPIKeyRoutes(r chi.Router, apiKeyService *auth.APIKeyService, authService *auth.AuthService, inventoryService *inventory.InventoryService) {
	r.Route("/api-keys", func(r chi.Router) {
		r.Use(authService.AuthMiddleware) // Protect API key routes

		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			userID := r.Context().Value(contextkey.UserIDKey).(string) // Use contextkey.UserIDKey
			keys, err := apiKeyService.ListAPIKeys(r.Context(), userID)
			if err != nil {
				http.Error(w, "Failed to list API keys", http.StatusInternalServerError)
				log.Printf("Error listing API keys: %v", err)
				return
			}
			json.NewEncoder(w).Encode(keys)
		})

		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			userID := r.Context().Value(contextkey.UserIDKey).(string) // Use contextkey.UserIDKey
			var req struct {
				Name string `json:"name"`
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "Invalid request body", http.StatusBadRequest)
				return
			}
			_, rawKey, err := apiKeyService.GenerateAPIKey(r.Context(), userID, req.Name) // Use _ for unused apiKey
			if err != nil {
				http.Error(w, "Failed to generate API key", http.StatusInternalServerError)
				log.Printf("Error generating API key: %v", err)
				return
			}
			// Return the generated raw key (only once)
			json.NewEncoder(w).Encode(map[string]string{"key": rawKey})
		})

		r.Delete("/{apiKeyID}", func(w http.ResponseWriter, r *http.Request) {
			userID := r.Context().Value(contextkey.UserIDKey).(string) // Use contextkey.UserIDKey
			apiKeyID := chi.URLParam(r, "apiKeyID")
			err := apiKeyService.RevokeAPIKey(r.Context(), apiKeyID, userID)
			if err != nil {
				if err == pgx.ErrNoRows { // Use pgx.ErrNoRows
					http.Error(w, "API key not found or does not belong to user", http.StatusNotFound)
				} else {
					http.Error(w, "Failed to revoke API key", http.StatusInternalServerError)
					log.Printf("Error revoking API key: %v", err)
				}
				return
			}
			w.WriteHeader(http.StatusNoContent)
		})
	})
}
