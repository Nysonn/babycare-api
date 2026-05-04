package auth

import (
	"database/sql"
	"log"
	"net/http"
	"strings"

	"babycare-api/internal/db"
	"babycare-api/internal/models"

	"github.com/gin-gonic/gin"
)

// ForgotPassword accepts an email and triggers Clerk's native password-reset email
// for the matching account. Always returns 200 to prevent email enumeration.
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "valid email is required"})
		return
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	queries := db.New(h.db)
	ctx := c.Request.Context()

	user, err := queries.GetUserByEmail(ctx, req.Email)
	if err == sql.ErrNoRows {
		// Return 200 regardless to prevent email enumeration.
		c.JSON(http.StatusOK, gin.H{"message": "If an account with that email exists, a reset link has been sent."})
		return
	}
	if err != nil {
		log.Printf("forgot_password: get user by email: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
		return
	}

	// Only users authenticated via Clerk can use this flow.
	if !user.ClerkUserID.Valid || user.ClerkUserID.String == "" {
		c.JSON(http.StatusOK, gin.H{"message": "If an account with that email exists, a reset link has been sent."})
		return
	}

	if err := h.clerkService.SendResetPasswordEmail(user.ClerkUserID.String); err != nil {
		log.Printf("forgot_password: clerk send reset email: %v", err)
		// Still return 200 — the caller should not know whether the operation failed.
		c.JSON(http.StatusOK, gin.H{"message": "If an account with that email exists, a reset link has been sent."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "If an account with that email exists, a reset link has been sent."})
}
