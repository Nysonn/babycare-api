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

// ApproveBabysitter marks a babysitter's profile as approved and notifies them by email.
func (h *AdminHandler) ApproveBabysitter(c *gin.Context) {
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
		log.Printf("admin approve_babysitter: get user: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
		return
	}

	if user.Role != db.UserRoleBabysitter {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "user is not a babysitter"})
		return
	}

	if _, err := queries.ApproveBabysitter(ctx, userID); err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "babysitter profile not found"})
		return
	} else if err != nil {
		log.Printf("admin approve_babysitter: approve query: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
		return
	}

	// Send approval notification email (non-fatal if it fails).
	body := fmt.Sprintf(
		"Hello %s, your babysitter account on BabyCare has been approved. "+
			"You can now log in and start receiving messages from parents.",
		user.FullName,
	)
	if err := h.emailService.SendEmail(
		user.Email,
		user.FullName,
		"Your BabyCare account has been approved",
		body,
	); err != nil {
		log.Printf("admin approve_babysitter: send email to %s: %v", user.Email, err)
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Message: "babysitter approved successfully"})
}
