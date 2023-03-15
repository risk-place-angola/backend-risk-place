package placetype_controllers

import (
	"net/http"

	"github.com/risk-place-angola/backend-risk-place/app/rest"
	placetype_presenter "github.com/risk-place-angola/backend-risk-place/app/rest/placetype/presenter"
	"github.com/risk-place-angola/backend-risk-place/usecase/placetype"
)

type PlaceTypeController interface {
	PlaceTypeCreateController(ctx placetype_presenter.PlaceTypePresenterCTX) error
	PlaceTypeFindAllController(ctx placetype_presenter.PlaceTypePresenterCTX) error
	PlaceTypeFindByIdController(ctx placetype_presenter.PlaceTypePresenterCTX) error
	PlaceTypeUpdateController(ctx placetype_presenter.PlaceTypePresenterCTX) error
	PlaceTypeDeleteController(ctx placetype_presenter.PlaceTypePresenterCTX) error
}

type PlaceTypeControllerImpl struct {
	placetypeUseCase placetype.PlaceTypeUseCase
}

func NewPlaceTypeController(placetypeUseCase placetype.PlaceTypeUseCase) PlaceTypeController {
	return &PlaceTypeControllerImpl{
		placetypeUseCase: placetypeUseCase,
	}
}

func (controller *PlaceTypeControllerImpl) PlaceTypeCreateController(ctx placetype_presenter.PlaceTypePresenterCTX) error {
	var placeType placetype.CreatePlaceTypeDTO
	if err := ctx.Bind(&placeType); err != nil {
		return ctx.JSON(http.StatusBadRequest, rest.ErrorResponse{Message: err.Error()})
	}

	placeTypeCreate, err := controller.placetypeUseCase.CreatePlaceType(placeType)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, rest.ErrorResponse{Message: err.Error()})
	}

	return ctx.JSON(http.StatusCreated, placeTypeCreate)
}

func (controller *PlaceTypeControllerImpl) PlaceTypeFindAllController(ctx placetype_presenter.PlaceTypePresenterCTX) error {
	placeTypeFindAll, err := controller.placetypeUseCase.FindAllPlaceTypes()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, rest.ErrorResponse{Message: err.Error()})
	}

	return ctx.JSON(http.StatusOK, placeTypeFindAll)
}

func (controller *PlaceTypeControllerImpl) PlaceTypeFindByIdController(ctx placetype_presenter.PlaceTypePresenterCTX) error {
	id := ctx.Param("id")

	placeTypeFindById, err := controller.placetypeUseCase.FindByIdPlaceType(id)
	if err != nil {
		return ctx.JSON(http.StatusFound, rest.ErrorResponse{Message: err.Error()})
	}

	return ctx.JSON(http.StatusOK, placeTypeFindById)
}

func (controller *PlaceTypeControllerImpl) PlaceTypeUpdateController(ctx placetype_presenter.PlaceTypePresenterCTX) error {
	id := ctx.Param("id")

	var placeType placetype.UpdatePlaceTypeDTO
	if err := ctx.Bind(&placeType); err != nil {
		return ctx.JSON(http.StatusBadRequest, rest.ErrorResponse{Message: err.Error()})
	}

	placeTypeUpdate, err := controller.placetypeUseCase.UpdatePlaceType(id, placeType)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, rest.ErrorResponse{Message: err.Error()})
	}

	return ctx.JSON(http.StatusOK, placeTypeUpdate)
}

func (controller *PlaceTypeControllerImpl) PlaceTypeDeleteController(ctx placetype_presenter.PlaceTypePresenterCTX) error {
	id := ctx.Param("id")

	if err := controller.placetypeUseCase.DeletePlaceType(id); err != nil {
		return ctx.JSON(http.StatusInternalServerError, rest.ErrorResponse{Message: err.Error()})
	}

	return ctx.JSON(http.StatusOK, rest.SuccessResponse{Message: "Place Type deleted successfully"})
}
