package reporting

import (
	"database/sql"

	"babycare-api/internal/config"
)

// ReportingHandler holds dependencies for all reporting endpoints.
type ReportingHandler struct {
	db  *sql.DB
	cfg *config.Config
}

// NewReportingHandler constructs a ReportingHandler.
func NewReportingHandler(db *sql.DB, cfg *config.Config) *ReportingHandler {
	return &ReportingHandler{db: db, cfg: cfg}
}
