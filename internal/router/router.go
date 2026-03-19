package router

import (
	"database/sql"
	"net/http"

	"babycare-api/internal/config"
	handlers_admin "babycare-api/internal/handlers/admin"
	handlers_auth "babycare-api/internal/handlers/auth"
	handlers_babysitter "babycare-api/internal/handlers/babysitter"
	handlers_messaging "babycare-api/internal/handlers/messaging"
	handlers_parent "babycare-api/internal/handlers/parent"
	"babycare-api/internal/middleware"
	services_auth "babycare-api/internal/services/auth"
	services_cache "babycare-api/internal/services/cache"
	services_email "babycare-api/internal/services/email"
	services_messaging "babycare-api/internal/services/messaging"
	"babycare-api/internal/services/storage"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Setup configures and returns a Gin engine with all middleware and routes.
func Setup(
	db *sql.DB,
	cfg *config.Config,
	clerkService *services_auth.ClerkService,
	storageService *storage.CloudinaryService,
	emailService *services_email.EmailService,
	cacheService *services_cache.CacheService, // may be nil if Redis is unavailable
	streamService *services_messaging.StreamService,
) *gin.Engine {
	// Switch Gin to release mode in production to suppress debug output.
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// CORS — allow all origins for Flutter web and mobile clients.
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{
		"http://localhost:5173",
		"http://localhost:5174",
		"https://babycare-f6f8e.web.app",
	}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Accept"}
	r.Use(cors.New(corsConfig))

	// Global middleware.
	r.Use(gin.Recovery()) // Recover from panics and return 500
	r.Use(gin.Logger())   // Request/response logging

	// Health check — used by Docker, load balancers, and uptime monitors.
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "babycare-api",
		})
	})

	// Initialise handlers.
	authHandler := handlers_auth.NewAuthHandler(db, clerkService, storageService, streamService, cfg)
	adminHandler := handlers_admin.NewAdminHandler(db, cfg, emailService)
	babysitterHandler := handlers_babysitter.NewBabysitterHandler(db, storageService, cacheService, cfg)
	parentHandler := handlers_parent.NewParentHandler(db, cfg)
	messagingHandler := handlers_messaging.NewMessagingHandler(db, streamService, emailService, cacheService, cfg)

	// Shared middleware factories.
	requireAuth := middleware.RequireAuth(clerkService)

	// Versioned API group.
	api := r.Group("/api/v1")

	// --- Auth routes (public) ---
	auth := api.Group("/auth")
	{
		auth.POST("/register/parent", authHandler.RegisterParent)
		auth.POST("/register/babysitter", authHandler.RegisterBabysitter)
		auth.POST("/login", authHandler.Login)
		auth.POST("/logout", requireAuth, authHandler.Logout)
	}

	// --- Admin routes (admin role required) ---
	admin := api.Group("/admin", requireAuth, middleware.RequireRole(db, "admin"))
	{
		admin.GET("/users", adminHandler.ListUsers)
		admin.GET("/users/:id", adminHandler.GetUser)
		admin.PUT("/babysitters/:id/approve", adminHandler.ApproveBabysitter)
		admin.PUT("/users/:id/suspend", adminHandler.SuspendUser)
		admin.DELETE("/users/:id", adminHandler.DeleteUser)
		admin.POST("/create", adminHandler.CreateAdmin)
		admin.GET("/activity", adminHandler.GetActivity)
	}

	// --- Babysitter routes ---

	// Public: browse approved babysitters.
	api.GET("/babysitters", babysitterHandler.ListBabysitters)

	// Babysitter-only: manage own profile.
	// Static paths registered before the :id wildcard group.
	babysitterSelf := api.Group("/babysitters", requireAuth, middleware.RequireRole(db, "babysitter"))
	{
		babysitterSelf.PUT("/profile", babysitterHandler.UpdateProfile)
		babysitterSelf.GET("/profile/views", babysitterHandler.GetProfileViews)
	}

	// Parents and babysitters can view a babysitter's profile.
	babysitterView := api.Group("/babysitters", requireAuth, middleware.RequireRole(db, "parent", "babysitter"))
	{
		babysitterView.GET("/:id", babysitterHandler.GetBabysitter)
	}

	// --- Parent routes (parent role required) ---
	parentRoutes := api.Group("/parents", requireAuth, middleware.RequireRole(db, "parent"))
	{
		parentRoutes.GET("/profile", parentHandler.GetProfile)
		parentRoutes.PUT("/profile", parentHandler.UpdateProfile)
	}

	// --- Messaging routes (parent or babysitter required) ---
	conversations := api.Group("/conversations", requireAuth, middleware.RequireRole(db, "parent", "babysitter"))
	{
		conversations.POST("", messagingHandler.StartConversation)
		conversations.GET("", messagingHandler.ListConversations)
		conversations.GET("/:conversation_id/messages", messagingHandler.ListMessages)
		conversations.POST("/:conversation_id/messages", messagingHandler.SendMessage)
	}

	return r
}
