package emergencycontact

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/risk-place-angola/backend-risk-place/internal/application/dto"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/model"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/repository"
)

const MaxPriorityContacts = 5

type EmergencyContactUseCase struct {
	repo repository.EmergencyContactRepository
}

func NewEmergencyContactUseCase(repo repository.EmergencyContactRepository) *EmergencyContactUseCase {
	return &EmergencyContactUseCase{
		repo: repo,
	}
}

func (uc *EmergencyContactUseCase) Create(ctx context.Context, userID uuid.UUID, input dto.CreateEmergencyContactInput) (*dto.EmergencyContactResponse, error) {
	if input.IsPriority {
		count, err := uc.repo.CountPriorityByUserID(ctx, userID)
		if err != nil {
			slog.Error("Error counting priority contacts", "user_id", userID, "error", err)
			return nil, errors.New("failed to validate priority contacts limit")
		}

		if count >= MaxPriorityContacts {
			return nil, errors.New("maximum priority contacts limit reached (5)")
		}
	}

	contact, err := model.NewEmergencyContact(
		userID,
		input.Name,
		input.Phone,
		model.RelationType(input.Relation),
		input.IsPriority,
	)
	if err != nil {
		slog.Error("Error creating emergency contact model", "user_id", userID, "error", err)
		return nil, err
	}

	if err := uc.repo.Save(ctx, contact); err != nil {
		slog.Error("Error saving emergency contact", "user_id", userID, "error", err)
		return nil, errors.New("failed to create emergency contact")
	}

	return &dto.EmergencyContactResponse{
		ID:         contact.ID.String(),
		UserID:     contact.UserID.String(),
		Name:       contact.Name,
		Phone:      contact.Phone,
		Relation:   string(contact.Relation),
		IsPriority: contact.IsPriority,
		CreatedAt:  contact.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:  contact.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

func (uc *EmergencyContactUseCase) GetAll(ctx context.Context, userID uuid.UUID) ([]dto.EmergencyContactResponse, error) {
	contacts, err := uc.repo.FindByUserID(ctx, userID)
	if err != nil {
		slog.Error("Error fetching emergency contacts", "user_id", userID, "error", err)
		return nil, errors.New("failed to fetch emergency contacts")
	}

	response := make([]dto.EmergencyContactResponse, 0, len(contacts))
	for _, contact := range contacts {
		response = append(response, dto.EmergencyContactResponse{
			ID:         contact.ID.String(),
			UserID:     contact.UserID.String(),
			Name:       contact.Name,
			Phone:      contact.Phone,
			Relation:   string(contact.Relation),
			IsPriority: contact.IsPriority,
			CreatedAt:  contact.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:  contact.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return response, nil
}

func (uc *EmergencyContactUseCase) GetByID(ctx context.Context, userID, contactID uuid.UUID) (*dto.EmergencyContactResponse, error) {
	contact, err := uc.repo.FindByUserIDAndID(ctx, userID, contactID)
	if err != nil {
		slog.Error("Error fetching emergency contact", "user_id", userID, "contact_id", contactID, "error", err)
		return nil, errors.New("failed to fetch emergency contact")
	}

	if contact == nil {
		return nil, errors.New("emergency contact not found")
	}

	return &dto.EmergencyContactResponse{
		ID:         contact.ID.String(),
		UserID:     contact.UserID.String(),
		Name:       contact.Name,
		Phone:      contact.Phone,
		Relation:   string(contact.Relation),
		IsPriority: contact.IsPriority,
		CreatedAt:  contact.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:  contact.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

func (uc *EmergencyContactUseCase) Update(ctx context.Context, userID, contactID uuid.UUID, input dto.UpdateEmergencyContactInput) (*dto.EmergencyContactResponse, error) {
	contact, err := uc.repo.FindByUserIDAndID(ctx, userID, contactID)
	if err != nil {
		slog.Error("Error fetching emergency contact for update", "user_id", userID, "contact_id", contactID, "error", err)
		return nil, errors.New("failed to fetch emergency contact")
	}

	if contact == nil {
		return nil, errors.New("emergency contact not found")
	}

	if input.IsPriority && !contact.IsPriority {
		count, err := uc.repo.CountPriorityByUserID(ctx, userID)
		if err != nil {
			slog.Error("Error counting priority contacts", "user_id", userID, "error", err)
			return nil, errors.New("failed to validate priority contacts limit")
		}

		if count >= MaxPriorityContacts {
			return nil, errors.New("maximum priority contacts limit reached (5)")
		}
	}

	if err := contact.Update(input.Name, input.Phone, model.RelationType(input.Relation), input.IsPriority); err != nil {
		slog.Error("Error updating emergency contact model", "contact_id", contactID, "error", err)
		return nil, err
	}

	if err := uc.repo.Update(ctx, contact); err != nil {
		slog.Error("Error saving updated emergency contact", "contact_id", contactID, "error", err)
		return nil, errors.New("failed to update emergency contact")
	}

	return &dto.EmergencyContactResponse{
		ID:         contact.ID.String(),
		UserID:     contact.UserID.String(),
		Name:       contact.Name,
		Phone:      contact.Phone,
		Relation:   string(contact.Relation),
		IsPriority: contact.IsPriority,
		CreatedAt:  contact.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:  contact.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

func (uc *EmergencyContactUseCase) Delete(ctx context.Context, userID, contactID uuid.UUID) error {
	contact, err := uc.repo.FindByUserIDAndID(ctx, userID, contactID)
	if err != nil {
		slog.Error("Error fetching emergency contact for deletion", "user_id", userID, "contact_id", contactID, "error", err)
		return errors.New("failed to fetch emergency contact")
	}

	if contact == nil {
		return errors.New("emergency contact not found")
	}

	if err := uc.repo.DeleteByUserIDAndID(ctx, userID, contactID); err != nil {
		slog.Error("Error deleting emergency contact", "contact_id", contactID, "error", err)
		return errors.New("failed to delete emergency contact")
	}

	return nil
}
