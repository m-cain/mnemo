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

1. Ensure Docker is running.
2. From the project root directory (`/Users/matt/src/mnemo`), start the Docker Compose services (database and backend):

   ```bash
   docker-compose up --build db backend
   ```

   This will build the backend development image, start the database container, and start the backend application with hot-reloading enabled via Air.

3. In a separate terminal, navigate to the `frontend` directory:
   ```bash
   cd frontend
   ```
4. Start the frontend development server:
   ```bash
   npm run dev # or yarn dev
   ```
   This will start the Vite development server, typically at `http://localhost:5173`.

The backend should be accessible at `http://localhost:8080` (or the port specified in `docker-compose.yml`), and the frontend development server will proxy API requests to the backend.

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
