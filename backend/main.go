package main

import (
	"context"
	"database/sql"
	"embed" // Add embed import
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/m-cain/mnemo/backend/auth"
	"github.com/m-cain/mnemo/backend/home"
	"github.com/m-cain/mnemo/backend/inventory"
	"github.com/m-cain/mnemo/backend/router"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

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

	// Run database migrations using goose
	dbMigrate, err := sql.Open("pgx", dbURL) // Use sql.Open with pgx stdlib driver
	if err != nil {
		log.Fatalf("Unable to connect to database for migrations: %v\n", err)
	}
	defer dbMigrate.Close()

	goose.SetBaseFS(embedMigrations) // Set the base file system

	if err := goose.SetDialect("postgres"); err != nil { // Set the dialect
		log.Fatalf("goose: failed to set dialect: %v", err)
	}

	if err := goose.Up(dbMigrate, "migrations"); err != nil { // Run migrations
		log.Fatalf("goose: migrations failed: %v", err)
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
