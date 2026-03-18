package admin

import (
	"database/sql"

	"babycare-api/internal/config"
	"babycare-api/internal/services/email"
)

// AdminHandler holds the dependencies needed by all admin endpoints.
type AdminHandler struct {
	db           *sql.DB
	cfg          *config.Config
	emailService *email.EmailService
}

// NewAdminHandler constructs an AdminHandler with all required dependencies.
func NewAdminHandler(db *sql.DB, cfg *config.Config, emailService *email.EmailService) *AdminHandler {
	return &AdminHandler{
		db:           db,
		cfg:          cfg,
		emailService: emailService,
	}
}
