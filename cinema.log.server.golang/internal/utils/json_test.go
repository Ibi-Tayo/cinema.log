package utils

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSendJSON(t *testing.T) {
	w := httptest.NewRecorder()
	data := map[string]string{"message": "hello"}
	
	SendJSON(w, data)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", contentType)
	}

	if !bytes.Contains(w.Body.Bytes(), []byte("hello")) {
		t.Error("expected body to contain 'hello'")
	}
}

func TestDecodeJSON_Success(t *testing.T) {
	type TestStruct struct {
		Name string `json:"name"`
	}

	jsonBody := `{"name":"John"}`
	req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBufferString(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	var result TestStruct
	err := DecodeJSON(req, &result)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Name != "John" {
		t.Errorf("expected name John, got %s", result.Name)
	}
}

func TestDecodeJSON_InvalidJSON(t *testing.T) {
	jsonBody := `{"name":"John"`
	req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBufferString(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	var result map[string]interface{}
	err := DecodeJSON(req, &result)

	if err != ErrDecoding {
		t.Errorf("expected ErrDecoding, got %v", err)
	}
}
