package parent

import (
	"database/sql"

	"babycare-api/internal/config"
)

// ParentHandler holds the dependencies needed by all parent endpoints.
type ParentHandler struct {
	db  *sql.DB
	cfg *config.Config
}

// NewParentHandler constructs a ParentHandler with all required dependencies.
func NewParentHandler(db *sql.DB, cfg *config.Config) *ParentHandler {
	return &ParentHandler{db: db, cfg: cfg}
}
