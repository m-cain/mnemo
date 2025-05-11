package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/m-cain/mnemo/backend/auth"
	"github.com/m-cain/mnemo/backend/home"
	"github.com/m-cain/mnemo/backend/inventory"
)

// NewRouter initializes and configures the main Chi router.
func NewRouter(dbPool *pgxpool.Pool, apiKeyService *auth.APIKeyService, authService *auth.AuthService, homeService *home.HomeService, inventoryService *inventory.InventoryService) http.Handler {
	r := chi.NewRouter()

	// Global Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)

	// CORS middleware (basic example, configure properly for production)
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*") // TODO: Restrict this in production
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
			if r.Method == "OPTIONS" {
				return
			}
			next.ServeHTTP(w, r)
		})
	})

	// Register API routes
	r.Route("/api/v1", func(r chi.Router) {
		// Register authentication routes (handled within auth package's RegisterRoutes)
		authService.RegisterRoutes(r)

		// Register other route groups
		RegisterAPIKeyRoutes(r, apiKeyService, authService, inventoryService) // Added inventoryService
		RegisterInventoryItemRoutes(r, inventoryService, authService, homeService)
		RegisterInventoryItemTypeRoutes(r, inventoryService, authService)
		RegisterHomeRoutes(r, homeService, authService, inventoryService) // Added inventoryService

		// Register location routes
		locationRouter := NewLocationRouter(inventoryService)
		r.Route("/locations", locationRouter.RegisterRoutes)
	})

	return r
}
