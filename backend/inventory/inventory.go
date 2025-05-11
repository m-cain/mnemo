package inventory

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/m-cain/mnemo/backend/apperrors"
	"github.com/m-cain/mnemo/backend/models"
)

// InventoryService handles operations related to inventory items and types.
type InventoryService struct {
	db *pgxpool.Pool
}

// NewInventoryService creates a new instance of InventoryService.
func NewInventoryService(db *pgxpool.Pool) *InventoryService {
	return &InventoryService{db: db}
}

// ListItemTypes retrieves all item types from the database.
func (s *InventoryService) ListItemTypes(ctx context.Context) ([]models.ItemType, error) {
	query := `SELECT id, name, created_at, updated_at FROM item_types`

	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query item types: %w", err)
	}
	defer rows.Close()

	var itemTypes []models.ItemType
	for rows.Next() {
		var itemType models.ItemType
		if err := rows.Scan(&itemType.ID, &itemType.Name, &itemType.CreatedAt, &itemType.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan item type row: %w", err)
		}
		itemTypes = append(itemTypes, itemType)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error after scanning item type rows: %w", err)
	}

	return itemTypes, nil
}

// ListItems retrieves all items from the database for a given home.
func (s *InventoryService) ListItems(ctx context.Context, homeID uuid.UUID) ([]models.Item, error) {
	// This query needs to be refined to join with locations and filter by home_id
	// For now, a basic query is used.
	query := `SELECT id, name, quantity, unit, location_id, item_type_id, created_at, updated_at FROM items`

	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query items: %w", err)
	}
	defer rows.Close()

	var items []models.Item
	for rows.Next() {
		var item models.Item
		if err := rows.Scan(&item.ID, &item.Name, &item.Quantity, &item.Unit, &item.LocationID, &item.ItemTypeID, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan item row: %w", err)
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error after scanning item rows: %w", err)
	}

	return items, nil
}

// CreateItemType creates a new item type in the database.
func (s *InventoryService) CreateItemType(ctx context.Context, name string) (*models.ItemType, error) {
	query := `INSERT INTO item_types (name) VALUES ($1) RETURNING id, name, created_at, updated_at`

	var itemType models.ItemType
	err := s.db.QueryRow(ctx, query, name).Scan(&itemType.ID, &itemType.Name, &itemType.CreatedAt, &itemType.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to insert item type: %w", err)
	}

	return &itemType, nil
}

// GetItemTypeByID retrieves an item type by its ID from the database.
func (s *InventoryService) GetItemTypeByID(ctx context.Context, id uuid.UUID) (*models.ItemType, error) {
	query := `SELECT id, name, created_at, updated_at FROM item_types WHERE id = $1`

	var itemType models.ItemType
	err := s.db.QueryRow(ctx, query, id).Scan(&itemType.ID, &itemType.Name, &itemType.CreatedAt, &itemType.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Item type not found
		}
		return nil, fmt.Errorf("failed to query item type by ID: %w", err)
	}

	return &itemType, nil
}

// UpdateItemType updates an existing item type in the database.
func (s *InventoryService) UpdateItemType(ctx context.Context, id uuid.UUID, name string) (*models.ItemType, error) {
	query := `UPDATE item_types SET name = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2 RETURNING id, name, created_at, updated_at`

	var itemType models.ItemType
	err := s.db.QueryRow(ctx, query, name, id).Scan(&itemType.ID, &itemType.Name, &itemType.CreatedAt, &itemType.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Item type not found
		}
		return nil, fmt.Errorf("failed to update item type: %w", err)
	}

	return &itemType, nil
}

// DeleteItemType deletes an item type by its ID from the database.
func (s *InventoryService) DeleteItemType(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM item_types WHERE id = $1`

	result, err := s.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete item type: %w", err)
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows // Item type not found
	}

	return nil
}

// ListLocationsByHome retrieves all top-level locations for a given home.
func (s *InventoryService) ListLocationsByHome(ctx context.Context, homeID uuid.UUID) ([]models.Location, error) {
	query := `SELECT id, name, parent_location_id, home_id, created_at, updated_at FROM locations WHERE home_id = $1 AND parent_location_id IS NULL`

	rows, err := s.db.Query(ctx, query, homeID)
	if err != nil {
		return nil, fmt.Errorf("failed to query top-level locations: %w", err)
	}
	defer rows.Close()

	var locations []models.Location
	for rows.Next() {
		var location models.Location
		if err := rows.Scan(&location.ID, &location.Name, &location.ParentLocationID, &location.HomeID, &location.CreatedAt, &location.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan location row: %w", err)
		}
		locations = append(locations, location)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error after scanning location rows: %w", err)
	}

	return locations, nil
}

// ListLocationsByParent retrieves direct child locations for a given parent location.
func (s *InventoryService) ListLocationsByParent(ctx context.Context, parentLocationID uuid.UUID) ([]models.Location, error) {
	query := `SELECT id, name, parent_location_id, home_id, created_at, updated_at FROM locations WHERE parent_location_id = $1`

	rows, err := s.db.Query(ctx, query, parentLocationID)
	if err != nil {
		return nil, fmt.Errorf("failed to query child locations: %w", err)
	}
	defer rows.Close()

	var locations []models.Location
	for rows.Next() {
		var location models.Location
		if err := rows.Scan(&location.ID, &location.Name, &location.ParentLocationID, &location.HomeID, &location.CreatedAt, &location.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan location row: %w", err)
		}
		locations = append(locations, location)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error after scanning location rows: %w", err)
	}

	return locations, nil
}

// CreateLocation creates a new location in the database.
func (s *InventoryService) CreateLocation(ctx context.Context, name string, parentLocationID *uuid.UUID, homeID uuid.UUID) (*models.Location, error) {
	query := `INSERT INTO locations (name, parent_location_id, home_id) VALUES ($1, $2, $3) RETURNING id, name, parent_location_id, home_id, created_at, updated_at`

	var location models.Location
	err := s.db.QueryRow(ctx, query, name, parentLocationID, homeID).Scan(&location.ID, &location.Name, &location.ParentLocationID, &location.HomeID, &location.CreatedAt, &location.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to insert location: %w", err)
	}

	return &location, nil
}

// DeleteLocation deletes a location by its ID from the database.
// It will return an error if the location has child locations.
func (s *InventoryService) DeleteLocation(ctx context.Context, id uuid.UUID) error {
	// Check if the location has any child locations
	countQuery := `SELECT COUNT(*) FROM locations WHERE parent_location_id = $1`
	var count int
	err := s.db.QueryRow(ctx, countQuery, id).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check for child locations: %w", err)
	}
	if count > 0 {
		return fmt.Errorf("location has child locations and cannot be deleted")
	}

	// Check if the location has any items
	itemCountQuery := `SELECT COUNT(*) FROM items WHERE location_id = $1`
	var itemCount int
	err = s.db.QueryRow(ctx, itemCountQuery, id).Scan(&itemCount)
	if err != nil {
		return fmt.Errorf("failed to check for items in location: %w", err)
	}
	if itemCount > 0 {
		return fmt.Errorf("location contains items and cannot be deleted")
	}

	query := `DELETE FROM locations WHERE id = $1`
	result, err := s.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete location: %w", err)
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows // Location not found
	}

	return nil
}

// UpdateLocation updates an existing location in the database.
func (s *InventoryService) UpdateLocation(ctx context.Context, id uuid.UUID, name string, parentLocationID *uuid.UUID) (*models.Location, error) {
	query := `UPDATE locations SET name = $1, parent_location_id = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3 RETURNING id, name, parent_location_id, home_id, created_at, updated_at`

	var location models.Location
	err := s.db.QueryRow(ctx, query, name, parentLocationID, id).Scan(&location.ID, &location.Name, &location.ParentLocationID, &location.HomeID, &location.CreatedAt, &location.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, apperrors.ErrNotFound // Location not found
		}
		return nil, fmt.Errorf("failed to update location: %w", err)
	}

	return &location, nil
}

// GetLocationByID retrieves a location by its ID from the database.
func (s *InventoryService) GetLocationByID(ctx context.Context, id uuid.UUID) (*models.Location, error) {
	query := `SELECT id, name, parent_location_id, home_id, created_at, updated_at FROM locations WHERE id = $1`

	var location models.Location
	err := s.db.QueryRow(ctx, query, id).Scan(&location.ID, &location.Name, &location.ParentLocationID, &location.HomeID, &location.CreatedAt, &location.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, apperrors.ErrNotFound // Location not found
		}
		return nil, fmt.Errorf("failed to query location by ID: %w", err)
	}

	return &location, nil
}

// UpdateItemQuantity updates the quantity of an existing item in the database.
func (s *InventoryService) UpdateItemQuantity(ctx context.Context, id uuid.UUID, quantity int) error {
	query := `UPDATE items SET quantity = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`

	result, err := s.db.Exec(ctx, query, quantity, id)
	if err != nil {
		return fmt.Errorf("failed to update item quantity: %w", err)
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows // Item not found
	}

	return nil
}

// CreateItem creates a new item in the database.
func (s *InventoryService) CreateItem(ctx context.Context, item models.Item) (*models.Item, error) {
	query := `INSERT INTO items (name, quantity, unit, location_id, item_type_id, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
			  RETURNING id, name, quantity, unit, location_id, item_type_id, created_at, updated_at`

	var createdItem models.Item
	err := s.db.QueryRow(ctx, query,
		item.Name,
		item.Quantity,
		item.Unit,
		item.LocationID,
		item.ItemTypeID,
	).Scan(
		&createdItem.ID,
		&createdItem.Name,
		&createdItem.Quantity,
		&createdItem.Unit,
		&createdItem.LocationID,
		&createdItem.ItemTypeID,
		&createdItem.CreatedAt,
		&createdItem.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert item: %w", err)
	}

	return &createdItem, nil
}

// GetItemByID retrieves an item by its ID from the database.
func (s *InventoryService) GetItemByID(ctx context.Context, id uuid.UUID) (*models.Item, error) {
	query := `SELECT id, name, quantity, unit, location_id, item_type_id, created_at, updated_at FROM items WHERE id = $1`

	var item models.Item
	err := s.db.QueryRow(ctx, query, id).Scan(&item.ID, &item.Name, &item.Quantity, &item.Unit, &item.LocationID, &item.ItemTypeID, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, apperrors.ErrNotFound // Item not found
		}
		return nil, fmt.Errorf("failed to query item by ID: %w", err)
	}

	return &item, nil
}

// UpdateItem updates an existing item in the database.
func (s *InventoryService) UpdateItem(ctx context.Context, id uuid.UUID, item models.Item) (*models.Item, error) {
	query := `UPDATE items SET name = $1, quantity = $2, unit = $3, location_id = $4, item_type_id = $5, updated_at = CURRENT_TIMESTAMP WHERE id = $6 RETURNING id, name, quantity, unit, location_id, item_type_id, created_at, updated_at`

	var updatedItem models.Item
	err := s.db.QueryRow(ctx, query,
		item.Name,
		item.Quantity,
		item.Unit,
		item.LocationID,
		item.ItemTypeID,
		id,
	).Scan(
		&updatedItem.ID,
		&updatedItem.Name,
		&updatedItem.Quantity,
		&updatedItem.Unit,
		&updatedItem.LocationID,
		&updatedItem.ItemTypeID,
		&updatedItem.CreatedAt,
		&updatedItem.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Item not found
		}
		return nil, fmt.Errorf("failed to update item: %w", err)
	}

	return &updatedItem, nil
}

// DeleteItem deletes an item by its ID from the database.
func (s *InventoryService) DeleteItem(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM items WHERE id = $1`

	result, err := s.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete item: %w", err)
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows // Item not found
	}

	return nil
}
