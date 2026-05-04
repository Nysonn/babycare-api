package messaging

import (
	"log"
	"net/http"

	"babycare-api/internal/db"
	"babycare-api/internal/models"

	"github.com/gin-gonic/gin"
)

// ListConversations returns all conversations the authenticated user is a participant in.
func (h *MessagingHandler) ListConversations(c *gin.Context) {
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

	rows, err := queries.ListConversationsForUser(c.Request.Context(), currentUser.ID)
	if err != nil {
		log.Printf("messaging list_conversations: db query: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
		return
	}

	response := make([]models.ConversationResponse, 0, len(rows))
	for _, row := range rows {
		picURL, _ := row.OtherUserProfilePictureUrl.(string)
		response = append(response, models.ConversationResponse{
			ID:                         row.ID.String(),
			OtherUserName:              row.OtherUserName,
			OtherUserProfilePictureURL: picURL,
			IsLocked:                   row.IsLocked,
			CreatedAt:                  row.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, response)
}
