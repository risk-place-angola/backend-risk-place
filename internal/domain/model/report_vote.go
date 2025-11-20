package model

import (
	"time"

	"github.com/google/uuid"
)

type VoteType string

const (
	VoteTypeUpvote   VoteType = "upvote"
	VoteTypeDownvote VoteType = "downvote"
)

type ReportVote struct {
	ID                 uuid.UUID
	ReportID           uuid.UUID
	UserID             *uuid.UUID
	AnonymousSessionID *uuid.UUID
	VoteType           VoteType
	CreatedAt          time.Time
}
