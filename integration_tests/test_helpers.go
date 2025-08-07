package integration_tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

const (
	BaseURL        = "http://localhost:50222"
	EverflowAPIURL = "https://api.eflow.team/v1"
	JWTSecret      = "gDxsm/JerlPJiOObQLtfjViLBQF2ggmJpYCNW+9LPwL2QJksmiYlzRCJCKseCLxJtGysx+awZvoiS0MF0pLjnw=="
)

// TestConfig holds configuration for integration tests
type TestConfig struct {
	EverflowAPIKey string
	BaseURL        string
	HTTPClient     *http.Client
	JWTToken       string
	UserID         string
}

// NewTestConfig creates a new test configuration
func NewTestConfig() *TestConfig {
	jwtToken, userID := GenerateTestJWT()
	return &TestConfig{
		EverflowAPIKey: os.Getenv("EVERFLOW_API_KEY"),
		BaseURL:        BaseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		JWTToken: jwtToken,
		UserID:   userID,
	}
}

// GenerateTestJWT generates a valid JWT token for testing with a random UUID
func GenerateTestJWT() (string, string) {
	// Generate a random UUID for this test session
	userID := uuid.New().String()
	
	// Create the claims - use a valid UUID for the subject
	claims := jwt.MapClaims{
		"sub":        userID,
		"email":      "test@example.com",
		"session_id": "test-session-12345",
		"iat":        time.Now().Unix(),
		"exp":        time.Now().Add(time.Hour * 24).Unix(), // 24 hours from now
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret
	tokenString, err := token.SignedString([]byte(JWTSecret))
	if err != nil {
		panic(fmt.Sprintf("Failed to generate test JWT: %v", err))
	}

	return tokenString, userID
}

// APIRequest represents a generic API request
type APIRequest struct {
	Method  string
	URL     string
	Body    interface{}
	Headers map[string]string
}

// APIResponse represents a generic API response
type APIResponse struct {
	StatusCode int
	Body       []byte
	Headers    http.Header
}

// MakeAPIRequest makes an HTTP request and returns the response
func (tc *TestConfig) MakeAPIRequest(t *testing.T, req APIRequest) *APIResponse {
	var bodyReader io.Reader
	if req.Body != nil {
		bodyBytes, err := json.Marshal(req.Body)
		require.NoError(t, err)
		bodyReader = bytes.NewReader(bodyBytes)
	}

	httpReq, err := http.NewRequest(req.Method, req.URL, bodyReader)
	require.NoError(t, err)

	// Set default headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	// Set custom headers
	for key, value := range req.Headers {
		httpReq.Header.Set(key, value)
	}

	resp, err := tc.HTTPClient.Do(httpReq)
	require.NoError(t, err)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return &APIResponse{
		StatusCode: resp.StatusCode,
		Body:       body,
		Headers:    resp.Header,
	}
}

// EverflowAPIRequest makes a request to the Everflow API
func (tc *TestConfig) EverflowAPIRequest(t *testing.T, method, endpoint string, body interface{}) *APIResponse {
	return tc.MakeAPIRequest(t, APIRequest{
		Method: method,
		URL:    EverflowAPIURL + endpoint,
		Body:   body,
		Headers: map[string]string{
			"X-Eflow-API-Key": tc.EverflowAPIKey,
		},
	})
}

// PlatformAPIRequest makes a request to our platform API
func (tc *TestConfig) PlatformAPIRequest(t *testing.T, method, endpoint string, body interface{}) *APIResponse {
	return tc.MakeAPIRequest(t, APIRequest{
		Method: method,
		URL:    tc.BaseURL + endpoint,
		Body:   body,
		Headers: map[string]string{
			"Authorization": "Bearer " + tc.JWTToken,
		},
	})
}

// PlatformPublicAPIRequest makes a request to our platform's public API (no auth required)
func (tc *TestConfig) PlatformPublicAPIRequest(t *testing.T, method, endpoint string, body interface{}) *APIResponse {
	return tc.MakeAPIRequest(t, APIRequest{
		Method: method,
		URL:    tc.BaseURL + endpoint,
		Body:   body,
		Headers: map[string]string{
			// No Authorization header for public endpoints
		},
	})
}

// WaitForSync waits for synchronization to complete
func (tc *TestConfig) WaitForSync(t *testing.T, maxWait time.Duration) {
	time.Sleep(2 * time.Second) // Give some time for async operations
}

// ParseJSONResponse parses a JSON response into the provided struct
func ParseJSONResponse(t *testing.T, resp *APIResponse, target interface{}) {
	err := json.Unmarshal(resp.Body, target)
	require.NoError(t, err, "Failed to parse JSON response: %s", string(resp.Body))
}

// LogResponse logs the API response for debugging
func LogResponse(t *testing.T, label string, resp *APIResponse) {
	t.Logf("%s Response - Status: %d, Body: %s", label, resp.StatusCode, string(resp.Body))
}

// GenerateTestName generates a unique test name with timestamp
func GenerateTestName(prefix string) string {
	return fmt.Sprintf("%s_test_%d", prefix, time.Now().Unix())
}

// GenerateTestEmail generates a unique test email
func GenerateTestEmail(prefix string) string {
	return fmt.Sprintf("%s_test_%d@example.com", prefix, time.Now().Unix())
}

// GenerateTestURL generates a unique test URL
func GenerateTestURL(prefix string) string {
	return fmt.Sprintf("https://%s-test-%d.example.com", prefix, time.Now().Unix())
}

// AssertSuccessResponse asserts that the response indicates success
func AssertSuccessResponse(t *testing.T, resp *APIResponse, expectedStatus int) {
	require.Equal(t, expectedStatus, resp.StatusCode, 
		"Expected status %d but got %d. Response: %s", 
		expectedStatus, resp.StatusCode, string(resp.Body))
}

// AssertErrorResponse asserts that the response indicates an error
func AssertErrorResponse(t *testing.T, resp *APIResponse, expectedStatus int) {
	require.Equal(t, expectedStatus, resp.StatusCode,
		"Expected error status %d but got %d. Response: %s",
		expectedStatus, resp.StatusCode, string(resp.Body))
}

// ExtractEverflowIDFromMapping extracts the Everflow ID from a provider mapping response
func ExtractEverflowIDFromMapping(t *testing.T, resp *APIResponse) int {
	// Parse the response as a generic map to handle different entity types
	var response map[string]interface{}
	ParseJSONResponse(t, resp, &response)
	
	// Handle different response structures - some have provider_mapping wrapper, others don't
	var providerMapping map[string]interface{}
	if mapping, exists := response["provider_mapping"].(map[string]interface{}); exists {
		// Advertiser format: {"provider_mapping": {...}}
		providerMapping = mapping
	} else {
		// Affiliate format: direct response without wrapper
		providerMapping = response
	}
	
	// Try to find the provider ID field - could be advertiser, affiliate, or campaign/offer
	var providerID string
	if advertiserID, exists := providerMapping["provider_advertiser_id"].(string); exists && advertiserID != "" {
		providerID = advertiserID
	} else if affiliateID, exists := providerMapping["provider_affiliate_id"].(string); exists && affiliateID != "" {
		providerID = affiliateID
	} else if campaignID, exists := providerMapping["provider_campaign_id"].(string); exists && campaignID != "" {
		providerID = campaignID
	} else if offerID, exists := providerMapping["provider_offer_id"].(string); exists && offerID != "" {
		providerID = offerID
	} else {
		t.Fatalf("Could not find provider ID field in mapping: %+v", providerMapping)
	}
	
	// Convert string ID to int
	everflowID, err := strconv.Atoi(providerID)
	if err != nil {
		t.Fatalf("Failed to convert provider ID '%s' to int: %v", providerID, err)
	}
	return everflowID
}

// EverflowAPIResponse represents a response from Everflow API
type EverflowAPIResponse struct {
	StatusCode int
	Body       string
	Headers    map[string][]string
}

// callEverflowAPI makes a direct call to Everflow API
func callEverflowAPI(t *testing.T, config *TestConfig, method, endpoint string, payload interface{}) *EverflowAPIResponse {
	baseURL := "https://api.eflow.team"
	url := baseURL + endpoint

	var body io.Reader
	if payload != nil {
		jsonData, err := json.Marshal(payload)
		if err != nil {
			t.Fatalf("Failed to marshal payload: %v", err)
		}
		body = bytes.NewBuffer(jsonData)
		t.Logf("Request payload: %s", string(jsonData))
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Eflow-API-Key", config.EverflowAPIKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	return &EverflowAPIResponse{
		StatusCode: resp.StatusCode,
		Body:       string(bodyBytes),
		Headers:    resp.Header,
	}
}