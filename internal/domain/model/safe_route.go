package model

import (
	"time"

	"github.com/google/uuid"
)

type RiskLevel string

const (
	RiskLevelVeryLow  RiskLevel = "very_low"
	RiskLevelLow      RiskLevel = "low"
	RiskLevelModerate RiskLevel = "moderate"
	RiskLevelHigh     RiskLevel = "high"
	RiskLevelVeryHigh RiskLevel = "very_high"

	// Safety score calculation constants
	perfectScore          = 100.0
	penaltyMultiplier     = 10.0
	maxPenalty            = 90.0
	veryLowRiskThreshold  = 80.0
	lowRiskThreshold      = 60.0
	moderateRiskThreshold = 40.0
	highRiskThreshold     = 20.0
)

type Waypoint struct {
	Latitude  float64
	Longitude float64
	Sequence  int
}

type IncidentNearRoute struct {
	ReportID     uuid.UUID
	RiskType     string
	RiskTopic    string
	Latitude     float64
	Longitude    float64
	DistanceKm   float64
	CreatedAt    time.Time
	DaysAgo      int
	WeightFactor float64
}

type SafeRoute struct {
	ID                uuid.UUID
	OriginLat         float64
	OriginLon         float64
	DestinationLat    float64
	DestinationLon    float64
	Waypoints         []Waypoint
	DistanceKm        float64
	EstimatedDuration int
	SafetyScore       float64
	RiskLevel         RiskLevel
	IncidentCount     int
	Incidents         []IncidentNearRoute
	CalculatedAt      time.Time
}

func NewSafeRoute(originLat, originLon, destLat, destLon float64) *SafeRoute {
	return &SafeRoute{
		ID:             uuid.New(),
		OriginLat:      originLat,
		OriginLon:      originLon,
		DestinationLat: destLat,
		DestinationLon: destLon,
		Waypoints:      []Waypoint{},
		Incidents:      []IncidentNearRoute{},
		CalculatedAt:   time.Now(),
	}
}

func (sr *SafeRoute) CalculateSafetyScore() {
	if sr.IncidentCount == 0 {
		sr.SafetyScore = perfectScore
		sr.RiskLevel = RiskLevelVeryLow
		return
	}

	totalWeight := 0.0
	for _, incident := range sr.Incidents {
		totalWeight += incident.WeightFactor
	}

	penalty := totalWeight * penaltyMultiplier
	if penalty > maxPenalty {
		penalty = maxPenalty
	}

	sr.SafetyScore = perfectScore - penalty

	switch {
	case sr.SafetyScore >= veryLowRiskThreshold:
		sr.RiskLevel = RiskLevelVeryLow
	case sr.SafetyScore >= lowRiskThreshold:
		sr.RiskLevel = RiskLevelLow
	case sr.SafetyScore >= moderateRiskThreshold:
		sr.RiskLevel = RiskLevelModerate
	case sr.SafetyScore >= highRiskThreshold:
		sr.RiskLevel = RiskLevelHigh
	default:
		sr.RiskLevel = RiskLevelVeryHigh
	}
}

func (sr *SafeRoute) AddWaypoint(lat, lon float64, sequence int) {
	sr.Waypoints = append(sr.Waypoints, Waypoint{
		Latitude:  lat,
		Longitude: lon,
		Sequence:  sequence,
	})
}

func (sr *SafeRoute) AddIncident(incident IncidentNearRoute) {
	sr.Incidents = append(sr.Incidents, incident)
	sr.IncidentCount = len(sr.Incidents)
}
