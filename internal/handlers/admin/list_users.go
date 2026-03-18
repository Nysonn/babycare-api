package admin

import (
	"log"
	"net/http"

	"babycare-api/internal/db"
	"babycare-api/internal/models"

	"github.com/gin-gonic/gin"
)

// ListUsers returns all non-admin, non-deleted users.
func (h *AdminHandler) ListUsers(c *gin.Context) {
	queries := db.New(h.db)

	users, err := queries.ListUsers(c.Request.Context())
	if err != nil {
		log.Printf("admin list_users: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
		return
	}

	response := make([]models.UserResponse, 0, len(users))
	for _, u := range users {
		response = append(response, models.UserResponse{
			ID:        u.ID.String(),
			FullName:  u.FullName,
			Email:     u.Email,
			Phone:     u.Phone.String,
			Role:      string(u.Role),
			Status:    string(u.Status),
			CreatedAt: u.CreatedAt,
		})
	}

	// Return an empty array rather than null when there are no users.
	if response == nil {
		response = []models.UserResponse{}
	}

	c.JSON(http.StatusOK, response)
}
