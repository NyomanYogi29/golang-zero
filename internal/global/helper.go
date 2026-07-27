package global

import (
	"encoding/json"
	"net/http"
)

func NewSuccessResponse(message string, data interface{}) SuccessResponse {
	return SuccessResponse{
		GlobalResponse: GlobalResponse{
			Success: true,
			Message: message,
		},
		Data: data,
	}
}

func NewErrorResponse(message string, errorCode string) ErrorResponse {
	var errDetail *ErrorDetail
	if errorCode != "" {
		errDetail = &ErrorDetail{
			Code: errorCode,
		}
	}

	return ErrorResponse{
		GlobalResponse: GlobalResponse{
			Success: false,
			Message: message,
		},
		Error: errDetail,
	}
}

func WriteJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}
