package auth

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"babycare-api/internal/db"
	"babycare-api/internal/models"
	services_auth "babycare-api/internal/services/auth"

	"github.com/gin-gonic/gin"
)

// Login authenticates a user and returns a session token.
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	queries := db.New(h.db)
	ctx := c.Request.Context()

	// Fetch user by email.
	user, err := queries.GetUserByEmail(ctx, req.Email)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "invalid credentials"})
		return
	}
	if err != nil {
		log.Printf("login: get user by email: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
		return
	}

	// Babysitters must be approved before they can log in.
	if user.Role == db.UserRoleBabysitter {
		profile, err := queries.GetBabysitterProfileByUserID(ctx, user.ID)
		if err != nil && err != sql.ErrNoRows {
			log.Printf("login: get babysitter profile: %v", err)
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
			return
		}
		if err == sql.ErrNoRows || !profile.IsApproved {
			c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "your account is pending admin approval"})
			return
		}
	}

	// Check account status.
	if user.Status == db.UserStatusSuspended {
		c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "your account has been suspended"})
		return
	}

	// Verify password.
	if !services_auth.CheckPassword(req.Password, user.PasswordHash.String) {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "invalid credentials"})
		return
	}

	expiresAt := time.Now().Add(90 * 24 * time.Hour)

	var tokenSubject string
	if user.Role == db.UserRoleAdmin {
		// Admin is seeded directly into the DB with no Clerk account.
		// Use the user's own UUID as the JWT subject.
		tokenSubject = user.ID.String()
	} else {
		// Guard against non-admin users with no Clerk ID (misconfigured accounts).
		if !user.ClerkUserID.Valid || user.ClerkUserID.String == "" {
			log.Printf("login: user %s has empty clerk_user_id", user.ID)
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "account configuration error, please contact support"})
			return
		}

		// Sync user into Stream Chat (non-fatal if it fails).
		if err := h.streamService.UpsertUser(user.ID.String(), user.FullName); err != nil {
			log.Printf("login: stream upsert user: %v", err)
		}

		tokenSubject = user.ClerkUserID.String
	}

	// Generate a signed JWT session token.
	token, err := h.clerkService.GenerateToken(tokenSubject, expiresAt)
	if err != nil {
		log.Printf("login: generate token: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "failed to create session"})
		return
	}

	c.JSON(http.StatusOK, models.LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User: models.UserResponse{
			ID:        user.ID.String(),
			FullName:  user.FullName,
			Email:     user.Email,
			Phone:     user.Phone.String,
			Role:      string(user.Role),
			Status:    string(user.Status),
			CreatedAt: user.CreatedAt,
		},
	})
}
