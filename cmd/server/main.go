package main

import (
	"fmt"
	"log"

	"babycare-api/internal/config"
	"babycare-api/internal/database"
	"babycare-api/internal/router"
	services_auth "babycare-api/internal/services/auth"
	services_cache "babycare-api/internal/services/cache"
	services_email "babycare-api/internal/services/email"
	services_messaging "babycare-api/internal/services/messaging"
	"babycare-api/internal/services/storage"
)

func main() {
	// Load all configuration from environment variables.
	cfg := config.Load()

	// Establish the database connection pool.
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("database connection established")

	if err := database.RunMigrations(db, "db/migrations"); err != nil {
		log.Fatalf("failed to run database migrations: %v", err)
	}
	log.Println("database migrations applied")

	// Seed the admin user (idempotent — safe to call on every startup).
	database.SeedAdmin(db, cfg)

	// Initialise third-party service clients.
	clerkService := services_auth.NewClerkService(cfg.ClerkSecretKey)
	cloudinaryService := storage.NewCloudinaryService(
		cfg.CloudinaryCloudName,
		cfg.CloudinaryAPIKey,
		cfg.CloudinaryAPISecret,
	)
	emailService := services_email.NewEmailService(cfg.SendGridAPIKey)

	// Redis cache — non-fatal if unavailable. The app degrades gracefully.
	var cacheService *services_cache.CacheService
	if cs, err := services_cache.NewCacheService(cfg.RedisURL); err != nil {
		log.Printf("WARNING: Redis unavailable, caching disabled: %v", err)
	} else {
		cacheService = cs
		log.Println("Redis cache connected")
	}

	// Stream Chat — critical service, fatal if unavailable.
	streamService, err := services_messaging.NewStreamService(cfg.StreamAPIKey, cfg.StreamAPISecret)
	if err != nil {
		log.Fatalf("failed to initialise Stream Chat: %v", err)
	}
	log.Println("Stream Chat client initialised")

	// Build the Gin router with all middleware and routes.
	r := router.Setup(db, cfg, clerkService, cloudinaryService, emailService, cacheService, streamService)

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("babycare-api starting on %s (env: %s)", addr, cfg.AppEnv)

	if err := r.Run(addr); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
