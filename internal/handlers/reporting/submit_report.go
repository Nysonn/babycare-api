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

var validReportTypes = map[string]bool{
	"spam":          true,
	"harassment":    true,
	"inappropriate": true,
	"other":         true,
}

// SubmitReport allows a parent or babysitter to report another user.
// POST /api/v1/reports
func (h *ReportingHandler) SubmitReport(c *gin.Context) {
	currentUserRaw, exists := c.Get("current_user")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "unauthorized"})
		return
	}
	reporter := currentUserRaw.(db.User)

	var req struct {
		ReportedUserID string `json:"reported_user_id" binding:"required"`
		ReportType     string `json:"report_type" binding:"required"`
		Description    string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	reportType := strings.ToLower(strings.TrimSpace(req.ReportType))
	if !validReportTypes[reportType] {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "report_type must be one of: spam, harassment, inappropriate, other",
		})
		return
	}

	reportedID, err := uuid.Parse(req.ReportedUserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid reported_user_id"})
		return
	}

	// Prevent self-reporting.
	if reporter.ID == reportedID {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "you cannot report yourself"})
		return
	}

	queries := db.New(h.db)
	ctx := c.Request.Context()

	// Verify the reported user exists.
	if _, err := queries.GetUserByID(ctx, reportedID); err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "reported user not found"})
		return
	} else if err != nil {
		log.Printf("submit_report: get reported user: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
		return
	}

	var desc sql.NullString
	if strings.TrimSpace(req.Description) != "" {
		desc = sql.NullString{String: strings.TrimSpace(req.Description), Valid: true}
	}

	report, err := queries.CreateReport(ctx, db.CreateReportParams{
		ReporterID:     reporter.ID,
		ReportedUserID: reportedID,
		ReportType:     reportType,
		Description:    desc,
	})
	if err != nil {
		log.Printf("submit_report: create report: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":   "Report submitted successfully",
		"report_id": report.ID,
	})
}
