package babysitter

import (
	"log"
	"net/http"

	"babycare-api/internal/db"
	"babycare-api/internal/models"

	"github.com/gin-gonic/gin"
)

// GetProfileViews returns the list of parents who viewed the babysitter's profile.
// Restricted parent details (email, phone, location) are only included when
// the parent has sent at least one message to this babysitter.
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
		entry := models.ProfileViewResponse{
			ID:          v.ID.String(),
			ParentID:    v.ParentID.String(),
			ParentName:  v.ParentName,
			Occupation:  v.Occupation.String,
			ViewedAt:    v.ViewedAt,
			HasMessaged: v.HasMessaged,
		}

		// Only expose restricted fields when the parent has actually sent a message.
		if v.HasMessaged {
			entry.Email = v.Email
			entry.Phone = v.Phone.String
			entry.PrimaryLocation = v.PrimaryLocation.String
			entry.PreferredHours = v.PreferredHours.String
		}

		response = append(response, entry)
	}

	c.JSON(http.StatusOK, response)
}
