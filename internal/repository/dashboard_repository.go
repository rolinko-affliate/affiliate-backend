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

// DashboardRepository defines the interface for dashboard caching operations
type DashboardRepository interface {
	// Cache methods
	GetCachedDashboardData(ctx context.Context, orgID int64, orgType domain.OrganizationType) (*domain.DashboardData, error)
	SetCachedDashboardData(ctx context.Context, orgID int64, orgType domain.OrganizationType, data *domain.DashboardData, ttl time.Duration) error
	InvalidateDashboardCache(ctx context.Context, orgID int64) error
	
	// Generic cache operations
	SetCache(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	GetCache(ctx context.Context, key string, dest interface{}) error
	DeleteCache(ctx context.Context, key string) error
}

// dashboardRepository implements DashboardRepository
type dashboardRepository struct {
	cache  *redis.Client
	logger *logger.Logger
}

// NewDashboardRepository creates a new dashboard repository
func NewDashboardRepository(cache *redis.Client) DashboardRepository {
	return &dashboardRepository{
		cache:  cache,
		logger: logger.GetDefault(),
	}
}

// Cache methods implementation

// Cache methods
func (r *dashboardRepository) GetCachedDashboardData(ctx context.Context, orgID int64, orgType domain.OrganizationType) (*domain.DashboardData, error) {
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

func (r *dashboardRepository) SetCachedDashboardData(ctx context.Context, orgID int64, orgType domain.OrganizationType, data *domain.DashboardData, ttl time.Duration) error {
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

func (r *dashboardRepository) InvalidateDashboardCache(ctx context.Context, orgID int64) error {
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
func (r *dashboardRepository) SetCache(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	if r.cache == nil {
		return fmt.Errorf("cache not available")
	}
	
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}
	
	return r.cache.Set(ctx, key, jsonData, ttl).Err()
}

func (r *dashboardRepository) GetCache(ctx context.Context, key string, dest interface{}) error {
	if r.cache == nil {
		return fmt.Errorf("cache not available")
	}
	
	data, err := r.cache.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	
	return json.Unmarshal([]byte(data), dest)
}

// Data retrieval methods
func (r *dashboardRepository) GetCampaignsByOrganization(ctx context.Context, orgID int64) ([]domain.CampaignInfo, error) {
	// Return empty slice for now
	return []domain.CampaignInfo{}, nil
}

func (r *dashboardRepository) GetCampaignDetail(ctx context.Context, campaignID int64) (*domain.CampaignDetail, error) {
	// Return basic campaign detail
	return &domain.CampaignDetail{
		Campaign: domain.CampaignInfo{
			ID:     campaignID,
			Name:   "Sample Campaign",
			Status: "active",
		},
		Performance: domain.CampaignMetrics{
			Clicks:      0,
			Conversions: 0,
			Revenue:     0,
			Impressions: 0,
		},
		DailyStats: []domain.DailyStat{},
	}, nil
}

func (r *dashboardRepository) GetBillingInfo(ctx context.Context, orgID int64) (*domain.BillingInfo, error) {
	// Return basic billing info
	return &domain.BillingInfo{
		CurrentBalance: 0,
		MonthlySpend:   0,
		Currency:       "USD",
	}, nil
}

func (r *dashboardRepository) GetOrganizationClients(ctx context.Context, orgID int64) ([]domain.ClientInfo, error) {
	// Return empty slice for now
	return []domain.ClientInfo{}, nil
}

func (r *dashboardRepository) GetDashboardMetrics(ctx context.Context, orgID int64, metricType string, from, to time.Time) ([]domain.DashboardMetric, error) {
	// Return empty metrics for now
	return []domain.DashboardMetric{}, nil
}

// Platform-specific methods
func (r *dashboardRepository) GetPlatformUserMetrics(ctx context.Context, from, to time.Time) (*domain.UserMetrics, error) {
	// Return basic user metrics
	return &domain.UserMetrics{
		ActiveUsers:    0,
		NewUsers:       0,
		UserGrowthRate: 0,
		UsersByType:    make(map[string]int),
	}, nil
}

func (r *dashboardRepository) GetRevenueMetrics(ctx context.Context, from, to time.Time) (*domain.RevenueMetrics, error) {
	// Return empty revenue metrics for now
	return &domain.RevenueMetrics{
		TotalRevenue:          0,
		RevenueGrowth:         0,
		AverageRevenuePerUser: 0,
		RevenueBySource:       []domain.RevenueBySource{},
	}, nil
}

func (r *dashboardRepository) GetSystemHealthMetrics(ctx context.Context) (*domain.SystemHealth, error) {
	// Return basic system health
	return &domain.SystemHealth{
		Uptime:            99.9,
		ResponseTime:      100,
		ErrorRate:         0.1,
		ActiveConnections: 10,
	}, nil
}

// DeleteCache deletes a cache key
func (r *dashboardRepository) DeleteCache(ctx context.Context, key string) error {
	if r.cache == nil {
		return fmt.Errorf("cache not available")
	}
	
	return r.cache.Del(ctx, key).Err()
}