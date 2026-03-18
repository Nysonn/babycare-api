package babysitter

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"babycare-api/internal/db"
	"babycare-api/internal/models"

	"github.com/gin-gonic/gin"
)

const maxUpdateUploadSize = 10 << 20 // 10 MB

// UpdateProfile allows the authenticated babysitter to update their own profile.
// Accepts multipart/form-data so a new profile picture can be uploaded optionally.
func (h *BabysitterHandler) UpdateProfile(c *gin.Context) {
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

	if err := c.Request.ParseMultipartForm(maxUpdateUploadSize); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "request too large or invalid multipart form"})
		return
	}

	queries := db.New(h.db)
	ctx := c.Request.Context()

	// Fetch current profile so unchanged fields retain their values.
	currentProfile, err := queries.GetBabysitterProfileByUserID(ctx, currentUser.ID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "profile not found"})
		return
	}
	if err != nil {
		log.Printf("babysitter update_profile: get current profile: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
		return
	}

	// --- Merge text fields: use form value if provided, else keep existing ---
	location := mergeString(c.Request.FormValue("location"), currentProfile.Location.String)
	paymentMethod := mergeString(c.Request.FormValue("payment_method"), currentProfile.PaymentMethod.String)
	rateAmountStr := mergeString(c.Request.FormValue("rate_amount"), currentProfile.RateAmount.String)

	// Languages
	languages := currentProfile.Languages
	if raw := c.Request.FormValue("languages"); raw != "" {
		languages = []string{}
		for _, lang := range strings.Split(raw, ",") {
			if trimmed := strings.TrimSpace(lang); trimmed != "" {
				languages = append(languages, trimmed)
			}
		}
	}

	// Numeric fields
	daysPerWeek := currentProfile.DaysPerWeek
	if raw := c.Request.FormValue("days_per_week"); raw != "" {
		if v, err := strconv.Atoi(raw); err == nil {
			daysPerWeek = sql.NullInt32{Int32: int32(v), Valid: true}
		}
	}
	hoursPerDay := currentProfile.HoursPerDay
	if raw := c.Request.FormValue("hours_per_day"); raw != "" {
		if v, err := strconv.Atoi(raw); err == nil {
			hoursPerDay = sql.NullInt32{Int32: int32(v), Valid: true}
		}
	}

	// Rate type
	rateType := currentProfile.RateType
	if raw := c.Request.FormValue("rate_type"); raw != "" {
		rateType = db.NullRateType{RateType: db.RateType(raw), Valid: true}
	}

	// Profile picture: upload new one if provided, otherwise keep existing URL.
	profilePictureURL := currentProfile.ProfilePictureUrl
	if _, fh, err := c.Request.FormFile("profile_picture"); err == nil {
		url, uploadErr := h.storageService.UploadFile(fh, "babycare/babysitters")
		if uploadErr != nil {
			log.Printf("babysitter update_profile: upload profile_picture: %v", uploadErr)
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "failed to upload profile picture"})
			return
		}
		profilePictureURL = sql.NullString{String: url, Valid: true}
	}

	updated, err := queries.UpdateBabysitterProfile(ctx, db.UpdateBabysitterProfileParams{
		UserID:            currentUser.ID,
		Location:          sql.NullString{String: location, Valid: location != ""},
		Languages:         languages,
		DaysPerWeek:       daysPerWeek,
		HoursPerDay:       hoursPerDay,
		RateType:          rateType,
		RateAmount:        sql.NullString{String: rateAmountStr, Valid: rateAmountStr != ""},
		PaymentMethod:     sql.NullString{String: paymentMethod, Valid: paymentMethod != ""},
		ProfilePictureUrl: profilePictureURL,
	})
	if err != nil {
		log.Printf("babysitter update_profile: db update: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
		return
	}

	// Invalidate stale cache entries.
	if h.cacheService != nil {
		cacheKey := fmt.Sprintf("babysitter:%s", currentUser.ID.String())
		if err := h.cacheService.Delete(ctx, cacheKey, cacheKeyBabysitterList); err != nil {
			log.Printf("babysitter update_profile: cache invalidate: %v", err)
		}
	}

	c.JSON(http.StatusOK, models.BabysitterProfileResponse{
		UserID:            currentUser.ID.String(),
		FullName:          currentUser.FullName,
		Email:             currentUser.Email,
		Phone:             currentUser.Phone.String,
		Location:          updated.Location.String,
		ProfilePictureURL: updated.ProfilePictureUrl.String,
		Languages:         updated.Languages,
		DaysPerWeek:       int(updated.DaysPerWeek.Int32),
		HoursPerDay:       int(updated.HoursPerDay.Int32),
		RateType:          string(updated.RateType.RateType),
		RateAmount:        parseFloat(updated.RateAmount.String),
		PaymentMethod:     updated.PaymentMethod.String,
		IsApproved:        updated.IsApproved,
	})
}

// mergeString returns next if non-empty, otherwise falls back to fallback.
func mergeString(next, fallback string) string {
	if next != "" {
		return next
	}
	return fallback
}
