package utils

import (
	"encoding/json"
	"errors"
	"net/http"
)

var (
	ErrEncoding = errors.New("error encoding response")
	ErrDecoding = errors.New("error decoding request body")
)

func SendJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, ErrEncoding.Error(), http.StatusInternalServerError)
	}
}

func DecodeJSON(r *http.Request, v any) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return ErrDecoding
	}
	return nil
}
