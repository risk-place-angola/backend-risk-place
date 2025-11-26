package handler

import (
	"context"
	"log/slog"
	"time"

	"github.com/risk-place-angola/backend-risk-place/internal/domain/service"
)

const dangerZoneRecalculationInterval = 30 * time.Minute

func StartDangerZoneCalculationJob(ctx context.Context, dangerZoneService service.DangerZoneService) {
	go func() {
		ticker := time.NewTicker(dangerZoneRecalculationInterval)
		defer ticker.Stop()

		slog.Info("starting danger zone calculation job", "interval", dangerZoneRecalculationInterval)

		if err := dangerZoneService.CalculateDangerZones(ctx); err != nil {
			slog.Error("initial danger zone calculation failed", "error", err)
		} else {
			slog.Info("initial danger zone calculation completed")
		}

		for {
			select {
			case <-ctx.Done():
				slog.Info("danger zone calculation job stopped")
				return
			case <-ticker.C:
				slog.Debug("recalculating danger zones")
				if err := dangerZoneService.CalculateDangerZones(ctx); err != nil {
					slog.Error("danger zone calculation failed", "error", err)
				} else {
					slog.Info("danger zone calculation completed")
				}
			}
		}
	}()
}
