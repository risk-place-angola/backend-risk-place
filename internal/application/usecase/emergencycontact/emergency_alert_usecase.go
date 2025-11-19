package emergencycontact

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/risk-place-angola/backend-risk-place/internal/application/dto"
	"github.com/risk-place-angola/backend-risk-place/internal/application/port"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/repository"
)

type EmergencyAlertUseCase struct {
	contactRepo repository.EmergencyContactRepository
	userRepo    repository.UserRepository
	smsNotifier port.NotifierSMSService
}

func NewEmergencyAlertUseCase(
	contactRepo repository.EmergencyContactRepository,
	userRepo repository.UserRepository,
	smsNotifier port.NotifierSMSService,
) *EmergencyAlertUseCase {
	return &EmergencyAlertUseCase{
		contactRepo: contactRepo,
		userRepo:    userRepo,
		smsNotifier: smsNotifier,
	}
}

func (uc *EmergencyAlertUseCase) SendEmergencyAlert(ctx context.Context, userID uuid.UUID, input dto.EmergencyAlertInput) (*dto.EmergencyAlertResponse, error) {
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		slog.Error("Error fetching user for emergency alert", "user_id", userID, "error", err)
		return nil, errors.New("failed to fetch user information")
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	contacts, err := uc.contactRepo.FindByUserID(ctx, userID)
	if err != nil {
		slog.Error("Error fetching emergency contacts", "user_id", userID, "error", err)
		return nil, errors.New("failed to fetch emergency contacts")
	}

	if len(contacts) == 0 {
		return nil, errors.New("no emergency contacts configured")
	}

	mapLink := fmt.Sprintf("https://maps.google.com/?q=%.6f,%.6f", input.Latitude, input.Longitude)
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	customMessage := input.Message
	if customMessage == "" {
		customMessage = "Preciso de ajuda urgente!"
	}

	smsMessage := fmt.Sprintf(
		"🚨 ALERTA DE EMERGÊNCIA 🚨\n\n"+
			"De: %s\n"+
			"Telefone: %s\n"+
			"Mensagem: %s\n"+
			"Data/Hora: %s\n"+
			"Localização: %s\n\n"+
			"Por favor, verifique se está tudo bem!",
		user.Name,
		user.Phone,
		customMessage,
		timestamp,
		mapLink,
	)

	notifiedContacts := make([]string, 0)
	successCount := 0

	for _, contact := range contacts {
		if contact.IsPriority {
			if err := uc.smsNotifier.NotifySMS(ctx, contact.Phone, smsMessage); err != nil {
				slog.Error("Failed to send SMS to priority contact",
					"contact_id", contact.ID,
					"contact_name", contact.Name,
					"phone", contact.Phone,
					"error", err)
				continue
			}

			notifiedContacts = append(notifiedContacts, contact.Name)
			successCount++

			slog.Info("Emergency alert sent successfully",
				"user_id", userID,
				"contact_name", contact.Name,
				"contact_phone", contact.Phone)
		}
	}

	if successCount == 0 {
		slog.Error("Failed to notify any emergency contacts", "user_id", userID)
		return nil, errors.New("failed to send emergency alerts to any contact")
	}

	return &dto.EmergencyAlertResponse{
		Success:          true,
		ContactsNotified: successCount,
		NotifiedContacts: notifiedContacts,
		Message:          fmt.Sprintf("Emergency alert sent to %d contact(s)", successCount),
	}, nil
}

func (uc *EmergencyAlertUseCase) SendEmergencyAlertToAll(ctx context.Context, userID uuid.UUID, input dto.EmergencyAlertInput) (*dto.EmergencyAlertResponse, error) {
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		slog.Error("Error fetching user for emergency alert", "user_id", userID, "error", err)
		return nil, errors.New("failed to fetch user information")
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	contacts, err := uc.contactRepo.FindByUserID(ctx, userID)
	if err != nil {
		slog.Error("Error fetching emergency contacts", "user_id", userID, "error", err)
		return nil, errors.New("failed to fetch emergency contacts")
	}

	if len(contacts) == 0 {
		return nil, errors.New("no emergency contacts configured")
	}

	mapLink := fmt.Sprintf("https://maps.google.com/?q=%.6f,%.6f", input.Latitude, input.Longitude)
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	customMessage := input.Message
	if customMessage == "" {
		customMessage = "Preciso de ajuda urgente!"
	}

	smsMessage := fmt.Sprintf(
		"🚨 ALERTA DE EMERGÊNCIA 🚨\n\n"+
			"De: %s\n"+
			"Telefone: %s\n"+
			"Mensagem: %s\n"+
			"Data/Hora: %s\n"+
			"Localização: %s\n\n"+
			"Por favor, verifique se está tudo bem!",
		user.Name,
		user.Phone,
		customMessage,
		timestamp,
		mapLink,
	)

	notifiedContacts := make([]string, 0)
	successCount := 0

	for _, contact := range contacts {
		if err := uc.smsNotifier.NotifySMS(ctx, contact.Phone, smsMessage); err != nil {
			slog.Error("Failed to send SMS to contact",
				"contact_id", contact.ID,
				"contact_name", contact.Name,
				"phone", contact.Phone,
				"error", err)
			continue
		}

		notifiedContacts = append(notifiedContacts, contact.Name)
		successCount++

		slog.Info("Emergency alert sent successfully",
			"user_id", userID,
			"contact_name", contact.Name,
			"contact_phone", contact.Phone)
	}

	if successCount == 0 {
		slog.Error("Failed to notify any emergency contacts", "user_id", userID)
		return nil, errors.New("failed to send emergency alerts to any contact")
	}

	return &dto.EmergencyAlertResponse{
		Success:          true,
		ContactsNotified: successCount,
		NotifiedContacts: notifiedContacts,
		Message:          fmt.Sprintf("Emergency alert sent to %d contact(s)", successCount),
	}, nil
}
