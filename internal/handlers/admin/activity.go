package admin

import (
	"log"
	"net/http"

	"babycare-api/internal/db"
	"babycare-api/internal/models"

	"github.com/gin-gonic/gin"
)

// activityLabel converts a 30-day message count to a human-readable label.
func activityLabel(count int64) string {
	switch {
	case count >= 51:
		return "High"
	case count >= 11:
		return "Medium"
	default:
		return "Low"
	}
}

// GetActivity returns all non-admin users with a 30-day message activity label.
func (h *AdminHandler) GetActivity(c *gin.Context) {
	queries := db.New(h.db)
	ctx := c.Request.Context()

	users, err := queries.ListUsers(ctx)
	if err != nil {
		log.Printf("admin get_activity: list users: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
		return
	}

	response := make([]models.ActivityResponse, 0, len(users))
	for _, u := range users {
		count, err := queries.GetUserMessageCount(ctx, u.ID)
		if err != nil {
			log.Printf("admin get_activity: message count for user %s: %v", u.ID, err)
			count = 0
		}

		response = append(response, models.ActivityResponse{
			UserID:        u.ID.String(),
			FullName:      u.FullName,
			Role:          string(u.Role),
			ActivityLabel: activityLabel(count),
			MessageCount:  count,
		})
	}

	c.JSON(http.StatusOK, response)
}
