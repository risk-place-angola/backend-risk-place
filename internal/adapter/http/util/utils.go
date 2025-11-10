package util

import (
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"strconv"
)

func ConvertStringToFloat64(value string) float64 {
	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		fmt.Errorf("error converting string to float64: %v", err)
		return 0.0
	}
	return floatValue
}

func ParseUUIDParam(r *http.Request, param string) (uuid.UUID, error) {
	value := r.PathValue(param)
	if value == "" {
		return uuid.Nil, fmt.Errorf("missing %s parameter", param)
	}
	id, err := uuid.Parse(value)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid %s parameter: %w", param, err)
	}
	return id, nil
}

func ParseUUIDQuery(r *http.Request, param string) (uuid.UUID, error) {
	value := r.URL.Query().Get(param)
	if value == "" {
		return uuid.Nil, fmt.Errorf("missing %s query parameter", param)
	}
	id, err := uuid.Parse(value)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid %s query parameter: %w", param, err)
	}
	return id, nil
}
