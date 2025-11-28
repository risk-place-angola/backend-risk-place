package util

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Success bool `json:"success"`
	Error   struct {
		Message   interface{} `json:"message"`
		Code      int         `json:"code"`
		ErrorCode string      `json:"error_code,omitempty"`
	} `json:"error"`
}

func Error(w http.ResponseWriter, message interface{}, code int) {
	ErrorWithCode(w, message, code, "")
}

func ErrorWithCode(w http.ResponseWriter, message interface{}, code int, errorCode string) {
	var payload = ErrorResponse{
		Success: false,
		Error: struct {
			Message   interface{} `json:"message"`
			Code      int         `json:"code"`
			ErrorCode string      `json:"error_code,omitempty"`
		}{
			Message:   message,
			Code:      code,
			ErrorCode: errorCode,
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(payload); err != nil {
		http.Error(w, `{"data":"failed to encode error response"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = w.Write(buf.Bytes())
}

type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func ResponseWithMessage(w http.ResponseWriter, message string, data interface{}, code int) {
	payload := SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(payload); err != nil {
		http.Error(w, `{"success":false,"error":{"message":"failed to encode response","code":500}}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = w.Write(buf.Bytes())
}
