package parent

import (
	"database/sql"
	"log"
	"net/http"

	"babycare-api/internal/db"
	"babycare-api/internal/models"

	"github.com/gin-gonic/gin"
)

// parentProfileResponse is the combined response shape for a parent profile.
type parentProfileResponse struct {
	ID              string `json:"id"`
	FullName        string `json:"full_name"`
	Email           string `json:"email"`
	Phone           string `json:"phone"`
	Status          string `json:"status"`
	Location        string `json:"location"`
	PrimaryLocation string `json:"primary_location"`
	Occupation      string `json:"occupation"`
	PreferredHours  string `json:"preferred_hours"`
	CreatedAt       string `json:"created_at"`
}

// GetProfile returns the authenticated parent's profile.
func (h *ParentHandler) GetProfile(c *gin.Context) {
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

	queries := db.New(h.db)

	profile, err := queries.GetParentProfileByUserID(c.Request.Context(), currentUser.ID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "profile not found"})
		return
	}
	if err != nil {
		log.Printf("parent get_profile: get profile: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
		return
	}

	c.JSON(http.StatusOK, buildParentResponse(currentUser, profile))
}

// buildParentResponse constructs the combined user + profile response.
func buildParentResponse(user db.User, profile db.ParentProfile) parentProfileResponse {
	return parentProfileResponse{
		ID:              user.ID.String(),
		FullName:        user.FullName,
		Email:           user.Email,
		Phone:           user.Phone.String,
		Status:          string(user.Status),
		Location:        profile.Location.String,
		PrimaryLocation: profile.PrimaryLocation.String,
		Occupation:      profile.Occupation.String,
		PreferredHours:  profile.PreferredHours.String,
		CreatedAt:       user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
