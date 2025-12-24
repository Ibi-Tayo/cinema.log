package database

import (
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
