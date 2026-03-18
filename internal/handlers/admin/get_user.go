package admin

import (
	"database/sql"
	"log"
	"net/http"

	"babycare-api/internal/db"
	"babycare-api/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// getUserDetailResponse is the combined shape returned by GetUser.
type getUserDetailResponse struct {
	models.UserResponse
	// Babysitter-only fields
	Location          *string  `json:"location,omitempty"`
	ProfilePictureURL *string  `json:"profile_picture_url,omitempty"`
	Languages         []string `json:"languages,omitempty"`
	DaysPerWeek       *int32   `json:"days_per_week,omitempty"`
	HoursPerDay       *int32   `json:"hours_per_day,omitempty"`
	RateType          *string  `json:"rate_type,omitempty"`
	RateAmount        *string  `json:"rate_amount,omitempty"`
	PaymentMethod     *string  `json:"payment_method,omitempty"`
	IsApproved        *bool    `json:"is_approved,omitempty"`
	// Parent-only fields
	Occupation     *string `json:"occupation,omitempty"`
	PreferredHours *string `json:"preferred_hours,omitempty"`
}

// GetUser returns a single user with their role-specific profile.
func (h *AdminHandler) GetUser(c *gin.Context) {
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

	queries := db.New(h.db)
	ctx := c.Request.Context()

	user, err := queries.GetUserByID(ctx, userID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "user not found"})
		return
	}
	if err != nil {
		log.Printf("admin get_user: get by id: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
		return
	}

	// Do not expose admin accounts.
	if user.Role == db.UserRoleAdmin {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "user not found"})
		return
	}

	resp := getUserDetailResponse{
		UserResponse: models.UserResponse{
			ID:        user.ID.String(),
			FullName:  user.FullName,
			Email:     user.Email,
			Phone:     user.Phone.String,
			Role:      string(user.Role),
			Status:    string(user.Status),
			CreatedAt: user.CreatedAt,
		},
	}

	switch user.Role {
	case db.UserRoleBabysitter:
		profile, err := queries.GetBabysitterProfileByUserID(ctx, user.ID)
		if err != nil && err != sql.ErrNoRows {
			log.Printf("admin get_user: get babysitter profile: %v", err)
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
			return
		}
		if err == nil {
			resp.Location = nullStringPtr(profile.Location)
			resp.ProfilePictureURL = nullStringPtr(profile.ProfilePictureUrl)
			resp.Languages = profile.Languages
			resp.DaysPerWeek = nullInt32Ptr(profile.DaysPerWeek)
			resp.HoursPerDay = nullInt32Ptr(profile.HoursPerDay)
			rateType := string(profile.RateType.RateType)
			if profile.RateType.Valid {
				resp.RateType = &rateType
			}
			resp.RateAmount = nullStringPtr(profile.RateAmount)
			resp.PaymentMethod = nullStringPtr(profile.PaymentMethod)
			resp.IsApproved = &profile.IsApproved
		}

	case db.UserRoleParent:
		profile, err := queries.GetParentProfileByUserID(ctx, user.ID)
		if err != nil && err != sql.ErrNoRows {
			log.Printf("admin get_user: get parent profile: %v", err)
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
			return
		}
		if err == nil {
			resp.Location = nullStringPtr(profile.Location)
			resp.Occupation = nullStringPtr(profile.Occupation)
			resp.PreferredHours = nullStringPtr(profile.PreferredHours)
		}
	}

	c.JSON(http.StatusOK, resp)
}

// nullStringPtr converts a sql.NullString to a *string (nil when not valid).
func nullStringPtr(ns sql.NullString) *string {
	if !ns.Valid {
		return nil
	}
	return &ns.String
}

// nullInt32Ptr converts a sql.NullInt32 to a *int32 (nil when not valid).
func nullInt32Ptr(ni sql.NullInt32) *int32 {
	if !ni.Valid {
		return nil
	}
	return &ni.Int32
}
