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

	var response []models.ReportResponse

	if statusFilter != "" {
		rows, err := queries.ListReportsByStatus(ctx, statusFilter)
		if err != nil {
			log.Printf("admin list_reports: %v", err)
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
			return
		}
		response = make([]models.ReportResponse, 0, len(rows))
		for _, r := range rows {
			var desc *string
			if r.Description.Valid {
				s := r.Description.String
				desc = &s
			}
			response = append(response, models.ReportResponse{
				ID:             r.ID.String(),
				ReporterID:     r.ReporterID.String(),
				ReportedUserID: r.ReportedUserID.String(),
				ReportType:     r.ReportType,
				Description:    desc,
				Status:         r.Status,
				CreatedAt:      r.CreatedAt,
				UpdatedAt:      r.UpdatedAt,
				ReporterName:   r.ReporterName,
				ReporterEmail:  r.ReporterEmail,
				ReportedName:   r.ReportedName,
				ReportedEmail:  r.ReportedEmail,
			})
		}
	} else {
		rows, err := queries.ListReports(ctx)
		if err != nil {
			log.Printf("admin list_reports: %v", err)
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
			return
		}
		response = make([]models.ReportResponse, 0, len(rows))
		for _, r := range rows {
			var desc *string
			if r.Description.Valid {
				s := r.Description.String
				desc = &s
			}
			response = append(response, models.ReportResponse{
				ID:             r.ID.String(),
				ReporterID:     r.ReporterID.String(),
				ReportedUserID: r.ReportedUserID.String(),
				ReportType:     r.ReportType,
				Description:    desc,
				Status:         r.Status,
				CreatedAt:      r.CreatedAt,
				UpdatedAt:      r.UpdatedAt,
				ReporterName:   r.ReporterName,
				ReporterEmail:  r.ReporterEmail,
				ReportedName:   r.ReportedName,
				ReportedEmail:  r.ReportedEmail,
			})
		}
	}

	c.JSON(http.StatusOK, response)
}
