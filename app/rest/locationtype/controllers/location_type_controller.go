package locationtype_controllers

import (
	"net/http"

	"github.com/risk-place-angola/backend-risk-place/app/rest"
	locationtype_presenter "github.com/risk-place-angola/backend-risk-place/app/rest/locationtype/presenter"
	"github.com/risk-place-angola/backend-risk-place/usecase/locationtype"
)

type LocationTypeController interface {
	LocationTypeCreateController(ctx locationtype_presenter.LocationTypePresenterCTX) error
	LocationTypeFindAllController(ctx locationtype_presenter.LocationTypePresenterCTX) error
	LocationTypeFindByIdController(ctx locationtype_presenter.LocationTypePresenterCTX) error 
	LocationTypeUpdateController(ctx locationtype_presenter.LocationTypePresenterCTX) error
}

type LocationTypeControllerImpl struct {
	locationtypeUseCase locationtype.LocationTypeUseCase
}

func NewLocationTypeController(locationtypeUseCase locationtype.LocationTypeUseCase) LocationTypeController {
	return &LocationTypeControllerImpl{
		locationtypeUseCase: locationtypeUseCase,
	}
}

func (controller *LocationTypeControllerImpl) LocationTypeCreateController(ctx locationtype_presenter.LocationTypePresenterCTX) error {
	var locationType locationtype.CreateLocationTypeDTO
	if err := ctx.Bind(&locationType); err != nil {
		return ctx.JSON(http.StatusBadRequest, rest.ErrorResponse{Message: err.Error()})
	}

	locationTypeCreate, err := controller.locationtypeUseCase.CreateLocationType(locationType)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, rest.ErrorResponse{Message: err.Error()})
	}

	return ctx.JSON(http.StatusCreated, locationTypeCreate)
}

func (controller *LocationTypeControllerImpl) LocationTypeFindAllController(ctx locationtype_presenter.LocationTypePresenterCTX) error {
	locationTypeFindAll, err := controller.locationtypeUseCase.FindAllLocationTypes()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, rest.ErrorResponse{Message: err.Error()})
	}

	return ctx.JSON(http.StatusOK, locationTypeFindAll)
}

func (controller *LocationTypeControllerImpl) LocationTypeFindByIdController(ctx locationtype_presenter.LocationTypePresenterCTX) error {
	id := ctx.Param("id")

	locationTypeFindById, err := controller.locationtypeUseCase.FindByIdLocationType(id)
	if err != nil {
		return ctx.JSON(http.StatusFound, rest.ErrorResponse{Message: err.Error()})
	}

	return ctx.JSON(http.StatusOK, locationTypeFindById)
}

func (controller *LocationTypeControllerImpl) LocationTypeUpdateController(ctx locationtype_presenter.LocationTypePresenterCTX) error {
	id := ctx.Param("id")

	var locationType locationtype.UpdateLocationTypeDTO
	if err := ctx.Bind(&locationType); err != nil {
		return ctx.JSON(http.StatusBadRequest, rest.ErrorResponse{Message: err.Error()})
	}

	locationTypeUpdate, err := controller.locationtypeUseCase.UpdateLocationType(id, locationType)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, rest.ErrorResponse{Message: err.Error()})
	}

	return ctx.JSON(http.StatusOK, locationTypeUpdate)
}