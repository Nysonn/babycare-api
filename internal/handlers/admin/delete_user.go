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

// DeleteUser soft-deletes a user account and notifies them by email.
func (h *AdminHandler) DeleteUser(c *gin.Context) {
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
		log.Printf("admin delete_user: get user: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
		return
	}

	if user.Role == db.UserRoleAdmin {
		c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "cannot delete an admin account"})
		return
	}

	if err := queries.SoftDeleteUser(ctx, userID); err != nil {
		log.Printf("admin delete_user: soft delete: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
		return
	}

	// Send deletion notification email (non-fatal if it fails).
	body := fmt.Sprintf(
		"Hello %s, your account on BabyCare has been removed. "+
			"If you believe this is a mistake please contact our support team.",
		user.FullName,
	)
	if err := h.emailService.SendEmail(
		user.Email,
		user.FullName,
		"Your BabyCare account has been removed",
		body,
	); err != nil {
		log.Printf("admin delete_user: send email to %s: %v", user.Email, err)
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Message: "user deleted successfully"})
}
