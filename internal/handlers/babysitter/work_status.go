package babysitter

import (
	"fmt"
	"log"
	"net/http"

	"babycare-api/internal/db"
	"babycare-api/internal/models"

	"github.com/gin-gonic/gin"
)

// SetWorkStatus allows the authenticated babysitter to toggle their availability.
// Unavailable babysitters are hidden from the public listing.
func (h *BabysitterHandler) SetWorkStatus(c *gin.Context) {
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

	var req models.SetWorkStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	queries := db.New(h.db)
	ctx := c.Request.Context()

	if err := queries.SetWorkStatus(ctx, db.SetWorkStatusParams{
		UserID:      currentUser.ID,
		IsAvailable: *req.IsAvailable,
	}); err != nil {
		log.Printf("babysitter work_status: db update: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
		return
	}

	// Invalidate cache so the list reflects the new status immediately.
	if h.cacheService != nil {
		cacheKey := fmt.Sprintf("babysitter:%s", currentUser.ID.String())
		if err := h.cacheService.Delete(ctx, cacheKey, cacheKeyBabysitterList); err != nil {
			log.Printf("babysitter work_status: cache invalidate: %v", err)
		}
	}

	status := "available"
	if !*req.IsAvailable {
		status = "unavailable"
	}
	c.JSON(http.StatusOK, models.SuccessResponse{Message: "work status set to " + status})
}
