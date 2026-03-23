package parent

import (
	"log"
	"net/http"

	"babycare-api/internal/db"
	"babycare-api/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SaveBabysitter adds a babysitter to the parent's saved list.
func (h *ParentHandler) SaveBabysitter(c *gin.Context) {
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

	var req models.SaveBabysitterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	babysitterID, err := uuid.Parse(req.BabysitterID)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid babysitter_id"})
		return
	}

	queries := db.New(h.db)
	ctx := c.Request.Context()

	// Verify the target user exists and is a babysitter.
	target, err := queries.GetUserByID(ctx, babysitterID)
	if err != nil || target.Role != db.UserRoleBabysitter {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "babysitter not found"})
		return
	}

	if err := queries.SaveBabysitter(ctx, db.SaveBabysitterParams{
		ParentID:     currentUser.ID,
		BabysitterID: babysitterID,
	}); err != nil {
		log.Printf("parent save_babysitter: db insert: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Message: "babysitter saved"})
}

// UnsaveBabysitter removes a babysitter from the parent's saved list.
func (h *ParentHandler) UnsaveBabysitter(c *gin.Context) {
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

	babysitterIDStr := c.Param("babysitter_id")
	babysitterID, err := uuid.Parse(babysitterIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid babysitter_id"})
		return
	}

	queries := db.New(h.db)

	if err := queries.UnsaveBabysitter(c.Request.Context(), db.UnsaveBabysitterParams{
		ParentID:     currentUser.ID,
		BabysitterID: babysitterID,
	}); err != nil {
		log.Printf("parent unsave_babysitter: db delete: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Message: "babysitter removed from saved list"})
}

// ListSavedBabysitters returns all babysitters saved by the authenticated parent.
func (h *ParentHandler) ListSavedBabysitters(c *gin.Context) {
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

	rows, err := queries.ListSavedBabysitters(c.Request.Context(), currentUser.ID)
	if err != nil {
		log.Printf("parent list_saved_babysitters: db query: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
		return
	}

	response := make([]models.BabysitterProfileResponse, 0, len(rows))
	for _, r := range rows {
		response = append(response, models.BabysitterProfileResponse{
			UserID:            r.ID.String(),
			FullName:          r.FullName,
			Email:             r.Email,
			Phone:             r.Phone.String,
			Location:          r.Location.String,
			ProfilePictureURL: r.ProfilePictureUrl.String,
			Languages:         r.Languages,
			DaysPerWeek:       int(r.DaysPerWeek.Int32),
			HoursPerDay:       int(r.HoursPerDay.Int32),
			RateType:          string(r.RateType.RateType),
			RateAmount:        parseFloat(r.RateAmount.String),
			PaymentMethod:     r.PaymentMethod.String,
			IsApproved:        r.IsApproved,
			Gender:            r.Gender.String,
			Availability:      r.Availability,
			Currency:          r.Currency,
			IsAvailable:       r.IsAvailable,
		})
	}

	c.JSON(http.StatusOK, response)
}
