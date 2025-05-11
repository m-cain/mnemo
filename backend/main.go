package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/golang-migrate/migrate/v4"
	migrate_postgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/m-cain/mnemo/backend/auth"
	"github.com/m-cain/mnemo/backend/home"
	"github.com/m-cain/mnemo/backend/inventory"
	"github.com/m-cain/mnemo/backend/router"
)

func main() {
	// Load environment variables
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable not set")
	}

	// Database connection pool for application services
	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		log.Fatalf("Unable to parse DATABASE_URL: %v\n", err)
	}

	dbPool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}
	defer dbPool.Close()

	// Check database connection
	err = dbPool.Ping(context.Background())
	if err != nil {
		log.Fatalf("Database connection failed: %v\n", err)
	}
	log.Println("Database connection pool established!")

	// Run database migrations using standard sql.DB connection
	dbMigrate, err := sql.Open("pgx", dbURL) // Use sql.Open with pgx stdlib driver
	if err != nil {
		log.Fatalf("Unable to connect to database for migrations: %v\n", err)
	}
	defer dbMigrate.Close()

	m, err := migrate.NewWithDatabaseInstance(
		"file://backend/migrations",
		"postgres", &migrate_postgres.Postgres{}) // Use the imported instance
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Printf("Database migrations failed or no change: %v", err)
	} else {
		log.Println("Database migrations applied successfully!")
	}

	// Initialize services
	apiKeyService := auth.NewAPIKeyService(dbPool)            // Initialize APIKeyService
	authService := auth.NewAuthService(dbPool, apiKeyService) // Pass dbPool and apiKeyService
	homeService := home.NewHomeService(dbPool)                // Initialize HomeService
	inventoryService := inventory.NewInventoryService(dbPool) // Initialize InventoryService

	// Setup router using the new router package
	r := router.NewRouter(dbPool, apiKeyService, authService, homeService, inventoryService)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}
	log.Printf("Server starting on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
