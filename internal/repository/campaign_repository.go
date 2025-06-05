package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/affiliate-backend/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CampaignRepository defines the interface for campaign operations
type CampaignRepository interface {
	CreateCampaign(ctx context.Context, campaign *domain.Campaign) error
	GetCampaignByID(ctx context.Context, id int64) (*domain.Campaign, error)
	UpdateCampaign(ctx context.Context, campaign *domain.Campaign) error
	ListCampaignsByAdvertiser(ctx context.Context, advertiserID int64, limit, offset int) ([]*domain.Campaign, error)
	ListCampaignsByOrganization(ctx context.Context, orgID int64, limit, offset int) ([]*domain.Campaign, error)
	DeleteCampaign(ctx context.Context, id int64) error
}

// pgxCampaignRepository implements CampaignRepository using pgx
type pgxCampaignRepository struct {
	db *pgxpool.Pool
}

// NewPgxCampaignRepository creates a new campaign repository
func NewPgxCampaignRepository(db *pgxpool.Pool) CampaignRepository {
	return &pgxCampaignRepository{db: db}
}

// CreateCampaign creates a new campaign in the database
func (r *pgxCampaignRepository) CreateCampaign(ctx context.Context, campaign *domain.Campaign) error {
	query := `INSERT INTO public.campaigns 
              (organization_id, advertiser_id, name, description, status, start_date, end_date, internal_notes,
               destination_url, thumbnail_url, preview_url, visibility, currency_id, conversion_method,
               session_definition, session_duration, terms_and_conditions, is_caps_enabled,
               daily_conversion_cap, weekly_conversion_cap, monthly_conversion_cap, global_conversion_cap,
               daily_click_cap, weekly_click_cap, monthly_click_cap, global_click_cap,
               payout_type, payout_amount, revenue_type, revenue_amount,
               created_at, updated_at)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18,
                      $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32)
              RETURNING campaign_id, created_at, updated_at`
	
	// Handle nullable fields
	var description, destinationURL, thumbnailURL, previewURL, visibility, currencyID sql.NullString
	var conversionMethod, sessionDefinition, termsAndConditions, payoutType, revenueType sql.NullString
	var startDate, endDate sql.NullTime
	var internalNotes sql.NullString
	var sessionDuration sql.NullInt32
	var isCapsEnabled sql.NullBool
	var dailyConversionCap, weeklyConversionCap, monthlyConversionCap, globalConversionCap sql.NullInt32
	var dailyClickCap, weeklyClickCap, monthlyClickCap, globalClickCap sql.NullInt32
	var payoutAmount, revenueAmount sql.NullFloat64
	
	// Set nullable fields
	if campaign.Description != nil {
		description = sql.NullString{String: *campaign.Description, Valid: true}
	}
	if campaign.StartDate != nil {
		startDate = sql.NullTime{Time: *campaign.StartDate, Valid: true}
	}
	if campaign.EndDate != nil {
		endDate = sql.NullTime{Time: *campaign.EndDate, Valid: true}
	}
	if campaign.InternalNotes != nil {
		internalNotes = sql.NullString{String: *campaign.InternalNotes, Valid: true}
	}
	if campaign.DestinationURL != nil {
		destinationURL = sql.NullString{String: *campaign.DestinationURL, Valid: true}
	}
	if campaign.ThumbnailURL != nil {
		thumbnailURL = sql.NullString{String: *campaign.ThumbnailURL, Valid: true}
	}
	if campaign.PreviewURL != nil {
		previewURL = sql.NullString{String: *campaign.PreviewURL, Valid: true}
	}
	if campaign.Visibility != nil {
		visibility = sql.NullString{String: *campaign.Visibility, Valid: true}
	}
	if campaign.CurrencyID != nil {
		currencyID = sql.NullString{String: *campaign.CurrencyID, Valid: true}
	}
	if campaign.ConversionMethod != nil {
		conversionMethod = sql.NullString{String: *campaign.ConversionMethod, Valid: true}
	}
	if campaign.SessionDefinition != nil {
		sessionDefinition = sql.NullString{String: *campaign.SessionDefinition, Valid: true}
	}
	if campaign.SessionDuration != nil {
		sessionDuration = sql.NullInt32{Int32: *campaign.SessionDuration, Valid: true}
	}
	if campaign.TermsAndConditions != nil {
		termsAndConditions = sql.NullString{String: *campaign.TermsAndConditions, Valid: true}
	}
	if campaign.IsCapsEnabled != nil {
		isCapsEnabled = sql.NullBool{Bool: *campaign.IsCapsEnabled, Valid: true}
	}
	if campaign.DailyConversionCap != nil {
		dailyConversionCap = sql.NullInt32{Int32: int32(*campaign.DailyConversionCap), Valid: true}
	}
	if campaign.WeeklyConversionCap != nil {
		weeklyConversionCap = sql.NullInt32{Int32: int32(*campaign.WeeklyConversionCap), Valid: true}
	}
	if campaign.MonthlyConversionCap != nil {
		monthlyConversionCap = sql.NullInt32{Int32: int32(*campaign.MonthlyConversionCap), Valid: true}
	}
	if campaign.GlobalConversionCap != nil {
		globalConversionCap = sql.NullInt32{Int32: int32(*campaign.GlobalConversionCap), Valid: true}
	}
	if campaign.DailyClickCap != nil {
		dailyClickCap = sql.NullInt32{Int32: int32(*campaign.DailyClickCap), Valid: true}
	}
	if campaign.WeeklyClickCap != nil {
		weeklyClickCap = sql.NullInt32{Int32: int32(*campaign.WeeklyClickCap), Valid: true}
	}
	if campaign.MonthlyClickCap != nil {
		monthlyClickCap = sql.NullInt32{Int32: int32(*campaign.MonthlyClickCap), Valid: true}
	}
	if campaign.GlobalClickCap != nil {
		globalClickCap = sql.NullInt32{Int32: int32(*campaign.GlobalClickCap), Valid: true}
	}
	if campaign.PayoutType != nil {
		payoutType = sql.NullString{String: *campaign.PayoutType, Valid: true}
	}
	if campaign.PayoutAmount != nil {
		payoutAmount = sql.NullFloat64{Float64: *campaign.PayoutAmount, Valid: true}
	}
	if campaign.RevenueType != nil {
		revenueType = sql.NullString{String: *campaign.RevenueType, Valid: true}
	}
	if campaign.RevenueAmount != nil {
		revenueAmount = sql.NullFloat64{Float64: *campaign.RevenueAmount, Valid: true}
	}
	
	now := time.Now()
	
	err := r.db.QueryRow(ctx, query,
		campaign.OrganizationID, campaign.AdvertiserID, campaign.Name, description, campaign.Status,
		startDate, endDate, internalNotes,
		destinationURL, thumbnailURL, previewURL, visibility, currencyID, conversionMethod,
		sessionDefinition, sessionDuration, termsAndConditions, isCapsEnabled,
		dailyConversionCap, weeklyConversionCap, monthlyConversionCap, globalConversionCap,
		dailyClickCap, weeklyClickCap, monthlyClickCap, globalClickCap,
		payoutType, payoutAmount, revenueType, revenueAmount,
		now, now,
	).Scan(&campaign.CampaignID, &campaign.CreatedAt, &campaign.UpdatedAt)
	
	if err != nil {
		return fmt.Errorf("failed to create campaign: %w", err)
	}
	
	return nil
}

// GetCampaignByID retrieves a campaign by its ID
func (r *pgxCampaignRepository) GetCampaignByID(ctx context.Context, id int64) (*domain.Campaign, error) {
	query := `SELECT campaign_id, organization_id, advertiser_id, name, description, status,
              start_date, end_date, internal_notes,
              destination_url, thumbnail_url, preview_url, visibility, currency_id, conversion_method,
              session_definition, session_duration, terms_and_conditions, is_caps_enabled,
              daily_conversion_cap, weekly_conversion_cap, monthly_conversion_cap, global_conversion_cap,
              daily_click_cap, weekly_click_cap, monthly_click_cap, global_click_cap,
              payout_type, payout_amount, revenue_type, revenue_amount,
              created_at, updated_at
              FROM public.campaigns WHERE campaign_id = $1`
	
	campaign := &domain.Campaign{}
	var description, destinationURL, thumbnailURL, previewURL, visibility, currencyID sql.NullString
	var conversionMethod, sessionDefinition, termsAndConditions, payoutType, revenueType sql.NullString
	var startDate, endDate sql.NullTime
	var internalNotes sql.NullString
	var sessionDuration sql.NullInt32
	var isCapsEnabled sql.NullBool
	var dailyConversionCap, weeklyConversionCap, monthlyConversionCap, globalConversionCap sql.NullInt32
	var dailyClickCap, weeklyClickCap, monthlyClickCap, globalClickCap sql.NullInt32
	var payoutAmount, revenueAmount sql.NullFloat64
	
	err := r.db.QueryRow(ctx, query, id).Scan(
		&campaign.CampaignID, &campaign.OrganizationID, &campaign.AdvertiserID,
		&campaign.Name, &description, &campaign.Status,
		&startDate, &endDate, &internalNotes,
		&destinationURL, &thumbnailURL, &previewURL, &visibility, &currencyID, &conversionMethod,
		&sessionDefinition, &sessionDuration, &termsAndConditions, &isCapsEnabled,
		&dailyConversionCap, &weeklyConversionCap, &monthlyConversionCap, &globalConversionCap,
		&dailyClickCap, &weeklyClickCap, &monthlyClickCap, &globalClickCap,
		&payoutType, &payoutAmount, &revenueType, &revenueAmount,
		&campaign.CreatedAt, &campaign.UpdatedAt,
	)
	
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("campaign not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get campaign: %w", err)
	}
	
	// Handle nullable fields
	if description.Valid {
		campaign.Description = &description.String
	}
	if startDate.Valid {
		campaign.StartDate = &startDate.Time
	}
	if endDate.Valid {
		campaign.EndDate = &endDate.Time
	}
	if internalNotes.Valid {
		campaign.InternalNotes = &internalNotes.String
	}
	if destinationURL.Valid {
		campaign.DestinationURL = &destinationURL.String
	}
	if thumbnailURL.Valid {
		campaign.ThumbnailURL = &thumbnailURL.String
	}
	if previewURL.Valid {
		campaign.PreviewURL = &previewURL.String
	}
	if visibility.Valid {
		campaign.Visibility = &visibility.String
	}
	if currencyID.Valid {
		campaign.CurrencyID = &currencyID.String
	}
	if conversionMethod.Valid {
		campaign.ConversionMethod = &conversionMethod.String
	}
	if sessionDefinition.Valid {
		campaign.SessionDefinition = &sessionDefinition.String
	}
	if sessionDuration.Valid {
		campaign.SessionDuration = &sessionDuration.Int32
	}
	if termsAndConditions.Valid {
		campaign.TermsAndConditions = &termsAndConditions.String
	}
	if isCapsEnabled.Valid {
		campaign.IsCapsEnabled = &isCapsEnabled.Bool
	}
	if dailyConversionCap.Valid {
		val := int(dailyConversionCap.Int32)
		campaign.DailyConversionCap = &val
	}
	if weeklyConversionCap.Valid {
		val := int(weeklyConversionCap.Int32)
		campaign.WeeklyConversionCap = &val
	}
	if monthlyConversionCap.Valid {
		val := int(monthlyConversionCap.Int32)
		campaign.MonthlyConversionCap = &val
	}
	if globalConversionCap.Valid {
		val := int(globalConversionCap.Int32)
		campaign.GlobalConversionCap = &val
	}
	if dailyClickCap.Valid {
		val := int(dailyClickCap.Int32)
		campaign.DailyClickCap = &val
	}
	if weeklyClickCap.Valid {
		val := int(weeklyClickCap.Int32)
		campaign.WeeklyClickCap = &val
	}
	if monthlyClickCap.Valid {
		val := int(monthlyClickCap.Int32)
		campaign.MonthlyClickCap = &val
	}
	if globalClickCap.Valid {
		val := int(globalClickCap.Int32)
		campaign.GlobalClickCap = &val
	}
	if payoutType.Valid {
		campaign.PayoutType = &payoutType.String
	}
	if payoutAmount.Valid {
		campaign.PayoutAmount = &payoutAmount.Float64
	}
	if revenueType.Valid {
		campaign.RevenueType = &revenueType.String
	}
	if revenueAmount.Valid {
		campaign.RevenueAmount = &revenueAmount.Float64
	}
	
	return campaign, nil
}

// UpdateCampaign updates an existing campaign
func (r *pgxCampaignRepository) UpdateCampaign(ctx context.Context, campaign *domain.Campaign) error {
	query := `UPDATE public.campaigns SET 
              organization_id = $2, advertiser_id = $3, name = $4, description = $5, status = $6,
              start_date = $7, end_date = $8, internal_notes = $9,
              destination_url = $10, thumbnail_url = $11, preview_url = $12, visibility = $13,
              currency_id = $14, conversion_method = $15, session_definition = $16, session_duration = $17,
              terms_and_conditions = $18, is_caps_enabled = $19,
              daily_conversion_cap = $20, weekly_conversion_cap = $21, monthly_conversion_cap = $22, global_conversion_cap = $23,
              daily_click_cap = $24, weekly_click_cap = $25, monthly_click_cap = $26, global_click_cap = $27,
              payout_type = $28, payout_amount = $29, revenue_type = $30, revenue_amount = $31,
              updated_at = $32
              WHERE campaign_id = $1`
	
	// Handle nullable fields (same as CreateCampaign)
	var description, destinationURL, thumbnailURL, previewURL, visibility, currencyID sql.NullString
	var conversionMethod, sessionDefinition, termsAndConditions, payoutType, revenueType sql.NullString
	var startDate, endDate sql.NullTime
	var internalNotes sql.NullString
	var sessionDuration sql.NullInt32
	var isCapsEnabled sql.NullBool
	var dailyConversionCap, weeklyConversionCap, monthlyConversionCap, globalConversionCap sql.NullInt32
	var dailyClickCap, weeklyClickCap, monthlyClickCap, globalClickCap sql.NullInt32
	var payoutAmount, revenueAmount sql.NullFloat64
	
	// Set nullable fields (same logic as CreateCampaign)
	if campaign.Description != nil {
		description = sql.NullString{String: *campaign.Description, Valid: true}
	}
	if campaign.StartDate != nil {
		startDate = sql.NullTime{Time: *campaign.StartDate, Valid: true}
	}
	if campaign.EndDate != nil {
		endDate = sql.NullTime{Time: *campaign.EndDate, Valid: true}
	}
	if campaign.InternalNotes != nil {
		internalNotes = sql.NullString{String: *campaign.InternalNotes, Valid: true}
	}
	if campaign.DestinationURL != nil {
		destinationURL = sql.NullString{String: *campaign.DestinationURL, Valid: true}
	}
	if campaign.ThumbnailURL != nil {
		thumbnailURL = sql.NullString{String: *campaign.ThumbnailURL, Valid: true}
	}
	if campaign.PreviewURL != nil {
		previewURL = sql.NullString{String: *campaign.PreviewURL, Valid: true}
	}
	if campaign.Visibility != nil {
		visibility = sql.NullString{String: *campaign.Visibility, Valid: true}
	}
	if campaign.CurrencyID != nil {
		currencyID = sql.NullString{String: *campaign.CurrencyID, Valid: true}
	}
	if campaign.ConversionMethod != nil {
		conversionMethod = sql.NullString{String: *campaign.ConversionMethod, Valid: true}
	}
	if campaign.SessionDefinition != nil {
		sessionDefinition = sql.NullString{String: *campaign.SessionDefinition, Valid: true}
	}
	if campaign.SessionDuration != nil {
		sessionDuration = sql.NullInt32{Int32: *campaign.SessionDuration, Valid: true}
	}
	if campaign.TermsAndConditions != nil {
		termsAndConditions = sql.NullString{String: *campaign.TermsAndConditions, Valid: true}
	}
	if campaign.IsCapsEnabled != nil {
		isCapsEnabled = sql.NullBool{Bool: *campaign.IsCapsEnabled, Valid: true}
	}
	if campaign.DailyConversionCap != nil {
		dailyConversionCap = sql.NullInt32{Int32: int32(*campaign.DailyConversionCap), Valid: true}
	}
	if campaign.WeeklyConversionCap != nil {
		weeklyConversionCap = sql.NullInt32{Int32: int32(*campaign.WeeklyConversionCap), Valid: true}
	}
	if campaign.MonthlyConversionCap != nil {
		monthlyConversionCap = sql.NullInt32{Int32: int32(*campaign.MonthlyConversionCap), Valid: true}
	}
	if campaign.GlobalConversionCap != nil {
		globalConversionCap = sql.NullInt32{Int32: int32(*campaign.GlobalConversionCap), Valid: true}
	}
	if campaign.DailyClickCap != nil {
		dailyClickCap = sql.NullInt32{Int32: int32(*campaign.DailyClickCap), Valid: true}
	}
	if campaign.WeeklyClickCap != nil {
		weeklyClickCap = sql.NullInt32{Int32: int32(*campaign.WeeklyClickCap), Valid: true}
	}
	if campaign.MonthlyClickCap != nil {
		monthlyClickCap = sql.NullInt32{Int32: int32(*campaign.MonthlyClickCap), Valid: true}
	}
	if campaign.GlobalClickCap != nil {
		globalClickCap = sql.NullInt32{Int32: int32(*campaign.GlobalClickCap), Valid: true}
	}
	if campaign.PayoutType != nil {
		payoutType = sql.NullString{String: *campaign.PayoutType, Valid: true}
	}
	if campaign.PayoutAmount != nil {
		payoutAmount = sql.NullFloat64{Float64: *campaign.PayoutAmount, Valid: true}
	}
	if campaign.RevenueType != nil {
		revenueType = sql.NullString{String: *campaign.RevenueType, Valid: true}
	}
	if campaign.RevenueAmount != nil {
		revenueAmount = sql.NullFloat64{Float64: *campaign.RevenueAmount, Valid: true}
	}
	
	now := time.Now()
	
	result, err := r.db.Exec(ctx, query,
		campaign.CampaignID, campaign.OrganizationID, campaign.AdvertiserID,
		campaign.Name, description, campaign.Status,
		startDate, endDate, internalNotes,
		destinationURL, thumbnailURL, previewURL, visibility, currencyID, conversionMethod,
		sessionDefinition, sessionDuration, termsAndConditions, isCapsEnabled,
		dailyConversionCap, weeklyConversionCap, monthlyConversionCap, globalConversionCap,
		dailyClickCap, weeklyClickCap, monthlyClickCap, globalClickCap,
		payoutType, payoutAmount, revenueType, revenueAmount,
		now,
	)
	
	if err != nil {
		return fmt.Errorf("failed to update campaign: %w", err)
	}
	
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("campaign not found: not found")
	}
	
	campaign.UpdatedAt = now
	return nil
}

// ListCampaignsByAdvertiser retrieves campaigns for a specific advertiser with pagination
func (r *pgxCampaignRepository) ListCampaignsByAdvertiser(ctx context.Context, advertiserID int64, limit, offset int) ([]*domain.Campaign, error) {
	query := `SELECT campaign_id, organization_id, advertiser_id, name, description, status,
              start_date, end_date, internal_notes,
              destination_url, thumbnail_url, preview_url, visibility, currency_id, conversion_method,
              session_definition, session_duration, terms_and_conditions, is_caps_enabled,
              daily_conversion_cap, weekly_conversion_cap, monthly_conversion_cap, global_conversion_cap,
              daily_click_cap, weekly_click_cap, monthly_click_cap, global_click_cap,
              payout_type, payout_amount, revenue_type, revenue_amount,
              created_at, updated_at
              FROM public.campaigns WHERE advertiser_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	
	rows, err := r.db.Query(ctx, query, advertiserID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list campaigns by advertiser: %w", err)
	}
	defer rows.Close()
	
	var campaigns []*domain.Campaign
	for rows.Next() {
		campaign := &domain.Campaign{}
		var description, destinationURL, thumbnailURL, previewURL, visibility, currencyID sql.NullString
		var conversionMethod, sessionDefinition, termsAndConditions, payoutType, revenueType sql.NullString
		var startDate, endDate sql.NullTime
		var internalNotes sql.NullString
		var sessionDuration sql.NullInt32
		var isCapsEnabled sql.NullBool
		var dailyConversionCap, weeklyConversionCap, monthlyConversionCap, globalConversionCap sql.NullInt32
		var dailyClickCap, weeklyClickCap, monthlyClickCap, globalClickCap sql.NullInt32
		var payoutAmount, revenueAmount sql.NullFloat64
		
		err := rows.Scan(
			&campaign.CampaignID, &campaign.OrganizationID, &campaign.AdvertiserID,
			&campaign.Name, &description, &campaign.Status,
			&startDate, &endDate, &internalNotes,
			&destinationURL, &thumbnailURL, &previewURL, &visibility, &currencyID, &conversionMethod,
			&sessionDefinition, &sessionDuration, &termsAndConditions, &isCapsEnabled,
			&dailyConversionCap, &weeklyConversionCap, &monthlyConversionCap, &globalConversionCap,
			&dailyClickCap, &weeklyClickCap, &monthlyClickCap, &globalClickCap,
			&payoutType, &payoutAmount, &revenueType, &revenueAmount,
			&campaign.CreatedAt, &campaign.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan campaign: %w", err)
		}
		
		// Handle nullable fields (same logic as GetCampaignByID)
		if description.Valid {
			campaign.Description = &description.String
		}
		if startDate.Valid {
			campaign.StartDate = &startDate.Time
		}
		if endDate.Valid {
			campaign.EndDate = &endDate.Time
		}
		if internalNotes.Valid {
			campaign.InternalNotes = &internalNotes.String
		}
		if destinationURL.Valid {
			campaign.DestinationURL = &destinationURL.String
		}
		if thumbnailURL.Valid {
			campaign.ThumbnailURL = &thumbnailURL.String
		}
		if previewURL.Valid {
			campaign.PreviewURL = &previewURL.String
		}
		if visibility.Valid {
			campaign.Visibility = &visibility.String
		}
		if currencyID.Valid {
			campaign.CurrencyID = &currencyID.String
		}
		if conversionMethod.Valid {
			campaign.ConversionMethod = &conversionMethod.String
		}
		if sessionDefinition.Valid {
			campaign.SessionDefinition = &sessionDefinition.String
		}
		if sessionDuration.Valid {
			campaign.SessionDuration = &sessionDuration.Int32
		}
		if termsAndConditions.Valid {
			campaign.TermsAndConditions = &termsAndConditions.String
		}
		if isCapsEnabled.Valid {
			campaign.IsCapsEnabled = &isCapsEnabled.Bool
		}
		if dailyConversionCap.Valid {
			val := int(dailyConversionCap.Int32)
			campaign.DailyConversionCap = &val
		}
		if weeklyConversionCap.Valid {
			val := int(weeklyConversionCap.Int32)
			campaign.WeeklyConversionCap = &val
		}
		if monthlyConversionCap.Valid {
			val := int(monthlyConversionCap.Int32)
			campaign.MonthlyConversionCap = &val
		}
		if globalConversionCap.Valid {
			val := int(globalConversionCap.Int32)
			campaign.GlobalConversionCap = &val
		}
		if dailyClickCap.Valid {
			val := int(dailyClickCap.Int32)
			campaign.DailyClickCap = &val
		}
		if weeklyClickCap.Valid {
			val := int(weeklyClickCap.Int32)
			campaign.WeeklyClickCap = &val
		}
		if monthlyClickCap.Valid {
			val := int(monthlyClickCap.Int32)
			campaign.MonthlyClickCap = &val
		}
		if globalClickCap.Valid {
			val := int(globalClickCap.Int32)
			campaign.GlobalClickCap = &val
		}
		if payoutType.Valid {
			campaign.PayoutType = &payoutType.String
		}
		if payoutAmount.Valid {
			campaign.PayoutAmount = &payoutAmount.Float64
		}
		if revenueType.Valid {
			campaign.RevenueType = &revenueType.String
		}
		if revenueAmount.Valid {
			campaign.RevenueAmount = &revenueAmount.Float64
		}
		
		campaigns = append(campaigns, campaign)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate campaigns: %w", err)
	}
	
	return campaigns, nil
}

// ListCampaignsByOrganization retrieves campaigns for a specific organization with pagination
func (r *pgxCampaignRepository) ListCampaignsByOrganization(ctx context.Context, orgID int64, limit, offset int) ([]*domain.Campaign, error) {
	query := `SELECT campaign_id, organization_id, advertiser_id, name, description, status,
              start_date, end_date, internal_notes,
              destination_url, thumbnail_url, preview_url, visibility, currency_id, conversion_method,
              session_definition, session_duration, terms_and_conditions, is_caps_enabled,
              daily_conversion_cap, weekly_conversion_cap, monthly_conversion_cap, global_conversion_cap,
              daily_click_cap, weekly_click_cap, monthly_click_cap, global_click_cap,
              payout_type, payout_amount, revenue_type, revenue_amount,
              created_at, updated_at
              FROM public.campaigns WHERE organization_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	
	rows, err := r.db.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list campaigns by organization: %w", err)
	}
	defer rows.Close()
	
	var campaigns []*domain.Campaign
	for rows.Next() {
		campaign := &domain.Campaign{}
		var description, destinationURL, thumbnailURL, previewURL, visibility, currencyID sql.NullString
		var conversionMethod, sessionDefinition, termsAndConditions, payoutType, revenueType sql.NullString
		var startDate, endDate sql.NullTime
		var internalNotes sql.NullString
		var sessionDuration sql.NullInt32
		var isCapsEnabled sql.NullBool
		var dailyConversionCap, weeklyConversionCap, monthlyConversionCap, globalConversionCap sql.NullInt32
		var dailyClickCap, weeklyClickCap, monthlyClickCap, globalClickCap sql.NullInt32
		var payoutAmount, revenueAmount sql.NullFloat64
		
		err := rows.Scan(
			&campaign.CampaignID, &campaign.OrganizationID, &campaign.AdvertiserID,
			&campaign.Name, &description, &campaign.Status,
			&startDate, &endDate, &internalNotes,
			&destinationURL, &thumbnailURL, &previewURL, &visibility, &currencyID, &conversionMethod,
			&sessionDefinition, &sessionDuration, &termsAndConditions, &isCapsEnabled,
			&dailyConversionCap, &weeklyConversionCap, &monthlyConversionCap, &globalConversionCap,
			&dailyClickCap, &weeklyClickCap, &monthlyClickCap, &globalClickCap,
			&payoutType, &payoutAmount, &revenueType, &revenueAmount,
			&campaign.CreatedAt, &campaign.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan campaign: %w", err)
		}
		
		// Handle nullable fields (same logic as GetCampaignByID)
		if description.Valid {
			campaign.Description = &description.String
		}
		if startDate.Valid {
			campaign.StartDate = &startDate.Time
		}
		if endDate.Valid {
			campaign.EndDate = &endDate.Time
		}
		if internalNotes.Valid {
			campaign.InternalNotes = &internalNotes.String
		}
		if destinationURL.Valid {
			campaign.DestinationURL = &destinationURL.String
		}
		if thumbnailURL.Valid {
			campaign.ThumbnailURL = &thumbnailURL.String
		}
		if previewURL.Valid {
			campaign.PreviewURL = &previewURL.String
		}
		if visibility.Valid {
			campaign.Visibility = &visibility.String
		}
		if currencyID.Valid {
			campaign.CurrencyID = &currencyID.String
		}
		if conversionMethod.Valid {
			campaign.ConversionMethod = &conversionMethod.String
		}
		if sessionDefinition.Valid {
			campaign.SessionDefinition = &sessionDefinition.String
		}
		if sessionDuration.Valid {
			campaign.SessionDuration = &sessionDuration.Int32
		}
		if termsAndConditions.Valid {
			campaign.TermsAndConditions = &termsAndConditions.String
		}
		if isCapsEnabled.Valid {
			campaign.IsCapsEnabled = &isCapsEnabled.Bool
		}
		if dailyConversionCap.Valid {
			val := int(dailyConversionCap.Int32)
			campaign.DailyConversionCap = &val
		}
		if weeklyConversionCap.Valid {
			val := int(weeklyConversionCap.Int32)
			campaign.WeeklyConversionCap = &val
		}
		if monthlyConversionCap.Valid {
			val := int(monthlyConversionCap.Int32)
			campaign.MonthlyConversionCap = &val
		}
		if globalConversionCap.Valid {
			val := int(globalConversionCap.Int32)
			campaign.GlobalConversionCap = &val
		}
		if dailyClickCap.Valid {
			val := int(dailyClickCap.Int32)
			campaign.DailyClickCap = &val
		}
		if weeklyClickCap.Valid {
			val := int(weeklyClickCap.Int32)
			campaign.WeeklyClickCap = &val
		}
		if monthlyClickCap.Valid {
			val := int(monthlyClickCap.Int32)
			campaign.MonthlyClickCap = &val
		}
		if globalClickCap.Valid {
			val := int(globalClickCap.Int32)
			campaign.GlobalClickCap = &val
		}
		if payoutType.Valid {
			campaign.PayoutType = &payoutType.String
		}
		if payoutAmount.Valid {
			campaign.PayoutAmount = &payoutAmount.Float64
		}
		if revenueType.Valid {
			campaign.RevenueType = &revenueType.String
		}
		if revenueAmount.Valid {
			campaign.RevenueAmount = &revenueAmount.Float64
		}
		
		campaigns = append(campaigns, campaign)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate campaigns: %w", err)
	}
	
	return campaigns, nil
}

// DeleteCampaign deletes a campaign by its ID
func (r *pgxCampaignRepository) DeleteCampaign(ctx context.Context, id int64) error {
	query := `DELETE FROM public.campaigns WHERE campaign_id = $1`
	
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete campaign: %w", err)
	}
	
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("campaign not found: not found")
	}
	
	return nil
}