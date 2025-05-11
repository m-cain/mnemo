# Active Context

This document tracks the current focus of work, recent changes, and immediate next steps. It captures active decisions, considerations, important patterns, and recent learnings or insights.

## Current Focus

- Implementing backend core functionality for the Mnemo home inventory tracking system
- Setting up database schema and migrations
- Developing authentication system
- Implementing inventory management APIs
- Creating frontend components for core functionality
- Modularizing backend routing by moving route definitions and handlers from `backend/main.go` to a new `backend/router` package

## Recent Changes

- Completed Phase 1: Project Initialization & Setup
- Created project repository structure with essential files
- Initialized backend (Go module) and frontend (React with TypeScript and Vite)
- Set up development environment with Docker Compose for database and backend services
- Configured linting and formatting tools (ESLint, Prettier for frontend)
- Established CI/CD pipeline using GitHub Actions
- Created documentation for development setup
- Completed database schema migrations for `homes`, `home_users`, and `locations` tables.
- Implemented core authentication logic (registration and login) in `backend/auth/auth.go`.
- Set up base backend API structure in `backend/main.go` with Chi router and registered auth routes.
- Corrected Go module name to "github.com/m-cain/mnemo" and updated import paths.
- Completed database schema migrations for `item_types`, `items`, and `api_keys` tables.
- Implemented backend API key management (generation, listing, revocation).
- Implemented backend home and home user management (CRUD for homes, listing home users, inviting users, updating roles, removing users).
- Refactored context key usage in backend to use defined constants.
- Implemented backend CRUD operations for ItemType (List, Create, Get by ID, Update, Delete).
- Started implementing backend CRUD operations for Item (Implemented Create and Get by ID).
- Integrated new inventory handlers and routes into `backend/main.go`.
- Modularized backend routing by moving all route definitions and handlers from `backend/main.go` to a new `backend/router` package with separate files per route group.
- Refactored `backend/main.go` to be minimal, focusing on setup and server start using the new router.

## Next Steps

- Backend core development for authentication and home/user management is largely complete.
- Implement core API services for inventory management.
- Implement location hierarchy logic.
- **Frontend setup is currently blocked by an environment issue preventing `tailwindcss` execution.** Need to resolve this before proceeding with Shadcn setup and frontend component development.
- **Resolve frontend `tailwindcss` environment issue (Highest Priority).** This is blocking Shadcn setup and frontend component development.
- If frontend remains blocked, continue backend development:
  - Implement remaining CRUD operations for Item (Update, Delete).
  - Implement location hierarchy logic.
- Investigate backend `.golangci.yml` schema validation issue.

## Active Decisions and Considerations

- Adopted a self-referential location model to allow for arbitrary hierarchical structures
- Chosen direct SQL over ORM for better performance and control
- Decided on a progressive implementation approach, starting with core functionality
- Planning to implement location hierarchy with recursive SQL queries
- Considering materialized paths for optimizing location hierarchy queries
- Modularized backend routing to improve maintainability and optimize for LLM context window minimization

## Important Patterns and Preferences

- Backend using Chi router for RESTful API endpoints
- Service-oriented architecture separating business logic
- Direct use of pgx for database access without ORM layers
- React with TypeScript and Shadcn for frontend development
- Test-driven development approach for core functionalities

## Learnings and Insights

- Self-referential database design is essential for supporting flexible location hierarchies
- Progressive implementation will enable faster iteration and validation
- Clear separation between authentication, inventory, and location services will improve maintainability
- Regular documentation updates in the memory bank will be crucial for project continuity
- Correct configuration of the Go module path at the project root is essential for resolving local package imports.
- Modularizing backend routing significantly improves code organization and LLM context window efficiency.
