package response

import "time"

type APIResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
	Timestamp string      `json:"timestamp"`
	Error     string      `json:"error,omitempty"`
}

func Success(message string, data interface{}) APIResponse {
	return APIResponse{
		Success:   true,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().Format(time.RFC3339),
	}
}

func Error(message string, err string) APIResponse {
	return APIResponse{
		Success:   false,
		Message:   message,
		Timestamp: time.Now().Format(time.RFC3339),
		Error:     err,
	}
}
