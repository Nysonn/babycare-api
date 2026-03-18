package admin

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"babycare-api/internal/db"
	"babycare-api/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SuspendUser sets a user's status to suspended and locks all their conversations.
func (h *AdminHandler) SuspendUser(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "id is required"})
		return
	}

	userID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid user id"})
		return
	}

	queries := db.New(h.db)
	ctx := c.Request.Context()

	user, err := queries.GetUserByID(ctx, userID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "user not found"})
		return
	}
	if err != nil {
		log.Printf("admin suspend_user: get user: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
		return
	}

	if user.Role == db.UserRoleAdmin {
		c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "cannot suspend an admin account"})
		return
	}

	if _, err := queries.UpdateUserStatus(ctx, db.UpdateUserStatusParams{
		ID:     userID,
		Status: db.UserStatusSuspended,
	}); err != nil {
		log.Printf("admin suspend_user: update status: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
		return
	}

	// Lock all conversations this user is part of.
	if err := queries.LockConversationsByUser(ctx, userID); err != nil {
		log.Printf("admin suspend_user: lock conversations: %v", err)
	}

	// Send suspension notification email (non-fatal if it fails).
	body := fmt.Sprintf(
		"Hello %s, your account on BabyCare has been suspended. "+
			"If you believe this is a mistake please contact our support team.",
		user.FullName,
	)
	if err := h.emailService.SendEmail(
		user.Email,
		user.FullName,
		"Your BabyCare account has been suspended",
		body,
	); err != nil {
		log.Printf("admin suspend_user: send email to %s: %v", user.Email, err)
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Message: "user suspended successfully"})
}
