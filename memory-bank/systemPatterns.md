# System Patterns

This document details the system architecture, key technical decisions, and design patterns used in the project. It describes the relationships between components and critical implementation paths.

## Architecture Overview

The system adopts a modern, modular architecture, clearly separating frontend, backend, and data layers. The backend is built using Go 1.24, providing robust, performant RESTful APIs served via chi v5. The data persistence layer leverages PostgreSQL 17, interacting through direct SQL queries managed via pgx v5. The frontend uses React with TypeScript and Shadcn components, and is delivered as a Progressive Web App (PWA).

Potential future integration points, including machine learning models for image recognition and object detection, are considered within the architecture to allow easy extension.

### Implementation Status (Phase 1 Complete)

With Phase 1 complete, we have established the foundational infrastructure for the architecture:

- Project structure with separate backend and frontend directories
- Backend Go module initialization with planned Chi router integration
- Frontend React with TypeScript and Vite setup (awaiting Shadcn integration)
- Development environment with Docker Compose for PostgreSQL and backend services
- Hot-reloading configured for both backend (Air) and frontend (Vite)
- CI/CD pipeline with GitHub Actions for automated testing and building

## Key Technical Decisions

- **Backend Language:** Go 1.24, chosen for its performance, simplicity, concurrency capabilities, and robust ecosystem.
- **Database:** PostgreSQL 17, selected for reliability, robust data integrity, advanced querying capabilities, and performance.
- **Migration Tool:** golang-migrate, chosen for straightforward version control and automation of database schema changes.
- **Frontend Stack:** React with TypeScript for maintainability, reliability, and developer experience, combined with Shadcn for rapid, consistent UI development.
- **HTTP Framework:** Chi v5, selected due to its minimalistic design, ease of use, middleware support, and strong community.
- **Data Access:** Direct usage of pgx v5, avoiding ORM and repository patterns to simplify debugging, optimize performance, and increase control over SQL.

## Design Patterns

- **Service Layer Pattern:** Clear separation of business logic into distinct, reusable services.
- **Stateless REST API:** Emphasis on stateless interactions for scalability and ease of integration.
- **Command Query Responsibility Segregation (CQRS):** Clear delineation between read (queries) and write (commands) operations in business logic.
- **Dependency Injection:** Used to enhance testability and modularity in backend services.
- **Component-Based UI Architecture:** React components structured clearly, with state management and logic encapsulation for maintainability.

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
