package babysitter

import (
	"log"
	"net/http"

	"babycare-api/internal/db"
	"babycare-api/internal/models"

	"github.com/gin-gonic/gin"
)

// GetWeeklyViews returns the total number of profile views in the last 7 days
// for the authenticated babysitter.
func (h *BabysitterHandler) GetWeeklyViews(c *gin.Context) {
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

	count, err := queries.GetWeeklyViewCount(c.Request.Context(), currentUser.ID)
	if err != nil {
		log.Printf("babysitter weekly_views: db query: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
		return
	}

	c.JSON(http.StatusOK, models.WeeklyViewsResponse{
		BabysitterID: currentUser.ID.String(),
		ViewCount:    count,
		PeriodDays:   7,
	})
}
