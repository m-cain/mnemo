# System Patterns

This document details the system architecture, key technical decisions, and design patterns used in the project. It describes the relationships between components and critical implementation paths.

## Architecture Overview

The system adopts a modern, modular architecture, clearly separating frontend, backend, and data layers. The backend is built using Go 1.24, providing robust, performant RESTful APIs served via chi v5. The data persistence layer leverages PostgreSQL 17, interacting through direct SQL queries managed via pgx v5. The frontend uses React with TypeScript and Shadcn components, and is delivered as a Progressive Web App (PWA).

Potential future integration points, including machine learning models for image recognition and object detection, are considered within the architecture to allow easy extension.

### Implementation Status (Phase 1 Complete)

With Phase 1 complete, we have established the foundational infrastructure for the architecture:

- Project structure with separate backend and frontend directories
- Backend Go module initialization with planned Chi router integration
- Frontend React with TypeScript and Vite setup, including Shadcn UI and Tailwind CSS integration
- Development environment with Docker Compose for PostgreSQL, backend run locally with Air, and orchestrated by a Makefile.
- Hot-reloading configured for both backend (Air) and frontend (Vite)
- CI/CD pipeline with GitHub Actions for automated testing and building

## Key Technical Decisions

- **Backend Language:** Go 1.24, chosen for its performance, simplicity, concurrency capabilities, and robust ecosystem.
- **Database:** PostgreSQL 17, selected for reliability, robust data integrity, advanced querying capabilities, and performance.
- **Migration Tool:** golang-migrate, chosen for straightforward version control and automation of database schema changes.
- **Frontend Stack:** React with TypeScript for maintainability, reliability, and developer experience, combined with Shadcn UI and Tailwind CSS for rapid, consistent UI development and styling.
- **HTTP Framework:** Chi v5, selected due to its minimalistic design, ease of use, middleware support, and strong community.
- **Data Access:** Direct usage of pgx v5, avoiding ORM and repository patterns to simplify debugging, optimize performance, and increase control over SQL.

## Design Patterns

- **Service Layer Pattern:** Clear separation of business logic into distinct, reusable services.
- **Stateless REST API:** Emphasis on stateless interactions for scalability and ease of integration.
- **Command Query Responsibility Segregation (CQRS):** Clear delineation between read (queries) and write (commands) operations in business logic.
- **Dependency Injection:** Used to enhance testability and modularity in backend services.
- **Component-Based UI Architecture:** React components structured clearly, with state management, data fetching using TanStack Query, and routing handled by TanStack Router for maintainability and performance.

## Component Relationships

- **API Layer:** Handles routing, request validation, and delegates business logic operations to services.
- **Auth Layer:** Integrates within API middleware for authentication and authorization, generating and managing tokens and API keys.
- **Data Layer:** Direct interaction with PostgreSQL through pgx, with queries embedded within service-layer logic.
- **Frontend (PWA):** Communicates exclusively with backend API, providing intuitive UI for inventory management and data visualization.
- **Future ML Layer:** Interfaces via API endpoints, providing modular integration capability for object detection and barcode recognition enhancements.

## Critical Implementation Paths

- **User Authentication & Authorization:** Secure registration, login, and API key generation flows.
- **Inventory Management Workflow:** Efficient and accurate item tracking through barcode scanning, updates via API endpoints, and visual spatial management.
- **Frontend Integration Path:** Seamless PWA operation, robust API communication, and responsive state updates.
- **Data Migration and Backup Path:** Regular, automated migration and backup operations ensuring data consistency and recoverability.
- **ML Integration Path (Future):** Clearly defined endpoints and data contracts to facilitate straightforward integration of image and barcode recognition services.

## Standardized "Not Found" Error Handling

To ensure consistency and clarity when a requested resource is not found, the backend uses a standardized error: `apperrors.ErrNotFound`.

- **Definition:** `apperrors.ErrNotFound` is a sentinel error defined in `backend/apperrors/errors.go`.
- **Purpose:** It provides a distinct error value to signal that a specific resource could not be located, avoiding ambiguous `nil, nil` return values in such scenarios.
- **Usage in Service Layer:**
  - Service methods that attempt to retrieve a single resource by ID (or similar unique identifier) should return `nil, apperrors.ErrNotFound` if the resource is not found in the data layer (e.g., if a database query returns `pgx.ErrNoRows`).
  - Example:
    ```go
    func (s *InventoryService) GetItemByID(ctx context.Context, id uuid.UUID) (*models.Item, error) {
        // ... database query ...
        if err == pgx.ErrNoRows {
            return nil, apperrors.ErrNotFound // Resource not found
        }
        if err != nil {
            return nil, fmt.Errorf("failed to query item: %w", err) // Handle other errors
        }
        // ... return resource, nil
    }
    ```
- **Usage in Handler Layer (or calling code):** - Code that calls service methods should check for `apperrors.ErrNotFound` using `errors.Is()`. - If the error is `apperrors.ErrNotFound`, the calling code should handle it appropriately, such as returning an HTTP 404 Not Found response from an HTTP handler. - Example:
  `go
    func getItemByIDHandler(inventoryService *inventory.InventoryService) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
            // ... get itemID ...
            item, err := inventoryService.GetItemByID(r.Context(), itemID)
            if err != nil {
                if errors.Is(err, apperrors.ErrNotFound) {
                    http.Error(w, "Item not found", http.StatusNotFound) // Return 404
                    return
                }
                http.Error(w, "Internal server error", http.StatusInternalServerError) // Handle other errors
                log.Printf("Error getting item by ID: %v", err)
                return
            }
            // ... respond with item
        }
    }
    `
  This standard promotes consistency and makes error handling for "not found" scenarios explicit and easier to manage.
