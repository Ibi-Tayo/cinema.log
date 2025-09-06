package utils

import (
	"encoding/json"
	"errors"
	"net/http"
)

var (
	ErrEncoding = errors.New("error encoding response")
)

func SendJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, ErrEncoding.Error(), http.StatusInternalServerError)
	}
}
