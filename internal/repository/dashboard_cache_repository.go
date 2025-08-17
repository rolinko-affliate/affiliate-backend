package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/affiliate-backend/internal/domain"
	"github.com/affiliate-backend/internal/platform/logger"
)

// DashboardCacheRepository defines the interface for dashboard caching operations
type DashboardCacheRepository interface {
	// Dashboard-specific cache methods
	GetCachedDashboardData(ctx context.Context, orgID int64, orgType domain.OrganizationType) (*domain.DashboardData, error)
	SetCachedDashboardData(ctx context.Context, orgID int64, orgType domain.OrganizationType, data *domain.DashboardData, ttl time.Duration) error
	InvalidateDashboardCache(ctx context.Context, orgID int64) error
	
	// Generic cache operations
	SetCache(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	GetCache(ctx context.Context, key string, dest interface{}) error
	DeleteCache(ctx context.Context, key string) error
}

// dashboardCacheRepository implements DashboardCacheRepository
type dashboardCacheRepository struct {
	cache  *redis.Client
	logger *logger.Logger
}

// NewDashboardCacheRepository creates a new dashboard cache repository
func NewDashboardCacheRepository(cache *redis.Client) DashboardCacheRepository {
	return &dashboardCacheRepository{
		cache:  cache,
		logger: logger.GetDefault(),
	}
}

// Dashboard-specific cache methods
func (r *dashboardCacheRepository) GetCachedDashboardData(ctx context.Context, orgID int64, orgType domain.OrganizationType) (*domain.DashboardData, error) {
	// TODO: Re-enable Redis caching - currently disabled
	if r.cache == nil {
		return nil, fmt.Errorf("cache not available")
	}
	
	key := fmt.Sprintf("dashboard:%s:org:%d", orgType, orgID)
	val, err := r.cache.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	
	var dashboard domain.DashboardData
	err = json.Unmarshal([]byte(val), &dashboard)
	if err != nil {
		return nil, err
	}
	
	return &dashboard, nil
}

func (r *dashboardCacheRepository) SetCachedDashboardData(ctx context.Context, orgID int64, orgType domain.OrganizationType, data *domain.DashboardData, ttl time.Duration) error {
	// TODO: Re-enable Redis caching - currently disabled
	if r.cache == nil {
		return nil // No-op if cache not available
	}
	
	key := fmt.Sprintf("dashboard:%s:org:%d", orgType, orgID)
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	
	return r.cache.Set(ctx, key, jsonData, ttl).Err()
}

func (r *dashboardCacheRepository) InvalidateDashboardCache(ctx context.Context, orgID int64) error {
	// TODO: Re-enable Redis caching - currently disabled
	if r.cache == nil {
		return nil // No-op if cache not available
	}
	
	pattern := fmt.Sprintf("dashboard:*:org:%d", orgID)
	keys, err := r.cache.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}
	
	if len(keys) > 0 {
		return r.cache.Del(ctx, keys...).Err()
	}
	
	return nil
}

// Generic cache operations
func (r *dashboardCacheRepository) SetCache(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	// TODO: Re-enable Redis caching - currently disabled
	if r.cache == nil {
		return nil // No-op if cache not available
	}
	
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal cache value: %w", err)
	}

	return r.cache.Set(ctx, key, data, ttl).Err()
}

func (r *dashboardCacheRepository) GetCache(ctx context.Context, key string, dest interface{}) error {
	// TODO: Re-enable Redis caching - currently disabled
	if r.cache == nil {
		return fmt.Errorf("cache not available")
	}
	
	data, err := r.cache.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(data), dest)
}

func (r *dashboardCacheRepository) DeleteCache(ctx context.Context, key string) error {
	// TODO: Re-enable Redis caching - currently disabled
	if r.cache == nil {
		return nil // No-op if cache not available
	}
	
	return r.cache.Del(ctx, key).Err()
}