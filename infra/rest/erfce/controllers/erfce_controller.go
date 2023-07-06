package erfce_controller

import (
	"net/http"

	"github.com/risk-place-angola/backend-risk-place/infra/rest"
	erfce_presenter "github.com/risk-place-angola/backend-risk-place/infra/rest/erfce/presenter"

	account "github.com/risk-place-angola/backend-risk-place/usecase/erfce"
)

type ErfceController interface {
	ErfceCreateController(ctx erfce_presenter.ErfcePresenterCTX) error
	ErfceUpdateController(ctx erfce_presenter.ErfcePresenterCTX) error
	ErfceFindAllController(ctx erfce_presenter.ErfcePresenterCTX) error
	ErfceFindByIdController(ctx erfce_presenter.ErfcePresenterCTX) error
	ErfceDeleteController(ctx erfce_presenter.ErfcePresenterCTX) error
	ErfceLoginController(ctx erfce_presenter.ErfcePresenterCTX) error
	FindAllErfceWarningsController(ctx erfce_presenter.ErfcePresenterCTX) error
	FindWarningByErfceIDController(ctx erfce_presenter.ErfcePresenterCTX) error
}

type ErfceontrollerImpl struct {
	erfceUseCase account.ErfceUseCase
}

func NewErfceController(erfceRepo account.ErfceUseCase) ErfceController {
	return &ErfceontrollerImpl{
		erfceUseCase: erfceRepo,
	}
}

// @Summary Create Erfce
// @Description Create Erfce
// @Tags Erfce
// @Accept  json
// @Produce  json
// @Param erfce body account.CreateErfceDTO true "Erfce"
// @Success 201 {object} account.DTO
// @Failure 500 {object} rest.ErrorResponse
// @Router /api/v1/erfce [post]
func (controller *ErfceontrollerImpl) ErfceCreateController(ctx erfce_presenter.ErfcePresenterCTX) error {
	var erfce account.CreateErfceDTO
	if err := ctx.Bind(&erfce); err != nil {
		return ctx.JSON(http.StatusBadRequest, rest.ErrorResponse{Message: err.Error()})
	}

	erfceCreate, err := controller.erfceUseCase.CreateErfce(&erfce)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, rest.ErrorResponse{Message: err.Error()})
	}

	return ctx.JSON(http.StatusCreated, erfceCreate)
}

// @Summary Find All Erfce
// @Description Find All Erfce
// @Tags Erfce
// @Accept  json
// @Produce  json
// @Success 200 {object} []account.DTO
// @Failure 500 {object} rest.ErrorResponse
// @Router /api/v1/erfce [get]
func (controller *ErfceontrollerImpl) ErfceFindAllController(ctx erfce_presenter.ErfcePresenterCTX) error {
	erfces, err := controller.erfceUseCase.FindAllErfce()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, rest.ErrorResponse{Message: err.Error()})
	}
	return ctx.JSON(http.StatusOK, erfces)
}

// @Summary Find Erfce By ID
// @Description Find Erfce By ID
// @Tags Erfce
// @Accept  json
// @Produce  json
// @Param id path string true "Erfce ID"
// @Success 200 {object} account.DTO
// @Failure 500 {object} rest.ErrorResponse
// @Router /api/v1/erfce/{id} [get]
func (controller *ErfceontrollerImpl) ErfceFindByIdController(ctx erfce_presenter.ErfcePresenterCTX) error {
	id := ctx.Param("id")

	erfceId, err := controller.erfceUseCase.FindErfceByID(id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, rest.ErrorResponse{Message: err.Error()})
	}
	return ctx.JSON(http.StatusOK, erfceId)
}

// @Summary Update Erfce
// @Description Update Erfce
// @Tags Erfce
// @Accept  json
// @Produce  json
// @Param id path string true "Erfce ID"
// @Param erfce body account.UpdateErfceDTO true "Erfce"
// @Success 200 {object} account.DTO
// @Failure 500 {object} rest.ErrorResponse
// @Router /api/v1/erfce/{id} [put]
func (controller *ErfceontrollerImpl) ErfceUpdateController(ctx erfce_presenter.ErfcePresenterCTX) error {
	id := ctx.Param("id")

	var erfce account.UpdateErfceDTO
	if err := ctx.Bind(&erfce); err != nil {
		return ctx.JSON(http.StatusBadRequest, rest.ErrorResponse{Message: err.Error()})
	}

	erfceUpdate, err := controller.erfceUseCase.UpdateErfce(id, &erfce)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, rest.ErrorResponse{Message: err.Error()})
	}

	return ctx.JSON(http.StatusOK, erfceUpdate)
}

// @Summary Delete Erfce
// @Description Delete Erfce
// @Tags Erfce
// @Accept  json
// @Produce  json
// @Param id path string true "Erfce ID"
// @Success 200 {object} rest.SuccessResponse
// @Failure 500 {object} rest.ErrorResponse
// @Router /api/v1/erfce/{id} [delete]
func (controller *ErfceontrollerImpl) ErfceDeleteController(ctx erfce_presenter.ErfcePresenterCTX) error {
	id := ctx.Param("id")

	err := controller.erfceUseCase.RemoveErfce(id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, rest.ErrorResponse{Message: err.Error()})
	}
	return ctx.JSON(http.StatusOK, rest.SuccessResponse{Message: "Erfce deleted successfully"})
}

// @Summary Login Erfce
// @Description Login Erfce
// @Tags Erfce
// @Accept  json
// @Produce  json
// @Param erfce body account.LoginDTO true "Erfce"
// @Success 200 {object} account.DTO
// @Failure 500 {object} rest.ErrorResponse
// @Router /api/v1/erfce/login [post]
func (controller *ErfceontrollerImpl) ErfceLoginController(ctx erfce_presenter.ErfcePresenterCTX) error {
	var credentials account.LoginDTO
	if err := ctx.Bind(&credentials); err != nil {
		return ctx.JSON(http.StatusBadRequest, rest.ErrorResponse{Message: err.Error()})
	}

	data, err := controller.erfceUseCase.Login(&credentials)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, rest.ErrorResponse{Message: err.Error()})
	}

	return ctx.JSON(http.StatusOK, data)
}

// @Summary Find All Erfce Warnings
// @Description Find All Erfce Warnings
// @Tags Erfce
// @Accept  json
// @Produce  json
// @Success 200 {object} []account.DTO
// @Failure 500 {object} rest.ErrorResponse
// @Router /api/v1/erfce/warning [get]
func (controller *ErfceontrollerImpl) FindAllErfceWarningsController(ctx erfce_presenter.ErfcePresenterCTX) error {
	warnings, err := controller.erfceUseCase.FindAllUErfceWarnings()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, rest.ErrorResponse{Message: err.Error()})
	}
	return ctx.JSON(http.StatusOK, warnings)
}

// @Summary Find Erfce Warnings By ID
// @Description Find Erfce Warnings By ID
// @Tags Erfce
// @Accept  json
// @Produce  json
// @Param id path string true "Erfce ID"
// @Success 200 {object} []account.DTO
// @Failure 500 {object} rest.ErrorResponse
// @Router /api/v1/erfce/warning/{id} [get]
func (controller *ErfceontrollerImpl) FindWarningByErfceIDController(ctx erfce_presenter.ErfcePresenterCTX) error {
	id := ctx.Param("id")

	warnings, err := controller.erfceUseCase.FindWarningByErfceID(id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, rest.ErrorResponse{Message: err.Error()})
	}
	return ctx.JSON(http.StatusOK, warnings)
}
