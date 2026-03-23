package parent

import (
	"database/sql"
	"log"
	"net/http"

	"babycare-api/internal/db"
	"babycare-api/internal/models"

	"github.com/gin-gonic/gin"
)

// UpdateProfile allows the authenticated parent to update their profile fields.
func (h *ParentHandler) UpdateProfile(c *gin.Context) {
	currentUserRaw, exists := c.Get("current_user")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "unauthorised"})
		return
	}
	currentUser, ok := currentUserRaw.(db.User)
	if !ok {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "unauthorised"})
		return
	}

	var req models.UpdateParentProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	queries := db.New(h.db)
	ctx := c.Request.Context()

	updated, err := queries.UpdateParentProfile(ctx, db.UpdateParentProfileParams{
		UserID:          currentUser.ID,
		Location:        sql.NullString{String: req.Location, Valid: req.Location != ""},
		Occupation:      sql.NullString{String: req.Occupation, Valid: req.Occupation != ""},
		PreferredHours:  sql.NullString{String: req.PreferredHours, Valid: req.PreferredHours != ""},
		PrimaryLocation: sql.NullString{String: req.PrimaryLocation, Valid: req.PrimaryLocation != ""},
	})
	if err != nil {
		log.Printf("parent update_profile: db update: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
		return
	}

	c.JSON(http.StatusOK, buildParentResponse(currentUser, updated))
}
