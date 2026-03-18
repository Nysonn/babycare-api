package babysitter

import (
	"log"
	"net/http"
	"time"

	"babycare-api/internal/db"
	"babycare-api/internal/models"

	"github.com/gin-gonic/gin"
)

const cacheKeyBabysitterList = "babysitters:list"
const cacheTTLList = 5 * time.Minute

// ListBabysitters returns all approved, active babysitters.
// Results are cached in Redis for 5 minutes.
func (h *BabysitterHandler) ListBabysitters(c *gin.Context) {
	ctx := c.Request.Context()

	// Try the cache first.
	if h.cacheService != nil {
		if cached, err := h.cacheService.Get(ctx, cacheKeyBabysitterList); err == nil {
			c.Data(http.StatusOK, "application/json; charset=utf-8", []byte(cached))
			return
		}
	}

	queries := db.New(h.db)

	rows, err := queries.ListApprovedBabysitters(ctx)
	if err != nil {
		log.Printf("babysitter list: db query: %v", err)
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
		})
	}

	// Populate the cache (failure is non-fatal).
	if h.cacheService != nil {
		if err := h.cacheService.Set(ctx, cacheKeyBabysitterList, response, cacheTTLList); err != nil {
			log.Printf("babysitter list: cache set: %v", err)
		}
	}

	c.JSON(http.StatusOK, response)
}
