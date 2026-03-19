package messaging

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"babycare-api/internal/db"
	"babycare-api/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// buildChannelID returns a deterministic channel ID under 64 characters.
// The two IDs are concatenated in a fixed order so the same pair always
// produces the same hash regardless of who initiated the conversation.
func buildChannelID(parentID, babysitterID string) string {
	combined := parentID + babysitterID
	hash := md5.Sum([]byte(combined))
	return fmt.Sprintf("bc-%x", hash)
}

type startConversationRequest struct {
	BabysitterID string `json:"babysitter_id" binding:"required"`
}

// StartConversation opens a new conversation between the authenticated parent and a babysitter.
// Only parents may initiate conversations.
func (h *MessagingHandler) StartConversation(c *gin.Context) {
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

	if currentUser.Role != db.UserRoleParent {
		c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "only parents can start conversations"})
		return
	}

	var req startConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	babysitterID, err := uuid.Parse(req.BabysitterID)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid babysitter_id"})
		return
	}

	queries := db.New(h.db)
	ctx := c.Request.Context()

	// Verify the target user exists and is a babysitter.
	babysitter, err := queries.GetUserByID(ctx, babysitterID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "babysitter not found"})
		return
	}
	if err != nil {
		log.Printf("messaging start_conversation: get babysitter: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
		return
	}
	if babysitter.Role != db.UserRoleBabysitter {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "target user is not a babysitter"})
		return
	}

	// Return existing conversation if one already exists (idempotent).
	existing, err := queries.GetConversationByParticipants(ctx, db.GetConversationByParticipantsParams{
		ParentID:     currentUser.ID,
		BabysitterID: babysitterID,
	})
	if err == nil {
		c.JSON(http.StatusOK, models.ConversationResponse{
			ID:        existing.ID.String(),
			IsLocked:  existing.IsLocked,
			CreatedAt: existing.CreatedAt,
		})
		return
	}
	if err != sql.ErrNoRows {
		log.Printf("messaging start_conversation: check existing: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
		return
	}

	// Ensure both users exist in Stream Chat before creating the channel.
	if err := h.streamService.UpsertUser(currentUser.ID.String(), currentUser.FullName); err != nil {
		log.Printf("start_conversation: upsert parent to stream: %v", err)
	}
	if err := h.streamService.UpsertUser(babysitterID.String(), babysitter.FullName); err != nil {
		log.Printf("start_conversation: upsert babysitter to stream: %v", err)
	}

	// Build a deterministic Stream channel ID under 64 characters.
	channelID := buildChannelID(currentUser.ID.String(), babysitterID.String())

	// Create the Stream Chat channel.
	if _, err := h.streamService.CreateChannel(channelID, currentUser.ID.String(), babysitterID.String()); err != nil {
		log.Printf("messaging start_conversation: create stream channel: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "failed to create chat channel"})
		return
	}

	// Persist the conversation record.
	conv, err := queries.CreateConversation(ctx, db.CreateConversationParams{
		ParentID:        currentUser.ID,
		BabysitterID:    babysitterID,
		StreamChannelID: sql.NullString{String: channelID, Valid: true},
	})
	if err != nil {
		log.Printf("messaging start_conversation: db create: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "failed to save conversation"})
		return
	}

	c.JSON(http.StatusCreated, models.ConversationResponse{
		ID:            conv.ID.String(),
		OtherUserName: babysitter.FullName,
		IsLocked:      conv.IsLocked,
		CreatedAt:     conv.CreatedAt,
	})
}
