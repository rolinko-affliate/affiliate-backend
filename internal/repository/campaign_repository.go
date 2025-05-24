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
	ListCampaignsByOrganization(ctx context.Context, orgID int64, limit, offset int) ([]*domain.Campaign, error)
	ListCampaignsByAdvertiser(ctx context.Context, advertiserID int64, limit, offset int) ([]*domain.Campaign, error)
	DeleteCampaign(ctx context.Context, id int64) error
	
	// Provider offer methods
	CreateCampaignProviderOffer(ctx context.Context, offer *domain.CampaignProviderOffer) error
	GetCampaignProviderOfferByID(ctx context.Context, id int64) (*domain.CampaignProviderOffer, error)
	UpdateCampaignProviderOffer(ctx context.Context, offer *domain.CampaignProviderOffer) error
	ListCampaignProviderOffersByCampaign(ctx context.Context, campaignID int64) ([]*domain.CampaignProviderOffer, error)
	DeleteCampaignProviderOffer(ctx context.Context, id int64) error
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
              (organization_id, advertiser_id, name, description, status, start_date, end_date,
               destination_url, thumbnail_url, preview_url, visibility, currency_id, 
               conversion_method, session_definition, session_duration, internal_notes,
               terms_and_conditions, is_force_terms_and_conditions, is_caps_enabled,
               daily_conversion_cap, weekly_conversion_cap, monthly_conversion_cap, global_conversion_cap,
               daily_click_cap, weekly_click_cap, monthly_click_cap, global_click_cap,
               payout_type, payout_amount, revenue_type, revenue_amount, offer_config,
               created_at, updated_at)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19,
                      $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34)
              RETURNING campaign_id, created_at, updated_at`
	
	var description sql.NullString
	var startDate, endDate sql.NullTime
	var destinationURL, thumbnailURL, previewURL, visibility, currencyID sql.NullString
	var conversionMethod, sessionDefinition, internalNotes, termsAndConditions sql.NullString
	var payoutType, revenueType, offerConfig sql.NullString
	var sessionDuration sql.NullInt32
	var isForceTermsAndConditions, isCapsEnabled sql.NullBool
	var dailyConversionCap, weeklyConversionCap, monthlyConversionCap, globalConversionCap sql.NullInt32
	var dailyClickCap, weeklyClickCap, monthlyClickCap, globalClickCap sql.NullInt32
	var payoutAmount, revenueAmount sql.NullFloat64
	
	if campaign.Description != nil {
		description = sql.NullString{String: *campaign.Description, Valid: true}
	}
	
	if campaign.StartDate != nil {
		startDate = sql.NullTime{Time: *campaign.StartDate, Valid: true}
	}
	
	if campaign.EndDate != nil {
		endDate = sql.NullTime{Time: *campaign.EndDate, Valid: true}
	}
	
	// Handle offer-specific fields
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
		sessionDuration = sql.NullInt32{Int32: int32(*campaign.SessionDuration), Valid: true}
	}
	if campaign.InternalNotes != nil {
		internalNotes = sql.NullString{String: *campaign.InternalNotes, Valid: true}
	}
	if campaign.TermsAndConditions != nil {
		termsAndConditions = sql.NullString{String: *campaign.TermsAndConditions, Valid: true}
	}
	if campaign.IsForceTermsAndConditions != nil {
		isForceTermsAndConditions = sql.NullBool{Bool: *campaign.IsForceTermsAndConditions, Valid: true}
	}
	if campaign.IsCapsEnabled != nil {
		isCapsEnabled = sql.NullBool{Bool: *campaign.IsCapsEnabled, Valid: true}
	}

	// Handle caps
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

	// Handle payout and revenue
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
	if campaign.OfferConfig != nil {
		offerConfig = sql.NullString{String: *campaign.OfferConfig, Valid: true}
	}

	now := time.Now()
	err := r.db.QueryRow(ctx, query,
		campaign.OrganizationID,
		campaign.AdvertiserID,
		campaign.Name,
		description,
		campaign.Status,
		startDate,
		endDate,
		destinationURL,
		thumbnailURL,
		previewURL,
		visibility,
		currencyID,
		conversionMethod,
		sessionDefinition,
		sessionDuration,
		internalNotes,
		termsAndConditions,
		isForceTermsAndConditions,
		isCapsEnabled,
		dailyConversionCap,
		weeklyConversionCap,
		monthlyConversionCap,
		globalConversionCap,
		dailyClickCap,
		weeklyClickCap,
		monthlyClickCap,
		globalClickCap,
		payoutType,
		payoutAmount,
		revenueType,
		revenueAmount,
		offerConfig,
		now,
		now,
	).Scan(
		&campaign.CampaignID,
		&campaign.CreatedAt,
		&campaign.UpdatedAt,
	)
	
	if err != nil {
		return fmt.Errorf("error creating campaign: %w", err)
	}
	return nil
}

// GetCampaignByID retrieves a campaign by ID
func (r *pgxCampaignRepository) GetCampaignByID(ctx context.Context, id int64) (*domain.Campaign, error) {
	query := `SELECT campaign_id, organization_id, advertiser_id, name, description, status, start_date, end_date,
              destination_url, thumbnail_url, preview_url, visibility, currency_id, 
              conversion_method, session_definition, session_duration, caps_timezone_id, project_id,
              date_live_until, html_description, internal_notes, terms_and_conditions, 
              is_using_explicit_terms_and_conditions, is_force_terms_and_conditions, is_caps_enabled,
              is_whitelist_check_enabled, is_view_through_enabled, server_side_url, 
              view_through_destination_url, is_description_plain_text, is_use_direct_linking,
              app_identifier, daily_conversion_cap, weekly_conversion_cap, monthly_conversion_cap, 
              global_conversion_cap, daily_click_cap, weekly_click_cap, monthly_click_cap, 
              global_click_cap, encoded_value, today_clicks, today_revenue, time_created, time_saved,
              payout_type, payout_amount, revenue_type, revenue_amount, offer_config,
              created_at, updated_at
              FROM public.campaigns WHERE campaign_id = $1`
	
	var campaign domain.Campaign
	var description sql.NullString
	var startDate, endDate, dateLiveUntil sql.NullTime
	var destinationURL, thumbnailURL, previewURL, visibility, currencyID sql.NullString
	var conversionMethod, sessionDefinition, internalNotes, termsAndConditions sql.NullString
	var projectID, htmlDescription, serverSideURL, viewThroughDestinationURL sql.NullString
	var appIdentifier, encodedValue, todayRevenue sql.NullString
	var payoutType, revenueType, offerConfig sql.NullString
	var sessionDuration, capsTimezoneID sql.NullInt32
	var dailyConversionCap, weeklyConversionCap, monthlyConversionCap, globalConversionCap sql.NullInt32
	var dailyClickCap, weeklyClickCap, monthlyClickCap, globalClickCap sql.NullInt32
	var todayClicks, timeCreated, timeSaved sql.NullInt32
	var isUsingExplicitTermsAndConditions, isForceTermsAndConditions, isCapsEnabled sql.NullBool
	var isWhitelistCheckEnabled, isViewThroughEnabled, isDescriptionPlainText, isUseDirectLinking sql.NullBool
	var payoutAmount, revenueAmount sql.NullFloat64
	
	err := r.db.QueryRow(ctx, query, id).Scan(
		&campaign.CampaignID,
		&campaign.OrganizationID,
		&campaign.AdvertiserID,
		&campaign.Name,
		&description,
		&campaign.Status,
		&startDate,
		&endDate,
		&destinationURL,
		&thumbnailURL,
		&previewURL,
		&visibility,
		&currencyID,
		&conversionMethod,
		&sessionDefinition,
		&sessionDuration,
		&capsTimezoneID,
		&projectID,
		&dateLiveUntil,
		&htmlDescription,
		&internalNotes,
		&termsAndConditions,
		&isUsingExplicitTermsAndConditions,
		&isForceTermsAndConditions,
		&isCapsEnabled,
		&isWhitelistCheckEnabled,
		&isViewThroughEnabled,
		&serverSideURL,
		&viewThroughDestinationURL,
		&isDescriptionPlainText,
		&isUseDirectLinking,
		&appIdentifier,
		&dailyConversionCap,
		&weeklyConversionCap,
		&monthlyConversionCap,
		&globalConversionCap,
		&dailyClickCap,
		&weeklyClickCap,
		&monthlyClickCap,
		&globalClickCap,
		&encodedValue,
		&todayClicks,
		&todayRevenue,
		&timeCreated,
		&timeSaved,
		&payoutType,
		&payoutAmount,
		&revenueType,
		&revenueAmount,
		&offerConfig,
		&campaign.CreatedAt,
		&campaign.UpdatedAt,
	)
	
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("campaign not found: %w", domain.ErrNotFound)
		}
		return nil, fmt.Errorf("error getting campaign by ID: %w", err)
	}
	
	if description.Valid {
		desc := description.String
		campaign.Description = &desc
	}
	
	if startDate.Valid {
		date := startDate.Time
		campaign.StartDate = &date
	}
	
	if endDate.Valid {
		date := endDate.Time
		campaign.EndDate = &date
	}
	
	// Assign nullable fields for Everflow offer support
	if destinationURL.Valid {
		url := destinationURL.String
		campaign.DestinationURL = &url
	}
	
	if thumbnailURL.Valid {
		url := thumbnailURL.String
		campaign.ThumbnailURL = &url
	}
	
	if previewURL.Valid {
		url := previewURL.String
		campaign.PreviewURL = &url
	}
	
	if visibility.Valid {
		vis := visibility.String
		campaign.Visibility = &vis
	}
	
	if currencyID.Valid {
		curr := currencyID.String
		campaign.CurrencyID = &curr
	}
	
	if conversionMethod.Valid {
		method := conversionMethod.String
		campaign.ConversionMethod = &method
	}
	
	if sessionDefinition.Valid {
		def := sessionDefinition.String
		campaign.SessionDefinition = &def
	}
	
	if sessionDuration.Valid {
		duration := int(sessionDuration.Int32)
		campaign.SessionDuration = &duration
	}
	
	if capsTimezoneID.Valid {
		timezone := int(capsTimezoneID.Int32)
		campaign.CapsTimezoneID = &timezone
	}
	
	if projectID.Valid {
		project := projectID.String
		campaign.ProjectID = &project
	}
	
	if dateLiveUntil.Valid {
		date := dateLiveUntil.Time
		campaign.DateLiveUntil = &date
	}
	
	if htmlDescription.Valid {
		desc := htmlDescription.String
		campaign.HTMLDescription = &desc
	}
	
	if termsAndConditions.Valid {
		terms := termsAndConditions.String
		campaign.TermsAndConditions = &terms
	}
	
	if serverSideURL.Valid {
		url := serverSideURL.String
		campaign.ServerSideURL = &url
	}
	
	if viewThroughDestinationURL.Valid {
		url := viewThroughDestinationURL.String
		campaign.ViewThroughDestinationURL = &url
	}
	
	if appIdentifier.Valid {
		app := appIdentifier.String
		campaign.AppIdentifier = &app
	}
	
	if dailyConversionCap.Valid {
		cap := int(dailyConversionCap.Int32)
		campaign.DailyConversionCap = &cap
	}
	
	if weeklyConversionCap.Valid {
		cap := int(weeklyConversionCap.Int32)
		campaign.WeeklyConversionCap = &cap
	}
	
	if monthlyConversionCap.Valid {
		cap := int(monthlyConversionCap.Int32)
		campaign.MonthlyConversionCap = &cap
	}
	
	if globalConversionCap.Valid {
		cap := int(globalConversionCap.Int32)
		campaign.GlobalConversionCap = &cap
	}
	
	if dailyClickCap.Valid {
		cap := int(dailyClickCap.Int32)
		campaign.DailyClickCap = &cap
	}
	
	if weeklyClickCap.Valid {
		cap := int(weeklyClickCap.Int32)
		campaign.WeeklyClickCap = &cap
	}
	
	if monthlyClickCap.Valid {
		cap := int(monthlyClickCap.Int32)
		campaign.MonthlyClickCap = &cap
	}
	
	if globalClickCap.Valid {
		cap := int(globalClickCap.Int32)
		campaign.GlobalClickCap = &cap
	}
	
	if encodedValue.Valid {
		encoded := encodedValue.String
		campaign.EncodedValue = &encoded
	}
	
	if todayClicks.Valid {
		clicks := int(todayClicks.Int32)
		campaign.TodayClicks = &clicks
	}
	
	if todayRevenue.Valid {
		revenue := todayRevenue.String
		campaign.TodayRevenue = &revenue
	}
	
	if timeCreated.Valid {
		created := int(timeCreated.Int32)
		campaign.TimeCreated = &created
	}
	
	if timeSaved.Valid {
		saved := int(timeSaved.Int32)
		campaign.TimeSaved = &saved
	}
	
	// Assign boolean fields
	if isUsingExplicitTermsAndConditions.Valid {
		val := isUsingExplicitTermsAndConditions.Bool
		campaign.IsUsingExplicitTermsAndConditions = &val
	}
	
	if isForceTermsAndConditions.Valid {
		val := isForceTermsAndConditions.Bool
		campaign.IsForceTermsAndConditions = &val
	}
	
	if isCapsEnabled.Valid {
		val := isCapsEnabled.Bool
		campaign.IsCapsEnabled = &val
	}
	
	if isWhitelistCheckEnabled.Valid {
		val := isWhitelistCheckEnabled.Bool
		campaign.IsWhitelistCheckEnabled = &val
	}
	
	if isViewThroughEnabled.Valid {
		val := isViewThroughEnabled.Bool
		campaign.IsViewThroughEnabled = &val
	}
	
	if isDescriptionPlainText.Valid {
		val := isDescriptionPlainText.Bool
		campaign.IsDescriptionPlainText = &val
	}
	
	if isUseDirectLinking.Valid {
		val := isUseDirectLinking.Bool
		campaign.IsUseDirectLinking = &val
	}
	
	return &campaign, nil
}

// UpdateCampaign updates a campaign in the database
func (r *pgxCampaignRepository) UpdateCampaign(ctx context.Context, campaign *domain.Campaign) error {
	query := `UPDATE public.campaigns
              SET name = $1, description = $2, status = $3, start_date = $4, end_date = $5,
                  destination_url = $6, thumbnail_url = $7, preview_url = $8, visibility = $9,
                  currency_id = $10, conversion_method = $11, session_definition = $12, session_duration = $13,
                  caps_timezone_id = $14, project_id = $15, date_live_until = $16, html_description = $17,
                  terms_and_conditions = $18, is_using_explicit_terms_and_conditions = $19,
                  is_force_terms_and_conditions = $20, is_caps_enabled = $21, is_whitelist_check_enabled = $22,
                  is_view_through_enabled = $23, server_side_url = $24, view_through_destination_url = $25,
                  is_description_plain_text = $26, is_use_direct_linking = $27, app_identifier = $28,
                  daily_conversion_cap = $29, weekly_conversion_cap = $30, monthly_conversion_cap = $31,
                  global_conversion_cap = $32, daily_click_cap = $33, weekly_click_cap = $34,
                  monthly_click_cap = $35, global_click_cap = $36, encoded_value = $37,
                  today_clicks = $38, today_revenue = $39, time_created = $40, time_saved = $41
              WHERE campaign_id = $42
              RETURNING updated_at`
	
	// Convert nullable fields to sql.Null types
	var description, destinationURL, thumbnailURL, previewURL, visibility, currencyID sql.NullString
	var conversionMethod, sessionDefinition, projectID, htmlDescription, termsAndConditions sql.NullString
	var serverSideURL, viewThroughDestinationURL, appIdentifier, encodedValue, todayRevenue sql.NullString
	var startDate, endDate, dateLiveUntil sql.NullTime
	var sessionDuration, capsTimezoneID, dailyConversionCap, weeklyConversionCap sql.NullInt32
	var monthlyConversionCap, globalConversionCap, dailyClickCap, weeklyClickCap sql.NullInt32
	var monthlyClickCap, globalClickCap, todayClicks, timeCreated, timeSaved sql.NullInt32
	
	if campaign.Description != nil {
		description = sql.NullString{String: *campaign.Description, Valid: true}
	}
	
	if campaign.StartDate != nil {
		startDate = sql.NullTime{Time: *campaign.StartDate, Valid: true}
	}
	
	if campaign.EndDate != nil {
		endDate = sql.NullTime{Time: *campaign.EndDate, Valid: true}
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
		sessionDuration = sql.NullInt32{Int32: int32(*campaign.SessionDuration), Valid: true}
	}
	
	if campaign.CapsTimezoneID != nil {
		capsTimezoneID = sql.NullInt32{Int32: int32(*campaign.CapsTimezoneID), Valid: true}
	}
	
	if campaign.ProjectID != nil {
		projectID = sql.NullString{String: *campaign.ProjectID, Valid: true}
	}
	
	if campaign.DateLiveUntil != nil {
		dateLiveUntil = sql.NullTime{Time: *campaign.DateLiveUntil, Valid: true}
	}
	
	if campaign.HTMLDescription != nil {
		htmlDescription = sql.NullString{String: *campaign.HTMLDescription, Valid: true}
	}
	
	if campaign.TermsAndConditions != nil {
		termsAndConditions = sql.NullString{String: *campaign.TermsAndConditions, Valid: true}
	}
	
	if campaign.ServerSideURL != nil {
		serverSideURL = sql.NullString{String: *campaign.ServerSideURL, Valid: true}
	}
	
	if campaign.ViewThroughDestinationURL != nil {
		viewThroughDestinationURL = sql.NullString{String: *campaign.ViewThroughDestinationURL, Valid: true}
	}
	
	if campaign.AppIdentifier != nil {
		appIdentifier = sql.NullString{String: *campaign.AppIdentifier, Valid: true}
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
	
	if campaign.EncodedValue != nil {
		encodedValue = sql.NullString{String: *campaign.EncodedValue, Valid: true}
	}
	
	if campaign.TodayClicks != nil {
		todayClicks = sql.NullInt32{Int32: int32(*campaign.TodayClicks), Valid: true}
	}
	
	if campaign.TodayRevenue != nil {
		todayRevenue = sql.NullString{String: *campaign.TodayRevenue, Valid: true}
	}
	
	if campaign.TimeCreated != nil {
		timeCreated = sql.NullInt32{Int32: int32(*campaign.TimeCreated), Valid: true}
	}
	
	if campaign.TimeSaved != nil {
		timeSaved = sql.NullInt32{Int32: int32(*campaign.TimeSaved), Valid: true}
	}
	
	err := r.db.QueryRow(ctx, query, 
		campaign.Name, 
		description, 
		campaign.Status, 
		startDate, 
		endDate,
		destinationURL,
		thumbnailURL,
		previewURL,
		visibility,
		currencyID,
		conversionMethod,
		sessionDefinition,
		sessionDuration,
		capsTimezoneID,
		projectID,
		dateLiveUntil,
		htmlDescription,
		termsAndConditions,
		campaign.IsUsingExplicitTermsAndConditions,
		campaign.IsForceTermsAndConditions,
		campaign.IsCapsEnabled,
		campaign.IsWhitelistCheckEnabled,
		campaign.IsViewThroughEnabled,
		serverSideURL,
		viewThroughDestinationURL,
		campaign.IsDescriptionPlainText,
		campaign.IsUseDirectLinking,
		appIdentifier,
		dailyConversionCap,
		weeklyConversionCap,
		monthlyConversionCap,
		globalConversionCap,
		dailyClickCap,
		weeklyClickCap,
		monthlyClickCap,
		globalClickCap,
		encodedValue,
		todayClicks,
		todayRevenue,
		timeCreated,
		timeSaved,
		campaign.CampaignID,
	).Scan(&campaign.UpdatedAt)
	
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("campaign not found: %w", domain.ErrNotFound)
		}
		return fmt.Errorf("error updating campaign: %w", err)
	}
	
	return nil
}

// ListCampaignsByOrganization retrieves a list of campaigns for an organization with pagination
func (r *pgxCampaignRepository) ListCampaignsByOrganization(ctx context.Context, orgID int64, limit, offset int) ([]*domain.Campaign, error) {
	query := `SELECT campaign_id, organization_id, advertiser_id, name, description, status, start_date, end_date, created_at, updated_at
              FROM public.campaigns
              WHERE organization_id = $1
              ORDER BY campaign_id
              LIMIT $2 OFFSET $3`
	
	return r.listCampaigns(ctx, query, orgID, limit, offset)
}

// ListCampaignsByAdvertiser retrieves a list of campaigns for an advertiser with pagination
func (r *pgxCampaignRepository) ListCampaignsByAdvertiser(ctx context.Context, advertiserID int64, limit, offset int) ([]*domain.Campaign, error) {
	query := `SELECT campaign_id, organization_id, advertiser_id, name, description, status, start_date, end_date, created_at, updated_at
              FROM public.campaigns
              WHERE advertiser_id = $1
              ORDER BY campaign_id
              LIMIT $2 OFFSET $3`
	
	return r.listCampaigns(ctx, query, advertiserID, limit, offset)
}

// listCampaigns is a helper function to list campaigns with a given query
func (r *pgxCampaignRepository) listCampaigns(ctx context.Context, query string, param int64, limit, offset int) ([]*domain.Campaign, error) {
	rows, err := r.db.Query(ctx, query, param, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error listing campaigns: %w", err)
	}
	defer rows.Close()
	
	var campaigns []*domain.Campaign
	for rows.Next() {
		var campaign domain.Campaign
		var description sql.NullString
		var startDate, endDate sql.NullTime
		
		if err := rows.Scan(
			&campaign.CampaignID,
			&campaign.OrganizationID,
			&campaign.AdvertiserID,
			&campaign.Name,
			&description,
			&campaign.Status,
			&startDate,
			&endDate,
			&campaign.CreatedAt,
			&campaign.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("error scanning campaign row: %w", err)
		}
		
		if description.Valid {
			desc := description.String
			campaign.Description = &desc
		}
		
		if startDate.Valid {
			date := startDate.Time
			campaign.StartDate = &date
		}
		
		if endDate.Valid {
			date := endDate.Time
			campaign.EndDate = &date
		}
		
		campaigns = append(campaigns, &campaign)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating campaign rows: %w", err)
	}
	
	return campaigns, nil
}

// DeleteCampaign deletes a campaign from the database
func (r *pgxCampaignRepository) DeleteCampaign(ctx context.Context, id int64) error {
	query := `DELETE FROM public.campaigns WHERE campaign_id = $1`
	
	commandTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting campaign: %w", err)
	}
	
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("campaign not found: %w", domain.ErrNotFound)
	}
	
	return nil
}

// CreateCampaignProviderOffer creates a new campaign provider offer in the database
func (r *pgxCampaignRepository) CreateCampaignProviderOffer(ctx context.Context, offer *domain.CampaignProviderOffer) error {
	query := `INSERT INTO public.campaign_provider_offers 
              (campaign_id, provider_type, provider_offer_ref, provider_offer_config, is_active_on_provider, last_synced_at, created_at, updated_at)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
              RETURNING campaign_provider_offer_id, created_at, updated_at`
	
	var providerOfferRef, providerOfferConfig sql.NullString
	var lastSyncedAt sql.NullTime
	
	if offer.ProviderOfferRef != nil {
		providerOfferRef = sql.NullString{String: *offer.ProviderOfferRef, Valid: true}
	}
	
	if offer.ProviderOfferConfig != nil {
		providerOfferConfig = sql.NullString{String: *offer.ProviderOfferConfig, Valid: true}
	}
	
	if offer.LastSyncedAt != nil {
		lastSyncedAt = sql.NullTime{Time: *offer.LastSyncedAt, Valid: true}
	}
	
	now := time.Now()
	err := r.db.QueryRow(ctx, query, 
		offer.CampaignID, 
		offer.ProviderType, 
		providerOfferRef, 
		providerOfferConfig, 
		offer.IsActiveOnProvider, 
		lastSyncedAt, 
		now, 
		now,
	).Scan(
		&offer.CampaignProviderOfferID, 
		&offer.CreatedAt, 
		&offer.UpdatedAt,
	)
	
	if err != nil {
		return fmt.Errorf("error creating campaign provider offer: %w", err)
	}
	
	return nil
}

// GetCampaignProviderOfferByID retrieves a campaign provider offer by ID
func (r *pgxCampaignRepository) GetCampaignProviderOfferByID(ctx context.Context, id int64) (*domain.CampaignProviderOffer, error) {
	query := `SELECT campaign_provider_offer_id, campaign_id, provider_type, provider_offer_ref, provider_offer_config, 
              is_active_on_provider, last_synced_at, created_at, updated_at
              FROM public.campaign_provider_offers WHERE campaign_provider_offer_id = $1`
	
	var offer domain.CampaignProviderOffer
	var providerOfferRef, providerOfferConfig sql.NullString
	var lastSyncedAt sql.NullTime
	
	err := r.db.QueryRow(ctx, query, id).Scan(
		&offer.CampaignProviderOfferID,
		&offer.CampaignID,
		&offer.ProviderType,
		&providerOfferRef,
		&providerOfferConfig,
		&offer.IsActiveOnProvider,
		&lastSyncedAt,
		&offer.CreatedAt,
		&offer.UpdatedAt,
	)
	
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("campaign provider offer not found: %w", domain.ErrNotFound)
		}
		return nil, fmt.Errorf("error getting campaign provider offer by ID: %w", err)
	}
	
	if providerOfferRef.Valid {
		ref := providerOfferRef.String
		offer.ProviderOfferRef = &ref
	}
	
	if providerOfferConfig.Valid {
		config := providerOfferConfig.String
		offer.ProviderOfferConfig = &config
	}
	
	if lastSyncedAt.Valid {
		synced := lastSyncedAt.Time
		offer.LastSyncedAt = &synced
	}
	
	return &offer, nil
}

// UpdateCampaignProviderOffer updates a campaign provider offer in the database
func (r *pgxCampaignRepository) UpdateCampaignProviderOffer(ctx context.Context, offer *domain.CampaignProviderOffer) error {
	query := `UPDATE public.campaign_provider_offers
              SET provider_offer_ref = $1, provider_offer_config = $2, is_active_on_provider = $3, last_synced_at = $4
              WHERE campaign_provider_offer_id = $5
              RETURNING updated_at`
	
	var providerOfferRef, providerOfferConfig sql.NullString
	var lastSyncedAt sql.NullTime
	
	if offer.ProviderOfferRef != nil {
		providerOfferRef = sql.NullString{String: *offer.ProviderOfferRef, Valid: true}
	}
	
	if offer.ProviderOfferConfig != nil {
		providerOfferConfig = sql.NullString{String: *offer.ProviderOfferConfig, Valid: true}
	}
	
	if offer.LastSyncedAt != nil {
		lastSyncedAt = sql.NullTime{Time: *offer.LastSyncedAt, Valid: true}
	}
	
	err := r.db.QueryRow(ctx, query, 
		providerOfferRef, 
		providerOfferConfig, 
		offer.IsActiveOnProvider, 
		lastSyncedAt, 
		offer.CampaignProviderOfferID,
	).Scan(&offer.UpdatedAt)
	
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("campaign provider offer not found: %w", domain.ErrNotFound)
		}
		return fmt.Errorf("error updating campaign provider offer: %w", err)
	}
	
	return nil
}

// ListCampaignProviderOffersByCampaign retrieves a list of campaign provider offers for a campaign
func (r *pgxCampaignRepository) ListCampaignProviderOffersByCampaign(ctx context.Context, campaignID int64) ([]*domain.CampaignProviderOffer, error) {
	query := `SELECT campaign_provider_offer_id, campaign_id, provider_type, provider_offer_ref, provider_offer_config, 
              is_active_on_provider, last_synced_at, created_at, updated_at
              FROM public.campaign_provider_offers
              WHERE campaign_id = $1
              ORDER BY campaign_provider_offer_id`
	
	rows, err := r.db.Query(ctx, query, campaignID)
	if err != nil {
		return nil, fmt.Errorf("error listing campaign provider offers: %w", err)
	}
	defer rows.Close()
	
	var offers []*domain.CampaignProviderOffer
	for rows.Next() {
		var offer domain.CampaignProviderOffer
		var providerOfferRef, providerOfferConfig sql.NullString
		var lastSyncedAt sql.NullTime
		
		if err := rows.Scan(
			&offer.CampaignProviderOfferID,
			&offer.CampaignID,
			&offer.ProviderType,
			&providerOfferRef,
			&providerOfferConfig,
			&offer.IsActiveOnProvider,
			&lastSyncedAt,
			&offer.CreatedAt,
			&offer.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("error scanning campaign provider offer row: %w", err)
		}
		
		if providerOfferRef.Valid {
			ref := providerOfferRef.String
			offer.ProviderOfferRef = &ref
		}
		
		if providerOfferConfig.Valid {
			config := providerOfferConfig.String
			offer.ProviderOfferConfig = &config
		}
		
		if lastSyncedAt.Valid {
			synced := lastSyncedAt.Time
			offer.LastSyncedAt = &synced
		}
		
		offers = append(offers, &offer)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating campaign provider offer rows: %w", err)
	}
	
	return offers, nil
}

// DeleteCampaignProviderOffer deletes a campaign provider offer from the database
func (r *pgxCampaignRepository) DeleteCampaignProviderOffer(ctx context.Context, id int64) error {
	query := `DELETE FROM public.campaign_provider_offers WHERE campaign_provider_offer_id = $1`
	
	commandTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting campaign provider offer: %w", err)
	}
	
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("campaign provider offer not found: %w", domain.ErrNotFound)
	}
	
	return nil
}