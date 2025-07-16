package handlers

// Common response messages
const (
	// Success messages
	MsgListCreated        = "Favorite publisher list created successfully"
	MsgListsRetrieved     = "Favorite publisher lists retrieved successfully"
	MsgListRetrieved      = "Favorite publisher list retrieved successfully"
	MsgListUpdated        = "Favorite publisher list updated successfully"
	MsgListDeleted        = "Favorite publisher list deleted successfully"
	MsgPublisherAdded     = "Publisher added to list successfully"
	MsgPublisherRemoved   = "Publisher removed from list successfully"
	MsgPublisherUpdated   = "Publisher notes updated successfully"
	MsgItemsRetrieved     = "List items retrieved successfully"
	MsgListsWithPublisher = "Lists containing publisher retrieved successfully"
	MsgPublisherRetrieved = "Publisher retrieved successfully"

	// Error messages
	ErrUnauthorized              = "Unauthorized"
	ErrInvalidRequestBody        = "Invalid request body"
	ErrInvalidInput              = "Invalid input"
	ErrInternalServer            = "Internal server error"
	ErrListNotFound              = "List not found"
	ErrPublisherNotFound         = "Publisher not found"
	ErrPublisherNotInList        = "Publisher not found in list"
	ErrInvalidListID             = "Invalid list ID"
	ErrInvalidDomain             = "Invalid domain"
	ErrFailedToCreateList        = "Failed to create favorite publisher list"
	ErrFailedToRetrieveLists     = "Failed to retrieve favorite publisher lists"
	ErrFailedToRetrieveList      = "Failed to retrieve favorite publisher list"
	ErrFailedToUpdateList        = "Failed to update favorite publisher list"
	ErrFailedToDeleteList        = "Failed to delete favorite publisher list"
	ErrFailedToAddPublisher      = "Failed to add publisher to list"
	ErrFailedToRemovePublisher   = "Failed to remove publisher from list"
	ErrFailedToUpdatePublisher   = "Failed to update publisher in list"
	ErrFailedToRetrieveItems     = "Failed to retrieve list items"
	ErrFailedToRetrievePublisher = "Failed to retrieve publisher"

	// Error details
	DetailOrgIDNotFound       = "Organization ID not found in context"
	DetailInvalidOrgIDType    = "Invalid organization ID type"
	DetailListIDRequired      = "List ID must be a valid integer"
	DetailDomainRequired      = "Domain parameter is required"
	DetailDomainQueryRequired = "Domain query parameter is required"
	DetailListNotBelongToOrg  = "List does not belong to your organization"
	DetailPublisherNotInList  = "No publisher found with the specified domain in this list"
)
