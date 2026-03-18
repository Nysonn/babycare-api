package auth

import (
	"log"
	"net/http"
	"strings"

	"babycare-api/internal/models"

	"github.com/gin-gonic/gin"
)

// Logout revokes the caller's Clerk session. Best-effort — always returns 200.
func (h *AuthHandler) Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")

	if token != "" {
		if err := h.clerkService.RevokeSession(token); err != nil {
			// Non-fatal: log and continue so the client-side session is still cleared.
			log.Printf("logout: revoke clerk session: %v", err)
		}
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Message: "logged out successfully"})
}
