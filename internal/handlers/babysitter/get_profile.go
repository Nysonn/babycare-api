package babysitter

import (
	"database/sql"
	"log"
	"net/http"

	"babycare-api/internal/db"
	"babycare-api/internal/models"

	"github.com/gin-gonic/gin"
)

// GetProfile returns the authenticated babysitter's own profile.
func (h *BabysitterHandler) GetProfile(c *gin.Context) {
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

	profile, err := queries.GetBabysitterProfileByUserID(c.Request.Context(), currentUser.ID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "profile not found"})
		return
	}
	if err != nil {
		log.Printf("babysitter get_profile: get profile: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
		return
	}

	c.JSON(http.StatusOK, models.BabysitterProfileResponse{
		UserID:            currentUser.ID.String(),
		FullName:          currentUser.FullName,
		Email:             currentUser.Email,
		Phone:             currentUser.Phone.String,
		Location:          profile.Location.String,
		ProfilePictureURL: profile.ProfilePictureUrl.String,
		Languages:         profile.Languages,
		DaysPerWeek:       int(profile.DaysPerWeek.Int32),
		HoursPerDay:       int(profile.HoursPerDay.Int32),
		RateType:          string(profile.RateType.RateType),
		RateAmount:        parseFloat(profile.RateAmount.String),
		PaymentMethod:     profile.PaymentMethod.String,
		IsApproved:        profile.IsApproved,
		Gender:            profile.Gender.String,
		Availability:      profile.Availability,
		Currency:          profile.Currency,
		IsAvailable:       profile.IsAvailable,
	})
}
