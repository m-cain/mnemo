# Development Setup

This document provides instructions on how to set up and run the Mnemo project for development.

## Prerequisites

- Docker and Docker Compose
- Go 1.24
- Node.js (LTS version recommended)
- npm or yarn

## Backend Setup

1. Navigate to the `backend` directory:
   ```bash
   cd backend
   ```
2. Initialize Go modules (if not already done):
   ```bash
   go mod init mnemo/backend
   ```
3. Clean up dependencies and generate go.sum:
   ```bash
   go mod tidy
   ```
4. Install Air for hot-reloading:
   ```bash
   go install github.com/cosmtrek/air@latest
   ```

## Frontend Setup

1. Navigate to the `frontend` directory:
   ```bash
   cd frontend
   ```
2. Install frontend dependencies:
   ```bash
   npm install # or yarn install
   ```

## Running the Development Environment

The project now uses a Makefile to simplify starting the various components. Ensure Docker is running before starting the database.

From the project root directory (`/Users/matt/src/mnemo`), you can use the following `make` commands:

- `make db`: Starts only the PostgreSQL database container using Docker Compose.
- `make backend`: Navigates to the `backend` directory and starts the Go backend server using `air` for hot-reloading.
- `make frontend`: Navigates to the `frontend` directory and starts the frontend development server using `npm run dev`.
- `make all`: Starts the database, backend, and frontend concurrently.
- `make clean`: Stops and removes the Docker containers defined in `docker-compose.yml`.

To start the full development environment, run:

```bash
make all
```

The backend will be accessible at `http://localhost:8080` (or the port configured in the backend code), and the frontend development server will typically be at `http://localhost:5173`. The frontend development server will proxy API requests to the backend.

## Database Migrations

Database schema migrations are handled by `golang-migrate`. Instructions for running migrations will be added in a later phase.

## Linting and Formatting

- **Frontend:**
  - Lint: `cd frontend && npm run lint`
  - Format: `cd frontend && npm run format`
- **Backend:**
  - Lint: `cd backend && golangci-lint run ./...` (Note: There is a known linter schema issue with `.golangci.yml` in the current environment.)
  - Format: `cd backend && goimports -w .`

## CI/CD

A basic GitHub Actions workflow is configured in `.github/workflows/ci.yml` to build and test the backend and frontend on push and pull request events.
