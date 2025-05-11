package auth

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/m-cain/mnemo/backend/contextkey"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("your_jwt_secret_key") // TODO: Load from configuration

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

// AuthService handles user authentication and authorization.
type AuthService struct {
	db            *pgxpool.Pool
	apiKeyService *APIKeyService // Inject APIKeyService
}

// NewAuthService creates a new AuthService.
func NewAuthService(db *pgxpool.Pool, apiKeyService *APIKeyService) *AuthService {
	return &AuthService{db: db, apiKeyService: apiKeyService}
}

// RegisterRoutes registers authentication routes.
func (s *AuthService) RegisterRoutes(r chi.Router) {
	r.Post("/register", s.handleRegister)
	r.Post("/login", s.handleLogin)
}

// AuthMiddleware is a middleware to authenticate requests using JWT or API Key.
func (s *AuthService) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check for JWT in Authorization header (Bearer token)
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && parts[0] == "Bearer" {
				tokenString := parts[1]
				token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"]) // Use fmt.Errorf
					}
					return jwtSecret, nil
				})

				if err == nil && token.Valid {
					if claims, ok := token.Claims.(jwt.MapClaims); ok {
						userID, ok := claims["sub"].(string)
						if ok {
							// Set userID in context
							ctx := context.WithValue(r.Context(), contextkey.UserIDKey, userID) // Use contextkey.UserIDKey
							next.ServeHTTP(w, r.WithContext(ctx))
							return
						}
					}
				}
			}
		}

		// Check for API Key in X-API-Key header
		apiKeyHeader := r.Header.Get("X-API-Key")
		if apiKeyHeader != "" {
			user, err := s.apiKeyService.ValidateAPIKey(r.Context(), apiKeyHeader)
			if err == nil {
				// Set userID in context
				ctx := context.WithValue(r.Context(), contextkey.UserIDKey, user.ID.String()) // Use contextkey.UserIDKey and user.ID
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		}

		// If neither JWT nor API Key is valid, return Unauthorized
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}

func (s *AuthService) handleRegister(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Basic validation
	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// Insert user into database
	userID := uuid.New()
	createdAt := time.Now()
	updatedAt := time.Now()

	query := `
		INSERT INTO users (id, email, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, email
	`
	var resp RegisterResponse
	// Use dbPool.QueryRow
	err = s.db.QueryRow(r.Context(), query, userID, req.Email, hashedPassword, createdAt, updatedAt).Scan(&resp.ID, &resp.Email)
	if err != nil {
		// TODO: Handle duplicate email error specifically
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (s *AuthService) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Basic validation
	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	// Retrieve user from database
	var userID uuid.UUID
	var hashedPassword string
	query := `SELECT id, password_hash FROM users WHERE email = $1`
	// Use dbPool.QueryRow
	err := s.db.QueryRow(r.Context(), query, req.Email).Scan(&userID, &hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Failed to retrieve user", http.StatusInternalServerError)
		return
	}

	// Compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID.String(),
		"exp": time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{Token: tokenString})
}
