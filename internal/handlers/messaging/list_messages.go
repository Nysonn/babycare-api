package messaging

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"babycare-api/internal/db"
	"babycare-api/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const cacheTTLMessages = 2 * time.Minute

// ListMessages returns all messages in a conversation and marks them as read.
// Results are cached in Redis for 2 minutes.
func (h *MessagingHandler) ListMessages(c *gin.Context) {
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

	convIDStr := c.Param("conversation_id")
	if convIDStr == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "conversation_id is required"})
		return
	}
	convID, err := uuid.Parse(convIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid conversation_id"})
		return
	}

	queries := db.New(h.db)
	ctx := c.Request.Context()

	conv, err := queries.GetConversationByID(ctx, convID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "conversation not found"})
		return
	}
	if err != nil {
		log.Printf("messaging list_messages: get conversation: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
		return
	}

	// Verify the current user is a participant.
	if currentUser.ID != conv.ParentID && currentUser.ID != conv.BabysitterID {
		c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "you are not a participant in this conversation"})
		return
	}

	cacheKey := fmt.Sprintf("messages:%s", convIDStr)

	// Try the cache first.
	if h.cacheService != nil {
		if cached, err := h.cacheService.Get(ctx, cacheKey); err == nil {
			// Still mark messages as read even when serving from cache.
			markMessagesRead(h.db, convID, currentUser.ID)
			c.Data(http.StatusOK, "application/json; charset=utf-8", []byte(cached))
			return
		}
	}

	messages, err := queries.ListMessagesByConversation(ctx, convID)
	if err != nil {
		log.Printf("messaging list_messages: db query: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
		return
	}

	response := make([]models.MessageResponse, 0, len(messages))
	for _, m := range messages {
		response = append(response, models.MessageResponse{
			ID:             m.ID.String(),
			ConversationID: m.ConversationID.String(),
			SenderID:       m.SenderID.String(),
			Content:        m.Content,
			IsRead:         m.IsRead,
			SentAt:         m.SentAt,
		})
	}

	// Populate the cache (failure is non-fatal).
	if h.cacheService != nil {
		if err := h.cacheService.Set(ctx, cacheKey, response, cacheTTLMessages); err != nil {
			log.Printf("messaging list_messages: cache set: %v", err)
		}
	}

	// Mark messages from the other participant as read.
	markMessagesRead(h.db, convID, currentUser.ID)

	c.JSON(http.StatusOK, response)
}

// markMessagesRead marks all unread messages in the conversation sent by the other
// participant as read. Runs synchronously but errors are only logged, never surfaced.
func markMessagesRead(database *sql.DB, convID, userID uuid.UUID) {
	queries := db.New(database)
	if err := queries.MarkMessagesAsRead(context.Background(), db.MarkMessagesAsReadParams{
		ConversationID: convID,
		SenderID:       userID,
	}); err != nil {
		log.Printf("messaging mark_read: %v", err)
	}
}
