package auth

import (
	"net/http"

	"babycare-api/internal/models"

	"github.com/gin-gonic/gin"
)

// ForgotPassword has been moved to the Clerk client SDK.
// Clerk's reset_password_email_code flow must be initiated from the client,
// not from this API.
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	c.JSON(http.StatusGone, models.ErrorResponse{
		Error: "forgot password must be initiated from the Clerk client SDK using the reset_password_email_code flow",
	})
}
