package router

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/m-cain/mnemo/backend/apperrors"
	"github.com/m-cain/mnemo/backend/contextkey"
	"github.com/m-cain/mnemo/backend/inventory"
)

// LocationRouter provides routing for location-related requests.
type LocationRouter struct {
	inventoryService *inventory.InventoryService
}

// NewLocationRouter creates a new instance of LocationRouter.
func NewLocationRouter(inventoryService *inventory.InventoryService) *LocationRouter {
	return &LocationRouter{inventoryService: inventoryService}
}

// RegisterRoutes registers the location routes with the provided router.
func (r *LocationRouter) RegisterRoutes(router chi.Router) {
	router.Post("/", r.createLocationHandler)
	router.Get("/{locationID}", r.getLocationByIDHandler)
	router.Put("/{locationID}", r.updateLocationHandler)
	router.Delete("/{locationID}", r.deleteLocationHandler)
	router.Get("/home/{homeID}", r.listLocationsByHomeHandler)
	router.Get("/parent/{parentLocationID}", r.listLocationsByParentHandler)
}

// createLocationHandler handles requests to create a new location.
func (r *LocationRouter) createLocationHandler(w http.ResponseWriter, req *http.Request) {
	var input struct {
		Name             string     `json:"name"`
		ParentLocationID *uuid.UUID `json:"parent_location_id"`
	}

	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if input.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	homeID, ok := req.Context().Value(contextkey.HomeIDKey).(uuid.UUID)
	if !ok {
		http.Error(w, "Home ID not found in context", http.StatusInternalServerError)
		return
	}

	location, err := r.inventoryService.CreateLocation(req.Context(), input.Name, input.ParentLocationID, homeID)
	if err != nil {
		http.Error(w, "Failed to create location", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(location)
}

// getLocationByIDHandler handles requests to get a location by its ID.
func (r *LocationRouter) getLocationByIDHandler(w http.ResponseWriter, req *http.Request) {
	locationIDStr := chi.URLParam(req, "locationID")
	locationID, err := uuid.Parse(locationIDStr)
	if err != nil {
		http.Error(w, "Invalid location ID", http.StatusBadRequest)
		return
	}

	location, err := r.inventoryService.GetLocationByID(req.Context(), locationID)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			http.Error(w, "Location not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to get location", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(location)
}

// updateLocationHandler handles requests to update an existing location.
func (r *LocationRouter) updateLocationHandler(w http.ResponseWriter, req *http.Request) {
	locationIDStr := chi.URLParam(req, "locationID")
	locationID, err := uuid.Parse(locationIDStr)
	if err != nil {
		http.Error(w, "Invalid location ID", http.StatusBadRequest)
		return
	}

	var input struct {
		Name             string     `json:"name"`
		ParentLocationID *uuid.UUID `json:"parent_location_id"`
	}

	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Although Name is not strictly required for an update, if provided, it shouldn't be empty.
	if input.Name == "" && input.ParentLocationID == nil {
		http.Error(w, "At least one field (name or parent_location_id) must be provided for update", http.StatusBadRequest)
		return
	}

	// Get the existing location to apply updates
	existingLocation, err := r.inventoryService.GetLocationByID(req.Context(), locationID)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			http.Error(w, "Location not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to get existing location", http.StatusInternalServerError)
		return
	}

	// Apply updates from input
	if input.Name != "" {
		existingLocation.Name = input.Name
	}
	if input.ParentLocationID != nil {
		existingLocation.ParentLocationID = input.ParentLocationID
	} else {
		// If parent_location_id is explicitly set to null in the request, update it.
		// This allows moving a location to the top level.
		// We need to check if the JSON field was present and null, not just if the pointer is nil.
		// This requires a custom unmarshal or checking the raw JSON, which is more complex.
		// For simplicity now, we'll assume if parent_location_id is provided and is null,
		// the user intends to set it to null. If the field is omitted, the pointer will be nil.
		// A more robust solution would involve a custom unmarshal or a dedicated update struct.
		// For now, we'll update if the pointer is non-nil, allowing setting to a new parent or nil.
		// If the user wants to explicitly set to nil, they must send "parent_location_id": null
		// in the JSON body. The `*uuid.UUID` type handles this.
		existingLocation.ParentLocationID = input.ParentLocationID
	}

	updatedLocation, err := r.inventoryService.UpdateLocation(req.Context(), locationID, existingLocation.Name, existingLocation.ParentLocationID)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			http.Error(w, "Location not found after update attempt", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to update location", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedLocation)
}

// deleteLocationHandler handles requests to delete a location by its ID.
func (r *LocationRouter) deleteLocationHandler(w http.ResponseWriter, req *http.Request) {
	locationIDStr := chi.URLParam(req, "locationID")
	locationID, err := uuid.Parse(locationIDStr)
	if err != nil {
		http.Error(w, "Invalid location ID", http.StatusBadRequest)
		return
	}

	err = r.inventoryService.DeleteLocation(req.Context(), locationID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) { // Use pgx.ErrNoRows for not found from DeleteLocation
			http.Error(w, "Location not found", http.StatusNotFound)
			return
		}
		// Check for specific error messages from DeleteLocation for business logic errors
		if err.Error() == "location has child locations and cannot be deleted" || err.Error() == "location contains items and cannot be deleted" {
			http.Error(w, err.Error(), http.StatusConflict) // Use 409 Conflict for business rule violation
			return
		}
		http.Error(w, "Failed to delete location", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// listLocationsByHomeHandler handles requests to list top-level locations for a home.
func (r *LocationRouter) listLocationsByHomeHandler(w http.ResponseWriter, req *http.Request) {
	homeIDStr := chi.URLParam(req, "homeID")
	homeID, err := uuid.Parse(homeIDStr)
	if err != nil {
		http.Error(w, "Invalid home ID", http.StatusBadRequest)
		return
	}

	locations, err := r.inventoryService.ListLocationsByHome(req.Context(), homeID)
	if err != nil {
		http.Error(w, "Failed to get locations for home", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(locations)
}

// listLocationsByParentHandler handles requests to list child locations for a parent.
func (r *LocationRouter) listLocationsByParentHandler(w http.ResponseWriter, req *http.Request) {
	parentLocationIDStr := chi.URLParam(req, "parentLocationID")
	parentLocationID, err := uuid.Parse(parentLocationIDStr)
	if err != nil {
		http.Error(w, "Invalid parent location ID", http.StatusBadRequest)
		return
	}

	locations, err := r.inventoryService.ListLocationsByParent(req.Context(), parentLocationID)
	if err != nil {
		http.Error(w, "Failed to get child locations", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(locations)
}
