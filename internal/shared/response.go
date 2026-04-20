package shared

import (
	"encoding/json"
	"net/http"
)

type BaseResponse struct {
	Error      bool   `json:"error"`
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}

type ResponseError struct {
	BaseResponse
}

type ResponseSuccess struct {
	BaseResponse
	Data any `json:"data"`
}

type ResponseErrorValidation struct {
	BaseResponse
	Field map[string]string `json:"fields,omitempty"`
}

func WriteError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ResponseError{
		BaseResponse: BaseResponse{
			Error:      true,
			Message:    message,
			StatusCode: statusCode,
		},
	})
}

func WriteSuccess(w http.ResponseWriter, message string, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ResponseSuccess{
		BaseResponse: BaseResponse{
			Error:      false,
			Message:    message,
			StatusCode: http.StatusOK,
		},
		Data: data,
	})
}

func WriteValidationError(w http.ResponseWriter, statusCode int, validationErrors map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ResponseErrorValidation{
		BaseResponse: BaseResponse{
			Error:      true,
			Message:    "There is something wrong from your request",
			StatusCode: statusCode,
		},
		Field: validationErrors,
	})

}
