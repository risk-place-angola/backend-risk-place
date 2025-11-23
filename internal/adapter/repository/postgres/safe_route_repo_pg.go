package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/risk-place-angola/backend-risk-place/internal/adapter/repository/postgres/sqlc"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/model"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/repository"
)

const (
	earthRadiusKm         = 6371.0
	corridorWidthKm       = 2.0
	avgSpeedKmh           = 40.0
	highRiskIncidentTypes = "robbery,assault,armed_robbery"
	gridCellSizeKm        = 0.5

	// Time and distance constants
	minutesPerHour           = 60.0
	hoursPerDay              = 24.0
	degreesInCircle          = 180.0
	kmPerDegreeLat           = 111.0
	maxIntermediatePoints    = 10
	intermediatePointSpacing = 2.0
	halfDivisor              = 2.0

	// Temporal decay constants
	recentIncidentDays   = 7
	mediumIncidentDays   = 30
	oldIncidentDays      = 60
	proximityThresholdKm = 0.5

	// Weight calculation constants
	baseWeightDivisor = 10.0
	maxWeightCap      = 10.0
	minWaypointCount  = 2
)

type SafeRouteRepoPG struct {
	q  sqlc.Querier
	db *sql.DB
}

func NewSafeRouteRepoPG(db *sql.DB) repository.SafeRouteRepository {
	return &SafeRouteRepoPG{
		q:  sqlc.New(db),
		db: db,
	}
}

func (r *SafeRouteRepoPG) CalculateSafeRoute(ctx context.Context, params repository.RouteCalculationParams) (*model.SafeRoute, error) {
	route := model.NewSafeRoute(params.OriginLat, params.OriginLon, params.DestinationLat, params.DestinationLon)

	waypoints := r.generateWaypoints(params.OriginLat, params.OriginLon, params.DestinationLat, params.DestinationLon)
	for i, wp := range waypoints {
		route.AddWaypoint(wp.Latitude, wp.Longitude, i)
	}

	route.DistanceKm = r.calculateTotalDistance(waypoints)
	route.EstimatedDuration = int((route.DistanceKm / avgSpeedKmh) * minutesPerHour)

	incidents, err := r.GetIncidentsForRoute(ctx, route.Waypoints, corridorWidthKm)
	if err != nil {
		return nil, fmt.Errorf("failed to get incidents for route: %w", err)
	}

	for _, incident := range incidents {
		route.AddIncident(incident)
	}

	route.CalculateSafetyScore()

	return route, nil
}

func (r *SafeRouteRepoPG) GetIncidentsForRoute(ctx context.Context, waypoints []model.Waypoint, corridorWidthKm float64) ([]model.IncidentNearRoute, error) {
	if len(waypoints) == 0 {
		return []model.IncidentNearRoute{}, nil
	}

	minLat, maxLat, minLon, maxLon := r.calculateBoundingBox(waypoints, corridorWidthKm)

	lineString := r.buildLineString(waypoints)
	// #nosec G202 - lineString is generated from validated waypoints, not user input
	query := `
		SELECT 
			r.id,
			rt.name as risk_type,
			rtopic.name as risk_topic,
			r.latitude,
			r.longitude,
			r.created_at,
			ST_Distance(
				ST_MakePoint(r.longitude, r.latitude)::geography,
				ST_MakeLine(ARRAY[` + lineString + `])::geography
			) / 1000.0 as distance_km
		FROM reports r
		JOIN risk_types rt ON r.risk_type_id = rt.id
		JOIN risk_topics rtopic ON r.risk_topic_id = rtopic.id
		WHERE r.status = 'verified'
			AND r.latitude BETWEEN $1 AND $2
			AND r.longitude BETWEEN $3 AND $4
			AND r.created_at > NOW() - INTERVAL '90 days'
		ORDER BY distance_km ASC
		LIMIT 50
	`

	rows, err := r.db.QueryContext(ctx, query, minLat, maxLat, minLon, maxLon)
	if err != nil {
		return nil, fmt.Errorf("failed to query incidents: %w", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			err = closeErr
		}
	}()

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	incidents := []model.IncidentNearRoute{}
	now := time.Now()

	for rows.Next() {
		var incident model.IncidentNearRoute
		var createdAt time.Time
		var riskType, riskTopic string

		err := rows.Scan(
			&incident.ReportID,
			&riskType,
			&riskTopic,
			&incident.Latitude,
			&incident.Longitude,
			&createdAt,
			&incident.DistanceKm,
		)
		if err != nil {
			continue
		}

		incident.RiskType = riskType
		incident.RiskTopic = riskTopic
		incident.CreatedAt = createdAt
		incident.DaysAgo = int(now.Sub(createdAt).Hours() / hoursPerDay)
		incident.WeightFactor = r.calculateIncidentWeight(riskType, incident.DaysAgo, incident.DistanceKm)

		if incident.DistanceKm <= corridorWidthKm {
			incidents = append(incidents, incident)
		}
	}

	return incidents, nil
}

func (r *SafeRouteRepoPG) GetIncidentsHeatmap(ctx context.Context, params repository.IncidentHeatmapParams) ([]repository.HeatmapPoint, error) {
	query := `
		SELECT 
			r.latitude,
			r.longitude,
			rt.name as risk_type,
			COUNT(*) as report_count
		FROM reports r
		JOIN risk_types rt ON r.risk_type_id = rt.id
		WHERE r.status = 'verified'
			AND r.latitude BETWEEN $1 AND $2
			AND r.longitude BETWEEN $3 AND $4
	`

	args := []interface{}{params.SouthWestLat, params.NorthEastLat, params.SouthWestLon, params.NorthEastLon}
	argIndex := 5

	if params.StartDate != "" {
		query += fmt.Sprintf(" AND r.created_at >= $%d", argIndex)
		args = append(args, params.StartDate)
		argIndex++
	}

	if params.EndDate != "" {
		query += fmt.Sprintf(" AND r.created_at <= $%d", argIndex)
		args = append(args, params.EndDate)
		argIndex++
	}

	if params.RiskTypeID != "" {
		query += fmt.Sprintf(" AND r.risk_type_id = $%d", argIndex)
		riskTypeUUID, err := uuid.Parse(params.RiskTypeID)
		if err == nil {
			args = append(args, riskTypeUUID)
		}
	}

	query += `
		GROUP BY r.latitude, r.longitude, rt.name
		ORDER BY report_count DESC
		LIMIT 500
	`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query heatmap: %w", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			err = closeErr
		}
	}()

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	points := []repository.HeatmapPoint{}

	for rows.Next() {
		var point repository.HeatmapPoint
		var reportCount int

		err := rows.Scan(
			&point.Latitude,
			&point.Longitude,
			&point.IncidentType,
			&reportCount,
		)
		if err != nil {
			continue
		}

		point.ReportCount = reportCount
		point.Weight = r.calculateHeatmapWeight(point.IncidentType, reportCount)
		points = append(points, point)
	}

	return points, nil
}

func (r *SafeRouteRepoPG) generateWaypoints(originLat, originLon, destLat, destLon float64) []model.Waypoint {
	waypoints := []model.Waypoint{
		{Latitude: originLat, Longitude: originLon, Sequence: 0},
	}

	distance := r.haversineDistance(originLat, originLon, destLat, destLon)
	numIntermediatePoints := int(distance / intermediatePointSpacing)

	if numIntermediatePoints > maxIntermediatePoints {
		numIntermediatePoints = maxIntermediatePoints
	}

	if numIntermediatePoints > 0 {
		for i := 1; i <= numIntermediatePoints; i++ {
			ratio := float64(i) / float64(numIntermediatePoints+1)
			lat := originLat + ratio*(destLat-originLat)
			lon := originLon + ratio*(destLon-originLon)
			waypoints = append(waypoints, model.Waypoint{
				Latitude:  lat,
				Longitude: lon,
				Sequence:  i,
			})
		}
	}

	waypoints = append(waypoints, model.Waypoint{
		Latitude:  destLat,
		Longitude: destLon,
		Sequence:  len(waypoints),
	})

	return waypoints
}

func (r *SafeRouteRepoPG) haversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	lat1Rad := lat1 * math.Pi / degreesInCircle
	lat2Rad := lat2 * math.Pi / degreesInCircle
	deltaLat := (lat2 - lat1) * math.Pi / degreesInCircle
	deltaLon := (lon2 - lon1) * math.Pi / degreesInCircle

	a := math.Sin(deltaLat/halfDivisor)*math.Sin(deltaLat/halfDivisor) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLon/halfDivisor)*math.Sin(deltaLon/halfDivisor)

	c := halfDivisor * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadiusKm * c
}

func (r *SafeRouteRepoPG) calculateTotalDistance(waypoints []model.Waypoint) float64 {
	if len(waypoints) < minWaypointCount {
		return 0
	}

	totalDistance := 0.0
	for i := range len(waypoints) - 1 {
		distance := r.haversineDistance(
			waypoints[i].Latitude, waypoints[i].Longitude,
			waypoints[i+1].Latitude, waypoints[i+1].Longitude,
		)
		totalDistance += distance
	}

	return totalDistance
}

func (r *SafeRouteRepoPG) calculateBoundingBox(waypoints []model.Waypoint, bufferKm float64) (float64, float64, float64, float64) {
	if len(waypoints) == 0 {
		return 0, 0, 0, 0
	}

	minLat := waypoints[0].Latitude
	maxLat := waypoints[0].Latitude
	minLon := waypoints[0].Longitude
	maxLon := waypoints[0].Longitude

	for _, wp := range waypoints {
		if wp.Latitude < minLat {
			minLat = wp.Latitude
		}
		if wp.Latitude > maxLat {
			maxLat = wp.Latitude
		}
		if wp.Longitude < minLon {
			minLon = wp.Longitude
		}
		if wp.Longitude > maxLon {
			maxLon = wp.Longitude
		}
	}

	latBuffer := bufferKm / kmPerDegreeLat
	lonBuffer := bufferKm / (kmPerDegreeLat * math.Cos(minLat*math.Pi/degreesInCircle))

	minLat -= latBuffer
	maxLat += latBuffer
	minLon -= lonBuffer
	maxLon += lonBuffer

	return minLat, maxLat, minLon, maxLon
}

func (r *SafeRouteRepoPG) buildLineString(waypoints []model.Waypoint) string {
	if len(waypoints) == 0 {
		return ""
	}

	lineString := ""
	for i, wp := range waypoints {
		if i > 0 {
			lineString += ", "
		}
		lineString += fmt.Sprintf("ST_MakePoint(%f, %f)", wp.Longitude, wp.Latitude)
	}

	return lineString
}

func (r *SafeRouteRepoPG) calculateIncidentWeight(riskType string, daysAgo int, distanceKm float64) float64 {
	weight := 1.0

	switch {
	case daysAgo <= recentIncidentDays:
		weight *= 3.0
	case daysAgo <= mediumIncidentDays:
		weight *= 2.0
	case daysAgo <= oldIncidentDays:
		weight *= 1.5
	default:
		weight *= 1.0
	}

	if containsRiskType(riskType, highRiskIncidentTypes) {
		weight *= 2.5
	} else {
		weight *= 1.0
	}

	if distanceKm <= proximityThresholdKm {
		weight *= 2.0
	} else if distanceKm <= 1.0 {
		weight *= 1.5
	}

	return weight
}

func (r *SafeRouteRepoPG) calculateHeatmapWeight(incidentType string, reportCount int) float64 {
	baseWeight := float64(reportCount) / baseWeightDivisor

	if containsRiskType(incidentType, highRiskIncidentTypes) {
		baseWeight *= 1.5
	}

	if baseWeight > maxWeightCap {
		baseWeight = maxWeightCap
	}

	return baseWeight
}

func containsRiskType(riskType, types string) bool {
	return len(riskType) > 0 && len(types) > 0
}
