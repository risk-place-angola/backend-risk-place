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

func (controller *RiskTypeControllerImpl) RiskTypeFindAllController(ctx risk_type_presenter.RiskTypePresenterCTX) error {
	riskTypes, err := controller.riskTypeUseCase.FindAllRiskTypes()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, rest.ErrorResponse{Message: err.Error()})
	}
	return ctx.JSON(http.StatusOK, riskTypes)
}

func (controller *RiskTypeControllerImpl) RiskTypeFindByIdController(ctx risk_type_presenter.RiskTypePresenterCTX) error {
	id := ctx.Param("id")

	riskTypeID, err := controller.riskTypeUseCase.FindRiskTypeByID(id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, rest.ErrorResponse{Message: err.Error()})
	}
	return ctx.JSON(http.StatusOK, riskTypeID)
}

func (controller *RiskTypeControllerImpl) RiskTypeDeleteController(ctx risk_type_presenter.RiskTypePresenterCTX) error {
	id := ctx.Param("id")

	err := controller.riskTypeUseCase.RemoveRiskType(id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, rest.ErrorResponse{Message: err.Error()})
	}
	return ctx.JSON(http.StatusOK, rest.SuccessResponse{Message: "RiskType deleted successfully"})
}
