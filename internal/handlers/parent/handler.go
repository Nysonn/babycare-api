package parent

import (
	"database/sql"

	"babycare-api/internal/config"
	"babycare-api/internal/services/storage"
)

// ParentHandler holds the dependencies needed by all parent endpoints.
type ParentHandler struct {
	db             *sql.DB
	storageService *storage.CloudinaryService
	cfg            *config.Config
}

// NewParentHandler constructs a ParentHandler with all required dependencies.
func NewParentHandler(db *sql.DB, storageService *storage.CloudinaryService, cfg *config.Config) *ParentHandler {
	return &ParentHandler{db: db, storageService: storageService, cfg: cfg}
}
