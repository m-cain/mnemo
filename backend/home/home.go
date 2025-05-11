package home

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/m-cain/mnemo/backend/models"
)

// HomeService handles operations related to homes and home users.
type HomeService struct {
	db *pgxpool.Pool
}

// NewHomeService creates a new HomeService.
func NewHomeService(db *pgxpool.Pool) *HomeService {
	return &HomeService{db: db}
}

// CreateHome creates a new home and assigns the creating user as the owner.
func (s *HomeService) CreateHome(ctx context.Context, name string, ownerID string) (*models.Home, error) {
	homeID := uuid.New()
	ownerUUID, err := uuid.Parse(ownerID)
	if err != nil {
		return nil, fmt.Errorf("invalid owner ID: %w", err)
	}
	createdAt := time.Now()
	updatedAt := time.Now()

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx) // Rollback if not committed

	// Insert home
	homeQuery := `
		INSERT INTO homes (id, name, owner_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, name, owner_id, created_at, updated_at
	`
	home := &models.Home{}
	err = tx.QueryRow(ctx, homeQuery, homeID, name, ownerUUID, createdAt, updatedAt).Scan(
		&home.ID, &home.Name, &home.OwnerID, &home.CreatedAt, &home.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert home: %w", err)
	}

	// Insert home_user entry for the owner
	homeUserQuery := `
		INSERT INTO home_users (home_id, user_id, role, joined_at)
		VALUES ($1, $2, $3, $4)
	`
	_, err = tx.Exec(ctx, homeUserQuery, homeID, ownerUUID, "owner", time.Now()) // Assign 'owner' role
	if err != nil {
		return nil, fmt.Errorf("failed to insert home owner user: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return home, nil
}

// ListHomes lists all homes a user is a member of.
func (s *HomeService) ListHomes(ctx context.Context, userID string) ([]models.Home, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	query := `
		SELECT h.id, h.name, h.owner_id, h.created_at, h.updated_at
		FROM homes h
		JOIN home_users hu ON h.id = hu.home_id
		WHERE hu.user_id = $1
	`
	rows, err := s.db.Query(ctx, query, userUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to list homes for user %s: %w", userID, err)
	}
	defer rows.Close()

	var homes []models.Home
	for rows.Next() {
		var home models.Home
		err := rows.Scan(&home.ID, &home.Name, &home.OwnerID, &home.CreatedAt, &home.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan home: %w", err)
		}
		homes = append(homes, home)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during list homes rows iteration: %w", err)
	}

	return homes, nil
}

// GetHomeByID retrieves a specific home by its ID.
func (s *HomeService) GetHomeByID(ctx context.Context, homeID string) (*models.Home, error) {
	homeUUID, err := uuid.Parse(homeID)
	if err != nil {
		return nil, fmt.Errorf("invalid home ID: %w", err)
	}

	query := `
		SELECT id, name, owner_id, created_at, updated_at
		FROM homes
		WHERE id = $1
	`
	home := &models.Home{}
	err = s.db.QueryRow(ctx, query, homeUUID).Scan(
		&home.ID, &home.Name, &home.OwnerID, &home.CreatedAt, &home.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Home not found
		}
		return nil, fmt.Errorf("failed to get home by ID %s: %w", homeID, err)
	}

	return home, nil
}

// UpdateHome updates an existing home.
func (s *HomeService) UpdateHome(ctx context.Context, homeID string, name string) (*models.Home, error) {
	homeUUID, err := uuid.Parse(homeID)
	if err != nil {
		return nil, fmt.Errorf("invalid home ID: %w", err)
	}
	updatedAt := time.Now()

	query := `
		UPDATE homes
		SET name = $1, updated_at = $2
		WHERE id = $3
		RETURNING id, name, owner_id, created_at, updated_at
	`
	home := &models.Home{}
	err = s.db.QueryRow(ctx, query, name, updatedAt, homeUUID).Scan(
		&home.ID, &home.Name, &home.OwnerID, &home.CreatedAt, &home.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Home not found
		}
		return nil, fmt.Errorf("failed to update home %s: %w", homeID, err)
	}

	return home, nil
}

// DeleteHome deletes a home.
func (s *HomeService) DeleteHome(ctx context.Context, homeID string) error {
	homeUUID, err := uuid.Parse(homeID)
	if err != nil {
		return fmt.Errorf("invalid home ID: %w", err)
	}

	// Deleting a home should also delete associated home_users, locations, and items
	// Due to foreign key constraints with ON DELETE CASCADE, deleting the home
	// should automatically handle the related records.

	query := `
		DELETE FROM homes
		WHERE id = $1
	`
	result, err := s.db.Exec(ctx, query, homeUUID)
	if err != nil {
		return fmt.Errorf("failed to delete home %s: %w", homeID, err)
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows // Home not found
	}

	return nil
}

// ListHomeUsers lists all users associated with a home.
func (s *HomeService) ListHomeUsers(ctx context.Context, homeID string) ([]models.HomeUser, error) {
	homeUUID, err := uuid.Parse(homeID)
	if err != nil {
		return nil, fmt.Errorf("invalid home ID: %w", err)
	}

	query := `
		SELECT home_id, user_id, role, joined_at
		FROM home_users
		WHERE home_id = $1
	`
	rows, err := s.db.Query(ctx, query, homeUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to list home users for home %s: %w", homeID, err)
	}
	defer rows.Close()

	var homeUsers []models.HomeUser
	for rows.Next() {
		var homeUser models.HomeUser
		err := rows.Scan(&homeUser.HomeID, &homeUser.UserID, &homeUser.Role, &homeUser.JoinedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan home user: %w", err)
		}
		homeUsers = append(homeUsers, homeUser)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during list home users rows iteration: %w", err)
	}

	return homeUsers, nil
}

// InviteUserToHome invites a user to a home with a specific role.
// This is a simplified version; a real implementation would involve invitations,
// email notifications, and acceptance flows.
func (s *HomeService) InviteUserToHome(ctx context.Context, homeID string, userID string, role string) error {
	homeUUID, err := uuid.Parse(homeID)
	if err != nil {
		return fmt.Errorf("invalid home ID: %w", err)
	}
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	// Check if user is already a member
	checkQuery := `SELECT COUNT(*) FROM home_users WHERE home_id = $1 AND user_id = $2`
	var count int
	err = s.db.QueryRow(ctx, checkQuery, homeUUID, userUUID).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check existing home user: %w", err)
	}
	if count > 0 {
		return fmt.Errorf("user %s is already a member of home %s", userID, homeID)
	}

	// Insert home_user entry
	insertQuery := `
		INSERT INTO home_users (home_id, user_id, role, joined_at)
		VALUES ($1, $2, $3, $4)
	`
	_, err = s.db.Exec(ctx, insertQuery, homeUUID, userUUID, role, time.Now())
	if err != nil {
		return fmt.Errorf("failed to invite user to home: %w", err)
	}

	return nil
}

// UpdateHomeUserRole updates the role of a user in a home.
func (s *HomeService) UpdateHomeUserRole(ctx context.Context, homeID string, userID string, role string) error {
	homeUUID, err := uuid.Parse(homeID)
	if err != nil {
		return fmt.Errorf("invalid home ID: %w", err)
	}
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	query := `
		UPDATE home_users
		SET role = $1
		WHERE home_id = $2 AND user_id = $3
	`
	result, err := s.db.Exec(ctx, query, role, homeUUID, userUUID)
	if err != nil {
		return fmt.Errorf("failed to update home user role: %w", err)
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows // Home user not found
	}

	return nil
}

// RemoveUserFromHome removes a user from a home.
func (s *HomeService) RemoveUserFromHome(ctx context.Context, homeID string, userID string) error {
	homeUUID, err := uuid.Parse(homeID)
	if err != nil {
		return fmt.Errorf("invalid home ID: %w", err)
	}
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	query := `
		DELETE FROM home_users
		WHERE home_id = $1 AND user_id = $2
	`
	result, err := s.db.Exec(ctx, query, homeUUID, userUUID)
	if err != nil {
		return fmt.Errorf("failed to remove user from home: %w", err)
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows // Home user not found
	}

	return nil
}

// CheckHomeMembership checks if a user is a member of a home and returns their role.
func (s *HomeService) CheckHomeMembership(ctx context.Context, homeID string, userID string) (string, error) {
	homeUUID, err := uuid.Parse(homeID)
	if err != nil {
		return "", fmt.Errorf("invalid home ID: %w", err)
	}
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return "", fmt.Errorf("invalid user ID: %w", err)
	}

	query := `
		SELECT role
		FROM home_users
		WHERE home_id = $1 AND user_id = $2
	`
	var role string
	err = s.db.QueryRow(ctx, query, homeUUID, userUUID).Scan(&role)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", nil // Not a member
		}
		return "", fmt.Errorf("failed to check home membership: %w", err)
	}

	return role, nil
}
