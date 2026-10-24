package models

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

func WriteSuccess(w http.ResponseWriter, status int, message string, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(APIResponse{
		Message: message,
		Data:    data,
	})
}
func WriteError(w http.ResponseWriter, status int, errMessage string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(APIResponse{
		Message: "request failed",
		Error:   errMessage,
	})
}
