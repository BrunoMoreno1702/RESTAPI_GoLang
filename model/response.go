package model

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type APIResponse struct {
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
}

func SuccessResponse(data interface{}, message string) APIResponse {
	return APIResponse{
		Data:    data,
		Message: message,
	}
}

func ErrorResponse(code string, message string) APIResponse {
	return APIResponse{
		Error: &APIError{
			Code:    code,
			Message: message,
		},
	}
}
