package repository

import (
	"context"
	"testing"
	"time"

	"github.com/affiliate-backend/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFavoritePublisherListRepository_CreateList(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewFavoritePublisherListRepository(db)
	ctx := context.Background()

	list := &domain.FavoritePublisherList{
		OrganizationID: 1,
		Name:           "Test List",
		Description:    stringPtr("Test description"),
	}

	err := repo.CreateList(ctx, list)
	require.NoError(t, err)
	assert.NotZero(t, list.ListID)
	assert.NotZero(t, list.CreatedAt)
	assert.NotZero(t, list.UpdatedAt)
}

func TestFavoritePublisherListRepository_GetListByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewFavoritePublisherListRepository(db)
	ctx := context.Background()

	// Create a test list
	list := &domain.FavoritePublisherList{
		OrganizationID: 1,
		Name:           "Test List",
		Description:    stringPtr("Test description"),
	}
	err := repo.CreateList(ctx, list)
	require.NoError(t, err)

	// Retrieve the list
	retrieved, err := repo.GetListByID(ctx, list.ListID)
	require.NoError(t, err)
	assert.Equal(t, list.ListID, retrieved.ListID)
	assert.Equal(t, list.OrganizationID, retrieved.OrganizationID)
	assert.Equal(t, list.Name, retrieved.Name)
	assert.Equal(t, list.Description, retrieved.Description)

	// Test non-existent list
	_, err = repo.GetListByID(ctx, 99999)
	assert.Equal(t, domain.ErrNotFound, err)
}

func TestFavoritePublisherListRepository_GetListsByOrganization(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewFavoritePublisherListRepository(db)
	ctx := context.Background()

	// Create test lists for organization 1
	list1 := &domain.FavoritePublisherList{
		OrganizationID: 1,
		Name:           "List 1",
		Description:    stringPtr("Description 1"),
	}
	list2 := &domain.FavoritePublisherList{
		OrganizationID: 1,
		Name:           "List 2",
		Description:    stringPtr("Description 2"),
	}
	list3 := &domain.FavoritePublisherList{
		OrganizationID: 2, // Different organization
		Name:           "List 3",
		Description:    stringPtr("Description 3"),
	}

	err := repo.CreateList(ctx, list1)
	require.NoError(t, err)
	err = repo.CreateList(ctx, list2)
	require.NoError(t, err)
	err = repo.CreateList(ctx, list3)
	require.NoError(t, err)

	// Add a publisher to list1 to test the count
	item := &domain.FavoritePublisherListItem{
		ListID:          list1.ListID,
		PublisherDomain: "example.com",
		Notes:           stringPtr("Test notes"),
	}
	err = repo.AddPublisherToList(ctx, item)
	require.NoError(t, err)

	// Get lists for organization 1
	lists, err := repo.GetListsByOrganization(ctx, 1)
	require.NoError(t, err)
	assert.Len(t, lists, 2)

	// Find list1 and check publisher count
	var foundList *domain.FavoritePublisherListWithStats
	for _, l := range lists {
		if l.ListID == list1.ListID {
			foundList = l
			break
		}
	}
	require.NotNil(t, foundList)
	assert.Equal(t, int64(1), foundList.PublisherCount)

	// Get lists for organization 2
	lists, err = repo.GetListsByOrganization(ctx, 2)
	require.NoError(t, err)
	assert.Len(t, lists, 1)
	assert.Equal(t, list3.ListID, lists[0].ListID)
}

func TestFavoritePublisherListRepository_UpdateList(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewFavoritePublisherListRepository(db)
	ctx := context.Background()

	// Create a test list
	list := &domain.FavoritePublisherList{
		OrganizationID: 1,
		Name:           "Original Name",
		Description:    stringPtr("Original description"),
	}
	err := repo.CreateList(ctx, list)
	require.NoError(t, err)

	originalUpdatedAt := list.UpdatedAt

	// Wait a bit to ensure timestamp difference
	time.Sleep(10 * time.Millisecond)

	// Update the list
	list.Name = "Updated Name"
	list.Description = stringPtr("Updated description")
	err = repo.UpdateList(ctx, list)
	require.NoError(t, err)
	assert.True(t, list.UpdatedAt.After(originalUpdatedAt))

	// Verify the update
	retrieved, err := repo.GetListByID(ctx, list.ListID)
	require.NoError(t, err)
	assert.Equal(t, "Updated Name", retrieved.Name)
	assert.Equal(t, "Updated description", *retrieved.Description)

	// Test updating non-existent list
	nonExistentList := &domain.FavoritePublisherList{
		ListID: 99999,
		Name:   "Non-existent",
	}
	err = repo.UpdateList(ctx, nonExistentList)
	assert.Equal(t, domain.ErrNotFound, err)
}

func TestFavoritePublisherListRepository_DeleteList(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewFavoritePublisherListRepository(db)
	ctx := context.Background()

	// Create a test list
	list := &domain.FavoritePublisherList{
		OrganizationID: 1,
		Name:           "Test List",
		Description:    stringPtr("Test description"),
	}
	err := repo.CreateList(ctx, list)
	require.NoError(t, err)

	// Delete the list
	err = repo.DeleteList(ctx, list.ListID)
	require.NoError(t, err)

	// Verify it's deleted
	_, err = repo.GetListByID(ctx, list.ListID)
	assert.Equal(t, domain.ErrNotFound, err)

	// Test deleting non-existent list
	err = repo.DeleteList(ctx, 99999)
	assert.Equal(t, domain.ErrNotFound, err)
}

func TestFavoritePublisherListRepository_AddPublisherToList(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewFavoritePublisherListRepository(db)
	ctx := context.Background()

	// Create a test list
	list := &domain.FavoritePublisherList{
		OrganizationID: 1,
		Name:           "Test List",
		Description:    stringPtr("Test description"),
	}
	err := repo.CreateList(ctx, list)
	require.NoError(t, err)

	// Add a publisher to the list
	item := &domain.FavoritePublisherListItem{
		ListID:          list.ListID,
		PublisherDomain: "example.com",
		Notes:           stringPtr("Test notes"),
	}
	err = repo.AddPublisherToList(ctx, item)
	require.NoError(t, err)
	assert.NotZero(t, item.ItemID)
	assert.NotZero(t, item.AddedAt)
}

func TestFavoritePublisherListRepository_RemovePublisherFromList(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewFavoritePublisherListRepository(db)
	ctx := context.Background()

	// Create a test list and add a publisher
	list := &domain.FavoritePublisherList{
		OrganizationID: 1,
		Name:           "Test List",
		Description:    stringPtr("Test description"),
	}
	err := repo.CreateList(ctx, list)
	require.NoError(t, err)

	item := &domain.FavoritePublisherListItem{
		ListID:          list.ListID,
		PublisherDomain: "example.com",
		Notes:           stringPtr("Test notes"),
	}
	err = repo.AddPublisherToList(ctx, item)
	require.NoError(t, err)

	// Remove the publisher
	err = repo.RemovePublisherFromList(ctx, list.ListID, "example.com")
	require.NoError(t, err)

	// Verify it's removed
	items, err := repo.GetListItems(ctx, list.ListID)
	require.NoError(t, err)
	assert.Len(t, items, 0)

	// Test removing non-existent publisher
	err = repo.RemovePublisherFromList(ctx, list.ListID, "nonexistent.com")
	assert.Equal(t, domain.ErrNotFound, err)
}

func TestFavoritePublisherListRepository_GetListItems(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewFavoritePublisherListRepository(db)
	ctx := context.Background()

	// Create a test list
	list := &domain.FavoritePublisherList{
		OrganizationID: 1,
		Name:           "Test List",
		Description:    stringPtr("Test description"),
	}
	err := repo.CreateList(ctx, list)
	require.NoError(t, err)

	// Add publishers to the list
	item1 := &domain.FavoritePublisherListItem{
		ListID:          list.ListID,
		PublisherDomain: "example1.com",
		Notes:           stringPtr("Notes 1"),
	}
	item2 := &domain.FavoritePublisherListItem{
		ListID:          list.ListID,
		PublisherDomain: "example2.com",
		Notes:           stringPtr("Notes 2"),
	}

	err = repo.AddPublisherToList(ctx, item1)
	require.NoError(t, err)
	err = repo.AddPublisherToList(ctx, item2)
	require.NoError(t, err)

	// Get list items
	items, err := repo.GetListItems(ctx, list.ListID)
	require.NoError(t, err)
	assert.Len(t, items, 2)

	// Items should be ordered by added_at DESC (newest first)
	assert.Equal(t, "example2.com", items[0].PublisherDomain)
	assert.Equal(t, "example1.com", items[1].PublisherDomain)
}

func TestFavoritePublisherListRepository_IsPublisherInList(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewFavoritePublisherListRepository(db)
	ctx := context.Background()

	// Create a test list
	list := &domain.FavoritePublisherList{
		OrganizationID: 1,
		Name:           "Test List",
		Description:    stringPtr("Test description"),
	}
	err := repo.CreateList(ctx, list)
	require.NoError(t, err)

	// Check non-existent publisher
	exists, err := repo.IsPublisherInList(ctx, list.ListID, "example.com")
	require.NoError(t, err)
	assert.False(t, exists)

	// Add a publisher
	item := &domain.FavoritePublisherListItem{
		ListID:          list.ListID,
		PublisherDomain: "example.com",
		Notes:           stringPtr("Test notes"),
	}
	err = repo.AddPublisherToList(ctx, item)
	require.NoError(t, err)

	// Check existing publisher
	exists, err = repo.IsPublisherInList(ctx, list.ListID, "example.com")
	require.NoError(t, err)
	assert.True(t, exists)
}

func TestFavoritePublisherListRepository_UpdatePublisherInList(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewFavoritePublisherListRepository(db)
	ctx := context.Background()

	// Create a test list and add a publisher
	list := &domain.FavoritePublisherList{
		OrganizationID: 1,
		Name:           "Test List",
		Description:    stringPtr("Test description"),
	}
	err := repo.CreateList(ctx, list)
	require.NoError(t, err)

	item := &domain.FavoritePublisherListItem{
		ListID:          list.ListID,
		PublisherDomain: "example.com",
		Notes:           stringPtr("Original notes"),
	}
	err = repo.AddPublisherToList(ctx, item)
	require.NoError(t, err)

	// Update the notes
	newNotes := "Updated notes"
	err = repo.UpdatePublisherInList(ctx, list.ListID, "example.com", &newNotes)
	require.NoError(t, err)

	// Verify the update
	items, err := repo.GetListItems(ctx, list.ListID)
	require.NoError(t, err)
	assert.Len(t, items, 1)
	assert.Equal(t, "Updated notes", *items[0].Notes)

	// Test updating non-existent publisher
	err = repo.UpdatePublisherInList(ctx, list.ListID, "nonexistent.com", &newNotes)
	assert.Equal(t, domain.ErrNotFound, err)
}

func TestFavoritePublisherListRepository_GetListsContainingPublisher(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewFavoritePublisherListRepository(db)
	ctx := context.Background()

	// Create test lists
	list1 := &domain.FavoritePublisherList{
		OrganizationID: 1,
		Name:           "List 1",
		Description:    stringPtr("Description 1"),
	}
	list2 := &domain.FavoritePublisherList{
		OrganizationID: 1,
		Name:           "List 2",
		Description:    stringPtr("Description 2"),
	}
	list3 := &domain.FavoritePublisherList{
		OrganizationID: 2, // Different organization
		Name:           "List 3",
		Description:    stringPtr("Description 3"),
	}

	err := repo.CreateList(ctx, list1)
	require.NoError(t, err)
	err = repo.CreateList(ctx, list2)
	require.NoError(t, err)
	err = repo.CreateList(ctx, list3)
	require.NoError(t, err)

	// Add the same publisher to list1 and list3
	item1 := &domain.FavoritePublisherListItem{
		ListID:          list1.ListID,
		PublisherDomain: "example.com",
		Notes:           stringPtr("Notes 1"),
	}
	item3 := &domain.FavoritePublisherListItem{
		ListID:          list3.ListID,
		PublisherDomain: "example.com",
		Notes:           stringPtr("Notes 3"),
	}

	err = repo.AddPublisherToList(ctx, item1)
	require.NoError(t, err)
	err = repo.AddPublisherToList(ctx, item3)
	require.NoError(t, err)

	// Get lists containing the publisher for organization 1
	lists, err := repo.GetListsContainingPublisher(ctx, 1, "example.com")
	require.NoError(t, err)
	assert.Len(t, lists, 1)
	assert.Equal(t, list1.ListID, lists[0].ListID)

	// Get lists containing the publisher for organization 2
	lists, err = repo.GetListsContainingPublisher(ctx, 2, "example.com")
	require.NoError(t, err)
	assert.Len(t, lists, 1)
	assert.Equal(t, list3.ListID, lists[0].ListID)

	// Get lists containing a non-existent publisher
	lists, err = repo.GetListsContainingPublisher(ctx, 1, "nonexistent.com")
	require.NoError(t, err)
	assert.Len(t, lists, 0)
}

// Helper function to create a string pointer
func stringPtr(s string) *string {
	return &s
}