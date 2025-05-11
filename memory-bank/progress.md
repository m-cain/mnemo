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
  - Frontend React/TypeScript with Vite setup
  - Docker Compose configuration for development environment
  - Linting and formatting tools configured
  - CI/CD pipeline established with GitHub Actions
  - Development documentation created
- Database schema migrations for `users`, `homes`, `home_users`, and `locations` tables are created and populated.
- Core authentication logic (user registration and login with JWT) is implemented.
- Base backend API structure is set up with Chi router and registered authentication routes.
- Go module name is correctly configured as "github.com/m-cain/mnemo" at the project root.

## What's Left to Build

- Remaining system implementation:
  - Database schema implementation and migrations (remaining tables: `items`, `item_types`, `api_keys`)
  - Authentication system (API keys, user management, etc.)
  - Core API services for inventory management
  - Location hierarchy implementation with recursive queries
  - Frontend PWA with React and Shadcn components
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
- Project repository created with the following components:
  - Go 1.24 backend with initial module setup
  - React/TypeScript frontend with Vite
  - PostgreSQL 17 container in Docker Compose
  - Development environment configuration
- Ready to continue with Phase 2 and begin Phase 3.

## Known Issues

- Schema validation issue with the `.golangci.yml` file in the backend linter configuration. This issue persists and requires further investigation.
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
- Frontend architectural decisions guided by the need for PWA capabilities and responsive design
- Corrected the Go module name to "github.com/m-cain/mnemo" and relocated the `go.mod` file to the project root for proper module management.
