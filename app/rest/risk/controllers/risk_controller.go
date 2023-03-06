package risk_controller

import (
	"net/http"

	risk_presenter "github.com/risk-place-angola/backend-risk-place/app/rest/risk/presenter"
	risk_usecase "github.com/risk-place-angola/backend-risk-place/usecase/risk"
)

type RiskController interface {
	RiskCreateController(ctx risk_presenter.RiskPresenterCTX) error
}

type RiskControllerImpl struct {
	riskUseCase risk_usecase.RiskUseCase
}

func NewRiskController(riskUseCase risk_usecase.RiskUseCase) RiskController {
	return &RiskControllerImpl{
		riskUseCase: riskUseCase,
	}
}

func (controller *RiskControllerImpl) RiskCreateController(ctx risk_presenter.RiskPresenterCTX) error {
	var risk risk_usecase.CreateRiskDTO
	if err := ctx.Bind(&risk); err != nil {
		return ctx.JSON(http.StatusBadRequest, risk_presenter.ErrorResponse{Message: err.Error()})
	}

	riskCreate, err := controller.riskUseCase.CreateRisk(risk)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, risk_presenter.ErrorResponse{Message: err.Error()})
	}

	return ctx.JSON(http.StatusCreated, riskCreate)
}
