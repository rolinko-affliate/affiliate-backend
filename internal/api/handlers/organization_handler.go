package handlers

import (
	"net/http"
	"strconv"

	"github.com/affiliate-backend/internal/service"
	"github.com/gin-gonic/gin"
)

// OrganizationHandler handles organization-related requests
type OrganizationHandler struct {
	organizationService service.OrganizationService
}

// NewOrganizationHandler creates a new organization handler
func NewOrganizationHandler(os service.OrganizationService) *OrganizationHandler {
	return &OrganizationHandler{organizationService: os}
}

// CreateOrganizationRequest defines the request for creating an organization
type CreateOrganizationRequest struct {
	Name string `json:"name" binding:"required"`
}

// CreateOrganization creates a new organization
// @Summary      Create a new organization
// @Description  Creates a new organization with the given name
// @Tags         organizations
// @Accept       json
// @Produce      json
// @Param        request  body      CreateOrganizationRequest  true  "Organization details"
// @Success      201      {object}  domain.Organization        "Created organization"
// @Failure      400      {object}  map[string]string          "Invalid request"
// @Failure      500      {object}  map[string]string          "Internal server error"
// @Security     BearerAuth
// @Router       /organizations [post]
func (h *OrganizationHandler) CreateOrganization(c *gin.Context) {
	var req CreateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	organization, err := h.organizationService.CreateOrganization(c.Request.Context(), req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create organization: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, organization)
}

// GetOrganization retrieves an organization by ID
// @Summary      Get organization by ID
// @Description  Retrieves an organization by its ID
// @Tags         organizations
// @Accept       json
// @Produce      json
// @Param        id   path      int                   true  "Organization ID"
// @Success      200  {object}  domain.Organization  "Organization details"
// @Failure      400  {object}  map[string]string    "Invalid organization ID"
// @Failure      404  {object}  map[string]string    "Organization not found"
// @Failure      500  {object}  map[string]string    "Internal server error"
// @Security     BearerAuth
// @Router       /organizations/{id} [get]
func (h *OrganizationHandler) GetOrganization(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	organization, err := h.organizationService.GetOrganizationByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "organization not found: not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Organization not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get organization: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, organization)
}

// UpdateOrganizationRequest defines the request for updating an organization
type UpdateOrganizationRequest struct {
	Name string `json:"name" binding:"required"`
}

// UpdateOrganization updates an organization
// @Summary      Update organization
// @Description  Updates an organization with the given details
// @Tags         organizations
// @Accept       json
// @Produce      json
// @Param        id       path      int                       true  "Organization ID"
// @Param        request  body      UpdateOrganizationRequest  true  "Organization details"
// @Success      200      {object}  domain.Organization        "Updated organization"
// @Failure      400      {object}  map[string]string          "Invalid request"
// @Failure      404      {object}  map[string]string          "Organization not found"
// @Failure      500      {object}  map[string]string          "Internal server error"
// @Security     BearerAuth
// @Router       /organizations/{id} [put]
func (h *OrganizationHandler) UpdateOrganization(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	var req UpdateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// Get existing organization
	organization, err := h.organizationService.GetOrganizationByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "organization not found: not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Organization not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get organization: " + err.Error()})
		return
	}

	// Update organization
	organization.Name = req.Name
	if err := h.organizationService.UpdateOrganization(c.Request.Context(), organization); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update organization: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, organization)
}

// ListOrganizations retrieves a list of organizations
// @Summary      List organizations
// @Description  Retrieves a list of organizations with pagination
// @Tags         organizations
// @Accept       json
// @Produce      json
// @Param        page      query     int                      false  "Page number (default: 1)"
// @Param        pageSize  query     int                      false  "Page size (default: 10)"
// @Success      200       {array}   domain.Organization      "List of organizations"
// @Failure      500       {object}  map[string]string        "Internal server error"
// @Security     BearerAuth
// @Router       /organizations [get]
func (h *OrganizationHandler) ListOrganizations(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	organizations, err := h.organizationService.ListOrganizations(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list organizations: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, organizations)
}

// DeleteOrganization deletes an organization
// @Summary      Delete organization
// @Description  Deletes an organization by its ID
// @Tags         organizations
// @Accept       json
// @Produce      json
// @Param        id   path      int                true  "Organization ID"
// @Success      204  {object}  nil                "No content"
// @Failure      400  {object}  map[string]string  "Invalid organization ID"
// @Failure      404  {object}  map[string]string  "Organization not found"
// @Failure      500  {object}  map[string]string  "Internal server error"
// @Security     BearerAuth
// @Router       /organizations/{id} [delete]
func (h *OrganizationHandler) DeleteOrganization(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	if err := h.organizationService.DeleteOrganization(c.Request.Context(), id); err != nil {
		if err.Error() == "organization not found: not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Organization not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete organization: " + err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}