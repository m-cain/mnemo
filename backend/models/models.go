package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system.
type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Exclude from JSON output
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Home represents a home entity.
type Home struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	OwnerID   uuid.UUID `json:"owner_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// HomeUser represents the relationship between a home and a user.
type HomeUser struct {
	HomeID   uuid.UUID `json:"home_id"`
	UserID   uuid.UUID `json:"user_id"`
	Role     string    `json:"role"`
	JoinedAt time.Time `json:"joined_at"`
}

// Location represents a location within a home.
type Location struct {
	ID               uuid.UUID  `json:"id"`
	HomeID           uuid.UUID  `json:"home_id"`
	ParentLocationID *uuid.UUID `json:"parent_location_id"` // Use pointer for nullable FK
	Name             string     `json:"name"`
	Type             string     `json:"type"`
	Metadata         []byte     `json:"metadata"` // Use []byte for JSONB
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// Item represents an inventory item.
type Item struct {
	ID         uuid.UUID  `json:"id"`
	Name       string     `json:"name"`
	Quantity   int        `json:"quantity"`
	Unit       string     `json:"unit"`
	LocationID *uuid.UUID `json:"location_id"`  // Use pointer for nullable FK
	ItemTypeID *uuid.UUID `json:"item_type_id"` // Use pointer for nullable FK
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// ItemType represents a type of item.
type ItemType struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// APIKey represents an API key for a user.
type APIKey struct {
	ID        uuid.UUID `json:"id"`
	UserID    string    `json:"user_id"` // Use string for UUID to match apikey.go
	Name      string    `json:"name"`
	Key       string    `json:"key"` // Store hashed key
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
