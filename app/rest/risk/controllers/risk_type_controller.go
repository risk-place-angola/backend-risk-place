package risk_controller

import (
	"net/http"

	"github.com/risk-place-angola/backend-risk-place/app/rest"
	risk_presenter "github.com/risk-place-angola/backend-risk-place/app/rest/risk/presenter"
	"github.com/risk-place-angola/backend-risk-place/usecase/risk/risktype"
)

type RiskTypeController interface {
	RiskTypeCreateController(ctx risk_presenter.RiskPresenterCTX) error
	RiskTypeUpdateController(ctx risk_presenter.RiskPresenterCTX) error
}

type RiskTypeControllerImpl struct {
	riskTypeUseCase risktype.RiskTypeUseCase
}

func NewRiskTypeController(riskTypeUseCase risktype.RiskTypeUseCase) RiskTypeController {
	return &RiskTypeControllerImpl{
		riskTypeUseCase: riskTypeUseCase,
	}
}

func (controller *RiskTypeControllerImpl) RiskTypeCreateController(ctx risk_presenter.RiskPresenterCTX) error {
	var risktype risktype.CreateRiskTypeDTO
	if err := ctx.Bind(&risktype); err != nil {
		return ctx.JSON(http.StatusBadRequest, rest.ErrorResponse{Message: err.Error()})
	}

	risktypeCreate, err := controller.riskTypeUseCase.CreateRiskType(&risktype)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, rest.ErrorResponse{Message: err.Error()})
	}

	return ctx.JSON(http.StatusCreated, risktypeCreate)
}

func (controller *RiskTypeControllerImpl) RiskTypeUpdateController(ctx risk_presenter.RiskPresenterCTX) error {
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
