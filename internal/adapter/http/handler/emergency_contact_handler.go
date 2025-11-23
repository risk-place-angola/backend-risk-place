package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/risk-place-angola/backend-risk-place/internal/adapter/http/util"
	"github.com/risk-place-angola/backend-risk-place/internal/application"
	"github.com/risk-place-angola/backend-risk-place/internal/application/dto"
)

type EmergencyContactHandler struct {
	app *application.Application
}

func NewEmergencyContactHandler(app *application.Application) *EmergencyContactHandler {
	return &EmergencyContactHandler{
		app: app,
	}
}

// GetEmergencyContacts godoc
// @Summary Get all emergency contacts for the current user
// @Description Retrieve all emergency contacts configured by the authenticated user
// @Tags emergency-contacts
// @Security BearerAuth
// @Produce json
// @Success 200 {array} dto.EmergencyContactResponse
// @Failure 401 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /users/me/emergency-contacts [get]
func (h *EmergencyContactHandler) GetEmergencyContacts(w http.ResponseWriter, r *http.Request) {
	uid, ok := util.ExtractAndValidateUserID(w, r)
	if !ok {
		return
	}

	contacts, err := h.app.EmergencyContactUseCase.GetAll(r.Context(), uid)
	if err != nil {
		slog.Error("error fetching emergency contacts", "user_id", uid, "error", err)
		util.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	util.Response(w, contacts, http.StatusOK)
}

// CreateEmergencyContact godoc
// @Summary Create a new emergency contact
// @Description Add a new emergency contact for the authenticated user
// @Tags emergency-contacts
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param contact body dto.CreateEmergencyContactInput true "Emergency contact data"
// @Success 201 {object} dto.EmergencyContactResponse
// @Failure 400 {object} util.ErrorResponse
// @Failure 401 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /users/me/emergency-contacts [post]
func (h *EmergencyContactHandler) CreateEmergencyContact(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		util.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	uid, err := uuid.Parse(userID)
	if err != nil {
		util.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}

	var input dto.CreateEmergencyContactInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		util.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	contact, err := h.app.EmergencyContactUseCase.Create(r.Context(), uid, input)
	if err != nil {
		slog.Error("error creating emergency contact", "user_id", userID, "error", err)
		util.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	util.Response(w, contact, http.StatusCreated)
}

// UpdateEmergencyContact godoc
// @Summary Update an emergency contact
// @Description Update an existing emergency contact for the authenticated user
// @Tags emergency-contacts
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Emergency Contact ID"
// @Param contact body dto.UpdateEmergencyContactInput true "Updated emergency contact data"
// @Success 200 {object} dto.EmergencyContactResponse
// @Failure 400 {object} util.ErrorResponse
// @Failure 401 {object} util.ErrorResponse
// @Failure 404 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /users/me/emergency-contacts/{id} [put]
func (h *EmergencyContactHandler) UpdateEmergencyContact(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		util.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	uid, err := uuid.Parse(userID)
	if err != nil {
		util.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}

	contactID := r.PathValue("id")
	if contactID == "" {
		util.Error(w, "contact ID is required", http.StatusBadRequest)
		return
	}

	cid, err := uuid.Parse(contactID)
	if err != nil {
		util.Error(w, "invalid contact ID", http.StatusBadRequest)
		return
	}

	var input dto.UpdateEmergencyContactInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		util.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	contact, err := h.app.EmergencyContactUseCase.Update(r.Context(), uid, cid, input)
	if err != nil {
		slog.Error("error updating emergency contact", "user_id", userID, "contact_id", contactID, "error", err)
		if err.Error() == "emergency contact not found" {
			util.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		util.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	util.Response(w, contact, http.StatusOK)
}

// DeleteEmergencyContact godoc
// @Summary Delete an emergency contact
// @Description Delete an emergency contact for the authenticated user
// @Tags emergency-contacts
// @Security BearerAuth
// @Produce json
// @Param id path string true "Emergency Contact ID"
// @Success 204
// @Failure 400 {object} util.ErrorResponse
// @Failure 401 {object} util.ErrorResponse
// @Failure 404 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /users/me/emergency-contacts/{id} [delete]
func (h *EmergencyContactHandler) DeleteEmergencyContact(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		util.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	uid, err := uuid.Parse(userID)
	if err != nil {
		util.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}

	contactID := r.PathValue("id")
	if contactID == "" {
		util.Error(w, "contact ID is required", http.StatusBadRequest)
		return
	}

	cid, err := uuid.Parse(contactID)
	if err != nil {
		util.Error(w, "invalid contact ID", http.StatusBadRequest)
		return
	}

	if err := h.app.EmergencyContactUseCase.Delete(r.Context(), uid, cid); err != nil {
		slog.Error("error deleting emergency contact", "user_id", userID, "contact_id", contactID, "error", err)
		if err.Error() == "emergency contact not found" {
			util.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		util.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// SendEmergencyAlert godoc
// @Summary Send emergency alert to all priority contacts
// @Description Send an emergency SMS alert with location to all priority emergency contacts
// @Tags emergency
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param alert body dto.EmergencyAlertInput true "Emergency alert data with location"
// @Success 200 {object} dto.EmergencyAlertResponse
// @Failure 400 {object} util.ErrorResponse
// @Failure 401 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /emergency/alert [post]
func (h *EmergencyContactHandler) SendEmergencyAlert(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		util.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	uid, err := uuid.Parse(userID)
	if err != nil {
		util.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}

	var input dto.EmergencyAlertInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		util.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if input.Latitude == 0 || input.Longitude == 0 {
		util.Error(w, "latitude and longitude are required", http.StatusBadRequest)
		return
	}

	result, err := h.app.EmergencyAlertUseCase.SendEmergencyAlert(r.Context(), uid, input)
	if err != nil {
		slog.Error("error sending emergency alert", "user_id", userID, "error", err)
		util.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	util.Response(w, result, http.StatusOK)
}
