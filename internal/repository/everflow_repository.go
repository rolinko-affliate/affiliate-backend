package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/affiliate-backend/internal/domain"
	"github.com/affiliate-backend/internal/platform/logger"
)

// EverflowRepository defines the interface for Everflow API operations
type EverflowRepository interface {
	// Entity reporting
	GetEntityReport(ctx context.Context, req *domain.EverflowEntityRequest) (*domain.EverflowEntityResponse, error)
	
	// Conversion reporting
	GetConversions(ctx context.Context, req *domain.EverflowConversionRequest) (*domain.EverflowConversionResponse, error)
	GetConversionByID(ctx context.Context, conversionID string) (*domain.EverflowConversion, error)
	ExportConversions(ctx context.Context, req *domain.EverflowExportRequest) (io.Reader, error)
	
	// Dashboard summary
	GetDashboardSummary(ctx context.Context, timezoneID int) (*domain.EverflowDashboardSummary, error)
	
	// Cache operations
	SetCache(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	GetCache(ctx context.Context, key string, dest interface{}) error
	DeleteCache(ctx context.Context, key string) error
}

// everflowRepository implements EverflowRepository
type everflowRepository struct {
	client   *http.Client
	baseURL  string
	apiKey   string
	cache    *redis.Client
	cacheTTL time.Duration
	logger   *logger.Logger
}

// NewEverflowRepository creates a new Everflow repository
func NewEverflowRepository(baseURL, apiKey string, cache *redis.Client, logger *logger.Logger) EverflowRepository {
	return &everflowRepository{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL:  baseURL,
		apiKey:   apiKey,
		cache:    cache,
		cacheTTL: 5 * time.Minute, // Default cache TTL
		logger:   logger,
	}
}

// GetEntityReport retrieves entity reporting data from Everflow
func (r *everflowRepository) GetEntityReport(ctx context.Context, req *domain.EverflowEntityRequest) (*domain.EverflowEntityResponse, error) {
	// Check cache first
	cacheKey := fmt.Sprintf("everflow:entity:%s:%s", req.From, req.To)
	var cachedResponse domain.EverflowEntityResponse
	if err := r.GetCache(ctx, cacheKey, &cachedResponse); err == nil {
		r.logger.Debug("Cache hit for entity report", "key", cacheKey)
		return &cachedResponse, nil
	}

	// Make API call
	url := fmt.Sprintf("%s/advertisers/reporting/entity", r.baseURL)
	response, err := r.makeAPICall(ctx, "POST", url, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get entity report: %w", err)
	}

	var result domain.EverflowEntityResponse
	if err := json.Unmarshal(response, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal entity response: %w", err)
	}

	// Cache the response
	if err := r.SetCache(ctx, cacheKey, result, r.cacheTTL); err != nil {
		r.logger.Warn("Failed to cache entity report", "error", err)
	}

	return &result, nil
}

// GetConversions retrieves conversion data from Everflow
func (r *everflowRepository) GetConversions(ctx context.Context, req *domain.EverflowConversionRequest) (*domain.EverflowConversionResponse, error) {
	url := fmt.Sprintf("%s/advertisers/reporting/conversions", r.baseURL)
	response, err := r.makeAPICall(ctx, "POST", url, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversions: %w", err)
	}

	var result domain.EverflowConversionResponse
	if err := json.Unmarshal(response, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal conversion response: %w", err)
	}

	return &result, nil
}

// GetConversionByID retrieves a specific conversion by ID
func (r *everflowRepository) GetConversionByID(ctx context.Context, conversionID string) (*domain.EverflowConversion, error) {
	url := fmt.Sprintf("%s/advertisers/reporting/conversions/%s", r.baseURL, conversionID)
	response, err := r.makeAPICall(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversion by ID: %w", err)
	}

	var result domain.EverflowConversion
	if err := json.Unmarshal(response, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal conversion: %w", err)
	}

	return &result, nil
}

// ExportConversions exports conversion data
func (r *everflowRepository) ExportConversions(ctx context.Context, req *domain.EverflowExportRequest) (io.Reader, error) {
	url := fmt.Sprintf("%s/advertisers/reporting/conversions/export", r.baseURL)
	response, err := r.makeAPICall(ctx, "POST", url, req)
	if err != nil {
		return nil, fmt.Errorf("failed to export conversions: %w", err)
	}

	return bytes.NewReader(response), nil
}

// GetDashboardSummary retrieves dashboard summary from Everflow
func (r *everflowRepository) GetDashboardSummary(ctx context.Context, timezoneID int) (*domain.EverflowDashboardSummary, error) {
	// Check cache first
	cacheKey := fmt.Sprintf("everflow:dashboard:summary:%d", timezoneID)
	var cachedSummary domain.EverflowDashboardSummary
	if err := r.GetCache(ctx, cacheKey, &cachedSummary); err == nil {
		r.logger.Debug("Cache hit for dashboard summary", "key", cacheKey)
		return &cachedSummary, nil
	}

	// Make API call
	url := fmt.Sprintf("%s/advertisers/dashboard/summary", r.baseURL)
	req := map[string]int{"timezone_id": timezoneID}
	response, err := r.makeAPICall(ctx, "POST", url, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get dashboard summary: %w", err)
	}

	var result domain.EverflowDashboardSummary
	if err := json.Unmarshal(response, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal dashboard summary: %w", err)
	}

	// Cache the response with shorter TTL for dashboard data
	if err := r.SetCache(ctx, cacheKey, result, 2*time.Minute); err != nil {
		r.logger.Warn("Failed to cache dashboard summary", "error", err)
	}

	return &result, nil
}

// makeAPICall makes an HTTP request to Everflow API
func (r *everflowRepository) makeAPICall(ctx context.Context, method, url string, payload interface{}) ([]byte, error) {
	var body io.Reader
	if payload != nil {
		jsonData, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request: %w", err)
		}
		body = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-Eflow-API-Key", r.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(responseBody))
	}

	return responseBody, nil
}

// Cache operations
func (r *everflowRepository) SetCache(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal cache value: %w", err)
	}

	return r.cache.Set(ctx, key, data, ttl).Err()
}

func (r *everflowRepository) GetCache(ctx context.Context, key string, dest interface{}) error {
	data, err := r.cache.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(data), dest)
}

func (r *everflowRepository) DeleteCache(ctx context.Context, key string) error {
	return r.cache.Del(ctx, key).Err()
}