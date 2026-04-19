package shared

import (
	"encoding/json"
	"net/http"
)

type ResponseError struct {
	Error      bool   `json:"error"`
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}

type ResponseSuccess struct {
	Error      bool   `json:"error"`
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
	Data       any    `json:"data"`
}

func WriteError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ResponseError{
		Error:      true,
		Message:    message,
		StatusCode: statusCode,
	})
}

func WriteSuccess(w http.ResponseWriter, message string, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ResponseSuccess{
		Error:      false,
		Message:    message,
		Data:       data,
		StatusCode: http.StatusOK,
	})
}
