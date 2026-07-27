package global

type GlobalResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type SuccessResponse struct {
	GlobalResponse
	Data interface{} `json:"data,omitempty"`
}

type ErrorDetail struct {
	Code string `json:"code"`
}

type ErrorResponse struct {
	GlobalResponse
	Error *ErrorDetail `json:"error,omitempty"`
}
