package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	httputil "github.com/risk-place-angola/backend-risk-place/internal/adapter/http/util"
	"github.com/risk-place-angola/backend-risk-place/internal/application/dto"
	"github.com/risk-place-angola/backend-risk-place/internal/application/usecase/device"
)

type DeviceHandler struct {
	registerDeviceUC       *device.RegisterDeviceUseCase
	updateDeviceLocationUC *device.UpdateDeviceLocationUseCase
}

func NewDeviceHandler(
	registerDeviceUC *device.RegisterDeviceUseCase,
	updateDeviceLocationUC *device.UpdateDeviceLocationUseCase,
) *DeviceHandler {
	return &DeviceHandler{
		registerDeviceUC:       registerDeviceUC,
		updateDeviceLocationUC: updateDeviceLocationUC,
	}
}

// RegisterDevice godoc
// @Summary Register anonymous device
// @Description Register or update an anonymous device for receiving notifications without authentication
// @Tags devices
// @Accept json
// @Produce json
// @Param request body dto.RegisterDeviceRequest true "Device registration data"
// @Success 200 {object} dto.DeviceResponse
// @Failure 400 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /devices/register [post]
func (h *DeviceHandler) RegisterDevice(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterDeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("failed to decode request", slog.Any("error", err))
		httputil.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := h.registerDeviceUC.Execute(r.Context(), req)
	if err != nil {
		slog.Error("failed to register device", slog.Any("error", err))
		httputil.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	httputil.Response(w, resp, http.StatusOK)
}

// UpdateDeviceLocation godoc
// @Summary Update device location
// @Description Update the location of an anonymous device for proximity-based notifications
// @Tags devices
// @Accept json
// @Produce json
// @Param request body dto.UpdateDeviceLocationRequest true "Location update data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /devices/location [put]
func (h *DeviceHandler) UpdateDeviceLocation(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateDeviceLocationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("failed to decode request", slog.Any("error", err))
		httputil.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := h.updateDeviceLocationUC.Execute(r.Context(), req)
	if err != nil {
		slog.Error("failed to update device location", slog.Any("error", err))
		httputil.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	httputil.Response(w, map[string]string{
		"message": "Location updated successfully",
	}, http.StatusOK)
}
