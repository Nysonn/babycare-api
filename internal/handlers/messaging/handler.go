package messaging

import (
	"database/sql"

	"babycare-api/internal/config"
	services_cache "babycare-api/internal/services/cache"
	services_email "babycare-api/internal/services/email"
	services_messaging "babycare-api/internal/services/messaging"
)

// MessagingHandler holds dependencies for all conversation and message endpoints.
type MessagingHandler struct {
	db            *sql.DB
	streamService *services_messaging.StreamService
	emailService  *services_email.EmailService
	cacheService  *services_cache.CacheService // may be nil if Redis is unavailable
	cfg           *config.Config
}

// NewMessagingHandler constructs a MessagingHandler with all required dependencies.
func NewMessagingHandler(
	db *sql.DB,
	streamService *services_messaging.StreamService,
	emailService *services_email.EmailService,
	cacheService *services_cache.CacheService,
	cfg *config.Config,
) *MessagingHandler {
	return &MessagingHandler{
		db:            db,
		streamService: streamService,
		emailService:  emailService,
		cacheService:  cacheService,
		cfg:           cfg,
	}
}
