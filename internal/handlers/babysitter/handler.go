package babysitter

import (
	"database/sql"

	"babycare-api/internal/config"
	services_cache "babycare-api/internal/services/cache"
	services_storage "babycare-api/internal/services/storage"
)

// BabysitterHandler holds the dependencies needed by all babysitter endpoints.
type BabysitterHandler struct {
	db             *sql.DB
	storageService *services_storage.CloudinaryService
	cacheService   *services_cache.CacheService // may be nil if Redis is unavailable
	cfg            *config.Config
}

// NewBabysitterHandler constructs a BabysitterHandler.
// cacheService may be nil — all cache operations are guarded against nil.
func NewBabysitterHandler(
	db *sql.DB,
	storageService *services_storage.CloudinaryService,
	cacheService *services_cache.CacheService,
	cfg *config.Config,
) *BabysitterHandler {
	return &BabysitterHandler{
		db:             db,
		storageService: storageService,
		cacheService:   cacheService,
		cfg:            cfg,
	}
}
