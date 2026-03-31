package parent

import (
	"database/sql"
	"log"
	"net/http"

	"babycare-api/internal/db"
	"babycare-api/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// publicParentResponse is the subset of parent data exposed to babysitters.
type publicParentResponse struct {
	ID                string `json:"id"`
	FullName          string `json:"full_name"`
	Location          string `json:"location"`
	PrimaryLocation   string `json:"primary_location"`
	Occupation        string `json:"occupation"`
	PreferredHours    string `json:"preferred_hours"`
	ProfilePictureURL string `json:"profile_picture_url"`
}

// GetParent returns a parent's public profile by user ID.
// Accessible to authenticated babysitters.
func (h *ParentHandler) GetParent(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "id is required"})
		return
	}

	userID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid user id"})
		return
	}

	ctx := c.Request.Context()
	queries := db.New(h.db)

	user, err := queries.GetUserByID(ctx, userID)
	if err == sql.ErrNoRows || (err == nil && user.Role != db.UserRoleParent) {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "parent not found"})
		return
	}
	if err != nil {
		log.Printf("parent get: get user: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
		return
	}

	profile, err := queries.GetParentProfileByUserID(ctx, userID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "parent profile not found"})
		return
	}
	if err != nil {
		log.Printf("parent get: get profile: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
		return
	}

	c.JSON(http.StatusOK, publicParentResponse{
		ID:                user.ID.String(),
		FullName:          user.FullName,
		Location:          profile.Location.String,
		PrimaryLocation:   profile.PrimaryLocation.String,
		Occupation:        profile.Occupation.String,
		PreferredHours:    profile.PreferredHours.String,
		ProfilePictureURL: profile.ProfilePictureUrl.String,
	})
}
