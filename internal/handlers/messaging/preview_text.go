package messaging

import (
	"log"
	"net/http"

	"babycare-api/internal/db"
	"babycare-api/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetConversationPreviews returns the authenticated user's conversations enriched
// with the most recent message in each thread. Conversations that have no messages
// yet are omitted from the response.
//
// The preview_text field is suitable for displaying directly in a conversation
// list (e.g. "Hey, are you available this weekend?"). Clients should use
// is_read together with last_sender_id to decide whether to show an unread
// indicator: a message is unread when is_read == false AND
// last_sender_id != <current user's ID>.
func (h *MessagingHandler) GetConversationPreviews(c *gin.Context) {
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
	ctx := c.Request.Context()

	// Fetch all conversations the user participates in (already ordered by updated_at DESC).
	conversations, err := queries.ListConversationsForUser(ctx, currentUser.ID)
	if err != nil {
		log.Printf("messaging get_conversation_previews: list conversations: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
		return
	}

	// Fetch the last message for every conversation this user is in.
	lastMessages, err := queries.GetLastMessagePerConversation(ctx, currentUser.ID)
	if err != nil {
		log.Printf("messaging get_conversation_previews: get last messages: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
		return
	}

	// Build a lookup map: conversation_id → last message.
	lastMsgByConv := make(map[uuid.UUID]db.Message, len(lastMessages))
	for _, m := range lastMessages {
		lastMsgByConv[m.ConversationID] = m
	}

	response := make([]models.ConversationPreviewResponse, 0, len(conversations))
	for _, conv := range conversations {
		msg, hasMsg := lastMsgByConv[conv.ID]
		if !hasMsg {
			// Conversation has no messages yet — skip.
			continue
		}

		picURL, _ := conv.OtherUserProfilePictureUrl.(string)
		response = append(response, models.ConversationPreviewResponse{
			ConversationID:             conv.ID.String(),
			OtherUserName:              conv.OtherUserName,
			OtherUserProfilePictureURL: picURL,
			IsLocked:                   conv.IsLocked,
			LastMessageID:              msg.ID.String(),
			LastSenderID:               msg.SenderID.String(),
			PreviewText:                msg.Content,
			IsRead:                     msg.IsRead,
			LastMessageSent:            msg.SentAt,
		})
	}

	c.JSON(http.StatusOK, response)
}
