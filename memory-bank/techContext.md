# Technical Context

This document outlines the technologies used, the development setup, technical constraints, and project dependencies. It also notes patterns in tool usage.

## Technologies Used

- **Backend:** Go 1.24
- **HTTP Framework:** Chi v5
- **Database:** PostgreSQL 17
- **Database Access:** pgx v5
- **Migration Tool:** golang-migrate
- **Frontend:** React with TypeScript
- **UI Components:** Shadcn
- **Styling:** Tailwind CSS
- **Frontend Deployment:** Progressive Web App (PWA)

## Development Setup

- **Backend Environment:**

  - Go 1.24 installed locally
  - PostgreSQL 17 database running in Docker container
  - Air for hot-reloading during development
  - Docker Compose for orchestrating services
  - GitHub Actions for CI/CD pipeline

- **Frontend Environment:**

  - Node.js with npm for package management
  - React 19 and TypeScript for frontend development
  - Vite as the build tool and development server
  - ESLint and Prettier for code quality

- **Tooling:**

  - Git for version control
  - Docker and Docker Compose for containerization
  - VS Code (recommended IDE)
  - golangci-lint for backend code quality (with configuration issues noted that require further investigation)

## Technical Constraints

- **Backend Constraints:**

  - Avoid using ORMs or repository patterns; direct SQL queries preferred
  - Maintain minimalistic middleware to ensure performance and simplicity

- **Frontend Constraints:**

  - Ensure compatibility with modern browsers and mobile devices as a PWA
  - Leverage Shadcn components for consistent UI design

- **Security Constraints:**

  - Data encryption and secure handling of API keys and tokens

## Dependencies

- **Backend:**

  - chi v5: minimal HTTP router
  - pgx v5: PostgreSQL driver
  - golang-migrate: database schema management

- **Frontend:**

  - React: core frontend library
  - TypeScript: type safety and improved developer experience
  - Tailwind CSS: utility-first CSS framework
  - Shadcn: UI component library

- **External APIs:**

  - Google Books API for data enrichment (books barcode scanning)

## Tool Usage Patterns

- **Version Control:**

  - Frequent, descriptive commits and branching strategy for feature development and bug fixes

- **Local Development:**

  - Backend developed and tested locally with hot-reloading using Air or similar tools
  - Frontend development leveraging React dev tools and hot module replacement

- **Database Management:**

  - Use golang-migrate to handle schema versioning and migrations
  - Regular automated backups for data security

- **Testing Patterns:**

  - Unit tests for backend logic using Go testing libraries
  - Integration tests to ensure reliable interactions between components
