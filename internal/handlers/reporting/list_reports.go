package reporting

import (
	"log"
	"net/http"

	"babycare-api/internal/db"
	"babycare-api/internal/models"

	"github.com/gin-gonic/gin"
)

// ListReports returns all reports, optionally filtered by status.
// GET /api/v1/admin/reports?status=pending|resolved|dismissed
func (h *ReportingHandler) ListReports(c *gin.Context) {
	queries := db.New(h.db)
	ctx := c.Request.Context()

	statusFilter := c.Query("status")

	if statusFilter != "" {
		reports, err := queries.ListReportsByStatus(ctx, statusFilter)
		if err != nil {
			log.Printf("admin list_reports: %v", err)
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
			return
		}
		c.JSON(http.StatusOK, reports)
		return
	}

	reports, err := queries.ListReports(ctx)
	if err != nil {
		log.Printf("admin list_reports: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
		return
	}
	c.JSON(http.StatusOK, reports)
}
