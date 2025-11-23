package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/model"
)

type ReportVerificationService interface {
	VoteReport(ctx context.Context, reportID uuid.UUID, userID *uuid.UUID, anonymousSessionID *uuid.UUID, voteType model.VoteType) error
	CheckDuplicates(ctx context.Context, lat, lon float64, riskTypeID uuid.UUID) ([]*model.Report, error)
	ExpireOldPendingReports(ctx context.Context) error
	CalculateReportExpiryTime() time.Time
}
