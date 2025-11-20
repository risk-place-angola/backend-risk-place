package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/model"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/repository"
)

const (
	trustScoreMin            = 0
	trustScoreMax            = 100
	trustScoreDefault        = 50
	trustScorePerUpvote      = 2
	trustScorePerDownvote    = -3
	trustScorePerVerified    = 5
	autoVerifyThreshold      = 3
	duplicateRadiusMeters    = 50.0
	duplicateTimeWindowHours = 24
	reportExpiryHours        = 48
)

type ReportVerificationService struct {
	reportRepo repository.ReportRepository
}

func NewReportVerificationService(reportRepo repository.ReportRepository) *ReportVerificationService {
	return &ReportVerificationService{
		reportRepo: reportRepo,
	}
}

func (s *ReportVerificationService) VoteReport(ctx context.Context, reportID uuid.UUID, userID *uuid.UUID, anonymousSessionID *uuid.UUID, voteType model.VoteType) error {
	vote := &model.ReportVote{
		ID:                 uuid.New(),
		ReportID:           reportID,
		UserID:             userID,
		AnonymousSessionID: anonymousSessionID,
		VoteType:           voteType,
		CreatedAt:          time.Now(),
	}

	if err := s.reportRepo.AddVote(ctx, vote); err != nil {
		slog.Error("failed to add vote", "error", err)
		return err
	}

	if err := s.recalculateVerification(ctx, reportID); err != nil {
		slog.Error("failed to recalculate verification", "error", err)
		return err
	}

	if userID != nil {
		if err := s.updateVoterTrustScore(ctx, *userID, voteType); err != nil {
			slog.Warn("failed to update voter trust score", "error", err)
		}
	}

	return nil
}

func (s *ReportVerificationService) recalculateVerification(ctx context.Context, reportID uuid.UUID) error {
	report, err := s.reportRepo.GetByID(ctx, reportID)
	if err != nil {
		return err
	}

	upvotes := report.VerificationCount
	downvotes := report.RejectionCount

	if err := s.reportRepo.UpdateVerificationCounts(ctx, reportID, upvotes, downvotes); err != nil {
		return err
	}

	if upvotes >= autoVerifyThreshold && report.Status == model.ReportStatusPending {
		if err := s.reportRepo.VerifyReport(ctx, reportID, uuid.Nil); err != nil {
			slog.Error("failed to auto-verify report", "error", err)
			return err
		}

		if err := s.reportRepo.IncrementReportsVerified(ctx, report.UserID); err != nil {
			slog.Warn("failed to increment verified count", "error", err)
		}

		if err := s.updateCreatorTrustScore(ctx, report.UserID, true); err != nil {
			slog.Warn("failed to update creator trust score", "error", err)
		}

		slog.Info("report auto-verified", "reportID", reportID, "upvotes", upvotes)
	}

	return nil
}

func (s *ReportVerificationService) updateVoterTrustScore(ctx context.Context, userID uuid.UUID, voteType model.VoteType) error {
	delta := trustScorePerUpvote
	if voteType == model.VoteTypeDownvote {
		delta = trustScorePerDownvote
	}

	return s.adjustTrustScore(ctx, userID, delta)
}

func (s *ReportVerificationService) updateCreatorTrustScore(ctx context.Context, userID uuid.UUID, verified bool) error {
	delta := trustScorePerVerified
	if !verified {
		delta = -trustScorePerVerified
	}

	return s.adjustTrustScore(ctx, userID, delta)
}

func (s *ReportVerificationService) adjustTrustScore(ctx context.Context, userID uuid.UUID, delta int) error {
	user, err := s.reportRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	newScore := clamp(user.VerificationCount+delta, trustScoreMin, trustScoreMax)
	return s.reportRepo.UpdateTrustScore(ctx, userID, newScore)
}

func (s *ReportVerificationService) CheckDuplicates(ctx context.Context, lat, lon float64, riskTypeID uuid.UUID) ([]*model.Report, error) {
	since := time.Now().Add(-duplicateTimeWindowHours * time.Hour)
	return s.reportRepo.FindDuplicates(ctx, lat, lon, riskTypeID, duplicateRadiusMeters, since)
}

func (s *ReportVerificationService) ExpireOldPendingReports(ctx context.Context) error {
	before := time.Now().Add(-reportExpiryHours * time.Hour)
	return s.reportRepo.ExpireOldReports(ctx, before)
}

func (s *ReportVerificationService) CalculateReportExpiryTime() time.Time {
	return time.Now().Add(reportExpiryHours * time.Hour)
}

func clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
