package database

import (
	"context"
	"log"
	"os"
	"testing"

	"cinema.log.server.golang/internal/utils"
)

var testDbSetup *utils.TestDatabase

func TestMain(m *testing.M) {
	var err error
	testDbSetup, err = utils.StartTestPostgres()
	if err != nil {
		log.Fatalf("could not start test database: %v", err)
	}

	code := m.Run()
	testDbSetup.Close()
	os.Exit(code)
}

func TestNew(t *testing.T) {
	srv := New()
	if srv == nil {
		t.Fatal("New() returned nil")
	}
}

// func TestHealth(t *testing.T) {
// 	srv := New()

// 	stats := srv.Health()

// 	if stats["status"] != "up" {
// 		t.Fatalf("expected status to be up, got %s", stats["status"])
// 	}

// 	if _, ok := stats["error"]; ok {
// 		t.Fatalf("expected error not to be present")
// 	}

// 	if stats["message"] != "It's healthy" {
// 		t.Fatalf("expected message to be 'It's healthy', got %s", stats["message"])
// 	}
// }

func TestClose(t *testing.T) {
	srv := New()

	if srv.Close() != nil {
		t.Fatalf("expected Close() to return nil")
	}
}

func TestService_Health(t *testing.T) {
	svc := &service{db: testDbSetup.DB}

	stats := svc.Health()

	if stats["status"] != "up" {
		t.Errorf("expected status to be 'up', got %s", stats["status"])
	}

	if _, ok := stats["error"]; ok {
		t.Error("expected error not to be present when database is healthy")
	}

	// Check that stats contain expected keys
	expectedKeys := []string{"status", "message", "open_connections", "in_use", "idle"}
	for _, key := range expectedKeys {
		if _, ok := stats[key]; !ok {
			t.Errorf("expected stats to contain key %s", key)
		}
	}
}

func TestService_Query(t *testing.T) {
	svc := &service{db: testDbSetup.DB}

	// Simple query to test Query method
	rows, err := svc.Query(context.Background(), "SELECT 1 as test_value")
	if err != nil {
		t.Fatalf("expected no error from Query, got %v", err)
	}
	defer rows.Close()

	// Verify we can read the result
	if !rows.Next() {
		t.Fatal("expected at least one row from query")
	}

	var testValue int
	if err := rows.Scan(&testValue); err != nil {
		t.Fatalf("expected no error scanning row, got %v", err)
	}

	if testValue != 1 {
		t.Errorf("expected test_value to be 1, got %d", testValue)
	}
}

func TestService_Close(t *testing.T) {
	// Create a new connection just for this test
	db := New()
	svc := &service{db: db}

	err := svc.Close()
	if err != nil {
		t.Errorf("expected no error from Close, got %v", err)
	}
}
