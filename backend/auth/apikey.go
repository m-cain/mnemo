package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/m-cain/mnemo/backend/models"
	"golang.org/x/crypto/bcrypt"
)

// APIKeyService handles operations related to API keys.
type APIKeyService struct {
	db *pgxpool.Pool
}

// NewAPIKeyService creates a new APIKeyService.
func NewAPIKeyService(db *pgxpool.Pool) *APIKeyService {
	return &APIKeyService{db: db}
}

// GenerateAPIKey generates a new API key for a user.
func (s *APIKeyService) GenerateAPIKey(ctx context.Context, userID string, name string) (*models.APIKey, string, error) {
	// Generate a random key
	keyBytes := make([]byte, 32)
	_, err := rand.Read(keyBytes)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate random key: %w", err)
	}
	rawKey := base64.URLEncoding.EncodeToString(keyBytes)

	// Hash the key for storage
	hashedKey, err := bcrypt.GenerateFromPassword([]byte(rawKey), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", fmt.Errorf("failed to hash API key: %w", err)
	}

	// Insert into database
	apiKey := &models.APIKey{
		UserID:    userID,
		Name:      name,
		Key:       string(hashedKey), // Store the hashed key
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	query := `
		INSERT INTO api_keys (user_id, name, key, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, user_id, name, key, is_active, created_at, updated_at
	`
	err = s.db.QueryRow(ctx, query,
		apiKey.UserID,
		apiKey.Name,
		apiKey.Key,
		apiKey.IsActive,
		apiKey.CreatedAt,
		apiKey.UpdatedAt,
	).Scan(
		&apiKey.ID,
		&apiKey.UserID,
		&apiKey.Name,
		&apiKey.Key, // This will scan the hashed key
		&apiKey.IsActive,
		&apiKey.CreatedAt,
		&apiKey.UpdatedAt,
	)
	if err != nil {
		return nil, "", fmt.Errorf("failed to insert API key: %w", err)
	}

	// Return the APIKey model (with hashed key) and the raw key
	return apiKey, rawKey, nil
}

// ValidateAPIKey validates an API key.
func (s *APIKeyService) ValidateAPIKey(ctx context.Context, rawKey string) (*models.User, error) {
	// Retrieve the API key and associated user from the database
	query := `
		SELECT
			ak.id, ak.user_id, ak.name, ak.key, ak.is_active, ak.created_at, ak.updated_at,
			u.id, u.email, u.password_hash, u.created_at, u.updated_at
		FROM api_keys ak
		JOIN users u ON ak.user_id = u.id
		WHERE ak.is_active = TRUE
	`
	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query API keys for validation: %w", err)
	}
	defer rows.Close()

	// Iterate through active keys and compare the hash
	for rows.Next() {
		var apiKey models.APIKey
		var user models.User
		err := rows.Scan(
			&apiKey.ID, &apiKey.UserID, &apiKey.Name, &apiKey.Key, &apiKey.IsActive, &apiKey.CreatedAt, &apiKey.UpdatedAt,
			&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan API key and user during validation: %w", err)
		}

		err = bcrypt.CompareHashAndPassword([]byte(apiKey.Key), []byte(rawKey))
		if err == nil {
			// Match found, return the user
			return &user, nil
		} else if err != bcrypt.ErrMismatchedHashAndPassword {
			// Other bcrypt error
			return nil, fmt.Errorf("failed to compare API key hash: %w", err)
		}
		// Mismatched hash, continue to the next key
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during API key validation rows iteration: %w", err)
	}

	// No active key matched
	return nil, pgx.ErrNoRows
}

// ListAPIKeys lists all API keys for a user.
func (s *APIKeyService) ListAPIKeys(ctx context.Context, userID string) ([]models.APIKey, error) {
	query := `
		SELECT id, user_id, name, key, is_active, created_at, updated_at
		FROM api_keys
		WHERE user_id = $1
	`
	rows, err := s.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list API keys for user %s: %w", userID, err)
	}
	defer rows.Close()

	var apiKeys []models.APIKey
	for rows.Next() {
		var apiKey models.APIKey
		err := rows.Scan(&apiKey.ID, &apiKey.UserID, &apiKey.Name, &apiKey.Key, &apiKey.IsActive, &apiKey.CreatedAt, &apiKey.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan API key: %w", err)
		}
		apiKeys = append(apiKeys, apiKey)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during API key list rows iteration: %w", err)
	}

	return apiKeys, nil
}

// RevokeAPIKey revokes an API key by setting is_active to false.
func (s *APIKeyService) RevokeAPIKey(ctx context.Context, apiKeyID string, userID string) error {
	query := `
		UPDATE api_keys
		SET is_active = FALSE, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND user_id = $2
	`
	result, err := s.db.Exec(ctx, query, apiKeyID, userID)
	if err != nil {
		return fmt.Errorf("failed to revoke API key %s for user %s: %w", apiKeyID, userID, err)
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows // API key not found or doesn't belong to the user
	}

	return nil
}
