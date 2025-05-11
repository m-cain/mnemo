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
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/m-cain/mnemo/backend/auth"
	"github.com/m-cain/mnemo/backend/contextkey"
	"github.com/m-cain/mnemo/backend/home"
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
