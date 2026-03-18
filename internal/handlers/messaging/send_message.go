package messaging

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"babycare-api/internal/db"
	"babycare-api/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SendMessage sends a message in the given conversation.
func (h *MessagingHandler) SendMessage(c *gin.Context) {
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
		log.Printf("messaging send_message: get conversation: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
		return
	}

	// Verify the current user is a participant.
	if currentUser.ID != conv.ParentID && currentUser.ID != conv.BabysitterID {
		c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "you are not a participant in this conversation"})
		return
	}

	if conv.IsLocked {
		c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "this conversation has been locked"})
		return
	}

	var req models.SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	// Send via Stream Chat.
	if err := h.streamService.SendMessage(
		conv.StreamChannelID.String,
		currentUser.ID.String(),
		req.Content,
	); err != nil {
		log.Printf("messaging send_message: stream send: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "failed to send message"})
		return
	}

	// Persist in our DB.
	msg, err := queries.CreateMessage(ctx, db.CreateMessageParams{
		ConversationID: convID,
		SenderID:       currentUser.ID,
		Content:        req.Content,
	})
	if err != nil {
		log.Printf("messaging send_message: db create: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "failed to save message"})
		return
	}

	// Invalidate the messages cache for this conversation.
	if h.cacheService != nil {
		cacheKey := fmt.Sprintf("messages:%s", convIDStr)
		if err := h.cacheService.Delete(ctx, cacheKey); err != nil {
			log.Printf("messaging send_message: cache invalidate: %v", err)
		}
	}

	// Determine the recipient (the other participant).
	recipientID := conv.BabysitterID
	if currentUser.ID == conv.BabysitterID {
		recipientID = conv.ParentID
	}

	// Send an unread notification email asynchronously.
	go h.sendUnreadNotification(
		recipientID,
		currentUser.FullName,
		conv.StreamChannelID.String,
	)

	c.JSON(http.StatusCreated, models.MessageResponse{
		ID:             msg.ID.String(),
		ConversationID: msg.ConversationID.String(),
		SenderID:       msg.SenderID.String(),
		Content:        msg.Content,
		IsRead:         msg.IsRead,
		SentAt:         msg.SentAt,
	})
}

// sendUnreadNotification emails the recipient if this is their first unread message
// in the channel (unread count == 1), to avoid notification spam.
func (h *MessagingHandler) sendUnreadNotification(
	recipientID uuid.UUID,
	senderFullName string,
	channelID string,
) {
	queries := db.New(h.db)
	ctx := context.Background()

	recipient, err := queries.GetUserByID(ctx, recipientID)
	if err != nil {
		log.Printf("messaging unread_notify: get recipient %s: %v", recipientID, err)
		return
	}

	unread, err := h.streamService.GetUnreadCount(recipientID.String(), channelID)
	if err != nil {
		log.Printf("messaging unread_notify: get unread count: %v", err)
		return
	}

	// Only send on the very first unread message to avoid spamming.
	if unread != 1 {
		return
	}

	body := fmt.Sprintf(
		"Hello %s, you have a new unread message from %s on BabyCare. "+
			"Please log in to read it.",
		recipient.FullName,
		senderFullName,
	)
	if err := h.emailService.SendEmail(
		recipient.Email,
		recipient.FullName,
		"You have a new message on BabyCare",
		body,
	); err != nil {
		log.Printf("messaging unread_notify: send email to %s: %v", recipient.Email, err)
	}
}
