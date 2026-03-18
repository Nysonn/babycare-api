package auth

import (
	"database/sql"

	"babycare-api/internal/config"
	services_auth "babycare-api/internal/services/auth"
	services_messaging "babycare-api/internal/services/messaging"
	"babycare-api/internal/services/storage"
)

// AuthHandler holds the dependencies needed by all auth endpoints.
type AuthHandler struct {
	db             *sql.DB
	clerkService   *services_auth.ClerkService
	storageService *storage.CloudinaryService
	streamService  *services_messaging.StreamService
	cfg            *config.Config
}

// NewAuthHandler constructs an AuthHandler with all required dependencies.
func NewAuthHandler(
	db *sql.DB,
	clerkService *services_auth.ClerkService,
	storageService *storage.CloudinaryService,
	streamService *services_messaging.StreamService,
	cfg *config.Config,
) *AuthHandler {
	return &AuthHandler{
		db:             db,
		clerkService:   clerkService,
		storageService: storageService,
		streamService:  streamService,
		cfg:            cfg,
	}
}
