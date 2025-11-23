package util

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func Response(w http.ResponseWriter, object interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")

	if object == nil && code == http.StatusOK {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("[]"))
		return
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(object); err != nil {
		http.Error(w, `{"data":"failed to encode response"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(code)
	_, _ = w.Write(buf.Bytes())
}
