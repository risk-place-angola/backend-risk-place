package model

import (
	"time"

	"github.com/google/uuid"
)

type DangerZone struct {
	ID           uuid.UUID `json:"id"`
	CellLat      float64   `json:"cell_lat"`
	CellLon      float64   `json:"cell_lon"`
	GridCellID   string    `json:"grid_cell_id"`
	IncidentCount int       `json:"incident_count"`
	RiskScore    float64   `json:"risk_score"`
	RiskLevel    string    `json:"risk_level"`
	CalculatedAt time.Time `json:"calculated_at"`
	ExpiresAt    time.Time `json:"expires_at"`
}

func NewDangerZone(cellLat, cellLon float64, gridCellID string) *DangerZone {
	return &DangerZone{
		ID:           uuid.New(),
		CellLat:      cellLat,
		CellLon:      cellLon,
		GridCellID:   gridCellID,
		IncidentCount: 0,
		RiskScore:    0,
		RiskLevel:    "low",
		CalculatedAt: time.Now(),
		ExpiresAt:    time.Now().Add(1 * time.Hour),
	}
}

const (
	riskScoreCritical = 8.0
	riskScoreHigh     = 6.0
	riskScoreMedium   = 4.0
)

func (dz *DangerZone) CalculateRiskLevel() {
	switch {
	case dz.RiskScore >= riskScoreCritical:
		dz.RiskLevel = "critical"
	case dz.RiskScore >= riskScoreHigh:
		dz.RiskLevel = "high"
	case dz.RiskScore >= riskScoreMedium:
		dz.RiskLevel = "medium"
	default:
		dz.RiskLevel = "low"
	}
}

func (dz *DangerZone) IsExpired() bool {
	return time.Now().After(dz.ExpiresAt)
}
