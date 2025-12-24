package users

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"cinema.log.server.golang/internal/domain"
	"cinema.log.server.golang/internal/utils"
	"github.com/google/uuid"
)

// Test environment variables
var (
	testDB      *sql.DB
	testHandler *Handler
	testService UserService
	testStore   Store
	testDbSetup *utils.TestDatabase
)

// Setup and tear down happen in TestMain (testing.M)
func TestMain(m *testing.M) {
	var err error
	testDbSetup, err = utils.StartTestPostgres()
	if err != nil {
		log.Fatalf("could not start test database: %v", err)
	}

	testDB = testDbSetup.DB

	// Set up test dependencies
	testStore = NewStore(testDB)
	testService = NewService(testStore)
	testHandler = NewHandler(testService)

	// Run tests
	code := m.Run()

	// Cleanup
	testDbSetup.Close()

	os.Exit(code)
}

func TestCreateUserIntegration(t *testing.T) {
	// Create test user
	testUser := &domain.User{
		Name:          "Test User",
		Username:      "testuser",
		GithubId:      12345,
		ProfilePicURL: "https://example.com/avatar.jpg",
	}

	// Convert to JSON
	userJSON, err := json.Marshal(testUser)
	if err != nil {
		t.Fatalf("failed to marshal user: %v", err)
	}

	// Create HTTP request
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(userJSON))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()

	// Call handler
	testHandler.CreateUser(w, req)

	// Check response
	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
		t.Errorf("response body: %s", w.Body.String())
	}

	// Parse response
	var createdUser domain.User
	if err := json.NewDecoder(w.Body).Decode(&createdUser); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	// Verify user was created
	if createdUser.Name != testUser.Name {
		t.Errorf("expected name %s, got %s", testUser.Name, createdUser.Name)
	}
	if createdUser.Username != testUser.Username {
		t.Errorf("expected username %s, got %s", testUser.Username, createdUser.Username)
	}
	if createdUser.ID == uuid.Nil {
		t.Error("expected non-nil UUID")
	}
}

func TestGetUserByIdIntegration(t *testing.T) {
	// First create a user
	testUser := &domain.User{
		Name:          "Get Test User",
		Username:      "gettestuser",
		GithubId:      54321,
		ProfilePicURL: "https://example.com/avatar2.jpg",
	}

	createdUser, err := testService.CreateUser(context.Background(), testUser)
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	// Create HTTP request
	req := httptest.NewRequest(http.MethodGet, "/users/"+createdUser.ID.String(), nil)
	req.SetPathValue("id", createdUser.ID.String())

	// Create response recorder
	w := httptest.NewRecorder()

	// Call handler
	testHandler.GetUserById(w, req)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		t.Errorf("response body: %s", w.Body.String())
	}

	// Parse response
	var retrievedUser domain.User
	if err := json.NewDecoder(w.Body).Decode(&retrievedUser); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	// Verify user data
	if retrievedUser.ID != createdUser.ID {
		t.Errorf("expected ID %s, got %s", createdUser.ID, retrievedUser.ID)
	}
	if retrievedUser.Name != testUser.Name {
		t.Errorf("expected name %s, got %s", testUser.Name, retrievedUser.Name)
	}
}

func TestGetAllUsersIntegration(t *testing.T) {
	// Create multiple test users
	users := []*domain.User{
		{Name: "User One", Username: "userone", GithubId: 111},
		{Name: "User Two", Username: "usertwo", GithubId: 222},
	}

	// Create users in database
	for _, user := range users {
		_, err := testService.CreateUser(context.Background(), user)
		if err != nil {
			t.Fatalf("failed to create test user: %v", err)
		}
	}

	// Create HTTP request
	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()

	// Call handler
	testHandler.GetAllUsers(w, req)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Parse response
	var retrievedUsers []*domain.User
	if err := json.NewDecoder(w.Body).Decode(&retrievedUsers); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	// Should have at least the users we created (may have more from other tests)
	if len(retrievedUsers) < 2 {
		t.Errorf("expected at least 2 users, got %d", len(retrievedUsers))
	}
}

func TestUpdateUserIntegration(t *testing.T) {
	// Create initial user
	testUser := &domain.User{
		Name:          "Original Name",
		Username:      "originaluser",
		GithubId:      99999,
		ProfilePicURL: "https://example.com/original.jpg",
	}

	createdUser, err := testService.CreateUser(context.Background(), testUser)
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	// Update user data
	updatedUser := *createdUser
	updatedUser.Name = "Updated Name"
	updatedUser.Username = "updateduser"

	// Convert to JSON
	userJSON, err := json.Marshal(updatedUser)
	if err != nil {
		t.Fatalf("failed to marshal user: %v", err)
	}

	// Create HTTP request
	req := httptest.NewRequest(http.MethodPut, "/users", bytes.NewBuffer(userJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	// Call handler
	testHandler.UpdateUser(w, req)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		t.Errorf("response body: %s", w.Body.String())
	}

	// Verify update in database
	retrievedUser, err := testService.GetUserById(context.Background(), createdUser.ID)
	if err != nil {
		t.Fatalf("failed to retrieve updated user: %v", err)
	}

	if retrievedUser.Name != "Updated Name" {
		t.Errorf("expected updated name, got %s", retrievedUser.Name)
	}
}

func TestDeleteUserIntegration(t *testing.T) {
	// Create user to delete
	testUser := &domain.User{
		Name:          "Delete Me",
		Username:      "deleteme",
		GithubId:      77777,
		ProfilePicURL: "https://example.com/delete.jpg",
	}

	createdUser, err := testService.CreateUser(context.Background(), testUser)
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	// Create delete request
	req := httptest.NewRequest(http.MethodDelete, "/users/"+createdUser.ID.String(), nil)
	req.SetPathValue("id", createdUser.ID.String())

	w := httptest.NewRecorder()

	// Call handler
	testHandler.DeleteUser(w, req)

	// Check response
	if w.Code != http.StatusNoContent {
		t.Errorf("expected status %d, got %d", http.StatusNoContent, w.Code)
	}

	// Verify user was deleted
	_, err = testService.GetUserById(context.Background(), createdUser.ID)
	if err != ErrUserNotFound {
		t.Errorf("expected user to be deleted, but got error: %v", err)
	}
}

// Test validation errors
func TestCreateUserValidationIntegration(t *testing.T) {
	// Test with invalid name (too short)
	testUser := &domain.User{
		Name:          "Bad", // Too short
		Username:      "baduser",
		GithubId:      88888,
		ProfilePicURL: "https://example.com/bad.jpg",
	}

	userJSON, _ := json.Marshal(testUser)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(userJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	testHandler.CreateUser(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d for validation error, got %d", http.StatusBadRequest, w.Code)
	}
}
