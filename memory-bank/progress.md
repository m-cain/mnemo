# Progress

This document summarizes the current status of the project, what is working, what is left to build, and any known issues. It also reflects on the evolution of project decisions.

## What Works

- Project planning and architecture design are complete
- Data model design with self-referential location hierarchy
- Implementation phases and timeline have been established
- Component structure for both backend and frontend has been defined
- Project initialization and environment setup (Phase 1) complete:
  - Repository structure with essential files (README, LICENSE, .gitignore)
  - Backend Go module initialized
  - Frontend React/TypeScript with Vite, Tailwind CSS, and Shadcn UI setup
  - Docker Compose configuration for development environment
  - Linting and formatting tools configured
  - CI/CD pipeline established with GitHub Actions
  - Development documentation created
- Database schema migrations for `users`, `homes`, `home_users`, and `locations` tables are created and populated.
- Database schema migrations for `item_types`, `items`, and `api_keys` tables are created and populated.
- Core authentication logic (user registration and login with JWT) is implemented.
- Base backend API structure is set up with Chi router and registered authentication routes.
- Go module name is correctly configured as "github.com/m-cain/mnemo" at the project root.
- Backend API key management (generation, listing, revocation) is implemented.
- Backend home and home user management (CRUD for homes, listing home users, inviting users, updating roles, removing users) is implemented.
- Context key usage in backend is refactored to use defined constants.
- Backend CRUD operations for ItemType (List, Create, Get by ID, Update, Delete) are implemented.
- Backend CRUD operations for Item (Create and Get by ID) are implemented.
- Backend routing modularized by moving all route definitions and handlers from `backend/main.go` to a new `backend/router` package with separate files per route group.
- Refactored `backend/main.go` to be minimal, focusing on setup and server start using the new router.

## What's Left to Build

- Remaining system implementation:
  - Backend core development for authentication and home/user management is largely complete.
- Complete frontend PWA development:
  - Inventory management views (List, Detail, Create, Edit)
  - Location management views (List, Detail, Create, Edit)
  - API key management views
- Barcode scanning and camera integration
- Search functionality with NLP capabilities
  - Visual location management system
  - Integration with external APIs (e.g., Google Books)
  - Performance optimization and security hardening
  - Documentation and deployment guides
- Implement location hierarchy with recursive queries
- Frontend PWA development using React, TypeScript, Tailwind CSS, and Shadcn UI
- Barcode scanning and camera integration
- Search functionality with NLP capabilities
  - Visual location management system
  - Integration with external APIs (e.g., Google Books)
  - Performance optimization and security hardening
  - Documentation and deployment guides

## Current Status

- Phase 1 (Project Initialization & Setup) complete
- Progress made in Phase 2 (Backend Core Development):
  - Core database schema migrations completed
  - Core authentication logic implemented
  - Base API structure established
  - Implemented CRUD operations for ItemType.
  - Implemented Create, Get by ID, Update, and Delete operations for Item.
  - Modularized backend routing.
  - Defined and implemented a standard for "Not Found" errors using `apperrors.ErrNotFound`.
  - Implemented backend location hierarchy logic (CRUD and listing).
- Project repository created with the following components:
  - Go 1.24 backend with initial module setup
  - React/TypeScript frontend with Vite
  - PostgreSQL 17 container in Docker Compose
  - Development environment configuration
- Ready to continue with Phase 2 and begin Phase 3.

## Known Issues

- Persistent compiler error (`undefined: homeID`) in the `listItemsHandler` function in `backend/main.go` despite correct parsing and top-level import. This issue needs further investigation if it persists during build/testing.
- Potential challenges identified:
  - Performance optimization for recursive location queries
  - Camera API integration across different devices
  - Ensuring seamless offline functionality for PWA
  - Secure handling of API keys and sensitive data

## Evolution of Decisions

- Initial data model evolved to support flexible location hierarchies using a self-referential design
- Original storage concept was merged into a unified location entity for better flexibility
- Implementation approach shifted to a more progressive model, focusing on core functionality first
- Decision to use direct SQL over ORM made based on performance and debugging considerations
- Frontend architectural decisions guided by the need for PWA capabilities and responsive design, leveraging Tailwind CSS for styling and Shadcn UI for components.
- Corrected the Go module name to "github.com/m-cain/mnemo" and relocated the `go.mod` file to the project root for proper module management.
- Modularized backend routing to improve maintainability and optimize for LLM context window minimization
