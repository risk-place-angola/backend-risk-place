package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"math"

	"github.com/risk-place-angola/backend-risk-place/internal/domain/model"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/repository"
)

const (
	earthRadiusKmDangerZone = 6371.0
	gridSizeKmDefault       = 0.5
	minIncidentsDefault     = 5
	daysBackDefault         = 30
)

type DangerZoneRepoPG struct {
	db *sql.DB
}

func NewDangerZoneRepoPG(db *sql.DB) repository.DangerZoneRepository {
	return &DangerZoneRepoPG{db: db}
}

func (r *DangerZoneRepoPG) CalculateDangerZones(ctx context.Context, gridSizeKm float64, minIncidents int, daysBack int) ([]repository.DangerZoneCalculationResult, error) {
	query := `
		WITH incident_grid AS (
			SELECT 
				FLOOR(latitude / $1) * $1 AS grid_lat,
				FLOOR(longitude / $1) * $1 AS grid_lon,
				COUNT(*) AS incident_count,
				SUM(
					CASE 
						WHEN created_at > NOW() - INTERVAL '7 days' THEN 3.0
						WHEN created_at > NOW() - INTERVAL '14 days' THEN 2.0
						WHEN created_at > NOW() - INTERVAL '30 days' THEN 1.5
						ELSE 1.0
					END *
					CASE 
						WHEN status = 'verified' THEN 2.0
						ELSE 1.0
					END
				) AS risk_score
			FROM reports
			WHERE 
				created_at > NOW() - INTERVAL '1 day' * $3
				AND status IN ('verified', 'pending')
			GROUP BY grid_lat, grid_lon
			HAVING COUNT(*) >= $2
		)
		SELECT 
			CONCAT(grid_lat::text, ',', grid_lon::text) AS grid_cell_id,
			grid_lat + ($1 / 2) AS cell_center_lat,
			grid_lon + ($1 / 2) AS cell_center_lon,
			incident_count,
			LEAST(risk_score, 10.0) AS risk_score
		FROM incident_grid
		ORDER BY risk_score DESC
		LIMIT 1000
	`

	rows, err := r.db.QueryContext(ctx, query, gridSizeKm, minIncidents, daysBack)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate danger zones: %w", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			err = closeErr
		}
	}()

	var results []repository.DangerZoneCalculationResult
	for rows.Next() {
		var result repository.DangerZoneCalculationResult
		err := rows.Scan(
			&result.GridCellID,
			&result.CellLat,
			&result.CellLon,
			&result.IncidentCount,
			&result.RiskScore,
		)
		if err != nil {
			continue
		}
		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating danger zone rows: %w", err)
	}

	return results, nil
}

func (r *DangerZoneRepoPG) GetNearbyDangerZones(ctx context.Context, lat, lon, radiusKm float64) ([]*model.DangerZone, error) {
	gridSize := gridSizeKmDefault
	latMin := lat - (radiusKm / kmPerDegreeLat)
	latMax := lat + (radiusKm / kmPerDegreeLat)
	lonMin := lon - (radiusKm / (kmPerDegreeLat * math.Cos(lat*math.Pi/degreesInCircle)))
	lonMax := lon + (radiusKm / (kmPerDegreeLat * math.Cos(lat*math.Pi/degreesInCircle)))

	query := `
		WITH incident_grid AS (
			SELECT 
				FLOOR(latitude / $1) * $1 AS grid_lat,
				FLOOR(longitude / $1) * $1 AS grid_lon,
				COUNT(*) AS incident_count,
				SUM(
					CASE 
						WHEN created_at > NOW() - INTERVAL '7 days' THEN 3.0
						WHEN created_at > NOW() - INTERVAL '14 days' THEN 2.0
						WHEN created_at > NOW() - INTERVAL '30 days' THEN 1.5
						ELSE 1.0
					END *
					CASE 
						WHEN status = 'verified' THEN 2.0
						ELSE 1.0
					END
				) AS risk_score
			FROM reports
			WHERE 
				created_at > NOW() - INTERVAL '30 days'
				AND status IN ('verified', 'pending')
				AND latitude BETWEEN $2 AND $3
				AND longitude BETWEEN $4 AND $5
			GROUP BY grid_lat, grid_lon
			HAVING COUNT(*) >= $6
		)
		SELECT 
			CONCAT(grid_lat::text, ',', grid_lon::text) AS grid_cell_id,
			grid_lat + ($1 / 2) AS cell_center_lat,
			grid_lon + ($1 / 2) AS cell_center_lon,
			incident_count,
			LEAST(risk_score, 10.0) AS risk_score
		FROM incident_grid
	`

	rows, err := r.db.QueryContext(ctx, query, gridSize, latMin, latMax, lonMin, lonMax, minIncidentsDefault)
	if err != nil {
		return nil, fmt.Errorf("failed to get nearby danger zones: %w", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			err = closeErr
		}
	}()

	var zones []*model.DangerZone
	for rows.Next() {
		var gridCellID string
		var cellLat, cellLon float64
		var incidentCount int
		var riskScore float64

		err := rows.Scan(&gridCellID, &cellLat, &cellLon, &incidentCount, &riskScore)
		if err != nil {
			continue
		}

		zone := model.NewDangerZone(cellLat, cellLon, gridCellID)
		zone.IncidentCount = incidentCount
		zone.RiskScore = riskScore
		zone.CalculateRiskLevel()

		zones = append(zones, zone)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating nearby danger zones: %w", err)
	}

	return zones, nil
}
