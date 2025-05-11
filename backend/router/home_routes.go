package router

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/m-cain/mnemo/backend/auth"
	"github.com/m-cain/mnemo/backend/contextkey"
	"github.com/m-cain/mnemo/backend/home"
)

// RegisterHomeRoutes registers the home related routes.
func RegisterHomeRoutes(r chi.Router, homeService *home.HomeService, authService *auth.AuthService) {
	r.Route("/homes", func(r chi.Router) {
		r.Use(authService.AuthMiddleware) // Protect home routes

		// Middleware to check home membership and set homeID in context for /homes/{homeID} routes
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

		r.Post("/", createHomeHandler(homeService))
		r.Get("/", listHomesHandler(homeService))

		r.Route("/{homeID}", func(r chi.Router) {
			r.Get("/", getHomeByIDHandler(homeService))
			r.Put("/", updateHomeHandler(homeService))
			r.Delete("/", deleteHomeHandler(homeService))

			// Home User Management Routes
			r.Route("/users", func(r chi.Router) {
				r.Get("/", listHomeUsersHandler(homeService))
				r.Post("/invite", inviteUserToHomeHandler(homeService))
				r.Put("/{userID}", updateHomeUserRoleHandler(homeService))
				r.Delete("/{userID}", removeUserFromHomeHandler(homeService))
			})
		})
	})
}

// createHomeHandler returns a http.HandlerFunc that creates a new home.
func createHomeHandler(homeService *home.HomeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
	}
}

// listHomesHandler returns a http.HandlerFunc that lists homes for the authenticated user.
func listHomesHandler(homeService *home.HomeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(contextkey.UserIDKey).(string) // Use contextkey.UserIDKey
		homes, err := homeService.ListHomes(r.Context(), userID)
		if err != nil {
			http.Error(w, "Failed to list homes", http.StatusInternalServerError)
			log.Printf("Error listing homes: %v", err)
			return
		}
		json.NewEncoder(w).Encode(homes)
	}
}

// getHomeByIDHandler returns a http.HandlerFunc that retrieves a home by ID.
func getHomeByIDHandler(homeService *home.HomeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
	}
}

// updateHomeHandler returns a http.HandlerFunc that updates an existing home.
func updateHomeHandler(homeService *home.HomeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
	}
}

// deleteHomeHandler returns a http.HandlerFunc that deletes a home by ID.
func deleteHomeHandler(homeService *home.HomeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
	}
}

// listHomeUsersHandler returns a http.HandlerFunc that lists users for a given home.
func listHomeUsersHandler(homeService *home.HomeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		homeID := chi.URLParam(r, "homeID")
		homeUsers, err := homeService.ListHomeUsers(r.Context(), homeID)
		if err != nil {
			http.Error(w, "Failed to list home users", http.StatusInternalServerError)
			log.Printf("Error listing home users: %v", err)
			return
		}
		json.NewEncoder(w).Encode(homeUsers)
	}
}

// inviteUserToHomeHandler returns a http.HandlerFunc that invites a user to a home.
func inviteUserToHomeHandler(homeService *home.HomeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
	}
}

// updateHomeUserRoleHandler returns a http.HandlerFunc that updates a user's role in a home.
func updateHomeUserRoleHandler(homeService *home.HomeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
	}
}

// removeUserFromHomeHandler returns a http.HandlerFunc that removes a user from a home.
func removeUserFromHomeHandler(homeService *home.HomeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
	}
}
