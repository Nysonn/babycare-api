package admin

import (
	"database/sql"
	"log"
	"net/http"

	"babycare-api/internal/db"
	"babycare-api/internal/models"
	services_auth "babycare-api/internal/services/auth"

	"github.com/gin-gonic/gin"
)

// CreateAdmin creates a new admin user account.
func (h *AdminHandler) CreateAdmin(c *gin.Context) {
	var req models.CreateAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	queries := db.New(h.db)
	ctx := c.Request.Context()

	// Check for duplicate email.
	_, err := queries.GetUserByEmail(ctx, req.Email)
	if err == nil {
		c.JSON(http.StatusConflict, models.ErrorResponse{Error: "email already registered"})
		return
	}
	if err != sql.ErrNoRows {
		log.Printf("admin create_admin: check email: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
		return
	}

	hash, err := services_auth.HashPassword(req.Password)
	if err != nil {
		log.Printf("admin create_admin: hash password: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
		return
	}

	user, err := queries.CreateAdminUser(ctx, db.CreateAdminUserParams{
		FullName:     req.FullName,
		Email:        req.Email,
		PasswordHash: sql.NullString{String: hash, Valid: true},
	})
	if err != nil {
		log.Printf("admin create_admin: db insert: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "failed to create admin"})
		return
	}

	c.JSON(http.StatusCreated, models.UserResponse{
		ID:        user.ID.String(),
		FullName:  user.FullName,
		Email:     user.Email,
		Phone:     user.Phone.String,
		Role:      string(user.Role),
		Status:    string(user.Status),
		CreatedAt: user.CreatedAt,
	})
}
