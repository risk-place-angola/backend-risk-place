package placetype_controllers

import (
	"net/http"

	"github.com/risk-place-angola/backend-risk-place/api/rest"
	placetype_presenter "github.com/risk-place-angola/backend-risk-place/api/rest/placetype/presenter"
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

// @Summary Create PlaceType
// @Description Create PlaceType
// @Tags PlaceType
// @Accept  json
// @Produce  json
// @Param placeType body placetype.CreatePlaceTypeDTO true "PlaceType"
// @Success 201 {object} placetype.PlaceTypeDTO
// @Failure 500 {object} rest.ErrorResponse
// @Router /api/v1/placetype [post]
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

// @Summary Find All PlaceType
// @Description Find All PlaceType
// @Tags PlaceType
// @Accept  json
// @Produce  json
// @Success 200 {object} []placetype.PlaceTypeDTO
// @Failure 500 {object} rest.ErrorResponse
// @Router /api/v1/placetype [get]
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

// @Summary Update PlaceType
// @Description Update PlaceType
// @Tags PlaceType
// @Accept  json
// @Produce  json
// @Param id path string true "PlaceType ID"
// @Param placeType body placetype.UpdatePlaceTypeDTO true "PlaceType"
// @Success 200 {object} placetype.PlaceTypeDTO
// @Failure 500 {object} rest.ErrorResponse
// @Router /api/v1/placetype/{id} [put]
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

// @Summary Delete PlaceType
// @Description Delete PlaceType
// @Tags PlaceType
// @Accept  json
// @Produce  json
// @Param id path string true "PlaceType ID"
// @Success 200 {object} rest.SuccessResponse
// @Failure 500 {object} rest.ErrorResponse
// @Router /api/v1/placetype/{id} [delete]
func (controller *PlaceTypeControllerImpl) PlaceTypeDeleteController(ctx placetype_presenter.PlaceTypePresenterCTX) error {
	id := ctx.Param("id")

	if err := controller.placetypeUseCase.DeletePlaceType(id); err != nil {
		return ctx.JSON(http.StatusInternalServerError, rest.ErrorResponse{Message: err.Error()})
	}

	return ctx.JSON(http.StatusOK, rest.SuccessResponse{Message: "Place Type deleted successfully"})
}
