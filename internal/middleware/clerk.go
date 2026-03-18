package middleware

import (
	"database/sql"
	"net/http"
	"strings"

	"babycare-api/internal/db"
	services_auth "babycare-api/internal/services/auth"

	"github.com/gin-gonic/gin"
)

// RequireAuth validates the Bearer token in the Authorization header via Clerk.
// On success it sets "clerk_user_id" in the Gin context for downstream handlers.
func RequireAuth(clerkService *services_auth.ClerkService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorised"})
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		clerkUserID, err := clerkService.VerifyToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		c.Set("clerk_user_id", clerkUserID)
		c.Next()
	}
}

// RequireRole checks that the authenticated user exists, is active, and holds one
// of the permitted roles. It must be chained after RequireAuth.
// On success it sets "current_user" (db.User) in the Gin context.
func RequireRole(database *sql.DB, roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		clerkUserID, exists := c.Get("clerk_user_id")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorised"})
			return
		}

		queries := db.New(database)
		user, err := queries.GetUserByClerkID(c.Request.Context(), sql.NullString{
			String: clerkUserID.(string),
			Valid:  true,
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorised"})
			return
		}

		// Deleted users are filtered by the query (deleted_at IS NULL), so a
		// not-found error already covers that case above.

		if user.Status == db.UserStatusSuspended {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "account suspended"})
			return
		}

		allowed := false
		for _, role := range roles {
			if string(user.Role) == role {
				allowed = true
				break
			}
		}
		if !allowed {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}

		c.Set("current_user", user)
		c.Next()
	}
}
