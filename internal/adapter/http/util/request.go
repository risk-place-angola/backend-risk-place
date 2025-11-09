package util

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Success bool `json:"success"`
	Error   struct {
		Message interface{} `json:"message"`
		Code    int         `json:"code"`
	} `json:"error"`
}

func Error(w http.ResponseWriter, message interface{}, code int) {
	var payload = ErrorResponse{
		Success: false,
		Error: struct {
			Message interface{} `json:"message"`
			Code    int         `json:"code"`
		}{
			Message: message,
			Code:    code,
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
