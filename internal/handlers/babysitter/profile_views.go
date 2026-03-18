package babysitter

import (
	"log"
	"net/http"

	"babycare-api/internal/db"
	"babycare-api/internal/models"

	"github.com/gin-gonic/gin"
)

// GetProfileViews returns a list of parents who viewed the authenticated babysitter's profile.
func (h *BabysitterHandler) GetProfileViews(c *gin.Context) {
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

	views, err := queries.ListProfileViewsForBabysitter(c.Request.Context(), currentUser.ID)
	if err != nil {
		log.Printf("babysitter profile_views: list views: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
		return
	}

	response := make([]models.ProfileViewResponse, 0, len(views))
	for _, v := range views {
		response = append(response, models.ProfileViewResponse{
			ID:         v.ID.String(),
			ParentID:   v.ParentID.String(),
			ParentName: v.ParentName,
			ViewedAt:   v.ViewedAt,
		})
	}

	c.JSON(http.StatusOK, response)
}
