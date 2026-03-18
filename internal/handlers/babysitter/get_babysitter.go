package babysitter

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"babycare-api/internal/db"
	"babycare-api/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const cacheTTLProfile = 10 * time.Minute

// GetBabysitter returns a single babysitter's full profile.
// Accessible to both parents and babysitters. Profile views are recorded
// asynchronously when the requester is a parent.
func (h *BabysitterHandler) GetBabysitter(c *gin.Context) {
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
	cacheKey := fmt.Sprintf("babysitter:%s", idStr)

	// Try the cache first.
	if h.cacheService != nil {
		if cached, err := h.cacheService.Get(ctx, cacheKey); err == nil {
			// Record profile view asynchronously before returning cached response.
			h.recordViewAsync(c)
			c.Data(http.StatusOK, "application/json; charset=utf-8", []byte(cached))
			return
		}
	}

	queries := db.New(h.db)

	user, err := queries.GetUserByID(ctx, userID)
	if err == sql.ErrNoRows || (err == nil && user.Role != db.UserRoleBabysitter) {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "babysitter not found"})
		return
	}
	if err != nil {
		log.Printf("babysitter get: get user: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
		return
	}

	profile, err := queries.GetBabysitterProfileByUserID(ctx, userID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "babysitter profile not found"})
		return
	}
	if err != nil {
		log.Printf("babysitter get: get profile: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
		return
	}

	response := models.BabysitterProfileResponse{
		UserID:            user.ID.String(),
		FullName:          user.FullName,
		Email:             user.Email,
		Phone:             user.Phone.String,
		Location:          profile.Location.String,
		ProfilePictureURL: profile.ProfilePictureUrl.String,
		Languages:         profile.Languages,
		DaysPerWeek:       int(profile.DaysPerWeek.Int32),
		HoursPerDay:       int(profile.HoursPerDay.Int32),
		RateType:          string(profile.RateType.RateType),
		RateAmount:        parseFloat(profile.RateAmount.String),
		PaymentMethod:     profile.PaymentMethod.String,
		IsApproved:        profile.IsApproved,
	}

	// Populate cache (failure is non-fatal).
	if h.cacheService != nil {
		if err := h.cacheService.Set(ctx, cacheKey, response, cacheTTLProfile); err != nil {
			log.Printf("babysitter get: cache set: %v", err)
		}
	}

	// Record profile view asynchronously if the requester is a parent.
	h.recordViewAsync(c)

	c.JSON(http.StatusOK, response)
}

// recordViewAsync records a profile view in the background if the requester is a parent.
func (h *BabysitterHandler) recordViewAsync(c *gin.Context) {
	// Extract babysitter ID from path.
	babysitterIDStr := c.Param("id")
	babysitterID, err := uuid.Parse(babysitterIDStr)
	if err != nil {
		return
	}

	// Check if the current requester is a parent.
	currentUserRaw, exists := c.Get("current_user")
	if !exists {
		return
	}
	currentUser, ok := currentUserRaw.(db.User)
	if !ok || currentUser.Role != db.UserRoleParent {
		return
	}

	parentID := currentUser.ID
	database := h.db

	go func() {
		queries := db.New(database)
		_, err := queries.RecordProfileView(context.Background(), db.RecordProfileViewParams{
			BabysitterID: babysitterID,
			ParentID:     parentID,
		})
		if err != nil {
			log.Printf("babysitter get: record profile view: %v", err)
		}
	}()
}
