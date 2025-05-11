package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/golang-migrate/migrate/v4"
	migrate_postgres "github.com/golang-migrate/migrate/v4/database/postgres" // Import with alias
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/m-cain/mnemo/backend/auth"
	"github.com/m-cain/mnemo/backend/contextkey"
	"github.com/m-cain/mnemo/backend/home"
	"github.com/m-cain/mnemo/backend/inventory"
	"github.com/m-cain/mnemo/backend/models"
)

func main() {
	// Load environment variables
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable not set")
	}

	// Database connection pool for application services
	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		log.Fatalf("Unable to parse DATABASE_URL: %v\n", err)
	}

	dbPool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}
	defer dbPool.Close()

	// Check database connection
	err = dbPool.Ping(context.Background())
	if err != nil {
		log.Fatalf("Database connection failed: %v\n", err)
	}
	log.Println("Database connection pool established!")

	// Run database migrations using standard sql.DB connection
	dbMigrate, err := sql.Open("pgx", dbURL) // Use sql.Open with pgx stdlib driver
	if err != nil {
		log.Fatalf("Unable to connect to database for migrations: %v\n", err)
	}
	defer dbMigrate.Close()

	m, err := migrate.NewWithDatabaseInstance(
		"file://backend/migrations",
		"postgres", &migrate_postgres.Postgres{}) // Use the imported instance
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Printf("Database migrations failed or no change: %v\n", err)
	} else {
		log.Println("Database migrations applied successfully!")
	}

	// Initialize services
	apiKeyService := auth.NewAPIKeyService(dbPool)            // Initialize APIKeyService
	authService := auth.NewAuthService(dbPool, apiKeyService) // Pass dbPool and apiKeyService
	homeService := home.NewHomeService(dbPool)                // Initialize HomeService
	inventoryService := inventory.NewInventoryService(dbPool) // Initialize InventoryService

	// Setup router
	r := chi.NewRouter()

	// Middleware
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
		authService.RegisterRoutes(r)

		// API Key Routes (Protected)
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

		// Inventory Routes (Protected)
		r.Route("/items", func(r chi.Router) {
			r.Use(authService.AuthMiddleware)    // Protect inventory routes
			r.Use(homeIDMiddleware(homeService)) // Middleware to check home membership and set homeID in context

			r.Get("/", listItemsHandler(inventoryService))
			r.Post("/", createItemHandler(inventoryService))
			// TODO: Add other item routes (GET by ID, PUT, DELETE, PUT quantity)
		})

		r.Route("/item-types", func(r chi.Router) {
			r.Use(authService.AuthMiddleware) // Protect item type routes
			r.Get("/", listItemTypesHandler(inventoryService))
			r.Post("/", createItemTypeHandler(inventoryService))
			r.Get("/{itemTypeID}", getItemTypeHandler(inventoryService))
			r.Put("/{itemTypeID}", updateItemTypeHandler(inventoryService))
			r.Delete("/{itemTypeID}", deleteItemTypeHandler(inventoryService))
		})

		// Home Routes (Protected)
		r.Route("/homes", func(r chi.Router) {
			r.Use(authService.AuthMiddleware) // Protect home routes

			// Middleware to check home membership and set homeID in context
			r.Use(func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					homeID := chi.URLParam(r, "homeID")
					if homeID != "" {
						userID := r.Context().Value(contextkey.UserIDKey).(string) // Use contextkey.UserIDKey
						role, err := homeService.CheckHomeMembership(r.Context(), homeID, userID)
						if err != nil {
							http.Error(w, "Failed to check home membership", http.StatusInternalServerError)
							log.Printf("Error checking home membership: %v", err)
							return
						}
						if role == "" {
							http.Error(w, "Forbidden", http.StatusForbidden) // User is not a member
							return
						}
						ctx := context.WithValue(r.Context(), contextkey.HomeIDKey, homeID) // Use contextkey.HomeIDKey
						ctx = context.WithValue(ctx, contextkey.UserRoleKey, role)          // Use contextkey.UserRoleKey
						next.ServeHTTP(w, r.WithContext(ctx))
						return
					}
					next.ServeHTTP(w, r) // Proceed if homeID is not in URL (e.g., POST /homes)
				})
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
				home, err := homeService.CreateHome(r.Context(), req.Name, userID)
				if err != nil {
					http.Error(w, "Failed to create home", http.StatusInternalServerError)
					log.Printf("Error creating home: %v", err)
					return
				}
				w.WriteHeader(http.StatusCreated)
				json.NewEncoder(w).Encode(home)
			})

			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				userID := r.Context().Value(contextkey.UserIDKey).(string) // Use contextkey.UserIDKey
				homes, err := homeService.ListHomes(r.Context(), userID)
				if err != nil {
					http.Error(w, "Failed to list homes", http.StatusInternalServerError)
					log.Printf("Error listing homes: %v", err)
					return
				}
				json.NewEncoder(w).Encode(homes)
			})

			r.Route("/{homeID}", func(r chi.Router) {
				r.Get("/", func(w http.ResponseWriter, r *http.Request) {
					homeID := chi.URLParam(r, "homeID")
					home, err := homeService.GetHomeByID(r.Context(), homeID)
					if err != nil {
						http.Error(w, "Failed to get home", http.StatusInternalServerError)
						log.Printf("Error getting home: %v", err)
						return
					}
					if home == nil {
						http.Error(w, "Home not found", http.StatusNotFound)
						return
					}
					json.NewEncoder(w).Encode(home)
				})

				r.Put("/", func(w http.ResponseWriter, r *http.Request) {
					homeID := chi.URLParam(r, "homeID")
					var req struct {
						Name string `json:"name"`
					}
					if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
						http.Error(w, "Invalid request body", http.StatusBadRequest)
						return
					}
					home, err := homeService.UpdateHome(r.Context(), homeID, req.Name)
					if err != nil {
						http.Error(w, "Failed to update home", http.StatusInternalServerError)
						log.Printf("Error updating home: %v", err)
						return
					}
					if home == nil {
						http.Error(w, "Home not found", http.StatusNotFound)
						return
					}
					json.NewEncoder(w).Encode(home)
				})

				r.Delete("/", func(w http.ResponseWriter, r *http.Request) {
					homeID := chi.URLParam(r, "homeID")
					err := homeService.DeleteHome(r.Context(), homeID)
					if err != nil {
						if err == pgx.ErrNoRows {
							http.Error(w, "Home not found", http.StatusNotFound)
						} else {
							http.Error(w, "Failed to delete home", http.StatusInternalServerError)
							log.Printf("Error deleting home: %v", err)
						}
						return
					}
					w.WriteHeader(http.StatusNoContent)
				})

				// Home User Management Routes
				r.Route("/users", func(r chi.Router) {
					r.Get("/", func(w http.ResponseWriter, r *http.Request) {
						homeID := chi.URLParam(r, "homeID")
						homeUsers, err := homeService.ListHomeUsers(r.Context(), homeID)
						if err != nil {
							http.Error(w, "Failed to list home users", http.StatusInternalServerError)
							log.Printf("Error listing home users: %v", err)
							return
						}
						json.NewEncoder(w).Encode(homeUsers)
					})

					r.Post("/invite", func(w http.ResponseWriter, r *http.Request) {
						homeID := chi.URLParam(r, "homeID")
						var req struct {
							UserID string `json:"user_id"` // User ID to invite
							Role   string `json:"role"`
						}
						if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
							http.Error(w, "Invalid request body", http.StatusBadRequest)
							return
						}
						err := homeService.InviteUserToHome(r.Context(), homeID, req.UserID, req.Role)
						if err != nil {
							http.Error(w, "Failed to invite user", http.StatusInternalServerError)
							log.Printf("Error inviting user: %v", err)
							return
						}
						w.WriteHeader(http.StatusNoContent)
					})

					r.Put("/{userID}", func(w http.ResponseWriter, r *http.Request) {
						homeID := chi.URLParam(r, "homeID")
						userID := chi.URLParam(r, "userID")
						var req struct {
							Role string `json:"role"`
						}
						if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
							http.Error(w, "Invalid request body", http.StatusBadRequest)
							return
						}
						err := homeService.UpdateHomeUserRole(r.Context(), homeID, userID, req.Role)
						if err != nil {
							if err == pgx.ErrNoRows {
								http.Error(w, "Home user not found", http.StatusNotFound)
							} else {
								http.Error(w, "Failed to update home user role", http.StatusInternalServerError)
								log.Printf("Error updating home user role: %v", err)
							}
							return
						}
						w.WriteHeader(http.StatusNoContent)
					})

					r.Delete("/{userID}", func(w http.ResponseWriter, r *http.Request) {
						homeID := chi.URLParam(r, "homeID")
						userID := chi.URLParam(r, "userID")
						err := homeService.RemoveUserFromHome(r.Context(), homeID, userID)
						if err != nil {
							if err == pgx.ErrNoRows {
								http.Error(w, "Home user not found", http.StatusNotFound)
							} else {
								http.Error(w, "Failed to remove user from home", http.StatusInternalServerError)
								log.Printf("Error removing user from home: %v", err)
							}
							return
						}
						w.WriteHeader(http.StatusNoContent)
					})
				})
			})
		})
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}
	log.Printf("Server starting on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

// homeIDMiddleware is a middleware that extracts the homeID from the URL
// and checks if the authenticated user is a member of that home.
// If the user is a member, it adds the homeID and user's role to the request context.
func homeIDMiddleware(homeService *home.HomeService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			homeIDStr := chi.URLParam(r, "homeID")
			if homeIDStr == "" {
				// If homeID is not in the URL, proceed without adding to context
				next.ServeHTTP(w, r)
				return
			}

			userID, ok := r.Context().Value(contextkey.UserIDKey).(string)
			if !ok {
				// UserID not found in context, AuthMiddleware should have handled this
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Check if the user is a member of the home
			role, err := homeService.CheckHomeMembership(r.Context(), homeIDStr, userID)
			if err != nil {
				log.Printf("Error checking home membership in middleware: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			if role == "" {
				// User is not a member of this home
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			// Add homeID and user role to the context
			ctx := context.WithValue(r.Context(), contextkey.HomeIDKey, homeIDStr)
			ctx = context.WithValue(ctx, contextkey.UserRoleKey, role)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
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

		_, err := uuid.Parse(homeIDStr) // Use blank identifier as homeID is not used directly yet
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

// createItemHandler returns a http.HandlerFunc that creates a new item.
func createItemHandler(inventoryService *inventory.InventoryService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		homeIDStr, ok := r.Context().Value(contextkey.HomeIDKey).(string)
		if !ok {
			http.Error(w, "Home ID not found in context", http.StatusInternalServerError)
			return
		}

		homeID, err := uuid.Parse(homeIDStr)
		if err != nil {
			http.Error(w, "Invalid Home ID format", http.StatusBadRequest)
			return
		}

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
