package auth

import (
	"database/sql"
	"log"
	"net/http"

	"babycare-api/internal/db"
	"babycare-api/internal/models"
	services_auth "babycare-api/internal/services/auth"

	"github.com/gin-gonic/gin"
)

// RegisterParent creates a new parent account with profile.
func (h *AuthHandler) RegisterParent(c *gin.Context) {
	var req models.RegisterParentRequest
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
		log.Printf("register_parent: check email: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
		return
	}

	// Create user in Clerk.
	clerkUserID, err := h.clerkService.CreateUser(req.Email, req.Password, req.FullName)
	if err != nil {
		log.Printf("register_parent: clerk create user: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "failed to create account"})
		return
	}

	// Hash password for local storage.
	hash, err := services_auth.HashPassword(req.Password)
	if err != nil {
		log.Printf("register_parent: hash password: %v", err)
		_ = h.clerkService.DeleteUser(clerkUserID) // rollback Clerk
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
		return
	}

	// Persist the user row.
	user, err := queries.CreateUser(ctx, db.CreateUserParams{
		FullName:     req.FullName,
		Email:        req.Email,
		Phone:        sql.NullString{String: req.Phone, Valid: req.Phone != ""},
		Role:         db.UserRoleParent,
		PasswordHash: sql.NullString{String: hash, Valid: true},
		ClerkUserID:  sql.NullString{String: clerkUserID, Valid: true},
	})
	if err != nil {
		log.Printf("register_parent: db create user: %v", err)
		_ = h.clerkService.DeleteUser(clerkUserID) // rollback Clerk
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "failed to save user"})
		return
	}

	// Create the parent profile.
	_, err = queries.CreateParentProfile(ctx, db.CreateParentProfileParams{
		UserID:          user.ID,
		Location:        sql.NullString{String: req.Location, Valid: req.Location != ""},
		Occupation:      sql.NullString{String: req.Occupation, Valid: req.Occupation != ""},
		PreferredHours:  sql.NullString{String: req.PreferredHours, Valid: req.PreferredHours != ""},
		PrimaryLocation: sql.NullString{String: req.PrimaryLocation, Valid: req.PrimaryLocation != ""},
	})
	if err != nil {
		log.Printf("register_parent: create parent profile: %v", err)
		// User row exists but profile failed — non-fatal; client can retry profile update.
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
