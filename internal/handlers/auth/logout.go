package auth

import (
	"net/http"

	"babycare-api/internal/models"

	"github.com/gin-gonic/gin"
)

// Logout instructs the client to clear the session token. Token removal is handled client-side.
func (h *AuthHandler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, models.SuccessResponse{Message: "logged out successfully"})
}
