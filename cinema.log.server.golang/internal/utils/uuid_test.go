package utils

import (
	"testing"

	"github.com/google/uuid"
)

func TestParseUUID_Success(t *testing.T) {
	validUUID := uuid.New()
	parsed, err := ParseUUID(validUUID.String())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if parsed != validUUID {
		t.Errorf("expected UUID %v, got %v", validUUID, parsed)
	}
}

func TestParseUUID_EmptyString(t *testing.T) {
	_, err := ParseUUID("")
	if err == nil {
		t.Fatal("expected error for empty string")
	}
}

func TestParseUUID_InvalidFormat(t *testing.T) {
	_, err := ParseUUID("not-a-uuid")
	if err == nil {
		t.Fatal("expected error for invalid format")
	}
}

func TestValidateUUID(t *testing.T) {
	validUUID := uuid.New().String()
	if !ValidateUUID(validUUID) {
		t.Error("expected ValidateUUID to return true for valid UUID")
	}

	if ValidateUUID("invalid") {
		t.Error("expected ValidateUUID to return false for invalid UUID")
	}
}

func TestGenerateUUID(t *testing.T) {
	id := GenerateUUID()
	if id == uuid.Nil {
		t.Error("GenerateUUID returned nil UUID")
	}
}

func TestIsNilUUID(t *testing.T) {
	if !IsNilUUID(uuid.Nil) {
		t.Error("expected IsNilUUID to return true for uuid.Nil")
	}

	validUUID := uuid.New()
	if IsNilUUID(validUUID) {
		t.Error("expected IsNilUUID to return false for valid UUID")
	}
}
