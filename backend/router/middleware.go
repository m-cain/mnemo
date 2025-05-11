package router

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/m-cain/mnemo/backend/contextkey"
	"github.com/m-cain/mnemo/backend/home"
)

// homeIDMiddleware is a middleware that extracts the homeID from the URL
// and checks if the authenticated user is a member of that home.
// If the user is a member, it adds the homeID and user's role to the request context.
func homeIDMiddleware(homeService *home.HomeService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			homeIDStr := chi.URLParam(r, "homeID")
			if homeIDStr == "" {
				// If homeID is not in the URL, proceed without adding to context
				// This is important for routes like POST /api/v1/items where homeID might not be a URL param
				// but expected in context from a higher-level middleware or derived differently.
				// However, for routes specifically under /homes/{homeID}/items, this middleware implies homeIDStr should exist.
				// The current logic seems to handle cases where it's optional at this middleware layer.
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
