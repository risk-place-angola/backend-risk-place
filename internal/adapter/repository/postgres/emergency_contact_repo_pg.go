package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/risk-place-angola/backend-risk-place/internal/adapter/repository/postgres/sqlc"
	domainErrors "github.com/risk-place-angola/backend-risk-place/internal/domain/errors"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/model"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/repository"
)

type emergencyContactRepoPG struct {
	q sqlc.Querier
}

func NewEmergencyContactRepository(db *sql.DB) repository.EmergencyContactRepository {
	return &emergencyContactRepoPG{q: sqlc.New(db)}
}

func (r *emergencyContactRepoPG) Save(ctx context.Context, contact *model.EmergencyContact) error {
	return r.q.CreateEmergencyContact(ctx, sqlc.CreateEmergencyContactParams{
		ID:         contact.ID,
		UserID:     contact.UserID,
		Name:       contact.Name,
		Phone:      contact.Phone,
		Relation:   string(contact.Relation),
		IsPriority: contact.IsPriority,
		CreatedAt:  contact.CreatedAt,
		UpdatedAt:  contact.UpdatedAt,
	})
}

func (r *emergencyContactRepoPG) Update(ctx context.Context, contact *model.EmergencyContact) error {
	return r.q.UpdateEmergencyContact(ctx, sqlc.UpdateEmergencyContactParams{
		ID:         contact.ID,
		Name:       contact.Name,
		Phone:      contact.Phone,
		Relation:   string(contact.Relation),
		IsPriority: contact.IsPriority,
		UpdatedAt:  time.Now(),
	})
}

func (r *emergencyContactRepoPG) FindByID(ctx context.Context, id uuid.UUID) (*model.EmergencyContact, error) {
	row, err := r.q.GetEmergencyContactByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domainErrors.ErrNotFound
		}
		return nil, err
	}

	return r.toDomain(row), nil
}

func (r *emergencyContactRepoPG) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*model.EmergencyContact, error) {
	rows, err := r.q.GetEmergencyContactsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	contacts := make([]*model.EmergencyContact, 0, len(rows))
	for _, row := range rows {
		contacts = append(contacts, r.toDomain(row))
	}

	return contacts, nil
}

func (r *emergencyContactRepoPG) FindByUserIDAndID(ctx context.Context, userID, contactID uuid.UUID) (*model.EmergencyContact, error) {
	row, err := r.q.GetEmergencyContactByUserIDAndID(ctx, sqlc.GetEmergencyContactByUserIDAndIDParams{
		UserID: userID,
		ID:     contactID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domainErrors.ErrNotFound
		}
		return nil, err
	}

	return r.toDomain(row), nil
}

func (r *emergencyContactRepoPG) FindPriorityByUserID(ctx context.Context, userID uuid.UUID) ([]*model.EmergencyContact, error) {
	rows, err := r.q.GetPriorityEmergencyContactsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	contacts := make([]*model.EmergencyContact, 0, len(rows))
	for _, row := range rows {
		contacts = append(contacts, r.toDomain(row))
	}

	return contacts, nil
}

func (r *emergencyContactRepoPG) CountPriorityByUserID(ctx context.Context, userID uuid.UUID) (int, error) {
	count, err := r.q.CountPriorityEmergencyContactsByUserID(ctx, userID)
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (r *emergencyContactRepoPG) Delete(ctx context.Context, id string) error {
	contactID, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return r.q.DeleteEmergencyContact(ctx, contactID)
}

func (r *emergencyContactRepoPG) DeleteByUserIDAndID(ctx context.Context, userID, contactID uuid.UUID) error {
	return r.q.DeleteEmergencyContactByUserIDAndID(ctx, sqlc.DeleteEmergencyContactByUserIDAndIDParams{
		UserID: userID,
		ID:     contactID,
	})
}

func (r *emergencyContactRepoPG) FindAll(ctx context.Context) ([]*model.EmergencyContact, error) {
	return nil, errors.New("FindAll not implemented for emergency contacts")
}

func (r *emergencyContactRepoPG) toDomain(row sqlc.EmergencyContact) *model.EmergencyContact {
	return &model.EmergencyContact{
		ID:         row.ID,
		UserID:     row.UserID,
		Name:       row.Name,
		Phone:      row.Phone,
		Relation:   model.RelationType(row.Relation),
		IsPriority: row.IsPriority,
		CreatedAt:  row.CreatedAt,
		UpdatedAt:  row.UpdatedAt,
	}
}
