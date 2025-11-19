package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type RelationType string

const (
	RelationFamily    RelationType = "family"
	RelationFriend    RelationType = "friend"
	RelationColleague RelationType = "colleague"
	RelationNeighbor  RelationType = "neighbor"
	RelationOther     RelationType = "other"
)

func isValidRelation(r RelationType) bool {
	validRelations := map[RelationType]bool{
		RelationFamily:    true,
		RelationFriend:    true,
		RelationColleague: true,
		RelationNeighbor:  true,
		RelationOther:     true,
	}
	return validRelations[r]
}

type EmergencyContact struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	Name       string
	Phone      string
	Relation   RelationType
	IsPriority bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func NewEmergencyContact(userID uuid.UUID, name, phone string, relation RelationType, isPriority bool) (*EmergencyContact, error) {
	if userID == uuid.Nil {
		return nil, errors.New("user ID is required")
	}

	if name == "" {
		return nil, errors.New("name is required")
	}

	if phone == "" {
		return nil, errors.New("phone is required")
	}

	if !isValidRelation(relation) {
		return nil, errors.New("invalid relation type")
	}

	now := time.Now()
	return &EmergencyContact{
		ID:         uuid.New(),
		UserID:     userID,
		Name:       name,
		Phone:      phone,
		Relation:   relation,
		IsPriority: isPriority,
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil
}

func (ec *EmergencyContact) Update(name, phone string, relation RelationType, isPriority bool) error {
	if name == "" {
		return errors.New("name is required")
	}

	if phone == "" {
		return errors.New("phone is required")
	}

	if !isValidRelation(relation) {
		return errors.New("invalid relation type")
	}

	ec.Name = name
	ec.Phone = phone
	ec.Relation = relation
	ec.IsPriority = isPriority
	ec.UpdatedAt = time.Now()

	return nil
}

func IsValidRelation(relation string) bool {
	return isValidRelation(RelationType(relation))
}
