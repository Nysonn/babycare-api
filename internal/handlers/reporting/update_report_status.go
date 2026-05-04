package reporting

import (
	"database/sql"
	"log"
	"net/http"
	"strings"

	"babycare-api/internal/db"
	"babycare-api/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var validReportStatuses = map[string]bool{
	"resolved":  true,
	"dismissed": true,
}

// UpdateReportStatus allows an admin to mark a report as resolved or dismissed.
// PUT /api/v1/admin/reports/:id
func (h *ReportingHandler) UpdateReportStatus(c *gin.Context) {
	idStr := c.Param("id")
	reportID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid report id"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	status := strings.ToLower(strings.TrimSpace(req.Status))
	if !validReportStatuses[status] {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "status must be one of: resolved, dismissed",
		})
		return
	}

	queries := db.New(h.db)
	ctx := c.Request.Context()

	report, err := queries.UpdateReportStatus(ctx, db.UpdateReportStatusParams{
		ID:     reportID,
		Status: status,
	})
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "report not found"})
		return
	}
	if err != nil {
		log.Printf("admin update_report_status: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Report status updated",
		"report":  report,
	})
}
