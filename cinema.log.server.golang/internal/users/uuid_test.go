package users_test

import (
	"context"
	"database/sql"
	"testing"

	"cinema.log.server.golang/internal/domain"
	"cinema.log.server.golang/internal/utils"
)

// Mock database for testing
type mockDB struct{}

func (m *mockDB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return nil, nil
}

func (m *mockDB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return nil
}

func (m *mockDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return nil, nil
}

func TestUUIDGeneration(t *testing.T) {
	// Test UUID generation
	id1 := utils.GenerateUUID()
	id2 := utils.GenerateUUID()
	
	if id1 == id2 {
		t.Error("Generated UUIDs should be unique")
	}
	
	if utils.IsNilUUID(id1) {
		t.Error("Generated UUID should not be nil")
	}
	
	if utils.IsNilUUID(id2) {
		t.Error("Generated UUID should not be nil")
	}
}

func TestUUIDValidation(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"550e8400-e29b-41d4-a716-446655440000", true},
		{"550e8400e29b41d4a716446655440000", true}, // Without hyphens is also valid
		{"invalid-uuid", false},
		{"", false},
		{"550e8400-e29b-41d4-a716-44665544000", false}, // Wrong length
		{"zzz", false}, // Completely invalid
		{"550e8400-e29b-41d4-a716-44665544000g", false}, // Invalid character
	}
	
	for _, test := range tests {
		result := utils.ValidateUUID(test.input)
		if result != test.expected {
			t.Errorf("ValidateUUID(%s) = %t, expected %t", test.input, result, test.expected)
		}
	}
}

func TestUUIDParsing(t *testing.T) {
	validUUID := "550e8400-e29b-41d4-a716-446655440000"
	
	parsed, err := utils.ParseUUID(validUUID)
	if err != nil {
		t.Errorf("ParseUUID(%s) returned error: %v", validUUID, err)
	}
	
	if parsed.String() != validUUID {
		t.Errorf("Parsed UUID %s doesn't match original %s", parsed.String(), validUUID)
	}
	
	// Test invalid UUID
	_, err = utils.ParseUUID("invalid-uuid")
	if err == nil {
		t.Error("ParseUUID should return error for invalid UUID")
	}
}

func TestUserStructWithUUID(t *testing.T) {
	user := &domain.User{
		ID:            utils.GenerateUUID(),
		GithubId:      789776,
		Name:          "Test User",
		Username:      "testuser",
		ProfilePicURL: "https://example.com/pic.jpg",
	}
	
	if utils.IsNilUUID(user.ID) {
		t.Error("User ID should not be nil")
	}
	
	if user.GithubId != 789776 {
		t.Error("GithubID not set correctly")
	}
}
