package database

import (
	"context"
	"database/sql"
	"log"

	"babycare-api/internal/config"
	"babycare-api/internal/db"
	services_auth "babycare-api/internal/services/auth"
)

// SeedAdmin ensures that an admin user exists in the database.
// If the admin already exists the function returns silently (idempotent).
func SeedAdmin(database *sql.DB, cfg *config.Config) {
	queries := db.New(database)
	ctx := context.Background()

	// Check if the admin already exists.
	_, err := queries.GetUserByEmail(ctx, cfg.AdminEmail)
	if err == nil {
		// User already exists — nothing to do.
		return
	}
	if err != sql.ErrNoRows {
		log.Printf("seed: error checking for admin user: %v", err)
		return
	}

	// Hash the admin password.
	hash, err := services_auth.HashPassword(cfg.AdminPassword)
	if err != nil {
		log.Printf("seed: failed to hash admin password: %v", err)
		return
	}

	// Insert the admin user.
	_, err = queries.CreateAdminUser(ctx, db.CreateAdminUserParams{
		FullName: "BabyCare Admin",
		Email:    cfg.AdminEmail,
		PasswordHash: sql.NullString{
			String: hash,
			Valid:  true,
		},
	})
	if err != nil {
		log.Printf("seed: failed to create admin user: %v", err)
		return
	}

	log.Println("Admin user seeded successfully")
}
