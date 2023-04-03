package risktype_controller

import (
	"net/http"

	"github.com/risk-place-angola/backend-risk-place/api/rest"
	risk_type_presenter "github.com/risk-place-angola/backend-risk-place/api/rest/risktype/presenter"
	"github.com/risk-place-angola/backend-risk-place/usecase/risktype"
)

type RiskTypeController interface {
	RiskTypeCreateController(ctx risk_type_presenter.RiskTypePresenterCTX) error
	RiskTypeUpdateController(ctx risk_type_presenter.RiskTypePresenterCTX) error
	RiskTypeFindAllController(ctx risk_type_presenter.RiskTypePresenterCTX) error
	RiskTypeFindByIdController(ctx risk_type_presenter.RiskTypePresenterCTX) error
	RiskTypeDeleteController(ctx risk_type_presenter.RiskTypePresenterCTX) error
}

type RiskTypeControllerImpl struct {
	riskTypeUseCase risktype.RiskTypeUseCase
}

func NewRiskTypeController(riskTypeUseCase risktype.RiskTypeUseCase) RiskTypeController {
	return &RiskTypeControllerImpl{
		riskTypeUseCase: riskTypeUseCase,
	}
}

// @Summary Create RiskType
// @Description Create RiskType
// @Tags RiskType
// @Accept  json
// @Produce  json
// @Param risktype body risktype.CreateRiskTypeDTO true "RiskType"
// @Success 201 {object} risktype.RiskTypeDTO
// @Failure 500 {object} rest.ErrorResponse
// @Router /api/v1/risktype [post]
func (controller *RiskTypeControllerImpl) RiskTypeCreateController(ctx risk_type_presenter.RiskTypePresenterCTX) error {
	var risktype risktype.CreateRiskTypeDTO
	if err := ctx.Bind(&risktype); err != nil {
		return ctx.JSON(http.StatusBadRequest, rest.ErrorResponse{Message: err.Error()})
	}

	risktypeCreate, err := controller.riskTypeUseCase.CreateRiskType(&risktype)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, rest.ErrorResponse{Message: err.Error()})
	}

	return ctx.JSON(http.StatusCreated, risktypeCreate)
}

// @Summary Update RiskType
// @Description Update RiskType
// @Tags RiskType
// @Accept  json
// @Produce  json
// @Param id path string true "RiskType ID"
// @Param risktype body risktype.UpdateRiskTypeDTO true "RiskType"
// @Success 200 {object} risktype.RiskTypeDTO
// @Failure 500 {object} rest.ErrorResponse
// @Router /api/v1/risktype/{id} [put]
func (controller *RiskTypeControllerImpl) RiskTypeUpdateController(ctx risk_type_presenter.RiskTypePresenterCTX) error {
	id := ctx.Param("id")

	var risktype risktype.UpdateRiskTypeDTO
	if err := ctx.Bind(&risktype); err != nil {
		return ctx.JSON(http.StatusBadRequest, rest.ErrorResponse{Message: err.Error()})
	}

	risktypeUpdate, err := controller.riskTypeUseCase.UpdateRiskType(id, &risktype)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, rest.ErrorResponse{Message: err.Error()})
	}

	return ctx.JSON(http.StatusOK, risktypeUpdate)
}

// @Summary Find All RiskType
// @Description Find All RiskType
// @Tags RiskType
// @Accept  json
// @Produce  json
// @Success 200 {object} []risktype.RiskTypeDTO
// @Failure 500 {object} rest.ErrorResponse
// @Router /api/v1/risktype [get]
func (controller *RiskTypeControllerImpl) RiskTypeFindAllController(ctx risk_type_presenter.RiskTypePresenterCTX) error {
	riskTypes, err := controller.riskTypeUseCase.FindAllRiskTypes()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, rest.ErrorResponse{Message: err.Error()})
	}
	return ctx.JSON(http.StatusOK, riskTypes)
}

// @Summary Find RiskType By ID
// @Description Find RiskType By ID
// @Tags RiskType
// @Accept  json
// @Produce  json
// @Param id path string true "RiskType ID"
// @Success 200 {object} risktype.RiskTypeDTO
// @Failure 500 {object} rest.ErrorResponse
// @Router /api/v1/risktype/{id} [get]
func (controller *RiskTypeControllerImpl) RiskTypeFindByIdController(ctx risk_type_presenter.RiskTypePresenterCTX) error {
	id := ctx.Param("id")

	riskTypeID, err := controller.riskTypeUseCase.FindRiskTypeByID(id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, rest.ErrorResponse{Message: err.Error()})
	}
	return ctx.JSON(http.StatusOK, riskTypeID)
}

// @Summary Delete RiskType
// @Description Delete RiskType
// @Tags RiskType
// @Accept  json
// @Produce  json
// @Param id path string true "RiskType ID"
// @Success 200 {object} rest.SuccessResponse
// @Failure 500 {object} rest.ErrorResponse
// @Router /api/v1/risktype/{id} [delete]
func (controller *RiskTypeControllerImpl) RiskTypeDeleteController(ctx risk_type_presenter.RiskTypePresenterCTX) error {
	id := ctx.Param("id")

	err := controller.riskTypeUseCase.RemoveRiskType(id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, rest.ErrorResponse{Message: err.Error()})
	}
	return ctx.JSON(http.StatusOK, rest.SuccessResponse{Message: "RiskType deleted successfully"})
}
